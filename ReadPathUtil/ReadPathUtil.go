package ReadPathUtil

import (
	"GoAutoController/Application"
	"GoAutoController/FileUtil"
	"GoAutoController/bean"
	"strings"
)

// 路由列表
var PathList []bean.PathBean

func Make() {
	for _, fPath := range Application.GoFileList {

		// 读取go文件里的路由配置
		list := readControllerPath(fPath)
		PathList = append(PathList[:], list[:]...)
	}
}

// 读取go文件里的路由配置
func readControllerPath(path string) []bean.PathBean {

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

	var pathList []bean.PathBean
	index := 0
	for index < len(lines) {
		line := lines[index]

		// 解析路由
		pathBean := readPath(line)
		if pathBean != nil {

			//设置包所在路径
			pathBean.PackagePath = packagePath

			// 读取参数
			pathBean.Parameters = readParameter(pathBean.PackagePath, lines[index+1])
			pathBean.FuncName = readFuncName(lines[index+1])
			pathBean.ReturnType = readReturnType(lines[index+1])
			pathList = append(pathList, *pathBean)
		}
		index++
	}
	return pathList
}

// 解析路由
func readPath(line string) *bean.PathBean {
	trimLine := strings.ReplaceAll(line, " ", "")
	trimLine = strings.ReplaceAll(trimLine, "\t", "")

	//标记改行是否有路由标记
	var pathBean *bean.PathBean
	if strings.HasPrefix(trimLine, "//POST:") {
		pathBean = &bean.PathBean{
			Method: "POST",
			Path:   trimLine[7:],
		}
	} else if strings.HasPrefix(trimLine, "//GET:") {
		pathBean = &bean.PathBean{
			Method: "GET",
			Path:   trimLine[6:],
		}
	} else if strings.HasPrefix(trimLine, "//REQUEST:") {
		pathBean = &bean.PathBean{
			Method: "REQUEST",
			Path:   trimLine[10:],
		}
	}
	return pathBean
}

// 读取参数
func readParameter(goPackagePath string, line string) []bean.ParamBean {
	paramStr := line[strings.Index(line, "(")+1 : strings.Index(line, ")")]
	if strings.TrimSpace(paramStr) == "" { //不需要参数
		return []bean.ParamBean{}
	}
	paramArr := strings.Split(paramStr, ",")
	var paramList []bean.ParamBean
	for _, param := range paramArr {
		paramInfoArr := strings.Split(strings.TrimSpace(param), " ")
		varType := paramInfoArr[1]
		packagePath := ""
		if strings.HasPrefix(varType, "form.") {
			varType = varType[5:]
			packagePath = goPackagePath + "/form"
		}
		paramBean := bean.ParamBean{
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
