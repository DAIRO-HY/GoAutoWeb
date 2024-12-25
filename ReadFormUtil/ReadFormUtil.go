package ReadFormUtil

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/FileUtil"
	"regexp"
	"strings"
)

// 表单列表
var FormList []FormBean

func Make() {
	for _, fPath := range Application.GoFileList {
		if !strings.HasSuffix(fPath, "Form.go") { //过滤文件后缀
			continue
		}

		// 读取go文件里的路由配置
		formBean := readControllerPath(fPath)
		FormList = append(FormList, formBean)
	}
}

func readControllerPath(path string) FormBean {

	//包所在路径
	packagePath := path[len(Application.RootProject):]
	packagePath = strings.ReplaceAll(packagePath, "\\", "/")
	packagePath = packagePath[:strings.LastIndex(packagePath, "/")]

	//读取代码
	goCode := FileUtil.ReadText(path)

	//先统一换行符
	goCode = strings.ReplaceAll(goCode, "\r", "\r\n")
	goCode = strings.ReplaceAll(goCode, "\r\n", "\n")
	for strings.Contains(goCode, "\n\n") { //将连续的换行符替换成一个
		goCode = strings.ReplaceAll(goCode, "\n\n", "\n")
	}

	lines := strings.Split(goCode, "\n")

	//该结构体中的属性列表
	var properties []PropertyBean

	//表单验证列表
	valids := []string{}

	index := 0
	for index < len(lines) {
		line := lines[index]
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//@") { //这是一个表单验证
			valid := line[3:]
			valids = append(valids, valid)
		}
		matched, _ := regexp.MatchString("[A-Z]]", line)
		if matched { //这是一个结构体属性
			lineArr := strings.Split(line, " ")

			//结构体属性
			property := PropertyBean{
				VarType: lineArr[1],

				//参数名
				Name: lineArr[0],

				//表单验证列表
				Valids: valids,
			}
			properties = append(properties, property)

			//重新初始化表单验证列表
			valids = []string{}
		}
		index++
	}
	formName := regexp.MustCompile("type .* struct").FindAllString(goCode, -1)[0]
	formName = formName[4 : len(formName)-6]
	return FormBean{

		//FORM包所在路径
		PackagePath: packagePath,

		//属性列表
		Properties: properties,

		//结构体名
		Name: formName,
	}
}

// 读取函数名
func readFuncName(line string) string {
	funcName := line[strings.Index(line, "func")+5 : strings.Index(line, "(")]
	return funcName
}

// 读取返回值
func readReturnType(line string) string {
	returnType := line[strings.Index(line, ")")+1 : strings.Index(line, "{")]
	return strings.TrimSpace(returnType)
}
