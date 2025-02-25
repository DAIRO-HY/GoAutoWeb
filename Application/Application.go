package Application

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// 程序启动参数
var Args appArgs

// 程序启动参数
type appArgs struct {

	// 源码路径
	SourceDir string

	// 生成目标代码类型，{web,flutter-api}
	TargetType string

	// 目标代码目录
	TargetDir string
}

func Init() {
	parseArgs()
}

// 解析参数
func parseArgs() {
	fmt.Println("------------------------------------------------------------------------")
	fmt.Println(strings.Join(os.Args, " "))
	fmt.Println("------------------------------------------------------------------------")
	Args = appArgs{
		TargetType: "web",
	}
	if len(os.Args) == 1 { //没有设置任何参数时
		source := ""
		source = os.Args[0]
		source = strings.ReplaceAll(source, "\\", "/")
		source = source[:strings.LastIndex(source, "/")]
		Args.SourceDir = source
	}
	argsElem := reflect.ValueOf(&Args).Elem()
	for i := 0; i < len(os.Args); i++ {
		key := os.Args[i]
		if !strings.HasPrefix(key, "--") {
			continue
		}

		filedName := ""
		for _, it := range strings.Split(key[2:], "-") {
			if len(it) == 0 {
				continue
			}
			// 将字符串转换为 rune 切片以便处理 Unicode 字符
			r := []rune(it)

			// 将首字母大写
			r[0] = unicode.ToUpper(r[0])
			filedName += string(r)
		}
		field := argsElem.FieldByName(filedName)
		if !field.IsValid() { //如果该字段不存在
			continue
		}
		if i+1 > len(os.Args)-1 { //已经没有下一个元素
			break
		}

		//参数值
		value := os.Args[i+1]
		switch field.Kind() {

		//整数类型转换
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Panicf("参数%s=%s发生了转换错误:%q", key, value, err)
			}
			field.SetInt(v)
		case reflect.Float32, reflect.Float64:
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Panicf("参数%s=%s发生了转换错误:%q", key, value, err)
			}
			field.SetFloat(v)
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				log.Panicf("参数%s=%s发生了转换错误:%q", key, value, err)
			}
			field.SetBool(v)
		case reflect.String:
			field.SetString(value)
		}
		i++
	}
}
