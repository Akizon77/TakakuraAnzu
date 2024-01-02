package message

import (
	"fmt"
	logs "github.com/Akizon77/TakakuraAnzu/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

func InitMessageTransfer(tgbot *tgbotapi.BotAPI) {
	bot = tgbot
}
func Transfer(sender string, form string, content string, to int64) {
	msg := tgbotapi.NewMessage(to, fmt.Sprint(sender, "(", form, ")", ":\n", content))
	logs.SendMsg(msg, bot)
}
