package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ndrwk/picker"
)

func main() {

	// usage1 - example -yml /home/drew/go/src/github.com/ndrwk/picker/example/example.yml -source_path /home/drew/go/src/github.com/ndrwk/picker -make_upload
	// usage2 - example -yml /home/drew/go/src/github.com/ndrwk/picker/example/example.yml -run

	ymlFile := flag.String("yml", "example.yml", "yaml config file name")
	makeFlag := flag.Bool("make_upload", false, "make source & upload flag")
	runFlag := flag.Bool("run", false, "run picker flag")
	sourcePath := flag.String("source_path", "..", "source path")
	flag.Parse()

	pickerError := picker.LoadConfig(*ymlFile)
	if pickerError != nil {
		log.Fatalf("error: %v", pickerError)
	}

	if *makeFlag {
		makeError := picker.MakeFirmWare(*sourcePath)
		if makeError != nil {
			log.Fatalf("error: %v", makeError)
		}
	} else if *runFlag {
		pickerError := picker.Create()
		defer picker.Destroy()
		if pickerError != nil {
			log.Fatalf("error: %v", pickerError)
		}

		values := make(chan picker.Message, 1)
		picker.Run(values)
		//picker.RunAll(values, 5)
		//picker.RunOne(values, "ds18b20", 5)
		//picker.ReadSensor(values, "ds18b20")

		for res := range values {
			if res.Error != nil {
				log.Fatalf("error: %v", res.Error)
			}
			fmt.Println(res.TimeStamp.String())
			fmt.Println("Device:", res.DeviceAddress)
			fmt.Printf("%+v\n", res.Sensor)
		}
	}
}
