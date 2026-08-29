package lz77

import "encoding/binary"

const (
	windowSize = 4096
	minMatch   = 3
	maxMatch   = 255
)

type token struct {
	match bool
	lit   byte
	off   uint16
	n     uint8
}

func Encode(src []byte) []byte {
	if len(src) == 0 {
		return []byte{}
	}

	var tokens []token
	for i := 0; i < len(src); {
		off, n := findMatch(src, i)
		if n >= minMatch {
			tokens = append(tokens, token{match: true, off: off, n: n})
			i += int(n)
			continue
		}
		tokens = append(tokens, token{lit: src[i]})
		i++
	}

	out := make([]byte, 8, 8+len(src)+1)
	binary.LittleEndian.PutUint64(out, uint64(len(src)))

	for i := 0; i < len(tokens); i += 8 {
		var flag byte
		var chunk []byte
		n := 8
		if i+n > len(tokens) {
			n = len(tokens) - i
		}
		for b := 0; b < n; b++ {
			t := tokens[i+b]
			if t.match {
				flag |= 1 << b
				chunk = append(chunk, byte(t.off), byte(t.off>>8), t.n)
			} else {
				chunk = append(chunk, t.lit)
			}
		}
		out = append(out, flag)
		out = append(out, chunk...)
	}
	return out
}

func findMatch(src []byte, i int) (uint16, uint8) {
	start := i - windowSize
	if start < 0 {
		start = 0
	}
	bestN, bestOff := 0, 0
	for j := start; j < i; j++ {
		n := 0
		for i+n < len(src) && n < maxMatch && src[j+n] == src[i+n] {
			n++
		}
		if n > bestN {
			bestN = n
			bestOff = i - j
		}
	}
	if bestN < minMatch {
		return 0, 0
	}
	return uint16(bestOff), uint8(bestN)
}
