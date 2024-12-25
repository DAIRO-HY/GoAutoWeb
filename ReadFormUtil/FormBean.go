package ReadFormUtil

import "strings"

// 路由参数信息
type FormBean struct {

	//FORM包所在路径
	PackagePath string

	//属性列表
	Properties []PropertyBean

	//结构体名
	Name string
}

func (mine *FormBean) GetNickImport() string {
	if len(mine.PackagePath) == 0 {
		return ""
	}
	nick := strings.ReplaceAll(mine.PackagePath, "/", "")
	nick = strings.ReplaceAll(nick, "_", "")
	return nick
}
