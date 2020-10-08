package sniffer

import (
	"net"

	"github.com/kr328/snimark/stream"
)

func Tcp(_ *stream.ShadowConn, dst net.Addr) (string, error) {
	addr := dst.String()

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr, nil
	}

	return host, nil
}
