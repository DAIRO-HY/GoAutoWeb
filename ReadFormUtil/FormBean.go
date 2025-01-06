package ReadFormUtil

import (
	"fmt"
	"strings"
)

// 路由参数信息
type FormBean struct {

	//FORM包所在路径
	PackagePath string

	//属性列表
	Properties []PropertyBean

	//结构体函数列表
	Functions []FunctionBean

	//结构体名
	Name string
}

func (mine *FormBean) GetNickImport() string {
	if len(mine.PackagePath) == 0 {
		return ""
	}
	nick := strings.ReplaceAll(mine.PackagePath, "/", "")
	nick = strings.ReplaceAll(nick, "_", "")
	return nick
}

// 生成获取表单值的代码
func (mine *FormBean) MakeGetParameterSource(formStructName string, paramName string) string {
	source := ""
	source += fmt.Sprintf("\t\t%s:=%s{}\n", paramName, formStructName)
	for _, it := range mine.Properties {
		key := strings.ToLower(it.Name[:1]) + it.Name[1:]

		var callMethodName string
		if strings.HasSuffix(it.VarType, "int") { //int类型的变量
			callMethodName = "getIntArray"
		} else if strings.HasSuffix(it.VarType, "int8") { //int类型的变量
			callMethodName = "getInt8Array"
		} else if strings.HasSuffix(it.VarType, "int16") { //int16类型的变量
			callMethodName = "getInt16Array"
		} else if strings.HasSuffix(it.VarType, "int32") { //int32类型的变量
			callMethodName = "getInt32Array"
		} else if strings.HasSuffix(it.VarType, "int64") { //int64类型的变量
			callMethodName = "getInt64Array"
		} else if strings.HasSuffix(it.VarType, "float32") { //float32类型的变量
			callMethodName = "getFloat32Array"
		} else if strings.HasSuffix(it.VarType, "float64") { //float64类型的变量
			callMethodName = "getFloat64Array"
		} else if strings.HasSuffix(it.VarType, "bool") { //bool类型的变量
			callMethodName = "getBoolArray"
		} else { //字符串类型的变量
			callMethodName = "getStringArray"
		}
		source += fmt.Sprintf("\t\t%s := %s(query,postForm,\"%s\")\n", paramName+it.Name, callMethodName, key)
		source += fmt.Sprintf("\t\tif %s != nil {// 如果参数存在\n", paramName+it.Name)
		if strings.HasPrefix(it.VarType, "*") { //这是一个指针类型
			source += fmt.Sprintf("\t\t\t%s.%s = &%s[0]\n", paramName, it.Name, paramName+it.Name)
		} else if strings.HasPrefix(it.VarType, "[]") { //这是一个切片
			source += fmt.Sprintf("\t\t\t%s.%s = %s\n", paramName, it.Name, paramName+it.Name)
		} else {
			source += fmt.Sprintf("\t\t\t%s.%s = %s[0]\n", paramName, it.Name, paramName+it.Name)
		}
		source += "\t\t}\n\n"
	}
	return source
}

// 生成获取表单值的代码
func (mine *FormBean) MakeValidateSource(formName string) string {
	source := ""
	for _, property := range mine.Properties {
		for _, validBean := range property.valids {
			source += validBean.MakeValidSource(property.Name, formName+"."+property.Name)
		}
	}
	if source != "" {
		source = "\t\tfiledError := map[string]*[]string{}\n" + source
		source += "\t\tif len(filedError) > 0{\n"
		source += "\t\t\twriteFieldError(writer, filedError)\n\t\t\treturn\n"
		source += "\t\t}\n"
	}
	return source
}
