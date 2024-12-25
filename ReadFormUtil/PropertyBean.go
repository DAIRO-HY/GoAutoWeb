package ReadFormUtil

// 结构体属性Bean
type PropertyBean struct {

	//FORM包所在路径
	//PackagePath string

	//参数类型
	VarType string

	//参数名
	Name string

	//表单验证列表
	Valids []string
}

//func (mine *PropertyBean) GetNickImport() string {
//	if len(mine.PackagePath) == 0 {
//		return ""
//	}
//	nick := strings.ReplaceAll(mine.PackagePath, "/", "")
//	nick = strings.ReplaceAll(nick, "_", "")
//	return nick
//}
