package main

import (
	"github.com/ndrwk/picker"
	"fmt"
	"os"
)

var portName = "/dev/ttyUSB0"
//var portBaud = 9600
var portBaud = 115200
var portTimeout = 3000

var deviceAddress byte = 0

func main() {

	var pickerError error
	pickerError = picker.Create(portName, portBaud, portTimeout, deviceAddress)
	if pickerError != nil {
		fmt.Println(pickerError)
		os.Exit(1)
	}
	defer picker.Destroy()

	pickerError = picker.UpdateSensors()
	if pickerError != nil {
		fmt.Println(pickerError)
		os.Exit(1)
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

	for _, s := range *pickerSensors{
		fmt.Println(s)
	}

}

