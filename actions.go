package picker

type Actions interface {
	readName() string
	updateName(string)
	readValues() []float32
	updateValues([]float32)
	readAddr() []byte
	setAddr([]byte)
	toString() string
}

type sensors []Actions

