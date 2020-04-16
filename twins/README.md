# Twins

Service twins is used for CRUD and update of digital twins. Twin is a semantic
representation of a real world entity, be it device, application or something
else. It holds the sequence of attribute based definitions of a real world thing
and refers to the time series of definition based states that hold the
historical data about the represented real world thing.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                   | Description                                                          | Default               |
|----------------------------|----------------------------------------------------------------------|-----------------------|
| PD_TWINS_LOG_LEVEL         | Log level for twin service (debug, info, warn, error)                | error                 |
| PD_TWINS_HTTP_PORT         | Twins service HTTP port                                              | 9021                  |
| PD_TWINS_SERVER_CERT       | Path to server certificate in PEM format                             |                       |
| PD_TWINS_SERVER_KEY        | Path to server key in PEM format                                     |                       |
| PD_JAEGER_URL              | Jaeger server URL                                                    |                       |
| PD_TWINS_DB_NAME           | Database name                                                        | mainflux              |
| PD_TWINS_DB_HOST           | Database host address                                                | localhost             |
| PD_TWINS_DB_PORT           | Database host port                                                   | 27017                 |
| PD_TWINS_SINGLE_USER_EMAIL | User email for single user mode (no gRPC communication with users)   |                       |
| PD_TWINS_SINGLE_USER_TOKEN | User token for single user mode that should be passed in auth header |                       |
| PD_TWINS_CLIENT_TLS        | Flag that indicates if TLS should be turned on                       | false                 |
| PD_TWINS_CA_CERTS          | Path to trusted CAs in PEM format                                    |                       |
| PD_TWINS_MQTT_URL          | Mqtt broker URL for twin CRUD and states update notifications        | tcp://localhost:1883  |
| PD_TWINS_THING_ID          | ID of thing representing twins service & mqtt user                   |                       |
| PD_TWINS_THING_KEY         | Key of thing representing twins service & mqtt pass                  |                       |
| PD_TWINS_CHANNEL_ID        | Mqtt notifications topic                                             |                       |
| PD_NATS_URL                | Mainflux NATS broker URL                                             | nats://127.0.0.1:4222 |
| PD_AUTHN_GRPC_PORT         | Authn service gRPC port                                              | 8181                  |
| PD_AUTHN_TIMEOUT           | Authn gRPC request timeout in seconds                                | 1                     |
| PD_AUTHN_URL               | Authn service URL                                                    | localhost:8181        |

## Deployment

The service itself is distributed as Docker container. The following snippet
provides a compose file template that can be used to deploy the service container
locally:

```yaml
version: "3"
services:
  twins:
    image: twins:[version]
    container_name: [instance name]
    ports:
      - [host machine port]:[configured HTTP port]
    environment:
      PD_TWINS_LOG_LEVEL: [Twins log level]
      PD_TWINS_HTTP_PORT: [Service HTTP port]
      PD_TWINS_SERVER_CERT: [String path to server cert in pem format]
      PD_TWINS_SERVER_KEY: [String path to server key in pem format]
      PD_JAEGER_URL: [Jaeger server URL]
      PD_TWINS_DB_NAME: [Database name]
      PD_TWINS_DB_HOST: [Database host address]
      PD_TWINS_DB_PORT: [Database host port]
      PD_TWINS_SINGLE_USER_EMAIL: [User email for single user mode]
      PD_TWINS_SINGLE_USER_TOKEN: [User token for single user mode]
      PD_TWINS_CLIENT_TLS: [Flag that indicates if TLS should be turned on]
      PD_TWINS_CA_CERTS: [Path to trusted CAs in PEM format]
      PD_TWINS_MQTT_URL: [Mqtt broker URL for twin CRUD and states]
      PD_TWINS_THING_ID: [ID of thing representing twins service]
      PD_TWINS_THING_KEY: [Key of thing representing twins service]
      PD_TWINS_CHANNEL_ID: [Mqtt notifications topic]
      PD_NATS_URL: [Mainflux NATS broker URL]
      PD_AUTHN_GRPC_PORT: [Authn service gRPC port]
      PD_AUTHN_TIMEOUT: [Authn gRPC request timeout in seconds]
      PD_AUTHN_URL: [Authn service URL]
```

To start the service outside of the container, execute the following shell script:

```bash
# download the latest version of the service
go get github.com/cloustone/pandas/mainflux

cd $GOPATH/src/github.com/cloustone/pandas/mainflux

# compile the twins
make twins

# copy binary to bin
make install

# set the environment variables and run the service
PD_TWINS_LOG_LEVEL: [Twins log level] PD_TWINS_HTTP_PORT: [Service HTTP port] PD_TWINS_SERVER_CERT: [String path to server cert in pem format] PD_TWINS_SERVER_KEY: [String path to server key in pem format] PD_JAEGER_URL: [Jaeger server URL] PD_TWINS_DB_NAME: [Database name] PD_TWINS_DB_HOST: [Database host address] PD_TWINS_DB_PORT: [Database host port] PD_TWINS_SINGLE_USER_EMAIL: [User email for single user mode] PD_TWINS_SINGLE_USER_TOKEN: [User token for single user mode] PD_TWINS_CLIENT_TLS: [Flag that indicates if TLS should be turned on] PD_TWINS_CA_CERTS: [Path to trusted CAs in PEM format] PD_TWINS_MQTT_URL: [Mqtt broker URL for twin CRUD and states] PD_TWINS_THING_ID: [ID of thing representing twins service] PD_TWINS_THING_KEY: [Key of thing representing twins service] PD_TWINS_CHANNEL_ID: [Mqtt notifications topic] PD_NATS_URL: [Mainflux NATS broker URL] PD_AUTHN_GRPC_PORT: [Authn service gRPC port] PD_AUTHN_TIMEOUT: [Authn gRPC request timeout in seconds] PD_AUTHN_URL: [Authn service URL] $GOBIN/mainflux-twins
```

## Usage

### Starting twins service

The twins service publishes notifications on an mqtt topic of the format
`channels/<PD_TWINS_CHANNEL_ID>/messages/<twinID>/<crudOp>`, where `crudOp`
stands for the crud operation done on twin - create, update, delete or
retrieve - or state - save state. In order to use twin service, one must
inform it - via environment variables - about the Mainflux thing and
channel used for mqtt notification publishing. You can use an already existing
thing and channel - thing has to be connected to channel - or create new ones.

To set the environment variables, please go to `.env` file and set the following
variables:

```
PD_TWINS_THING_ID=
PD_TWINS_THING_KEY=
PD_TWINS_CHANNEL_ID=
```

with the corresponding values of the desired thing and channel. If you are
running mainflux natively, than do the same thing in the corresponding console
environment.

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yaml).

[doc]: http://mainflux.readthedocs.io
