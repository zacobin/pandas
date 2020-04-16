# Users service

Users service provides an HTTP API for managing users. Through this API clients
are able to do the following actions:

- register new accounts
- obtain access tokens
- verify access tokens

For in-depth explanation of the aforementioned scenarios, as well as thorough
understanding of Mainflux, please check out the [official documentation][doc].

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                  | Description                                                             | Default        |
|---------------------------|-------------------------------------------------------------------------|----------------|
| PD_USERS_LOG_LEVEL        | Log level for Users (debug, info, warn, error)                          | error          |
| PD_USERS_DB_HOST          | Database host address                                                   | localhost      |
| PD_USERS_DB_PORT          | Database host port                                                      | 5432           |
| PD_USERS_DB_USER          | Database user                                                           | mainflux       |
| PD_USERS_DB_PASSWORD      | Database password                                                       | mainflux       |
| PD_USERS_DB               | Name of the database used by the service                                | users          |
| PD_USERS_DB_SSL_MODE      | Database connection SSL mode (disable, require, verify-ca, verify-full) | disable        |
| PD_USERS_DB_SSL_CERT      | Path to the PEM encoded certificate file                                |                |
| PD_USERS_DB_SSL_KEY       | Path to the PEM encoded key file                                        |                |
| PD_USERS_DB_SSL_ROOT_CERT | Path to the PEM encoded root certificate file                           |                |
| PD_USERS_HTTP_PORT        | Users service HTTP port                                                 | 8180           |
| PD_USERS_SERVER_CERT      | Path to server certificate in pem format                                |                |
| PD_USERS_SERVER_KEY       | Path to server key in pem format                                        |                |
| PD_JAEGER_URL             | Jaeger server URL                                                       | localhost:6831 |
| PD_EMAIL_DRIVER           | Mail server driver, mail server for sending reset password token        | smtp           |
| PD_EMAIL_HOST             | Mail server host                                                        | localhost      |
| PD_EMAIL_PORT             | Mail server port                                                        | 25             |
| PD_EMAIL_USERNAME         | Mail server username                                                    |                |
| PD_EMAIL_PASSWORD         | Mail server password                                                    |                |
| PD_EMAIL_FROM_ADDRESS     | Email "from" address                                                    |                |
| PD_EMAIL_FROM_NAME        | Email "from" name                                                       |                |
| PD_EMAIL_TEMPLATE         | Email template for sending emails with password reset link              | email.tmpl     |
| PD_TOKEN_RESET_ENDPOINT   | Password request reset endpoint, for constructing link                  | /reset-request |

## Deployment

The service itself is distributed as Docker container. The following snippet
provides a compose file template that can be used to deploy the service container
locally:

```yaml
version: "2"
services:
  users:
    image: mainflux/users:[version]
    container_name: [instance name]
    ports:
      - [host machine port]:[configured HTTP port]
    environment:
      PD_USERS_LOG_LEVEL: [Users log level]
      PD_USERS_DB_HOST: [Database host address]
      PD_USERS_DB_PORT: [Database host port]
      PD_USERS_DB_USER: [Database user]
      PD_USERS_DB_PASS: [Database password]
      PD_USERS_DB: [Name of the database used by the service]
      PD_USERS_DB_SSL_MODE: [SSL mode to connect to the database with]
      PD_USERS_DB_SSL_CERT: [Path to the PEM encoded certificate file]
      PD_USERS_DB_SSL_KEY: [Path to the PEM encoded key file]
      PD_USERS_DB_SSL_ROOT_CERT: [Path to the PEM encoded root certificate file]
      PD_USERS_HTTP_PORT: [Service HTTP port]
      PD_USERS_SERVER_CERT: [String path to server certificate in pem format]
      PD_USERS_SERVER_KEY: [String path to server key in pem format]
      PD_JAEGER_URL: [Jaeger server URL]
      PD_EMAIL_DRIVER: [Mail server driver smtp]
      PD_EMAIL_HOST: [PD_EMAIL_HOST]
      PD_EMAIL_PORT: [PD_EMAIL_PORT]
      PD_EMAIL_USERNAME: [PD_EMAIL_USERNAME]
      PD_EMAIL_PASSWORD: [PD_EMAIL_PASSWORD]
      PD_EMAIL_FROM_ADDRESS: [PD_EMAIL_FROM_ADDRESS]
      PD_EMAIL_FROM_NAME: [PD_EMAIL_FROM_NAME]
      PD_EMAIL_TEMPLATE: [PD_EMAIL_TEMPLATE]
      PD_TOKEN_RESET_ENDPOINT: [PD_TOKEN_RESET_ENDPOINT]
```

To start the service outside of the container, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/cloustone/pandas

cd mainflux

# compile the service
make users

# copy binary to bin
make install

# set the environment variables and run the service
PD_USERS_LOG_LEVEL=[Users log level] PD_USERS_DB_HOST=[Database host address] PD_USERS_DB_PORT=[Database host port] PD_USERS_DB_USER=[Database user] PD_USERS_DB_PASS=[Database password] PD_USERS_DB=[Name of the database used by the service] PD_USERS_DB_SSL_MODE=[SSL mode to connect to the database with] PD_USERS_DB_SSL_CERT=[Path to the PEM encoded certificate file] PD_USERS_DB_SSL_KEY=[Path to the PEM encoded key file] PD_USERS_DB_SSL_ROOT_CERT=[Path to the PEM encoded root certificate file] PD_USERS_HTTP_PORT=[Service HTTP port] PD_USERS_SERVER_CERT=[Path to server certificate] PD_USERS_SERVER_KEY=[Path to server key] PD_JAEGER_URL=[Jaeger server URL] PD_EMAIL_DRIVER=[Mail server driver smtp] PD_EMAIL_HOST=[Mail server host] PD_EMAIL_PORT=[Mail server port] PD_EMAIL_USERNAME=[Mail server username] PD_EMAIL_PASSWORD=[Mail server password] PD_EMAIL_FROM_ADDRESS=[Email from address] PD_EMAIL_FROM_NAME=[Email from name] PD_EMAIL_TEMPLATE=[Email template file] PD_TOKEN_RESET_ENDPOINT=[Password reset token endpoint] $GOBIN/mainflux-users
```

If `PD_EMAIL_TEMPLATE` doesn't point to any file service will function but password reset functionality will not work.

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yaml).

[doc]: http://mainflux.readthedocs.io
