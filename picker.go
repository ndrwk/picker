package picker

import (
	"fmt"
)

var device Device
var port Port

func Create(portName string, portBaud int, portTimeout int, deviceAddress byte) error {
	port = Port{Name: portName, Baud: portBaud, Timeout: portTimeout}
	device = Device{Address: deviceAddress, Port: &port, Sensors: &Sensors{}}
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

func UpdateSensors() {
	tempErr := device.updateTempSensors()
	if tempErr != nil {
		fmt.Println(tempErr)
	}
	pressErr := device.updatePressureSensors()
	if pressErr != nil {
		fmt.Println(pressErr)
	}
}

func PrintSensors() {
	for _, v := range *device.Sensors {
		switch v.(type) {
		case TempSensor:
			fmt.Println("Температурный")
		case PressureSensor:
			fmt.Println("Давление")
		}
		fmt.Println("Имя ", v.ReadName())
		fmt.Println("Адрес ", v.ReadAddr())
		fmt.Println("Показание ", v.ReadValue())

	}
}
