package picker

import (
	"fmt"
	"github.com/tarm/serial"
	"time"
	"errors"
)

const (
	SLIP_END = 0xc0
	SLIP_ESC = 0xdb
	SLIP_ESC_END = 0xdc
	SLIP_ESC_ESC = 0xdd
)

type Buf []byte
type Crc [2]byte

func (b Buf) AddCrc() Buf {
	for _, v := range b.getCrc() {
		b = append(b, v)
	}
	return b
}

func (b Buf) RemoveCrc() Buf {
	b = b[0:len(b) - 2]
	return b
}

func (b Buf) CheckCrc() bool {
	msgCrc := b[len(b) - 2: ]
	crc := b.RemoveCrc().getCrc()
	for i := range crc {
		if crc[i] != msgCrc[i] {
			return false
		}
	}
	return true
}

func (b Buf) getCrc() Crc {
	res := 0xffff
	for _, c := range b {
		res ^= int(c)
		for j := 1; j <= 8; j++ {
			flag := res & 0x001
			res >>= 1
			if flag > 0 {
				res ^= 0xa001
			}
		}
	}
	tmp := res >> 8
	res = (res << 8) | tmp
	res &= 0xffff
	crc := [2]byte{}
	crc[0] = byte(res & 0xff)
	crc[1] = byte(res >> 8)
	return crc
}

func (b Buf) Slip() Buf {
	slipped := Buf{}
	slipped = append(slipped, SLIP_END)
	for _, v := range b {
		switch v {
		case SLIP_END:
			slipped = append(slipped, SLIP_ESC, SLIP_ESC_END)
		case SLIP_ESC:
			slipped = append(slipped, SLIP_ESC, SLIP_ESC_ESC)
		default:
			slipped = append(slipped, v)
		}
	}
	slipped = append(slipped, SLIP_END)
	return slipped
}

func (b Buf) UnSlip() (Buf, error) {
	unslipped := Buf{}
	packet := Buf{}
	started := false
	escaped := false
	for _, v := range b {
		switch v {
		case SLIP_END:
			if started {
				for i := range packet {
					unslipped = append(unslipped, packet[i])
				}
			} else {
				started = true
			}
			packet = packet[:0]
		case SLIP_ESC:
			escaped = true
		case SLIP_ESC_END:
			if escaped {
				packet = append(packet, SLIP_END)
				escaped = false
			} else {
				packet = append(packet, v)
			}
		case SLIP_ESC_ESC:
			if escaped {
				packet = append(packet, SLIP_ESC)
				escaped = false
			} else {
				packet = append(packet, v)
			}
		default:
			if escaped {
				packet = packet[:0]
				escaped = false
				return nil, errors.New("Protocol: error on UnSlip operation")
			} else {
				if started {
					packet = append(packet, v)
					started = true
				}
			}
		}
	}
	if len(unslipped) == 0 {
		return nil, errors.New("Protocol: error on UnSlip operation")
	}
	return unslipped, nil
}

func (b Buf) ToString() string {
	s := ""
	if len(b) > 0 {
		for i := 0; i < (len(b) - 1); i++ {
			s += fmt.Sprintf("%02X", b[i])
			s += fmt.Sprint(" ")
		}
		s += fmt.Sprintf("%02X", b[len(b) - 1])
	}
	return s
}

func (b Buf) PrintHex() {
	//printHexBuf(b)
	fmt.Println(b.ToString())
}

type Port struct {
	Name    string
	Baud    int
	Timeout int
	Serial  *serial.Port
}

func (p Port) Open() (*serial.Port, error) {
	return serial.OpenPort(&serial.Config{Name: p.Name, Baud: p.Baud, ReadTimeout: time.Duration(p.Timeout) * time.Millisecond})
}

func (p Port) Close() error{
	return p.Serial.Close()
}

func (p Port) Write(b Buf) error {
	n, err := p.Serial.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return errors.New("Device: wrong byte count on Write()")
	}
	return nil
}

func (p Port) Read() (Buf, error) {
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
