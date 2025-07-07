package MakeApiRequest

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/FileUtil"
	"GoAutoWeb/ReadFormUtil"
	"GoAutoWeb/ReadPathUtil"
	"encoding/json"
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
			importStr := "import " + Application.Args.TargetPackage + ".http.ApiHttp\n"
			for im := range importList {
				importStr += "import " + im + "\n"
			}

			source := "package " + Application.Args.TargetPackage + "\n\n"
			source += importStr
			source += "\nobject " + tempClassName + " {\n" + sourceBody + "}"
			save(tempClassName, source)
			sourceBody = ""
			importList = map[string]struct{}{}
		}
		tempClassName = className

		//参数调用部分代码
		callParamSource := makeParamSource(it)
		callApiHttp := makeCallHttpSource(it)
		returnType := makeReturnTypeSource(it)
		comment := makeComment(it)
		sourceBody += comment
		sourceBody += "\n  fun " + makeFuncName(it.LowerFuncName()) + "(" + callParamSource + ") : " + returnType + "{\n" + callApiHttp + "\n  }\n"
	}
}

// 生成函数名
func makeFuncName(name string) string {
	switch name {
	case "init":
		return "_init"
	default:
		return name
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
				paramType := goTypeToKotlinType(formMember.VarType)
				source += formMember.LowerName() + ": " + paramType + ","
			}
		} else {
			paramType := goTypeToKotlinType(it.VarType)
			source += it.Name + ": " + paramType + ","
		}
	}
	if source != "" {
		source = source[:len(source)-1]
	}
	return source
}

// 生成发起网络请求部分代码
func makeCallHttpSource(pb ReadPathUtil.PathBean) string {
	parameterSource := ""
	for _, it := range pb.Parameters {
		if strings.HasSuffix(it.VarType, "Form") { //如果这是个Form表单,则挨个解析表单内全部标量
			form, isExists := ReadFormUtil.FormMap[it.PackagePath+"/"+it.VarType]
			if !isExists {
				continue
			}
			for _, formMember := range form.Properties {
				parameterSource += "\"" + formMember.LowerName() + "\" to " + formMember.LowerName() + ","
			}
		} else {
			parameterSource += "\"" + it.Name + "\" to " + it.Name + ","
		}
	}
	if len(parameterSource) > 0 {
		parameterSource = parameterSource[:len(parameterSource)-1]
		parameterSource = ",parameter = mapOf(" + parameterSource + ")"
	}
	returnCls := ""
	if len(pb.ReturnType) > 0 { //如果有返回值
		returnType := goTypeToKotlinType(pb.ReturnType)
		if strings.HasPrefix(returnType, "List<") { //这是一个列表返回值
			returnType = returnType[5 : len(returnType)-1]
			returnCls = ", listCls = " + returnType + "::class.java"
		} else {
			returnCls = ", cls = " + returnType + "::class.java"
		}
	}

	constName := urlToConst(pb)
	return "    return ApiHttp(ApiConst." + constName + parameterSource + returnCls + ")"
}

// 生成返回值类型的代码
func makeReturnTypeSource(pb ReadPathUtil.PathBean) string {
	returnType := goTypeToKotlinType(pb.ReturnType)
	if len(returnType) == 0 {
		returnType = "Unit"
	}
	return "ApiHttp<" + returnType + ">"
}

// 生成注释部分的代码
func makeComment(pb ReadPathUtil.PathBean) string {
	comment := pb.Comment
	if comment == "" {
		return ""
	}
	cms := strings.Split(comment, "\n")
	return "\n  //" + strings.Join(cms, "\n  //")
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
		kotlinType := ""
		if strings.HasPrefix(varType, "[]") { //这是一个List数据类型
			listType := varType[2:]
			listType = goTypeToKotlinType(listType)
			kotlinType = "List<" + listType + ">"
		} else if strings.HasPrefix(varType, "map[") { //这是一个Map数据类型
			keyType := varType[strings.Index(varType, "[")+1 : strings.Index(varType, "]")]
			valueType := varType[strings.Index(varType, "]")+1:]

			keyType = goTypeToKotlinType(keyType)
			valueType = goTypeToKotlinType(valueType)
			kotlinType = "[" + keyType + " : " + valueType + "]"
		} else if strings.Contains(varType, ".") { //这个返回的类型包含了包名
			kotlinType = varType[strings.LastIndex(varType, ".")+1:]
			kotlinType = goTypeToKotlinType(kotlinType)
		} else if strings.HasSuffix(varType, "Form") {
			kotlinType = varType[:len(varType)-4] + "Model"
			importList[Application.Args.TargetPackage+".model."+kotlinType] = struct{}{}
		} else if varType == "" {
			return ""
		} else {
			importList[Application.Args.TargetPackage+".model."+varType] = struct{}{}
			return varType
		}
		return kotlinType
	}
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
	path := Application.Args.TargetDir + "/" + fileName + ".kt"
	fileContent := FileUtil.ReadText(path)
	if fileContent == source {
		return
	}
	FileUtil.WriteText(path, source)
}
