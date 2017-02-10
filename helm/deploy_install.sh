#!/bin/bash
set -euo pipefail
source ./common.sh

: ${RELEASE_NAME:=""}
# Append the release name during creation if set
if [[ -n $RELEASE_NAME ]]; then
  HELM_OPTS="$HELM_OPTS --name $RELEASE_NAME"
fi

log "Deploying release using images $UI_IMAGE:$UI_TAG and $API_IMAGE:$API_TAG"
set -x
helm install monocular --set $VALUES_OVERRIDE $HELM_OPTS
