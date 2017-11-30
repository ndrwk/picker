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
	Devices []DeviceConfig
}

func (env *Env) Configure(ymlData []byte) error {
	err := yaml.Unmarshal(ymlData, &env)
	if err != nil {
		return errors.New("Config: Yaml error: " + err.Error())
	}
	return nil
}
