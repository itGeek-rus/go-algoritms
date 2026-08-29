package huffman

import "fmt"

type bitWriter struct {
	buf []byte
	acc byte
	n   int
}

func (w *bitWriter) writeBit(b uint8) {
	w.acc = w.acc<<1 | (b & 1)
	w.n++
	if w.n == 8 {
		w.buf = append(w.buf, w.acc)
		w.acc = 0
		w.n = 0
	}
}

func (w *bitWriter) writeByte(b byte) {
	for i := 7; i >= 0; i-- {
		w.writeBit((b >> i) & 1)
	}
}

func (w *bitWriter) bytes() []byte {
	if w.n > 0 {
		w.acc <<= 8 - w.n
		w.buf = append(w.buf, w.acc)
		w.acc = 0
		w.n = 0
	}
	return w.buf
}

type bitReader struct {
	buf []byte
	i   int
	acc byte
	n   int
}

func (r *bitReader) readBit() (uint8, error) {
	if r.n == 0 {
		if r.i >= len(r.buf) {
			return 0, fmt.Errorf("huffman: unexpected eof")
		}
		r.acc = r.buf[r.i]
		r.i++
		r.n = 8
	}
	r.n--
	return (r.acc >> r.n) & 1, nil
}

func (r *bitReader) readByte() (byte, error) {
	var b byte
	for i := 0; i < 8; i++ {
		bit, err := r.readBit()
		if err != nil {
			return 0, err
		}
		b = b<<1 | bit
	}
	return b, nil
}
