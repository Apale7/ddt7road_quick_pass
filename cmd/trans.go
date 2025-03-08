package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// 定义结构体来保存 GAMEFORM 部分的数据
type GameForm struct {
	Name      string
	Platform  string
	URL       string
	Username  string
	Password  string
	PageNote  string
	GameSpeed float64
	AutoLogin int
	LoginNum  int
}

// 定义结构体来保存整个文件的数据
type DatFile struct {
	Encry   string
	Version string
	Forms   []GameForm
}

// 文件格式
/*
encry=
version=0

[GAMEFORM]
name=弹弹堂官网经典版3区
platform=wan.com
url={http://www.wan.com/game/play/id/8665.html}
username=201675F43470FFE663F43EB7AC65BAA6
pwd=cWlkYW85OTgyNDQzNTM=
pagenote=
gamespeed=1.00
autologin=1
loginnum=6

[GAMEFORM]
name=弹弹堂官网经典版3区
platform=wan.com
url={http://www.wan.com/game/play/id/8665.html}
username=Apale7
pwd=bXVtdW11OTg3
pagenote=
gamespeed=1.00
autologin=1
loginnum=2
*/

func main() {
	// 打开文件
	file, err := os.Open("tg_account.dat")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	bs, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("无法读取文件:", err)
		return
	}

	// 换行符是crlf
	blocks := strings.Split(string(bs), "\n\n")
	fmt.Println(len(blocks))
	// newBlocks := []string{blocks[0]}

	// 正则表达式获取username=和pwd=后面的内容
	usernameRegex := regexp.MustCompile(`username=([^\s]+)`)
	pwdRegex := regexp.MustCompile(`pwd=([^\s]+)`)
	// 正则表达式修改url, 在url后面加入username和pwd两个query参数
	// reURL := regexp.MustCompile(`url=\{(\S+)\}`)
	for _, block := range blocks[1:] {
		usernameMatch := usernameRegex.FindStringSubmatch(block)
		pwdMatch := pwdRegex.FindStringSubmatch(block)
		fmt.Println(usernameMatch, pwdMatch)
		// if len(res) > 2 {
		// 	block = reURL.ReplaceAllString(block, `url={$1?username=`+res[1]+`&password=`+res[2]+`}`)
		// }

		// newBlocks = append(newBlocks, block)
	}

	// result := strings.Join(newBlocks, "\n\n")
	// os.WriteFile("tg_account_new.dat", []byte(result), 0o644)
}
