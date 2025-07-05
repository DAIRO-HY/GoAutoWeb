package GoBean

type GoStruct struct {

	//结构体名
	Name string

	//注释
	Comment string

	//注解
	AnnotationMap map[string]GoAnnotation

	//成员变量
	Members []GoVariable
}
