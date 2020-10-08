package stream

import (
	"bytes"
	"io"
	"net"
	"time"
)

const (
	StateDirect State = 0
	StateShadow State = 1
	StateBlock  State = 2
)

type State int

type ShadowConn struct {
	net.Conn

	readState  State
	writeState State
	closeState State

	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
}

func (c *ShadowConn) Read(buf []byte) (int, error) {
	switch c.readState {
	case StateDirect:
		if c.readBuffer.Len() > 0 {
			return c.readBuffer.Read(buf)
		}

		return c.Conn.Read(buf)
	case StateShadow:
		n, err := c.Conn.Read(buf)
		if err == nil {
			c.readBuffer.Write(buf[:n])
		}

		return n, err
	case StateBlock:
		return 0, io.EOF
	}

	return 0, io.EOF
}

func (c *ShadowConn) Write(buf []byte) (int, error) {
	switch c.writeState {
	case StateDirect:
		if c.writeBuffer.Len() > 0 {
			_, err := c.writeBuffer.WriteTo(c.Conn)

			return 0, err
		}

		return c.Conn.Write(buf)
	case StateShadow:
		return c.writeBuffer.Write(buf)
	case StateBlock:
		return 0, io.ErrClosedPipe
	}

	return 0, io.ErrClosedPipe
}

func (c *ShadowConn) Close() error {
	switch c.closeState {
	case StateShadow, StateBlock:
		return nil
	case StateDirect:
		return c.Conn.Close()
	}

	return nil
}

func (c *ShadowConn) WriteTo(w io.Writer) (int64, error) {
	pending := int64(0)

	if c.readBuffer.Len() > 0 {
		if n, err := io.Copy(w, c.readBuffer); err != nil {
			return n, err
		} else {
			pending = n
		}
	}

	n, err := io.Copy(w, c.Conn)

	return pending + n, err
}

func (c *ShadowConn) ReadFrom(r io.Reader) (int64, error) {
	pending := int64(0)

	if c.writeBuffer.Len() > 0 {
		if n, err := io.Copy(c.Conn, c.writeBuffer); err != nil {
			return n, err
		} else {
			pending = n
		}
	}

	n, err := io.Copy(c.Conn, r)

	return pending + n, err
}

func (c *ShadowConn) SetReadState(state State) (original State) {
	original = c.readState

	c.readState = state

	return
}

func (c *ShadowConn) SetWriteState(state State) (original State) {
	original = c.writeState

	c.writeState = state

	return
}

func (c *ShadowConn) SetCloseState(state State) (original State) {
	original = c.closeState

	c.closeState = state

	return
}

func (c *ShadowConn) ResetDeadline() {
	_ = c.SetDeadline(time.Time{})
}

func NewShadowConn(conn net.Conn) *ShadowConn {
	return &ShadowConn{
		Conn:        conn,
		readState:   StateDirect,
		writeState:  StateDirect,
		closeState:  StateDirect,
		readBuffer:  bytes.NewBuffer(nil),
		writeBuffer: bytes.NewBuffer(nil),
	}
}
