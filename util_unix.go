//go:build unix

package go_webserver

import (
	"syscall"
)

// -----------------------------------------------------------------------------

func checkMaxFileDescriptors(value uint64) bool {
	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	return err == nil && rLimit.Cur >= value
}
