package listener

import (
	"net"
)

// -----------------------------------------------------------------------------

type gracefulConn struct {
	net.Conn
	ln *gracefulListener
}

// -----------------------------------------------------------------------------

func (c *gracefulConn) Close() error {
	err := c.Conn.Close()
	if err == nil {
		c.ln.connClosed()
	}

	return err
}
