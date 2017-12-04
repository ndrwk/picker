package main

import (
	"github.com/ndrwk/picker"
	"log"
	"fmt"
)

func main() {

	//data, err := ioutil.ReadFile("devices.yml")
	//if err != nil {
	//	log.Fatalf("error: %v", err)
	//}
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

	config := picker.Env{}
	configErr := config.Configure(yml)
	if configErr != nil {
		log.Fatal(configErr)
	}

	fmt.Printf("%+v\n", config)

	dev, err := config.Device.MakeConfigH()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dev)


}
