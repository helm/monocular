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
# USAGE: repo-sync.sh

GIT_URL=github.com/kubernetes-helm/monocular.git
REPO_URL=https://kubernetes-helm.github.io/monocular
REPO_DIR=$TRAVIS_BUILD_DIR
CHART_PATH="$REPO_DIR/deployment/monocular"
COMMIT_CHANGES=true

log () {
  echo -e "\033[0;33m$(date "+%H:%M:%S")\033[0;37m ==> $1."
}

travis_setup_git() {
  git config user.email "travis@travis-ci.org"
  git config user.name "Travis CI"
  git remote add upstream "https://$GH_TOKEN@$GIT_URL"
}

show_important_vars() {
  echo "  REPO_URL: $REPO_URL"
  echo "  BUILD_DIR: $BUILD_DIR"
  echo "  REPO_DIR: $REPO_DIR"
  echo "  TRAVIS: $TRAVIS"
  echo "  COMMIT_CHANGES: $COMMIT_CHANGES"
}

# https://github.com/bitnami/test-infra/blob/master/circle/docker-release-image.sh#L234
HELM_VERSION=2.5.1
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

update_chart_version() {
  CHART_VERSION=$(grep '^version:' $CHART_PATH/Chart.yaml | awk '{print $2}')
  CHART_VERSION_NEXT="${CHART_VERSION%.*}.$((${CHART_VERSION##*.}+1))"
  sed -i 's|^version:.*|version: '"$CHART_VERSION_NEXT"'|g' $CHART_PATH/Chart.yaml
  sed -i 's|^appVersion:.*|appVersion: '"$TRAVIS_TAG"'|g' $CHART_PATH/Chart.yaml
  sed -i '/bitnami\/monocular/{n; s/tag:.*/tag: '"$TRAVIS_TAG"'/}' $CHART_PATH/values.yaml

  if [ $COMMIT_CHANGES != "false" ]; then
    log "Commiting chart source changes to master branch"
    git add $CHART_PATH/Chart.yaml $CHART_PATH/values.yaml
    git commit --message "chart: bump to $CHART_VERSION_NEXT [skip ci]" --message "travis build #$TRAVIS_BUILD_NUMBER"
    git push -q upstream HEAD:master
  fi
}

show_important_vars

travis_setup_git
git fetch upstream

# Bump chart if this is a release
if [[ -n "$TRAVIS_TAG" ]]; then
  log "Updating chart version"
  update_chart_version
fi

BUILD_DIR=$(mktemp -d)
log "Initializing build directory with existing charts index"
if git show upstream/gh-pages:index.yaml >/dev/null 2>&1; then
  git show upstream/gh-pages:index.yaml > $BUILD_DIR/index.yaml
fi

# Skip repository sync if chart already exists in index
CHART_VERSION=$(grep '^version:' $CHART_PATH/Chart.yaml | awk '{print $2}')
if grep -q "version: $CHART_VERSION" $BUILD_DIR/index.yaml; then
  log "Chart version $CHART_VERSION already exists... skipping"
  exit 0
fi

install_helm

# Package all charts and update index in temporary buildDir
log "Packaging charts from source code"
pushd $BUILD_DIR
  log "Packaging chart"
  helm dep build $CHART_PATH
  helm package $CHART_PATH

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
  git commit --message "release $CHART_VERSION [skip ci]" --message "travis build #$TRAVIS_BUILD_NUMBER"
  git push -q upstream HEAD:gh-pages
fi

log "Repository cleanup and reset"
git reset --hard upstream/master
git clean -df .
