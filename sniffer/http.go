package sniffer

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"github.com/kr328/snimark/stream"
)

func Http(conn *stream.ShadowConn, _ net.Addr) (string, error) {
	r := conn.SetReadState(stream.StateShadow)
	w := conn.SetWriteState(stream.StateBlock)
	c := conn.SetCloseState(stream.StateBlock)
	_ = conn.SetDeadline(time.Now().Add(DefaultTimeout))

	defer conn.SetReadState(r)
	defer conn.SetWriteState(w)
	defer conn.SetCloseState(c)
	defer conn.ResetDeadline()

	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return "", err
	}

	host := request.Host

	if host == "" {
		return "", ErrUnknownHost
	}

	return host, nil
}
