package picker

import (
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v2"
)

type SensorConfig struct {
	Type   string
	Pins   []string
	Period int
}

type DeviceConfig struct {
	Ip       string
	Address  byte
	Port     string
	Baud     int
	TimeOut  int
	DTRReset bool
	Log      string
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
			if v.Pins == nil {
				return "", errors.New("Sensor config error: ds18b20 - empty Pins")
			}
			res += "#define DS1820ENABLE\n"
			if _, err := strconv.Atoi(v.Pins[0][1:]); err != nil {
				return "", errors.New("Sensor config error: ds18b20: " + err.Error())
			}
			res += "#define DS1820_PIN " + v.Pins[0][1:] + "\n"
		case "dht22":
			if v.Pins == nil {
				return "", errors.New("Sensor config error: dht22 - empty Pins")
			}
			res += "#define DHT22ENABLE\n"
			if _, err := strconv.Atoi(v.Pins[0][1:]); err != nil {
				return "", errors.New("Sensor config error: dht22: " + err.Error())
			}
			res += "#define DHT22_PIN " + v.Pins[0][1:] + "\n"
		case "bmp085":
			res += "#define BMP085ENABLE\n"
		case "analog":
			if v.Pins == nil {
				return "", errors.New("Analog inputs config error: Analog - empty Pins")
			}
			res += "#define ANALOGREADENABLE\n"
			res += "unsigned char analog_pins[] = {"
			strPins := ""
			for _, pin := range v.Pins {
				if _, err := strconv.Atoi(pin[1:]); err != nil {
					return "", errors.New("Analog input config error: Analog: " + err.Error())
				}
				strPins += pin[1:] + ","
			}
			res += strPins[:len(strPins)-1]
			res += "};\n"
		case "servo":
			if v.Pins == nil {
				return "", errors.New("Servo config error: Servo - empty Pins")
			}
			res += "#define SERVOENABLE\n"
			res += "unsigned char servo_pins[] = {"
			strPins := ""
			for _, pin := range v.Pins {
				if _, err := strconv.Atoi(pin[1:]); err != nil {
					return "", errors.New("Servo config error: Servo: " + err.Error())
				}
				strPins += pin[1:] + ","
			}
			res += strPins[:len(strPins)-1]
			res += "};\n"
			res += "const int servo_numbers = " + fmt.Sprint(len(v.Pins)) + ";\n"

		}
	}
	return res, nil
}
