package codec

import (
	"slices"
	"strings"
)

const base62Str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62Encoder(id int64, salt int64) string {

	// shift over 27 so that the id does not overlap with the salt
	combined := (id << 27) | salt

	return base62Encoding(combined)
}

// return the id and salt
func Base62Decoder(str string) (id int64, salt int64) {
	result := base62Decoding(str)

	if result == -1 {
		return -1, 0
	}

	salt = result & ((1 << 27) - 1)
	id = result >> 27

	return id, salt
}

func base62Encoding(value int64) string {
	if value == 0 {
		return "0"
	}

	compute := value
	var encodeStr []byte

	for compute != 0 {
		quotient := compute / 62
		remainder := compute % 62
		compute = quotient
		encodeStr = append(encodeStr, base62Str[remainder])
	}

	slices.Reverse(encodeStr)
	return string(encodeStr)
}

func base62Decoding(str string) int64 {
	result := int64(0)

	for idx := range str {
		pos := strings.IndexByte(base62Str, str[idx])
		if pos == -1 {
			return -1
		}
		result = (result * 62) + int64(pos)
	}

	return result
}
