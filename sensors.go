package picker

type Actions interface {
	ReadName() string
	UpdateName(string)
	ReadValue() float32
	//ReadValues() []float32
	UpdateValue(float32)
	//UpdateValues([]float32)
	ReadAddr() []byte
	SetAddr([]byte)
}

type sensors []Actions

