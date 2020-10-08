package sniffer

import (
	"errors"
	"net"
	"time"

	"github.com/kr328/snimark/stream"
)

type Sniffer func(conn *stream.ShadowConn, dst net.Addr) (string, error)

var (
	ErrUnknownHost = errors.New("unknown host")

	DefaultTimeout = time.Second * 2

	Sniffers = map[string]Sniffer{
		"http": Http,
		"tls":  Tls,
		"tcp":  Tcp,
	}
)
