#!/bin/bash
# generate and install sidecar
SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
set -x
$SCRIPTDIR/install-sidecar.sh gcr.io
