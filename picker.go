package picker

import "fmt"

var device Device
var port Port

func Create(yml []byte) error {

	env := Env{}
	config := env.Configure(yml)
	if config != nil {
		return config
	}
	fmt.Printf("%+v\n", env)

	port = Port{Name: env.Devices[0].Port, Baud: env.Devices[0].Baud, Timeout: env.Devices[0].TimeOut}
	device = Device{address: env.Devices[0].Address, port: &port, sensors: &sensors{}, dtrReset: env.Devices[0].DTRReset}
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

func UpdateSensors() error {
	tempErr := device.updateDS1820Sensors()
	if tempErr != nil {
		return tempErr
	}
	pressErr := device.updateBMP085Sensors()
	if pressErr != nil {
		return pressErr
	}
	return nil
}

func GetSensorsRef() *sensors{
	return device.sensors
}

