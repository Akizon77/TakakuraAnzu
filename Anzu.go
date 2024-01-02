/* Takakura Anzu
 *
 *
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Akizon77/TakakuraAnzu/log/logger"
	"github.com/Akizon77/TakakuraAnzu/qqbot/command"
	"github.com/tencent-connect/botgo/dto"
	qevent "github.com/tencent-connect/botgo/event"
	qtoken "github.com/tencent-connect/botgo/token"
	qws "github.com/tencent-connect/botgo/websocket"
	"log"
	"strings"
	"time"

	"github.com/Akizon77/TakakuraAnzu/config"
	"github.com/Akizon77/TakakuraAnzu/data/sql/TakakuraAnzu"
	"github.com/Akizon77/TakakuraAnzu/data/sql/TakakuraAnzu/whitelist"
	messageLogger "github.com/Akizon77/TakakuraAnzu/log"
	"github.com/Akizon77/TakakuraAnzu/message"
	"github.com/Akizon77/TakakuraAnzu/rss"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	qqbotapi "github.com/tencent-connect/botgo"
)

const (
	configPath = "./config.txt"
)

var (
	token           = config.Config.Token
	interval        = config.Config.Interval
	Version  string = ""
)

func main() {
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
	go loadTGBot()
	if config.Config.EnableQQBot {
		messageLogger.Info("正在启动QQ机器人")
		go loadQQBot()
	}
	// 永远卡死主线程
	for {
		ch := make(chan int)
		<-ch
	}
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
	message.InitMessageTransfer(bot)
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
func loadQQBot() {
	token := qtoken.BotToken(config.Config.QQ_App_id, config.Config.QQ_Token)
	api := qqbotapi.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Printf("%+v, err:%v", ws, err)
	}

	//监听哪类事件就需要实现哪类的 handler，定义：websocket/event_handler.go
	var atMessage qevent.ATMessageEventHandler = func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		spl := strings.Split(data.Content, " ")
		if len(spl) == 1 {
			return nil
		}
		if len(spl) == 2 && spl[1] == "" {
			return nil
		}
		err := command.RunCommand((*dto.Message)(data))
		if err != nil {
			messageLogger.Error("无法发送", err)
		}
		return nil
	}
	var allMessage qevent.MessageEventHandler = func(event *dto.WSPayload, data *dto.WSMessageData) error {
		//fmt.Println(event, data)
		//messageLogger.Debug(data.Author.Username)
		messageLogger.Info(fmt.Sprint("收到QQ用户", data.Author.Username, "的消息：", data.Content))
		if strings.Contains(data.Content, "/") {
			return nil
		}
		go message.Transfer(data.Author.Username, "QQ", data.Content, config.Config.QQ_Trans_To_TG_ChatID)
		return nil
	}

	qqbotapi.SetLogger(logger.NewLogger())
	intent := qws.RegisterHandlers(allMessage, atMessage)
	// 启动 session manager 进行 ws 连接的管理，如果接口返回需要启动多个 shard 的连接，这里也会自动启动多个
	_ = qqbotapi.NewSessionManager().Start(ws, token, &intent)
}
