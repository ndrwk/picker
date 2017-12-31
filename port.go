package picker

import (
	"errors"
	"github.com/tarm/serial"
	"time"
	"sync"
)

type Port struct {
	name    string
	baud    int
	timeout int
	inUse   sync.Mutex
	serial  *serial.Port
}

func (p *Port) openPort() error {
	var err error
	p.serial, err = serial.OpenPort(&serial.Config{Name: p.name, Baud: p.baud, ReadTimeout: time.Duration(p.timeout) * time.Millisecond})
	if err != nil {
		return err
	}
	return nil
}

func (p *Port) closePort() error {
	return p.serial.Close()
}

func (p *Port) write(b Buf) error {
	n, err := p.serial.Write(b)
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
		_, err = p.serial.Read(tmpBuf)
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
