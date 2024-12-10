package Application

import (
	"GoAutoWeb/FileUtil"
	"strings"
)

// 项目根目录
var RootProject string

// go代码文件列表
var GoFileList []string

// 项目的模块名
var ModuleName string

func Init(folder string) {
	RootProject = folder
	readModuleName()
	GoFileList = FileUtil.GetGoFile(RootProject)
}

// 读取项目的模块名
func readModuleName() {
	gomod := FileUtil.ReadText(RootProject + "/go.mod")
	gomod = strings.TrimSpace(gomod)
	gomod = strings.ReplaceAll(gomod, "\r\n", "\n")
	gomod = strings.ReplaceAll(gomod, "\n", " ")
	ModuleName = strings.Split(gomod, " ")[1]
}
