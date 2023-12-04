package config

type Secrets struct {
	Integrations map[string]IntegrationSecret `json:"integrations"`
}

type IntegrationSecret struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type TickerMaster struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type Twitter struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type Oauth2 struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}
