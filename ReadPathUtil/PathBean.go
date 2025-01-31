package ReadPathUtil

import (
	"GoAutoWeb/ReadInterceptorUtil"
	"GoAutoWeb/ReadTemplateUtil"
	"regexp"
	"slices"
	"strconv"
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

	//路由参数路径
	VariablePath string

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
	pathVariableList := mine.getPathVariableList()
	getFormParamtaSource := ""  //获取表单参数部分代码
	getPathVariableSource := "" //获取路由参数部分代码
	for _, parameter := range mine.Parameters {
		if slices.Contains(pathVariableList, parameter.Name) { //这是一个url路径参数
			index := slices.Index(pathVariableList, parameter.Name)
			getPathVariableSource += parameter.MakeGetPathVariableSource(index)
		} else {
			getFormParamtaSource += parameter.MakeGetParameterSource()
		}
	}
	if len(getFormParamtaSource) > 0 { //生成URL参数和Body参数变量代码
		getFormParamtaSource = "\t\trequestFormData := getRequestFormData(request) //获取表单数据\n" + getFormParamtaSource
	}
	if len(getPathVariableSource) > 0 { //如果有路径参数
		getPathVariableSource = mine.getPathVariableListSource() + getPathVariableSource
	}
	return getFormParamtaSource + getPathVariableSource
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

// 获取url参数名名称
func (mine PathBean) getPathVariableList() []string {
	pathVarList := make([]string, 0)
	findResults := regexp.MustCompile(`\{([^}]+)\}`).FindAllString(mine.VariablePath, -1)
	for _, it := range findResults {
		pathVarList = append(pathVarList, it[1:len(it)-1])
	}
	return pathVarList
}

// 获取路由参数数组部分的代码
func (mine PathBean) getPathVariableListSource() string {

	// 定义正则表达式，匹配所有被 {} 包围的内容
	re := regexp.MustCompile(`\{[^}]+\}`)

	// 使用 ReplaceAllString 方法将所有匹配的内容替换为空格
	replacePath := re.ReplaceAllString(mine.VariablePath, " ")
	pathVariableSplits := strings.Split(replacePath, " ")
	pathVariableSplitStr := strings.Join(pathVariableSplits, `", "`)

	hasPrefixSource := ""
	hasSuffixSource := ""
	if pathVariableSplits[0] != "" {
		hasPrefixSource = `if !strings.HasPrefix(varPath, "` + pathVariableSplits[0] + `") { //如果不是以指定的字符串开头
			writer.WriteHeader(http.StatusNotFound)
			return
		}`
	}
	if pathVariableSplits[len(pathVariableSplits)-1] != "" {
		hasSuffixSource = `if !strings.HasSuffix(varPath, "` + pathVariableSplits[len(pathVariableSplits)-1] + `") {//如果不是以指定的字符串结尾
			writer.WriteHeader(http.StatusNotFound)
			return
		}`
	}

	return `
		pathVariables := make([]string, 0)
		varPath := request.URL.Path[` + strconv.Itoa(len(mine.Path)) + `:]
		pathVariableSplitArr := []string{"` + pathVariableSplitStr + `"}
		` + hasPrefixSource + `
		` + hasSuffixSource + `
		for i := 0; i < len(pathVariableSplitArr)-1; i++ {
			varPath = varPath[len(pathVariableSplitArr[i]):]
			if pathVariableSplitArr[i+1] == "" { //这已经是最后一个参数了
				pathVariables = append(pathVariables, varPath)
			} else {
				nextIndex := strings.Index(varPath, pathVariableSplitArr[i+1])
				if nextIndex == -1 {
					writer.WriteHeader(http.StatusNotFound)
					return
				}
				pathVariables = append(pathVariables, varPath[:nextIndex])
				varPath = varPath[nextIndex:]
			}
		}` + "\n"
}
