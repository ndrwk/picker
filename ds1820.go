package picker

import "fmt"

type DS1820 struct {
	Name    string
	Value   float32
	Address []byte
}

func (s DS1820) ReadValues() []float32 {
	return []float32{s.Value}
}

func (s DS1820) UpdateValues(values []float32) {
	s.Value = values[0]
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
	return "DS1820: " + "Temperature = " + fmt.Sprintf("%.2f", s.Value) + ", s/n " + ser.ToString()
}
