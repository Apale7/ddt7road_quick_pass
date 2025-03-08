package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	ErrorPasswordChange = errors.New("密码修改失败")
	ErrorLoginFailed    = errors.New("登录失败")
)

// 登录并返回cookie
func Login(ctx context.Context, username, password string) ([]string, error) {
	url := fmt.Sprintf("https://www.wan.com/index.php/accounts/checklogin.html?cn=%s&pwd=%s", username, password)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	req.Header.Add("authority", "www.wan.com")
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Add("referer", "https://ddt.wan.com/")
	req.Header.Add("sec-ch-ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Microsoft Edge\";v=\"114\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "script")
	req.Header.Add("sec-fetch-mode", "no-cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.51")

	res, err := client.Do(req)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	if !strings.Contains(string(body), `\u767b\u9646\u6210\u529f`) {
		fmt.Println(string(body))
		return nil, ErrorLoginFailed
	}

	// 获取cookies
	cookies := res.Header.Values("Set-Cookie")

	return cookies, nil
}
