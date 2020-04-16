# Postgres reader

Postgres reader provides message repository implementation for Postgres.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                            | Description                            | Default        |
|-------------------------------------|----------------------------------------|----------------|
| PD_THINGS_URL                       | Things service URL                     | things:8183    |
| PD_POSTGRES_READER_LOG_LEVEL        | Service log level                      | debug          |
| PD_POSTGRES_READER_PORT             | Service HTTP port                      | 9204           |
| PD_POSTGRES_READER_CLIENT_TLS       | TLS mode flag                          | false          |
| PD_POSTGRES_READER_CA_CERTS         | Path to trusted CAs in PEM format      |                |
| PD_POSTGRES_READER_DB_HOST          | Postgres DB host                       | postgres       |
| PD_POSTGRES_READER_DB_PORT          | Postgres DB port                       | 5432           |
| PD_POSTGRES_READER_DB_USER          | Postgres user                          | mainflux       |
| PD_POSTGRES_READER_DB_PASS          | Postgres password                      | mainflux       |
| PD_POSTGRES_READER_DB_NAME          | Postgres database name                 | messages       |
| PD_POSTGRES_READER_DB_SSL_MODE      | Postgres SSL mode                      | disabled       |
| PD_POSTGRES_READER_DB_SSL_CERT      | Postgres SSL certificate path          | ""             |
| PD_POSTGRES_READER_DB_SSL_KEY       | Postgres SSL key                       | ""             |
| PD_POSTGRES_READER_DB_SSL_ROOT_CERT | Postgres SSL root certificate path     | ""             |
| PD_JAEGER_URL                       | Jaeger server URL                      | localhost:6831 |
| PD_POSTGRES_READER_THINGS_TIMEOUT   | Things gRPC request timeout in seconds | 1              |

## Deployment

```yaml
  version: "3.7"
  postgres-writer:
    image: mainflux/postgres-writer:[version]
    container_name: [instance name]
    depends_on:
      - postgres
      - nats
    restart: on-failure
    environment:
      PD_NATS_URL: [NATS instance URL]
      PD_POSTGRES_READER_LOG_LEVEL: [Service log level]
      PD_POSTGRES_READER_PORT: [Service HTTP port]
      PD_POSTGRES_READER_DB_HOST: [Postgres host]
      PD_POSTGRES_READER_DB_PORT: [Postgres port]
      PD_POSTGRES_READER_DB_USER: [Postgres user]
      PD_POSTGRES_READER_DB_PASS: [Postgres password]
      PD_POSTGRES_READER_DB_NAME: [Postgres database name]
      PD_POSTGRES_READER_DB_SSL_MODE: [Postgres SSL mode]
      PD_POSTGRES_READER_DB_SSL_CERT: [Postgres SSL cert]
      PD_POSTGRES_READER_DB_SSL_KEY: [Postgres SSL key]
      PD_POSTGRES_READER_DB_SSL_ROOT_CERT: [Postgres SSL Root cert]
      PD_JAEGER_URL: [Jaeger server URL]
      PD_POSTGRES_READER_THINGS_TIMEOUT: [Things gRPC request timeout in seconds]
    ports:
      - 8903:8903
    networks:
      - docker_mainflux-base-net
```

To start the service, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/cloustone/pandas/mainflux

cd mainflux

# compile the postgres writer
make postgres-writer

# copy binary to bin
make install

# Set the environment variables and run the service
PD_THINGS_URL=[Things service URL] PD_POSTGRES_READER_LOG_LEVEL=[Service log level] PD_POSTGRES_READER_PORT=[Service HTTP port] PD_POSTGRES_READER_CLIENT_TLS =[TLS mode flag] PD_POSTGRES_READER_CA_CERTS=[Path to trusted CAs in PEM format] PD_POSTGRES_READER_DB_HOST=[Postgres host] PD_POSTGRES_READER_DB_PORT=[Postgres port] PD_POSTGRES_READER_DB_USER=[Postgres user] PD_POSTGRES_READER_DB_PASS=[Postgres password] PD_POSTGRES_READER_DB_NAME=[Postgres database name] PD_POSTGRES_READER_DB_SSL_MODE=[Postgres SSL mode] PD_POSTGRES_READER_DB_SSL_CERT=[Postgres SSL cert] PD_POSTGRES_READER_DB_SSL_KEY=[Postgres SSL key] PD_POSTGRES_READER_DB_SSL_ROOT_CERT=[Postgres SSL Root cert] PD_JAEGER_URL=[Jaeger server URL] PD_POSTGRES_READER_THINGS_TIMEOUT=[Things gRPC request timeout in seconds] $GOBIN/mainflux-postgres-reader
```

## Usage

Starting service will start consuming normalized messages in SenML format.
