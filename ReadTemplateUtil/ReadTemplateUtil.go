package ReadTemplateUtil

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/FileUtil"
	"strings"
)

// 模板名称对应的路径
var TemplateNameToPath = map[string]string{}

func Make() {
	for _, path := range Application.HtmlFileList {
		for _, it := range readName(path) {
			absPath := strings.ReplaceAll(path, "\\", "/")
			absPath = absPath[len(Application.RootProject)+1:]
			TemplateNameToPath[it] = absPath
		}
	}
}

// 获取模板名
func readName(path string) []string {
	content := FileUtil.ReadText(path)

	//先统一换行符
	content = strings.ReplaceAll(content, "\r", "\r\n")
	content = strings.ReplaceAll(content, "\r\n", "\n")

	var names []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "{{define") {
			name := line
			name = line[strings.Index(name, "\"")+1:]
			name = name[:strings.Index(name, "\"")]
			names = append(names, name)
		}
	}
	return names
}

// 从html文件中获取使用到模板名
func ReadUseTemplatesByHtml(html string) string {
	if strings.HasPrefix(html, "/") {
		html = html[1:] //去掉前面的斜杠
	}

	//获取页面html绝对路径
	path := Application.RootProject + "/resources/templates/" + html
	content := FileUtil.ReadText(path)

	//先统一换行符
	content = strings.ReplaceAll(content, "\r", "\r\n")
	content = strings.ReplaceAll(content, "\r\n", "\n")

	//使用到的模板名称
	var useTemplateNames []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "{{template") {
			name := line
			name = line[strings.Index(name, "\"")+1:]
			name = name[:strings.Index(name, "\"")]
			useTemplateNames = append(useTemplateNames, name)
		}
	}
	if len(useTemplateNames) == 0 {
		return ""
	}
	templatePathMap := map[string]bool{}
	for _, it := range useTemplateNames {
		templatePath := TemplateNameToPath[it]
		templatePathMap[templatePath] = true
	}

	//使用的模板相对路径列表
	var templatePaths []string
	for key := range templatePathMap {
		templatePaths = append(templatePaths, key)
	}

	// 使用的模板路径代码
	source := "\"resources/templates/" + html + "\""
	for _, it := range templatePaths {
		source += ", \"" + it + "\""
	}
	return source
}
