#include "config.h"

#ifdef DS1820ENABLE
#define ONE_WIRE_BUS DS1820_PIN
#define TEMPERATURE_PRECISION 12
#include <DallasTemperature.h>
#include <OneWire.h>
const unsigned char MAXNUMBERS = 10;
const int wait_ds1820 = 750 / (1 << (12 - TEMPERATURE_PRECISION));
OneWire one_wire(ONE_WIRE_BUS);
DallasTemperature sensors_ds1820(&one_wire);
DeviceAddress dallas_addresses[MAXNUMBERS];
unsigned char numbers = 0;
#endif

#ifdef BMP085ENABLE
#include <Adafruit_BMP085.h>
Adafruit_BMP085 bmp085;
int32_t bmp_pressure;
float bmp_temperature;
#endif

#ifdef DHT22ENABLE
#include <dht.h>
dht DHT;
float dht_humidity;
float dht_temperature;
#endif

#ifdef ANALOGREADENABLE
int analog_number = sizeof(analog_pins)/sizeof(unsigned char);
unsigned short value;
#endif


#define LED 13

// #define SERVO 9

const char SLIP_END = '\xC0';
const char SLIP_ESC = '\xDB';
const char SLIP_ESC_END = '\xDC';
const char SLIP_ESC_ESC = '\xDD';
const char CS_PING = '\x00';
const char CS_INFO = '\x01';
const char LOC_ADR = ADDRESS;
const long BAUD = BAUDRATE;

// Servo servo;

char read_write_buf[256];
int msglen = 0;

void transfer_data(char *buf, unsigned char cnt) {
  Serial.print(SLIP_END);
  for (int i = 0; i < cnt; i++) {
    switch (buf[i]) {
    case SLIP_END:
      Serial.print(SLIP_ESC);
      Serial.print(SLIP_ESC_END);
      break;
    case SLIP_ESC:
      Serial.print(SLIP_ESC);
      Serial.print(SLIP_ESC_ESC);
      break;
    default:
      Serial.print(buf[i]);
      break;
    }
  }
  Serial.print(SLIP_END);
}

unsigned short get_crc(char *buf, unsigned char cnt) {
  unsigned short crc, tmp_val, flag;
  crc = 0xFFFF;
  for (int i = 0; i < cnt; i++) {
    crc ^= (unsigned char)buf[i];
    for (int j = 1; j <= 8; j++) {
      flag = crc & 0x0001;
      crc >>= 1;
      if (flag)
        crc ^= 0xA001;
    }
  }
  tmp_val = crc >> 8;
  crc = (crc << 8) | tmp_val;
  crc &= 0xFFFF;
  return crc;
}

int add_crc(char *buf, unsigned char cnt) {
  unsigned short crc = get_crc(buf, cnt);
  memcpy(&buf[cnt], &crc, 2);
  return cnt + 2;
}

int read_command(char *buf) {
  int i = 0;
  bool escaped = false;
  char c = (char)Serial.read();
  if (c == SLIP_END) {
    bool beginflag = true;
    while (beginflag) {
      char c1 = (char)Serial.read();
      switch (c1) {
      case SLIP_END:
        return i;
        break;
      case SLIP_ESC:
        escaped = true;
        break;
      case SLIP_ESC_END:
        if (escaped) {
          buf[i] = SLIP_END;
          escaped = false;
        } else
          buf[i] = c1;
        i++;
        break;
      case SLIP_ESC_ESC:
        if (escaped) {
          buf[i] = SLIP_ESC;
          escaped = false;
        } else
          buf[i] = c1;
        i++;
        break;
      default:
        if (escaped) {
          return 0;
        } else
          buf[i] = c1;
        i++;
        break;
      }
    }
  }
  return i;
}

void serialEvent() {
  while (Serial.available()) {
    msglen = read_command(read_write_buf);
  }
}

void setup() {
  Serial.begin(BAUD);
#ifdef BMP085ENABLE
  bmp085.begin();
#endif
#ifdef DS1820ENABLE
  sensors_ds1820.begin();
#endif
  pinMode(LED, OUTPUT);
  // servo.attach(SERVO);
}

void loop() {
  if (msglen) {
    unsigned short msgcrc;
    memcpy(&msgcrc, &read_write_buf[msglen - 2], 2);
    unsigned short crc = get_crc(read_write_buf, msglen - 2);

    if (crc == msgcrc) {
      char adr = read_write_buf[0];
      char cs = read_write_buf[1];
      char mtd = read_write_buf[2];
      int len;
      if (adr == LOC_ADR) {
        switch (cs) {
        case CS_PING:
          read_write_buf[0] = LOC_ADR;
          read_write_buf[1] = '\x55';
          read_write_buf[2] = '\xAA';
          read_write_buf[3] = '\x55';
          read_write_buf[4] = '\xAA';
          len = add_crc(read_write_buf, 5);
          delay(100);
          transfer_data(read_write_buf, len);
          digitalWrite(LED, HIGH);
          delay(100);
          digitalWrite(LED, LOW);
          break;
        case CS_INFO:
          switch (mtd) {
#ifdef DS1820ENABLE
          case 0:
            numbers = 0;
            for (int i = 0; i < MAXNUMBERS; i++) {
              if (!sensors_ds1820.getAddress(dallas_addresses[i], i))
                break;
              numbers++;
            }
            read_write_buf[0] = LOC_ADR;
            read_write_buf[1] = numbers;
            len = add_crc(read_write_buf, 2);
            delay(100);
            transfer_data(read_write_buf, len);
            break;
          case 1:
            numbers = 0;
            for (int i = 0; i < MAXNUMBERS; i++) {
              if (!sensors_ds1820.getAddress(dallas_addresses[i], i))
                break;
              numbers++;
            }
            for (unsigned char i = 0; i < numbers; i++) {
              sensors_ds1820.setResolution(dallas_addresses[i], TEMPERATURE_PRECISION);
            }
            sensors_ds1820.requestTemperatures();
            delay(wait_ds1820);
            float temp;
            read_write_buf[0] = LOC_ADR;
            read_write_buf[1] = numbers;
            for (int i = 0; i < numbers; i++) {
              temp = sensors_ds1820.getTempC(dallas_addresses[i]);
              memcpy(&read_write_buf[i * 12 + 2], &temp, 4);
              memcpy(&read_write_buf[i * 12 + 6], &dallas_addresses[i], 8);
            }
            len = add_crc(read_write_buf, numbers * 12 + 2);
            delay(10);
            transfer_data(read_write_buf, len);
            break;
#endif
#ifdef BMP085ENABLE
          case 2:
            bmp_pressure = (int32_t)(bmp085.readPressure() / 133.3224);
            bmp_temperature = bmp085.readTemperature();
            read_write_buf[0] = LOC_ADR;
            memcpy(&read_write_buf[1], &bmp_pressure, 4);
            memcpy(&read_write_buf[5], &bmp_temperature, 4);
            read_write_buf[9] = 0;
            len = add_crc(read_write_buf, 10);
            delay(10);
            transfer_data(read_write_buf, len);
            break;
#endif
#ifdef DHT22ENABLE
          case 3:
            DHT.read22(DHT22_PIN);
            dht_humidity = DHT.humidity;
            dht_temperature = DHT.temperature;
            read_write_buf[0] = LOC_ADR;
            memcpy(&read_write_buf[1], &dht_humidity, 4);
            memcpy(&read_write_buf[5], &dht_temperature, 4);
            read_write_buf[9] = 0;
            len = add_crc(read_write_buf, 10);
            delay(10);
            transfer_data(read_write_buf, len);
            break;
#endif
#ifdef ANALOGREADENABLE
          case 4:
            read_write_buf[0] = LOC_ADR;
            read_write_buf[1] = analog_number;
            for (int i = 0; i < analog_number; i++){
              value = 0x03ff & analogRead(analog_pins[i]);
              memcpy(&read_write_buf[i * 3 + 2], &analog_pins[i], 1);
              memcpy(&read_write_buf[i * 3 + 3], &value, 2);
            }
            len = add_crc(read_write_buf, analog_number * 3 + 2);
            delay(10);
            transfer_data(read_write_buf, len);
            break;
#endif
          }
          break;
        }
      }
    }
    msglen = 0;
  }
  delay(100);
}
