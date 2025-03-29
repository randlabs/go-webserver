// See the LICENSE file for license details.

//go:build !unix

package util

// -----------------------------------------------------------------------------

func CheckMaxFileDescriptors(_ uint64) bool {
	return true
}
