#!/bin/bash

set -euox pipefail

here="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

rm -rf $here/../test_artifacts/coverage.out

cat $here/../test_artifacts/unit.out > $here/../test_artifacts/coverage.out
cat $here/../test_artifacts/int.out | awk 'NR > 1' >> $here/../test_artifacts/coverage.out

go tool cover -func $here/../test_artifacts/coverage.out