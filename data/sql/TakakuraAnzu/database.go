package TakakuraAnzu

import (
	"database/sql"
	"github.com/Akizon77/TakakuraAnzu/data/sql/TakakuraAnzu/whitelist"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const Path = "./data/TakakuraAnzu.db"
const Dir = "./data"

var Db *sql.DB = nil

func LoadDatabase() {
	// 如果文件不存在，创建新的数据库文件
	if _, err := os.Stat(Path); os.IsNotExist(err) {
		os.MkdirAll("./data", os.ModePerm)
		file, err := os.Create(Path)
		if err != nil {
			log.Println("无法创建数据库文件", Path, err)
			panic(err)
		}
		err = file.Close()
		if err != nil {
			log.Println("无法结束占用", Path, err)
			panic(err)
		}
	}

	// 连接到SQLite数据库
	var err error
	// 打开数据库文件
	Db, err = sql.Open("sqlite3", Path)
	if err != nil {
		log.Println("无法打开数据库文件:", Path)
		panic(err)
	}
	notify()
}
func notify() {
	whitelist.LoadWhitelist(Db)
}
