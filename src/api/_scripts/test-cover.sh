#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

# shellcheck disable=SC2046
pkgs=$(go list ./...)
COVERAGE_EXCLUDES=("github.com/kubernetes-helm/monocular/src/api/swagger/models" \
"github.com/kubernetes-helm/monocular/src/api/swagger/restapi" \
"github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations" \
"github.com/kubernetes-helm/monocular/src/api")

echo "" > coverage.txt
for p in $pkgs; do
    if [[ " ${COVERAGE_EXCLUDES[@]} " =~ " ${p} " ]]; then
      echo "skipping test coverage for ${p}"
    else
      go test -covermode=atomic -coverprofile=profile.out "$p"
      if [ -s profile.out ]; then
          cat profile.out >> coverage.txt
      fi
      rm -f profile.out
    fi
done
