package picker

import "fmt"

type DHT22 struct {
	Name     string
	Temperature float32
	Moisture float32
	Address  []byte
}

func (s DHT22) toString() string {
	return "DHT22: Moisture = " + fmt.Sprintf("%d", int(s.Moisture)) + ", Temperature = " + fmt.Sprintf("%.2f", s.Temperature)
}

func (s DHT22) readValues() []float32 {
	return []float32{s.Moisture, s.Temperature}
}

func (s *DHT22) updateValues(values []float32) {
	s.Moisture = values[0]
	s.Temperature = values[1]
}

func (s DHT22) readName() string {
	return s.Name
}

func (s *DHT22) updateName(name string) {
	s.Name = name
}

func (s DHT22) readAddr() []byte {
	return s.Address
}

func (s *DHT22) setAddr(addr []byte) {
	s.Address = s.Address[:0]
	for _, v := range addr {
		s.Address = append(s.Address, v)
	}
}

