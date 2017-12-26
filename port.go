package picker

import (
	"errors"
	"github.com/tarm/serial"
	"time"
)

type Port struct {
	Name    string
	Baud    int
	Timeout int
	Serial  *serial.Port
}

func (p *Port) openPort() error {
	var err error
	p.Serial, err = serial.OpenPort(&serial.Config{Name: p.Name, Baud: p.Baud, ReadTimeout: time.Duration(p.Timeout) * time.Millisecond})
	if err != nil {
		return err
	}
	return nil
}

func (p *Port) closePort() error {
	return p.Serial.Close()
}

func (p *Port) write(b Buf) error {
	n, err := p.Serial.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return errors.New("Device: wrong byte count on Write()")
	}
	return nil
}

func (p *Port) read() (Buf, error) {
	response := Buf{}
	tmpBuf := make([]byte, 1)
	packetStarted := false
	var err error
	for err == nil {
		_, err = p.Serial.Read(tmpBuf)
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
	return response, err
}
