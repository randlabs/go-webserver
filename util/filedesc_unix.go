// See the LICENSE file for license details.

//go:build unix

package util

import (
	"syscall"
)

// -----------------------------------------------------------------------------

func CheckMaxFileDescriptors(value uint64) bool {
	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	return err == nil && rLimit.Cur >= value
}
