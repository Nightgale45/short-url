package codec

import (
	"slices"
)

const base62Str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62Encoder(id int64, salt int64) string {

	// shift over 27 so that the id does not overlap with the salt
	combined := (id << 27) | salt

	return base62Encoding(combined)
}

func Base62Decoder(str string) string {
	return ""
}

func base62Encoding(value int64) string {
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
