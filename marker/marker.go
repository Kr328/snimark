package marker

import (
	"syscall"
)

type Marker interface {
	Name() string
	Mark(conn syscall.RawConn) error
}
