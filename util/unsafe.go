// See the LICENSE file for license details.

package util

import (
	"unsafe"
)

// -----------------------------------------------------------------------------

func UnsafeByteSlice2String(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func UnsafeString2ByteSlice(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
