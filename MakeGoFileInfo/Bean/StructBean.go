package Bean

type StructBean struct {

	//结构体名
	Name string

	//注释
	Comment string

	//注解
	AnnotationMap map[string]AnnotationBean

	//成员变量
	Members []VariableBean
}
