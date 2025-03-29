// See the LICENSE file for license details.

package middleware

// -----------------------------------------------------------------------------

func fastUint2Bytes(dst []byte, n uint64) int {
	var buf [24]byte

	ofs := 24
	q := uint64(0)
	for n >= 10 {
		q = n / 10
		ofs -= 1
		buf[ofs] = '0' + byte(n-q*10)
		n = q
	}
	ofs--
	buf[ofs] = '0' + byte(n)
	copy(dst, buf[ofs:])
	return 24 - ofs
}
