package ReadFormUtil

import "strings"

// 结构体属性Bean
type FormPropertyBean struct {

	//注释
	Comment string

	//参数类型
	VarType string

	//参数名
	Name string

	/** 表单验证列表 **/
	valids []FormValidateBean
}

// 获取小写的名称
func (mine FormPropertyBean) LowerName() string {
	return strings.ToLower(mine.Name[:1]) + mine.Name[1:]
}
