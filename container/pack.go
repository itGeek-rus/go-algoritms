package container

import "go-algoritms/mycodec"

// Pack wraps src in the GOAL container. Payload is always RLE for now.
func Pack(src []byte) []byte {
	payload := mycodec.Encode(src)
	out := writeHeader(AlgoRLE, uint64(len(src)))
	return append(out, payload...)
}
