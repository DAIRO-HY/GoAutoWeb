package MakeSourceUtil

import (
	"GoAutoController/Application"
	"GoAutoController/FileUtil"
	"GoAutoController/ReadInterceptorUtil"
	"GoAutoController/ReadPathUtil"
	"sort"
	"strings"
)

func init() {

	//初始化读取路由列表
	ReadPathUtil.Make()

	//初始化读取拦截器列表
	ReadInterceptorUtil.Make()
}

func Make() {
	autoWebCode := ""
	for _, pathBean := range ReadPathUtil.PathList {
		autoWebCode += `
	http.HandleFunc("` + pathBean.Path + `", func(writer http.ResponseWriter, request *http.Request) {` +
			ReadInterceptorUtil.MappingPre(pathBean.Path) +
			pathBean.GetControllerParamSource() + // 获取Controller参数部分的代码
			pathBean.GetCallMethodSource() + // 生成调用函数部分的代码
			`
	})`
	}
	autoWebSample := FileUtil.ReadText("./AutoWeb.go")
	autoWebCode = strings.ReplaceAll(autoWebSample, "//{BODY}", autoWebCode)
	autoWebCode = strings.ReplaceAll(autoWebCode, "//{IMPORT}", makeImportSource())
	FileUtil.WriteText(Application.RootProject+"/AutoWeb.go", autoWebCode)
}

// 生成导入包的代码
func makeImportSource() string {

	//读取项目的模块名
	moduleName := Application.ModuleName
	importSourceMap := make(map[string]bool)

	//遍历路由中所有用到类的Import
	for _, pathBean := range ReadPathUtil.PathList {
		for _, paramBean := range pathBean.Parameters { //函数参数的import
			if len(paramBean.PackagePath) == 0 {
				continue
			}
			formIm := "\t" + paramBean.GetNickImport() + " \"" + moduleName + paramBean.PackagePath + "\""
			importSourceMap[formIm] = true
		}
		im := "\t" + pathBean.GetNickImport() + " \"" + moduleName + pathBean.PackagePath + "\""
		importSourceMap[im] = true
	}

	//遍历拦截器中所有用到类的Import
	for _, interceptor := range ReadInterceptorUtil.InterceptorList {
		im := "\t" + interceptor.GetNickImport() + " \"" + moduleName + interceptor.PackagePath + "\""
		importSourceMap[im] = true
	}

	importList := make([]string, 0)
	for im := range importSourceMap {
		importList = append(importList, im)
	}

	//排序一下
	sort.Strings(importList)
	importSource := ""
	for _, im := range importList {
		importSource += im + "\n"
	}
	return importSource
}
