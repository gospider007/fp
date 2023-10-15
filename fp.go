package fp

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gospider007/gtls"
	"github.com/gospider007/ja3"
	"github.com/gospider007/net/http2"
	"github.com/gospider007/tools"
)

type Option struct {
	Addr        string
	Server      *http.Server
	CertFile    string
	KeyFile     string
	Certificate tls.Certificate
	TLSConfig   *tls.Config
	Handler     http.Handler
}

// https://tools.scrapfly.io/api/fp/anything?extended=1
func GinHandlerFunc(ctx *gin.Context) {
	fpData, ok := ja3.GetFpContextData(ctx.Request.Context())
	result := make(map[string]any)
	result["negotiatedProtocol"] = ctx.Request.TLS.NegotiatedProtocol
	result["tlsVersion"] = ctx.Request.TLS.Version
	result["userAgent"] = ctx.Request.UserAgent()
	rawClientHelloInfo, err := fpData.RawClientHelloInfo()
	if err == nil {
		clientHelloParseData := rawClientHelloInfo.Parse()
		result["tls"] = clientHelloParseData
		result["ja3"], result["ja3n"] = clientHelloParseData.Fp()
		ja4aStr := "t"
		switch ctx.Request.TLS.Version {
		case tls.VersionTLS10:
			ja4aStr += "10"
		case tls.VersionTLS11:
			ja4aStr += "11"
		case tls.VersionTLS12:
			ja4aStr += "12"
		case tls.VersionTLS13:
			ja4aStr += "13"
		default:
			ja4aStr += "00"
		}
		if ctx.Request.TLS.ServerName == "" {
			ja4aStr += "i"
		} else if _, addTyp := gtls.ParseHost(ctx.Request.TLS.ServerName); addTyp != 0 {
			ja4aStr += "i"
		} else {
			ja4aStr += "d"
		}
		ciphers := ja3.ClearGreas(clientHelloParseData.Ciphers)
		ja4aStr += fmt.Sprint(len(ciphers))
		exts := ja3.ClearGreas(clientHelloParseData.Extensions)
		ja4aStr += fmt.Sprint(len(exts))
		if ctx.Request.TLS.NegotiatedProtocol == "" {
			ja4aStr += "00"
		} else {
			ja4aStr += ctx.Request.TLS.NegotiatedProtocol
		}
		sort.Slice(ciphers, func(i, j int) bool { return ciphers[i] < ciphers[j] })
		sort.Slice(exts, func(i, j int) bool { return exts[i] < exts[j] })
		ja4bStr := tools.Hex(sha256.Sum256([]byte(tools.AnyJoin(ciphers, ""))))
		ja4cStr := tools.Hex(sha256.Sum256([]byte(tools.AnyJoin(exts, "") + tools.AnyJoin(clientHelloParseData.Algorithms, ""))))
		ja4 := tools.AnyJoin([]string{ja4aStr, ja4bStr, ja4cStr}, "_")
		result["ja4"] = ja4
	}
	h2Ja3Spec := fpData.H2Ja3Spec()
	result["http2"] = h2Ja3Spec
	result["akamai_fp"] = h2Ja3Spec.Fp()
	if ok {
		ctx.JSON(200, result)
	} else {
		ctx.JSON(200, map[string]any{
			"error": "指纹加载失败",
		})
	}
}

func Serve(handler http.Handler, options ...Option) (err error) {
	var option Option
	if len(options) > 0 {
		option = options[0]
	}
	if option.Addr == "" {
		if option.Server != nil && option.Server.Addr != "" {
			option.Addr = option.Server.Addr
		} else {
			option.Addr = ":0"
		}
	}
	var certificate tls.Certificate
	if option.Certificate.Certificate != nil {
		certificate = option.Certificate
	} else if option.CertFile != "" && option.KeyFile != "" {
		if certData, err := os.ReadFile(option.CertFile); err != nil {
			return err
		} else if cert, err := gtls.LoadCert(certData); err != nil {
			return err
		} else if keyData, err := os.ReadFile(option.KeyFile); err != nil {
			return err
		} else if key, err := gtls.LoadCertKey(keyData); err != nil {
			return err
		} else if certificate, err = gtls.CreateTlsCert(cert, key); err != nil {
			return err
		}
	} else if certificate, err = gtls.CreateProxyCertWithName("test"); err != nil {
		return err
	}
	var tLSConfig *tls.Config
	if option.TLSConfig == nil {
		tLSConfig = new(tls.Config)
	}
	if tLSConfig.Certificates == nil {
		tLSConfig.Certificates = []tls.Certificate{certificate}
	}
	if tLSConfig.GetConfigForClient == nil {
		tLSConfig.GetConfigForClient = GetConfigForClient
	}
	if tLSConfig.NextProtos == nil {
		tLSConfig.NextProtos = []string{"h2", "http/1.1"}
	}
	var server *http.Server
	if option.Server != nil {
		server = option.Server
	} else {
		server = new(http.Server)
	}
	if server.TLSNextProto == nil {
		server.TLSNextProto = map[string]func(*http.Server, *tls.Conn, http.Handler){
			"h2": TlsH2tProto,
		}
	}
	if server.ConnContext == nil {
		server.ConnContext = ConnContext
	}
	if server.Addr == "" {
		server.Addr = option.Addr
	}
	if server.Handler == nil {
		server.Handler = handler
	}
	if server.TLSConfig == nil {
		server.TLSConfig = tLSConfig
	}
	ln, err := Listen("tcp", option.Addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	return server.ServeTLS(ln, "", "")
}

func TlsH2tProto(s *http.Server, c *tls.Conn, h http.Handler) {
	ctx, fpContextData := ja3.CreateContext(context.TODO())
	if f, ok := c.NetConn().(interface{ ClientHelloData() []byte }); ok {
		fpContextData.SetClientHelloData(f.ClientHelloData())
	}
	(&http2.Server{}).ServeConn(c, &http2.ServeConnOpts{
		Context:    ctx,
		Handler:    h,
		BaseConfig: s,
	})
}
func ConnContext(ctx context.Context, c net.Conn) context.Context {
	return ja3.ConnContext(ctx, c)
}
func GetConfigForClient(chi *tls.ClientHelloInfo) (*tls.Config, error) {
	if fpContextData, ok := ja3.GetFpContextData(chi.Context()); ok {
		fpContextData.SetClientHelloInfo(*chi)
		if conn, ok := chi.Conn.(*Conn); ok {
			fpContextData.SetClientHelloData(conn.ClientHelloData())
		}
	}
	return nil, nil
}

func Listen(network string, address string) (*Listener, error) {
	lns, err := net.Listen(network, address)
	return &Listener{lns}, err
}

type Listener struct{ listen net.Listener }
type Conn struct {
	conn net.Conn
	raw  []byte
}

func (obj *Conn) ClientHelloData() []byte {
	return obj.raw
}
func (obj *Conn) Read(b []byte) (n int, err error) {
	i, err := obj.conn.Read(b)
	if err == nil && obj.raw == nil {
		obj.raw = make([]byte, i)
		copy(obj.raw, b)
	}
	return i, err
}
func (obj *Conn) Write(b []byte) (n int, err error)  { return obj.conn.Write(b) }
func (obj *Conn) Close() error                       { return obj.conn.Close() }
func (obj *Conn) LocalAddr() net.Addr                { return obj.conn.LocalAddr() }
func (obj *Conn) RemoteAddr() net.Addr               { return obj.conn.RemoteAddr() }
func (obj *Conn) SetDeadline(t time.Time) error      { return obj.conn.SetDeadline(t) }
func (obj *Conn) SetReadDeadline(t time.Time) error  { return obj.conn.SetReadDeadline(t) }
func (obj *Conn) SetWriteDeadline(t time.Time) error { return obj.conn.SetWriteDeadline(t) }

func (obj *Listener) Accept() (net.Conn, error) {
	conn, err := obj.listen.Accept()
	return &Conn{conn: conn}, err
}
func (obj *Listener) Close() error   { return obj.listen.Close() }
func (obj *Listener) Addr() net.Addr { return obj.listen.Addr() }
