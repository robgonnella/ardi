#!/bin/bash

set -euox pipefail

rm -rf coverage.out

cat unit.out > coverage.out
cat int.out | awk 'NR > 1' >> coverage.out

go tool cover -func coverage.out