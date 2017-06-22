#!/bin/bash
set -euo pipefail

# $IMAGE_TAG Required:
: ${API_IMAGE:="bitnami/monocular-api"}
: ${UI_IMAGE:="bitnami/monocular-ui"}

# Helm settings
: ${SKIP_UPGRADE_CONFIRMATION:="false"}
: ${VALUES_OVERRIDE:="api.image.repository=${API_IMAGE},api.image.tag=${API_TAG},ui.image.repository=${UI_IMAGE},ui.image.tag=${UI_TAG}"}
: ${HELM_OPTS:=""}

log () {
  echo -e "\033[0;33m$(date "+%H:%M:%S")\033[0;37m ==> $1."
}
