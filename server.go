package fp

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"net"
	"net/http"

	"golang.org/x/net/http2"
)

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
	tlsCtx, tlsCnl := context.WithCancel(preCtx)
	tlsConn := newTlsConn(tlsCnl, tlsClient, rawClientHello)
	if tlsClient.ConnectionState().NegotiatedProtocol == "h2" {
		(&http2.Server{}).ServeConn(tlsConn, &http2.ServeConnOpts{
			Context: ConnContext(preCtx, tlsConn),
			Handler: obj.handler,
		})
	} else {
		if tlsConn.connectionState.NegotiatedProtocol != "http/1.1" {
			tlsConn.saveOk = true
		}
		select {
		case <-preCtx.Done():
			return context.Cause(preCtx)
		case obj.connPip <- tlsConn:
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
