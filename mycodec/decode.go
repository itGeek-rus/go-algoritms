package mycodec

import "fmt"

// Decode restores the data compressed by Encode.
func Decode(src []byte) ([]byte, error) {
	if len(src)%2 != 0 {
		return nil, fmt.Errorf("mycodec: truncated rle stream: len=%d", len(src))
	}

	n := 0
	for i := 0; i < len(src); i += 2 {
		if src[i] == 0 {
			return nil, fmt.Errorf("mycodec: invalid rle count 0 at offset %d", i)
		}
		n += int(src[i])
	}

	out := make([]byte, 0, n)
	for i := 0; i < len(src); i += 2 {
		count := int(src[i])
		value := src[i+1]
		for j := 0; j < count; j++ {
			out = append(out, value)
		}
	}
	return out, nil
}
