package picker

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Message struct {
	TimeStamp     time.Time
	DeviceAddress byte
	Sensor
	Error         error
}

var device Device
var port Port
var config Env

func LoadConfig(ymlFile string) error {
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
	cmd0 += "platformio lib --global install 1336@None"        // DHTlib@None
	err := runCmd(cmd0)
	if err != nil {
		return err
	}
	hFile, err := config.Device.makeHeader()
	if err != nil {
		return err
	}
	hPath := filepath.Join(".picker", "src", "config.h")
	ioutil.WriteFile(hPath, []byte(hFile), 0777)
	templatePath := filepath.Join("..", "arduino", "device.cpp")
	source, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return err
	}
	cppPath := filepath.Join(".picker", "src", "device.cpp")
	ioutil.WriteFile(cppPath, source, 0777)
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
		fmt.Println(">>>>", line)
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
	device = Device{address: config.Device.Address, port: &port, sensors: sensors{}, dtrReset: config.Device.DTRReset}
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

func Run(valChan chan Message) {
	ticker := time.NewTicker(time.Second * time.Duration(config.Device.QueryPeriod))
	go func() {
		for range ticker.C {
			ReadSensors(valChan)
		}
	}()
}

func ReadSensors(valChan chan Message) error {
	for _, sensor := range config.Device.Sensors {
		sensorErr := updateSensor(sensor)
		if sensorErr != nil {
			valChan <- Message{Error: sensorErr}
			return sensorErr
		}
	}
	for _, s := range device.sensors {
			valChan <- Message{Sensor: s, Error: nil, TimeStamp: time.Now(), DeviceAddress: device.address}
	}
	return nil
}

func updateSensor(s SensorConfig) error {
	switch s.Type {
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
	return nil
}
