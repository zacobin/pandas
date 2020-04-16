# WebSocket adapter

WebSocket adapter provides an WebSocket API for sending and receiving messages through the platform.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                     | Description                                    | Default               |
|------------------------------|------------------------------------------------|-----------------------|
| PD_WS_ADAPTER_CLIENT_TLS     | Flag that indicates if TLS should be turned on | false                 |
| PD_WS_ADAPTER_CA_CERTS       | Path to trusted CAs in PEM format              |                       |
| PD_WS_ADAPTER_LOG_LEVEL      | Log level for the WS Adapter                   | error                 |
| PD_WS_ADAPTER_PORT           | Service WS port                                | 8180                  |
| PD_NATS_URL                  | NATS instance URL                              | nats://localhost:4222 |
| PD_THINGS_URL                | Things service URL                             | localhost:8181        |
| PD_JAEGER_URL                | Jaeger server URL                              | localhost:6831        |
| PD_WS_ADAPTER_THINGS_TIMEOUT | Things gRPC request timeout in seconds         | 1                     |

## Deployment

The service is distributed as Docker container. The following snippet provides
a compose file template that can be used to deploy the service container locally:

```yaml
version: "2"
services:
  ws:
    image: mainflux/ws:[version]
    container_name: [instance name]
    ports:
      - [host machine port]:[configured port]
    environment:
      PD_THINGS_URL: [Things service URL]
      PD_NATS_URL: [NATS instance URL]
      PD_WS_ADAPTER_PORT: [Service WS port]
      PD_WS_ADAPTER_LOG_LEVEL: [WS adapter log level]
      PD_WS_ADAPTER_CLIENT_TLS: [Flag that indicates if TLS should be turned on]
      PD_WS_ADAPTER_CA_CERTS: [Path to trusted CAs in PEM format]
      PD_JAEGER_URL: [Jaeger server URL]
      PD_WS_ADAPTER_THINGS_TIMEOUT: [Things gRPC request timeout in seconds]
```

To start the service outside of the container, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/cloustone/pandas/mainflux

cd mainflux

# compile the ws
make ws

# copy binary to bin
make install

# set the environment variables and run the service
PD_THINGS_URL=[Things service URL] PD_NATS_URL=[NATS instance URL] PD_WS_ADAPTER_PORT=[Service WS port] PD_WS_ADAPTER_LOG_LEVEL=[WS adapter log level] PD_WS_ADAPTER_CLIENT_TLS=[Flag that indicates if TLS should be turned on] PD_WS_ADAPTER_CA_CERTS=[Path to trusted CAs in PEM format] PD_JAEGER_URL=[Jaeger server URL] PD_WS_ADAPTER_THINGS_TIMEOUT=[Things gRPC request timeout in seconds] $GOBIN/mainflux-ws
```

## Usage

For more information about service capabilities and its usage, please check out
the [WebSocket paragraph](https://mainflux.readthedocs.io/en/latest/messaging/#websocket) in the Getting Started guide.
