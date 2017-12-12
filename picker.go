package picker

import (
	//"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"bytes"
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

	check := exec.Command("platformio", "-h")
	err := check.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = check.Wait()
	if err != nil {
		log.Fatal(err)
	}

	mkdir := exec.Command("mkdir", ".picker")
	err = mkdir.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = mkdir.Wait()
	if err != nil {
		log.Fatal(err)
	}

	create := exec.Command("platformio", "init", "--board", "nanoatmega328", "--project-dir", ".picker")
	var out bytes.Buffer
	create.Stdout = &out
	err = create.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())


	//cmd := exec.Command("platformio", "-h")
	////cmd.Stdin = strings.NewReader("some input")
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//err := cmd.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(out.String())
	//dev, err := config.Device.makeHeader()
	//if err != nil {
	//	return err
	//}
	//fmt.Println(dev)

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
