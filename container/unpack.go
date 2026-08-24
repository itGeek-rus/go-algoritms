package container

import (
	"fmt"
	"go-algoritms/mycodec"
)

func Unpack(src []byte) ([]byte, Header, error) {
	h, payload, err := parseHeader(src)
	if err != nil {
		return nil, Header{}, err
	}

	out, err := decodePayload(h.Algorithm, payload)
	if err != nil {
		return nil, h, err
	}
	if uint64(len(out)) != h.OrigSize {
		return nil, h, fmt.Errorf("container: size mismatch: header=%d got=%d", h.OrigSize, len(out))
	}
	return out, h, nil
}

func decodePayload(algo Algorithm, payload []byte) ([]byte, error) {
	switch algo {
	case AlgoRLE:
		return mycodec.Decode(payload)
	case AlgoNone:
		out := make([]byte, len(payload))
		copy(out, payload)
		return out, nil
	default:
		return nil, fmt.Errorf("container: unknown algorithm %d", algo)
	}
}
