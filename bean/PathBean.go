package bean

import "strings"

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
func (mine *PathBean) GetControllerParamSource() string {
	source := ""
	for _, parameter := range mine.Parameters {
		if parameter.VarType == "http.ResponseWriter" { //这不是一个URL参数
			continue
		}
		if parameter.VarType == "*http.Request" { //这不是一个URL参数
			continue
		}
		if parameter.VarType == "string" { //字符串类型
			source += "\n\t\t" + parameter.Name + " := getString(paramMap, \"" + parameter.Name + "\")"
		} else if parameter.VarType == "int" {
			source += "\n\t\t" + parameter.Name + " := getInt(paramMap, \"" + parameter.Name + "\")"
		} else if parameter.VarType == "int64" {
			source += "\n\t\t" + parameter.Name + " := getInt64(paramMap, \"" + parameter.Name + "\")"
		} else if parameter.VarType == "float32" {
			source += "\n\t\t" + parameter.Name + " := getFloat32(paramMap, \"" + parameter.Name + "\")"
		} else if parameter.VarType == "float64" {
			source += "\n\t\t" + parameter.Name + " := getFloat64(paramMap, \"" + parameter.Name + "\")"
		} else if strings.HasSuffix(parameter.VarType, "Form") { //这是一个结构体Form表单
			source += "\n\t\t" + parameter.Name + " := getForm[" + parameter.GetNickImport() + "." + parameter.VarType + "](paramMap)"
		}
	}
	if len(source) > 0 {
		source = "\t\tparamMap := makeParamMap(request)" + source + "\n"
	}
	return source
}

// 生成调用函数部分的代码
func (mine *PathBean) GetCallMethodSource() string {

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
func (mine *PathBean) GetEndSource() string {
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
