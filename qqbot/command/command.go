package command

import (
	logs "github.com/Akizon77/TakakuraAnzu/log"
	"github.com/Akizon77/TakakuraAnzu/minecraft"
	"github.com/Akizon77/TakakuraAnzu/qqbot"
	"github.com/tencent-connect/botgo/dto"
	"strings"
)

func RunCommand(msg *dto.Message) error {
	send := ""
	if strings.Contains(msg.Content, "/mcs") {
		send = mcsCommand(msg)
	}
	//
	if msg.ChannelID != "" {
		message := qqbot.NewMessage(send, msg.ID)
		err := qqbot.SendToChannel(msg.ChannelID, message)
		if err != nil {
			return err
		}
	}
	return nil
}
func mcsCommand(msg *dto.Message) string {
	split := strings.Split(msg.Content, " ")
	if len(split) == 2 {
		logs.Info("(2)查询 cc.akz.moe")
		return minecraft.GetPrettiedStringForQQ("cc.akz.moe")
	}
	if len(split) == 3 {
		logs.Info("(3)查询" + split[2])
		if split[2] == "" {
			logs.Info("(3->2)查询 cc.akz.moe")
			return minecraft.GetPrettiedStringForQQ("cc.akz.moe")
		}
		return minecraft.GetPrettiedStringForQQ(split[2])
	} else {
		return "Usage:/mcs host"
	}
}
