#!/bin/bash
set -e -o pipefail

: "${WORKDIR:=./pms/}"
: "${NOCOMPRESS:=false}"


GO_BINDATA="pushd ${WORKDIR} && \
                go-bindata-assetfs -pkg pms -nocompress=${NOCOMPRESS} dist/... && \
                popd"

bash -c "${GO_BINDATA}"
