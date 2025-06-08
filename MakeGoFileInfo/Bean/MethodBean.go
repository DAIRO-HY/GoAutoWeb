package Bean

// go函数信息
type MethodBean struct {

	//函数名
	Name string

	//返回值类型
	Returns []string

	//注释
	Comment string

	//函数参数列表
	Parameters []VariableBean

	//注解列表
	Annotations []AnnotationBean
}
