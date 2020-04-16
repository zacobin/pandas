#!/bin/bash
# Copyright (c) Mainflux
# SPDX-License-Identifier: Apache-2.0

###
# Runs all Mainflux microservices (must be previously built and installed).
#
# Expects that PostgreSQL and needed messaging DB are alredy running.
# Additionally, MQTT microservice demands that Redis is up and running.
#
###

BUILD_DIR=../build

# Kill all mainflux-* stuff
function cleanup {
    pkill mainflux
    pkill nats
}

###
# NATS
###
gnatsd &
counter=1
until nc -zv localhost 4222 1>/dev/null 2>&1; 
do
    sleep 0.5
    ((counter++))
    if [ ${counter} -gt 10 ]
    then
        echo -ne "gnatsd failed to start in 5 sec, exiting"
        exit 1
    fi
    echo -ne "Waiting for gnatsd"
done

###
# Users
###
PD_USERS_LOG_LEVEL=info PD_EMAIL_TEMPLATE=../docker/users/emailer/templates/email.tmpl $BUILD_DIR/mainflux-users &

###
# Things
###
PD_THINGS_LOG_LEVEL=info PD_THINGS_HTTP_PORT=8182 PD_THINGS_AUTH_GRPC_PORT=8183 PD_THINGS_AUTH_HTTP_PORT=8989 $BUILD_DIR/mainflux-things &

###
# HTTP
###
PD_HTTP_ADAPTER_LOG_LEVEL=info PD_HTTP_ADAPTER_PORT=8185 PD_THINGS_URL=localhost:8183 $BUILD_DIR/mainflux-http &

###
# WS
###
PD_WS_ADAPTER_LOG_LEVEL=info PD_WS_ADAPTER_PORT=8186 PD_THINGS_URL=localhost:8183 $BUILD_DIR/mainflux-ws &

###
# MQTT
###
PD_MQTT_ADAPTER_LOG_LEVEL=info PD_THINGS_URL=localhost:8183 $BUILD_DIR/mainflux-mqtt &

###
# CoAP
###
PD_COAP_ADAPTER_LOG_LEVEL=info PD_COAP_ADAPTER_PORT=5683 PD_THINGS_URL=localhost:8183 $BUILD_DIR/mainflux-coap &

###
# AUTHN
###
PD_AUTHN_LOG_LEVEL=debug PD_AUTHN_HTTP_PORT=8189 PD_AUTHN_GRPC_PORT=8181 PD_AUTHN_DB_PORT=5432 PD_AUTHN_DB_USER=mainflux PD_AUTHN_DB_PASS=mainflux PD_AUTHN_DB=authn PD_AUTHN_SECRET=secret $BUILD_DIR/mainflux-authn &

trap cleanup EXIT

while : ; do sleep 1 ; done
