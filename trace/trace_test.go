package trace

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/markchadwick/spec"
	"testing"
	"time"
)

var _ = spec.Suite("Traces", func(c *spec.C) {

	c.It("should initialize", func(c *spec.C) {
		t := Start()
		c.Assert(t).NotNil()

		c.Assert(t.Start).NotNil()
	})

	c.It("should report a duration", func(c *spec.C) {
		t := Start()
		t.Finish()

		c.Assert(t.Duration() > 0)
	})

	c.It("should not report a duration on an unfinished trace", func(c *spec.C) {
		t := Start()
		c.Assert(t.Duration()).Equals(time.Duration(-1))
	})

	c.It("should intialize with optional messages", func(c *spec.C) {
		t1 := Start()
		c.Assert(t1.Message).HasLen(0)

		t2 := Start("Hello", "There")
		c.Assert(t2.Message).HasLen(2)
		c.Assert(t2.Message[0]).Equals("Hello")
		c.Assert(t2.Message[1]).Equals("There")
	})

	c.It("should add a child event", func(c *spec.C) {
		parent := Start()
		child := Start()
		child.Finish()

		c.Assert(parent.Child).HasLen(0)
		parent.Add(child)
		c.Assert(parent.Child).HasLen(1)
		c.Assert(parent.Child[0]).Equals(child)
	})

	c.It("when inserting with 2 children", func(c *spec.C) {
		t := Start()
		t.Add(&Trace{Start: proto.Int64(5)})
		t.Add(&Trace{Start: proto.Int64(7)})

		c.Assert(t.Child).HasLen(2)
		c.Assert(*t.Child[0].Start).Equals(int64(5))
		c.Assert(*t.Child[1].Start).Equals(int64(7))

		c.It("should insert an earlier child first", func(c *spec.C) {
			t.Add(&Trace{Start: proto.Int64(4)})
			c.Assert(t.Child).HasLen(3)
			c.Assert(*t.Child[0].Start).Equals(int64(4))
			c.Assert(*t.Child[1].Start).Equals(int64(5))
			c.Assert(*t.Child[2].Start).Equals(int64(7))
		})

		c.It("should insert a later child last", func(c *spec.C) {
			t.Add(&Trace{Start: proto.Int64(11)})
			c.Assert(t.Child).HasLen(3)
			c.Assert(*t.Child[0].Start).Equals(int64(5))
			c.Assert(*t.Child[1].Start).Equals(int64(7))
			c.Assert(*t.Child[2].Start).Equals(int64(11))
		})

		c.It("should insert a middle child in the middle", func(c *spec.C) {
			t.Add(&Trace{Start: proto.Int64(6)})
			c.Assert(t.Child).HasLen(3)
			c.Assert(*t.Child[0].Start).Equals(int64(5))
			c.Assert(*t.Child[1].Start).Equals(int64(6))
			c.Assert(*t.Child[2].Start).Equals(int64(7))
		})
	})

	c.It("when inserting a child", func(c *spec.C) {
		t := Start()
		t.Add(&Trace{Start: proto.Int64(13)})
		t.Add(&Trace{Start: proto.Int64(5)})
		t.Add(&Trace{Start: proto.Int64(11)})
		t.Add(&Trace{Start: proto.Int64(7)})
		t.Add(&Trace{Start: proto.Int64(3)})

		c.It("should insert the lowest start first", func(c *spec.C) {
			c.Assert(t.indexOf(2)).Equals(0)
		})

		c.It("should insert the highest start last", func(c *spec.C) {
			c.Assert(t.indexOf(17)).Equals(5)
		})

		c.It("should be ordered", func(c *spec.C) {
			c.Assert(t.indexOf(1)).Equals(0)
			c.Assert(t.indexOf(2)).Equals(0)
			c.Assert(t.indexOf(3)).Equals(0)
			c.Assert(t.indexOf(4)).Equals(1)
			c.Assert(t.indexOf(5)).Equals(1)
			c.Assert(t.indexOf(6)).Equals(2)
			c.Assert(t.indexOf(7)).Equals(2)
			c.Assert(t.indexOf(8)).Equals(3)
			c.Assert(t.indexOf(9)).Equals(3)
			c.Assert(t.indexOf(10)).Equals(3)
			c.Assert(t.indexOf(11)).Equals(3)
			c.Assert(t.indexOf(12)).Equals(4)
			c.Assert(t.indexOf(13)).Equals(4)
			c.Assert(t.indexOf(14)).Equals(5)
		})
	})
})

func Test(t *testing.T) {
	spec.Run(t)
}
