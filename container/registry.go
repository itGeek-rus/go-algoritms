package container

import (
	"fmt"
	"go-algoritms/codec"
	"go-algoritms/huffman"
	"go-algoritms/lz77"
	"go-algoritms/mycodec"
)

func lookup(algo Algorithm) (codec.Codec, error) {
	switch algo {
	case AlgoNone:
		return codec.Func{Enc: storeEncode, Dec: storeDecode}, nil
	case AlgoRLE:
		return codec.Func{Enc: mycodec.Encode, Dec: mycodec.Decode}, nil
	case AlgoHuffman:
		return codec.Func{Enc: huffman.Encode, Dec: huffman.Decode}, nil
	case AlgoLZ77:
		return codec.Func{Enc: lz77.Encode, Dec: lz77.Decode}, nil
	default:
		return nil, fmt.Errorf("container: unknown algorithm %d", algo)
	}
}

func storeEncode(src []byte) []byte {
	out := make([]byte, len(src))
	copy(out, src)
	return out
}

func storeDecode(src []byte) ([]byte, error) {
	out := make([]byte, len(src))
	copy(out, src)
	return out, nil
}
