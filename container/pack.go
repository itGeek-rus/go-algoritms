package container

var packAlgos = []Algorithm{AlgoRLE, AlgoHuffman, AlgoLZ77}

// Pack encodes src with the smallest GOAL container.
// If no codec shrinks the blob versus storing the original, AlgoNone is used.
func Pack(src []byte) []byte {
	best, _ := PackWith(AlgoNone, src)
	for _, algo := range packAlgos {
		blob, err := PackWith(algo, src)
		if err != nil {
			continue
		}
		if len(blob) < len(best) {
			best = blob
		}
	}
	return best
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
