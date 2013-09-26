//
// +---------------------------------------+
// |   Length (24 bits)   | Frame Type (8) |
// +---------------------------------------+
// | Data....                              |
// +---------------------------------------+

package net

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const MaxMessageSize = 1<<24 - 1

type FrameType uint8

const (
	PingFrame FrameType = iota
)

type Framer struct {
	r        io.Reader
	w        io.Writer
	stopRead chan bool
	// stopWrite chan bool
	Ping chan (*Ping)
}

func NewFramer(r io.Reader, w io.Writer) *Framer {
	framer := &Framer{
		r:        r,
		w:        w,
		stopRead: make(chan bool),
		// stopWrite: make(chan bool),
		Ping: make(chan *Ping),
	}
	go framer.read()
	return framer
}

func (f *Framer) Close() {
	defer func() {
		if r := recover(); r != nil {
			// Ignore for now
		}
	}()
	close(f.Ping)
	f.stopRead <- true
	// f.stopWrite <- true
	close(f.stopRead)
}

func (f *Framer) read() {
	for {
		select {
		default:
			if err := f.readFrame(); err != nil {
				if err == io.EOF {
					f.Close()
				} else {
					log.Printf("Error reading frame: %s", err.Error())
				}
			}
		case <-f.stopRead:
			return
		}
	}
}

func (f Framer) writeFrame(typ FrameType, msg proto.Message) (err error) {
	bs, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	length := len(bs)
	if length > MaxMessageSize {
		return fmt.Errorf("Message too large: %d > %d", length, MaxMessageSize)
	}

	header := int32(length << 8)
	header |= int32(typ & 0xff)

	if err = binary.Write(f.w, binary.BigEndian, header); err != nil {
		return
	}
	_, err = f.w.Write(bs)
	return
}

func (f *Framer) readFrame() (err error) {
	var header int32
	if err = binary.Read(f.r, binary.BigEndian, &header); err != nil {
		return
	}
	length := header >> 8
	frameType := uint8(header & 0xff)

	// TODO: recycle these bufs -- see bufpool
	buf := make([]byte, length)
	_, err = io.ReadFull(f.r, buf)
	if err != nil {
		return err
	}
	go f.parseFrame(FrameType(frameType), buf)
	return
}

func (f *Framer) parseFrame(typ FrameType, buf []byte) {
	var err error
	switch typ {
	default:
		log.Printf("Unknown frame type: %d", typ)
	case PingFrame:
		msg := new(Ping)
		if err = proto.Unmarshal(buf, msg); err != nil {
			f.Ping <- msg
		}
	}
}
