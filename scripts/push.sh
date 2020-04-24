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
    push_image cloustone/pandas-$1
else
    push_image cloustone/pandas-rulechain
    push_image cloustone/pandas-lbs
    push_image cloustone/pandas-authn
    push_image cloustone/pandas-authz
    push_image cloustone/pandas-things
    push_image cloustone/pandas-bootstrap
    push_image cloustone/pandas-twins
    push_image cloustone/pandas-users
    push_image cloustone/pandas-vms
    push_image cloustone/pandas-pms
    push_image cloustone/pandas-realms
    push_image cloustone/pandas-swagger
    push_image cloustone/pandas-http
    push_image cloustone/pandas-ws
    push_image cloustone/pandas-coap
    push_image cloustone/pandas-lora
    push_image cloustone/pandas-opcua
    push_image cloustone/pandas-mqtt
    push_image cloustone/pandas-cli
    push_image cloustone/pandas-influxdb-writer
    push_image cloustone/pandas-influxdb-reader
    push_image cloustone/pandas-mongodb-writer
    push_image cloustone/pandas-mongodb-reader
    push_image cloustone/pandas-cassandra-writer
    push_image cloustone/pandas-cassandra-reader
    push_image cloustone/pandas-postgres-writer
    push_image cloustone/pandas-postgres-reader
fi
#push_image redis:alpine
#push_image bitnami/rabbitmq
#push_image postgres:latest
#push_image elcolio/etcd
