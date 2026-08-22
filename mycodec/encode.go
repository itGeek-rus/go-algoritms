package mycodec

// Encode compresses the source data using the RLE algorithm.
// Format: a sequence of [count][value] pairs, where count ranges from 1 to 255.
func Encode(src []byte) []byte {
	if len(src) == 0 {
		return []byte{}
	}

	out := make([]byte, 0, len(src))
	count := 1
	prev := src[0]

	flush := func() {
		for count > 0 {
			n := count
			if n > 255 {
				n = 255
			}
			out = append(out, byte(n), prev)
			count -= n
		}
	}

	for i := 1; i < len(src); i++ {
		if src[i] == prev && count < 255 {
			count++
			continue
		}
		flush()
		prev = src[i]
		count = 1
	}
	flush()

	return out
}
