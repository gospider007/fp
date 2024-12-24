package fp

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/gospider007/gtls"
	"github.com/gospider007/ja3"
	"github.com/gospider007/net/http2"
)

type Option struct {
	Addr        string
	CertFile    string
	KeyFile     string
	Certificate tls.Certificate
	NextProtos  []string
	Handler     http.Handler
	DomainNames []string
}

func newTlsConfig(option Option) (*tls.Config, error) {
	var tlsConfig *tls.Config
	if option.Certificate.Certificate != nil {
		tlsConfig = &tls.Config{Certificates: []tls.Certificate{option.Certificate}}
	} else if option.CertFile != "" && option.KeyFile != "" {
		if certData, err := os.ReadFile(option.CertFile); err != nil {
			return tlsConfig, err
		} else if cert, err := gtls.LoadCert(certData); err != nil {
			return tlsConfig, err
		} else if keyData, err := os.ReadFile(option.KeyFile); err != nil {
			return tlsConfig, err
		} else if key, err := gtls.LoadCertKey(keyData); err != nil {
			return tlsConfig, err
		} else if certificate, err := gtls.CreateTlsCert(cert, key); err != nil {
			return tlsConfig, err
		} else {
			tlsConfig = &tls.Config{Certificates: []tls.Certificate{certificate}}
		}
	} else if option.DomainNames != nil {
		return gtls.TLS(option.DomainNames)
	} else if certificate, err := gtls.CreateProxyCertWithName("test"); err != nil {
		return tlsConfig, err
	} else {
		tlsConfig = &tls.Config{Certificates: []tls.Certificate{certificate}}
	}
	if tlsConfig.NextProtos == nil {
		if option.NextProtos == nil {
			tlsConfig.NextProtos = []string{"h2", "http/1.1"}
		} else {
			tlsConfig.NextProtos = option.NextProtos
		}
	}
	return tlsConfig, nil
}

type server struct {
	connPip   chan net.Conn
	listener  net.Listener
	tlsConfig *tls.Config
	handler   http.Handler
	addr      net.Addr
	ctx       context.Context
	cnl       context.CancelFunc
	err       error
}

func (obj *server) close() {
	obj.cnl()
}
func (obj *server) mainHandle(preCtx context.Context, client net.Conn) (err error) {
	defer recover()
	if client == nil {
		return errors.New("client is nil")
	}
	defer client.Close()
	clientReader := bufio.NewReader(client)

	peekBytes, err := clientReader.Peek(1)
	if err != nil {
		return err
	}
	if peekBytes[0] != 22 {
		return errors.New("protol error")
	}
	if peekBytes, err = clientReader.Peek(5); err != nil {
		return err
	}
	if peekBytes, err = clientReader.Peek(int(binary.BigEndian.Uint16(peekBytes[3:5])) + 5); err != nil {
		return err
	}
	rawClientHello := make([]byte, len(peekBytes))
	copy(rawClientHello, peekBytes)
	tlsClient := tls.Server(newWrapCon(client, clientReader), obj.tlsConfig)
	defer tlsClient.Close()
	if err := tlsClient.HandshakeContext(preCtx); err != nil {
		return err
	}
	if tlsClient.ConnectionState().NegotiatedProtocol == "h2" {
		ja3Ctx, ja3Context := ja3.CreateContext(preCtx)
		ja3Context.SetClientHelloData(rawClientHello)
		ja3Context.SetConnectionState(tlsClient.ConnectionState())
		(&http2.Server{CloseCallBack: func() bool {
			select {
			case <-preCtx.Done():
				return true
			default:
				return false
			}
		}}).ServeConn(tlsClient, &http2.ServeConnOpts{
			Context: ja3Ctx,
			Handler: obj.handler,
		})
	} else {
		tlsCtx, tlsCnl := context.WithCancel(preCtx)
		select {
		case <-preCtx.Done():
			return context.Cause(preCtx)
		case obj.connPip <- newTlsConn(tlsCnl, tlsClient, rawClientHello):
			<-tlsCtx.Done()
		}
	}
	return nil
}
func (obj *server) listen() error {
	defer obj.close()
	ln := newListen(obj.ctx, obj.cnl, obj.connPip, obj.addr)
	defer ln.Close()
	return (&http.Server{ConnContext: ConnContext, Handler: obj.handler}).Serve(ln)
}
func (obj *server) serve() error {
	defer obj.close()
	for {
		select {
		case <-obj.ctx.Done():
			obj.err = context.Cause(obj.ctx)
			return obj.err
		default:
			client, err := obj.listener.Accept() //接受数据
			if err != nil {
				obj.err = err
				return err
			}
			go obj.mainHandle(obj.ctx, client)
		}
	}
}
func Server(ctx context.Context, handler http.Handler, options ...Option) (err error) {
	var option Option
	if len(options) > 0 {
		option = options[0]
	}
	if option.Addr == "" {
		option.Addr = ":0"
	}
	if ctx == nil {
		ctx = context.TODO()
	}
	server := &server{
		handler: handler,
		connPip: make(chan net.Conn),
	}
	server.ctx, server.cnl = context.WithCancel(ctx)
	if server.listener, err = net.Listen("tcp", option.Addr); err != nil {
		return err
	}
	if server.tlsConfig, err = newTlsConfig(option); err != nil {
		return err
	}
	go server.listen()
	return server.serve()
}
