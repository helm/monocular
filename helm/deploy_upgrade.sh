#!/bin/bash
set -euo pipefail
source ./common.sh

askForConfirmation() {
  log "Release \"$RELEASE_NAME\" already exists, you will upgrade it."
  read -p "Are you sure? " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]];then
    return 1
  fi
}

releaseExists() {
  RELEASES_NUMBER=`helm list -q $RELEASE_NAME | wc -l`
  [ $RELEASES_NUMBER -ne 0 ]
}

upgradeChart() {
  log "Upgrading release \"$RELEASE_NAME\" using images $UI_IMAGE:$UI_TAG and $API_IMAGE:$API_TAG"
  set -x
  helm upgrade $RELEASE_NAME monocular --set $VALUES_OVERRIDE $HELM_OPTS
}

if releaseExists; then
  if [[ $SKIP_UPGRADE_CONFIRMATION != "true" ]]; then
    askForConfirmation
  fi
  upgradeChart
else
  log "Release \"$RELEASE_NAME\" not found. Upgrade aborted"
fi
