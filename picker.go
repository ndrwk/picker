package picker
//
//import (
//	"fmt"
//	"os"
//)
//
//var sensors = Sensors{
//	TempSensor{
//		Name:    "Температура",
//		Value:   11.22,
//		Address: []byte{1, 2, 3, 4, 5, 6},
//		//Device:
//	},
//	PressureSensor{
//		Name:  "Давление",
//		Value: 760.0,
//	},
//}
////var port = protocol.Port{Name: "/dev/ttyUSB0", Baud: 115200, Timeout: 3000}
//var port = Port{Name: "/dev/ttyUSB0", Baud: 9600, Timeout: 3000}
//var device = Device{Sensors: sensors, Address: 0x00, Port: &port}
//
//func main() {
//	//for _, v := range device.Sensors {
//	//	fmt.Println(v.ReadName(), v.ReadValue())
//	//}
//	initDeviceError := device.Init()
//	if initDeviceError != nil {
//		fmt.Println(initDeviceError)
//		os.Exit(1)
//	}
//
//	pingRes, err := device.Ping()
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(pingRes)
//
//	tempErr := device.UpdateTempSensors()
//	if tempErr != nil {
//		fmt.Println(tempErr)
//	}
//
//	pressErr := device.UpdatePressureSensor()
//	if pressErr != nil {
//		fmt.Println(pressErr)
//	}
//
//	closeDeviceError := device.Close()
//	if closeDeviceError != nil {
//		fmt.Println(closeDeviceError)
//	}
//
//}
