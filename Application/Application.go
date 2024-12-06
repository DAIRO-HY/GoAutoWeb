package Application

import (
	"GoAutoController/FileUtil"
	"strings"
)

// 项目根目录
var RootProject = "C:\\develop\\project\\idea\\DairoNPS"

// go代码文件列表
var GoFileList = FileUtil.GetGoFile(RootProject)

// 项目的模块名
var ModuleName = readModuleName()

// 读取项目的模块名
func readModuleName() string {
	gomod := FileUtil.ReadText(RootProject + "/go.mod")
	gomod = strings.TrimSpace(gomod)
	gomod = strings.ReplaceAll(gomod, "\r\n", "\n")
	gomod = strings.ReplaceAll(gomod, "\n", " ")
	moduleName := strings.Split(gomod, " ")[1]
	return moduleName
}
