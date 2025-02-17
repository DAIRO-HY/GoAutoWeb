package main

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/MakeSourceUtil"
	"GoAutoWeb/ReadFormUtil"
	"GoAutoWeb/ReadInterceptorUtil"
	"GoAutoWeb/ReadPathUtil"
	"GoAutoWeb/ReadTemplateUtil"
	"fmt"
	"os"
	"strings"
	"time"
)

const VERSION = "1.0.0"

// 是否能匹配当前路由参数
func isPathVariable(path string, splitList []string) bool {
	if !strings.HasPrefix(path, splitList[0]) { //判断前缀是否一致
		return false
	}
	for _, it := range splitList { //挨个匹配路由
		index := strings.Index(path, it)
		if index == -1 {
			return false
		}
		path = path[index+len(it):]
	}
	return path == ""
}

func main() {
	start := time.Now()
	if len(os.Args) == 1 {
		currentFolder := os.Args[0]
		currentFolder = strings.ReplaceAll(currentFolder, "\\", "/")
		currentFolder = currentFolder[:strings.LastIndex(currentFolder, "/")]
		Application.Init(currentFolder)
	} else {
		Application.Init(os.Args[1])
	}
	fmt.Println(Application.RootProject)

	//初始化读取路由列表
	ReadPathUtil.Make()

	//初始化读取拦截器列表
	ReadInterceptorUtil.Make()

	//生成Form表单列表
	ReadFormUtil.Make()

	//生成模板数据
	ReadTemplateUtil.Make()

	//生成代码
	MakeSourceUtil.Make()
	fmt.Printf("本次耗时：%s", time.Since(start))
}
