package sniffer

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/kr328/snimark/stream"
)

func Tls(conn *stream.ShadowConn, _ net.Addr) (string, error) {
	r := conn.SetReadState(stream.StateShadow)
	w := conn.SetWriteState(stream.StateBlock)
	c := conn.SetCloseState(stream.StateBlock)
	_ = conn.SetDeadline(time.Now().Add(DefaultTimeout))

	defer conn.SetReadState(r)
	defer conn.SetWriteState(w)
	defer conn.SetCloseState(c)
	defer conn.ResetDeadline()

	tc := tls.Server(conn, &tls.Config{})
	_ = tc.Handshake()

	host := tc.ConnectionState().ServerName

	if host == "" {
		return "", ErrUnknownHost
	}

	return host, nil
}
