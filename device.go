package picker

import (
	"encoding/binary"
	"errors"
	"log"
	"math"
	"reflect"
	"time"
)

type valueSet map[string]interface{}

type Sensor struct {
	Name    string
	Address []byte
	Values  valueSet
}

type sensors []Sensor

type Device struct {
	sensors
	port     *Port
	address  byte
	dtrReset bool
	logger   *log.Logger
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
	var res Buf
	var err error
	d.logger.Println(">>> ", request)
	attempts := 5
	for attempts > 0 {
		err = nil
		writeError := d.port.write(request.addCrc().slip())
		if writeError != nil {
			err = errors.New("Device: Write: " + writeError.Error())
			attempts--
			d.logger.Println(err, "attempt - ", attempts)
			continue
		}
		response, readError := d.port.read()
		if readError != nil {
			err = errors.New("Device: Read: " + readError.Error())
			attempts--
			d.logger.Println(err, "attempt - ", attempts)
			continue
		}
		unslipped, unSlipError := response.unSlip()
		if unSlipError != nil {
			err = errors.New("Device: " + unSlipError.Error())
		}
		if !unslipped.checkCrc() {
			err = errors.New("Device: CRC error")
			attempts--
			d.logger.Println(err, "attempt - ", attempts)
			continue
		}
		res = unslipped.removeCrc()
		break
	}
	d.logger.Println("<<< ", res)
	if err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

func (d *Device) ping() error {
	const pingCommand byte = 0x00
	rightAnswer := []byte{d.address, 0x55, 0xAA, 0x55, 0xAA}
	request := Buf{d.address, pingCommand}
	d.port.inUse.Lock()
	d.logger.Println("Ping")
	msg, commError := d.communicate(request)
	d.port.inUse.Unlock()
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
	const class byte = 0x01
	const method byte = 0x01
	request := Buf{d.address, class, method}
	d.port.inUse.Lock()
	d.logger.Println("Read DS1820")
	msg, commError := d.communicate(request)
	d.port.inUse.Unlock()
	if commError != nil {
		return commError
	}
	number := int(msg[1])
	for i := 0; i < number; i++ {
		tempBits := binary.LittleEndian.Uint32(msg[i*12+2 : i*12+6])
		temperature := math.Float32frombits(tempBits)
		sernum := msg[i*12+6 : i*12+14]
		values := valueSet{"temperature": temperature}
		isExist := d.updateIfExist(sernum, values, "ds18b20")
		if !isExist {
			newSensor := Sensor{Values: values, Address: sernum, Name: "ds18b20"}
			d.sensors = append(d.sensors, newSensor)
		}
	}
	return nil
}

func (d *Device) updateIfExist(sernum []byte, values valueSet, name string) bool {
	for _, sensor := range d.sensors {
		if reflect.DeepEqual(sensor.Address, sernum) && sensor.Name == name {
			for k := range values {
				sensor.Values[k] = values[k]
			}
			return true
		}
	}
	return false
}

func (d *Device) updateDHT22() error {
	const class byte = 0x01
	const method byte = 0x03
	request := Buf{d.address, class, method}
	d.port.inUse.Lock()
	d.logger.Println("Read DHT22")
	msg, commError := d.communicate(request)
	d.port.inUse.Unlock()
	if commError != nil {
		return commError
	}
	tempBits := binary.LittleEndian.Uint32(msg[1:5])
	humidity := math.Float32frombits(tempBits)
	tempBits = binary.LittleEndian.Uint32(msg[5:9])
	temperature := math.Float32frombits(tempBits)
	var sernum Buf
	sernum = append(sernum, msg[9])
	values := valueSet{"temperature": temperature, "humidity": humidity}
	isExist := d.updateIfExist(sernum, values, "dht22")
	if !isExist {
		newSensor := Sensor{Values: values, Address: sernum, Name: "dht22"}
		d.sensors = append(d.sensors, newSensor)
	}
	return nil
}

func (d *Device) updateBMP085Sensors() error {
	const class byte = 0x01
	const method byte = 0x02
	request := Buf{d.address, class, method}
	d.port.inUse.Lock()
	d.logger.Println("Read BMP085")
	msg, commError := d.communicate(request)
	d.port.inUse.Unlock()
	if commError != nil {
		return commError
	}
	pressure := binary.LittleEndian.Uint32(msg[1:5])
	tempBits := binary.LittleEndian.Uint32(msg[5:9])
	temperature := math.Float32frombits(tempBits)
	var sernum Buf
	sernum = append(sernum, msg[9])
	values := valueSet{"pressure": pressure, "temperature": temperature}
	isExist := d.updateIfExist(sernum, values, "bmp085")
	if !isExist {
		newSensor := Sensor{Values: values, Address: sernum, Name: "bmp085"}
		d.sensors = append(d.sensors, newSensor)
	}
	return nil
}

func (d *Device) updateAnalogInputs() error {
	const class byte = 0x01
	const method byte = 0x04
	request := Buf{d.address, class, method}
	d.port.inUse.Lock()
	d.logger.Println("Read analog inputs")
	msg, commError := d.communicate(request)
	d.port.inUse.Unlock()
	if commError != nil {
		return commError
	}
	number := int(msg[1])
	for i := 0; i < number; i++ {
		analogValue := binary.LittleEndian.Uint16(msg[i*3+3 : i*3+5])
		sernum := msg[i*3+2 : i*3+3]
		values := valueSet{"value": analogValue}
		isExist := d.updateIfExist(sernum, values, "analog")
		if !isExist {
			newSensor := Sensor{Values: values, Address: sernum, Name: "analog"}
			d.sensors = append(d.sensors, newSensor)
		}
	}
	return nil
}

func (d *Device) writeServo(index byte, value byte) error {
	const class byte = 2
	const method byte = 0
	request := Buf{d.address, class, method, index, value}
	d.port.inUse.Lock()
	d.logger.Println("Write", value, "to servo", index)
	msg, commError := d.communicate(request)
	d.port.inUse.Unlock()
	if commError != nil {
		return commError
	}
	if msg[0] == d.address && msg[1] == 0 {
		values := valueSet{"angle": value}
		sernum := []byte{index}
		isExist := d.updateIfExist(sernum, values, "servo")
		if !isExist {
			newSensor := Sensor{Values: values, Address: sernum, Name: "servo"}
			d.sensors = append(d.sensors, newSensor)
		}
		return nil
	} else {
		return errors.New("Servo error")
	}
}
