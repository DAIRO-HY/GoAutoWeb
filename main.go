package main

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/Global"
	"GoAutoWeb/MakeFlutterApi"
	"GoAutoWeb/MakeSourceUtil"
	"GoAutoWeb/MakeSwiftApi"
	"GoAutoWeb/ReadFormUtil"
	"GoAutoWeb/ReadInterceptorUtil"
	"GoAutoWeb/ReadPathUtil"
	"GoAutoWeb/ReadTemplateUtil"
	"fmt"
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

	//初始化程序参数
	Application.Init()
	Global.Init()
	fmt.Println(Global.RootProject)

	//初始化读取路由列表
	ReadPathUtil.Make()

	//初始化读取拦截器列表
	ReadInterceptorUtil.Make()

	//生成Form表单列表
	ReadFormUtil.Make()

	//生成模板数据
	ReadTemplateUtil.Make()
	if Application.Args.TargetType == "web" { //生成web的controller代码
		MakeSourceUtil.Make()
	} else if Application.Args.TargetType == "flutter-api" {
		MakeFlutterApi.Make()
	} else if Application.Args.TargetType == "swift-api" {
		MakeSwiftApi.Make()
	} else {

	}
	fmt.Printf("本次耗时：%s", time.Since(start))
}
