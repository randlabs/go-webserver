//go:build !unix

package go_webserver

// -----------------------------------------------------------------------------

func checkMaxFileDescriptors(_ uint64) bool {
	return true
}
