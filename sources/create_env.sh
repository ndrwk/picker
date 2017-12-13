#!/bin/bash
dir=".picker/src"
if [ -d "$dir" ]
then
	echo "$dir has already existed."
else
    mkdir .picker
    platformio init --board nanoatmega328 --project-dir .picker
fi
cp ../sources/arduino/device.cpp .picker/src
cp ../sources/arduino/config.h .picker/src


