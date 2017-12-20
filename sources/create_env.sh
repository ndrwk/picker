#!/bin/bash
dir=".picker/src"
if [ -d "$dir" ]
then
	echo "$dir is existed."
else
    mkdir .picker
    platformio init --board nanoatmega328 --project-dir .picker
fi
cp ../sources/arduino/device.cpp .picker/src
cp ../sources/arduino/config.h .picker/src
#Adafruit BMP085 Library @ 1.0.0
platformio lib --global install 525@1.0.0
#OneWire @ 2.3.2
platformio lib --global install 54@3.7.7
#DHTlib@None
platformio lib --global install 1336@None