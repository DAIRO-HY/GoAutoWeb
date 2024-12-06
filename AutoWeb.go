/**
 * 代码为自动生成，请勿手动修改
 */
package main

import (
	//{IMPORT}
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func init() {
	//{BODY}
}

// 生成参数Map
func makeParamMap(request *http.Request) map[string][]string {
	query := request.URL.Query()

	//解析post表单
	request.ParseForm()
	postParams := request.PostForm

	//将参数转换成Map
	paramMap := make(map[string][]string)
	for key, v := range query {
		paramMap[key] = v
	}
	for key, v := range postParams {
		paramMap[key] = v
	}
	return paramMap
}

// 获取表单实例
func getForm[T any](paramMap map[string][]string) T {

	// 创建结构体实例
	targetForm := new(T)
	reflectForm := reflect.ValueOf(targetForm).Elem()
	argType := reflect.TypeOf(*targetForm)

	// 遍历结构体字段
	for j := 0; j < argType.NumField(); j++ {
		field := argType.Field(j)
		fieldName := field.Name

		//得到参数值
		value := paramMap[fieldName]
		if value == nil {
			//将首字母小写再去获取参数
			lowerKey := strings.ToLower(fieldName[:1]) + fieldName[1:]
			value = paramMap[lowerKey]
		}
		if value == nil {
			continue
		}

		// 设置字段值（这里我们设置为示例值）
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

			// 设置整数字段
			intValue, _ := strconv.ParseInt(value[0], 10, 64)
			reflectForm.Field(j).SetInt(intValue)
		case reflect.Float32, reflect.Float64:
			floatValue, _ := strconv.ParseFloat(value[0], 64)
			reflectForm.Field(j).SetFloat(floatValue)
		case reflect.String:
			reflectForm.Field(j).SetString(value[0]) // 设置字符串字段
		}
	}
	return *targetForm
}

// 获取string类型的参数
func getString(paramMap map[string][]string, key string) string {
	value := paramMap[key]
	if value == nil {
		return ""
	}
	rValue := value[0]
	return rValue
}

// 获取int类型的参数
func getInt(paramMap map[string][]string, key string) int {
	value := paramMap[key]
	if value == nil {
		return 0
	}
	rValue, _ := strconv.Atoi(value[0])
	return rValue
}

// 获取int类型的参数
func getInt64(paramMap map[string][]string, key string) int64 {
	value := paramMap[key]
	if value == nil {
		return 0
	}
	rValue, _ := strconv.ParseInt(value[0], 10, 64)
	return rValue
}

// 获取float32类型的参数
func getFloat32(paramMap map[string][]string, key string) float32 {
	value := paramMap[key]
	if value == nil {
		return 0
	}
	rValue, _ := strconv.ParseFloat(value[0], 32)
	return float32(rValue)
}

// 获取float64类型的参数
func getFloat64(paramMap map[string][]string, key string) float64 {
	value := paramMap[key]
	if value == nil {
		return 0
	}
	rValue, _ := strconv.ParseFloat(value[0], 64)
	return rValue
}

// 返回结果
func writeToResponse(writer http.ResponseWriter, body any) {
	if body == nil {
		return
	}
	if body == "" {
		return
	}

	// 设置 Content-Type 头部信息
	writer.Header().Set("Content-Type", "text/plain;charset=UTF-8")

	switch returnBody := body.(type) {
	case string:
		writer.Write([]uint8(returnBody))
	case int:
		writer.Write([]uint8(strconv.Itoa(returnBody)))
	case int8:
		writer.Write([]uint8(strconv.Itoa(int(returnBody))))
	case int16:
		writer.Write([]uint8(strconv.Itoa(int(returnBody))))
	case int32:
		writer.Write([]uint8(strconv.Itoa(int(returnBody))))
	case int64:
		writer.Write([]uint8(strconv.FormatInt(returnBody, 10)))
	case error:
		// 设置 HTTP 状态码
		writer.WriteHeader(http.StatusInternalServerError) // 设置状态码
		jsonData, _ := json.Marshal(body)
		writer.Write(jsonData)
	default:
		jsonData, _ := json.Marshal(body)
		writer.Write(jsonData)
	}
}
