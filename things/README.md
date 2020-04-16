# Things

Things service provides an HTTP API for managing platform resources: things and channels.
Through this API clients are able to do the following actions:

- provision new things
- create new channels
- "connect" things into the channels

For an in-depth explanation of the aforementioned scenarios, as well as thorough
understanding of Mainflux, please check out the [official documentation][doc].

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                    | Description                                                            | Default        |
|-----------------------------|------------------------------------------------------------------------|----------------|
| PD_THINGS_LOG_LEVEL         | Log level for Things (debug, info, warn, error)                        | error          |
| PD_THINGS_DB_HOST           | Database host address                                                  | localhost      |
| PD_THINGS_DB_PORT           | Database host port                                                     | 5432           |
| PD_THINGS_DB_USER           | Database user                                                          | mainflux       |
| PD_THINGS_DB_PASS           | Database password                                                      | mainflux       |
| PD_THINGS_DB                | Name of the database used by the service                               | things         |
| PD_THINGS_DB_SSL_MODE       | Database connection SSL mode (disable, require, verify-ca, verify-full)| disable        |
| PD_THINGS_DB_SSL_CERT       | Path to the PEM encoded certificate file                               |                |
| PD_THINGS_DB_SSL_KEY        | Path to the PEM encoded key file                                       |                |
| PD_THINGS_DB_SSL_ROOT_CERT  | Path to the PEM encoded root certificate file                          |                |
| PD_THINGS_CLIENT_TLS        | Flag that indicates if TLS should be turned on                         | false          |
| PD_THINGS_CA_CERTS          | Path to trusted CAs in PEM format                                      |                |
| PD_THINGS_CACHE_URL         | Cache database URL                                                     | localhost:6379 |
| PD_THINGS_CACHE_PASS        | Cache database password                                                |                |
| PD_THINGS_CACHE_DB          | Cache instance name                                                    | 0              |
| PD_THINGS_ES_URL            | Event store URL                                                        | localhost:6379 |
| PD_THINGS_ES_PASS           | Event store password                                                   |                |
| PD_THINGS_ES_DB             | Event store instance name                                              | 0              |
| PD_THINGS_HTTP_PORT         | Things service HTTP port                                               | 8180           |
| PD_THINGS_AUTH_HTTP_PORT    | Things service auth HTTP port                                          | 8989           |
| PD_THINGS_AUTH_GRPC_PORT    | Things service auth gRPC port                                          | 8181           |
| PD_THINGS_SERVER_CERT       | Path to server certificate in pem format                               |                |
| PD_THINGS_SERVER_KEY        | Path to server key in pem format                                       |                |
| PD_USERS_URL                | Users service URL                                                      | localhost:8181 |
| PD_THINGS_SINGLE_USER_EMAIL | User email for single user mode (no gRPC communication with users)     |                |
| PD_THINGS_SINGLE_USER_TOKEN | User token for single user mode that should be passed in auth header   |                |
| PD_JAEGER_URL               | Jaeger server URL                                                      | localhost:6831 |
| PD_THINGS_USERS_TIMEOUT     | Users gRPC request timeout in seconds                                  | 1              |

**Note** that if you want `things` service to have only one user locally, you should use `PD_THINGS_SINGLE_USER` env vars. By specifying these, you don't need `users` service in your deployment as it won't be used for authorization.

## Deployment

The service itself is distributed as Docker container. The following snippet
provides a compose file template that can be used to deploy the service container
locally:

```yaml
version: "2"
services:
  things:
    image: things:[version]
    container_name: [instance name]
    ports:
      - [host machine port]:[configured HTTP port]
    environment:
      PD_THINGS_LOG_LEVEL: [Things log level]
      PD_THINGS_DB_HOST: [Database host address]
      PD_THINGS_DB_PORT: [Database host port]
      PD_THINGS_DB_USER: [Database user]
      PD_THINGS_DB_PASS: [Database password]
      PD_THINGS_DB: [Name of the database used by the service]
      PD_THINGS_DB_SSL_MODE: [SSL mode to connect to the database with]
      PD_THINGS_DB_SSL_CERT: [Path to the PEM encoded certificate file]
      PD_THINGS_DB_SSL_KEY: [Path to the PEM encoded key file]
      PD_THINGS_DB_SSL_ROOT_CERT: [Path to the PEM encoded root certificate file]
      PD_THINGS_CA_CERTS: [Path to trusted CAs in PEM format]
      PD_THINGS_CACHE_URL: [Cache database URL]
      PD_THINGS_CACHE_PASS: [Cache database password]
      PD_THINGS_CACHE_DB: [Cache instance that should be used]
      PD_THINGS_ES_URL: [Event store URL]
      PD_THINGS_ES_PASS: [Event store password]
      PD_THINGS_ES_DB: [Event store instance name]
      PD_THINGS_HTTP_PORT: [Service HTTP port]
      PD_THINGS_AUTH_HTTP_PORT: [Service auth HTTP port]
      PD_THINGS_AUTH_GRPC_PORT: [Service auth gRPC port]
      PD_THINGS_SERVER_CERT: [String path to server cert in pem format]
      PD_THINGS_SERVER_KEY: [String path to server key in pem format]
      PD_USERS_URL: [Users service URL]
      PD_THINGS_SECRET: [String used for signing tokens]
      PD_THINGS_SINGLE_USER_EMAIL: [User email for single user mode (no gRPC communication with users)]
      PD_THINGS_SINGLE_USER_TOKEN: [User token for single user mode that should be passed in auth header]
      PD_JAEGER_URL: [Jaeger server URL]
      PD_THINGS_USERS_TIMEOUT: [Users gRPC request timeout in seconds]
```

To start the service outside of the container, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/cloustone/pandas/mainflux

cd mainflux

# compile the things
make things

# copy binary to bin
make install

# set the environment variables and run the service
PD_THINGS_LOG_LEVEL=[Things log level] PD_THINGS_DB_HOST=[Database host address] PD_THINGS_DB_PORT=[Database host port] PD_THINGS_DB_USER=[Database user] PD_THINGS_DB_PASS=[Database password] PD_THINGS_DB=[Name of the database used by the service] PD_THINGS_DB_SSL_MODE=[SSL mode to connect to the database with] PD_THINGS_DB_SSL_CERT=[Path to the PEM encoded certificate file] PD_THINGS_DB_SSL_KEY=[Path to the PEM encoded key file] PD_THINGS_DB_SSL_ROOT_CERT=[Path to the PEM encoded root certificate file] PD_HTTP_ADAPTER_CA_CERTS=[Path to trusted CAs in PEM format] PD_THINGS_CACHE_URL=[Cache database URL] PD_THINGS_CACHE_PASS=[Cache database password] PD_THINGS_CACHE_DB=[Cache instance name] PD_THINGS_ES_URL=[Event store URL] PD_THINGS_ES_PASS=[Event store password] PD_THINGS_ES_DB=[Event store instance name] PD_THINGS_HTTP_PORT=[Service HTTP port] PD_THINGS_AUTH_HTTP_PORT=[Service auth HTTP port] PD_THINGS_AUTH_GRPC_PORT=[Service auth gRPC port] PD_USERS_URL=[Users service URL] PD_THINGS_SERVER_CERT=[Path to server certificate] PD_THINGS_SERVER_KEY=[Path to server key] PD_THINGS_SINGLE_USER_EMAIL=[User email for single user mode (no gRPC communication with users)] PD_THINGS_SINGLE_USER_TOKEN=[User token for single user mode that should be passed in auth header] PD_JAEGER_URL=[Jaeger server URL] PD_THINGS_USERS_TIMEOUT=[Users gRPC request timeout in seconds] $GOBIN/mainflux-things
```

Setting `PD_THINGS_CA_CERTS` expects a file in PEM format of trusted CAs. This will enable TLS against the Users gRPC endpoint trusting only those CAs that are provided.

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yaml).

[doc]: http://mainflux.readthedocs.io
