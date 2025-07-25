package MakeModel

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/FileUtil"
	"GoAutoWeb/Global"
	"GoAutoWeb/MakeGoFileInfo/GoBean"
	"encoding/json"
	"strings"
)

// 导入的Model列表
var importMap = make(map[string]struct{})

// 生成Model类文件
func Make() {
	jsonData, _ := json.Marshal(Global.GoClassList)

	// 复制一个对象操作，避免指针操作修改到原始数据
	copyList := make([]GoBean.GoClass, 0)
	json.Unmarshal(jsonData, &copyList)
	for _, goClass := range copyList {
		for _, goStruct := range goClass.Structs {
			if !strings.HasSuffix(goStruct.Name, "Form") {
				continue
			}
			makeModelByForm(goStruct)
		}
	}
}

// 通过一个Form表单生成Model文件
func makeModelByForm(goStruct GoBean.GoStruct) {

	//修改类名
	goStruct.Name = goStruct.Name[0:len(goStruct.Name)-4] + "Model"
	body := makeVarBodySource(goStruct) // + makeConstructorSource(form)
	source := `/*工具自动生成代码,请勿手动修改*/
struct ` + goStruct.Name + ` : Codable {
` + body + "}\n"
	save(source, goStruct.Name)
	importMap = make(map[string]struct{})
}

// go数据类型转dart数据类型
func goTypeToSwiftType(varType string) string {
	swiftType := ""
	switch varType {
	case "int":
		swiftType = "Int"
	case "int8":
		swiftType = "Int8"
	case "int16":
		swiftType = "Int16"
	case "int32":
		swiftType = "Int32"
	case "int64":
		swiftType = "Int64"
	case "string":
		swiftType = "String"
	case "bool":
		swiftType = "Bool"
	default:
		if strings.HasPrefix(varType, "[]") {
			listFormName := varType[2:]
			if strings.HasSuffix(listFormName, "Form") { //如果是以Form结尾的类名
				listFormName = listFormName[:len(listFormName)-4] + "Model"
				importMap[listFormName] = struct{}{}
			} else {
				listFormName = goTypeToSwiftType(listFormName)
			}
			swiftType = "[" + listFormName + "]"
		} else if strings.HasSuffix(varType, "Form") {
			swiftType = varType[:len(varType)-4] + "Model"
			importMap[swiftType] = struct{}{}
		} else {
			swiftType = varType
		}
	}
	return swiftType
}

// 生成变量定义部分的代码
func makeVarBodySource(goStruct GoBean.GoStruct) string {
	source := ""
	for _, it := range goStruct.Members {
		source += "\n  " + it.Comment + "  var " + it.LowerName() + ": " + goTypeToSwiftType(it.Type) + "\n"
	}
	return source + "\n"
}

// 生成构造函数部分的代码
//func makeConstructorSource(form ReadFormUtil.FormBean) string {
//	source := ""
//	for _, it := range form.Properties {
//		source += "      required this." + it.LowerName() + ",\n"
//	}
//	source = source[:len(source)-2]
//	return "  " + form.Name + "(\n      {" + source + "});\n"
//}

// 保存文件
func save(source string, fileName string) {
	path := Application.Args.TargetDir + "/Model/" + fileName + ".swift"
	fileContent := FileUtil.ReadText(path)
	if fileContent == source {
		return
	}
	FileUtil.WriteText(path, source)
}
