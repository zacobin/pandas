FROM golang:1.13.8-alpine AS go-builder

# Replicate the APK repository override.
# If it is no longer necessary to avoid the CDN mirros we should consider dropping this as it is brittle.
RUN sed -e 's/dl-cdn[.]alpinelinux.org/mirrors.aliyun.com/g' -i~ /etc/apk/repositories
# Install our build time packages.
RUN apk add --no-cache build-base git openssh-client upx
# go proxy
RUN go env -w GOPROXY=https://goproxy.io
WORKDIR $GOPATH/src/github.com/cloustone/pandas
COPY . .
# make
RUN make shiro
RUN upx bin/shiro -o bin/main
RUN chmod +x bin/main

FROM alpine
LABEL MAINTAINER="kaisawind <wind.kaisa@gmail.com>"

COPY --from=go-builder /go/src/github.com/cloustone/pandas/bin/main /main

USER 6342627

CMD ["/main"]
