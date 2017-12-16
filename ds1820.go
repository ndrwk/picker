package picker

import "fmt"

type DS1820 struct {
	Name    string
	Value   float32
	Address []byte
}

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
