package utils

import "bytes"

var baseChars = []byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'm', 'n', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'2', '3', '4', '5', '6', '7', '8', '9',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

var BaseLen = uint64(56)

// singleToDec converts a single base character to the equivalent decimal value
func singleToDec(c byte) uint64 {

	// illegal characters
	if c == 'l' || c == 'o' || c == '1' || c == '0' || c == 'I' || c == 'O' {
		return ^uint64(0) // return the max value of this type
	}

	switch {
	case c >= 'a' && c <= 'k':
		return uint64(c - 'a') // 11 chars
	case c >= 'm' && c <= 'n':
		return uint64(c - 'm' + 11) // 13 chars
	case c >= 'p' && c <= 'z':
		return uint64(c - 'p' + 13) // 24 chars
	case c >= '2' && c <= '9':
		return uint64(c - '2' + 24) // 32 chars
	case c >= 'A' && c <= 'H':
		return uint64(c - 'A' + 32) // 40 chars
	case c >= 'J' && c <= 'N':
		return uint64(c - 'J' + 40) // 45 chars
	case c >= 'P' && c <= 'Z':
		return uint64(c - 'P' + 45) // 56 chars
	}

	return ^uint64(0)
}

// ToSafeBase converts a decimal number (base 10) to the safe base representation
func ToSafeBase(n uint64) string {
	if n < BaseLen {
		return string(baseChars[n])
	}

	var buff bytes.Buffer

	for n != 0 {
		buff.WriteByte(baseChars[n%BaseLen])
		n /= BaseLen
	}

	return buff.String()
}

// ToDec converts a safe base stirng to an equivalent decimal (base 10) value
func ToDec(safe string) uint64 {
	if len(safe) < 1 {
		return ^uint64(0)
	}

	var n, mul uint64 = 0, 1
	for i := range safe {
		n += singleToDec(safe[i]) * mul
		mul *= BaseLen
	}
	return n
}
