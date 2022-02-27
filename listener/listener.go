package listener

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

// -----------------------------------------------------------------------------

// gracefulListener defines a listener that we can gracefully stop
type gracefulListener struct {
	// inner listener
	ln net.Listener

	// maximum wait time for graceful shutdown
	maxShutdownTime time.Duration

	// this channel is closed during graceful shutdown on zero open connections.
	done chan struct{}

	// the number of open connections
	activeConnectionsCount uint64

	// becomes non-zero when graceful shutdown starts
	shutdown uint64
}

// -----------------------------------------------------------------------------

// NewGracefulListener wraps the given listener into 'graceful shutdown' listener.
func NewGracefulListener(ln net.Listener, maxShutdownTime time.Duration) net.Listener {
	return &gracefulListener{
		ln:              ln,
		maxShutdownTime: maxShutdownTime,
		done:            make(chan struct{}),
	}
}

// Accept creates a conn
func (ln *gracefulListener) Accept() (net.Conn, error) {
	var c net.Conn
	var err error

	c, err = ln.ln.Accept()
	if err != nil {
		return nil, err
	}

	atomic.AddUint64(&ln.activeConnectionsCount, 1)

	return &gracefulConn{
		Conn: c,
		ln:   ln,
	}, nil
}

// Addr returns the listen address
func (ln *gracefulListener) Addr() net.Addr {
	return ln.ln.Addr()
}

// Close closes the inner listener and waits until all the pending
// open connections are closed before returning.
func (ln *gracefulListener) Close() error {
	var err error

	err = ln.ln.Close()
	if err != nil {
		return err
	}

	return ln.waitForZeroConns()
}

// -----------------------------------------------------------------------------
// Private methods

func (ln *gracefulListener) waitForZeroConns() error {
	atomic.AddUint64(&ln.shutdown, 1)

	if atomic.LoadUint64(&ln.activeConnectionsCount) == 0 {
		close(ln.done)
		return nil
	}

	select {
	case <-ln.done:
		return nil
	case <-time.After(ln.maxShutdownTime):
		return fmt.Errorf("unable to complete graceful shutdown in %s", ln.maxShutdownTime)
	}
}

func (ln *gracefulListener) connClosed() {
	activeConnectionsCount := atomic.AddUint64(&ln.activeConnectionsCount, ^uint64(0))

	if atomic.LoadUint64(&ln.shutdown) != 0 && activeConnectionsCount == 0 {
		close(ln.done)
	}
}
