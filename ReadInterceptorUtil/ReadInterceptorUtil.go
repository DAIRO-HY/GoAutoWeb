package ReadInterceptorUtil

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/FileUtil"
	"sort"
	"strconv"
	"strings"
)

// InterceptorList 拦截器列表
var InterceptorList []InterceptorBean

func Make() {
	for _, fPath := range Application.GoFileList {

		// 读取go文件里的拦截器
		list := readInterceptor(fPath)
		InterceptorList = append(InterceptorList[:], list[:]...)
	}

	//按照设置的Order大小排序，值越小越靠前
	sort.Slice(InterceptorList, func(i, j int) bool {
		return InterceptorList[i].Order < InterceptorList[j].Order
	})
}

// 匹配执行前拦截器
func MappingBefore(path string) string {
	source := ""
	for _, interceptor := range InterceptorList {
		if interceptor.HandleFlag != "before" {
			continue
		}
		isInterceptor := mappingInterceptor(path, interceptor.Include, interceptor.Exclude)
		if isInterceptor { //匹配到了拦截器
			afterSource := MappingAfter(path) //执行后拦截器
			afterSource = strings.ReplaceAll(afterSource, "\t\t", "\t\t\t")
			source += "\t\t\tif !" + interceptor.GetNickImport() + "." + interceptor.FuncName + "(writer, request) {\n\n"
			source += "\t\t\t\t// 始终都要执行后的操作\n"
			source += afterSource
			source += "\t\t\t\treturn\n"
			source += "\t\t\t}\n"
		}
	}
	return source
}

// 匹配执行后拦截器
func MappingAfter(path string) string {
	source := ""
	for _, interceptor := range InterceptorList {
		if interceptor.HandleFlag != "after" {
			continue
		}
		isInterceptor := mappingInterceptor(path, interceptor.Include, interceptor.Exclude)
		if isInterceptor { //匹配到了拦截器
			source += "\t\t\t\tbody = " + interceptor.GetNickImport() + "." + interceptor.FuncName + "(writer, request, body)\n"
		}
	}
	return source
}

// 匹配执行前拦截器
func mappingInterceptor(path string, includes []string, excludes []string) bool {
	isInterceptor := false
	for _, include := range includes { //从包含路由中匹配
		if strings.Contains(include, "**") { //匹配所有子路由
			if strings.HasPrefix(path, include[:strings.Index(include, "**")]) && strings.HasSuffix(path, include[strings.Index(include, "**")+2:]) {
				isInterceptor = true
			}
		} else if strings.HasSuffix(include, "*") { //只匹配子路由
			preInclude := include[:len(include)-1]
			if strings.HasPrefix(path, preInclude) { //首先判断前缀是否一致
				afterPath := path[len(preInclude):]
				if !strings.Contains(afterPath, "/") { //再判断路由后面是否还有子路由
					isInterceptor = true
				}
			}
		} else { //完全匹配
			if path == include {
				isInterceptor = true
			}
		}
		if isInterceptor {
			break
		}
	}
	if !isInterceptor { //不包含路由，直接跳过
		return false
	}
	for _, exclude := range excludes { //从排除路由中匹配
		if strings.HasSuffix(exclude, "**") { //匹配所有子路由
			if strings.HasPrefix(path, exclude[:len(exclude)-2]) {
				isInterceptor = false
			}
		} else if strings.HasSuffix(exclude, "*") { //只匹配子路由
			preInclude := exclude[:len(exclude)-1]
			if strings.HasPrefix(path, preInclude) { //首先判断前缀是否一致
				afterPath := path[len(preInclude):]
				if !strings.Contains(afterPath, "/") { //再判断路由后面是否还有子路由
					isInterceptor = false
				}
			}
		} else { //完全匹配
			if path == exclude {
				isInterceptor = false
			}
		}
		if !isInterceptor {
			break
		}
	}
	return isInterceptor
}

// 读取go文件里的拦截器
func readInterceptor(path string) []InterceptorBean {

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

	var interceptorList []InterceptorBean
	index := 0
	for index < len(lines) {
		line := lines[index]
		trimLine := strings.ReplaceAll(line, " ", "")
		trimLine = strings.ReplaceAll(trimLine, "\t", "")
		if !strings.HasPrefix(trimLine, "//@interceptor:") {
			index++
			continue
		}
		lineArr := strings.Split(trimLine, ":")
		interceptorBean := &InterceptorBean{
			HandleFlag: lineArr[1],

			//设置包所在路径
			PackagePath: packagePath,
		}
		for { //继续解析其他属性
			index++
			if index >= len(lines) { //超出了范围
				break
			}
			trimLine = strings.ReplaceAll(lines[index], " ", "")
			trimLine = strings.ReplaceAll(trimLine, "\t", "")
			if strings.HasPrefix(trimLine, "//@include:") { //包含路由
				include := strings.Split(trimLine, ":")[1]

				//去除所有空格
				include = strings.ReplaceAll(include, " ", "")
				interceptorBean.Include = strings.Split(include, ",")
			} else if strings.HasPrefix(trimLine, "//@exclude:") { //排除路由
				exclude := strings.Split(trimLine, ":")[1]

				//去除所有空格
				exclude = strings.ReplaceAll(exclude, " ", "")
				interceptorBean.Exclude = strings.Split(exclude, ",")
			} else if strings.HasPrefix(trimLine, "//@order:") { //优先级
				interceptorBean.Order, _ = strconv.Atoi(strings.Split(trimLine, ":")[1])
			} else if strings.HasPrefix(trimLine, "func") { //拦截器函数名
				interceptorBean.FuncName = trimLine[strings.Index(trimLine, "func")+4 : strings.Index(trimLine, "(")]
			} else {
				break
			}
		}
		interceptorList = append(interceptorList, *interceptorBean)
	}
	return interceptorList
}
