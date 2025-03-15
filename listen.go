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

func NewListen(addr string, handler http.Handler, tlsConfig *tls.Config) (*Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Listener{
		listener:  listener,
		tlsConfig: tlsConfig,
		handler:   handler,
	}, nil
}

type Listener struct {
	listener  net.Listener
	tlsConfig *tls.Config
	handler   http.Handler
}

func (obj *Listener) Accept() (net.Conn, error) {
	for {
		con, err := obj.listener.Accept()
		if err != nil {
			return nil, err
		}
		wrapConn, err := obj.mainHandle(con)
		if err != nil {
			con.Close()
		} else if wrapConn != nil {
			return wrapConn, nil
		}
	}
}
func (obj *Listener) Close() error {
	return obj.listener.Close()
}
func (obj *Listener) Addr() net.Addr { return obj.listener.Addr() }
func (obj *Listener) mainHandle(client net.Conn) (net.Conn, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	clientReader := bufio.NewReader(client)
	peekBytes, err := clientReader.Peek(1)
	if err != nil {
		return nil, err
	}
	if peekBytes[0] != 22 {
		return nil, errors.New("protol error")
	}
	if peekBytes, err = clientReader.Peek(5); err != nil {
		return nil, err
	}
	if peekBytes, err = clientReader.Peek(int(binary.BigEndian.Uint16(peekBytes[3:5])) + 5); err != nil {
		return nil, err
	}
	rawClientHello := make([]byte, len(peekBytes))
	copy(rawClientHello, peekBytes)
	tlsClient := tls.Server(newWrapCon(client, clientReader), obj.tlsConfig)
	if err := tlsClient.Handshake(); err != nil {
		return nil, err
	}
	tlsConn := newTlsConn(tlsClient, rawClientHello)
	if tlsClient.ConnectionState().NegotiatedProtocol == "h2" {
		go func() {
			(&http2.Server{}).ServeConn(tlsConn, &http2.ServeConnOpts{
				Context: ConnContext(context.TODO(), tlsConn),
				Handler: obj.handler,
			})
			tlsClient.Close()
		}()
		return nil, nil
	} else {
		if tlsConn.connectionState.NegotiatedProtocol != "http/1.1" {
			tlsConn.saveOk = true
		}
		return tlsConn, nil
	}
}
