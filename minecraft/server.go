package minecraft

import (
	"encoding/json"
	"fmt"
	logs "github.com/Akizon77/TakakuraAnzu/log"
	"io"
	"net/http"
	"strings"
)

func GetPrettiedString(host string) string {
	req, err := http.NewRequest(http.MethodGet, "https://api.mcstatus.io/v2/status/java/"+host, nil)
	if err != nil {
		return fmt.Sprintf(":C 无法创建请求\n%s", err.Error())
	}
	c := &http.Client{}
	response, err := c.Do(req)
	if err != nil {
		return fmt.Sprintf(":C 无法请求\n%s", err.Error())
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Sprintf(":C 无法读取返回数据\n%s", err.Error())
	}
	//s := string(bytes)
	//logs.Debug(s)
	server := &Server{}
	//总会有人打错IP
	s := string(bytes)
	if strings.Contains(s, "Invalid address value") {
		return fmt.Sprintf(":C 带上一个正确的IP地址好不好嘛\n")
	}
	err = json.Unmarshal(bytes, server)
	if err != nil {
		return fmt.Sprintf(":C 无法读取返回数据\n%s", err.Error())
	}
	//判断在线状态
	var online string
	if server.Online {
		online = "在线"
	} else {
		online = "离线"
	}
	basic := fmt.Sprintf("服务器：`%s:%d` \n版本：%s \n状态：%s\n人数：%d/%d\nMotd：\n%s", server.Host, server.Port, server.Version.NameClean, online, server.Players.Online, server.Players.Max, server.MOTD.Clean)
	if server.EULABlocked {
		basic = basic + "\n> 此服务器已被 Mojang/Microsoft 封禁"
	}
	return basic
}
func GetPrettiedStringForQQ(host string) string {
	req, err := http.NewRequest(http.MethodGet, "https://api.mcstatus.io/v2/status/java/"+host, nil)
	if err != nil {
		return fmt.Sprintf(":C 无法创建请求\n%s", err.Error())
	}
	c := &http.Client{}
	response, err := c.Do(req)
	if err != nil {
		return fmt.Sprintf(":C 无法请求\n%s", err.Error())
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		logs.Warn(string(bytes))
		return fmt.Sprintf(":C 无法读取返回数据\n%s", err.Error())
	}
	//s := string(bytes)
	//logs.Debug(s)
	server := &Server{}
	//总会有人打错IP
	s := string(bytes)
	if strings.Contains(s, "Invalid address value") {
		return fmt.Sprintf(":C 带上一个正确的IP地址好不好嘛\n")
	}
	err = json.Unmarshal(bytes, server)
	if err != nil {
		logs.Warn(string(bytes))
		return fmt.Sprintf(":C 无法读取返回数据\n%s", err.Error())
	}
	//判断在线状态
	var online string
	if server.Online {
		online = "在线"
	} else {
		online = "离线"
	}
	basic := fmt.Sprintf("版本：%s \n状态：%s\n在线：%d/%d", server.Version.NameClean, online, server.Players.Online, server.Players.Max)
	if server.EULABlocked {
		basic = basic + "\n> 此服务器已被 Mojang/Microsoft 封禁"
	}
	return basic
}
