package huffman

import "encoding/binary"

func Encode(src []byte) []byte {
	if len(src) == 0 {
		return []byte{}
	}

	var freq [256]int
	for _, b := range src {
		freq[b]++
	}
	root := buildTree(freq)

	var codes [256]code
	assignCode(root, 0, 0, &codes)

	w := &bitWriter{}
	writeTree(w, root)
	if !root.leaf() {
		for _, b := range src {
			c := codes[b]
			for i := c.len - 1; i >= 0; i-- {
				w.writeBit(uint8((c.bits >> i) & 1))
			}
		}
	}

	out := make([]byte, 8, 8+len(w.buf)+1)
	binary.LittleEndian.PutUint64(out, uint64(len(src)))
	return append(out, w.bytes()...)
}
