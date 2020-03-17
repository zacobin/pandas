//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package auth

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// Keys used to sign and verify our tokens
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

type JwtClaims struct {
	jwt.StandardClaims
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	AccessId string   `json:"accessId"`
}

func InitPublicKey(keyPath string) error {
	// loads public keys to verify our tokens
	verifyKeyBuf, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return errors.New("Cannot load public key for tokens")
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyKeyBuf)
	if err != nil {
		return errors.New("Invalid public key for tokens")
	}
	return nil
}

func InitPrivateKey(keyPath string) error {
	signKeyBuf, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return errors.New("Cannot load private key for tokens")
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signKeyBuf)
	if err != nil {
		return errors.New("Invalid private key for tokens")
	}
	return nil
}

func GenToken(claims JwtClaims) (string, error) {
	// Creat token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	return token.SignedString(signKey)
}

func ParseAndCheckToken(token string) (*JwtClaims, error) {
	// the API key is a JWT signed by us with a claim to be a reseller
	parsedToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(parsedToken *jwt.Token) (interface{}, error) {
		// the key used to validate tokens
		return verifyKey, nil
	})

	if err == nil {
		if claims, ok := parsedToken.Claims.(*JwtClaims); ok && parsedToken.Valid {
			return claims, nil
		}
	}
	return nil, err
}
