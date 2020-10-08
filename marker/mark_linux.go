package marker

import (
	"syscall"
)

func MarkConn(conn syscall.RawConn, mark uint32) (err error) {
	err = conn.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, int(mark))
	})

	return
}
