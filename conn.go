package fp

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/gospider007/ja3"
	"github.com/gospider007/tools"
)

type TlsConn struct {
	conn           *tls.Conn
	rawClientHello []byte
	fpContextData  *ja3.FpContextData
	rawHeaders     []byte
	cnl            context.CancelFunc
}

func newTlsConn(cnl context.CancelFunc, conn *tls.Conn, rawClientHello []byte) *TlsConn {
	return &TlsConn{cnl: cnl, conn: conn, rawClientHello: rawClientHello, rawHeaders: []byte{}}
}
func (obj *TlsConn) Read(b []byte) (n int, err error) {
	i, err := obj.conn.Read(b)
	if obj.fpContextData != nil && obj.rawHeaders != nil && err == nil && i > 0 {
		obj.rawHeaders = append(obj.rawHeaders, b[:i]...)
		firstIndex := bytes.LastIndex(obj.rawHeaders, []byte{' ', 'H', 'T', 'T', 'P', '/', '1', '.', '1', '\r', '\n'})
		lastIndex := bytes.Index(obj.rawHeaders, []byte{'\r', '\n', '\r', '\n'})
		if firstIndex != -1 {
			if lastIndex > firstIndex+11 {
				headers := obj.rawHeaders[firstIndex+11 : lastIndex]
				orderHeaders := []string{}
				for _, kvs := range bytes.Split(headers, []byte{'\r', '\n'}) {
					if index := bytes.Index(kvs, []byte{':'}); index > 0 {
						orderHeaders = append(orderHeaders, tools.BytesToString(kvs[:index]))
					}
				}
				obj.fpContextData.SetOrderHeaders(orderHeaders)
				obj.rawHeaders = nil
				obj.fpContextData = nil
			} else if firstIndex > 0 {
				obj.rawHeaders = obj.rawHeaders[firstIndex:]
			}
		} else if total := len(obj.rawHeaders); total >= 11 {
			obj.rawHeaders = obj.rawHeaders[total-11 : total]

		}
	}
	return i, err
}
func (obj *TlsConn) Write(b []byte) (n int, err error) { return obj.conn.Write(b) }
func (obj *TlsConn) Close() error {
	obj.cnl()
	return obj.conn.Close()
}
func (obj *TlsConn) LocalAddr() net.Addr                { return obj.conn.LocalAddr() }
func (obj *TlsConn) RemoteAddr() net.Addr               { return obj.conn.RemoteAddr() }
func (obj *TlsConn) SetDeadline(t time.Time) error      { return obj.conn.SetDeadline(t) }
func (obj *TlsConn) SetReadDeadline(t time.Time) error  { return obj.conn.SetReadDeadline(t) }
func (obj *TlsConn) SetWriteDeadline(t time.Time) error { return obj.conn.SetWriteDeadline(t) }

type WrapCon struct {
	rawConn net.Conn
	reader  *bufio.Reader
}

func newWrapCon(conn net.Conn, reader *bufio.Reader) *WrapCon {
	return &WrapCon{rawConn: conn, reader: reader}
}
func (obj *WrapCon) Read(b []byte) (n int, err error)   { return obj.reader.Read(b) }
func (obj *WrapCon) Write(b []byte) (n int, err error)  { return obj.rawConn.Write(b) }
func (obj *WrapCon) Close() error                       { return obj.rawConn.Close() }
func (obj *WrapCon) LocalAddr() net.Addr                { return obj.rawConn.LocalAddr() }
func (obj *WrapCon) RemoteAddr() net.Addr               { return obj.rawConn.RemoteAddr() }
func (obj *WrapCon) SetDeadline(t time.Time) error      { return obj.rawConn.SetDeadline(t) }
func (obj *WrapCon) SetReadDeadline(t time.Time) error  { return obj.rawConn.SetReadDeadline(t) }
func (obj *WrapCon) SetWriteDeadline(t time.Time) error { return obj.rawConn.SetWriteDeadline(t) }
