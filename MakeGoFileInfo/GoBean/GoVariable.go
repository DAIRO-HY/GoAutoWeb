package GoBean

import "strings"

// go变量信息
type GoVariable struct {

	//变量名
	Name string

	//变量类型
	Type string

	//注释
	Comment string

	//初始值
	Value string

	//注解
	AnnotationMap map[string]GoAnnotation
}

// 获取小写的名称
func (mine GoVariable) LowerName() string {
	return strings.ToLower(mine.Name[:1]) + mine.Name[1:]
}
