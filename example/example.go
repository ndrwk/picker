package main

import (
	"fmt"
	"github.com/ndrwk/picker"
	//"io/ioutil"
	"log"
)

var portName = "/dev/ttyUSB0"

//var portBaud = 9600
var portBaud = 115200
var portTimeout = 3000

var deviceAddress byte = 0

func main() {

	//data, err := ioutil.ReadFile("devices.yml")
	//if err != nil {
	//	log.Fatalf("error: %v", err)
	//}
	var data = `
---
devices:
  - address: 0
    port: /dev/ttyUSB0
    baud: 115200
    timeout: 3000
    dtrreset: true
    sensors:
      - type: ds18b20
        pins: D10
  - address: 1
    port: /dev/ttyUSB1
    baud: 9600
    timeout: 3000
    dtrreset: false
    sensors:
      - type: ds18b20
        pins: D10
      - type: bmp085
        pins: i2c
`
	yml := []byte(data)

	var pickerError error
	pickerError = picker.Create(yml)
	if pickerError != nil {
		log.Fatalf("error: %v", pickerError)
	}
	defer picker.Destroy()

	pickerError = picker.UpdateSensors()
	if pickerError != nil {
		log.Fatalf("error: %v", pickerError)
	}

	var pickerSensors = picker.GetSensorsRef()

	for _, v := range *pickerSensors {
		switch v.(type) {
		case picker.DS1820:
			fmt.Println("Температурный")
		case picker.BMP085:
			fmt.Println("Давление")
		}
		fmt.Println("Имя ", v.ReadName())
		fmt.Println("Адрес ", v.ReadAddr())
		fmt.Println("Показание ", v.ReadValue())
		fmt.Println()
	}

	for _, s := range *pickerSensors {
		fmt.Println(s)
	}

}
