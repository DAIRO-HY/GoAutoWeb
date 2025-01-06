package ReadPathUtil

import (
	"GoAutoWeb/ReadInterceptorUtil"
	"strings"
)

// Controller路由信息
type PathBean struct {

	//包所在路径
	PackagePath string

	//请求方案
	Method string

	//路由路径
	Path string

	//函数名
	FuncName string

	//函数名
	ReturnType string

	//该路由的参数
	Parameters []ParamBean

	//使用template模板（html专用）
	Templates []string
}

// MakeHandleSource 生成Handle部分的代码
func (mine *PathBean) MakeHandleSource() string {
	return "\thttp.HandleFunc(\"" + mine.Path + "\", func(writer http.ResponseWriter, request *http.Request) {\n" +
		ReadInterceptorUtil.MappingBefore(mine.Path) + //执行前拦截器
		mine.getControllerParamSource() + // 获取Controller参数部分的代码
		mine.getCallMethodSource() + // 生成调用函数部分的代码
		ReadInterceptorUtil.MappingAfter(mine.Path) + //执行前拦截器
		mine.getEndSource() +
		"\t})\n"
}

// 获取导入昵称
func (mine *PathBean) GetNickImport() string {
	if len(mine.PackagePath) == 0 {
		return ""
	}
	nick := strings.ReplaceAll(mine.PackagePath, "/", "")
	nick = strings.ReplaceAll(nick, "_", "")
	return nick
}

// 获取Controller参数部分的代码
func (mine *PathBean) getControllerParamSource() string {
	source := ""
	for _, parameter := range mine.Parameters {
		source += parameter.MakeGetParameterSource()
	}
	if len(source) > 0 { //生成URL参数和Body参数变量代码
		queryAndPostFormVarSource := ""
		queryAndPostFormVarSource += "\t\tquery := request.URL.Query()\n"
		queryAndPostFormVarSource += "\t\t//解析post表单\n"
		queryAndPostFormVarSource += "\t\trequest.ParseForm()\n"
		queryAndPostFormVarSource += "\t\tpostForm := request.PostForm\n"

		source = queryAndPostFormVarSource + source
	}
	return source
}

// 生成调用函数部分的代码
func (mine *PathBean) getCallMethodSource() string {

	//函数参数部分的代码
	methodParamSource := ""
	for _, parameter := range mine.Parameters { //传递参数
		if parameter.VarType == "http.ResponseWriter" {
			methodParamSource += "writer, "
			continue
		}
		if parameter.VarType == "*http.Request" {
			methodParamSource += "request, "
			continue
		}
		methodParamSource += parameter.Name + ", "
	}
	if len(methodParamSource) > 0 {
		methodParamSource = methodParamSource[:len(methodParamSource)-2]
	}

	//调用Controller代码
	callMethodSource := mine.GetNickImport() + "." + mine.FuncName + "(" + methodParamSource + ")"

	source := "\t\tvar body any = nil\n"
	if len(mine.ReturnType) > 0 { //如果有返回值
		source += "\t\tbody = " + callMethodSource
	} else {
		source += "\t\t" + callMethodSource
	}
	return source + "\n"
}

// 获取结尾部分调用代码
func (mine *PathBean) getEndSource() string {
	if len(mine.Templates) > 0 { //这是一个html模板路由
		source := ""
		for _, template := range mine.Templates {
			if strings.HasSuffix(template, ".html") { //这是一个html模板
				source = "\t\ttemplates := append([]string{\"resources/templates/" + template + "\"}, COMMON_TEMPLATES...)\n"
				continue
			}
			source += "\t\ttemplates = append(templates, " + template + "...)\n"
		}
		source += "\t\twriteToTemplate(writer, templates, body)\n"
		return source
	} else {
		return "\t\twriteToResponse(writer, body)\n"
	}
}
