package qqbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Akizon77/TakakuraAnzu/config"
	logs "github.com/Akizon77/TakakuraAnzu/log"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendToChannel(channelID string, message *Message) error {
	url := "https://api.sgroup.qq.com/channels/" + channelID + "/messages"

	payload, _ := json.Marshal(message)

	req, _ := http.NewRequest("POST", url, strings.NewReader(string(payload)))

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bot "+fmt.Sprint(config.Config.QQ_App_id)+"."+config.Config.QQ_Token)

	res, err := http.DefaultClient.Do(req)

	defer res.Body.Close()
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(res.Body)
	serverReply := &SendMsgReply{}
	json.Unmarshal(body, &serverReply)
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprint(serverReply.ErrorCode, serverReply.ErrorMessage))
	}
	logs.Info("发送:" + message.Content)
	return nil
}
