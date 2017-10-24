package picker

var device Device
var port Port

func Create(portName string, portBaud int, portTimeout int, deviceAddress byte) error {
	port = Port{Name: portName, Baud: portBaud, Timeout: portTimeout}
	device = Device{address: deviceAddress, port: &port, sensors: &sensors{}}
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
	tempErr := device.updateTempSensors()
	if tempErr != nil {
		return tempErr
	}
	pressErr := device.updatePressureSensors()
	if pressErr != nil {
		return pressErr
	}
	return nil
}

func GetSensorsRef() *sensors{
	return device.sensors
}

