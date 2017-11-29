package picker

import "fmt"

type DS1820 struct {
	Name    string
	Value   float32
	Address []byte
}

type BMP085 struct {
	Name    string
	Value   float32
	Address []byte
}

type Communicator interface {
	ReadName() string
	UpdateName(string)
	ReadValue() float32
	UpdateValue(float32)
	ReadAddr() []byte
	SetAddr([]byte)
}

type sensors []Communicator

func (s DS1820) ReadValue() float32 {
	return s.Value
}

func (s DS1820) UpdateValue(value float32) {
	s.Value = value
}

func (s DS1820) ReadName() string {
	return s.Name
}

func (s DS1820) UpdateName(name string) {
	s.Name = name
}

func (s DS1820) ReadAddr() []byte {
	return s.Address
}

func (s DS1820) SetAddr(addr []byte) {
	s.Address = s.Address[:0]
	for _, v := range addr {
		s.Address = append(s.Address, v)
	}
}

func (s DS1820) String() string {
	var ser Buf = s.Address
	return "Temp: " + fmt.Sprintf("%.1f", s.Value) + " on " + ser.ToString()
}

func (s BMP085) String() string {
	var ser Buf = s.Address
	return "Pressure: " + fmt.Sprintf("%.1f", s.Value) + " on " + ser.ToString()
}

func (s BMP085) ReadValue() float32 {
	return s.Value
}

func (s BMP085) UpdateValue(value float32) {
	s.Value = value
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
