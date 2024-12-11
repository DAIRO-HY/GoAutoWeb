package ReadPathUtil

import "strings"

// 路由参数信息
type ParamBean struct {

	//FORM包所在路径
	PackagePath string

	//参数类型
	VarType string

	//参数名
	Name string
}

func (mine *ParamBean) GetNickImport() string {
	if len(mine.PackagePath) == 0 {
		return ""
	}
	nick := strings.ReplaceAll(mine.PackagePath, "/", "")
	nick = strings.ReplaceAll(nick, "_", "")
	return nick
}
