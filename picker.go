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

func M() {
	//for _, v := range device.Sensors {
	//	fmt.Println(v.ReadName(), v.ReadValue())
	//}

	tempErr := device.updateTempSensors()
	if tempErr != nil {
		fmt.Println(tempErr)
	}
	pressErr := device.updatePressureSensor()
	if pressErr != nil {
		fmt.Println(pressErr)
	}

}
