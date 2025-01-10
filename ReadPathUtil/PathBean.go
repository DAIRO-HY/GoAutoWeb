package ReadPathUtil

import (
	"GoAutoWeb/ReadInterceptorUtil"
	"GoAutoWeb/ReadTemplateUtil"
	"strings"
)

// Controller路由信息
type PathBean struct {

	//包所在路径
	PackagePath string

	//Http请求方法
	HttpMethod string

	//路由路径
	Path string

	//函数名
	FuncName string

	//函数名
	ReturnType string

	//该路由的参数
	Parameters []ParamBean

	//页面html相对路径
	Html string
}

// MakeHandleSource 生成Handle部分的代码
func (mine PathBean) MakeHandleSource() string {
	source := ""
	source += "\thttp.HandleFunc(\"" + mine.Path + "\", func(writer http.ResponseWriter, request *http.Request) {\n"
	if mine.HttpMethod != "REQUEST" { //需要指定请求方法的情况
		source += "\t\tif request.Method != \"" + mine.HttpMethod + "\" {\n"
		source += "\t\t\twriter.WriteHeader(http.StatusMethodNotAllowed) // 设置状态码\n"
		source += "\t\t\twriter.Write([]byte(\"Method Not Allowed\"))\n"
		source += "\t\t\treturn\n"
		source += "\t\t}\n"
	}
	source += ReadInterceptorUtil.MappingBefore(mine.Path) //执行前拦截器
	source += mine.getControllerParamSource()              // 获取Controller参数部分的代码
	source += mine.getCallMethodSource()                   // 生成调用函数部分的代码
	source += ReadInterceptorUtil.MappingAfter(mine.Path)  //执行前拦截器
	source += mine.makeWriteToSource()
	source += "\t})\n"
	return source
}

// 获取导入昵称
func (mine PathBean) GetNickImport() string {
	if len(mine.PackagePath) == 0 {
		return ""
	}
	nick := strings.ReplaceAll(mine.PackagePath, "/", "")
	nick = strings.ReplaceAll(nick, "_", "")
	return nick
}

// 获取Controller参数部分的代码
func (mine PathBean) getControllerParamSource() string {
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
func (mine PathBean) getCallMethodSource() string {

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

// 生成写入Respone部分的代码
func (mine PathBean) makeWriteToSource() string {
	templateSource := ReadTemplateUtil.ReadUseTemplatesByHtml(mine.Html)
	if len(templateSource) > 0 { //这是一个html模板路由
		return "\t\twriteToTemplate(writer, body, " + templateSource + ")\n"
	} else {
		return "\t\twriteToResponse(writer, body)\n"
	}
}
