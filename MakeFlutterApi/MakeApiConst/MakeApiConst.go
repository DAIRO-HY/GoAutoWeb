package MakeApiConst

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/ReadPathUtil"
	"os"
	"strings"
)

// 生成API常量文件
func Make() {
	source := ""
	for _, it := range ReadPathUtil.PathList {
		if !strings.HasSuffix(it.FileName, Application.Args.ApiSuffix+".go") {
			continue
		}
		url := it.Path + it.VariablePath
		source += "  static const " + urlToConst(it) + " = \"" + url + "\";\n"
	}
	source = "class Api{\n" + source + "}"
	save(source)
}

// 将路由转成常量名
func urlToConst(pb ReadPathUtil.PathBean) string {
	url := pb.Path + pb.VariablePath
	key := strings.ReplaceAll(url, "/", "_")
	key = strings.ReplaceAll(key, "{", "_")
	key = strings.ReplaceAll(key, "}", "_")
	key = strings.ReplaceAll(key, "__", "_")
	key = strings.ReplaceAll(key, "__", "_")
	key = strings.ReplaceAll(key, ".", "_")
	key = strings.ToUpper(key)
	return key[1:]
}

// 保存文件
func save(source string) {
	os.WriteFile(Application.Args.TargetDir+"/lib/api/API.dart", []byte(source), 0644)
}
