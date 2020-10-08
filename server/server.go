package server

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"syscall"

	"github.com/kr328/snimark/config"
	"github.com/kr328/snimark/platform"
	"github.com/kr328/snimark/stream"
)

type Server struct {
	listen    string
	endpoints map[uint16]*Endpoint

	closed   bool
	listener net.Listener
}

func New(cfg *config.Config) (*Server, error) {
	eps, err := ParseEndpoints(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		listen:    cfg.ListenAddress,
		endpoints: eps,
	}, nil
}

func (e *Server) Exec() error {
	listener, err := net.Listen("tcp", e.listen)
	if err != nil {
		return err
	}

	e.listener = listener

	log.Printf("[INFO] Listen at %s\n", listener.Addr().String())

	for !e.closed {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[ERRO] Accept connection: %s\n", err.Error())
			continue
		}

		go e.handleConn(conn)
	}

	return nil
}

func (e *Server) Close() {
	e.closed = true

	if l := e.listener; l != nil {
		_ = l.Close()
	}
}

func (e *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	dst, err := platform.DstAddr(conn)
	if err != nil {
		log.Printf("[ERRO] Obtain destination of %s: %s\n", conn.RemoteAddr().String(), err.Error())
		return
	}

	_, port, err := net.SplitHostPort(dst.String())
	if err != nil {
		log.Printf("[ERRO] Invalid destination: %s\n", dst.String())
		return
	}

	dialer := net.Dialer{}

	p, _ := strconv.Atoi(port)

	host := ""
	marker := "DIRECT"

	ep, existed := e.endpoints[uint16(p)]
	if existed {
		s := stream.NewShadowConn(conn)

		if h, err := ep.sniffer(s, dst); err == nil {
			host = h

			if ep.matcher.Match(h) {
				marker = ep.marker.Name()

				dialer.Control = func(network, address string, c syscall.RawConn) (err error) {
					if err = ep.marker.Mark(c); err != nil {
						log.Printf("[ERRO] mark %s: %s\n", conn.RemoteAddr().String(), err.Error())
					}

					return
				}
			}
		}

		conn = s
	}

	if host == "" {
		log.Printf("[INFO] %s -> %s using %s\n", conn.RemoteAddr().String(), dst.String(), marker)
	} else {
		log.Printf("[INFO] %s -> %s/%s using %s\n", conn.RemoteAddr().String(), host, dst.String(), marker)
	}

	remote, err := dialer.DialContext(context.Background(), "tcp", dst.String())
	if err != nil {
		log.Printf("[ERRO] %s\n", err.Error())
		return
	}

	relay(conn, remote)
}

func relay(left net.Conn, right net.Conn) {
	defer left.Close()
	defer right.Close()

	ch := make(chan struct{})

	go func() {
		if r, ok := left.(io.ReaderFrom); ok {
			_, _ = r.ReadFrom(right)
		} else {
			_, _ = io.Copy(left, right)
		}

		ch <- struct{}{}
	}()

	if w, ok := left.(io.WriterTo); ok {
		_, _ = w.WriteTo(right)
	} else {
		_, _ = io.Copy(right, left)
	}

	<-ch
}
