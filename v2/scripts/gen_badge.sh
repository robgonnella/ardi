#!/bin/bash

set -ex

toNum() {
  echo "${1%\%}"
}

request_badge() {
  local coverage=$(toNum $1)
  local color='red'
  if (( $(echo "${coverage} > 85.0" | bc -l) )); then
    color='brightgreen'
  fi
  cat <<-EOF > badge.json
  {
    "schemaVersion": 1,
    "label": "coverage",
    "message": "${coverage}",
    "color": "${color}"
  }
EOF
}

request_badge "$@"
