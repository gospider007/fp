package fp

import (
	"context"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/gospider007/ja3"
)

// example:  https://tools.scrapfly.io/api/fp/anything?extended=1
func GinHandlerFunc(ctx *gin.Context) {
	fpData, ok := ja3.GetFpContextData(ctx.Request.Context())
	result := make(map[string]any)
	connectionState := fpData.ConnectionState()
	result["negotiatedProtocol"] = connectionState.NegotiatedProtocol
	result["tlsVersion"] = connectionState.Version
	result["userAgent"] = ctx.Request.UserAgent()
	result["orderHeaders"] = fpData.OrderHeaders()
	clientHello, err := fpData.ClientHello()
	if err == nil {
		clientHelloParseData := clientHello.TlsData()
		result["tls"] = clientHelloParseData
		result["ja3"], result["ja3n"] = clientHelloParseData.Fp()
		result["ja4"] = fpData.Ja4()
	}
	h2Ja3Spec := fpData.H2Ja3Spec()
	result["http2"] = h2Ja3Spec
	result["akamai_fp"] = h2Ja3Spec.Fp()
	if ok {
		ctx.JSON(200, result)
	} else {
		ctx.JSON(200, map[string]any{
			"error": "fp load error",
		})
	}
}
func ConnContext(ctx context.Context, c net.Conn) context.Context {
	ja3Ctx, fpContextData := ja3.CreateContext(ctx)
	if conn, ok := c.(*TlsConn); ok {
		fpContextData.SetClientHelloData(conn.rawClientHello)
		fpContextData.SetConnectionState(conn.conn.ConnectionState())
		conn.fpContextData = fpContextData
	}
	return ja3Ctx
}
