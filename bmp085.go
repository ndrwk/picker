package picker

import (
	"encoding/json"
	"fmt"
)

type BMP085 struct {
	Name        string
	Temperature float32
	Pressure    float32
	Address     []byte
}

func (s BMP085) toString() string {
	return "BMP085: Pressure = " + fmt.Sprintf("%d", int(s.Pressure)) + ", Temperature = " + fmt.Sprintf("%.2f", s.Temperature)
}

func (s BMP085) readValues() []float32 {
	return []float32{s.Pressure, s.Temperature}
}

func (s *BMP085) updateValues(values []float32) {
	s.Pressure = values[0]
	s.Temperature = values[1]
}

func (s BMP085) readAddr() []byte {
	return s.Address
}

func (s *BMP085) setAddr(addr []byte) {
	s.Address = s.Address[:0]
	for _, v := range addr {
		s.Address = append(s.Address, v)
	}
}

func (s BMP085) toJSON() ([]byte, error) {
	return json.Marshal(s)
}
