package main

import (
	"GoAutoWeb/Application"
	"GoAutoWeb/MakeSourceUtil"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		currentFolder := os.Args[0]
		currentFolder = strings.ReplaceAll(currentFolder, "\\", "/")
		currentFolder = currentFolder[:strings.LastIndex(currentFolder, "/")]
		Application.Init(currentFolder)
	} else {
		Application.Init(os.Args[1])
	}
	fmt.Println(Application.RootProject)
	MakeSourceUtil.Make()
}
