package bean

import "strings"

// 拦截器信息
type InterceptorBean struct {

	//包所在路径
	PackagePath string

	//包含路由
	Include []string

	//排除路由
	Exclude []string

	//函数名
	FuncName string

	//拦截器执行时机 pre:进入controller之前 after:controller执行完成后
	HandleFlag string

	//执行优先顺序,值越小越优先
	Order int
}

func (mine *InterceptorBean) GetNickImport() string {
	if len(mine.PackagePath) == 0 {
		return ""
	}
	nick := strings.ReplaceAll(mine.PackagePath, "/", "")
	nick = strings.ReplaceAll(nick, "_", "")
	return nick
}
