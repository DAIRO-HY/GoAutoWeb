package main

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/MakeSourceUtil"
	"GoAutoWeb/ReadInterceptorUtil"
	"GoAutoWeb/ReadPathUtil"
	"fmt"
	"os"
	"strings"
)

const VERSION = "1.0.0"

func main() {
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

	//生成代码
	MakeSourceUtil.Make()
}
