package FileUtil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// 获取go文件列表
func GetGoFile(root string) []string {
	var goFileList []string

	// 遍历文件夹
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") {
			goFileList = append(goFileList, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err)
	}
	return goFileList
}
