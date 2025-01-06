package ReadFormUtil

import "fmt"

// 表单验证Bean
type ValidateBean struct {

	//验证规则名称
	Name string

	//参数
	Args map[string]string
}

// 生成表单验证部分代码
func (mine *ValidateBean) MakeValidSource(field string, value string) string {
	source := ""
	if mine.Name == "NOTEMPTY" { //非空验证
		msg := mine.Args["msg"]
		source += fmt.Sprintf("\t\tisNotEmpty(filedError, \"%s\", %s, \"%s\")// 非空验证\n", field, value, msg)
	} else if mine.Name == "LENGTH" {
		minlen := mine.Args["min"]
		maxlen := mine.Args["max"]
		msg := mine.Args["msg"]
		source += fmt.Sprintf("\t\tisLength(filedError, \"%s\", %s, %s, %s, \"%s\")// 输入长度验证\n", field, value, minlen, maxlen, msg)
	}
	return source
}
