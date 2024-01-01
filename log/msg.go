package log

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

var DebugMode bool = false

// 发送消息并记录，如出现空地址请排查是否启用的 Markdown 语法但标记不能匹配的清空
func SendMsg(c tgbotapi.Chattable, bot *tgbotapi.BotAPI) tgbotapi.Message {
	message, err := bot.Send(c)
	if message.MessageID == 0 {
		Error("无法发送消息", err)
		return message
	}
	if err != nil {
		Error("呀呀呀，发送 "+message.Text+" 时错误 详情", err)
	}
	var header string = ""
	if message.Chat.IsGroup() || message.Chat.IsSuperGroup() {
		header = "发送给群组 " + message.Chat.Title + " ："
	} else {
		header = "发送给 " + message.Chat.FirstName + message.Chat.LastName + "(@" + message.Chat.UserName + ")："
	}
	Info(header + message.Text)
	return message
}

func Log(message *tgbotapi.Message) {
	var header string
	//群组的话，前缀改成群
	if message.Chat.IsSuperGroup() || message.Chat.IsGroup() {
		header = "收到群组 " + message.Chat.Title + " 用户 " + message.From.FirstName + message.From.LastName + "(@" + message.From.UserName + ")" + " 的"
	} else {
		header = "收到 " + message.From.FirstName + message.From.LastName + "(@" + message.From.UserName + ")" + " 的"
	}
	if message.Sticker != nil {
		Info(header + "贴纸，对应Emoji是：" + message.Sticker.Emoji + ",贴纸包：" + message.Sticker.SetName)
		return
	}
	if message.Text != "" {
		Info(header + "文字信息：" + message.Text)
		return
	}
	if message.Photo != nil {
		Info(header + "图片")
	}
	if message.Document != nil {
		Info(header + "文件: " + message.Document.FileName)
	}
	Warn(header + "消息,无法预览")
}

func logs(msg string, header string) {
	log.Println(fmt.Sprintf("| %s | %s", header, msg))
}
func EnableDebugMode() {
	DebugMode = true
	Debug("Debug Mode Enabled")
}

func Debug(message string) {
	if DebugMode {
		logs(message, "debug")
	}
}
func Warn(message string) {
	logs(message, "warn")
}
func Info(message string) {
	logs(message, "info")
}
func Panic(message string) {
	logs(message, "panic")
}
func Error(message string, e error) {
	logs(message+"  "+e.Error(), "error")
}
