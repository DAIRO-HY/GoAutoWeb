package ReadAnnotationUtil

import (
	"GoAutoWeb/MakeGoFileInfo/Bean"
	"strings"
	"unicode"
)

// 获取某行代码上的注解
func Read(line string) Bean.AnnotationBean {
	line = strings.TrimSpace(line)
	annotation := Bean.AnnotationBean{}
	if !strings.HasSuffix(line, ")") { //没有参数
		name := line[strings.Index(line, "@")+1:]
		name = strings.TrimSpace(name)
		annotation.Name = name
		return annotation
	}

	//获取注解名
	name := line[strings.Index(line, "@")+1 : strings.Index(line, "(")]
	name = strings.TrimSpace(name)
	annotation.Name = name

	//获取参数部分的字符串
	parameterStr := line[strings.Index(line, "("):strings.Index(line, ")")]
	valueMap := make(map[string]string)
	for _, it := range strings.Split(parameterStr, ",") {
		if strings.Contains(it, "=") {
			kv := strings.Split(it, "=")
			valueMap[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		} else { //没有参数名,视为默认参数
			valueMap["default"] = strings.TrimSpace(it)
		}
	}
	annotation.ValueMap = valueMap
	return annotation
}

// 读取指定行以上的注解
func ReadAnnotationByTargetLineNo(lines []string, targetLineNo int) []Bean.AnnotationBean {
	var annotations []Bean.AnnotationBean
	for index := targetLineNo - 1; index >= 0; index-- {
		line := lines[index]
		findLine := strings.ReplaceAll(line, " ", "")
		findLine = strings.ReplaceAll(line, "\t", "")
		if findLine == "" { //忽略换行
			continue
		}
		if unicode.IsLetter([]rune(findLine)[0]) { //如果读取到一个字母开头的代码,则结束
			break
		}
		if strings.HasSuffix(findLine, "}") { //读取到一段代码结束标记,则结束
			break
		}
		if strings.HasSuffix(findLine, ")") { //读取到一段代码结束标记,则结束
			break
		}
		if !strings.HasPrefix(findLine, "//@") { //注解一定是以//@开头
			continue
		}

		//获取该行代码上的注解
		annotations = append(annotations, Read(line))
	}
	return annotations
}
