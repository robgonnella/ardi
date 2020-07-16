#!/bin/bash

set -ex

here="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null 2>&1 && pwd )"

function clean_up {
  rm -rf $here/../test_projects/pixie/build || true
  rm -rf $here/../.ardi $here/../ardi.json || true
}

trap "clean_up" EXIT

clean_up

go install $here/../

ardi project-init -v

ardi add platforms --all -v

ardi add lib "Adafruit Pixie" -v

ardi add build \
  --name pixie \
  --fqbn arduino:avr:mega \
  --sketch $here/../test_projects/pixie

ardi build pixie -v
