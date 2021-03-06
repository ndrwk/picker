package picker

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	Error error
}

type Task struct {
	Ticker *time.Ticker
	Sensor SensorConfig
}

var tasks []Task
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

func MakeFirmWare(sourcePath string) error {
	os.Mkdir(".picker", 0777)
	cmd0 := "platformio init --board nanoatmega328 --project-dir .picker" + "\n"
	cmd0 += "platformio lib --global install 525@1.0.0" + "\n" // Adafruit BMP085 Library @ 1.0.0
	cmd0 += "platformio lib --global install 54@3.7.7" + "\n"  // OneWire @ 2.3.2
	cmd0 += "platformio lib --global install 1336"             // DHTlib
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
	templatePath := filepath.Join(sourcePath, "arduino", "device.cpp")
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
	logContext := fmt.Sprint("#", config.Device.Address, ": ")
	var logFileName string
	var logFile io.Writer
	var logErr error
	if config.Device.Log != "" {
		if config.Device.Log == "off" {
			logFileName = os.DevNull
		} else {
			logFileName = config.Device.Log
		}
		logFile, logErr = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if logErr != nil {
			logFile = os.Stdout
		}
	} else {
		logFile = os.Stdout
	}
	logger := log.New(logFile, logContext, log.Ldate|log.Ltime)
	port = Port{name: config.Device.Port, baud: config.Device.Baud, timeout: config.Device.TimeOut}
	device = Device{address: config.Device.Address, port: &port, sensors: sensors{}, dtrReset: config.Device.DTRReset, logger: logger}
	initDeviceError := device.init()
	if initDeviceError != nil {
		return initDeviceError
	}
	pingError := device.ping()
	if pingError != nil {
		return pingError
	}
	go runWSServer(config.Device.Ip, logger)
	return nil
}

func Destroy() error {
	closeDeviceError := device.close()
	return closeDeviceError
}

func Run(valChan chan Message) {
	for _, s := range config.Device.Sensors {
		if s.Period != 0 {
			tasks = append(tasks, Task{Ticker: time.NewTicker(time.Second * time.Duration(s.Period)), Sensor: s})
		}
	}
	for _, t := range tasks {
		go doTask(t, valChan)
	}
}

func doTask(t Task, valChan chan Message) {
	for range t.Ticker.C {
		ReadSensor(valChan, t.Sensor.Type)
	}
}

//func RunAll(valChan chan Message, period int) {
//	ticker := time.NewTicker(time.Second * time.Duration(period))
//	go func() {
//		for range ticker.C {
//			ReadAllSensors(valChan)
//		}
//	}()
//}

func RunOne(valChan chan Message, sensorType string, period int) {
	ticker := time.NewTicker(time.Second * time.Duration(period))
	go func() {
		for range ticker.C {
			ReadSensor(valChan, sensorType)
		}
	}()
}

func ReadAllSensors(valChan chan Message) error {
	for _, sensor := range config.Device.Sensors {
		sensorErr := updateSensor(sensor.Type)
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

func ReadSensor(valChan chan Message, sensorType string) error {
	sensorErr := updateSensor(sensorType)
	if sensorErr != nil {
		valChan <- Message{Error: sensorErr}
		return sensorErr
	}
	for _, s := range device.sensors {
		if s.Name == sensorType {
			valChan <- Message{Sensor: s, Error: nil, TimeStamp: time.Now(), DeviceAddress: device.address}
		}
	}
	return nil
}

func updateSensor(sensorType string) error {
	var err error
	switch sensorType {
	case "ds18b20":
		err = device.updateDS1820Sensors()
	case "bmp085":
		err = device.updateBMP085Sensors()
	case "dht22":
		err = device.updateDHT22()
	case "analog":
		err = device.updateAnalogInputs()
	}
	return err
}

func WriteOutput(sensorType string, index byte, value byte) error {
	var err error
	switch sensorType {
	case "servo":
		err = device.writeServo(index, value)
	}
	return err
}
