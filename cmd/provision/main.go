package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/pkg/errors"
	"github.com/cloustone/pandas/pkg/logger"
	"github.com/cloustone/pandas/provision"
	"github.com/cloustone/pandas/provision/api"
	mfSDK "github.com/cloustone/pandas/sdk/go"
)

const (
	defLogLevel        = "debug"
	defConfigFile      = "config.toml"
	defTLS             = "false"
	defServerCert      = ""
	defServerKey       = ""
	defThingsLocation  = "http://localhost"
	defUsersLocation   = "http://localhost"
	defMQTTURL         = "localhost:1883"
	defHTTPPort        = "8091"
	defMfUser          = "test@example.com"
	defMfPass          = "test"
	defMfAPIKey        = ""
	defMfBSURL         = "http://localhost:8202/things/configs"
	defMfWhiteListURL  = "http://localhost:8202/things/state"
	defMfCertsURL      = "http://localhost:8204"
	defProvisionCerts  = "false"
	defProvisionBS     = "true"
	defBSAutoWhitelist = "true"
	defBSContent       = ""
	defCertsHoursValid = "2400h"
	defCertsKeyBits    = "4096"

	envConfigFile       = "PD_PROVISION_CONFIG_FILE"
	envLogLevel         = "PD_PROVISION_LOG_LEVEL"
	envHTTPPort         = "PD_PROVISION_HTTP_PORT"
	envTLS              = "PD_PROVISION_ENV_CLIENTS_TLS"
	envServerCert       = "PD_PROVISION_SERVER_CERT"
	envServerKey        = "PD_PROVISION_SERVER_KEY"
	envMQTTURL          = "PD_PROVISION_MQTT_URL"
	envUsersLocation    = "PD_PROVISION_USERS_LOCATION"
	envThingsLocation   = "PD_PROVISION_THINGS_LOCATION"
	envMfUser           = "PD_PROVISION_USER"
	envMfPass           = "PD_PROVISION_PASS"
	envMfAPIKey         = "PD_PROVISION_API_KEY"
	envMfCertsURL       = "PD_PROVISION_CERTS_SVC_URL"
	envProvisionCerts   = "PD_PROVISION_X509_PROVISIONING"
	envMfBSURL          = "PD_PROVISION_BS_SVC_URL"
	envMfBSWhiteListURL = "PD_PROVISION_BS_SVC_WHITELIST_URL"
	envProvisionBS      = "PD_PROVISION_BS_CONFIG_PROVISIONING"
	envBSAutoWhiteList  = "PD_PROVISION_BS_AUTO_WHITELIST"
	envBSContent        = "PD_PROVISION_BS_CONTENT"
	envCertsHoursValid  = "PD_PROVISION_CERTS_HOURS_VALID"
	envCertsKeyBits     = "PD_PROVISION_CERTS_RSA_BITS"
)

var (
	errMissingConfigFile            = errors.New("missing config file setting")
	errFailLoadingConfigFile        = errors.New("failed to load config from file")
	errFailGettingAutoWhiteList     = errors.New("failed to get auto whitelist setting")
	errFailGettingCertSettings      = errors.New("failed to get certificate file setting")
	errFailGettingTLSConf           = errors.New("failed to get TLS setting")
	errFailGettingProvBS            = errors.New("failed to get BS url setting")
	errFailSettingKeyBits           = errors.New("failed to set rsa number of bits")
	errFailedToReadBootstrapContent = errors.New("failed to read bootstrap content from envs")
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	logger, err := logger.New(os.Stdout, cfg.Server.LogLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if cfgFromFile, err := loadConfigFromFile(cfg.File); err != nil {
		logger.Warn(fmt.Sprintf("Continue with settings from env, failed to load from: %s: %s", cfg.File, err))
	} else {
		// Merge environment variables and file settings.
		mergeConfigs(&cfgFromFile, &cfg)
		cfg = cfgFromFile
		logger.Info("Continue with settings from file: " + cfg.File)
	}

	SDKCfg := mfSDK.Config{
		BaseURL:           cfg.Server.ThingsLocation,
		BootstrapURL:      cfg.Server.MfBSURL,
		CertsURL:          cfg.Server.MfCertsURL,
		HTTPAdapterPrefix: "http",
		MsgContentType:    "application/json",
		TLSVerification:   cfg.Server.TLS,
	}
	SDK := mfSDK.NewSDK(SDKCfg)

	svc := provision.New(cfg, SDK, logger)
	svc = api.NewLoggingMiddleware(svc, logger)

	errs := make(chan error, 2)

	go startHTTPServer(svc, cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Provision service terminated: %s", err))
}

func startHTTPServer(svc provision.Service, cfg provision.Config, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.Server.HTTPPort)
	if cfg.Server.ServerCert != "" || cfg.Server.ServerKey != "" {
		logger.Info(fmt.Sprintf("Provision service started using https on port %s with cert %s key %s",
			cfg.Server.HTTPPort, cfg.Server.ServerCert, cfg.Server.ServerKey))
		errs <- http.ListenAndServeTLS(p, cfg.Server.ServerCert, cfg.Server.ServerKey, api.MakeHandler(svc))
		return
	}
	logger.Info(fmt.Sprintf("Provision service started using http on port %s", cfg.Server.HTTPPort))
	errs <- http.ListenAndServe(p, api.MakeHandler(svc))
}

func loadConfigFromFile(file string) (provision.Config, error) {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return provision.Config{}, errors.Wrap(errMissingConfigFile, err)
	}
	c, err := provision.Read(file)
	if err != nil {
		return provision.Config{}, errors.Wrap(errFailLoadingConfigFile, err)
	}
	return c, nil
}

func loadConfig() (provision.Config, error) {
	tls, err := strconv.ParseBool(pandas.Env(envTLS, defTLS))
	if err != nil {
		return provision.Config{}, errors.Wrap(errFailGettingTLSConf, err)
	}
	provisionX509, err := strconv.ParseBool(pandas.Env(envProvisionCerts, defProvisionCerts))
	if err != nil {
		return provision.Config{}, errors.Wrap(errFailGettingCertSettings, err)
	}
	provisionBS, err := strconv.ParseBool(pandas.Env(envProvisionBS, defProvisionBS))
	if err != nil {
		return provision.Config{}, errors.Wrap(errFailGettingProvBS, fmt.Errorf(" for %s", envProvisionBS))
	}

	autoWhiteList, err := strconv.ParseBool(pandas.Env(envBSAutoWhiteList, defBSAutoWhitelist))
	if err != nil {
		return provision.Config{}, errors.Wrap(errFailGettingAutoWhiteList, fmt.Errorf(" for %s", envBSAutoWhiteList))
	}
	if autoWhiteList && !provisionBS {
		return provision.Config{}, errors.New("Can't auto whitelist if auto config save is off")
	}
	keyBits, err := strconv.Atoi(pandas.Env(envCertsKeyBits, defCertsKeyBits))
	if err != nil && provisionX509 == true {
		return provision.Config{}, errFailSettingKeyBits
	}

	var content map[string]interface{}
	if c := pandas.Env(envBSContent, defBSContent); c != "" {
		if err = json.Unmarshal([]byte(c), content); err != nil {
			return provision.Config{}, errFailedToReadBootstrapContent
		}
	}

	cfg := provision.Config{
		Server: provision.ServiceConf{
			LogLevel:       pandas.Env(envLogLevel, defLogLevel),
			ServerCert:     pandas.Env(envServerCert, defServerCert),
			ServerKey:      pandas.Env(envServerKey, defServerKey),
			HTTPPort:       pandas.Env(envHTTPPort, defHTTPPort),
			MfBSURL:        pandas.Env(envMfBSURL, defMfBSURL),
			MfWhiteListURL: pandas.Env(envMfBSWhiteListURL, defMfWhiteListURL),
			MfCertsURL:     pandas.Env(envMfCertsURL, defMfCertsURL),
			MfUser:         pandas.Env(envMfUser, defMfUser),
			MfPass:         pandas.Env(envMfPass, defMfPass),
			MfAPIKey:       pandas.Env(envMfAPIKey, defMfAPIKey),
			ThingsLocation: pandas.Env(envThingsLocation, defThingsLocation),
			UsersLocation:  pandas.Env(envUsersLocation, defUsersLocation),
			TLS:            tls,
		},
		Certs: provision.Certs{
			HoursValid: pandas.Env(envCertsHoursValid, defCertsHoursValid),
			KeyBits:    keyBits,
		},
		Bootstrap: provision.Bootstrap{
			X509Provision: provisionX509,
			Provision:     provisionBS,
			AutoWhiteList: autoWhiteList,
			Content:       content,
		},

		// This is default conf for provision if there is no config file
		Channels: []provision.Channel{
			{
				Name:     "control-channel",
				Metadata: map[string]interface{}{"type": "control"},
			}, {
				Name:     "data-channel",
				Metadata: map[string]interface{}{"type": "data"},
			},
		},
		Things: []provision.Thing{
			{
				Name:     "thing",
				Metadata: map[string]interface{}{"external_id": "xxxxxx"},
			},
		},
	}

	cfg.File = pandas.Env(envConfigFile, defConfigFile)
	return cfg, nil
}

func mergeConfigs(dst, src interface{}) interface{} {
	d := reflect.ValueOf(dst).Elem()
	s := reflect.ValueOf(src).Elem()

	for i := 0; i < d.NumField(); i++ {
		dField := d.Field(i)
		sField := s.Field(i)
		switch dField.Kind() {
		case reflect.Struct:
			dst := dField.Addr().Interface()
			src := sField.Addr().Interface()
			m := mergeConfigs(dst, src)
			val := reflect.ValueOf(m).Elem().Interface()
			dField.Set(reflect.ValueOf(val))
		case reflect.Slice:
		case reflect.Bool:
			if dField.Interface() == false {
				dField.Set(reflect.ValueOf(sField.Interface()))
			}
		case reflect.Int:
			if dField.Interface() == 0 {
				dField.Set(reflect.ValueOf(sField.Interface()))
			}
		case reflect.String:
			if dField.Interface() == "" {
				dField.Set(reflect.ValueOf(sField.Interface()))
			}
		}
	}
	return dst
}
