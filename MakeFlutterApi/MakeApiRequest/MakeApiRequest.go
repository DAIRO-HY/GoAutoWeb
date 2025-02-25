package MakeApiRequest

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/ReadPathUtil"
	"os"
	"strings"
)

// 生成API常量文件
func Make() {
	fileToSource := make(map[string]string)
	for _, it := range ReadPathUtil.PathList {
		className := strings.ReplaceAll(it.FileName, "Controller.go", "Api")
		body := fileToSource[className]
		body += "  static  " + it.FuncName + "(){}\n"
		fileToSource[className] = body
	}
	for key, value := range fileToSource {
		save(key, value)
	}
}

// 保存文件
func save(fileName string, source string) {
	os.WriteFile(Application.Args.TargetDir+"/lib/api/"+fileName+".New.dart", []byte(source), 0644)
}
