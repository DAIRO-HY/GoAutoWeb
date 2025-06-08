package ReadMethodUtil

import (
	"GoAutoWeb/MakeGoFileInfo/Bean"
	"GoAutoWeb/MakeGoFileInfo/ReadAnnotationUtil"
	"GoAutoWeb/MakeGoFileInfo/ReadCommentUtil"
	"strings"
)

// 读取go文件中的函数
func Read(lines []string) []Bean.MethodBean {
	var methods []Bean.MethodBean
	index := -1
	for {
		index++
		if index > len(lines)-1 { //超出了范围
			break
		}
		line := lines[index]
		findLine := strings.TrimSpace(line)
		if !strings.HasPrefix(findLine, "func ") { //这是一个函数开始行
			continue
		}
		method := findMethodByStartLineNo(lines, index)

		//读取函数的注解
		method.AnnotationMap = ReadAnnotationUtil.ReadAnnotationByTargetLineNo(lines, index)

		//读取指定行以上的注解
		method.Comment = ReadCommentUtil.Read(lines, index)
		methods = append(methods, method)

		// 查找一个函数的结束行
		index = findMethodEndLineNo(lines, index)
	}
	return methods
}

// 获取函数信息
func findMethodByStartLineNo(lines []string, start int) Bean.MethodBean {
	method := Bean.MethodBean{}

	//带有函数信息的字符串
	methodInfoStr := ""
	currentNo := start
	for {
		line := lines[currentNo]
		methodInfoStr += line
		if strings.Contains(line, "{") {
			break
		}
	}
	methodInfoStr = methodInfoStr[:strings.Index(methodInfoStr, "{")]

	findNameInfoStr := methodInfoStr
	findNameInfoStr = strings.ReplaceAll(findNameInfoStr, " ", "")
	findNameInfoStr = strings.ReplaceAll(findNameInfoStr, "\t", "")

	name := ""
	parameterStr := ""
	if strings.HasPrefix(findNameInfoStr, "func(") { //这是一个结构体扩展函数
		name = methodInfoStr
		name = name[strings.Index(name, ")")+1:]
		name = name[:strings.Index(name, "(")]

		//函数参数部分的字符串
		parameterStr = methodInfoStr
		parameterStr = parameterStr[strings.Index(parameterStr, ")")+1:]
		parameterStr = parameterStr[strings.Index(parameterStr, "(")+1 : strings.Index(parameterStr, ")")]

	} else {
		name = methodInfoStr[strings.Index(methodInfoStr, "func ")+5 : strings.Index(methodInfoStr, "(")]

		//函数参数部分的字符串
		parameterStr = methodInfoStr[strings.Index(methodInfoStr, "(")+1 : strings.Index(methodInfoStr, ")")]

	}
	method.Name = name
	var parameters []Bean.VariableBean
	for _, it := range strings.Split(parameterStr, ",") {
		it = strings.TrimSpace(it)
		if len(it) == 0 {
			continue
		}
		param := Bean.VariableBean{}
		param.Name = it[:strings.Index(it, " ")]
		param.Type = it[strings.Index(it, " "):]
		param.Type = strings.TrimSpace(param.Type)
		parameters = append(parameters, param)
	}
	method.Parameters = parameters

	//返回值类型
	returnTypeStr := methodInfoStr[strings.Index(methodInfoStr, ")")+1:]
	returnTypeStr = strings.TrimSpace(returnTypeStr)
	if strings.HasPrefix(returnTypeStr, "(") {
		returnTypeStr = returnTypeStr[1 : len(returnTypeStr)-1]
	}

	var returns []string
	for _, it := range strings.Split(returnTypeStr, ",") {
		it = strings.TrimSpace(it)
		if len(it) == 0 {
			continue
		}
		returns = append(returns, it)
	}

	method.Returns = returns
	return method
}

// 查找一个函数的结束行
func findMethodEndLineNo(lines []string, start int) int {

	//代码块开始[{]出现次数
	startBlockCount := 0

	//代码块结束[}]出现次数
	endBlockCount := 0

	index := start
	for {
		line := lines[index]
		for _, ch := range line { //遍历每个字符
			if string(ch) == "{" {
				startBlockCount++
			}
			if string(ch) == "}" {
				endBlockCount++
				if startBlockCount == endBlockCount {
					return index
				}
			}
		}
		index++
	}
}
