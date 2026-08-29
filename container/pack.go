package container

// Pack wraps src in the GOAL container. Payload is always RLE for now.
func Pack(src []byte) []byte {
	blob, _ := PackWith(AlgoRLE, src)
	return blob
}

func PackWith(algo Algorithm, src []byte) ([]byte, error) {
	c, err := lookup(algo)
	if err != nil {
		return nil, err
	}
	payload := c.Encode(src)
	out := writeHeader(algo, uint64(len(src)))
	return append(out, payload...), nil
}
