package util

import (
	"unsafe"
)

// -----------------------------------------------------------------------------

func UnsafeByteSlice2String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func UnsafeString2ByteSlice(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
