package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type GameForm struct {
	Name      string
	Platform  string
	URL       string
	Username  string
	Pwd       string
	PageNote  string
	GameSpeed float64
	AutoLogin int
	LoginNum  int
}

var tgConfDir string

func InitAccounts() {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		fmt.Println("无法找到 LOCALAPPDATA 环境变量")
		return
	}

	tgConfDir = filepath.Join(localAppData, "TGGAME")
	// logrus.Println(tgConfDir)
}

func readConf() []GameForm {
	file, err := os.Open(filepath.Join(tgConfDir, `accouts2.dat`))
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// 将 UTF-16LE 转换为 UTF-8
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	reader := transform.NewReader(file, decoder)

	// 读取文件内容
	content, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading file content:", err)
		return nil
	}

	// 动态重命名重复的 [GAMEFORM] 节
	modifiedContent := renameDuplicateSections(string(content))

	// 使用 go-ini/ini 解析修改后的 INI 内容
	cfg, err := ini.Load([]byte(modifiedContent))
	if err != nil {
		fmt.Println("Error parsing INI file:", err)
		return nil
	}

	// 提取所有 [GAMEFORM_*] 节
	var gameForms []GameForm
	for _, section := range cfg.Sections() {
		if strings.HasPrefix(section.Name(), "GAMEFORM") {
			// 解析每个 GAMEFORM 节
			var form GameForm
			form.Name = section.Key("name").String()
			form.Platform = section.Key("platform").String()
			form.URL = section.Key("url").String()
			form.Username = section.Key("username").String()
			form.Pwd = section.Key("pwd").String()
			form.PageNote = section.Key("pagenote").String()
			form.GameSpeed, _ = section.Key("gamespeed").Float64()
			form.AutoLogin, _ = section.Key("autologin").Int()
			form.LoginNum, _ = section.Key("loginnum").Int()

			// 添加到数组
			gameForms = append(gameForms, form)
		}
	}

	// // 打印结果
	// for i, form := range gameForms {
	// 	fmt.Printf("GameForm %d:\n", i+1)
	// 	fmt.Printf("  Name: %s\n", form.Name)
	// 	fmt.Printf("  Platform: %s\n", form.Platform)
	// 	fmt.Printf("  URL: %s\n", form.URL)
	// 	fmt.Printf("  Username: %s\n", form.Username)
	// 	fmt.Printf("  Pwd: %s\n", form.Pwd)
	// 	fmt.Printf("  PageNote: %s\n", form.PageNote)
	// 	fmt.Printf("  GameSpeed: %.2f\n", form.GameSpeed)
	// 	fmt.Printf("  AutoLogin: %d\n", form.AutoLogin)
	// 	fmt.Printf("  LoginNum: %d\n", form.LoginNum)
	// 	fmt.Println()
	// }
	return gameForms
}

func copyFile(src, dst string) error {
	// 读取源文件内容
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// 写入目标文件
	err = os.WriteFile(dst, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}

	return nil
}

func saveConf(gameForms []GameForm) error {
	filename := filepath.Join(tgConfDir, `accouts2.dat`)
	// 保留副本
	copyFile(filename, filepath.Join(tgConfDir, `accouts2_old.dat`))
	// 打开文件用于写入
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// 构造 INI 文件内容，使用 CRLF (\r\n) 换行符
	var result strings.Builder
	result.WriteString("encry=\r\nversion=0\r\n\r\n")

	for _, form := range gameForms {
		result.WriteString(fmt.Sprintf(
			"[GAMEFORM]\r\n"+
				"name=%s\r\n"+
				"platform=%s\r\n"+
				"url=%s\r\n"+
				"username=%s\r\n"+
				"pwd=%s\r\n"+
				"pagenote=%s\r\n"+
				"gamespeed=%.2f\r\n"+
				"autologin=%d\r\n"+
				"loginnum=%d\r\n\r\n",
			form.Name, form.Platform, form.URL, form.Username, form.Pwd,
			form.PageNote, form.GameSpeed, form.AutoLogin, form.LoginNum,
		))
	}

	// 将字符串转换为 UTF-16LE 编码
	buf := new(bytes.Buffer)

	// 添加 BOM（0xFF 0xFE）
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFE)

	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	writer := transform.NewWriter(buf, encoder)

	_, err = writer.Write([]byte(result.String()))
	if err != nil {
		fmt.Println("Error encoding to UTF-16LE:", err)
		return err
	}

	// 确保所有数据都被写入缓冲区
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return err
	}

	// 写入文件
	_, err = file.Write(buf.Bytes())
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}

// 在url后面添加username
func editConf(gameForms []GameForm) {
	for i, form := range gameForms {
		u := gameForms[i].URL[1 : len(gameForms[i].URL)-1]
		if strings.Contains(u, "username=") { // 已经处理好了
			continue
		}

		u = fmt.Sprintf("%s?username=%s", u, form.Username)
		gameForms[i].URL = fmt.Sprintf("{%s}", u)
		// fmt.Printf("  URL: %s\n", gameForms[i].URL)
	}
}

// 动态重命名重复的 [GAMEFORM] 节
func renameDuplicateSections(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	counter := 1

	for _, line := range lines {
		if strings.TrimSpace(line) == "[GAMEFORM]" {
			// 重命名为 [GAMEFORM_1], [GAMEFORM_2], ...
			result = append(result, fmt.Sprintf("[GAMEFORM_%d]", counter))
			counter++
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
