package main

import (
	"github.com/ndrwk/picker"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	ymlFile := flag.String("yml", "example.yml", "yaml config file name")
	makeFlag := flag.Bool("make_upload", false, "make source & upload flag")
	runFlag := flag.Bool("run", false, "run picker flag")
	flag.Parse()
	fmt.Println("yml:", *ymlFile)
	fmt.Println("make:", *makeFlag)
	fmt.Println("run:", *runFlag)

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
}
