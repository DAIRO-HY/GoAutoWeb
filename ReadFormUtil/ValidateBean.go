package ReadFormUtil

import (
	"fmt"
	"strings"
)

// 表单验证Bean
type ValidateBean struct {

	//验证规则名称
	Name string

	//参数
	Args map[string]string
}

// 生成表单验证部分代码
func (mine *ValidateBean) MakeValidSource(field string) string {
	msg := mine.Args["msg"]
	formField := strings.ToLower(field[:1]) + field[1:]
	source := ""
	if mine.Name == "NOTEMPTY" { //非空验证
		source += fmt.Sprintf("\t\tisNotEmpty(filedError, \"%s\", valid%s, \"%s\") // 非空验证\n", formField, field, msg)
	} else if mine.Name == "LENGTH" { //长度验证代码
		minlength := mine.Args["min"]
		maxlength := mine.Args["max"]
		if minlength == "" {
			minlength = "nil"
		} else {
			minlength = "intP(" + minlength + ")"
		}
		if maxlength == "" {
			maxlength = "nil"
		} else {
			maxlength = "intP(" + maxlength + ")"
		}
		source += fmt.Sprintf("\t\tisLength(filedError, \"%s\", valid%s, %s, %s, \"%s\")// 输入长度验证\n", formField, field, minlength, maxlength, msg)
	} else if mine.Name == "LIMIT" { //数值值区间验证
		minValue := mine.Args["min"]
		maxValue := mine.Args["max"]
		if minValue == "" {
			minValue = "nil"
		} else {
			minValue = "floatP(" + minValue + ")"
		}
		if maxValue == "" {
			maxValue = "nil"
		} else {
			maxValue = "floatP(" + maxValue + ")"
		}
		source += fmt.Sprintf("\t\tisLimit(filedError, \"%s\", valid%s, %s, %s, \"%s\")// 数值值区间验证\n", formField, field, minValue, maxValue, msg)
	} else if mine.Name == "EMAIL" {
		source += fmt.Sprintf("\t\tisEmail(filedError, \"%s\", valid%s, \"%s\") // 邮箱格式验证\n", formField, field, msg)
	}
	return source
}
