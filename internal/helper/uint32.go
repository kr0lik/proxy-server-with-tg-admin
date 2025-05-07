package helper

import "fmt"

const uint32MaxDigits = 10

func Uint32ToString(n uint32) string {
	var buf [uint32MaxDigits]byte
	i := len(buf)

	for n >= uint32MaxDigits {
		i--
		buf[i] = byte('0' + n%uint32MaxDigits)
		n /= uint32MaxDigits
	}
	i--
	buf[i] = byte('0' + n)

	// return slice without zeros
	return string(buf[i:])
}

func StringToUint32(s string) (uint32, error) {
	var n uint32
	for i := range s {
		d := s[i] - '0'
		if d > (uint32MaxDigits - 1) {
			return 0, fmt.Errorf("StringToUint32: invalid digit %q in string", s[i])
		}
		n = n*uint32MaxDigits + uint32(d)
	}

	return n, nil
}
