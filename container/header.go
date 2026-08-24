package container

import (
	"encoding/binary"
	"fmt"
)

const (
	Version    uint8 = 1
	headerSize       = 14
)

var Magic = [4]byte{'G', 'O', 'A', 'L'}

type Algorithm uint8

const (
	AlgoNone Algorithm = 0
	AlgoRLE  Algorithm = 1
)

type Header struct {
	Version   uint8
	Algorithm Algorithm
	OrigSize  uint64
}

func writeHeader(algo Algorithm, origSize uint64) []byte {
	buf := make([]byte, headerSize)
	copy(buf[0:4], Magic[:])
	buf[4] = Version
	buf[5] = byte(algo)
	binary.LittleEndian.PutUint64(buf[6:14], origSize)
	return buf
}

func parseHeader(src []byte) (Header, []byte, error) {
	if len(src) < headerSize {
		return Header{}, nil, fmt.Errorf("container: truncated header: len=%d", len(src))
	}
	if string(src[0:4]) != string(Magic[:]) {
		return Header{}, nil, fmt.Errorf("container: bad magic")
	}

	h := Header{
		Version:   src[4],
		Algorithm: Algorithm(src[5]),
		OrigSize:  binary.LittleEndian.Uint64(src[6:14]),
	}
	if h.Version != Version {
		return Header{}, nil, fmt.Errorf("container: unsupported version %d", h.Version)
	}

	return h, src[headerSize:], nil
}
