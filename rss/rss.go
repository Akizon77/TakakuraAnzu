package rss

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Akizon77/TakakuraAnzu/config"
	"github.com/Akizon77/TakakuraAnzu/log"
	"io"
	"net/http"
	"strings"
)

var api string
var token string

func init() {
	if config.Config.WebApiEndpoint[len(config.Config.WebApiEndpoint)-1:] != "/" {
		api = config.Config.WebApiEndpoint + "/"
	} else {
		api = config.Config.WebApiEndpoint
	}
	token = config.Config.WebToken
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 向数据库中添加一个表 格式为 rss_subs<ChatID>
func AddRssForChatID(chatID int64, title string, link string) error {
	var url = fmt.Sprint(api+"add/?token=", token, "&user=", chatID)
	body, _ := json.Marshal(struct {
		User  int64  `json:"user"`
		Title string `json:"title"`
		Link  string `json:"link"`
	}{
		User:  chatID,
		Title: title,
		Link:  link,
	})
	log.Debug("Post " + url)
	response, err := http.Post(url, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {

		return errors.New(fmt.Sprint(response.Body))
	}
	return nil
}

// RemoveRssForChatID 删除订阅链接
func RemoveRssForChatID(chatID int64, link string) error {
	var url = fmt.Sprint(api+"del/?token=", token, "&user=", chatID)
	body, _ := json.Marshal(struct {
		User int64  `json:"user"`
		Link string `json:"link"`
	}{
		User: chatID,
		Link: link,
	})
	log.Debug("Post " + url)
	response, err := http.Post(url, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {

		return errors.New(fmt.Sprint(response.Body))
	}
	return nil
}

type AllRssResponse struct {
	Code    int64   `json:"code"`
	Message string  `json:"message"`
	Data    RSSData `json:"data"`
}

type RSSData struct {
	User int64 `json:"user"`
	RSS  []RSS `json:"rss"`
}

type RSS struct {
	Title   string `json:"title"`
	SubLink string `json:"sub_link"`
}

func ListAllSubs(chatID int64) (string, error) {
	var url = fmt.Sprint(api+"all/?token=", token, "&user=", chatID)
	log.Debug("Get " + url)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	res, _ := io.ReadAll(response.Body)
	st := &AllRssResponse{}
	_ = json.Unmarshal(res, st)
	if response.StatusCode != 200 {
		return "", errors.New(st.Message)
	}
	var result string
	for _, r := range st.Data.RSS {
		result = result + fmt.Sprint("Title:", r.Title, "\n", "Link:`", r.SubLink, "`\n\n")
	}
	return result, nil
}

type AllUpdates struct {
	Code    int64      `json:"code"`
	Message string     `json:"message"`
	Data    UpdateData `json:"data"`
}

type UpdateData struct {
	User    int64     `json:"user"`
	Updates []Updates `json:"updates"`
}

type Updates struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func Update(chatid int64) (string, error) {
	var url = fmt.Sprint(api+"updates/?token=", token, "&user=", chatid)
	log.Debug("Get " + url)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	res, _ := io.ReadAll(response.Body)
	st := &AllUpdates{}
	_ = json.Unmarshal(res, st)
	if response.StatusCode != 200 {
		return "", errors.New(st.Message)
	}
	var result string
	if len(st.Data.Updates) > 30 {
		return "更新数量超过30，请手动查看", nil
	}
	for i, r := range st.Data.Updates {
		result = result + fmt.Sprint(i+1, ". <a href=\"", r.Link, "\">", r.Title, "</a>\n")
	}
	return result, nil
}
