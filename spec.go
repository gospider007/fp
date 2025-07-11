package fp

import (
	"crypto/tls"
	"errors"
	"strings"

	"github.com/gospider007/http2"
	"github.com/gospider007/ja3"
	"github.com/gospider007/requests"
)

type tlsConn struct {
	conn            *tls.Conn
	rawClientHello  []byte
	rawContent      []byte
	saveOk          bool
	connectionState tls.ConnectionState
	h2Spec          *http2.Spec
	h1Spec          *requests.Spec
}

func newTlsConn(conn *tls.Conn, rawClientHello []byte) *tlsConn {
	return &tlsConn{
		conn:            conn,
		connectionState: conn.ConnectionState(),
		rawClientHello:  rawClientHello,
		rawContent:      []byte{},
	}
}

func (obj *tlsConn) SetOk() {
	obj.saveOk = true
}
func (obj *tlsConn) ConnectionState() tls.ConnectionState {
	return obj.connectionState
}
func (obj *tlsConn) TLSSpec() *ja3.Spec {
	spec, _ := ja3.ParseSpec(obj.rawClientHello)
	return spec
}
func (obj *tlsConn) Content() []byte {
	return obj.rawContent
}
func (obj *tlsConn) H1Spec() *requests.Spec {
	return obj.h1Spec
}
func (obj *tlsConn) H2Spec() *http2.Spec {
	return obj.h2Spec
}
func (obj *tlsConn) GoSpiderSpec() string {
	results := []string{obj.TLSSpec().Hex()}
	if h1Spec := obj.H1Spec(); h1Spec != nil {
		results = append(results, h1Spec.Hex())
	} else {
		results = append(results, "")
	}
	if h2Spec := obj.H2Spec(); h2Spec != nil {
		results = append(results, h2Spec.Hex())
	} else {
		results = append(results, "")
	}
	return strings.Join(results, "@")
}
func (obj *tlsConn) init() error {
	switch obj.connectionState.NegotiatedProtocol {
	case "h2":
		return obj.initH2()
	case "http/1.1":
		return obj.initH1()
	}
	obj.saveOk = true
	return errors.New("unknown protocol")
}
func (obj *tlsConn) initH1() error {
	if obj.saveOk {
		return nil
	}
	spec, err := requests.ParseSpec(obj.rawContent)
	if err != nil {
		return err
	}
	obj.h1Spec = spec
	obj.saveOk = true
	return nil
}

func (obj *tlsConn) initH2() error {
	if obj.saveOk {
		return nil
	}
	spec, err := http2.ParseSpec(obj.rawContent)
	if err != nil {
		return err
	}
	obj.h2Spec = spec
	obj.saveOk = true
	return nil
}
