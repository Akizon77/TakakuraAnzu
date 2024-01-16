/* Takakura Anzu
 *
 *
 */
package main

import (
	"flag"
	"fmt"
	"github.com/Akizon77/TakakuraAnzu/config"
	messageLogger "github.com/Akizon77/TakakuraAnzu/log"
	"github.com/Akizon77/TakakuraAnzu/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const (
	configPath = "./config.txt"
)

var (
	token           = config.Config.Token
	interval        = config.Config.Interval
	Version  string = "1.5.0"
)

func main() {
	//加载数据库
	messageLogger.Info("Takakura Anzu " + Version)
	messageLogger.Info(fmt.Sprintf("读取到配置：Token: %s  Interval:%d", token, interval))
	debug := flag.Bool("debug", false, "")
	flag.Parse()
	if *debug {
		messageLogger.EnableDebugMode()
	}

	loadTGBot()

}
func loadTGBot() {
	// 新建bot 使用NewBotAPI函数，参数是Bot Token
	bot, err := tgbotapi.NewBotAPI(token)
	//无法实例化bot中断程序
	if err != nil {
		log.Println("无法创建Telegram Bot实例,请检查配置文件config.json")
		log.Panic(err)
		return
	}
	messageLogger.Info("实例已启动！")
	// 设置Bot的更新模式为长轮询
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// 获取Telegram Bot的更新通道
	updates, _ := bot.GetUpdatesChan(u)
	// 循环处理Telegram Bot的更新
	for {
		select {
		case update := <-updates:
			// 处理接收到的消息
			if update.Message != nil {
				go messageLogger.Log(update.Message)
				message.RunCommand(update.Message, bot)
			}
		}
	}
}
