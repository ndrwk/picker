package main

import (
	"flag"
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
	writeOut := flag.Bool("wr", false, "write to output")
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

		// go picker.RunServer()

		values := make(chan picker.Message, 1)
		picker.Run(values)
		//picker.RunAll(values, 5)
		//picker.RunOne(values, "ds18b20", 5)
		//picker.ReadSensor(values, "ds18b20")

		for res := range values {
			if res.Error != nil {
				log.Fatalf("error: %v", res.Error)
			}
			// fmt.Println(res.TimeStamp.String())
			// fmt.Println("Device:", res.DeviceAddress)
			// fmt.Printf("%+v\n", res.Sensor)
		}
	} else if *writeOut {
		pickerWrError := picker.Create()
		defer picker.Destroy()
		if pickerWrError != nil {
			log.Fatalf("error: %v", pickerError)
		}
		writeError := picker.WriteOutput("servo", 0, 126)
		if writeError != nil {
			log.Fatalf("error: %v", writeError)
		}

		// var a byte
		// for a = 0; a <= 125; a += 5 {
		// 	writeError := picker.WriteOutput("servo", 0, a)
		// 	if writeError != nil {
		// 		log.Fatalf("error: %v", writeError)
		// 	}
		// }
	}
}
