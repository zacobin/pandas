# Postgres writer

Postgres writer provides message repository implementation for Postgres.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                             | Description                                 | Default                |
|--------------------------------------|---------------------------------------------|------------------------|
| PD_NATS_URL                          | NATS instance URL                           | nats://localhost:4222  |
| PD_POSTGRES_WRITER_LOG_LEVEL         | Service log level                           | error                  |
| PD_POSTGRES_WRITER_PORT              | Service HTTP port                           | 9104                   |
| PD_POSTGRES_WRITER_DB_HOST           | Postgres DB host                            | postgres               |
| PD_POSTGRES_WRITER_DB_PORT           | Postgres DB port                            | 5432                   |
| PD_POSTGRES_WRITER_DB_USER           | Postgres user                               | mainflux               |
| PD_POSTGRES_WRITER_DB_PASS           | Postgres password                           | mainflux               |
| PD_POSTGRES_WRITER_DB_NAME           | Postgres database name                      | messages               |
| PD_POSTGRES_WRITER_DB_SSL_MODE       | Postgres SSL mode                           | disabled               |
| PD_POSTGRES_WRITER_DB_SSL_CERT       | Postgres SSL certificate path               | ""                     |
| PD_POSTGRES_WRITER_DB_SSL_KEY        | Postgres SSL key                            | ""                     |
| PD_POSTGRES_WRITER_DB_SSL_ROOT_CERT  | Postgres SSL root certificate path          | ""                     |
| PD_POSTGRES_WRITER_SUBJECTS_CONFIG   | Configuration file path with subjects list  | /config/subjects.toml  |

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
      PD_POSTGRES_WRITER_LOG_LEVEL: [Service log level]
      PD_POSTGRES_WRITER_PORT: [Service HTTP port]
      PD_POSTGRES_WRITER_DB_HOST: [Postgres host]
      PD_POSTGRES_WRITER_DB_PORT: [Postgres port]
      PD_POSTGRES_WRITER_DB_USER: [Postgres user]
      PD_POSTGRES_WRITER_DB_PASS: [Postgres password]
      PD_POSTGRES_WRITER_DB_NAME: [Postgres database name]
      PD_POSTGRES_WRITER_DB_SSL_MODE: [Postgres SSL mode]
      PD_POSTGRES_WRITER_DB_SSL_CERT: [Postgres SSL cert]
      PD_POSTGRES_WRITER_DB_SSL_KEY: [Postgres SSL key]
      PD_POSTGRES_WRITER_DB_SSL_ROOT_CERT: [Postgres SSL Root cert]
      PD_POSTGRES_WRITER_SUBJECTS_CONFIG: [Configuration file path with subjects list]
    ports:
      - 9104:9104
    networks:
      - docker_mainflux-base-net
    volume:
      - ./subjects.yaml:/config/subjects.yaml
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
PD_NATS_URL=[NATS instance URL] PD_POSTGRES_WRITER_LOG_LEVEL=[Service log level] PD_POSTGRES_WRITER_PORT=[Service HTTP port] PD_POSTGRES_WRITER_DB_HOST=[Postgres host] PD_POSTGRES_WRITER_DB_PORT=[Postgres port] PD_POSTGRES_WRITER_DB_USER=[Postgres user] PD_POSTGRES_WRITER_DB_PASS=[Postgres password] PD_POSTGRES_WRITER_DB_NAME=[Postgres database name] PD_POSTGRES_WRITER_DB_SSL_MODE=[Postgres SSL mode] PD_POSTGRES_WRITER_DB_SSL_CERT=[Postgres SSL cert] PD_POSTGRES_WRITER_DB_SSL_KEY=[Postgres SSL key] PD_POSTGRES_WRITER_DB_SSL_ROOT_CERT=[Postgres SSL Root cert] PD_POSTGRES_WRITER_SUBJECTS_CONFIG=[Configuration file path with subjects list] $GOBIN/mainflux-postgres-writer
```

## Usage

Starting service will start consuming normalized messages in SenML format.
