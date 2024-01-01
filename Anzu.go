/* Takakura Anzu
 *
 *
 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Akizon77/TakakuraAnzu/data/sql/TakakuraAnzu"
	"github.com/Akizon77/TakakuraAnzu/data/sql/TakakuraAnzu/whitelist"
	messageLogger "github.com/Akizon77/TakakuraAnzu/log"
	"github.com/Akizon77/TakakuraAnzu/message"
	"github.com/Akizon77/TakakuraAnzu/rss"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	configPath = "./config.txt"
	version    = "1.3.4"
)

var (
	token    = ""
	interval = 10
	Debug    = false
)

func main() {
	//加载数据库
	TakakuraAnzu.LoadDatabase()
	loadConfig()
	messageLogger.Info("Takakura Anzu " + version)
	messageLogger.Info(fmt.Sprintf("读取到配置：Token: %s  Interval:%d", token, interval))
	debug := flag.Bool("debug", false, "")
	addPermission := flag.Bool("addPermission", false, "")
	flag.Parse()
	if *debug {
		messageLogger.EnableDebugMode()
	}
	if *addPermission {
		whitelist.Add(1977354088, "开发者")
		whitelist.Add(-1001942218297, "世界，你好")
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
				if update.Message.From.ID != 1977354088 && messageLogger.DebugMode {
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

func loadConfig() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, _ := os.Create(configPath)
		file.WriteString("TOKEN = YOUR_TOKEN_HERE\nINTERVAL = 10")
		token = "YOUR_TOKEN_HERE"
		interval = 10
		defer file.Close()
		return
	}
	file, err := os.Open(configPath)
	if err != nil {
		messageLogger.Error("无法读取配置文件 "+configPath, err)
		return
	}
	scanner := bufio.NewScanner(file)

	// 逐行读取配置文件
	for scanner.Scan() {
		line := scanner.Text()
		// 使用等号分割每一行的键值对
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			messageLogger.Warn("配置文件中存在无效行")
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "TOKEN":
			token = value
		case "INTERVAL":
			interval, err = strconv.Atoi(value)
			if err != nil {
				messageLogger.Warn("无法读取INTERVAL,将使用默认值10分钟")
				interval = 10
			}
		default:

		}
	}

	if err := scanner.Err(); err != nil {
		messageLogger.Error("读取配置文件失败:", err)
		return
	}

}
