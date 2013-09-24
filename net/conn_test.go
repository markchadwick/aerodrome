package net

import (
	"github.com/markchadwick/spec"
	"net"
)

var _ = spec.Suite("A Connection", func(c *spec.C) {
	rawClient, rawServer := net.Pipe()
	client := NewConn(rawClient)
	server := NewConn(rawServer)

	c.It("should have a simple message written and read", func(c *spec.C) {
		go func() {
			n, err := client.Write([]byte("Hello, world!"))
			c.Assert(err).IsNil()
			c.Assert(n).Equals(13)
		}()

		bs := make([]byte, 7)
		n, err := server.Read(bs)
		c.Assert(err).IsNil()
		c.Assert(n).Equals(7)
		c.Assert(string(bs)).Equals("Hello, ")

		n, err = server.Read(bs)
		c.Assert(err).IsNil()
		c.Assert(n).Equals(6)
		c.Assert(string(bs[:n])).Equals("world!")
	})

  c.It("should read a fragmented message", func(c *spec.C) {
    c.Skip("pending")
  })
})
