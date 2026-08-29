package huffman

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncodeDecodeRoundtrip(t *testing.T) {
	cases := [][]byte{
		{},
		{0},
		{1, 2, 3, 4, 5},
		bytes.Repeat([]byte{'A'}, 10),
		[]byte("hello hello hello"),
		{1, 1, 2, 2, 2, 3},
	}
	for i, in := range cases {
		got, err := Decode(Encode(in))
		if err != nil {
			t.Fatalf("case %d: decode: %v", i, err)
		}
		if !bytes.Equal(got, in) {
			t.Fatalf("case %d: roundtrip mismatch", i)
		}
	}
}

func TestEncodeDecodeRandom(t *testing.T) {
	in := make([]byte, 4096)
	if _, err := rand.Read(in); err != nil {
		t.Fatal(err)
	}
	got, err := Decode(Encode(in))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, in) {
		t.Fatal("random roundtrip mismatch")
	}
}
