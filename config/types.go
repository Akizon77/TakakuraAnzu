package config

type ConfigStruct struct {
	Token           string `json:"token"`
	Interval        int    `json:"interval"`
	Owner           int64  `json:"owner"`
	MinecraftServer string `json:"minecraft_server"`
	DDNS_Interface  string `json:"DDNS_Interface"`
	DDNS_Email      string `json:"DDNS_Email"`
	DDNS_APIKEY     string `json:"DDNS_APIKEY"`
	WebToken        string `json:"web_token"`
	WebApiEndpoint  string `json:"web_api_endpoint"`
}
