package zip

import (
	"fmt"
	"testing"
)

func Test_concurrentZip(t *testing.T) {
	folders := getFirstSubDirs("D:\\BaiduNetdiskDownload\\掘金小册新2024\\")
	fmt.Println(folders)
	concurrentZip(folders, 5)
}
