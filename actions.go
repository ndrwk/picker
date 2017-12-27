package picker

type Actions interface {
	readValues() []float32
	updateValues([]float32)
	readAddr() []byte
	setAddr([]byte)
	toString() string
	toJSON() ([]byte, error)
}

type sensors []Actions
