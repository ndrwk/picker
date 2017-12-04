package picker

import (
	"errors"
	"gopkg.in/yaml.v2"
	"strconv"
	"fmt"
)

type SensorConfig struct {
	Type string
	Pins string
}

type DeviceConfig struct {
	Address  byte
	Port     string
	Baud     int
	TimeOut  int
	DTRReset bool
	Sensors  []SensorConfig
}

type Env struct {
	Device DeviceConfig
}

func (config *Env) Configure(ymlData []byte) error {
	err := yaml.Unmarshal(ymlData, &config)
	if err != nil {
		return errors.New("Config: Yaml error: " + err.Error())
	}
	return nil
}

func (dc DeviceConfig) MakeConfigH() (string, error) {
	res := ""
	res += "#define ADDRESS " + fmt.Sprintf("%d\n", dc.Address)
	res += "#define BAUDRATE " + fmt.Sprintf("%d\n", dc.Baud)
	for _, v := range dc.Sensors {
		switch v.Type {
		case "ds18b20":
			if v.Pins == "" {
				return "", errors.New("Sensor config error: ds18b20")
			}
			res += "#define DS1820ENABLE\n"
			if _, err := strconv.Atoi(v.Pins[1:]); err != nil {
				return "", errors.New("Sensor config error: ds18b20: " + err.Error())
			}
			res += "#define DS1820_PIN " + v.Pins[1:] + "\n"
		case "bmp085":
			res += "#define BMP085ENABLE\n"
		}
	}
	return res, nil
}
