package GoBean

// go函数信息
type GoMethod struct {

	//函数名
	Name string

	//返回值类型
	Returns []string

	//注释
	Comment string

	//函数参数列表
	Parameters []GoVariable

	//注解列表
	AnnotationMap map[string]GoAnnotation
}
