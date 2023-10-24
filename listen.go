package fp

import (
	"context"
	"net"
)

func newListen(ctx context.Context, cnl context.CancelFunc, connPip chan net.Conn, addr net.Addr) *Listener {
	return &Listener{ctx: ctx, cnl: cnl, connPip: connPip, addr: addr}
}

type Listener struct {
	connPip chan net.Conn
	ctx     context.Context
	cnl     context.CancelFunc
	addr    net.Addr
}

func (obj *Listener) Accept() (net.Conn, error) {
	select {
	case <-obj.ctx.Done():
		return nil, obj.ctx.Err()
	case conn := <-obj.connPip:
		return conn, nil
	}
}
func (obj *Listener) Close() error {
	obj.cnl()
	return nil
}
func (obj *Listener) Addr() net.Addr { return obj.addr }
