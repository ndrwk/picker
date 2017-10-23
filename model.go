package picker

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"
)

type TempSensor struct {
	Name    string
	Value   float32
	Address []byte
	Device  *Device
}

type PressureSensor struct {
	Name   string
	Value  float32
	Device *Device
}

type Communicator interface {
	ReadName() string
	UpdateName(string)
	ReadValue() float32
	UpdateValue(float32)
	ReadAddr() []byte
	SetAddr([]byte)
}

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

type Sensors []Communicator

type Device struct {
	*Sensors
	Port    *Port
	Address byte
}

func (d Device) init() error {
	var portError error
	d.Port.Serial, portError = d.Port.Open()
	if portError != nil {
		return errors.New("Device: Open port: " + portError.Error())
	}
	resetTimeout := time.Millisecond * 1500
	time.Sleep(resetTimeout)
	return nil
}

func (d Device) close() error {
	closeError := d.Port.Close()
	if closeError != nil {
		errors.New("Device: Close port: " + closeError.Error())
	}
	return nil
}

func (d Device) communicate(request Buf) (Buf, error) {
	writeError := d.Port.Write(request.AddCrc().Slip())
	if writeError != nil {
		return nil, errors.New("Device: Write: " + writeError.Error())
	}
	response, readError := d.Port.Read()
	if readError != nil {
		return nil, errors.New("Device: Read: " + readError.Error())
	}
	unslipped, unSlipError := response.UnSlip()
	if unSlipError != nil {
		return nil, errors.New("Device: " + unSlipError.Error())
	}
	if !unslipped.CheckCrc() {
		return nil, errors.New(" Device: CRC error")
	}
	return unslipped.RemoveCrc(), nil
}

func (d Device) ping() error {
	const PING_COMMAND byte = 0x00
	rightAnswer := []byte{d.Address, 0x55, 0xAA, 0x55, 0xAA}
	request := Buf{d.Address, PING_COMMAND}
	msg, commError := d.communicate(request)
	if commError != nil {
		return commError
	}
	if len(rightAnswer) != len(msg) {
		return errors.New("Device: too short answer")
	}
	for i := range rightAnswer {
		if rightAnswer[i] != msg[i] {
			return errors.New("Device: the answer doesn't match")
		}
	}
	return nil
}

func (d Device) updateTempSensors() error {
	const GET_TEMP_COMMAND_1 byte = 0x01
	const GET_TEMP_COMMAND_2 byte = 0x01
	request := Buf{d.Address, GET_TEMP_COMMAND_1, GET_TEMP_COMMAND_2}
	msg, commError := d.communicate(request)
	if commError != nil {
		return commError
	}
	number := int(msg[1])
	for i := 0; i < number; i++ {
		tempBits := binary.LittleEndian.Uint32(msg[i*12+2 : i*12+6])
		temperature := math.Float32frombits(tempBits)
		sernum := msg[i*12+6 : i*12+14]
		isExist := false
		for _, v := range *device.Sensors {
			isExist = true
			if v != nil {
				for i := range v.ReadAddr() {
					if v.ReadAddr()[i] != sernum[i] {
						isExist = false
						break
					}
				}
				if isExist {
					v.UpdateValue(temperature)
				}
			}
		}
		if !isExist {
			newTempSensor := TempSensor{Value: temperature, Address: sernum}
			*d.Sensors = append(*d.Sensors, newTempSensor)
		}
	}
	return nil
}

func (d Device) updatePressureSensor() error {
	const GET_PRESSURE_COMMAND_1 byte = 0x01
	const GET_PRESSURE_COMMAND_2 byte = 0x02
	request := Buf{d.Address, GET_PRESSURE_COMMAND_1, GET_PRESSURE_COMMAND_2}
	msg, commError := d.communicate(request)
	if commError != nil {
		return commError
	}
	msg.PrintHex()
	pressure := binary.LittleEndian.Uint32(msg[1:5])
	fmt.Println(pressure)
	sernum := msg[5]
	fmt.Println(sernum)
	return nil
}
