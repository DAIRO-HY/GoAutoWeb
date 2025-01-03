package ReadPathUtil

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/FileUtil"
	"strings"
)

// 路由列表
var PathList []PathBean

func Make() {
	for _, fPath := range Application.GoFileList {

		// 读取go文件里的路由配置
		list := readControllerPath(fPath)
		PathList = append(PathList[:], list[:]...)
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
	index := 0
	for index < len(lines) {
		line := lines[index]

		// 解析路由
		pathBean := readPath(line)
		if pathBean != nil {

			//设置包所在路径
			pathBean.PackagePath = packagePath
			for {
				index++
				line = lines[index]
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "func") {

					// 读取参数
					pathBean.Parameters = readParameter(pathBean.PackagePath, line)
					pathBean.FuncName = readFuncName(line)
					pathBean.ReturnType = readReturnType(line)
					break
				} else if strings.Contains(line, "@templates:") {
					pathBean.Templates = readTemplate(line)
				} else {
				}
			}
			pathList = append(pathList, *pathBean)
		}
		index++
	}
	return pathList
}

// 读取要使用template模板（html专用）
func readTemplate(line string) []string {
	lineArr := strings.Split(line, ":")
	templates := strings.Split(lineArr[1], ",")
	for i, template := range templates { //去掉空格
		templates[i] = strings.TrimSpace(template)
	}
	return templates
}

// 解析路由
func readPath(line string) *PathBean {
	trimLine := strings.ReplaceAll(line, " ", "")
	trimLine = strings.ReplaceAll(trimLine, "\t", "")

	//标记改行是否有路由标记
	var pathBean *PathBean
	if strings.HasPrefix(trimLine, "//@post:") {
		pathBean = &PathBean{
			Method: "post",
			Path:   trimLine[8:],
		}
	} else if strings.HasPrefix(trimLine, "//@get:") {
		pathBean = &PathBean{
			Method: "get",
			Path:   trimLine[7:],
		}
	} else if strings.HasPrefix(trimLine, "//@request:") {
		pathBean = &PathBean{
			Method: "request",
			Path:   trimLine[11:],
		}
	} else {
		return nil
	}
	return pathBean
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
		paramInfoArr := strings.Split(strings.TrimSpace(param), " ")
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
