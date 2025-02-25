package Global

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const count = 10000000

// 读取项目的模块名
func TestMakeFileList1(t *testing.T) {
	now := time.Now().UnixMilli()
	for i := 0; i < count; i++ {
		filepath.WalkDir(RootProject, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if strings.HasSuffix(path, ".idea") || strings.HasSuffix(path, ".git") {
					return filepath.SkipDir
				}
				return nil
			}
			return nil
		})
	}
	fmt.Println(time.Now().UnixMilli() - now)
}

// 读取项目的模块名
func TestMakeFileList2(t *testing.T) {
	now := time.Now().UnixMilli()
	for i := 0; i < count; i++ {
		filepath.Walk(RootProject, func(path string, d os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return nil
		})
	}
	fmt.Println(time.Now().UnixMilli() - now)
}
