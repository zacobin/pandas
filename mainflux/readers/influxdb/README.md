# InfluxDB reader

InfluxDB reader provides message repository implementation for InfluxDB.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                        | Description                                    | Default        |
|---------------------------------|------------------------------------------------|----------------|
| PD_INFLUX_READER_PORT           | Service HTTP port                              | 8180           |
| PD_INFLUX_READER_DB_NAME        | InfluxDB database name                         | mainflux       |
| PD_INFLUX_READER_DB_HOST        | InfluxDB host                                  | localhost      |
| PD_INFLUX_READER_DB_PORT        | Default port of InfluxDB database              | 8086           |
| PD_INFLUX_READER_DB_USER        | Default user of InfluxDB database              | mainflux       |
| PD_INFLUX_READER_DB_PASS        | Default password of InfluxDB user              | mainflux       |
| PD_INFLUX_READER_CLIENT_TLS     | Flag that indicates if TLS should be turned on | false          |
| PD_INFLUX_READER_CA_CERTS       | Path to trusted CAs in PEM format              |                |
| PD_INFLUX_READER_SERVER_CERT    | Path to server certificate in pem format       |                |
| PD_INFLUX_READER_SERVER_KEY     | Path to server key in pem format               |                |
| PD_JAEGER_URL                   | Jaeger server URL                              | localhost:6831 |
| PD_INFLUX_READER_THINGS_TIMEOUT | Things gRPC request timeout in seconds         | 1              |

## Deployment

```yaml
  version: "3.7"
  influxdb-reader:
    image: mainflux/influxdb-reader:[version]
    container_name: [instance name]
    restart: on-failure
    environment:
      PD_THINGS_URL: [Things service URL]
      PD_INFLUX_READER_PORT: [Service HTTP port]
      PD_INFLUX_READER_DB_NAME: [InfluxDB name]
      PD_INFLUX_READER_DB_HOST: [InfluxDB host]
      PD_INFLUX_READER_DB_PORT: [InfluxDB port]
      PD_INFLUX_READER_DB_USER: [InfluxDB admin user]
      PD_INFLUX_READER_DB_PASS: [InfluxDB admin password]
      PD_INFLUX_READER_CLIENT_TLS: [Flag that indicates if TLS should be turned on]
      PD_INFLUX_READER_CA_CERTS: [Path to trusted CAs in PEM format]
      PD_INFLUX_READER_SERVER_CERT: [String path to server cert in pem format]
      PD_INFLUX_READER_SERVER_KEY: [String path to server key in pem format]
      PD_JAEGER_URL: [Jaeger server URL]
      PD_INFLUX_READER_THINGS_TIMEOUT: [Things gRPC request timeout in seconds]
    ports:
      - [host machine port]:[configured HTTP port]
```

To start the service, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/cloustone/pandas/mainflux

cd mainflux

# compile the influxdb-reader
make influxdb-reader

# copy binary to bin
make install

# Set the environment variables and run the service
PD_THINGS_URL=[Things service URL] \
PD_INFLUX_READER_PORT=[Service HTTP port] \
PD_INFLUX_READER_DB_NAME=[InfluxDB database name] \
PD_INFLUX_READER_DB_HOST=[InfluxDB database host] \
PD_INFLUX_READER_DB_PORT=[InfluxDB database port] \
PD_INFLUX_READER_DB_USER=[InfluxDB admin user] \
PD_INFLUX_READER_DB_PASS=[InfluxDB admin password] \
PD_INFLUX_READER_CLIENT_TLS=[Flag that indicates if TLS should be turned on] \
PD_INFLUX_READER_CA_CERTS=[Path to trusted CAs in PEM format] \
PD_INFLUX_READER_SERVER_CERT=[Path to server pem certificate file] \
PD_INFLUX_READER_SERVER_KEY=[Path to server pem key file] \
PD_JAEGER_URL=[Jaeger server URL] \
PD_INFLUX_READER_THINGS_TIMEOUT=[Things gRPC request timeout in seconds] \
$GOBIN/mainflux-influxdb

```

### Using docker-compose

This service can be deployed using docker containers. Docker compose file is
available in `<project_root>/docker/addons/influxdb-reader/docker-compose.yml`.
In order to run all Mainflux core services, as well as mentioned optional ones,
execute following command:

```bash
docker-compose -f docker/docker-compose.yml up -d
docker-compose -f docker/addons/influxdb-reader/docker-compose.yml up -d
```

## Usage

Service exposes [HTTP API][doc] for fetching messages.

[doc]: ../swagger.yml
