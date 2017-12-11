package main

import (
	"fmt"
	"github.com/ndrwk/picker"
	//"io/ioutil"
	"log"
)

func main() {

	pickerError := picker.Init()
	if pickerError != nil {
		log.Fatalf("error: %v", pickerError)
	}
	defer picker.Destroy()

	var pickerSensors = picker.GetSensorsRef()

	pickerError = picker.ReadSensors()
	if pickerError != nil {
		log.Fatalf("error: %v", pickerError)
	}



	for _, v := range *pickerSensors {
		fmt.Println("Имя ", v.ReadName())
		fmt.Println("Адрес ", v.ReadAddr())
		fmt.Println("Показание ", v.ReadValue())
	}

	for _, s := range *pickerSensors {
		fmt.Println(s)
	}

}
