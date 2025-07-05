package MakeGoFileInfo

import (
	"GoAutoWeb/FileUtil"
	"GoAutoWeb/MakeGoFileInfo/GoBean"
	"GoAutoWeb/MakeGoFileInfo/ReadAnnotationUtil"
	"GoAutoWeb/MakeGoFileInfo/ReadMethodUtil"
	"GoAutoWeb/MakeGoFileInfo/ReadStructUtil"
	"fmt"
	"strings"
)

func Test(goList []GoBean.GoClass) {
	//for _, gb := range goList {
	//	println("--------------------------------------------->Path:" + gb.Path)
	//	println("-->Package:" + gb.Package)
	//	for _, imt := range gb.Imports {
	//		println("-->Import:" + imt)
	//	}
	//	for _, anno := range gb.Annotations {
	//		println("-->Annotation.Name:@" + anno.Name)
	//		for k, v := range anno.ValueMap {
	//			println("-->Annotation.ValueMap." + k + "=" + v)
	//		}
	//	}
	//	for _, mth := range gb.Methods {
	//		println("-->Method.Name:" + mth.Name + "()")
	//		println("-->Method.Comment:" + mth.Comment)
	//		println("-->Method.Returns:" + strings.Join(mth.Returns, ","))
	//		for _, anno := range mth.Annotations {
	//			println("-->Method.Annotation:@" + anno.Name)
	//		}
	//		for _, param := range mth.Parameters {
	//			println("-->Method.param:" + param.Name + " " + param.Type)
	//		}
	//	}
	//
	//	for _, stt := range gb.Structs {
	//		println("-->Struct.Name:" + stt.Name + "{}")
	//		println("-->Struct.Comment:" + stt.Comment)
	//		for _, anno := range stt.Annotations {
	//			println("-->Struct.Annotation:@" + anno.Name)
	//		}
	//		for _, mem := range stt.Members {
	//			println("-->Struct.Members:" + mem.Name + " " + mem.Type)
	//			println("-->Struct.Members.Comment:" + mem.Comment)
	//			for _, anno := range stt.Annotations {
	//				println("-->Struct.Members.Annotation:@" + anno.Name)
	//			}
	//		}
	//	}
	//}
}

// 读取Go代码信息
func ReadGoInfo(path string) GoBean.GoClass {
	goBean := GoBean.GoClass{}
	goBean.FilePath = strings.ReplaceAll(path, "\\", "/")
	goCode := FileUtil.ReadText(path)

	//先统一换行符
	goCode = strings.ReplaceAll(goCode, "\r", "\r\n")
	goCode = strings.ReplaceAll(goCode, "\r\n", "\n")

	lines := strings.Split(goCode, "\n")

	//获取包名
	goBean.Package = readPackage(lines)

	//获取导入的包
	goBean.Imports = readImport(lines)

	// 读取go文件上的注解
	goBean.AnnotationMap = readAnnotation(lines)

	//读取go文件中的函数
	goBean.Methods = ReadMethodUtil.Read(lines)

	//读取结构体
	goBean.Structs = ReadStructUtil.Read(lines)
	return goBean
}

// 获取包名
func readPackage(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "package ") {
			pkg := line[8:]
			return strings.TrimSpace(pkg)
		}
	}
	return ""
}

// 获取导入的包
func readImport(lines []string) []string {
	importStartLineNo := -1
	for i, line := range lines { //寻找import开始行
		if strings.HasPrefix(line, "import ") {
			importStartLineNo = i
			break
		}
	}
	if importStartLineNo == -1 { //没有找到import
		return nil
	}
	importSrc := strings.Join(lines, "\n")
	importSrc = importSrc[strings.Index(importSrc, "import ")+7:]
	importSrc = strings.TrimSpace(importSrc)
	if strings.HasPrefix(importSrc, "(") { //这是一个多行Import
		importSrc = importSrc[strings.Index(importSrc, "(")+1 : strings.Index(importSrc, ")")]
	} else { //只有一个Import
		importSrc = importSrc[strings.Index(importSrc, "\"")+1:]
		importSrc = importSrc[:strings.Index(importSrc, "\"")]
		importSrc = "\"" + importSrc + "\""
	}
	//
	//importStr := lines[importStartLineNo]
	//for index := importStartLineNo + 1; index < len(lines); index++ { //获取import部分的字符串
	//	line := lines[index]
	//	line = strings.TrimSpace(line)
	//	if len(line) == 0 {
	//		continue
	//	}
	//	if unicode.IsLetter([]rune(line)[0]) { //下一个代码块开始
	//		break
	//	}
	//	if strings.HasPrefix(line, "/") { //注解
	//		continue
	//	}
	//	if strings.HasPrefix(line, "*") { //注解
	//		continue
	//	}
	//	importStr += line + "\n"
	//}
	//if strings.Contains(importStr, "(") {
	//	importStr = importStr[strings.Index(importStr, "(")+1 : strings.Index(importStr, ")")]
	//} else { //只有一个import时没有()
	//	importStr = importStr[8:]
	//}

	imports := make([]string, 0)
	for _, line := range strings.Split(importSrc, "\n") {
		if !strings.Contains(line, "\"") {
			continue
		}
		if strings.Index(line, "\"") == 0 && strings.LastIndex(line, "\"") == 0 {
			fmt.Print("dfsf")
		}

		//取引号之间的字符串
		imt := line[strings.Index(line, "\"")+1 : strings.LastIndex(line, "\"")]
		imports = append(imports, imt)
	}
	return imports
}

// 读取go文件上的注解
func readAnnotation(lines []string) map[string]GoBean.GoAnnotation {
	for i, line := range lines { //寻找import开始行
		if strings.HasPrefix(strings.TrimSpace(line), "package ") { //忽略换行
			return ReadAnnotationUtil.ReadAnnotationByTargetLineNo(lines, i)
		}
	}
	return nil
}
