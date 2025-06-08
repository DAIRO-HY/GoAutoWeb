package MakeGoFileInfo

import (
	"GoAutoWeb/FileUtil"
	"GoAutoWeb/MakeGoFileInfo/Bean"
	"GoAutoWeb/MakeGoFileInfo/ReadAnnotationUtil"
	"GoAutoWeb/MakeGoFileInfo/ReadMethodUtil"
	"GoAutoWeb/MakeGoFileInfo/ReadStructUtil"
	"strings"
)

func Test(goList []Bean.GoBean) {
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
func ReadGoInfo(path string) Bean.GoBean {
	goBean := Bean.GoBean{}
	goBean.Path = strings.ReplaceAll(path, "\\", "/")
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
	importStr := ""
	for index := importStartLineNo; index < len(lines); index++ { //获取import部分的字符串
		line := lines[index]
		importStr += line + "\n"
		if strings.Contains(line, ")") {
			break
		}
	}
	importStr = importStr[strings.Index(importStr, "(")+1 : strings.Index(importStr, ")")]

	imports := make([]string, 0)
	for _, line := range strings.Split(importStr, "\n") {
		if !strings.Contains(line, "\"") {
			continue
		}

		//取引号之间的字符串
		imt := line[strings.Index(line, "\"")+1 : strings.LastIndex(line, "\"")]
		imports = append(imports, imt)
	}
	return imports
}

// 读取go文件上的注解
func readAnnotation(lines []string) map[string]Bean.AnnotationBean {
	for i, line := range lines { //寻找import开始行
		if strings.HasPrefix(strings.TrimSpace(line), "package ") { //忽略换行
			return ReadAnnotationUtil.ReadAnnotationByTargetLineNo(lines, i)
		}
	}
	return nil
}
