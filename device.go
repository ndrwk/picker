package picker

import (
	"encoding/binary"
	"errors"
	"math"
	"time"
)

type Device struct {
	*sensors
	port    *Port
	address byte
}

func (d Device) init() error {
	var portError error
	d.port.Serial, portError = d.port.Open()
	if portError != nil {
		return errors.New("Device: Open port: " + portError.Error())
	}
	resetTimeout := time.Millisecond * 1500
	time.Sleep(resetTimeout)
	return nil
}

func (d Device) close() error {
	closeError := d.port.Close()
	if closeError != nil {
		errors.New("Device: Close port: " + closeError.Error())
	}
	return nil
}

func (d Device) communicate(request Buf) (Buf, error) {
	writeError := d.port.Write(request.AddCrc().Slip())
	if writeError != nil {
		return nil, errors.New("Device: Write: " + writeError.Error())
	}
	response, readError := d.port.Read()
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
	const pingCommand byte = 0x00
	rightAnswer := []byte{d.address, 0x55, 0xAA, 0x55, 0xAA}
	request := Buf{d.address, pingCommand}
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
	const tempCommand1 byte = 0x01
	const tempCommand2 byte = 0x01
	request := Buf{d.address, tempCommand1, tempCommand2}
	msg, commError := d.communicate(request)
	if commError != nil {
		return commError
	}
	number := int(msg[1])
	for i := 0; i < number; i++ {
		tempBits := binary.LittleEndian.Uint32(msg[i*12+2 : i*12+6])
		temperature := math.Float32frombits(tempBits)
		sernum := msg[i*12+6 : i*12+14]
		//isExist := false
		isExist := updateIfNotExist(sernum, temperature)
		if !isExist {
			newTempSensor := TempSensor{Value: temperature, Address: sernum}
			*d.sensors = append(*d.sensors, newTempSensor)
		}
	}
	return nil
}
func updateIfNotExist(sernum Buf, temperature float32) bool {
	isExist := false
	for _, v := range *device.sensors {
		isExist = true
		if v != nil {
			switch v.(type) {
			case TempSensor:
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
	}
	return isExist
}

func (d Device) updatePressureSensors() error {
	const getPressureCommand1 byte = 0x01
	const getPressureCommand2 byte = 0x02
	request := Buf{d.address, getPressureCommand1, getPressureCommand2}
	msg, commError := d.communicate(request)
	if commError != nil {
		return commError
	}
	//msg.PrintHex()
	pressure := binary.LittleEndian.Uint32(msg[1:5])
	//fmt.Println(pressure)
	var sernum Buf
	sernum = append(sernum, msg[5])
	//fmt.Println(sernum)
	isExist := updateIfNotExist(sernum, float32(pressure))
	if !isExist {
		newPressureSensor := PressureSensor{Value: float32(pressure), Address: sernum}
		*d.sensors = append(*d.sensors, newPressureSensor)
	}
	return nil
}
