package platform

import (
	"encoding/binary"
	"net"
	"syscall"
	"unsafe"
)

const (
	SO_ORIGINAL_DST = 80 // from linux/include/uapi/linux/netfilter_ipv4.h
)

func DstAddr(conn net.Conn) (net.Addr, error) {
	c, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, ErrUnknownDst
	}

	sys, err := c.SyscallConn()
	if err != nil {
		return nil, ErrUnknownDst
	}

	var addr net.Addr

	err = sys.Control(func(fd uintptr) {
		raw := syscall.RawSockaddrInet4{}
		siz := unsafe.Sizeof(raw)

		if err = socketcall(
			GETSOCKOPT,
			fd,
			syscall.IPPROTO_IP,
			SO_ORIGINAL_DST,
			uintptr(unsafe.Pointer(&raw)),
			uintptr(unsafe.Pointer(&siz)),
			0); err != nil {
			return
		}

		port := []byte{
			*(*byte)(unsafe.Pointer(&raw.Port)),
			*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&raw.Port)) + 1)),
		}

		addr = &net.TCPAddr{
			IP:   raw.Addr[:],
			Port: int(binary.BigEndian.Uint16(port)),
			Zone: "",
		}
	})

	return addr, err
}
