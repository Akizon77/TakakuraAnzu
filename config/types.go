package config

type ConfigStruct struct {
	Token           string `json:"token"`
	Interval        int    `json:"interval"`
	Owner           int    `json:"owner"`
	MinecraftServer string `json:"minecraft_server"`
	DDNS_Interface  string `json:"DDNS_Interface"`
	DDNS_Email      string `json:"DDNS_Email"`
	DDNS_APIKEY     string `json:"DDNS_APIKEY"`
}
