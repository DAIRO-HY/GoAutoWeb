package MakeApiConst

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/FileUtil"
	"GoAutoWeb/ReadPathUtil"
	"strings"
)

// 生成API常量文件
func Make() {
	body := ""
	for _, it := range ReadPathUtil.PathList {
		if !strings.HasSuffix(it.FileName, Application.Args.ApiSuffix+".go") {
			continue
		}
		url := it.Path + it.VariablePath
		body += makeComment(it) + "\n"
		body += "  const val " + urlToConst(it) + " = \"" + url + "\"\n"
	}
	classSource := "package " + Application.Args.TargetPackage + "\n"
	classSource += "object ApiConst{\n" + body + "}"
	save(classSource)
}

// 生成注释部分的代码
func makeComment(pb ReadPathUtil.PathBean) string {
	comment := pb.Comment
	if comment == "" {
		return ""
	}
	cms := strings.Split(comment, "\n")
	return "\n  //" + strings.Join(cms, "\n  //")
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
	path := Application.Args.TargetDir + "/ApiConst.kt"
	fileContent := FileUtil.ReadText(path)
	if fileContent == source {
		return
	}
	FileUtil.WriteText(path, source)
}
