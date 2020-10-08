package marker

import (
	"syscall"
)

type MonoMarker struct {
	name string
	mark uint32
}

func (m *MonoMarker) Name() string {
	return m.name
}

func (m *MonoMarker) Mark(conn syscall.RawConn) error {
	return MarkConn(conn, m.mark)
}

func NewMono(name string, mark uint32) Marker {
	return &MonoMarker{
		name: name,
		mark: mark,
	}
}
