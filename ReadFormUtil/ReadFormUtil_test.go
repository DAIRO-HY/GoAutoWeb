package ReadFormUtil

import (
	"fmt"
	"testing"
)

func TestMake(t *testing.T) {
	Global.Init("/Users/zhoulq/dev/java/idea/DairoDFS")
	Make()
	fmt.Println(FormList)
}
