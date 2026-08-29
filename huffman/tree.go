package huffman

import "container/heap"

type node struct {
	freq        int
	sym         int
	left, right *node
}

func (n *node) leaf() bool {
	return n.left == nil && n.right == nil
}

type nodeHeap []*node

func (h nodeHeap) Len() int {
	return len(h)
}

func (h nodeHeap) Less(i, j int) bool {
	if h[i].freq == h[j].freq {
		return h[i].sym < h[j].sym
	}
	return h[i].freq < h[j].freq
}

func (h nodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *nodeHeap) Push(x any) {
	*h = append(*h, x.(*node))
}

func (h *nodeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func buildTree(freq [256]int) *node {
	h := nodeHeap{}
	for sym, n := range freq {
		if n > 0 {
			h = append(h, &node{freq: n, sym: sym})
		}
	}
	heap.Init(&h)
	if len(h) == 1 {
		return h[0]
	}
	for len(h) > 1 {
		a := heap.Pop(&h).(*node)
		b := heap.Pop(&h).(*node)
		heap.Push(&h, &node{freq: a.freq + b.freq, sym: -1, left: a, right: b})
	}
	return heap.Pop(&h).(*node)
}

type code struct {
	bits uint64
	len  int
}

func assignCode(n *node, bits uint64, length int, codes *[256]code) {
	if n.leaf() {
		if length == 0 {
			length = 1
		}
		codes[n.sym] = code{bits: bits, len: length}
		return
	}
	assignCode(n.left, bits<<1, length+1, codes)
	assignCode(n.right, bits<<1|1, length+1, codes)
}

func writeTree(w *bitWriter, n *node) {
	if n.leaf() {
		w.writeBit(1)
		w.writeByte(byte(n.sym))
		return
	}
	w.writeBit(0)
	writeTree(w, n.left)
	writeTree(w, n.right)
}

func readTree(r *bitReader) (*node, error) {
	bit, err := r.readBit()
	if err != nil {
		return nil, err
	}
	if bit == 1 {
		sym, err := r.readByte()
		if err != nil {
			return nil, err
		}
		return &node{sym: int(sym)}, nil
	}
	left, err := readTree(r)
	if err != nil {
		return nil, err
	}
	right, err := readTree(r)
	if err != nil {
		return nil, err
	}
	return &node{sym: -1, left: left, right: right}, nil
}
