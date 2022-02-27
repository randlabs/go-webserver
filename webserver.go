package go_webserver

import (
	"crypto/tls"
	"errors"
	"net"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/randlabs/go-webserver/listener"
	"github.com/valyala/fasthttp"
)

// -----------------------------------------------------------------------------

// RequestCtx ...
type RequestCtx = fasthttp.RequestCtx

// Router ...
type Router = fasthttprouter.Router

// Cookie ...
type Cookie = fasthttp.Cookie

// -----------------------------------------------------------------------------

type Server struct {
	fastserver             fasthttp.Server
	Router                 *fasthttprouter.Router
	bindAddress            net.IP
	bindPort               uint16
	listenErrorCallback    ListenErrorCallback
	state                  int32
	startShutdownSignal    chan struct{}
	shutdownCompleteSignal chan struct{}
}

type Options struct {
	// Address is the bind address to attach the server listener.
	Address string

	// Port is the port number the server will listen.
	Port uint16

	// ReadTimeout is the amount of time allowed to read
	// the full request including body. The connection's read
	// deadline is reset when the connection opens, or for
	// keep-alive connections after the first byte has been read.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset after the request handler
	// has returned.
	WriteTimeout time.Duration

	// Maximum number of concurrent client connections allowed per IP.
	MaxConnsPerIP int

	// Maximum number of requests served per connection.
	MaxRequestsPerConn int

	// Maximum request body size.
	MaxRequestBodySize int

	// Enable compression.
	EnableCompression bool

	// A callback to call if an error is encountered.
	ListenErrorCallback ListenErrorCallback

	// TLSConfig optionally provides a TLS configuration for use.
	TLSConfig *tls.Config
}

// ListenErrorCallback is a callback to call if an error is encountered.
type ListenErrorCallback func(srv *Server, err error)

// -----------------------------------------------------------------------------

const (
	DefaultReadTimeout        = 10 * time.Second
	DefaultWriteTimeout       = 10 * time.Second
	DefaultMaxRequestsPerConn = 8
	DefaultMaxRequestBodySize = 4 * 1048576 // 4MB

	stateNotStarted = 1
	stateStarting   = 2
	stateRunning    = 3
	stateStopping   = 4
	stateStopped    = 3
)

// -----------------------------------------------------------------------------

var strContentType = []byte("Content-Type")
var strApplicationJSON = []byte("application/json")

// -----------------------------------------------------------------------------

// Create creates a new webserver
func Create(options Options) (*Server, error) {
	// Check options
	if len(options.Address) == 0 {
		return nil, errors.New("invalid server bind address")
	}
	if options.Port < 1 || options.Port > 65535 {
		return nil, errors.New("invalid server port")
	}

	readTimeout := options.ReadTimeout
	if readTimeout < time.Duration(0) {
		return nil, errors.New("invalid read timeout")
	} else if readTimeout == time.Duration(0) {
		readTimeout = DefaultReadTimeout
	}

	writeTimeout := options.WriteTimeout
	if writeTimeout < time.Duration(0) {
		return nil, errors.New("invalid write timeout")
	} else if writeTimeout == time.Duration(0) {
		writeTimeout = DefaultWriteTimeout
	}

	maxConnsPerIP := options.MaxConnsPerIP
	if maxConnsPerIP < 0 {
		return nil, errors.New("invalid max connections per ip")
	}

	maxRequestsPerConn := options.MaxRequestsPerConn
	if maxRequestsPerConn < 0 {
		return nil, errors.New("invalid max requests per connections")
	} else if maxRequestsPerConn == 0 {
		maxRequestsPerConn = DefaultMaxRequestsPerConn
	}

	maxRequestBodySize := options.MaxRequestBodySize
	if maxRequestBodySize < 0 {
		return nil, errors.New("invalid max request body size")
	} else if maxRequestBodySize == 0 {
		maxRequestBodySize = DefaultMaxRequestBodySize
	}

	// Create a new server container
	srv := &Server{
		Router:                 fasthttprouter.New(),
		bindAddress:            net.ParseIP(options.Address),
		bindPort:               options.Port,
		listenErrorCallback:    options.ListenErrorCallback,
		state:                  stateNotStarted,
		startShutdownSignal:    make(chan struct{}, 1),
		shutdownCompleteSignal: make(chan struct{}, 1),
	}
	if srv.bindAddress == nil {
		return nil, errors.New("invalid server bind address")
	}
	if p4 := srv.bindAddress.To4(); len(p4) == net.IPv4len {
		srv.bindAddress = p4
	}

	// Setup compression
	h := srv.Router.Handler
	if options.EnableCompression {
		h = fasthttp.CompressHandler(h)
	}

	// Create FastHTTP server
	srv.fastserver = fasthttp.Server{
		Handler:            h,
		ReadTimeout:        readTimeout,
		WriteTimeout:       writeTimeout,
		MaxConnsPerIP:      maxConnsPerIP,
		MaxRequestsPerConn: maxRequestsPerConn,
		DisableKeepalive:   true,
		MaxRequestBodySize: maxRequestBodySize,
		TLSConfig:          options.TLSConfig,
		Logger:             newLoggerBridge(srv.logCallback),
	}

	// Done
	return srv, nil
}

// Start initiates listening
func (srv *Server) Start() error {
	if !atomic.CompareAndSwapInt32(&srv.state, stateNotStarted, stateStarting) {
		return errors.New("invalid state")
	}

	// Create the listener
	var network string

	// "tcp" network is not supported by all platforms
	address := srv.bindAddress.String()
	if len(srv.bindAddress) == net.IPv4len {
		network = "tcp4"
	} else {
		network = "tcp6"
		address = "[" + address + "]"
	}

	ln, err := createListener(network, address+":"+strconv.Itoa(int(srv.bindPort)))
	if err != nil {
		atomic.StoreInt32(&srv.state, stateNotStarted)
		return err
	}

	// Wrap listener into a TLS listener if a TLS configuration was specified
	if srv.fastserver.TLSConfig != nil {
		ln = tls.NewListener(ln, srv.fastserver.TLSConfig.Clone())
	}

	// Wrap listener inside a graceful shutdown listener
	ln = listener.NewGracefulListener(ln, 5*time.Second)

	// Start accepting connections
	errorChannel := make(chan error, 1)
	go func() {
		errorChannel <- srv.fastserver.Serve(ln)
	}()

	// Set new state
	atomic.StoreInt32(&srv.state, stateRunning)

	// Run in background until shutdown or error
	go func() {
		select {
		case errCh := <-errorChannel:
			atomic.StoreInt32(&srv.state, stateStopping)

			// Web server is no longer able to accept more connections
			if srv.listenErrorCallback != nil {
				srv.listenErrorCallback(srv, errCh)
			}

		// handle termination signal
		case <-srv.startShutdownSignal:
			atomic.StoreInt32(&srv.state, stateStopping)

			// Attempt the graceful shutdown by closing the listener
			// and completing all inflight requests.
			_ = srv.fastserver.Shutdown()
		}

		atomic.StoreInt32(&srv.state, stateStopped)
	}()

	// Done
	return nil
}

// Stop shuts down the web server
func (srv *Server) Stop() {
	if atomic.CompareAndSwapInt32(&srv.state, stateRunning, stateStopping) {
		srv.startShutdownSignal <- struct{}{}

		// Spin until server is really stopped
		for atomic.LoadInt32(&srv.state) == stateStopping {
			runtime.Gosched()
		}
	}
}

// -----------------------------------------------------------------------------
// Private methods

func (srv *Server) logCallback(format string, args ...interface{}) {
	// Nothing for now
}
