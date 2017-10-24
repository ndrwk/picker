package main

import (
	"github.com/ndrwk/picker"
	"fmt"
	"os"
)

var portName = "/dev/ttyUSB0"
var portBaud = 9600
var portTimeout = 3000

var deviceAddress byte = 0

var pickerError error

func main() {

	pickerError = picker.Create(portName, portBaud, portTimeout, deviceAddress)
	if pickerError != nil {
		fmt.Println(pickerError)
		os.Exit(1)
	}
	defer picker.Destroy()

	picker.UpdateSensors()
	picker.PrintSensors()

}

