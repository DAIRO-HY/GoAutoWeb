package Global

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/FileUtil"
	"GoAutoWeb/MakeGoFileInfo"
	"GoAutoWeb/MakeGoFileInfo/Bean"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 项目根目录
var RootProject string

// go代码文件列表
var GoFileList []string

// go代码信息列表
var GoBeanList []Bean.GoBean

// html模板文件列表
var HtmlFileList []string

// 项目的模块名
var ModuleName string

func Init() {
	RootProject = Application.Args.SourceDir
	readModuleName()
	makeFileList()
	for _, it := range GoFileList {
		GoBeanList = append(GoBeanList, MakeGoFileInfo.ReadGoInfo(it))
	}
}

// 读取项目的模块名
func readModuleName() {
	gomod := FileUtil.ReadText(RootProject + "/go.mod")
	gomod = strings.TrimSpace(gomod)
	gomod = strings.ReplaceAll(gomod, "\r\n", "\n")
	gomod = strings.ReplaceAll(gomod, "\n", " ")
	ModuleName = strings.Split(gomod, " ")[1]
}

// 获取go文件列表
func makeFileList() {

	// 遍历文件夹
	err := filepath.WalkDir(RootProject, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasSuffix(path, ".idea") || strings.HasSuffix(path, ".git") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			GoFileList = append(GoFileList, path)
		}
		if strings.HasSuffix(path, ".html") {
			HtmlFileList = append(HtmlFileList, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", RootProject, err)
	}
}
