package MakeSourceUtil

import (
	"GoAutoController/Application"
	"GoAutoController/FileUtil"
	"GoAutoController/ReadPathUtil"
	"GoAutoController/bean"
	"strings"
)

func Start() {
	goFileList := FileUtil.GetGoFile(Application.RootProject)

	var pathList []bean.PathBean
	for _, fPath := range goFileList {

		// 读取go文件里的路由配置
		list := MakeSourceUtil.ReadControllerPath(fPath)
		pathList = append(pathList[:], list[:]...)
	}

	autoWebCode := ""
	for _, pathBean := range pathList {

		//生成获取参数部分的代码
		getParamSource := makeGetParamSource(pathBean.Parameters)

		// 生成调用函数部分的代码
		callMethodSource := makeCallMethodSource(pathBean)

		autoWebCode += `
	http.HandleFunc("` + pathBean.Path + `", func(writer http.ResponseWriter, request *http.Request) {` + getParamSource + callMethodSource + `
	})`
	}

	autoWebSample := FileUtil.ReadText("./AutoWeb.go")
	autoWebCode = strings.ReplaceAll(autoWebSample, "//{BODY}", autoWebCode)
	FileUtil.WriteText(Application.RootProject+"/AutoWeb.go", autoWebCode)
}

// 生成获取参数部分的代码
func makeGetParamSource(parameters []bean.ParamBean) string {
	source := ""
	for _, parameter := range parameters {
		if parameter.VarType == "http.ResponseWriter" {
			continue
		}
		if parameter.VarType == "*http.Request" {
			continue
		}
		if parameter.VarType == "string" { //字符串类型
			source += "\n\t\t" + parameter.Name + " := query.Get(\"" + parameter.Name + "\")"
		} else if parameter.VarType == "int" {
			source += "\n\t\t" + parameter.Name + ",_ := strconv.Atoi(query.Get(\"" + parameter.Name + "\"))"
		} else if parameter.VarType == "int64" {
			source += "\n\t\t" + parameter.Name + ",_ := strconv.ParseInt(query.Get(\"" + parameter.Name + "\"),10,64)"
		} else if parameter.VarType == "float32" {
			source += "\n\t\t" + parameter.Name + "_64,_ := strconv.ParseFloat(query.Get(\"" + parameter.Name + "\"),32)"
			source += "\n\t\t" + parameter.Name + " := float32(" + parameter.Name + "_64)"
		} else if parameter.VarType == "float64" {
			source += "\n\t\t" + parameter.Name + ",_ := strconv.ParseFloat(query.Get(\"" + parameter.Name + "\"),64)"
		} else if strings.HasSuffix(parameter.VarType, "Form") { //这是一个结构体Form表单
			source += "\n\t\t" + parameter.Name + " := getForm[" + parameter.VarType + "](request)"
		}
	}
	if len(source) > 0 {
		source = "\n\t\tquery := request.URL.Query()" + source
	}
	return source
}

// 生成调用函数部分的代码
func makeCallMethodSource(pathBean bean.PathBean) string {
	source := pathBean.FuncName + "("
	for _, parameter := range pathBean.Parameters { //传递参数
		if parameter.VarType == "http.ResponseWriter" {
			source += "writer, "
			continue
		}
		if parameter.VarType == "*http.Request" {
			source += "request, "
			continue
		}
		source += parameter.Name + ", "
	}
	source = source[:len(source)-2]
	source += ")"

	if len(pathBean.ReturnType) > 0 { //如果有返回值
		source = "body := " + source
		source += "\n\t\twriteToResponse(writer, body)"
	}
	source = "\n\t\t" + source
	return source
}
