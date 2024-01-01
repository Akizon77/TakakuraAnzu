package minecraft

type Server struct {
	Online      bool          `json:"online"`
	Host        string        `json:"host"`
	Port        int64         `json:"port"`
	IPAddress   string        `json:"ip_address"`
	EULABlocked bool          `json:"eula_blocked"`
	RetrievedAt int64         `json:"retrieved_at"`
	ExpiresAt   int64         `json:"expires_at"`
	SrvRecord   interface{}   `json:"srv_record"`
	Version     Version       `json:"version"`
	Players     Players       `json:"players"`
	MOTD        MOTD          `json:"motd"`
	Icon        string        `json:"icon"`
	Mods        []interface{} `json:"mods"`
	Software    interface{}   `json:"software"`
	Plugins     []interface{} `json:"plugins"`
}

type MOTD struct {
	Raw   string `json:"raw"`
	Clean string `json:"clean"`
	HTML  string `json:"html"`
}

type Players struct {
	Online int64         `json:"online"`
	Max    int64         `json:"max"`
	List   []interface{} `json:"list"`
}

type Version struct {
	NameRaw   string `json:"name_raw"`
	NameClean string `json:"name_clean"`
	NameHTML  string `json:"name_html"`
	Protocol  int64  `json:"protocol"`
}
