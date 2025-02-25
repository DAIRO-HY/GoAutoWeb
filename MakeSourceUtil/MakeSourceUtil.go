package MakeSourceUtil

import (
	"GoAutoWeb/FileUtil"
	"GoAutoWeb/Global"
	"GoAutoWeb/ReadInterceptorUtil"
	"GoAutoWeb/ReadPathUtil"
	_ "embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed AutoWeb.go.txt
var autoWebSample string

func Make() {
	autoWebCode := autoWebSample
	autoWebCode = strings.ReplaceAll(autoWebCode, "//{BODY}", makeHandleBody())
	autoWebCode = strings.ReplaceAll(autoWebCode, "//{IMPORT}", makeImportSource())
	if FileUtil.ReadText(Global.RootProject+"/AutoWeb.go") != autoWebCode { //避免重复写入
		FileUtil.WriteText(Global.RootProject+"/AutoWeb.go", autoWebCode)
	}
}

// 生成Handle部分代码
func makeHandleBody() string {
	source := ""
	//for _, pathBean := range ReadPathUtil.PathList {
	//	source += pathBean.MakeHandleSource()
	//}

	//遍历所有的路由
	for path, list := range ReadPathUtil.PathMap {
		source += fmt.Sprintf("\thttp.HandleFunc(\"%s\", func(writer http.ResponseWriter, request *http.Request) {\n", path)
		for _, pb := range list {
			source += pb.MakeHandleSource()
		}
		source += `		writer.WriteHeader(http.StatusNotFound) // 设置状态码"
	})
`
	}
	return source
}

// 生成导入包的代码
func makeImportSource() string {

	//读取项目的模块名
	moduleName := Global.ModuleName
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
