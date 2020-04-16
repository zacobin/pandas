# Mainflux CoAP Adapter

Mainflux CoAP adapter provides an [CoAP](http://coap.technology/) API for sending messages through the
platform.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                       | Description                                            | Default               |
|--------------------------------|--------------------------------------------------------|-----------------------|
| PD_COAP_ADAPTER_PORT           | Service listening port                                 | 5683                  |
| PD_NATS_URL                    | NATS instance URL                                      | nats://localhost:4222 |
| PD_THINGS_URL                  | Things service URL                                     | localhost:8181        |
| PD_COAP_ADAPTER_LOG_LEVEL      | Service log level                                      | error                 |
| PD_COAP_ADAPTER_CLIENT_TLS     | Flag that indicates if TLS should be turned on         | false                 |
| PD_COAP_ADAPTER_CA_CERTS       | Path to trusted CAs in PEM format                      |                       |
| PD_COAP_ADAPTER_PING_PERIOD    | Hours between 1 and 24 to ping client with ACK message | 12                    |
| PD_JAEGER_URL                  | Jaeger server URL                                      | localhost:6831        |
| PD_COAP_ADAPTER_THINGS_TIMEOUT | Things gRPC request timeout in seconds                 | 1                     |

## Deployment

The service is distributed as Docker container. The following snippet provides
a compose file template that can be used to deploy the service container locally:

```yaml
version: "2"
services:
  adapter:
    image: mainflux/coap:[version]
    container_name: [instance name]
    ports:
      - [host machine port]:[configured port]
    environment:
      PD_COAP_ADAPTER_PORT: [Service HTTP port]
      PD_NATS_URL: [NATS instance URL]
      PD_THINGS_URL: [Things service URL]
      PD_COAP_ADAPTER_LOG_LEVEL: [Service log level]
      PD_COAP_ADAPTER_CLIENT_TLS: [Flag that indicates if TLS should be turned on]
      PD_COAP_ADAPTER_CA_CERTS: [Path to trusted CAs in PEM format]
      PD_COAP_ADAPTER_PING_PERIOD: [Hours between 1 and 24 to ping client with ACK message]
      PD_JAEGER_URL: [Jaeger server URL]
      PD_COAP_ADAPTER_THINGS_TIMEOUT: [Things gRPC request timeout in seconds]
```

Running this service outside of container requires working instance of the NATS service.
To start the service outside of the container, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/cloustone/pandas/mainflux

cd mainflux

# compile the http
make coap

# copy binary to bin
make install

# set the environment variables and run the service
PD_THINGS_URL=[Things service URL] PD_NATS_URL=[NATS instance URL] PD_COAP_ADAPTER_PORT=[Service HTTP port] PD_COAP_ADAPTER_LOG_LEVEL=[Service log level] PD_COAP_ADAPTER_CLIENT_TLS=[Flag that indicates if TLS should be turned on] PD_COAP_ADAPTER_CA_CERTS=[Path to trusted CAs in PEM format]  PD_COAP_ADAPTER_PING_PERIOD: [Hours between 1 and 24 to ping client with ACK message] PD_JAEGER_URL=[Jaeger server URL] PD_COAP_ADAPTER_THINGS_TIMEOUT=[Things gRPC request timeout in seconds] $GOBIN/mainflux-coap
```

## Usage

If CoAP adapter is running locally (on default 5683 port), a valid URL would be: `coap://localhost/channels/<channel_id>/messages?authorization=<thing_auth_key>`.
Since CoAP protocol does not support `Authorization` header (option) and options have limited size, in order to send CoAP messages, valid `authorization` value (a valid Thing key) must be present in `Uri-Query` option.
