package ReadPathUtil

import (
	"GoAutoWeb/ReadFormUtil"
	"fmt"
	"strings"
)

// 路由参数信息
type ParamBean struct {

	//FORM包所在路径
	PackagePath string

	//参数类型
	VarType string

	//参数名
	Name string
}

func (mine *ParamBean) GetNickImport() string {
	if len(mine.PackagePath) == 0 {
		return ""
	}
	nick := strings.ReplaceAll(mine.PackagePath, "/", "")
	nick = strings.ReplaceAll(nick, "_", "")
	return nick
}

// 生成获取参数的代码
func (mine *ParamBean) MakeGetParameterSource() string {
	if mine.VarType == "http.ResponseWriter" { //这不是一个URL参数
		return ""
	}
	if mine.VarType == "*http.Request" { //这不是一个URL参数
		return ""
	}
	source := ""

	if strings.HasSuffix(mine.VarType, "Form") { //这是一个结构体Form表单

		//生成表单相关验证代码
		formBean := ReadFormUtil.FormMap[mine.PackagePath+"/"+mine.VarType]

		//source += "\n\t\t" + mine.Name + " := getForm[" + mine.GetNickImport() + "." + mine.VarType + "](paramMap)"
		source += formBean.MakeValidateSource()
		source += formBean.MakeGetParameterSource(mine.GetNickImport()+"."+mine.VarType, mine.Name)
		//source += "\t\tvalidBody := validateForm(" + mine.Name + ")\n"
		//source += "\t\tif validBody != nil {\n"
		//source += "\t\t\twriteFieldError(writer, validBody)\n"
		//source += "\t\t\treturn\n"
		//source += "\t\t}\n"
		for _, function := range formBean.Functions {
			source += function.MakeFormCheckSource(mine.Name)
		}
	} else {

		//初始化变量
		source += fmt.Sprintf("\t\tvar %s %s // 初始化变量\n", mine.Name, mine.VarType)
		var callMethodName string
		if strings.HasSuffix(mine.VarType, "int") { //int类型的变量
			callMethodName = "getIntArray"
		} else if strings.HasSuffix(mine.VarType, "int8") { //int类型的变量
			callMethodName = "getInt8Array"
		} else if strings.HasSuffix(mine.VarType, "int16") { //int16类型的变量
			callMethodName = "getInt16Array"
		} else if strings.HasSuffix(mine.VarType, "int32") { //int32类型的变量
			callMethodName = "getInt32Array"
		} else if strings.HasSuffix(mine.VarType, "int64") { //int64类型的变量
			callMethodName = "getInt64Array"
		} else if strings.HasSuffix(mine.VarType, "float32") { //float32类型的变量
			callMethodName = "getFloat32Array"
		} else if strings.HasSuffix(mine.VarType, "float64") { //float64类型的变量
			callMethodName = "getFloat64Array"
		} else if strings.HasSuffix(mine.VarType, "bool") { //bool类型的变量
			callMethodName = "getBoolArray"
		} else { //字符串类型的变量
			callMethodName = "getStringArray"
		}
		source += fmt.Sprintf("\t\t%sArr := %s(query, postForm, \"%s\")\n", mine.Name, callMethodName, mine.Name)
		source += fmt.Sprintf("\t\tif %sArr != nil { // 如果参数存在\n", mine.Name)
		if strings.HasPrefix(mine.VarType, "*") { //这是一个指针类型
			source += fmt.Sprintf("\t\t\t%s = &%sArr[0]\n", mine.Name, mine.Name)
		} else if strings.HasPrefix(mine.VarType, "[]") { //这是一个数组类型
			source += fmt.Sprintf("\t\t\t%s = %sArr\n", mine.Name, mine.Name)
		} else {
			source += fmt.Sprintf("\t\t\t%s = %sArr[0]\n", mine.Name, mine.Name)
		}
		source += "\t\t}\n"
	}
	return source
}
