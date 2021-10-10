package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloustone/pandas/kuiper"
)

const kuiperEndpoint = "kuiper"

type KuiperPluginType string

const (
	KuiperPluginSink   KuiperPluginType = "sink"
	KuiperPluginSource KuiperPluginType = "source"
)

// CreateKuiperStream create a stream in kuiper
func (sdk mfSDK) CreateKuiperStream(stream string, token string) (string, error) {
	data, err := json.Marshal(stream)
	if err != nil {
		return "", ErrInvalidArgs
	}

	endpoint := "streams"
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return string(body), ErrInvalidArgs
		case http.StatusForbidden:
			return string(body), ErrUnauthorized
		default:
			return string(body), ErrFailedCreation
		}
	}

	return "", nil
}

// CreateKuiperRule register a rule
func (sdk mfSDK) CreateKuiperRule(desc kuiper.Rule, token string) (string, error) {
	data, err := json.Marshal(desc)
	if err != nil {
		return "", ErrInvalidArgs
	}

	endpoint := "rules"
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return string(body), ErrInvalidArgs
		case http.StatusForbidden:
			return string(body), ErrUnauthorized
		default:
			return string(body), ErrFailedCreation
		}
	}
	return "", nil
}

// KuiperStream retrurn specified stream info
func (sdk mfSDK) KuiperStream(streamName string, token string) (string, error) {
	endpoint := fmt.Sprintf("streams/%s", streamName)
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return string(body), ErrInvalidArgs
		case http.StatusForbidden:
			return string(body), ErrUnauthorized
		default:
			return string(body), ErrFailedCreation
		}
	}
	return string(body), nil
}

// KuiperRule return specified rule info
func (sdk mfSDK) KuiperRule(ruleName string, token string) (string, error) {
	endpoint := fmt.Sprintf("rules/%s", ruleName)
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return string(body), ErrInvalidArgs
		case http.StatusForbidden:
			return string(body), ErrUnauthorized
		default:
			return string(body), ErrFailedCreation
		}
	}
	return string(body), nil

}

// KuiperStreams retrurn streams info
func (sdk mfSDK) KuiperStreams(token string) (string, error) {
	endpoint := "streams"
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return string(body), ErrInvalidArgs
		case http.StatusForbidden:
			return string(body), ErrUnauthorized
		default:
			return string(body), ErrFailedCreation
		}
	}
	return string(body), nil

}

// KuiperRule return specified rule info
func (sdk mfSDK) KuiperRules(token string) (string, error) {
	endpoint := "rules"
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return string(body), ErrInvalidArgs
		case http.StatusForbidden:
			return string(body), ErrUnauthorized
		default:
			return string(body), ErrFailedCreation
		}
	}
	return string(body), nil

}

// KupirRuleStatus return a rule's status
func (sdk mfSDK) KuiperRuleStatus(ruleName, token string) (string, error) {
	endpoint := fmt.Sprintf("rules/%s/status", ruleName)
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return string(body), ErrInvalidArgs
		case http.StatusForbidden:
			return string(body), ErrUnauthorized
		default:
			return string(body), ErrFailedCreation
		}
	}
	return string(body), nil
}

func buildPluginEndpoint(pluginType KuiperPluginType) string {
	if pluginType == KuiperPluginSource {
		return "plugin/sources"
	} else {
		return "plugin/sinks"
	}
}

// KuiperPlugin return specified plugin info
func (sdk mfSDK) KuiperPlugins(pluginType KuiperPluginType, token string) (string, error) {
	endpoint := buildPluginEndpoint(pluginType)
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return string(body), ErrInvalidArgs
		case http.StatusForbidden:
			return string(body), ErrUnauthorized
		default:
			return string(body), ErrFailedCreation
		}
	}
	return string(body), nil
}

// KuiperPlugin return specified plugin info
func (sdk mfSDK) KuiperPlugin(pluginType KuiperPluginType, id string, token string) (string, error) {
	endpoint := buildPluginEndpoint(pluginType)
	endpoint = fmt.Sprintf("%s/%s", endpoint, id)

	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return string(body), ErrInvalidArgs
		case http.StatusForbidden:
			return string(body), ErrUnauthorized
		default:
			return string(body), ErrFailedCreation
		}
	}
	return string(body), nil
}

// DeleteKuiperStream remove kuiper stream
func (sdk mfSDK) DeleteKuiperStream(streamName string, token string) error {
	endpoint := fmt.Sprintf("streams/%s", streamName)
	url := createURL(sdk.baseURL, sdk.thingsPrefix, endpoint)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return ErrInvalidArgs
		case http.StatusForbidden:
			return ErrUnauthorized
		default:
			return ErrFailedUpdate
		}
	}
	return nil
}

// DeleteKuiperRule remove kuiper rule
func (sdk mfSDK) DeleteKuiperRule(ruleName string, token string) error {
	endpoint := fmt.Sprintf("rules/%s", ruleName)
	url := createURL(sdk.baseURL, sdk.thingsPrefix, endpoint)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return ErrInvalidArgs
		case http.StatusForbidden:
			return ErrUnauthorized
		default:
			return ErrFailedUpdate
		}
	}
	return nil
}

// DeleteKuiperPlugin remove kuiper plugin
func (sdk mfSDK) DeleteKuiperPlugin(desc kuiper.Plugin, token string) error {
	return ErrNotFound
}

// StartKuiperRule start an already existed rule in kuiper
func (sdk mfSDK) StartKuiperRule(ruleName string, token string) error {
	endpoint := fmt.Sprintf("rules/%s/start", ruleName)
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return ErrInvalidArgs
		case http.StatusForbidden:
			return ErrUnauthorized
		default:
			return ErrFailedCreation
		}
	}
	return nil
}

// StopKuiperRule stop an already existed rule in kuiper
func (sdk mfSDK) StopKuiperRule(ruleName string, token string) error {
	endpoint := fmt.Sprintf("rules/%s/stop", ruleName)
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return ErrInvalidArgs
		case http.StatusForbidden:
			return ErrUnauthorized
		default:
			return ErrFailedCreation
		}
	}
	return nil

}

// RestartKuiperRule restart an already existed rule in kuiper
func (sdk mfSDK) RestartKuiperRule(ruleName string, token string) error {
	endpoint := fmt.Sprintf("rules/%s/restart", ruleName)
	url := createURL(sdk.baseURL, sdk.kuiperPrefix, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return ErrInvalidArgs
		case http.StatusForbidden:
			return ErrUnauthorized
		default:
			return ErrFailedCreation
		}
	}
	return nil

}
