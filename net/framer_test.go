package net

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"github.com/markchadwick/spec"
	"io/ioutil"
	"time"
)

var _ = spec.Suite("Framer", func(c *spec.C) {
	rBuf := new(bytes.Buffer)
	wBuf := new(bytes.Buffer)

	client := NewFramer(rBuf, wBuf)
	defer client.Close()

	server := NewFramer(wBuf, rBuf)
	defer server.Close()

	c.It("should write a ping", func(c *spec.C) {
		ping := &Ping{
			Pong: proto.Bool(false),
			Id:   proto.Uint32(22),
		}

		pingBytes, _ := proto.Marshal(ping)

		err := client.writeFrame(PingFrame, ping)
		c.Assert(err).IsNil()

		wBytes, err := ioutil.ReadAll(wBuf)
		c.Assert(err).IsNil()

		w := bytes.NewBuffer(wBytes)
		var header int32
		err = binary.Read(w, binary.BigEndian, &header)
		c.Assert(err).IsNil()

		c.Assert(header >> 8).Equals(int32(len(pingBytes)))
		c.Assert(FrameType(header & 0xff)).Equals(PingFrame)
	})

	c.It("should read a ping", func(c *spec.C) {
		ping := &Ping{
			Pong: proto.Bool(false),
			Id:   proto.Uint32(15),
		}
		err := client.writeFrame(PingFrame, ping)
		c.Assert(err).IsNil()

		err = server.readFrame()
		c.Assert(err).IsNil()

		timeout := time.NewTimer(100 * time.Millisecond)
		select {
		case <-timeout.C:
			c.Failf("Never got my ping!")
		case p := <-server.Ping:
			c.Assert(*p.Pong).IsFalse()
			c.Assert(*p.Id).Equals(uint32(15))
		}
	})

	c.It("should detect and invalid frame type", func(c *spec.C) {
		c.Skip("pending")
	})
})
