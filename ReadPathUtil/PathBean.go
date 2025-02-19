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
	source += "\t\t\tvar body any = nil\n"
	source += ReadInterceptorUtil.MappingBefore(mine.Path) // 执行前拦截器
	source += mine.getControllerParamSource()              // 获取Controller参数部分的代码
	source += mine.makeDeferSource()                       // 生成最终执行代码
	source += mine.getCallMethodSource()                   // 生成调用函数部分的代码
	//source += ReadInterceptorUtil.MappingAfter(mine.Path) // 执行后拦截器
	//source += mine.makeWriteToSource()
	source += "\t\t\treturn\n"

	if len(mine.getPathVariableList()) > 0 { //如果有路由参数
		pathVariableSplitStr := strings.Join(mine.getPathVariableSplitArr(), "\", \"")
		source = `
		pathVariableSplitArr := []string{"` + pathVariableSplitStr + `"}
		varPath := request.URL.Path[` + strconv.Itoa(len(mine.Path)) + `:]
		if isPathVariable(varPath, pathVariableSplitArr){// 判断是否满足定义的路由参数规则
			` + source + `
		}
`
	}

	if mine.HttpMethod != "REQUEST" { //需要指定请求方法的情况
		source = "\t\t\tif request.Method == \"" + mine.HttpMethod + "\" {\n" +
			source +
			"\t\t}\n"
	}
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
	getFormParamtaSource := ""  //获取表单参数部分代码
	getPathVariableSource := "" //获取路由参数部分代码
	pathVariableList := mine.getPathVariableList()
	for _, parameter := range mine.Parameters {
		if slices.Contains(pathVariableList, parameter.Name) { //这是一个url路由参数
			index := slices.Index(pathVariableList, parameter.Name)
			getPathVariableSource += parameter.MakeGetPathVariableParameterSource(index)
		} else {
			getFormParamtaSource += parameter.MakeGetParameterSource()
		}
	}
	if len(getFormParamtaSource) > 0 { //生成URL参数和Body参数变量代码
		getFormParamtaSource = "\t\t\trequestFormData := getRequestFormData(request) //获取表单数据\n" + getFormParamtaSource
	}
	if len(getPathVariableSource) > 0 { //如果这是一个路由参数
		getPathVariableSource = makePathVariableParameterSource + getPathVariableSource
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

	//source := "\t\tvar body any = nil\n"
	source := ""
	if len(mine.ReturnType) > 0 { //如果有返回值
		source += "\t\t\tbody = " + callMethodSource
	} else {
		source += "\t\t\t" + callMethodSource
	}
	return source + "\n"
}

// 生成最终执行代码
func (mine PathBean) makeDeferSource() string {

	// 执行后拦截器
	afterSource := ReadInterceptorUtil.MappingAfter(mine.Path)

	// 生成写入Response部分的代码
	writeSource := mine.makeWriteToSource()
	return `
			defer func() {
				if panicErr := recover(); panicErr != nil { // 程序终止异常全局捕获
					body = panicErr
				}
` + afterSource + writeSource + `			}()
`
}

// 生成写入Response部分的代码
func (mine PathBean) makeWriteToSource() string {
	templateSource := ReadTemplateUtil.ReadUseTemplatesByHtml(mine.Html)
	if len(templateSource) > 0 { //这是一个html模板路由
		return "\t\t\t\twriteToTemplate(writer, body, " + templateSource + ")\n"
	} else {
		return "\t\t\t\twriteToResponse(writer, body)\n"
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

// 获取路由参数分割参数的字符串列表
func (mine PathBean) getPathVariableSplitArr() []string {

	// 定义正则表达式，匹配所有被 {} 包围的内容
	re := regexp.MustCompile(`\{[^}]+\}`)

	// 使用 ReplaceAllString 方法将所有匹配的内容替换为空格
	replacePath := re.ReplaceAllString(mine.VariablePath, " ")
	return strings.Split(replacePath, " ")
}

// 获取路由参数数组部分的代码
const makePathVariableParameterSource = `
		pathVariables := make([]string, 0)
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
