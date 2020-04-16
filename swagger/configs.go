package swagger

type DownstreamSwagger struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	SwaggerUrl string `json:"swaggerUrl"`
}

type SwaggerInfo struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type Configs struct {
	DownstreamSwaggers []DownstreamSwagger `json:"downstreamSwaggers"`
	Info               SwaggerInfo         `json:"info"`
}
