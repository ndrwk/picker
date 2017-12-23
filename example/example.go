package main

import (
	"flag"
	"fmt"
	"github.com/ndrwk/picker"
	"log"
)

func main() {

	ymlFile := flag.String("yml", "example.yml", "yaml config file name")
	makeFlag := flag.Bool("make_upload", false, "make source & upload flag")
	runFlag := flag.Bool("run", false, "run picker flag")
	flag.Parse()

	pickerError := picker.Init(*ymlFile)
	if pickerError != nil {
		log.Fatalf("error: %v", pickerError)
	}

	if *makeFlag {
		makeError := picker.MakeFirmWare()
		if makeError != nil {
			log.Fatalf("error: %v", makeError)
		}
	} else if *runFlag {
		pickerError := picker.Create()
		defer picker.Destroy()
		if pickerError != nil {
			log.Fatalf("error: %v", pickerError)
		}

		var pickerSensors = picker.GetSensorsRef()

		//pickerError = picker.ReadSensors()
		//if pickerError != nil {
		//	log.Fatalf("error: %v", pickerError)
		//}
		//
		//for _, v := range *pickerSensors {
		//	fmt.Println("Имя ", v.ReadName())
		//	fmt.Println("Адрес ", v.ReadAddr())
		//	fmt.Println("Показание ", v.ReadValues())
		//}

		fmt.Println()
		for i := 0; i < 5; i++ {
			pickerError = picker.ReadSensors()
			if pickerError != nil {
				log.Fatalf("error: %v", pickerError)
			}
			for _, s := range *pickerSensors {
				fmt.Println(s)
			}
		}
	}
}
