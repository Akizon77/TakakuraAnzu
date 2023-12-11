package rss

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Akizon77/TakakuraAnzu/data/sql/rss_subs"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/rss"
	"log"
	"regexp"
	"strconv"
	"strings"
)

//var Db *sql.DB = rss_subs.Db

type RssConfig struct {
	ChatID int64    `json:"chat_id"`
	Rss    []string `json:"rss"`
	Cache  []string `json:"cache"`
}

// 向数据库中添加一个表 格式为 rss_subs<ChatID>
func AddRssForChatID(chatID int64, link string) error {
	// 正则表达式 匹配网址
	re := regexp.MustCompile(`^(http(s)?:\/\/)?(www\.)?[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`)
	regResult := re.FindString(link)
	if regResult == "" {
		return errors.New(fmt.Sprintf("请认真的发送RSS订阅地址！"))
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS 'rss_subs%d' (links TEXT PRIMARY KEY,cache TEXT)", chatID)
	//创建表（如果不存在）
	_, err := rss_subs.Db.Exec(query)
	if err != nil {
		return errors.New(fmt.Sprintf("SQL create table: %s\nQuery: %s", err.Error(), query))
	}

	query = fmt.Sprintf("INSERT INTO 'rss_subs%d' (links,cache) VALUES ('%s',NULL);", chatID, link)
	_, err = rss_subs.Db.Exec(query)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New(fmt.Sprintf("同一个地址可以不用发两遍啦"))
		}
		return errors.New(fmt.Sprintf("SQL insert: %s\nQuery: %s", err.Error(), query))
	}
	return nil
}

// RemoveRssForChatID 删除订阅链接
func RemoveRssForChatID(chatID int64, link string, bot *tgbotapi.BotAPI) error {
	table := fmt.Sprintf("rss_subs%d", chatID)
	query := fmt.Sprintf("DELETE FROM '%s' WHERE links = '%s'", table, link)
	_, err := rss_subs.Db.Exec(query)
	if err != nil {
		if strings.Contains(err.Error(), "no such") {
			return errors.New("没有这条订阅哦")
		}
		return errors.New(fmt.Sprintf("SQL delete row: %s\nQuery: %s", err.Error(), query))
	}
	return nil
}

func ListAllSubs(chatID int64) (string, error) {
	table := fmt.Sprintf("rss_subs%d", chatID)
	query := fmt.Sprintf("SELECT links FROM '%s'", table)
	rows, err := rss_subs.Db.Query(query)
	//会导致markdown解析失败
	if err != nil {
		log.Println(err)
		return "", errors.New(fmt.Sprintf("SQL query: %s\nQuery: %s", err.Error(), query))
	}
	defer rows.Close()
	var result = ""
	for rows.Next() {
		var link string
		err = rows.Scan(&link)
		if err != nil {
			log.Println(err)
			return "", errors.New(fmt.Sprintf("SQL walk row: %s", err.Error()))
		}
		result += "`" + link + "`\n"
	}
	return result, nil
}
func GetAllSubs(chatID int64) ([]string, error) {
	table := fmt.Sprintf("rss_subs%d", chatID)
	query := fmt.Sprintf("SELECT links FROM '%s'", table)
	rows, err := rss_subs.Db.Query(query)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("SQL query: %s\nQuery: %s", err.Error(), query))
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var link string
		err = rows.Scan(&link)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("SQL walk row: %s", err.Error()))
		}
		result = append(result, link)
	}
	return result, nil
}
func SQLRun(query string) (sql.Result, error) {
	return rss_subs.Db.Exec(query)
}

// 刷新并发送
func RefreshAndSend(chatID int64, bot *tgbotapi.BotAPI) error {
	items, err := Refresh(chatID, bot)
	SendRssUpdate(chatID, items, bot)
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "空数据") {
			return nil
		}
		return err
	}
	if len(items) == 0 {
		return errors.New("暂时还没有更新哦")
	}
	return nil
}
func SendRssUpdate(chatID int64, items []rss.Item, bot *tgbotapi.BotAPI) {
	for _, item := range items {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(item.Title+"\n"+item.Link))
		_, err := bot.Send(msg)
		if err != nil {
			return
		}
	}
}
func Refresh(chatID int64, bot *tgbotapi.BotAPI) ([]rss.Item, error) {
	links, err := GetAllSubs(chatID)
	var newItems []rss.Item
	var allNewItems []rss.Item

	if err != nil {
		log.Println(err)
		if bot != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("无法更新"+err.Error()))
			bot.Send(msg)
		}
		return nil, err
	}
	//判断是否是空的，是的话顺便删除
	if len(links) == 0 {
		log.Println(chatID, " 已使用但是无订阅，尝试删除空的表")
		err = DeleteTable(chatID)
		if err == nil {
			return nil, errors.New("已使用但是无订阅，正在删除空数据")
		}
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	// 获取新的订阅链接
	// 获取所有订阅地址
	// 拿到每一个的新订阅
	// 对比cache 看看有哪些新增的
	// 写入每一个的cache
	for _, link := range links {
		var newCache string
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(link)
		if err != nil {
			log.Println("无法解析RSS源:", err)
			log.Println(chatID, "出现错误！\n尝试解析RSS失败\n地址：", "\n详细信息：\n", link, err.Error())
			if bot != nil {
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("无法更新 %s\n 因为%s", link, err.Error()))
				bot.Send(msg)
			}
			//return nil, errors.New("出现错误！\n尝试解析RSS失败\n地址：" + link + "\n详细信息：\n" + err.Error())
		}
		cache, _ := getcache(chatID, link)
		//那Link下单Cache 用函数查询SQL
		//拿到每一个item
		//对比和cache的不同
		if feed == nil {
			continue
		}
		if feed.Items == nil {
			continue
		}
		// 遍历获取到的RSS 并添加到allNewItems
		// newItems 仅在一次循环中生效
		for _, item := range feed.Items {
			//如果不包含，那么就是新的文章
			if !strings.Contains(cache, item.Link) {
				allNewItems = append(allNewItems, rss.Item{Link: item.Link, Title: item.Title})
				newItems = append(newItems, rss.Item{Link: item.Link, Title: item.Title})
			}
			//放在if外，放在里面会导致下次刷新重复获取
			newCache += item.Link + ","

		}
		//写入数据库
		table := fmt.Sprintf("rss_subs%d", chatID)
		query := fmt.Sprintf("UPDATE '%s' SET cache = '%s' WHERE links = '%s'", table, newCache, link)
		_, errq := rss_subs.Db.Exec(query)
		if errq != nil {
			log.Println(errq)
			return nil, errors.New("SQL update: " + errq.Error())
		}
		//清空进行下一次
		newItems = nil
	}

	return allNewItems, nil
}
func RefreshAll(bot *tgbotapi.BotAPI) error {
	log.Println("开始处理所有人的订阅")
	users := AllUser()
	for _, user := range users {
		items, err := Refresh(user, nil)
		SendRssUpdate(user, items, bot)
		if err != nil {
			log.Println("无法处理 " + strconv.FormatInt(user, 10) + "的订阅")
		}
	}
	return nil
}

func getcache(chatID int64, link string) (string, error) {
	table := fmt.Sprintf("rss_subs%d", chatID)
	query := fmt.Sprintf("SELECT cache FROM '%s' WHERE links = '%s'", table, link)
	var resule string
	err := rss_subs.Db.QueryRow(query).Scan(&resule)
	if err != nil {
		return "", errors.New(fmt.Sprintf("SQL query: %s\nQuery: %s", err.Error(), query))
	}

	return resule, nil

}
func AllUser() []int64 {
	var users []int64
	rows, err := rss_subs.Db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	re := regexp.MustCompile(`[-]?(\d+)`)
	var tableName string
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Println(err)
		}
		tableName = re.FindString(tableName)
		num, err := strconv.ParseInt(tableName, 10, 64)
		users = append(users, num)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return users
}
func DeleteTable(chatID int64) error {
	table := fmt.Sprintf("rss_subs%d", chatID)
	query := fmt.Sprintf("DROP TABLE '%s'", table)
	_, err := rss_subs.Db.Exec(query)
	if err != nil {
		log.Println(err)
		return errors.New(fmt.Sprintf("SQL delete table: %s\nQuery: %s", err.Error(), query))
	}

	return nil
}
