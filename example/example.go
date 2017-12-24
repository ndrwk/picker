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

	pickerError := picker.LoadConfig(*ymlFile)
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

		values := make(chan string)
		go picker.ReadSensors(values)
		fmt.Println(<-values)

		picker.Run(values)

		for i := 0; i < 500; i++ {
			//go picker.ReadSensors(values)
			sensors := <-values
			fmt.Println(sensors)
		}
	}
}
