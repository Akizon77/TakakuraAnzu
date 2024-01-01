/* Takakura Anzu
 *
 *
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Akizon77/TakakuraAnzu/config"
	"github.com/Akizon77/TakakuraAnzu/data/sql/TakakuraAnzu"
	"github.com/Akizon77/TakakuraAnzu/data/sql/TakakuraAnzu/whitelist"
	messageLogger "github.com/Akizon77/TakakuraAnzu/log"
	"github.com/Akizon77/TakakuraAnzu/message"
	"github.com/Akizon77/TakakuraAnzu/rss"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	configPath = "./config.txt"
)

var (
	Version string = ""
)

func main() {
	token := config.Config.Token
	interval := config.Config.Interval
	//加载数据库
	TakakuraAnzu.LoadDatabase()
	messageLogger.Info("Takakura Anzu " + Version)
	messageLogger.Info(fmt.Sprintf("读取到配置：Token: %s  Interval:%d", token, interval))
	debug := flag.Bool("debug", false, "")
	addPermission := flag.Bool("addPermission", false, "")
	flag.Parse()
	if *debug {
		messageLogger.EnableDebugMode()
	}
	if *addPermission {
		whitelist.Add(int64(config.Config.Owner), "开发者")
		return
	}
	// 测试用 记得删

	// 新建bot 使用NewBotAPI函数，参数是Bot Token
	bot, err := tgbotapi.NewBotAPI(token)
	//无法实例化bot中断程序
	if err != nil {
		log.Println("无法创建Telegram Bot实例")
		log.Panic(err)
		return
	}
	messageLogger.Info("实例已启动！")
	// 设置Bot的更新模式为长轮询
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// 获取Telegram Bot的更新通道
	updates, _ := bot.GetUpdatesChan(u)
	duration, _ := time.ParseDuration(fmt.Sprintf("%dm", interval))

	// 定义一个定时器
	ticker := time.NewTicker(duration)
	// 循环处理Telegram Bot的更新
	for {
		select {
		case update := <-updates:
			// 处理接收到的消息
			if update.Message != nil {
				go messageLogger.Log(update.Message)
				if update.Message.From.ID != config.Config.Owner && messageLogger.DebugMode {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "开发者正在调试程序，请稍后再试")
					messageLogger.SendMsg(msg, bot)
				} else {
					message.RunCommand(update.Message, bot)
				}

			}
		case <-ticker.C:
			// 定时从RSS源获取最新的内容
			go rss.RefreshAll(bot)
		}
	}
}
