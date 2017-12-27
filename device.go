package picker

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	"time"
)

type Sensor struct {
	Name    string
	Address []byte
	Values  map[string]float32
}

type sensors []Sensor

type Device struct {
	sensors
	port     *Port
	address  byte
	dtrReset bool
}

func (d *Device) init() error {
	portError := d.port.openPort()
	if portError != nil {
		return errors.New("Device: Open port: " + portError.Error())
	}
	if d.dtrReset {
		resetTimeout := time.Millisecond * 1500
		time.Sleep(resetTimeout)
	}
	return nil
}

func (d *Device) close() error {
	closeError := d.port.closePort()
	if closeError != nil {
		return errors.New("Device: Close port: " + closeError.Error())
	}
	return nil
}

func (d *Device) communicate(request Buf) (Buf, error) {
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

func (d *Device) ping() error {
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

func (d *Device) updateDS1820Sensors() error {
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
		values := make(map[string]float32)
		values["temperature"] = temperature
		isExist := updateIfExist(sernum, values, "DS1820")
		if !isExist {
			newSensor := Sensor{Values: values, Address: sernum, Name: "DS1820"}
			d.sensors = append(d.sensors, newSensor)
		}
	}
	return nil
}

func updateIfExist(sernum []byte, values map[string]float32, name string) bool {
	for _, sensor := range device.sensors {
		if reflect.DeepEqual(sensor.Address, sernum) && sensor.Name == name {
			for k := range values{
				sensor.Values[k] = values[k]
			}
			return true
		}
	}
	return false
}

func (d *Device) updateDHT22() error {
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
	values := make(map[string]float32)
	values["temperature"] = temperature
	values["humidity"] = humidity
	isExist := updateIfExist(sernum, values, "DHT22")
	if !isExist {
		newSensor := Sensor{Values: values, Address: sernum, Name: "DHT22"}
		d.sensors = append(d.sensors, newSensor)
	}
	return nil
}

func (d *Device) updateBMP085Sensors() error {
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
	values := make(map[string]float32)
	values["pressure"] = float32(pressure)
	values["temperature"] = float32(temperature)
	isExist := updateIfExist(sernum, values, "BMP085")
	if !isExist {
		newSensor := Sensor{Values: values, Address: sernum, Name: "BMP085"}
		d.sensors = append(d.sensors, newSensor)
	}
	return nil
}
