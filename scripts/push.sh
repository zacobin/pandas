#!/usr/bin/env bash

registry=10.4.47.129:5000
#registry=192.168.64.1:5000

pull_image() {
    docker pull $registry/$1
    docker tag $registry/$1 $1
    docker rmi $registry/$1
}

push_image() {
    docker tag $1 $registry/$1
    docker push $registry/$1
    docker rmi $registry/$1
}

if [ $# -ge 1 ]; then
    push_image pandas/pandas-$1
else
    push_image pandas/pandas-rulechain
    push_image pandas/pandas-lbs
    push_image pandas/pandas-authn
    push_image pandas/pandas-authz
    push_image pandas/pandas-things
    push_image pandas/pandas-bootstrap
    push_image pandas/pandas-twins
    push_image pandas/pandas-users
    push_image pandas/pandas-vms
    push_image pandas/pandas-pms
    push_image pandas/pandas-realms
    push_image pandas/pandas-swagger
    push_image pandas/pandas-http
    push_image pandas/pandas-ws
    push_image pandas/pandas-coap
    push_image pandas/pandas-lora
    push_image pandas/pandas-opcua
    push_image pandas/pandas-mqtt
    push_image pandas/pandas-cli
    push_image pandas/pandas-influxdb-writer
    push_image pandas/pandas-influxdb-reader
    push_image pandas/pandas-mongodb-writer
    push_image pandas/pandas-mongodb-reader
    push_image pandas/pandas-cassandra-writer
    push_image pandas/pandas-cassandra-reader
    push_image pandas/pandas-postgres-writer
    push_image pandas/pandas-postgres-reader
fi
#push_image redis:alpine
#push_image bitnami/rabbitmq
#push_image postgres:latest
#push_image elcolio/etcd
