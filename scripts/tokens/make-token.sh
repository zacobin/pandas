#! /bin/bash

# jti: ID
# iss: Issuer
# roles: custom claim
#
# Token for user
token='token-apikey-user.jwt'
echo \
    '{"jti": "test", "iss": "example.com", "roles": [ "user" ]}' |
    jwt -key ./keys/apiKey.prv -alg RS256 -sign - >${token}
jwt -key ./keys/apiKey.pem -alg RS256 -verify ${token}
jwt -show ${token}

# Token for API Keys
token='token-apikey-admin.jwt'
echo \
    '{"jti": "test", "iss": "example.com", "roles": [ "admin" ]}' |
    jwt -key ./keys/apiKey.prv -alg RS256 -sign - >${token}
jwt -key ./keys/apiKey.pem -alg RS256 -verify ${token}
jwt -show ${token}
