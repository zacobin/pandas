#!/bin/bash
set -e -o pipefail

: "${WORKDIR:=./vms/}"
: "${NOCOMPRESS:=false}"


GO_BINDATA="pushd ${WORKDIR} && \
                go-bindata-assetfs -pkg vms -nocompress=${NOCOMPRESS} dist/... && \
                popd"

bash -c "${GO_BINDATA}"
