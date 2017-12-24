package picker

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	"time"
)

type Device struct {
	*sensors
	port     *Port
	address  byte
	dtrReset bool
}

func (d Device) init() error {
	var portError error
	d.port.Serial, portError = d.port.openPort()
	if portError != nil {
		return errors.New("Device: Open port: " + portError.Error())
	}
	if d.dtrReset {
		resetTimeout := time.Millisecond * 1500
		time.Sleep(resetTimeout)
	}
	return nil
}

func (d Device) close() error {
	closeError := d.port.closePort()
	if closeError != nil {
		return errors.New("Device: Close port: " + closeError.Error())
	}
	return nil
}

func (d Device) communicate(request Buf) (Buf, error) {
	writeError := d.port.write(request.addCrc().slip())
	if writeError != nil {
		return nil, errors.New("Device: Write: " + writeError.Error())
	}
	response, readError := d.port.read()
	if readError != nil {
		return nil, errors.New("Device: Read: " + readError.Error())
	}
	unslipped, unSlipError := response.unSlip()
	if unSlipError != nil {
		return nil, errors.New("Device: " + unSlipError.Error())
	}
	if !unslipped.checkCrc() {
		return nil, errors.New(" Device: CRC error")
	}
	return unslipped.removeCrc(), nil
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

func (d Device) updateDS1820Sensors() error {
	const getCommand1 byte = 0x01
	const getCommand2 byte = 0x01
	request := Buf{d.address, getCommand1, getCommand2}
	msg, commError := d.communicate(request)
	if commError != nil {
		return commError
	}
	number := int(msg[1])
	for i := 0; i < number; i++ {
		tempBits := binary.LittleEndian.Uint32(msg[i*12+2 : i*12+6])
		temperature := math.Float32frombits(tempBits)
		sernum := msg[i*12+6 : i*12+14]
		isExist := updateIfExist(sernum, []float32{temperature}, &DS1820{})
		if !isExist {
			newTempSensor := DS1820{Value: temperature, Address: sernum, Name: "DS1820"}
			*d.sensors = append(*d.sensors, &newTempSensor)
		}
	}
	return nil
}

func updateIfExist(sernum []byte, values []float32, s Actions) bool {
	for _, sensor := range *device.sensors {
		if sensor != nil {
			addr := sensor.readAddr()
			if reflect.DeepEqual(addr, sernum) && reflect.TypeOf(s) == reflect.TypeOf(sensor) {
				sensor.updateValues(values)
				return true
			}
		}
	}
	return false
}

func (d Device) updateDHT22() error {
	const getCommand1 byte = 0x01
	const getCommand2 byte = 0x03
	request := Buf{d.address, getCommand1, getCommand2}
	msg, commError := d.communicate(request)
	if commError != nil {
		return commError
	}
	tempBits := binary.LittleEndian.Uint32(msg[1:5])
	humidity := math.Float32frombits(tempBits)
	tempBits = binary.LittleEndian.Uint32(msg[5:9])
	temperature := math.Float32frombits(tempBits)
	var sernum Buf
	sernum = append(sernum, msg[9])
	isExist := updateIfExist(sernum, []float32{float32(humidity), float32(temperature)}, &DHT22{})
	if !isExist {
		newPressureSensor := DHT22{Moisture: float32(humidity), Temperature: float32(temperature), Address: sernum, Name: "DHT22"}
		*d.sensors = append(*d.sensors, &newPressureSensor)
	}
	return nil
}

func (d Device) updateBMP085Sensors() error {
	const getCommand1 byte = 0x01
	const getCommand2 byte = 0x02
	request := Buf{d.address, getCommand1, getCommand2}
	msg, commError := d.communicate(request)
	if commError != nil {
		return commError
	}
	pressure := binary.LittleEndian.Uint32(msg[1:5])
	temperature := 100.0
	var sernum Buf
	sernum = append(sernum, msg[5])
	isExist := updateIfExist(sernum, []float32{float32(pressure), float32(temperature)}, &BMP085{})
	if !isExist {
		newPressureSensor := BMP085{Pressure: float32(pressure), Temperature: float32(temperature), Address: sernum, Name: "BMP085"}
		*d.sensors = append(*d.sensors, &newPressureSensor)
	}
	return nil
}
