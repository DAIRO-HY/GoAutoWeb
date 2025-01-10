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
