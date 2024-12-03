package bean

// Controller路由信息
type PathBean struct {

	//包所在路径
	PackagePath string

	//请求方案
	Method string

	//路由路径
	Path string

	//函数名
	FuncName string

	//函数名
	ReturnType string

	//该路由的参数
	Parameters []ParamBean
}
