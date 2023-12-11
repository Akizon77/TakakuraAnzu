package log

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

// 发送消息并记录，如出现空地址请排查是否启用的 Markdown 语法但标记不能匹配的清空
func SendMsg(c tgbotapi.Chattable, bot *tgbotapi.BotAPI) tgbotapi.Message {
	message, err := bot.Send(c)
	if message.MessageID == 0 {
		log.Println("无法发送消息\n" + err.Error())
		return message
	}
	if err != nil {
		log.Println("呀呀呀，发送 " + message.Text + " 时错误 详情\n" + err.Error())
	}
	var header string = ""
	if message.Chat.IsGroup() || message.Chat.IsSuperGroup() {
		header = "发送给群组 " + message.Chat.Title + " ："
	} else {
		header = "发送给 " + message.Chat.FirstName + message.Chat.LastName + "(@" + message.Chat.UserName + ")："
	}
	log.Println(header + message.Text)
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
		log.Println(header + "贴纸，对应Emoji是：" + message.Sticker.Emoji + ",贴纸包：" + message.Sticker.SetName)
		return
	}
	if message.Text != "" {
		log.Println(header + "文字信息：" + message.Text)
		return
	}
	if message.Photo != nil {
		log.Println(header + "图片")
	}
	if message.Document != nil {
		log.Println(header + "文件: " + message.Document.FileName)
	}
	log.Println(header + "消息,无法预览")
}
