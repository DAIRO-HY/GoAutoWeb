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
	source := "/*工具自动生成代码,请勿手动修改*/\n"
	source += "package " + Application.Args.TargetPackage + ".model\n\n"
	source += "class " + goStruct.Name + "{\n" + body + "}\n"
	save(source, goStruct.Name)
	importMap = make(map[string]struct{})
}

// go数据类型转dart数据类型
func goTypeToKotlinType(varType string) string {
	switch varType {
	case "int":
		return "Int"
	case "int8":
		return "Byte"
	case "int16":
		return "Short"
	case "int32":
		return "Int"
	case "int64":
		return "Long"
	case "float32":
		return "Float"
	case "float64":
		return "Double"
	case "string":
		return "String"
	case "bool":
		return "Boolean"
	case "time.Time":
		return "String"
	case "any":
		return "Any"
	default:
		if strings.HasPrefix(varType, "[]") {
			listFormName := varType[2:]
			//if strings.HasSuffix(listFormName, "Form") { //如果是以Form结尾的类名
			//	listFormName = listFormName[:len(listFormName)-4] + "Model"
			//	importMap[listFormName] = struct{}{}
			//} else {
			//	listFormName = goTypeToKotlinType(listFormName)
			//}
			listFormName = goTypeToKotlinType(listFormName)
			return "Array<" + listFormName + ">"
		} else if strings.HasSuffix(varType, "Form") {
			//kotlinType := ""
			//importMap[kotlinType] = struct{}{}
			if strings.Contains(varType, ".") {
				varType = varType[strings.LastIndex(varType, ".")+1:]
			}
			return varType[:len(varType)-4] + "Model"
		} else {
			return varType
		}
	}
}

// 生成变量定义部分的代码
func makeVarBodySource(goStruct GoBean.GoStruct) string {
	source := ""
	for _, goVariable := range goStruct.Members {
		kotlinType := goTypeToKotlinType(goVariable.Type)
		source += "\n  " + goVariable.Comment + "  var " + goVariable.LowerName() + ": " + kotlinType + " = " + getKotlinDefaultValue(kotlinType) + "\n"
	}
	return source + "\n"
}

// go数据类型转dart数据类型
func getKotlinDefaultValue(varType string) string {
	switch varType {
	case "Int":
		return "0"
	case "Byte":
		return "0"
	case "Short":
		return "0"
	case "Long":
		return "0"
	case "Float":
		return "0"
	case "Double":
		return "0.0"
	case "String":
		return "\"\""
	case "Boolean":
		return "false"
	case "Any":
		return "0"
	default:
		if strings.HasPrefix(varType, "Array<") {
			return "emptyArray()"
		} else if strings.HasSuffix(varType, "Model") {
			return varType + "()"
		} else {
			return ""
		}
	}
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
	path := Application.Args.TargetDir + "/model/" + fileName + ".kt"
	fileContent := FileUtil.ReadText(path)
	if fileContent == source {
		return
	}
	FileUtil.WriteText(path, source)
}
