package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func writeTgProxyConf() error {
	// 文件路径
	fPath := filepath.Join(tgConfDir, "userproxy.txt")

	// 打开文件
	file, err := os.Open(fPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// 将 UTF-16LE 转换为 UTF-8
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	reader := transform.NewReader(file, decoder)

	// 读取文件内容
	content, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading file content:", err)
		return err
	}

	// 检查是否已包含目标配置
	if strings.Contains(string(content), "免验证码") {
		logrus.Println("已配置代理")
		return nil
	}

	// 定义要插入的新配置
	const conf = "\r\n[USERPROXY]\r\nname=免验证码\r\nip=127.0.0.1\r\nport=8888\r\nproxyType=http\r\npasswordCheck=0\r\nusername=\r\npassword=\r\nisActive=1\r\n"

	// 将文件内容按行分割
	lines := strings.Split(string(content), "\r\n")

	// 在第二行插入新配置
	var newContent string
	if len(lines) > 1 {
		newContent = strings.Join(lines[:1], "\r\n") + "\r\n" + conf + strings.Join(lines[1:], "\r\n")
	} else {
		// 如果文件只有一行或为空，直接追加新配置
		newContent = string(content) + "\r\n" + conf
	}

	// 将新内容写回文件，使用 UTF-16LE 编码并添加 BOM
	err = writeFileWithUTF16LE(fPath, newContent)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	logrus.Println("代理配置已成功写入文件")
	return nil
}

// 辅助函数：以带 BOM 的 UTF-16LE 编码写入文件
func writeFileWithUTF16LE(filePath, content string) error {
	// 打开文件用于写入
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	// 创建缓冲区
	buf := new(bytes.Buffer)

	// 添加 BOM（0xFF 0xFE）
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFE)

	// 使用 UTF-16LE 编码器转换内容
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	writer := transform.NewWriter(buf, encoder)

	_, err = writer.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to encode content to UTF-16LE: %w", err)
	}

	// 确保所有数据被写入缓冲区
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// 写入文件
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
