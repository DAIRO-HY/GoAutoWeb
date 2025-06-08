package Bean

// go变量信息
type VariableBean struct {

	//变量名
	Name string

	//变量类型
	Type string

	//注释
	Comment string

	//初始值
	Value string

	//注解
	Annotations []AnnotationBean
}
