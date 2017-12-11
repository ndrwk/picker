package picker

import (
	"fmt"
	"io/ioutil"
	"errors"
	"os/exec"
	"bytes"
	"log"
	"flag"
)

var device Device
var port Port
var config Env

func Init() error {
	ymlFile := flag.String("yml", "example.yml", "yaml config file name")
	makeFlag := flag.Bool("make_upload", false, "make source & upload flag")
	runFlag := flag.Bool("run", false, "run picker flag")
	flag.Parse()
	fmt.Println("yml:", *ymlFile)
	fmt.Println("make:", *makeFlag)
	fmt.Println("run:", *runFlag)
	data, err := ioutil.ReadFile(*ymlFile)
	if err != nil {
		return errors.New("Picker: Open file: " + err.Error())
	}
	config = Env{}
	configErr := config.configure([]byte(data))
	if configErr != nil {
		return configErr
	}
	fmt.Printf("%+v\n", config)
	if *runFlag {
		pickerError := create()
		if pickerError != nil{
			return pickerError
		}
		//defer destroy()
	}
	return nil
}

func makeFirmWare(yml []byte) error {
	dev, err := config.Device.makeHeader()
	if err != nil {
		return err
	}
	fmt.Println(dev)


	cmd := exec.Command("pwd")
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())

	return nil
}

func create() error {
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
