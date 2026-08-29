package lz77

import (
	"encoding/binary"
	"fmt"
)

func Decode(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return []byte{}, nil
	}
	if len(src) < 8 {
		return nil, fmt.Errorf("lz77: truncated stream")
	}

	need := binary.LittleEndian.Uint64(src[:8])
	out := make([]byte, 0, need)
	i := 8

	readByte := func() (byte, error) {
		if i >= len(src) {
			return 0, fmt.Errorf("lz77: unexpected eof")
		}
		b := src[i]
		i++
		return b, nil
	}

	for uint64(len(out)) < need {
		flag, err := readByte()
		if err != nil {
			return nil, err
		}
		for b := 0; b < 8 && uint64(len(out)) < need; b++ {
			if flag&(1<<b) == 0 {
				lit, err := readByte()
				if err != nil {
					return nil, err
				}
				out = append(out, lit)
				continue
			}
			lo, err := readByte()
			if err != nil {
				return nil, err
			}
			hi, err := readByte()
			if err != nil {
				return nil, err
			}
			n, err := readByte()
			if err != nil {
				return nil, err
			}
			off := int(lo) | int(hi)<<8
			if off == 0 || off > len(out) || n < minMatch {
				return nil, fmt.Errorf("lz77: invalid match off=%d n=%d", off, n)
			}
			for k := 0; k < int(n); k++ {
				out = append(out, out[len(out)-off])
			}
		}
	}
	return out, nil
}
