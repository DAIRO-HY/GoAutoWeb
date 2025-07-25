package ReadFormUtil

import (
	"GoAutoWeb/FileUtil"
	"GoAutoWeb/Global"
	"regexp"
	"strings"
)

// 表单列表
var FormList []FormBean

// 表单Map
var FormMap = map[string]FormBean{}

func Make() {
	for _, fPath := range Global.GoFileList {
		if !strings.HasSuffix(fPath, "Form.go") { //过滤文件后缀
			continue
		}

		// 读取go文件里的路由配置
		formBean := readControllerPath(fPath)
		FormList = append(FormList, formBean)

		FormMap[formBean.PackagePath+"/"+formBean.Name] = formBean
	}
}

func readControllerPath(path string) FormBean {

	//包所在路径
	packagePath := path[len(Global.RootProject):]
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
	for strings.Contains(goCode, "  ") { //将连续的空格替换成一个
		goCode = strings.ReplaceAll(goCode, "  ", " ")
	}

	lines := strings.Split(goCode, "\n")

	//该结构体中的属性列表
	var properties []FormPropertyBean

	//该结构体中的函数列表
	var functions []FormFunctionBean

	//表单验证
	valids := make([]FormValidateBean, 0)

	//注释内容
	comment := ""
	index := 0
	for index < len(lines) {
		line := lines[index]
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			index++
			continue
		}

		//读取表单验证规则
		validBean := readValidate(line)
		if validBean != nil {
			valids = append(valids, *validBean)
		} else if line[0] >= 65 && line[0] <= 90 && strings.Contains(line, " ") { //首字母是大写,判定这是一个属性
			property := readProperty(line)
			property.valids = valids
			property.Comment = comment

			//清空当前属性的验证规则
			valids = make([]FormValidateBean, 0)
			properties = append(properties, property)
			comment = ""
		} else if strings.HasPrefix(line, "func") { //这是一个函数
			function := readFunction(line)
			if function != nil {
				functions = append(functions, *function)
			}
		} else if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			comment += line + "\n"
		}
		index++
	}
	formName := regexp.MustCompile("type .* struct").FindAllString(goCode, -1)[0]
	formName = formName[4 : len(formName)-6]
	formName = strings.TrimSpace(formName)
	return FormBean{

		//FORM包所在路径
		PackagePath: packagePath,

		//属性列表
		Properties: properties,

		//结构体函数列表
		Functions: functions,

		//结构体名
		Name: formName,
	}
}

// 读取结构体属性
func readProperty(line string) FormPropertyBean {

	//通过空格分隔单词
	words := strings.Fields(line)

	//结构体属性
	return FormPropertyBean{

		//参数名
		Name: words[0],

		//参数类型
		VarType: words[1],
	}
}

// 读取表单验证规则
func readValidate(line string) *FormValidateBean {

	//去掉表单验证规则中多余的空格
	line = strings.ReplaceAll(line, "// @", "//@")
	if !strings.HasPrefix(line, "//@") { //这不是一个表单验证
		return nil
	}
	line = line[3:]

	//验证规则名称
	validName := ""

	//验证规则的参数
	validArgs := make(map[string]string)
	startKHIndex := strings.Index(line, "(")
	if startKHIndex == -1 { //没有括号
		validName = line
	} else {
		validName = line[0:startKHIndex]
		validParamStr := line[startKHIndex+1 : strings.Index(line, ")")]
		validParamArr := strings.Split(validParamStr, ",")
		for _, it := range validParamArr {
			kv := strings.Split(it, "=")
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			validArgs[key] = value
		}
	}
	return &FormValidateBean{
		Name: strings.ToUpper(validName),
		Args: validArgs,
	}
}

// 读取结构体函数
func readFunction(line string) *FormFunctionBean {
	regResults := regexp.MustCompile("Form\\).+\\(").FindAllString(line, -1)
	if len(regResults) == 0 { //这不是一个结构体函数
		return nil
	}
	name := regResults[0]
	name = name[strings.Index(name, ")")+1 : len(name)-1]
	name = strings.TrimSpace(name)

	returnType := regexp.MustCompile("\\(\\s*\\).*{").FindAllString(line, -1)[0]
	returnType = returnType[strings.Index(returnType, ")")+1 : len(returnType)-1]
	returnType = strings.TrimSpace(returnType)

	//结构体函数
	return &FormFunctionBean{

		//函数名
		Name: name,

		//返回值类型
		ReturnType: returnType,
	}
}
