package Bean

type GoBean struct {

	//文件路径
	FilePath string

	//包名
	Package string

	//导入的go模块
	Imports []string

	//结构体列表
	Structs []StructBean

	//函数列表
	Methods []MethodBean

	//注解列表
	AnnotationMap map[string]AnnotationBean
}
