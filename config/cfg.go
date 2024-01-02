package config

import (
	"encoding/json"
	log "github.com/Akizon77/TakakuraAnzu/log"
	"os"
)

const (
	configPath = "./config.json"
)

var (
	Config = &ConfigStruct{Interval: 10}
)

func init() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		bytes, _ := json.Marshal(Config)
		file, _ := os.Create(configPath)
		file.Write(bytes)
		defer file.Close()
		return
	}
	cfg, err := os.ReadFile(configPath)
	if err != nil {
		log.Error("无法读取配置文件 "+configPath, err)
		return
	}
	err = json.Unmarshal(cfg, &Config)
	if err != nil {
		log.Error("无法解析配置文件 "+configPath, err)
		return
	}
	log.Debug("Config：" + string(cfg))

}
