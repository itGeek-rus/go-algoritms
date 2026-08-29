package container

import (
	"fmt"
)

func Unpack(src []byte) ([]byte, Header, error) {
	h, payload, err := parseHeader(src)
	if err != nil {
		return nil, Header{}, err
	}

	c, err := lookup(h.Algorithm)
	if err != nil {
		return nil, h, err
	}
	out, err := c.Decode(payload)
	if err != nil {
		return nil, h, err
	}
	if uint64(len(out)) != h.OrigSize {
		return nil, h, fmt.Errorf("container: size mismatch: header=%d got=%d", h.OrigSize, len(out))
	}
	return out, h, nil
}
