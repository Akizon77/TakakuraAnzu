package rss_subs

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var Db *sql.DB
var rssfile = "./data/rss.db"

// 打开SQL数据库
func init() {
	// 如果文件不存在，创建新的数据库文件
	if _, err := os.Stat(rssfile); os.IsNotExist(err) {
		os.MkdirAll("./data", os.ModePerm)
		file, err := os.Create(rssfile)
		if err != nil {
			log.Println("无法创建数据库RSS文件")
			panic(err)
		}
		err = file.Close()
		if err != nil {
			log.Println("无法结束占用RSS")
			panic(err)
		}
	}

	// 连接到SQLite数据库
	var err error
	// 打开数据库文件
	Db, err = sql.Open("sqlite3", rssfile)
	if err != nil {
		log.Println("无法打开数据库文件:", rssfile)
		panic(err)
	}
}

/*
func Exec(query string, args ...any) (sql.Result, error) {
	return Db.Exec(query, args)
}
func Query(query string, args ...any) (*sql.Rows, error) {
	return Db.Query(query, args)
}
func QueryRow(query string, args ...any) *sql.Row {
	return Db.QueryRow(query, args)
}
func QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return Db.QueryRowContext(ctx, query, args)
}
func Prepare(query string) (*sql.Stmt, error) {
	return Db.Prepare(query)
}
*/
