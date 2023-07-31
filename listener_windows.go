//go:build windows

package go_webserver

import (
	"net"
)

// -----------------------------------------------------------------------------

func createListener(network string, address string) (net.Listener, error) {
	return net.Listen(network, address)
}
