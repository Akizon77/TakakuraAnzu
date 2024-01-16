package message

import (
	"fmt"
	"github.com/Akizon77/TakakuraAnzu/config"
	"github.com/Akizon77/TakakuraAnzu/data/whitelist"
	messageLogger "github.com/Akizon77/TakakuraAnzu/log"
	"github.com/Akizon77/TakakuraAnzu/minecraft"
	"github.com/Akizon77/TakakuraAnzu/network"
	"github.com/Akizon77/TakakuraAnzu/rss"
	"github.com/Akizon77/TakakuraAnzu/status"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

const (
	noPermission = "呀呀呀，权限好像不够！\n去问问 @AkizonChan "
)

func RunCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if !message.IsCommand() {
		return
	}
	if message.Chat.IsGroup() || message.Chat.IsSuperGroup() {
		if !strings.Contains(message.Text, bot.Self.UserName) {
			return
		}
	}
	switch message.Command() {
	case "start":
		messageLogger.Debug("startCommand")
		go startCommand(message, bot)
	case "refresh":
		messageLogger.Debug("refreshCommand")
		go refreshCommand(message, bot)
	case "chatid":
		messageLogger.Debug("chatIdCommand")
		go chatIdCommand(message, bot)
	case "status":
		messageLogger.Debug("statusCommand")
		go statusCommand(message, bot)
	case "ip":
		messageLogger.Debug("ipCommand")
		go ipCommand(message, bot)
	case "ddns":
		messageLogger.Debug("ddnsCommand")
		go ddnsCommand(message, bot)
	case "add":
		messageLogger.Debug("addCommand")
		go addCommand(message, bot)
	case "remove":
		messageLogger.Debug("removeCommand")
		go removeCommand(message, bot)
	//case "sql":
	//	messageLogger.Debug("sqlCommand")
	//	go sqlCommand(message, bot)
	case "list":
		messageLogger.Debug("listCommand")
		go listCommand(message, bot)
	case "whitelist":
		messageLogger.Debug("whitelistCommand")
		go whitelistCommand(message, bot)
	case "mcs":
		messageLogger.Debug("mcsCommand")
		go mcsCommand(message, bot)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("杏铃不认识 %s 哦", message.Command()))
		msg.ReplyToMessageID = message.MessageID
		msg.ParseMode = tgbotapi.ModeMarkdown
		messageLogger.SendMsg(msg, bot)
	}
}
func whitelistCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Usage:\nadd - 添加白名单\ndel - 删除白名单\ndrop - 清空白名单"))
	msg.ReplyToMessageID = message.MessageID
	msg.DisableWebPagePreview = true
	msg.ParseMode = tgbotapi.ModeMarkdown

	arg := message.CommandArguments()
	args := strings.Split(arg, " ")
	if len(args) <= 1 {
		messageLogger.SendMsg(msg, bot)
		return
	}
	message.Text = "/whitelist " + args[1]
	switch args[0] {
	case "add":
		addWhitelistCommand(message, bot)
	case "del":
		delWhitelistCommand(message, bot)
	case "drop":
		dropWhitelistCommand(message, bot)
	default:
		messageLogger.SendMsg(msg, bot)
	}
}
func mcsCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Error"))
	msg.ReplyToMessageID = message.MessageID
	msg.DisableWebPagePreview = true
	msg.ParseMode = tgbotapi.ModeMarkdown
	server := message.CommandArguments()
	if server == "" {
		server = config.Config.MinecraftServer
	}
	info := minecraft.GetPrettiedString(server)
	msg.Text = info
	messageLogger.SendMsg(msg, bot)
}
func dropWhitelistCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Error"))
	msg.ReplyToMessageID = message.MessageID
	msg.DisableWebPagePreview = true
	msg.ParseMode = tgbotapi.ModeMarkdown

	if message.From.ID != 1977354088 {
		msg.Text = "可惜，只有 @AkizonChan 能drop白名单啦，有需要的话就快去找他吧"
		messageLogger.SendMsg(msg, bot)
		return
	}
	err := whitelist.Clear()
	if err != nil {
		msg.Text = fmt.Sprintf("搞不定，因为 %s", err.Error())
	} else {
		msg.Text = fmt.Sprintf("啊，白名单就这么Drop掉了")
	}
	messageLogger.SendMsg(msg, bot)
}

func delWhitelistCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Error"))
	msg.ReplyToMessageID = message.MessageID
	msg.DisableWebPagePreview = true
	msg.ParseMode = tgbotapi.ModeMarkdown

	if message.From.ID != 1977354088 {
		msg.Text = "可惜，只有 @AkizonChan 能移除白名单啦，有需要的话就快去找他吧"
		messageLogger.SendMsg(msg, bot)
		return
	}
	if message.CommandArguments() == "" {
		msg.Text = "你都不告诉我要把谁删除白名单"
		messageLogger.SendMsg(msg, bot)
		return
	}
	arg, err := strconv.ParseInt(message.CommandArguments(), 10, 64)
	if err != nil {
		msg.Text = "这好像不是一个ChatID？"
		messageLogger.SendMsg(msg, bot)
		return
	}
	err = whitelist.Remove(arg)
	if err != nil {
		msg.Text = fmt.Sprintf(err.Error())
		messageLogger.SendMsg(msg, bot)
		return
	}
	msg.Text = "已经把他删掉了"
	messageLogger.SendMsg(msg, bot)
}

func addWhitelistCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Error"))
	msg.ReplyToMessageID = message.MessageID
	msg.DisableWebPagePreview = true
	msg.ParseMode = tgbotapi.ModeMarkdown

	if message.From.ID != 1977354088 {
		msg.Text = "可惜，只有 @AkizonChan 能添加白名单啦，有需要的话就快去找他吧"
		messageLogger.SendMsg(msg, bot)
		return
	}
	if message.CommandArguments() == "" {
		msg.Text = "你都不告诉我要把谁添加到白名单"
		messageLogger.SendMsg(msg, bot)
		return
	}
	arg, err := strconv.ParseInt(message.CommandArguments(), 10, 64)
	if err != nil {
		msg.Text = "这好像不是一个ChatID？"
		messageLogger.SendMsg(msg, bot)
		return
	}
	if message.Chat.IsGroup() || message.Chat.IsSuperGroup() {
		err = whitelist.Add(arg, message.Chat.Title)
	} else {
		err = whitelist.Add(arg, message.From.UserName)
	}
	if err != nil {
		msg.Text = fmt.Sprintf("无法添加，因为 %s", err.Error())
		messageLogger.SendMsg(msg, bot)
		return
	}

	msg.Text = "成功添加到白名单"
	messageLogger.SendMsg(msg, bot)
	return

}

// 列出所有订阅链接
func listCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("好像列不出来"))
	msg.ReplyToMessageID = message.MessageID
	msg.DisableWebPagePreview = true
	msg.ParseMode = tgbotapi.ModeMarkdown

	subs, err := rss.ListAllSubs(message.Chat.ID)
	if err != nil {
		if strings.Contains(err.Error(), "no such") {
			msg.Text = "还没有订阅呢，要不要用 /add 添加一个？"
		} else {
			msg.Text = fmt.Sprintf("啊呀！出错了\n%s", err.Error())
		}

	} else if subs == "" {
		msg.Text = "还没有订阅呢，要不要用 /add 添加一个？"
	} else {
		msg.Text = subs
	}

	messageLogger.SendMsg(msg, bot)
}

// 远程执行SQL指令

// 添加数据库的RSS订阅
func addCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("已经加上了哦"))
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = tgbotapi.ModeMarkdown

	arg := message.CommandArguments()
	if arg == "" {
		msg.Text = fmt.Sprint("或许可以订阅空气呢\nUsage:`/add title link`")
		messageLogger.SendMsg(msg, bot)
		return
	}
	//TODO 处理title
	args := strings.Split(arg, " ")
	messageLogger.Debug(fmt.Sprint(args))
	if len(args) < 2 {
		msg.Text = fmt.Sprint("好像识别不出来捏\nUsage:`/add title link`")
		messageLogger.SendMsg(msg, bot)
		return
	}
	err := rss.AddRssForChatID(message.Chat.ID, args[0], args[1])
	if err != nil {
		msg.Text = fmt.Sprintf("无法添加订阅\n" + err.Error())
	}

	messageLogger.SendMsg(msg, bot)
}

// 移除数据库中的RSS订阅
func removeCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("已经帮你删掉啦"))
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = tgbotapi.ModeMarkdown

	arg := message.CommandArguments()
	if arg == "" {
		msg.Text = fmt.Sprint("杏铃不知道你要删掉什么哦\nUsage:`/remove link`")
		messageLogger.SendMsg(msg, bot)
		return
	}
	log.Println("指令参数：" + arg)

	err := rss.RemoveRssForChatID(message.Chat.ID, arg)
	if err != nil {
		msg.Text = fmt.Sprintf("删不掉呜呜呜\n" + err.Error())
	}

	messageLogger.SendMsg(msg, bot)
}

// 手动刷新Rss订阅
func refreshCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("已经尝试刷新了哦"))
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML
	update, err := rss.Update(message.Chat.ID)
	if err != nil {
		msg.Text = fmt.Sprint(err)
	}
	msg.Text = update
	messageLogger.SendMsg(msg, bot)
}

// 欢迎界面，同时将chatID记录到控制台
func startCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "你好，这里是鹰仓杏铃 Takakura Anzu ，是阿卡西的小助理。请多多关照啦")
	msg.ReplyToMessageID = message.MessageID
	messageLogger.SendMsg(msg, bot)
	if message.Chat.IsGroup() || message.Chat.IsSuperGroup() {
		log.Printf("群组 " + message.Chat.Title + " 的ChatID是：" + fmt.Sprintf("%d", message.Chat.ID))
	} else {
		log.Printf("与 " + message.From.FirstName + message.From.LastName + " 的ChatID是：" + fmt.Sprintf("%d", message.Chat.ID))
	}

}

// 发送当前聊天ID
func chatIdCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("`%d`", message.Chat.ID))
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = tgbotapi.ModeMarkdown
	messageLogger.SendMsg(msg, bot)
}

// 发送IP地址
func ipCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "未知错误，位于function ipCommand")
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyToMessageID = message.MessageID

	if !whitelist.IsWhitelist(message.Chat.ID) {
		msg.Text = noPermission
		messageLogger.SendMsg(msg, bot)
		return
	}
	ip, err := network.GetIPv4()
	if err != nil {
		msg.Text = "呜呜呜， 拿不到IP地址！\n" + err.Error()
	} else {
		msg.Text = fmt.Sprintf("`%s`", ip)
	}

	messageLogger.SendMsg(msg, bot)
}

// 更新DDNS
func ddnsCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "未知错误，位于function ddnsCommand")

	msg.ReplyToMessageID = message.MessageID
	if !whitelist.IsWhitelist(message.Chat.ID) {
		msg.Text = noPermission
		messageLogger.SendMsg(msg, bot)
		return
	}

	ip, err := network.GetIPv4()
	if err != nil {
		msg.Text = fmt.Sprintf("呀呀呀，发生错误\n%s", fmt.Sprintf(err.Error()))
		messageLogger.SendMsg(msg, bot)
		return
	}
	res, err := network.UpdateDDNS(ip)
	if err != nil {
		msg.Text = fmt.Sprintf("呀呀呀，发生错误\n%s", fmt.Sprintf(err.Error()))
		messageLogger.SendMsg(msg, bot)
		return
	}
	log.Println(res)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.Text = fmt.Sprintf("欸嘿嘿，已成功更新`%s`到`cc.akz.moe`", ip)
	messageLogger.SendMsg(msg, bot)
}

// 发送当前系统信息
func statusCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, status.GetStatusFormattedString())
	msg.ReplyToMessageID = message.MessageID
	messageLogger.SendMsg(msg, bot)
}
