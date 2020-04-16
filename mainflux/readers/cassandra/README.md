# Cassandra reader

Cassandra reader provides message repository implementation for Cassandra.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                           | Description                                    | Default        |
|------------------------------------|------------------------------------------------|----------------|
| PD_CASSANDRA_READER_PORT           | Service HTTP port                              | 8180           |
| PD_CASSANDRA_READER_DB_CLUSTER     | Cassandra cluster comma separated addresses    | 127.0.0.1      |
| PD_CASSANDRA_READER_DB_KEYSPACE    | Cassandra keyspace name                        | mainflux       |
| PD_CASSANDRA_READER_DB_USERNAME    | Cassandra DB username                          |                |
| PD_CASSANDRA_READER_DB_PASSWORD    | Cassandra DB password                          |                |
| PD_CASSANDRA_READER_DB_PORT        | Cassandra DB port                              | 9042           |
| PD_THINGS_URL                      | Things service URL                             | localhost:8181 |
| PD_CASSANDRA_READER_CLIENT_TLS     | Flag that indicates if TLS should be turned on | false          |
| PD_CASSANDRA_READER_CA_CERTS       | Path to trusted CAs in PEM format              |                |
| PD_CASSANDRA_READER_SERVER_CERT    | Path to server certificate in pem format       |                |
| PD_CASSANDRA_READER_SERVER_KEY     | Path to server key in pem format               |                |
| PD_JAEGER_URL                      | Jaeger server URL                              | localhost:6831 |
| PD_CASSANDRA_READER_THINGS_TIMEOUT | Things gRPC request timeout in seconds         | 1              |


## Deployment

```yaml
  version: "3.7"
  cassandra-reader:
    image: mainflux/cassandra-reader:[version]
    container_name: [instance name]
    expose:
      - [Service HTTP port]
    restart: on-failure
    environment:
      PD_THINGS_URL: [Things service URL]
      PD_CASSANDRA_READER_PORT: [Service HTTP port]
      PD_CASSANDRA_READER_DB_CLUSTER: [Cassandra cluster comma separated addresses]
      PD_CASSANDRA_READER_DB_KEYSPACE: [Cassandra keyspace name]
      PD_CASSANDRA_READER_DB_USERNAME: [Cassandra DB username]
      PD_CASSANDRA_READER_DB_PASSWORD: [Cassandra DB password]
      PD_CASSANDRA_READER_DB_PORT: [Cassandra DB port]
      PD_CASSANDRA_READER_CLIENT_TLS: [Flag that indicates if TLS should be turned on]
      PD_CASSANDRA_READER_CA_CERTS: [Path to trusted CAs in PEM format]
      PD_CASSANDRA_READER_SERVER_CERT: [String path to server cert in pem format]
      PD_CASSANDRA_READER_SERVER_KEY: [String path to server key in pem format]
      PD_JAEGER_URL: [Jaeger server URL]
      PD_CASSANDRA_READER_THINGS_TIMEOUT: [Things gRPC request timeout in seconds]
    ports:
      - [host machine port]:[configured HTTP port]
```

To start the service, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/cloustone/pandas/mainflux

cd mainflux

# compile the cassandra
make cassandra-reader

# copy binary to bin
make install

# Set the environment variables and run the service
PD_THINGS_URL=[Things service URL] \
PD_CASSANDRA_READER_PORT=[Service HTTP port] \
PD_CASSANDRA_READER_DB_CLUSTER=[Cassandra cluster comma separated addresses] \
PD_CASSANDRA_READER_DB_KEYSPACE=[Cassandra keyspace name] \
PD_CASSANDRA_READER_DB_USERNAME=[Cassandra DB username] \
PD_CASSANDRA_READER_DB_PASSWORD=[Cassandra DB password] \
PD_CASSANDRA_READER_DB_PORT=[Cassandra DB port] \
PD_CASSANDRA_READER_CLIENT_TLS=[Flag that indicates if TLS should be turned on] \
PD_CASSANDRA_READER_CA_CERTS=[Path to trusted CAs in PEM format] \
PD_CASSANDRA_READER_SERVER_CERT=[Path to server pem certificate file] \
PD_CASSANDRA_READER_SERVER_KEY=[Path to server pem key file] \
PD_JAEGER_URL=[Jaeger server URL] \
PD_CASSANDRA_READER_THINGS_TIMEOUT=[Things gRPC request timeout in seconds] \
$GOBIN/mainflux-cassandra-reader

```

### Using docker-compose

This service can be deployed using docker containers. Docker compose file is
available in `<project_root>/docker/addons/cassandra-reader/docker-compose.yml`.
In order to run all Mainflux core services, as well as mentioned optional ones,
execute following command:

```bash
docker-compose -f docker/docker-compose.yml up -d
./docker/addons/cassandra-writer/init.sh
docker-compose -f docker/addons/casandra-reader/docker-compose.yml up -d
```

## Usage

Service exposes [HTTP API][doc]  for fetching messages.

[doc]: ../swagger.yml
