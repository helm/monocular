#!/bin/bash
set -e
INIT_SEM=/tmp/initialized.sem
PACKAGE_FILE=Gopkg.lock

log () {
  echo -e "\033[0;33m$(date "+%H:%M:%S")\033[0;37m ==> $1."
}

dependencies_up_to_date() {
  # It it up to date if the package file is older than
  # the last time the container was initialized
  [ ! $PACKAGE_FILE -nt $INIT_SEM ]
}

if [ "$1" == "gin" -a "$3" == "run" ]; then
	if ! dependencies_up_to_date; then
		log "Packages updating..."
		dep ensure
		log "Packages updated"
	fi
  touch $INIT_SEM

  # Set env vars if .env file exists
  if [ -f .env ]; then
    export $(egrep -v '^#' .env | xargs)
  fi
fi

exec "$@"
