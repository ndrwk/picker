package picker

import "fmt"

type BMP085 struct {
	Name     string
	Pressure float32
	Address  []byte
}

func (s BMP085) String() string {
	var ser Buf = s.Address
	return "Pressure: " + fmt.Sprintf("%.1f", s.Pressure) + " on " + ser.ToString()
}

//func (s BMP085) ReadValue() float32 {
//	return s.Pressure
//}

func (s BMP085) ReadValues() []float32 {
	return []float32{s.Pressure}
}

func (s BMP085) UpdateValue(value float32) {
	s.Pressure = value
}

func (s BMP085) ReadName() string {
	return s.Name
}

func (s BMP085) UpdateName(name string) {
	s.Name = name
}

func (s BMP085) ReadAddr() []byte {
	return s.Address
}

func (s BMP085) SetAddr(addr []byte) {
	s.Address = s.Address[:0]
	for _, v := range addr {
		s.Address = append(s.Address, v)
	}
}

