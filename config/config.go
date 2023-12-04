package config

type Secrets struct {
	Integrations Integration `json:"integrations"`
}

type Integration struct {
	TickerMaster    APIKeySecret              `json:"ticketmaster"`
	Twitter         APIKeySecret              `json:"twitter"`
	TwitterAccounts map[string]TwitterAccount `json:"twitterAccounts"`
	Oauth2          Oauth2                    `json:"oauth2"`
}

type APIKeySecret struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type TwitterAccount struct {
	BearerToken       string `json:"bearerToken"`
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
}

type Oauth2 struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}
