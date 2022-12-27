#!/bin/bash

set -ex

here="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null 2>&1 && pwd )"
top="$(dirname ${here})"

function clean_up {
  rm -rf $here/../test_projects/pixie/build
  rm -rf $here/../.ardi $here/../ardi.json
}

trap "clean_up" EXIT

clean_up

ESP8266=https://arduino.esp8266.com/stable/package_esp8266com_index.json

go install ${top}

ardi project-init

ardi add platforms arduino:avr
ardi add lib "Adafruit Pixie"
ardi add board-url ${ESP8266}
ardi add build \
  --name pixie \
  --fqbn arduino:avr:mega \
  --sketch $top/test_projects/pixie

ardi list platforms
ardi search platforms

ardi list libs
ardi search libs "Adafruit Pixie"

ardi exec arduino-cli -- compile --fqbn arduino:avr:mega $top/test_projects/pixie

ardi remove platform arduino:avr
ardi remove lib Adafruit_Pixie
ardi remove board-url ${ESP8266}
ardi remove build pixie
