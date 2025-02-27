package MakeModel

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/ReadFormUtil"
	"encoding/json"
	"os"
	"strings"
)

// 导入的Model列表
var importList []string

// 生成Model类文件
func Make() {
	jsonData, _ := json.Marshal(ReadFormUtil.FormList)

	// 复制一个对象操作，避免指针操作修改到原始数据
	copyList := make([]ReadFormUtil.FormBean, 0)
	json.Unmarshal(jsonData, &copyList)
	for _, it := range copyList {
		makeModelByForm(it)
	}
}

// 通过一个Form表单生成Model文件
func makeModelByForm(form ReadFormUtil.FormBean) {

	//修改类名
	form.Name = form.Name[0:len(form.Name)-4] + "Model"
	body := makeVarBodySource(form) + makeConstructorSource(form) + makeToJsonSource(form) + makeFromSource(form)

	importStr := ""
	for _, it := range importList {
		importStr += "import '" + it + ".dart';\n"
	}
	if len(importStr) > 0 {
		importStr = importStr[:len(importStr)-1]
	}
	source := `/*工具自动生成代码,请勿手动修改*/

import 'dart:convert';

import '../../util/JsonSerialize.dart';
` + importStr + `

class ` + form.Name + ` extends JsonSerialize {
` + body + "}\n"
	save(source, form.Name)
	importList = make([]string, 0)
}

// go数据类型转dart数据类型
func goTypeToDartType(varType string) string {
	dartType := ""
	switch varType {
	case "int", "int8", "int16", "int32", "int64":
		dartType = "int"
	case "string":
		dartType = "String"
	default:
		if strings.HasPrefix(varType, "[]") {
			listFormName := varType[2:]
			if strings.HasSuffix(listFormName, "Form") { //如果是以Form结尾的类名
				listFormName = listFormName[:len(listFormName)-4] + "Model"
				importList = append(importList, listFormName)
			} else {
				listFormName = goTypeToDartType(listFormName)
			}
			dartType = "List<" + listFormName + ">"
		} else if strings.HasSuffix(varType, "Form") {
			dartType = varType[:len(varType)-4] + "Model"
			importList = append(importList, dartType)
		} else {
			dartType = varType
		}
	}
	return dartType
}

// 生成变量定义部分的代码
func makeVarBodySource(form ReadFormUtil.FormBean) string {
	source := ""
	for _, it := range form.Properties {
		source += "\n  " + it.Comment + "  " + goTypeToDartType(it.VarType) + " " + it.LowerName() + ";\n"
	}
	return source + "\n"
}

// 生成构造函数部分的代码
func makeConstructorSource(form ReadFormUtil.FormBean) string {
	source := ""
	for _, it := range form.Properties {
		source += "required this." + it.LowerName() + ", "
	}
	source = source[:len(source)-2]
	return "  " + form.Name + "({" + source + "});\n"
}

// 生成转Json部分的代码
func makeToJsonSource(form ReadFormUtil.FormBean) string {
	source := ""
	for _, it := range form.Properties {
		source += "        \"" + it.LowerName() + "\": this." + it.LowerName() + ",\n"
	}
	source = source[0 : len(source)-1]
	return `
  /// 将model转Json
  @override
  toJson() => {
` + source + `
      };
`
}

// 生成转换成model部分的代码
func makeFromSource(form ReadFormUtil.FormBean) string {
	source := ""
	for _, it := range form.Properties {
		if strings.HasPrefix(it.VarType, "[]") { //如果这是一个List数据类型
			if strings.HasSuffix(it.ListType(), "Form") {
				source += "        " + it.LowerName() + ": " + goTypeToDartType(it.ListType()) + ".fromMapList(map[\"" + it.LowerName() + "\"]),\n"
			} else {
				//TODO：待实现
			}
		} else {
			source += "        " + it.LowerName() + ": map[\"" + it.LowerName() + "\"],\n"
		}
	}
	source = source[:len(source)-2]
	source = `
  /// 将json字符串转{Model}对象
  static {Model} fromJson(String json) {
    Map<String, dynamic> map = jsonDecode(json);
    return {Model}.fromMap(map);
  }

  /// 将Map对象转{Model}对象
  static {Model} fromMap(Map<String, dynamic> map) {
    return {Model}(
` + source + `);
  }

  /// 将Json字符串转{Model}对象列表
  static List<{Model}> fromJsonList(String json) {
    List<dynamic> list = jsonDecode(json);
    return {Model}.fromMapList(list);
  }

  /// 将List<Map>对象转{Model}对象列表
  static List<{Model}> fromMapList(List<dynamic> list) {
    return list.map((map) => {Model}.fromMap(map)).toList();
  }
`
	return strings.ReplaceAll(source, "{Model}", form.Name)
}

// 保存文件
func save(source string, fileName string) {
	os.WriteFile(Application.Args.TargetDir+"/lib/api/model/"+fileName+".dart", []byte(source), 0644)
}
