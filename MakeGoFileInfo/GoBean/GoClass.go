package GoBean

type GoClass struct {

	//文件路径
	FilePath string

	//包名
	Package string

	//导入的go模块
	Imports []string

	//结构体列表
	Structs []GoStruct

	//函数列表
	Methods []GoMethod

	//注解列表
	AnnotationMap map[string]GoAnnotation
}
