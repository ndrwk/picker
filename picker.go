package picker

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

var device Device
var port Port
var config Env

func Init(ymlFile string) error {
	data, err := ioutil.ReadFile(ymlFile)
	if err != nil {
		return errors.New("Picker: Open file: " + err.Error())
	}
	config = Env{}
	configErr := config.configure([]byte(data))
	if configErr != nil {
		return configErr
	}
	fmt.Printf("%+v\n", config)
	return nil
}

func MakeFirmWare() error {

	cmd := "platformio" + "\n"
	cmd += "../sources/create_env.sh" + "\n"
	cmd += "platformio run --target upload --project-dir .picker --upload-port " + config.Device.Port

	hFile, err := config.Device.makeHeader()
	if err != nil {
		return err
	}
	ioutil.WriteFile("../sources/arduino/config.h", []byte(hFile), 0666)

	cmds := strings.Split(cmd, "\n")
	for _, line := range cmds {
		fmt.Println(line)
		c := strings.Split(line, " ")
		create := exec.Command(c[0], c[1:]...)
		var out bytes.Buffer
		create.Stdout = &out
		err := create.Run()
		if err != nil {
			return errors.New(err.Error() + ":\n" + out.String())
		}
		fmt.Println(out.String())

	}
	return nil
}

func Create() error {
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
