package picker

import "fmt"

type DHT22 struct {
	Name     string
	Temperature float32
	Moisture float32
	Address  []byte
}

func (s DHT22) String() string {
	return "DHT22: Moisture = " + fmt.Sprintf("%d", int(s.Moisture)) + ", Temperature = " + fmt.Sprintf("%.2f", s.Temperature)
}

func (s DHT22) ReadValues() []float32 {
	return []float32{s.Moisture, s.Temperature}
}

func (s *DHT22) UpdateValues(values []float32) {
	s.Moisture = values[0]
	s.Temperature = values[1]
}

func (s DHT22) ReadName() string {
	return s.Name
}

func (s *DHT22) UpdateName(name string) {
	s.Name = name
}

func (s DHT22) ReadAddr() []byte {
	return s.Address
}

func (s *DHT22) SetAddr(addr []byte) {
	s.Address = s.Address[:0]
	for _, v := range addr {
		s.Address = append(s.Address, v)
	}
}

