package picker

import "testing"

func TestEnv_Configure(t *testing.T) {
	var data = `
---
device:
  address: 0
  port: /dev/ttyUSB0
  baud: 115200
  timeout: 3000
  dtrreset: true
  sensors:
    - type: ds18b20
      pins: D10
    - type: bmp085
      pins: i2c
`
	yml := []byte(data)
	config := Env{}
	configErr := config.Configure(yml)
	if configErr != nil {
		t.Error("Expected nil, got", configErr)
	}
	if config.Device.Port != "/dev/ttyUSB0" || config.Device.Baud != 115200 || config.Device.TimeOut != 3000 {
		t.Error("Wrong unmarshalling")
	}
}
