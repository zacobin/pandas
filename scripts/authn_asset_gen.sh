#!/bin/bash
set -e -o pipefail

: "${WORKDIR:=./authn/}"
: "${NOCOMPRESS:=false}"


GO_BINDATA="pushd ${WORKDIR} && \
                go-bindata-assetfs -pkg authn -nocompress=${NOCOMPRESS} dist/... && \
                popd"

bash -c "${GO_BINDATA}"
