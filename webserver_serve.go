package go_webserver

import (
	"net"
	"sync/atomic"
)

// -----------------------------------------------------------------------------

const (
	stateNotStarted = 1
	stateStarting   = 2
	stateRunning    = 3
	stateStopping   = 4
	stateStopped    = 5
)

// -----------------------------------------------------------------------------

func (srv *Server) serve(ln net.Listener) {
	ch := make(chan error, 1)

	go func(ln net.Listener) {
		ch <- srv.fastserver.Serve(ln)
	}(ln)

	// Set new state
	srv.setState(stateRunning)

	// Run in background until shutdown or error
	go srv.serveLoop(ch)
}

func (srv *Server) serveLoop(ch chan error) {
	select {
	case err := <-ch:
		srv.setState(stateStopping)

		// Web server is no longer able to accept more connections
		if srv.listenErrorHandler != nil && err != nil {
			srv.listenErrorHandler(srv, err)
		}

	// handle termination signal
	case <-srv.startShutdownSignal:
		srv.setState(stateStopping)

		// Attempt the graceful shutdown by closing the listener
		// and completing all inflight requests.
		_ = srv.fastserver.Shutdown()
	}

	srv.setState(stateStopped)
}

func (srv *Server) setState(newState int32) {
	atomic.StoreInt32(&srv.state, newState)
}
