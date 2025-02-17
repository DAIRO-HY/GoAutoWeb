package ReadPathUtil

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/FileUtil"
	"strings"
)

// 路由列表
var PathList []PathBean

// 路由对应的路由信息列表
var PathMap map[string][]PathBean

func Make() {
	for _, fPath := range Application.GoFileList {

		// 读取go文件里的路由配置
		list := readControllerPath(fPath)
		PathList = append(PathList[:], list[:]...)
	}
	PathMap = make(map[string][]PathBean)

	// 生成路由对应的路由信息列表
	for _, it := range PathList {
		list := PathMap[it.Path]
		if list != nil {
			list = append(list, it)
		} else {
			list = []PathBean{it}
		}
		PathMap[it.Path] = list
	}
}

// 读取go文件里的路由配置
func readControllerPath(path string) []PathBean {

	//包所在路径
	packagePath := path[len(Application.RootProject):]
	packagePath = strings.ReplaceAll(packagePath, "\\", "/")
	packagePath = packagePath[:strings.LastIndex(packagePath, "/")]

	//读取代码
	goCode := FileUtil.ReadText(path)

	//先统一换行符
	goCode = strings.ReplaceAll(goCode, "\r", "\r\n")
	goCode = strings.ReplaceAll(goCode, "\r\n", "\n")
	for strings.Contains(goCode, "\n\n") { //将连续的换行符替换成一个
		goCode = strings.ReplaceAll(goCode, "\n\n", "\n")
	}
	lines := strings.Split(goCode, "\n")

	var pathList []PathBean
	pathBean := &PathBean{}
	group := "" //分组名
	for index, line := range lines {
		if group == "" {
			group = readGroup(line)
		}
		readPath(line, pathBean)                          // 解析路由
		readHtml(line, pathBean)                          //解析html页面
		readFunction(index, lines, packagePath, pathBean) //读取调用的函数
		if pathBean.FuncName != "" {                      //已经读取到了调用的函数
			if pathBean.HttpMethod == "" && pathBean.Html != "" { //如果有配置html页面，但没有配置HTTP请求方式，则将html路径作为路由
				pathBean.HttpMethod = "GET"
				if strings.HasPrefix(pathBean.Html, ".") {
					pathBean.Path = pathBean.Html
				} else {
					pathBean.Path = "/" + pathBean.Html
				}
				pathBean.Html = group + pathBean.Html
			}
			if pathBean.HttpMethod != "" { //路由标记的才是对象
				pathBean.PackagePath = packagePath

				//判断path中有没有参数路由
				path = group + pathBean.Path
				pathVariableStartSplitCharIndex := strings.Index(path, "{")
				if pathVariableStartSplitCharIndex != -1 { //有路由参数

					//获取路由参数之前的最后一个路径分隔符位置
					lastPathIndex := strings.LastIndex(path[:pathVariableStartSplitCharIndex], "/") + 1
					pathBean.Path = path[:lastPathIndex]
					pathBean.VariablePath = path[lastPathIndex:]
				} else {
					pathBean.Path = path
				}
				pathList = append(pathList, *pathBean)
			}
			pathBean = &PathBean{}
		}
	}
	return pathList
}

// 解析路由
func readGroup(line string) string {
	trimLine := strings.ReplaceAll(line, " ", "")
	trimLine = strings.TrimSpace(trimLine)
	trimLineUppercase := strings.ToUpper(trimLine)         //忽略大小写
	if strings.HasPrefix(trimLineUppercase, "//@GROUP:") { //读取到一个分组
		return trimLine[strings.Index(trimLine, ":")+1:]
	}
	return ""
}

// 解析路由
func readPath(line string, bean *PathBean) {
	trimLine := strings.ReplaceAll(line, " ", "")
	trimLine = strings.TrimSpace(trimLine)
	trimLineUppercase := strings.ToUpper(trimLine) //忽略大小写

	//标记改行是否有路由标记
	if strings.HasPrefix(trimLineUppercase, "//@POST:") {
		bean.HttpMethod = "POST"
		bean.Path = trimLine[strings.Index(trimLine, ":")+1:]
	} else if strings.HasPrefix(trimLineUppercase, "//@GET:") {
		bean.HttpMethod = "GET"
		bean.Path = trimLine[strings.Index(trimLine, ":")+1:]
	} else if strings.HasPrefix(trimLineUppercase, "//@REQUEST:") {
		bean.HttpMethod = "REQUEST"
		bean.Path = trimLine[strings.Index(trimLine, ":")+1:]
	} else {
	}
}

// 读取页面html配置
func readHtml(line string, bean *PathBean) {
	lineUppercase := strings.ToUpper(line) //忽略大小写
	lineUppercase = strings.ReplaceAll(lineUppercase, " ", "")
	if !strings.Contains(lineUppercase, "@HTML:") {
		return
	}
	html := line[strings.Index(line, ":")+1:]
	bean.Html = strings.TrimSpace(html)
}

// 读取调用的函数
func readFunction(index int, lines []string, packagePath string, bean *PathBean) {
	line := strings.TrimSpace(lines[index])
	if !strings.HasPrefix(line, "func") {
		return
	}
	for !strings.Contains(line, "{") { //如果该行没有{，说明函数的参数已换行处理
		index++
		line += strings.TrimSpace(lines[index])
	}

	// 读取参数
	bean.Parameters = readParameter(packagePath, line)
	bean.FuncName = readFuncName(line)
	bean.ReturnType = readReturnType(line)
}

// 读取参数
func readParameter(goPackagePath string, line string) []ParamBean {
	paramStr := line[strings.Index(line, "(")+1 : strings.Index(line, ")")]
	if strings.TrimSpace(paramStr) == "" { //不需要参数
		return []ParamBean{}
	}
	paramArr := strings.Split(paramStr, ",")
	var paramList []ParamBean
	for _, param := range paramArr {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}
		paramInfoArr := strings.Split(param, " ")
		if len(paramInfoArr) < 2 { //这不是一个正常的参数
			continue
		}
		varType := paramInfoArr[1]
		packagePath := ""
		if strings.HasPrefix(varType, "form.") {
			varType = varType[5:]
			packagePath = goPackagePath + "/form"
		}
		paramBean := ParamBean{
			PackagePath: packagePath,
			VarType:     varType,
			Name:        paramInfoArr[0],
		}
		paramList = append(paramList, paramBean)
	}
	return paramList
}

// 读取函数名
func readFuncName(line string) string {
	funcName := line[strings.Index(line, "func")+5 : strings.Index(line, "(")]
	return funcName
}

// 读取返回值
func readReturnType(line string) string {
	returnType := line[strings.Index(line, ")")+1 : strings.Index(line, "{")]
	return strings.TrimSpace(returnType)
}
