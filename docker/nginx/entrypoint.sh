#!/bin/ash

if [ -z "$PD_MQTT_CLUSTER" ]
then
      envsubst '${PD_MQTT_ADAPTER_PORT}' < /etc/nginx/snippets/mqtt-upstream-single.conf > /etc/nginx/snippets/mqtt-upstream.conf
      envsubst '${PD_MQTT_ADAPTER_WS_PORT}' < /etc/nginx/snippets/mqtt-ws-upstream-single.conf > /etc/nginx/snippets/mqtt-ws-upstream.conf
else
      envsubst '${PD_MQTT_ADAPTER_PORT}' < /etc/nginx/snippets/mqtt-upstream-cluster.conf > /etc/nginx/snippets/mqtt-upstream.conf
      envsubst '${PD_MQTT_ADAPTER_WS_PORT}' < /etc/nginx/snippets/mqtt-ws-upstream-cluster.conf > /etc/nginx/snippets/mqtt-ws-upstream.conf
fi

envsubst '
    ${PD_USERS_HTTP_PORT}
    ${PD_THINGS_HTTP_PORT}
    ${PD_THINGS_HTTP_PORT}
    ${PD_HTTP_ADAPTER_PORT}
    ${PD_WS_ADAPTER_PORT}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf

exec nginx -g "daemon off;"