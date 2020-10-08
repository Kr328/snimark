package marker

import (
	"sync/atomic"
	"syscall"
)

type RoundMarker struct {
	name    string
	marks   []uint32
	current uint32
}

func (r *RoundMarker) Name() string {
	return r.name
}

func (r *RoundMarker) Mark(conn syscall.RawConn) error {
	c := atomic.AddUint32(&r.current, 1) % uint32(len(r.marks))

	return MarkConn(conn, r.marks[c])
}

func NewRound(name string, marks []uint32) Marker {
	return &RoundMarker{
		name:    name,
		marks:   marks,
		current: 0,
	}
}
