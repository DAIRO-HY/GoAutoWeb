package ReadCommentUtil

import (
	"strings"
)

// 读取指定行以上的注解
func Read(lines []string, targetIndex int) string {
	var comments []string

	for index := targetIndex - 1; index >= 0; index-- {
		line := lines[index]
		findLine := strings.ReplaceAll(line, " ", "")
		findLine = strings.ReplaceAll(line, "\t", "")
		if findLine == "\n" { //忽略换行
			continue
		}
		if strings.HasPrefix(findLine, "//@") { //注解一定是以//@开头,跳过
			continue
		}
		if strings.HasPrefix(findLine, "//") { //单行注解
			comments = append(comments, line)
		} else if strings.HasSuffix(findLine, "*/") { //多行注释
			for ; index >= 0; index-- {
				line = lines[index]
				comments = append(comments, line)
				findLine = strings.ReplaceAll(line, " ", "")
				findLine = strings.ReplaceAll(line, "\t", "")
				if strings.HasPrefix(findLine, "/*") {
					break
				}
			}
		} else { //这已经不是注释代码
			break
		}
	}
	comment := ""
	for i := len(comments) - 1; i >= 0; i-- {
		comment += comments[i] + "\n"
	}
	return comment
}
