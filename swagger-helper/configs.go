package swagger_helper

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mainflux/mainflux/logger"
)

type DownstreamSwagger struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	SwaggerUrl string `json:"swaggerUrl"`
}

type SwaggerInfo struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type SwaggerHelperConfigs struct {
	DownstreamSwaggers []DownstreamSwagger `json:"downstreamSwaggers"`
	Info               SwaggerInfo         `json:"info"`
}

// LoadDownstreamSwaggers load downstream swagger from config file
func LoadDownstreamSwaggers(fullFilePath string, logger logger.Logger) (SwaggerHelperConfigs, error) {
	buf, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		logger.Debug("open downstream swagger helper config file failed")
		return SwaggerHelperConfigs{}, err
	}
	configs := SwaggerHelperConfigs{}
	if err := json.Unmarshal(buf, &configs); err != nil {
		logger.Debug("illegal downstream swagger config file")
		return SwaggerHelperConfigs{}, err
	}
	return configs, nil
}
