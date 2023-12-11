package network

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
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

	return "The source code contains the private key, you need to manually complete this code", errors.New("The source code contains the private key, you need to manually complete this code")
}
