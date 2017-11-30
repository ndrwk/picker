package picker

var device Device
var port Port

func Create(yml []byte) error {

	config := Env{}
	configErr := config.Configure(yml)
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

func GetSensorsRef() *sensors {
	return device.sensors
}
