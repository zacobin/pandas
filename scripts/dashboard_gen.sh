#!/bin/sh
set -e -o pipefail

: "${WORKDIR:=dashboard/}"
: "${NOCOMPRESS:=false}"
: "${pushd:=cd}"
: "${popd:=cd ..}"

YARN_INSTALL="${pushd} ${WORKDIR} && \
                yarn && \
                ${popd}"

YARN_BUILD_PROD="${pushd} ${WORKDIR} && \
                yarn build:prod && \
                ${popd}"

YARN_BUILD_SIT="${pushd} ${WORKDIR} && \
                yarn build:stage && \
                ${popd}"

GO_BINDATA="${pushd} ${WORKDIR} && \
                go-bindata-assetfs -pkg dashboard -nocompress=${NOCOMPRESS} dist/... && \
                ${popd}"

sh -c "${YARN_INSTALL}"
sh -c "${YARN_BUILD_PROD}"
sh -c "${GO_BINDATA}"
