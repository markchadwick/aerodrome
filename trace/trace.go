package trace

import (
	"code.google.com/p/goprotobuf/proto"
	"time"
)

func Start(msg ...string) *Trace {
	return &Trace{
		Start:   proto.Int64(time.Now().UnixNano()),
		Message: msg,
	}
}

func (t *Trace) Finish() {
	t.End = proto.Int64(time.Now().UnixNano())
}

func (t *Trace) Duration() time.Duration {
	if t.Start == nil || t.End == nil {
		return time.Duration(-1)
	}
	return time.Duration(*t.End - *t.Start)
}

// Add a trace as a child to this trace. The child will be inserted relative to
// the current set of children.
func (t *Trace) Add(child *Trace) {
	idx := t.indexOf(*child.Start)
	children := make([]*Trace, len(t.Child)+1)
	copy(children[0:idx], t.Child[0:idx])
	copy(children[idx+1:], t.Child[idx:])
	children[idx] = child
	t.Child = children
}

// Find the index in the `child` array for the given start time. This will not
// insert the element.
func (t *Trace) indexOf(v int64) int {
	children := len(t.Child)
	if children == 0 {
		return 0
	}
	return t.indexIn(v, 0, children-1)
}

// Find the appropriate index for a start date element in the given range
// assuming the current array of children is sorted.
func (t *Trace) indexIn(v int64, min, max int) int {
	if max < min {
		return min
	}
	mid := (min + max) / 2
	midValue := *t.Child[mid].Start

	if midValue > v {
		return t.indexIn(v, min, mid-1)
	} else if midValue < v {
		return t.indexIn(v, mid+1, max)
	} else {
		return mid
	}
}
