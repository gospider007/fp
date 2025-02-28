package fp

import (
	"context"
	"net"

	"github.com/gin-gonic/gin"
)

// example:  https://tools.scrapfly.io/api/fp/anything?extended=1
func GinHandlerFunc(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	rawConn := GetRawConn(ctx.Request.Context())
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
	ctx.JSON(200, results)
}

type keyPrincipalIDT string

const keyPrincipalID keyPrincipalIDT = "FpContextData"

func ConnContext(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, keyPrincipalID, c.(*tlsConn))
}

func GetRawConn(ctx context.Context) *tlsConn {
	return ctx.Value(keyPrincipalID).(*tlsConn)
}
