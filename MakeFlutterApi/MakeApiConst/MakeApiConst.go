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
		url := it.Path + it.VariablePath
		key := strings.ReplaceAll(url, "/", "_")
		key = strings.ReplaceAll(key, "{", "_")
		key = strings.ReplaceAll(key, "}", "_")
		key = strings.ReplaceAll(key, "__", "_")
		key = strings.ReplaceAll(key, "__", "_")
		key = strings.ReplaceAll(key, ".", "_")
		key = strings.ToUpper(key)
		key = key[1:]
		source += "  static const " + key + " = \"" + url + "\";\n"
	}
	source = "class Api{\n" + source + "}"
	save(source)
}

// 保存文件
func save(source string) {
	os.WriteFile(Application.Args.TargetDir+"/lib/api/API.dart", []byte(source), 0644)
}
