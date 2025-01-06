package ReadFormUtil

import "strings"

// 结构体函数Bean
type FunctionBean struct {

	//函数名
	Name string

	//返回值类型
	ReturnType string
}

// 生成相关验证代码
func (mine *FunctionBean) MakeFormCheckSource(paramName string) string {
	if !strings.HasPrefix(mine.Name, "Is") { //如果函数名不是以Is开头
		return ""
	}

	//获取要验证的参数名
	fieldStr := mine.Name[2:] //去掉前面的Is
	fieldArr := strings.Split(fieldStr, "And")
	fields := ""
	for _, it := range fieldArr {

		//将首字母小写之后作为要验证的字段名
		fields += " \"" + strings.ToLower(it[:1]) + it[1:] + "\"" + ","
	}
	fields = fields[:len(fields)-1]

	source := ""
	msgVar := paramName + mine.Name + "Msg"
	source += "\t\t" + msgVar + " := " + paramName + "." + mine.Name + "()\n"
	source += "\t\tif " + msgVar + " != nil { // 表单相关验证失败\n"
	source += "\t\t\twriteFieldFormError(writer, *" + msgVar + "," + fields + ")\n"
	source += "\t\t\treturn\n"
	source += "\t\t}\n"
	return source
}
