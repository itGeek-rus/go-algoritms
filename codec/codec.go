package codec

type Codec interface {
	Encode(src []byte) []byte
	Decode(src []byte) ([]byte, error)
}

type Func struct {
	Enc func([]byte) []byte
	Dec func([]byte) ([]byte, error)
}

func (f Func) Encode(src []byte) []byte {
	return f.Enc(src)
}

func (f Func) Decode(src []byte) ([]byte, error) {
	return f.Dec(src)
}
