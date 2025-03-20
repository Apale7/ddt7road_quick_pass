package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/sirupsen/logrus"
)

func main() {
	Init()
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8888", "proxy listen address")
	flag.Parse()

	// err := SetCA(cert, key)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	proxy := goproxy.NewProxyHttpServer()
	proxy.Tr.TLSClientConfig.MinVersion = tls.VersionTLS10
	proxy.Verbose = *verbose
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	ctx := context.TODO()
	proxy.OnRequest().DoFunc(func(req *http.Request, c *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, nil
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, c *goproxy.ProxyCtx) *http.Response {
		u := c.Req.URL
		// fmt.Println(u.String())
		if strings.HasPrefix(u.String(), "http://www.wan.com/game/play/id/8665.html") || strings.HasPrefix(u.String(), "http://www.wan.com/game/play/id/4048.html") {

			// cookiesStr, _ := json.Marshal(resp.Cookies())
			// logrus.Infof("%s", string(cookiesStr))

			for _, cookie := range resp.Cookies() {
				if cookie.Name == "name" && cookie.Value != "deleted" {
					// logrus.Info(cookie.Name + ": " + cookie.Value)
					// logrus.Info("已登录, 无需登录")
					return resp
				}
			}

			// 从url解析出username和password
			username := u.Query().Get("username")
			password, ok := accountMap[username]
			if !ok {
				logrus.Println("未知账号")
			}
			// password := u.Query().Get("pwd")
			// password做base64反解码
			passwordBytes, err := base64.StdEncoding.DecodeString(password)
			if err != nil {
				logrus.Errorln("base64解码失败: ", err) // base64解码失败，可能是明文密码，继续往下跑
			} else {
				password = string(passwordBytes)
			}

			logrus.Infof("username: %s, password: %s", username, password)
			cookies, err := Login(ctx, username, password)
			if err != nil {
				logrus.Errorln("自动登录失败: ", err)
				return resp
			}
			logrus.Info("自动登录")
			// logrus.Infof("setcookie: %s", cookies)
			for _, cookie := range cookies {
				resp.Header.Add("Set-Cookie", cookie)
			}
			logrus.Info("setcookie")
		}
		return resp
	})

	logrus.Info("自动登录插件已启动")
	logrus.Fatal(http.ListenAndServe(*addr, proxy))
}
