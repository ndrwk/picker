package picker

import "fmt"

type TempSensor struct {
	Name    string
	Value   float32
	Address []byte
}

type PressureSensor struct {
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

func (s TempSensor) ReadValue() float32 {
	return s.Value
}

func (s TempSensor) UpdateValue(value float32) {
	s.Value = value
}

func (s TempSensor) ReadName() string {
	return s.Name
}

func (s TempSensor) UpdateName(name string) {
	s.Name = name
}

func (s TempSensor) ReadAddr() []byte {
	return s.Address
}

func (s TempSensor) SetAddr(addr []byte) {
	s.Address = s.Address[:0]
	for _, v := range addr {
		s.Address = append(s.Address, v)
	}
}

func (s TempSensor) String() string {
	var ser Buf = s.Address
	return "Temp: " + fmt.Sprintf("%.1f", s.Value) + " on " + ser.ToString()
}

func (s PressureSensor) String() string {
	var ser Buf = s.Address
	return "Pressure: " + fmt.Sprintf("%.1f", s.Value) + " on " + ser.ToString()
}

func (s PressureSensor) ReadValue() float32 {
	return s.Value
}

func (s PressureSensor) UpdateValue(value float32) {
	s.Value = value
}

func (s PressureSensor) ReadName() string {
	return s.Name
}

func (s PressureSensor) UpdateName(name string) {
	s.Name = name
}

func (s PressureSensor) ReadAddr() []byte {
	return s.Address
}

func (s PressureSensor) SetAddr(addr []byte) {
	s.Address = s.Address[:0]
	for _, v := range addr {
		s.Address = append(s.Address, v)
	}
}
