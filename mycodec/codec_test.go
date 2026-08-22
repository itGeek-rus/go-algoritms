package mycodec

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
		bytes.Repeat([]byte{0xFF}, 300), // серия длиннее 255
		{1, 1, 2, 2, 2, 3},
	}
	for i, in := range cases {
		got, err := Decode(Encode(in))
		if err != nil {
			t.Fatalf("case %d: decode: %v", i, err)
		}
		if !bytes.Equal(got, in) {
			t.Fatalf("case %d: roundtrip mismatch\nwant %q\ngot  %q", i, in, got)
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
func TestEncodeShrinksRuns(t *testing.T) {
	in := bytes.Repeat([]byte{'Z'}, 100)
	out := Encode(in)
	if len(out) >= len(in) {
		t.Fatalf("expected shrink, in=%d out=%d", len(in), len(out))
	}
}
func TestDecodeRejectsOddLength(t *testing.T) {
	if _, err := Decode([]byte{1}); err == nil {
		t.Fatal("expected error for truncated stream")
	}
}
func TestDecodeRejectsZeroCount(t *testing.T) {
	if _, err := Decode([]byte{0, 'A'}); err == nil {
		t.Fatal("expected error for count=0")
	}
}
