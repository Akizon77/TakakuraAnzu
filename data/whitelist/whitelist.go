package whitelist

import (
	"database/sql"
	"errors"
	"github.com/Akizon77/TakakuraAnzu/config"
)

var (
	db    *sql.DB = nil
	table         = "whitelist"
)

func IsWhitelist(chatID int64) bool {
	if chatID == config.Config.Owner {
		return true
	}
	return false
}

func Add(chatID int64, name string) error {
	return errors.New("v2.0已取消白名单规则")
}
func Clear() error {
	return errors.New("v2.0已取消白名单规则")
}
func Remove(chatID int64) error {
	return errors.New("v2.0已取消白名单规则")
}
