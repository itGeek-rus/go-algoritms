package container

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"testing"
)

func TestPackUnpackRoundtrip(t *testing.T) {
	cases := [][]byte{
		{},
		{0},
		{1, 2, 3, 4, 5},
		bytes.Repeat([]byte{'A'}, 10),
		bytes.Repeat([]byte{0xFF}, 300),
		{1, 1, 2, 2, 2, 3},
	}
	for i, in := range cases {
		got, h, err := Unpack(Pack(in))
		if err != nil {
			t.Fatalf("case %d: unpack: %v", i, err)
		}
		if h.Algorithm != AlgoRLE {
			t.Fatalf("case %d: algo=%d, want RLE", i, h.Algorithm)
		}
		if h.OrigSize != uint64(len(in)) {
			t.Fatalf("case %d: orig size=%d, want %d", i, h.OrigSize, len(in))
		}
		if !bytes.Equal(got, in) {
			t.Fatalf("case %d: roundtrip mismatch", i)
		}
	}
}
func TestPackUnpackRandom(t *testing.T) {
	in := make([]byte, 4096)
	if _, err := rand.Read(in); err != nil {
		t.Fatal(err)
	}
	got, _, err := Unpack(Pack(in))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, in) {
		t.Fatal("random roundtrip mismatch")
	}
}
func TestPackWritesHeader(t *testing.T) {
	blob := Pack([]byte("aaa"))
	if len(blob) < headerSize {
		t.Fatal("blob shorter than header")
	}
	if string(blob[0:4]) != "GOAL" {
		t.Fatalf("magic=%q", blob[0:4])
	}
	if blob[4] != Version {
		t.Fatalf("version=%d", blob[4])
	}
	if Algorithm(blob[5]) != AlgoRLE {
		t.Fatalf("algo=%d", blob[5])
	}
	if binary.LittleEndian.Uint64(blob[6:14]) != 3 {
		t.Fatal("orig size != 3")
	}
}
func TestUnpackAlgoNone(t *testing.T) {
	raw := []byte("hello")
	blob := writeHeader(AlgoNone, uint64(len(raw)))
	blob = append(blob, raw...)
	got, h, err := Unpack(blob)
	if err != nil {
		t.Fatal(err)
	}
	if h.Algorithm != AlgoNone {
		t.Fatalf("algo=%d", h.Algorithm)
	}
	if !bytes.Equal(got, raw) {
		t.Fatal("store payload mismatch")
	}
}
func TestUnpackRejectsShort(t *testing.T) {
	if _, _, err := Unpack([]byte("GO")); err == nil {
		t.Fatal("expected truncated header error")
	}
}
func TestUnpackRejectsBadMagic(t *testing.T) {
	blob := Pack([]byte("x"))
	blob[0] = 'X'
	if _, _, err := Unpack(blob); err == nil {
		t.Fatal("expected bad magic error")
	}
}
func TestUnpackRejectsUnknownAlgo(t *testing.T) {
	blob := Pack([]byte("x"))
	blob[5] = 99
	if _, _, err := Unpack(blob); err == nil {
		t.Fatal("expected unknown algorithm error")
	}
}
func TestUnpackRejectsSizeMismatch(t *testing.T) {
	blob := Pack([]byte("aaa"))
	binary.LittleEndian.PutUint64(blob[6:14], 99)
	if _, _, err := Unpack(blob); err == nil {
		t.Fatal("expected size mismatch error")
	}
}

func TestPackWithRoundtrip(t *testing.T) {
	algos := []Algorithm{AlgoNone, AlgoRLE, AlgoHuffman, AlgoLZ77}
	cases := [][]byte{
		{},
		{0},
		[]byte("abcabcabc AAA hello hello"),
		bytes.Repeat([]byte{'Z'}, 100),
	}
	for _, algo := range algos {
		for i, in := range cases {
			blob, err := PackWith(algo, in)
			if err != nil {
				t.Fatalf("algo=%d case %d: pack: %v", algo, i, err)
			}
			got, h, err := Unpack(blob)
			if err != nil {
				t.Fatalf("algo=%d case %d: unpack: %v", algo, i, err)
			}
			if h.Algorithm != algo {
				t.Fatalf("algo=%d case %d: header algo=%d", algo, i, h.Algorithm)
			}
			if !bytes.Equal(got, in) {
				t.Fatalf("algo=%d case %d: roundtrip mismatch", algo, i)
			}
		}
	}
}
func TestPackWithUnknown(t *testing.T) {
	if _, err := PackWith(99, []byte("x")); err == nil {
		t.Fatal("expected unknown algorithm error")
	}
}
