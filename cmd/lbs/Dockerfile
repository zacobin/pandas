FROM golang:1.13.8-alpine AS go-builder

# Replicate the APK repository override.
# If it is no longer necessary to avoid the CDN mirros we should consider dropping this as it is brittle.
RUN sed -e 's/dl-cdn[.]alpinelinux.org/mirrors.aliyun.com/g' -i~ /etc/apk/repositories
# Install our build time packages.
RUN apk add --no-cache build-base git openssh-client upx npm
RUN npm install -g yarn
# go proxy
RUN go env -w GOPROXY=https://goproxy.io
WORKDIR $GOPATH/src/github.com/cloustone
RUN go get -u github.com/jteeuwen/go-bindata/...
RUN go get github.com/elazarl/go-bindata-assetfs/...
COPY . .
# make
RUN cd pandas \
    && git rm -r --cached dashboard \
    && rm -rf dashboard \
    && git submodule add --force https://github.com/cloustone/pandas-dashboard dashboard \
    && ./scripts/dashboard_gen.sh \
    && make lbs \
    && upx bin/lbs -o bin/main \
    && chmod +x bin/main

FROM alpine
LABEL MAINTAINER="kaisawind <wind.kaisa@gmail.com>"

COPY --from=go-builder /go/src/github.com/cloustone/pandas/bin/main /main

USER root

CMD ["/main"]
