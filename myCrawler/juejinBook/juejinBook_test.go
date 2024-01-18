package juejinBook

import (
	"fmt"
	"testing"
)

func TestGetAllXiaoces(t *testing.T) {
	GetAllBookListSortLatestSaveToJSON()
}

func TestDownload2Markdown(t *testing.T) {

	c := Config{
		Sessionid: "a115bed3665f8f179576610d73a3eae0",
		BookIDs: []string{
			"7302990019642261567",
		},
		SaveDir: "",
	}

	juejin, err := NewConfig(c)
	if err != nil {
		fmt.Println(err)
	}
	juejin.Download()
}
