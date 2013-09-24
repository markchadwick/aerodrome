package net

import (
	"bytes"
	"net"
	// "code.google.com/p/snappy-go/snappy"
	"encoding/binary"
	"io"
	"log"
	"time"
)

// This assumes that all frames for a message will be squential. This may or may
// not be the case
// TODO: interleave messages above the frame size to see what happens
type Conn struct {
	conn net.Conn
}

// Wrap a connection in an aerodrome connection. This will add metadata to each
// frame and expect the same from the endpoint its talking to.
func NewConn(conn net.Conn) *Conn {
	return &Conn{
		conn: conn,
	}
}

func (c *Conn) Read(b []byte) (int, error) {
	var length uint32

	log.Printf("Read: reading length")
	if err := binary.Read(c.conn, binary.BigEndian, &length); err != nil {
		return 0, err
	}
	log.Printf("Read: length: %d", length)

	b = make([]byte, length)
	log.Printf("Read: read fully")
	if n, err := io.ReadFull(c.conn, b); err != nil {
		log.Printf("Read: couldn't read fully %s", err.Error())
		return n, err
	}

	log.Printf("Read: OK! %s", b)
	return int(length), nil
}

func (c *Conn) read(length int) []byte {
	return nil
}

// Do a blocking write to the underlying connection until all frames have been
// written. This may translate to many underlying blocking operations. Each
// frame will be encoded is per the spec above.
func (c *Conn) Write(b []byte) (int, error) {
	buf := new(bytes.Buffer)
	length := uint32(len(b))

	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return 0, err
	}
	if _, err := buf.Write(b); err != nil {
		return 0, err
	}
	if _, err := io.Copy(c.conn, buf); err != nil {
		return 0, err
	}
	return int(length), nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// Ensure that our Conn can indeed pass as a net.Conn
var _ net.Conn = &Conn{}
