#!/bin/bash
set -e -o pipefail

: "${WORKDIR:=./lbs/}"
: "${NOCOMPRESS:=false}"


GO_BINDATA="pushd ${WORKDIR} && \
                go-bindata-assetfs -pkg lbs -nocompress=${NOCOMPRESS} dist/... && \
                popd"

bash -c "${GO_BINDATA}"
