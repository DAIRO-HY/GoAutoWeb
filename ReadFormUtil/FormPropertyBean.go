package ReadFormUtil

// 结构体属性Bean
type FormPropertyBean struct {

	//参数类型
	VarType string

	//参数名
	Name string

	/** 表单验证列表 **/
	valids []FormValidateBean
}
