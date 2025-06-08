package ReadStructUtil

import (
	"GoAutoWeb/MakeGoFileInfo/Bean"
	"GoAutoWeb/MakeGoFileInfo/ReadAnnotationUtil"
	"GoAutoWeb/MakeGoFileInfo/ReadCommentUtil"
	"strings"
	"unicode"
)

// 读取go文件中的函数
func Read(lines []string) []Bean.StructBean {
	var stts []Bean.StructBean
	index := -1
	for {
		index++
		if index > len(lines)-1 { //超出了范围
			break
		}
		line := lines[index]
		findLine := strings.TrimSpace(line)
		if strings.HasPrefix(findLine, "type ") && strings.Contains(findLine, " struct") { //这是一个结构体开始
			stt := Bean.StructBean{}

			//读取结构体名
			stt.Name = readStructNameByStartLineNo(lines, index)

			//读取结构体注解
			stt.AnnotationMap = ReadAnnotationUtil.ReadAnnotationByTargetLineNo(lines, index)

			//读取结构体注解
			stt.Comment = ReadCommentUtil.Read(lines, index)

			//读取成员变量
			stt.Members = readStructMemberByStartLineNo(lines, index)
			stts = append(stts, stt)

			// 查找一个函数的结束行
			index = findStructEndLineNo(lines, index)
		}
	}
	return stts
}

// 读取结构体名
func readStructNameByStartLineNo(lines []string, start int) string {

	//带有函数信息的字符串
	structInfoStr := ""
	currentNo := start
	for {
		line := lines[currentNo]
		structInfoStr += line
		if strings.Contains(line, "{") {
			break
		}
	}
	structInfoStr = structInfoStr[:strings.Index(structInfoStr, "{")]

	//得到结构体名
	name := structInfoStr[strings.Index(structInfoStr, "type ")+5 : strings.Index(structInfoStr, " struct")]
	name = strings.TrimSpace(name)
	return name
}

// 读取成员变量
func readStructMemberByStartLineNo(lines []string, start int) []Bean.VariableBean {
	var members []Bean.VariableBean

	//结构体结束行
	structEndLineNo := findStructEndLineNo(lines, start)
	for i := start + 1; i < structEndLineNo; i++ {
		line := lines[i]
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if !unicode.IsLetter([]rune(line)[0]) { // 如果第一个字符是字母,则代表这是一个变量）
			continue
		}
		member := Bean.VariableBean{}
		memberStrArr := strings.Split(line, " ")
		member.Name = memberStrArr[0]
		member.Type = strings.TrimSpace(memberStrArr[0])
		member.Comment = ReadCommentUtil.Read(lines, i)                                  //读取到注释
		member.AnnotationMap = ReadAnnotationUtil.ReadAnnotationByTargetLineNo(lines, i) //读取注解
		members = append(members, member)
	}
	return members
}

// 查找一个函数的结束行
func findStructEndLineNo(lines []string, start int) int {

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
