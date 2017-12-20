package picker

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
	return nil
}

func MakeFirmWare() error {
	os.Mkdir(".picker", 0777)
	cmd0 := "platformio init --board nanoatmega328 --project-dir .picker" + "\n"
	cmd0 += "platformio lib --global install 525@1.0.0" + "\n" // Adafruit BMP085 Library @ 1.0.0
	cmd0 += "platformio lib --global install 54@3.7.7" + "\n"  // OneWire @ 2.3.2
	cmd0 += "platformio lib --global install 1336@None" // DHTlib@None
	err := runCmd(cmd0)
	if err != nil {
		return err
	}
	hFile, err := config.Device.makeHeader()
	if err != nil {
		return err
	}
	ioutil.WriteFile(".picker/src/config.h", []byte(hFile), 0777)
	source, err := ioutil.ReadFile("../arduino/device.cpp")
	if err != nil {
		return err
	}
	ioutil.WriteFile(".picker/src/device.cpp", source, 0777)
	cmd1 := "platformio run --target upload --project-dir .picker --upload-port " + config.Device.Port
	err = runCmd(cmd1)
	if err != nil {
		return err
	}
	return nil
}

func runCmd(cmdStr string) error {
	cmds := strings.Split(cmdStr, "\n")
	for _, line := range cmds {
		fmt.Println(">>>>>", line)
		c := strings.Split(line, " ")
		createCmd := exec.Command(c[0], c[1:]...)
		var out bytes.Buffer
		createCmd.Stdout = &out
		createCmd.Stderr = &out
		err := createCmd.Run()
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
		case "dht22":
			err := device.updateDHT22()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetSensorsRef() *sensors {
	return device.sensors
}
