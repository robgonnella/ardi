#!/bin/bash

set -ex

here="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null 2>&1 && pwd )"

function clean_up {
  rm -rf $here/../test_projects/pixie/build || true
  rm -rf $here/../.ardi $here/../ardi.json || true
}

trap "clean_up" EXIT

clean_up

GO111MODULE=on go get github.com/robgonnella/ardi/v2

ardi project init -v

ardi project add platforms --all -v

ardi project add lib "Adafruit Pixie" -v

ardi project add build \
  --name pixie \
  --platform arduino:avr \
  --fqbn arduino:avr:mega \
  --sketch $here/../test_projects/pixie

ardi project build pixie
