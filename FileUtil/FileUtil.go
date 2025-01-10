package FileUtil

import (
	"fmt"
	"os"
)

// 读取文本文件
func ReadText(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// 写入文件
func WriteText(path string, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
