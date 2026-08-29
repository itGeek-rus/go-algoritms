# go-algoritms

Lossless compression library for Go 1.27. It wraps data in a small GOAL container, tries several codecs, and keeps the original payload when compression does not shrink the result.

There is no universal lossless algorithm that makes every input smaller. Already compressed or random data is stored as-is (`AlgoNone`). Unpack always reads the codec id from the header.

## Install

```bash
go get go-algoritms
```

## Usage

```go
package main

import (
	"fmt"

	"go-algoritms/container"
)

func main() {
	blob := container.Pack([]byte("hello hello hello"))
	out, h, err := container.Unpack(blob)
	if err != nil {
		panic(err)
	}
	fmt.Println(h.Algorithm, len(blob), string(out))
}
```

- `Pack` — try RLE, Huffman, LZ77; pick the smallest GOAL blob; if none is strictly smaller than store, use `AlgoNone`.
- `PackWith(algo, src)` — force one codec, even when the blob grows. For tests and debugging.
- `Unpack` — read the header and decode with that codec.

## Container format

```
offset  size  field
0       4     magic "GOAL"
4       1     version (1)
5       1     algorithm id
6       8     original size, little-endian uint64
14      …     payload
```

| id | Name | Package |
|----|------|---------|
| 0 | store (original bytes) | `AlgoNone` |
| 1 | RLE | `mycodec` |
| 2 | Huffman | `huffman` |
| 3 | LZ77 | `lz77` |

The output is always a GOAL blob. For incompressible input it is 14 bytes larger than the raw source because of the header; that is what makes `Unpack` possible.

## Packages

| Package | Role |
|---------|------|
| `container` | GOAL format, `Pack` / `PackWith` / `Unpack` |
| `codec` | `Encode` / `Decode` interface |
| `mycodec` | run-length encoding |
| `huffman` | Huffman coding |
| `lz77` | LZ77 matches |

## Test

```bash
go test ./...
```
