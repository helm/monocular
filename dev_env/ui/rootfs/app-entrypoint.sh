#!/bin/bash -e

# Color Palette
RESET='\033[0m'
MAGENTA='\033[38;5;5m'

log() {
  echo -e "${MAGENTA}$(date "+%T.%2N ")${RESET}${@}" >&2
}

INIT_SEM=/tmp/initialized.sem
PACKAGE_FILE=/app/package.json

fresh_container() {
  [ ! -f $INIT_SEM ]
}

dependencies_up_to_date() {
  # It is up to date if the package file is older than
  # the last time the container was initialized
  [ ! $PACKAGE_FILE -nt $INIT_SEM ]
}

if [ "$1" == ng -a "$2" == "serve" ]; then
  if ! dependencies_up_to_date; then
    log "Installing/Updating Angular dependencies (yarn)"
    yarn
    log "Dependencies updated"
  fi

  if ! fresh_container; then
    echo "#########################################################################"
    echo "                                                                       "
    echo " App initialization skipped:"
    echo " Delete the file $INIT_SEM and restart the container to reinitialize"
    echo " You can alternatively run specific commands using docker-compose exec"
    echo " e.g docker-compose exec myapp npm install angular"
    echo "                                                                       "
    echo "#########################################################################"
  else
    log "Initialization finished"
  fi

  touch $INIT_SEM
fi

exec "$@"
