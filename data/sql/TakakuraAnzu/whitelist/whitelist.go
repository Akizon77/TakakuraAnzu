package whitelist

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Akizon77/TakakuraAnzu/data/sql/rss_subs"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

var (
	whitelistFile         = "./data/config.db"
	db            *sql.DB = nil
	table                 = "whitelist"
)

func LoadWhitelist(d *sql.DB) {
	db = d //引入本文件
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS '%s' (id BIGINT PRIMARY KEY,name TEXT)", table)
	_, err := db.Exec(query)
	if err != nil {
		log.Println("无法创建表", table)
	}
}
func IsWhitelist(chatID int64) bool {
	query := fmt.Sprintf("SELECT COUNT(*) FROM '%s' WHERE id = ?", table)

	var count int
	err := db.QueryRow(query, chatID).Scan(&count)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	} else {
		return false
	}

}

func Add(chatID int64, name string) error {
	query := fmt.Sprintf("INSERT INTO '%s' (id,name) VALUES ('%d','%s');", table, chatID, name)
	_, err := db.Exec(query)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New(fmt.Sprintf("这个ID (%d) 已经添加到了白名单了(ALREADY EXISTS)", chatID))
		}
		return err
	}
	log.Println("添加", chatID, "到白名单")
	return nil
}
func Clear() error {
	query := fmt.Sprintf("DROP TABLE '%s'", table)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS '%s' (id BIGINT PRIMARY KEY,name TEXT)", table)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	log.Println("表", table, "已清空")
	return nil
}
func Remove(chatID int64) error {
	query := fmt.Sprintf("DELETE FROM '%s' WHERE id = '%d'", table, chatID)
	_, err := rss_subs.Db.Exec(query)
	if err != nil {
		if strings.Contains(err.Error(), "no such") {
			return errors.New("这个ID本身就不在白名单哦")
		}
		return errors.New(fmt.Sprintf("SQL delete row: %s\nQuery: %s", err.Error(), query))
	}
	return nil
}
