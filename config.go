package picker

import (
	"errors"
	"gopkg.in/yaml.v2"
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
