package fp

import (
	"bufio"
	"net"
)

func (obj *tlsConn) Read(b []byte) (n int, err error) {
	i, err := obj.Conn.Read(b)
	if !obj.saveOk {
		obj.rawContent = append(obj.rawContent, b[:i]...)
		obj.init()
	}
	return i, err
}

type WrapCon struct {
	net.Conn
	reader *bufio.Reader
}

func newWrapCon(conn net.Conn, reader *bufio.Reader) *WrapCon {
	return &WrapCon{Conn: conn, reader: reader}
}
func (obj *WrapCon) Read(b []byte) (n int, err error) { 
	return obj.reader.Read(b) 
}
