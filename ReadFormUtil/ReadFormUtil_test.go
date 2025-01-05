package ReadFormUtil

import (
	"GoAutoWeb/Application"
	"fmt"
	"testing"
)

func TestMake(t *testing.T) {
	Application.Init("/Users/zhoulq/dev/java/idea/DairoDFS")
	Make()
	fmt.Println(FormList)
}
