#!/bin/bash

# Copyright 2017 Istio Authors
# Copied from Istio installation  update.sh
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at

#       http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License

ROOT="."
VERSION_FILE="${ROOT}/istio.VERSION"
TEMP_DIR="/tmp"
GIT_COMMIT=false
CHECK_GIT_STATUS=false

set -o errexit
set -o pipefail
set -x

function usage() {
  [[ -n "${1}" ]] && echo "${1}"

  cat <<EOF
usage: ${BASH_SOURCE[0]} [options ...]"
  options:
    -i ... URL to download istioctl binary
    -p ... <hub>,<tag> for the pilot docker image
    -x ... <hub>,<tag> for the mixer docker image
    -c ... <hub>,<tag> for the istio-ca docker image
    -r ... tag for proxy debian package
    -g ... create a git commit for the changes
    -n ... <namespace> namespace in which to install Istio control plane components
    -s ... check if template files have been updated with this tool
    -A ... URL to download auth debian packages
    -P ... URL to download pilot debian packages
    -E ... URL to download proxy debian packages
EOF
  exit 2
}

source "$VERSION_FILE" || error_exit "Could not source versions"

while getopts :gi:n:p:x:c:r:sA:P:E: arg; do
  case ${arg} in
    i) ISTIOCTL_URL="${OPTARG}";;
    n) ISTIO_NAMESPACE="${OPTARG}";;
    p) PILOT_HUB_TAG="${OPTARG}";; # Format: "<hub>,<tag>"
    x) MIXER_HUB_TAG="${OPTARG}";; # Format: "<hub>,<tag>"
    c) CA_HUB_TAG="${OPTARG}";; # Format: "<hub>,<tag>"
    r) PROXY_TAG="${OPTARG}";;
    g) GIT_COMMIT=true;;
    s) CHECK_GIT_STATUS=true;;
    A) AUTH_DEBIAN_URL="${OPTARG}";;
    P) PILOT_DEBIAN_URL="${OPTARG}";;
    E) PROXY_DEBIAN_URL="${OPTARG}";;
    *) usage;;
  esac
done

if [[ -n ${PILOT_HUB_TAG} ]]; then
    PILOT_HUB="$(echo ${PILOT_HUB_TAG}|cut -f1 -d,)"
    PILOT_TAG="$(echo ${PILOT_HUB_TAG}|cut -f2 -d,)"
fi

if [[ -n ${MIXER_HUB_TAG} ]]; then
    MIXER_HUB="$(echo ${MIXER_HUB_TAG}|cut -f1 -d,)"
    MIXER_TAG="$(echo ${MIXER_HUB_TAG}|cut -f2 -d,)"
fi

if [[ -n ${CA_HUB_TAG} ]]; then
    CA_HUB="$(echo ${CA_HUB_TAG}|cut -f1 -d,)"
    CA_TAG="$(echo ${CA_HUB_TAG}|cut -f2 -d,)"
fi

function error_exit() {
  # ${BASH_SOURCE[1]} is the file name of the caller.
  echo "${BASH_SOURCE[1]}: line ${BASH_LINENO[0]}: ${1:-Unknown Error.} (exit ${2:-1})" 1>&2
  exit ${2:-1}
}

function set_git() {
  if [[ ! -e "${HOME}/.gitconfig" ]]; then
    cat > "${HOME}/.gitconfig" << EOF
[user]
  name = istio-testing
  email = istio.testing@gmail.com
EOF
  fi
}


function create_commit() {
  set_git
  # If nothing to commit skip
  check_git_status && return

  echo 'Creating a commit'
  git commit -a -m "Updating istio version" \
    || error_exit 'Could not create a commit'

}

function check_git_status() {
  local git_files="$(git status -s)"
  [[ -z "${git_files}" ]] && return 0
  return 1
}


function update_istio_install() {

  SRC=$ROOT/install/kubernetes/templates
  DEST=$ROOT/install/kubernetes

  ISTIO_INITIALIZER=$DEST/istio-initializer.yaml

  cp $SRC/istio-initializer.yaml.tmpl $ISTIO_INITIALIZER


  echo "# GENERATED FILE. Use with Kubernetes 1.7+" > $ISTIO_INITIALIZER
  echo "# TO UPDATE, modify files in install/kubernetes/templates and run install/updateVersion.sh" >> $ISTIO_INITIALIZER
  cat ${SRC}/istio-initializer.yaml.tmpl >> $ISTIO_INITIALIZER
  sed -i .bak "s|{ISTIO_NAMESPACE}|${ISTIO_NAMESPACE}|" $ISTIO_INITIALIZER
  sed -i .bak "s|{PILOT_HUB}|${PILOT_HUB}|" $ISTIO_INITIALIZER
  sed -i .bak "s|{PILOT_TAG}|${PILOT_TAG}|" $ISTIO_INITIALIZER
  sed -i .bak "s|{PROXY_HUB}|${PROXY_HUB}|" $ISTIO_INITIALIZER
  sed -i .bak "s|{PROXY_TAG}|${PROXY_TAG}|" $ISTIO_INITIALIZER

}




if [[ ${GIT_COMMIT} == true ]]; then
    check_git_status \
      || error_exit "You have modified files. Please commit or reset your workspace."
fi


update_istio_install


if [[ ${GIT_COMMIT} == true ]]; then
    create_commit
fi

if [[ ${CHECK_GIT_STATUS} == true ]]; then
 check_git_status \
   || { echo "Need to update template and run install/updateVersion.sh"; git diff; exit 1; }
fi
