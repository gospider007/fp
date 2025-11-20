package fp

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gospider007/gtls"
)

type Option struct {
	Addr      string
	TLSConfig *tls.Config
}

func Server(handler http.Handler, options ...Option) (err error) {
	var option Option
	if len(options) > 0 {
		option = options[0]
	}
	if option.Addr == "" {
		option.Addr = ":8999"
		log.Print("Starting server on https://localhost:8999")
		// log.Print("root certificate on https://localhost:8999/cert")
	}
	if option.TLSConfig == nil {
		option.TLSConfig = &tls.Config{
			GetCertificate:     func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) { return gtls.GetCertificate(chi, nil, nil) },
			InsecureSkipVerify: true,
			NextProtos:         []string{"h2", "http/1.1"},
		}
	}
	ln, err := NewListen(option.Addr, handler, option.TLSConfig)
	if err != nil {
		return err
	}
	defer ln.Close()
	return (&http.Server{ConnContext: ConnContext, Handler: handler}).Serve(ln)
}

func fingerprintHandler(c *gin.Context) {
	// 设置 CORS 和连接头
	c.Header("Connection", "close")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Content-Security-Policy", "default-src 'self'; connect-src https://localhost:8999")

	c.Header("Access-Control-Allow-Headers", "*")

	// 从 context 获取原始连接对象
	rawConn := GetRawConn(c.Request.Context())

	tlsSpec := rawConn.TLSSpec()
	h1Spec := rawConn.H1Spec()
	h2Spec := rawConn.H2Spec()

	results := map[string]any{
		"tls":          tlsSpec.Map(),
		"goSpiderSpec": rawConn.GoSpiderSpec(),
	}

	if h1Spec != nil {
		results["h1"] = h1Spec.Map()
	}
	if h2Spec != nil {
		results["h2"] = h2Spec.Map()
	}
	// 返回 JSON 响应
	c.JSON(200, results)
}
func certHander(c *gin.Context) {
	certPath, _, _ := gtls.GetRootCertWithLocal(true)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/x-pem-file")
	c.Header("Content-Disposition", "attachment; filename=cert.pem")
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(certPath)
}
func optionHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Status(200)

}

func Start(option ...Option) error {
	ginServer := gin.Default()
	ginServer.GET("/", fingerprintHandler)
	ginServer.GET("/cert", certHander)
	ginServer.OPTIONS("/*path", optionHandler)
	return Server(ginServer, option...)
}
