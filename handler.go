package fp

import (
	"context"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/gospider007/ja3"
	"github.com/gospider007/requests"
)

// example:  https://tools.scrapfly.io/api/fp/anything?extended=1
func GinHandlerFunc(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	fpData, ok := ja3.GetFpContextData(ctx.Request.Context())
	result := make(map[string]any)
	connectionState := fpData.ConnectionState()
	result["negotiatedProtocol"] = connectionState.NegotiatedProtocol
	result["tlsVersion"] = connectionState.Version
	result["userAgent"] = ctx.Request.UserAgent()
	result["orderHeaders"] = fpData.OrderHeaders()
	result["cookies"] = requests.Cookies(ctx.Request.Cookies()).String()
	tlsData, err := fpData.TlsData()
	if err == nil {
		clientHelloParseData := tlsData
		result["tls"] = clientHelloParseData
		result["ja3"], result["ja3n"] = tlsData.Fp()
		result["ja4"] = tlsData.Ja4()
		result["ja4h"] = fpData.Ja4H(ctx.Request)
	}
	h2Ja3Spec := fpData.H2Spec()
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
