package picker

type Actions interface {
	ReadName() string
	UpdateName(string)
	ReadValues() []float32
	UpdateValues([]float32)
	ReadAddr() []byte
	SetAddr([]byte)
}

type sensors []Actions

