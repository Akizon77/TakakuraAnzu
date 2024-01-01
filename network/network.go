package network

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func GetIPv4() (string, error) {
	resp, err := http.Get("https://ddns.oray.com/checkip")
	if err != nil {
		return "发送请求失败", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "无法从接口读取数据", err
	}

	str := string(body)
	re := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	ip := re.FindString(str)
	return ip, nil
}
func UpdateDDNS(ip string) (string, error) {
	res, err := PUT(ip)
	if err != nil {
		return "", err
	}
	if strings.Contains(res, "\"success\":false") {
		return "res", errors.New(fmt.Sprintf("请求失败\n%s", res))
	}
	return res, nil
}

func PUT(ip string) (string, error) {
	req, err := http.NewRequest("PUT",
		"https://api.cloudflare.com/client/v4/zones/8794b43213d9257d8dc9bb1f184073ea/dns_records/ceaaa9087b2f219457cb12a13e4ab538/",
		strings.NewReader(fmt.Sprintf("{\r\n  \"type\": \"A\",\r\n  \"name\": \"cc.akz.moe\",\r\n  \"content\": \"%s\",\r\n  \"ttl\": 60,\r\n  \"proxied\": false,\r\n  \"comment\": \"Updated by Takakura Anzu at %s\"\r\n}", ip, time.Now().Format(time.RFC3339))))
	if err != nil {
		return "", err
	}
	req.Header = map[string][]string{
		"X-Auth-Email": {"hytf233g@gmail.com"},
		"X-Auth-Key":   {"faa35a3bc4cf5613f626ef8d86109a86f199a"},
	}

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	content, err := io.ReadAll(res.Body)
	return string(content), nil
	//return "The source code contains the private key, you need to manually complete this code", errors.New("The source code contains the private key, you need to manually complete this code")
}
