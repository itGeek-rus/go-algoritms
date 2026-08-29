package huffman

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Decode(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return []byte{}, nil
	}
	if len(src) < 8 {
		return nil, fmt.Errorf("huffman: truncated stream")
	}

	n := binary.LittleEndian.Uint64(src[:8])
	r := &bitReader{buf: src[8:]}
	root, err := readTree(r)
	if err != nil {
		return nil, err
	}

	if root.leaf() {
		return bytes.Repeat([]byte{byte(root.sym)}, int(n)), nil
	}

	out := make([]byte, 0, int(n))
	for uint64(len(out)) < n {
		cur := root
		for !cur.leaf() {
			bit, err := r.readBit()
			if err != nil {
				return nil, err
			}
			if bit == 0 {
				cur = cur.left
			} else {
				cur = cur.right
			}
		}
		out = append(out, byte(cur.sym))
	}
	return out, nil
}
