package fp

import (
	"context"
	"net"
)

type keyPrincipalIDT string

const keyPrincipalID keyPrincipalIDT = "FpContextData"

func ConnContext(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, keyPrincipalID, c.(*tlsConn))
}

func GetRawConn(ctx context.Context) *tlsConn {
	return ctx.Value(keyPrincipalID).(*tlsConn)
}
