package picker

import (
	"time"
	"errors"
	"github.com/tarm/serial"
)

type Port struct {
	Name    string
	Baud    int
	Timeout int
	Serial  *serial.Port
}

func (p Port) openPort() (*serial.Port, error) {
	return serial.OpenPort(&serial.Config{Name: p.Name, Baud: p.Baud, ReadTimeout: time.Duration(p.Timeout) * time.Millisecond})
}

func (p Port) closePort() error {
	return p.Serial.Close()
}

func (p Port) write(b Buf) error {
	n, err := p.Serial.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return errors.New("Device: wrong byte count on Write()")
	}
	return nil
}

func (p Port) read() (Buf, error) {
	response := Buf{}
	tmpBuf := make([]byte, 1)
	packetStarted := false
	var err error
	for err == nil {
		_, err = p.Serial.Read(tmpBuf)
		//fmt.Printf("%02X", tmpBuf[0])
		//fmt.Print(" ")
		if tmpBuf[0] == 0xC0 {
			response = append(response, tmpBuf[0])
			if packetStarted {
				break
			} else {
				packetStarted = true
			}

		} else {
			if packetStarted {
				response = append(response, tmpBuf[0])
			}
		}
	}
	//fmt.Println()
	return response, err
}
