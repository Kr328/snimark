// +build !linux

package platform

import (
	"net"
)

func DstAddr(conn net.Conn) (net.Addr, error) {
	return nil, ErrUnsupported
}
