package picker

import (
	"fmt"
)

var device Device
var port Port
var config Env

func MakeFirmWare(yml []byte) error {
	config = Env{}
	configErr := config.configure(yml)
	if configErr != nil {
		return configErr
	}
	fmt.Printf("%+v\n", config)
	dev, err := config.Device.makeHeader()
	if err != nil {
		return err
	}
	fmt.Println(dev)
	return nil
}

func Create(yml []byte) error {

	config = Env{}
	configErr := config.configure(yml)
	if configErr != nil {
		return configErr
	}

	port = Port{Name: config.Device.Port, Baud: config.Device.Baud, Timeout: config.Device.TimeOut}
	device = Device{address: config.Device.Address, port: &port, sensors: &sensors{}, dtrReset: config.Device.DTRReset}
	initDeviceError := device.init()
	if initDeviceError != nil {
		return initDeviceError
	}
	pingError := device.ping()
	if pingError != nil {
		return pingError
	}
	return nil
}

func Destroy() error {
	closeDeviceError := device.close()
	return closeDeviceError
}

func ReadSensors() error {
	for _, sensor := range config.Device.Sensors {
		switch sensor.Type {
		case "ds18b20":
			tempErr := device.updateDS1820Sensors()
			if tempErr != nil {
				return tempErr
			}
		case "bmp085":
			pressErr := device.updateBMP085Sensors()
			if pressErr != nil {
				return pressErr
			}
		}
	}
	return nil
}

func GetSensorsRef() *sensors {
	return device.sensors
}
