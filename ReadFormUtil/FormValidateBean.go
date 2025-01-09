package ReadFormUtil

import (
	"fmt"
	"strings"
)

// 表单验证Bean
type FormValidateBean struct {

	//验证规则名称
	Name string

	//参数
	Args map[string]string
}

// 生成表单验证部分代码
func (mine FormValidateBean) MakeValidSource(field string) string {
	formField := strings.ToLower(field[:1]) + field[1:]
	source := ""
	if mine.Name == "NOTEMPTY" { //非空验证
		source += fmt.Sprintf("\t\tisNotEmpty(filedError, \"%s\", valid%s) // 非空验证\n", formField, field)
	} else if mine.Name == "NOTBLANK" { //非空白验证
		source += fmt.Sprintf("\t\tisNotBlank(filedError, \"%s\", valid%s) // 非空白验证\n", formField, field)
	} else if mine.Name == "LENGTH" { //长度验证代码
		minlength := mine.Args["min"]
		maxlength := mine.Args["max"]
		if minlength == "" {
			minlength = "-1"
		}
		if maxlength == "" {
			maxlength = "-1"
		}
		source += fmt.Sprintf("\t\tisLength(filedError, \"%s\", valid%s, %s, %s)// 输入长度验证\n", formField, field, minlength, maxlength)
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
		source += fmt.Sprintf("\t\tisLimit(filedError, \"%s\", valid%s, %s, %s)// 数值值区间验证\n", formField, field, minValue, maxValue)
	} else if mine.Name == "DIGITS" { //数值验证
		integer := mine.Args["integer"]
		fraction := mine.Args["fraction"]
		if integer == "" {
			integer = "0"
		}
		if fraction == "" {
			fraction = "0"
		}
		source += fmt.Sprintf("\t\tisDigits(filedError, \"%s\", valid%s, %s, %s)// 数值值区间验证\n", formField, field, integer, fraction)
	} else if mine.Name == "HALF" { //半角验证

		upper := mine.Args["upper"]
		lower := mine.Args["lower"]
		number := mine.Args["number"]
		symbol := mine.Args["symbol"]
		if upper != "false" {
			upper = "true"
		}
		if lower != "false" {
			lower = "true"
		}
		if number != "false" {
			number = "true"
		}
		if symbol != "false" {
			symbol = "true"
		}
		source += fmt.Sprintf("\t\tisHalf(filedError, \"%s\", valid%s, %s, %s, %s, %s)// 半角字符验证\n", formField, field, upper, lower, number, symbol)
	} else if mine.Name == "EMAIL" {
		source += fmt.Sprintf("\t\tisEmail(filedError, \"%s\", valid%s) // 邮箱格式验证\n", formField, field)
	}
	return source
}
