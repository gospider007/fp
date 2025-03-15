package fp

import (
	"bufio"
	"net"
	"time"
)

func (obj *tlsConn) Read(b []byte) (n int, err error) {
	i, err := obj.conn.Read(b)
	if !obj.saveOk {
		obj.rawContent = append(obj.rawContent, b[:i]...)
		obj.init()
	}
	return i, err
}
func (obj *tlsConn) Write(b []byte) (n int, err error) { return obj.conn.Write(b) }
func (obj *tlsConn) Close() error {
	return obj.conn.Close()
}
func (obj *tlsConn) LocalAddr() net.Addr                { return obj.conn.LocalAddr() }
func (obj *tlsConn) RemoteAddr() net.Addr               { return obj.conn.RemoteAddr() }
func (obj *tlsConn) SetDeadline(t time.Time) error      { return obj.conn.SetDeadline(t) }
func (obj *tlsConn) SetReadDeadline(t time.Time) error  { return obj.conn.SetReadDeadline(t) }
func (obj *tlsConn) SetWriteDeadline(t time.Time) error { return obj.conn.SetWriteDeadline(t) }

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
