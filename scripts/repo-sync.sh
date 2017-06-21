#!/bin/bash -e
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Based on https://github.com/migmartri/helm-hack-night-charts/blob/master/repo-sync.sh
# USAGE: repo-sync.sh <commit-changes?>

log () {
  echo -e "\033[0;33m$(date "+%H:%M:%S")\033[0;37m ==> $1."
}

travis_setup_git() {
  git config user.email "travis@travis-ci.org"
  git config user.name "Travis CI"
  COMMIT_MSG="Updating chart repository, travis build #$TRAVIS_BUILD_NUMBER"
  git remote add upstream "https://$GH_TOKEN@github.com/helm/monocular.git"
}

show_important_vars() {
    echo "  REPO_URL: $REPO_URL"
    echo "  BUILD_DIR: $BUILD_DIR"
    echo "  REPO_DIR: $REPO_DIR"
    echo "  TRAVIS: $TRAVIS"
    echo "  COMMIT_CHANGES: $COMMIT_CHANGES"
}

# https://github.com/bitnami/test-infra/blob/master/circle/docker-release-image.sh#L234
HELM_VERSION=2.4.2
install_helm() {
  if ! which helm >/dev/null ; then
    log "Downloading helm..."
    if ! wget -q https://storage.googleapis.com/kubernetes-helm/helm-v${HELM_VERSION}-linux-amd64.tar.gz; then
      log "Could not download helm..."
      return 1
    fi

    log "Installing helm..."
    if ! tar zxf helm-v${HELM_VERSION}-linux-amd64.tar.gz --strip 1 linux-amd64/helm; then
      log "Could not install helm..."
      return 1
    fi
    chmod +x helm
    sudo mv helm /usr/local/bin/helm

    if ! helm version --client; then
      return 1
    fi

    if ! helm init --client-only >/dev/null; then
      return 1
    fi
  fi
}

COMMIT_CHANGES="${1}"
: ${COMMIT_CHANGES:=false}
: ${TRAVIS:=false}
REPO_URL=https://helm.github.io/monocular
BUILD_DIR=$(mktemp -d)
# Current directory
REPO_DIR="$( cd "$(dirname $(dirname "${BASH_SOURCE[0]}"))" && pwd )"
CHART_PATH="$REPO_DIR/deployment/monocular"
COMMIT_MSG="Updating chart repository"

show_important_vars

if [ $TRAVIS != "false" ]; then
  log "Configuring git for Travis-ci"
  travis_setup_git
else
  git remote add upstream git@github.com:helm/monocular.git || true
fi

git fetch upstream
git checkout gh-pages

log "Initializing build directory with existing charts index"
if [ -f index.yaml ]; then
  cp index.yaml $BUILD_DIR
fi

git checkout master

update_chart_version() {
  CHART_VERSION=$(grep '^version:' $1/Chart.yaml | awk '{print $2}')
  CHART_VERSION_NEXT="${CHART_VERSION%.*}.$((${CHART_VERSION##*.}+1))"
  sed -i 's|^version:.*|version: '"$CHART_VERSION_NEXT"'|g' $1/Chart.yaml
  sed -i 's|^appVersion:.*|appVersion: '"$2"'|g' $1/Chart.yaml
  sed -i '/bitnami\/monocular/{n; s/tag:.*/tag: '"$2"'/}' $1/values.yaml

  if [ $COMMIT_CHANGES != "false" ]; then
    log "Commiting chart source changes to master branch"
    git add $1/Chart.yaml $1/values.yaml
    git commit --message "chart: bump to $CHART_VERSION_NEXT"
    git push -q upstream HEAD:master
  fi
}

log "Updating chart version"
update_chart_version $CHART_PATH $TRAVIS_TAG

install_helm

# Package all charts and update index in temporary buildDir
log "Packaging charts from source code"
pushd $BUILD_DIR
  log "Packaging chart"
  helm package $REPO_DIR/deployment/monocular

  log "Indexing repository"
  if [ -f index.yaml ]; then
    helm repo index --url ${REPO_URL} --merge index.yaml .
  else
    helm repo index --url ${REPO_URL} .
  fi
popd

git reset upstream/gh-pages
cp $BUILD_DIR/* $REPO_DIR

# Commit changes are not enabled during PR
if [ $COMMIT_CHANGES != "false" ]; then
  log "Commiting changes to gh-pages branch"
  git add *.tgz index.yaml
  git commit --message "$COMMIT_MSG"
  git push -q upstream HEAD:gh-pages
fi

log "Repository cleanup and reset"
git reset --hard upstream/master
git clean -df .