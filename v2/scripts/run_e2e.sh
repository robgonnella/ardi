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

go install ${top}

ardi project-init -v

ardi add platforms arduino:avr -v

ardi add lib "Adafruit Pixie" -v

ardi add build \
  --name pixie \
  --fqbn arduino:avr:mega \
  --sketch $top/test_projects/pixie

ardi compile pixie -v
