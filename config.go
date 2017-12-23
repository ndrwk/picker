package picker

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"strconv"
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
	QueryPeriod int
	Sensors  []SensorConfig
}

type Env struct {
	Device DeviceConfig
}

func (config *Env) configure(ymlData []byte) error {
	err := yaml.Unmarshal(ymlData, &config)
	if err != nil {
		return errors.New("Config: Yaml error: " + err.Error())
	}
	return nil
}

func (dc DeviceConfig) makeHeader() (string, error) {
	res := ""
	res += "#define ADDRESS " + fmt.Sprintf("%d\n", dc.Address)
	res += "#define BAUDRATE " + fmt.Sprintf("%d\n", dc.Baud)
	for _, v := range dc.Sensors {
		switch v.Type {
		case "ds18b20":
			if v.Pins == "" {
				return "", errors.New("Sensor config error: ds18b20 - empty Pins")
			}
			res += "#define DS1820ENABLE\n"
			if _, err := strconv.Atoi(v.Pins[1:]); err != nil {
				return "", errors.New("Sensor config error: ds18b20: " + err.Error())
			}
			res += "#define DS1820_PIN " + v.Pins[1:] + "\n"
		case "dht22":
			if v.Pins == "" {
				return "", errors.New("Sensor config error: dht22 - empty Pins")
			}
			res += "#define DHT22ENABLE\n"
			if _, err := strconv.Atoi(v.Pins[1:]); err != nil {
				return "", errors.New("Sensor config error: dht22: " + err.Error())
			}
			res += "#define DHT22_PIN " + v.Pins[1:] + "\n"
		case "bmp085":
			res += "#define BMP085ENABLE\n"
		}
	}
	return res, nil
}
