package MakeApiRequest

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/ReadFormUtil"
	"GoAutoWeb/ReadPathUtil"
	"encoding/json"
	"os"
	"strings"
)

// 导入的Model列表
var importList = map[string]struct{}{}

// 生成API常量文件
func Make() {
	copyList := fixPathList()
	tempClassName := ""
	sourceBody := ""
	for _, it := range copyList {
		if !strings.HasSuffix(it.FileName, Application.Args.ApiSuffix+".go") {
			continue
		}
		className := strings.ReplaceAll(it.FileName, Application.Args.ApiSuffix+".go", "Api")
		if tempClassName != "" && tempClassName != className { //上一个文件的代码生成完成，先保存
			importStr := ""
			for im := range importList {
				importStr += im + "\n"
			}

			source := ""
			source += "import 'API.dart';\n"
			source += importStr
			source += "\nclass " + tempClassName + " {\n" + sourceBody + "}"
			save(tempClassName, source)
			sourceBody = ""
			importList = map[string]struct{}{}
		}
		tempClassName = className

		//参数调用部分代码
		callParamSource := makeParamSource(it)
		callHttpSource := makeCallHttpSource(it)
		returnSource := makeReturnTypeSource(it)
		comment := makeComment(it)
		sourceBody += comment
		sourceBody += "\n  static " + returnSource + " " + it.LowerFuncName() + "(" + callParamSource + "){\n" + callHttpSource + "\n  }\n"
	}
}

// 修正一些数据
func fixPathList() []ReadPathUtil.PathBean {
	jsonData, _ := json.Marshal(ReadPathUtil.PathList)

	// 复制一个对象操作，避免指针操作修改到原始数据
	copyList := make([]ReadPathUtil.PathBean, 0)
	json.Unmarshal(jsonData, &copyList)

	for i, _ := range copyList {
		pb := &copyList[i]
		newParameters := make([]ReadPathUtil.ParamBean, 0)
		for _, it := range pb.Parameters {
			if it.VarType == "http.ResponseWriter" || it.VarType == "*http.Request" {
				continue
			}
			if strings.HasPrefix(it.Name, "_") { //以_开头的参数不需要自动生成
				continue
			}
			newParameters = append(newParameters, it)
		}
		pb.Parameters = newParameters
	}
	return copyList
}

// 生成参数部分代码
func makeParamSource(pb ReadPathUtil.PathBean) string {
	source := ""
	for _, it := range pb.Parameters {
		if strings.HasSuffix(it.VarType, "Form") {
			form, isExists := ReadFormUtil.FormMap[it.PackagePath+"/"+it.VarType]
			if !isExists {
				continue
			}
			for _, formMember := range form.Properties {
				paramType := goTypeToDartType(formMember.VarType)
				source += "required " + paramType + " " + formMember.LowerName() + ","
			}
		} else {
			paramType := goTypeToDartType(it.VarType)
			source += "required " + paramType + " " + it.Name + ","
		}
	}
	if source != "" {
		source = "{" + source[:len(source)-1] + "}"
	}
	return source
}

// 生成发起网络请求部分代码
func makeCallHttpSource(pb ReadPathUtil.PathBean) string {
	source := ""
	for _, it := range pb.Parameters {
		if strings.HasSuffix(it.VarType, "Form") {
			form, isExists := ReadFormUtil.FormMap[it.PackagePath+"/"+it.VarType]
			if !isExists {
				continue
			}
			for _, formMember := range form.Properties {
				source += ".add(\"" + formMember.LowerName() + "\"," + formMember.LowerName() + ")"
			}
		} else {
			source += ".add(\"" + it.Name + "\"," + it.Name + ")"
		}
	}
	constName := urlToConst(pb)
	returnSource := makeReturnTypeSource(pb)
	toModelSource := makeToModelSource(pb) //转换Model对象的代码
	source = "    return " + returnSource + "(Api." + constName + toModelSource + ")" + source + ";"
	return source
}

// 生成返回值类型的代码
func makeReturnTypeSource(pb ReadPathUtil.PathBean) string {
	returnType := goTypeToDartType(pb.ReturnType)
	if returnType == "" {
		importList["import '../util/http/VoidApiHttp.dart';"] = struct{}{}
		return "VoidApiHttp"
	} else {
		importList["import '../util/http/ReturnApiHttp.dart';"] = struct{}{}
		return "ReturnApiHttp<" + returnType + ">"
	}
}

// 生成转换成model对象的代码
func makeToModelSource(pb ReadPathUtil.PathBean) string {
	returnType := goTypeToDartType(pb.ReturnType)
	if returnType == "" {
		return ""
	}
	if strings.HasSuffix(returnType, "Model") {
		return ", " + returnType + ".fromJson"
	} else if strings.HasPrefix(returnType, "List<") {
		tType := returnType[strings.Index(returnType, "<")+1 : strings.Index(returnType, ">")]
		if strings.HasSuffix(tType, "Model") {
			return ", " + tType + ".fromJsonList"
		} else {
			return ""
		}
	} else {
		return ""
	}
}

// 生成注释部分的代码
func makeComment(pb ReadPathUtil.PathBean) string {
	comment := pb.Comment
	if comment == "" {
		return ""
	}
	cms := strings.Split(comment, "\n")
	return "  //" + strings.Join(cms, "\n  //")
}

// go数据类型转dart数据类型
func goTypeToDartType(varType string) string {
	dartType := ""
	switch varType {
	case "int", "int8", "int16", "int32", "int64":
		dartType = "int"
	case "string":
		dartType = "String"
	case "any":
		dartType = "Object"
	case "error":
		dartType = ""
	default:
		listFormName := varType
		if strings.HasPrefix(varType, "[]") {
			listFormName = listFormName[2:]
			listFormName = goTypeToDartType(listFormName)
			dartType = "List<" + listFormName + ">"
		} else if strings.HasPrefix(varType, "map") { //这是一个map数据
			keyType := varType[strings.Index(varType, "[")+1 : strings.Index(varType, "]")]
			valueType := varType[strings.Index(varType, "]")+1:]
			dartType = "Map<" + goTypeToDartType(keyType) + "," + goTypeToDartType(valueType) + ">"
		} else if strings.HasSuffix(varType, "Form") { //如果是以Form结尾的类名
			if strings.Contains(listFormName, ".") { //只取点以后的字符串
				listFormName = listFormName[strings.LastIndex(listFormName, ".")+1:]
			}
			listFormName = listFormName[:len(listFormName)-4] + "Model"
			dartType = listFormName
			importList["import 'model/"+listFormName+".dart';"] = struct{}{}
		} else {
			dartType = listFormName
		}
	}
	return dartType
}

// 将路由转成常量名
func urlToConst(pb ReadPathUtil.PathBean) string {
	url := pb.Path + pb.VariablePath
	key := strings.ReplaceAll(url, "/", "_")
	key = strings.ReplaceAll(key, "{", "_")
	key = strings.ReplaceAll(key, "}", "_")
	key = strings.ReplaceAll(key, "__", "_")
	key = strings.ReplaceAll(key, "__", "_")
	key = strings.ReplaceAll(key, ".", "_")
	key = strings.ToUpper(key)
	return key[1:]
}

// 保存文件
func save(fileName string, source string) {
	os.WriteFile(Application.Args.TargetDir+"/lib/api/"+fileName+".dart", []byte(source), 0644)
}
