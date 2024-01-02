package config

type ConfigStruct struct {
	Token                 string `json:"token"`
	Interval              int    `json:"interval"`
	Owner                 int    `json:"owner"`
	MinecraftServer       string `json:"minecraft_server"`
	DDNS_Interface        string `json:"DDNS_Interface"`
	DDNS_Email            string `json:"DDNS_Email"`
	DDNS_APIKEY           string `json:"DDNS_APIKEY"`
	EnableQQBot           bool   `json:"useQQ"`
	QQ_App_id             uint64 `json:"QQ_App_Id"`
	QQ_Secret             string `json:"QQ_Secret"`
	QQ_Token              string `json:"QQ_Token"`
	QQ_Trans_To_TG_ChatID int64  `json:"transTo"`
}
