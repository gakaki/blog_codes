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

func TestParseMarkdownImagePath(t *testing.T) {
	// 定义正则表达式
	text := " Electron 应用场景的分布\n\n使用 `Electron` 开发的应用品类非常丰富，" +
		"我们看看官网的一些案例展示中的统计数据：\n\n" +
		"<p align=center><img src=\"https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/94cb5ee0174642acb4a24cf2c0fe1fad~tplv-k3u1fbpfcp-jj-mark:0:0:0:0:q75.image#?w=1486&h=1193&s=53343&e=png&a=1&b=546ec6\" alt=\"image.png\"  /></p>\n\n> 数据来源：[Electron ShowCase](https://www.electronjs.org/apps)\n\n可以看到，在使用 Electron 开发的 APP 中，开发者工具、效率应用占据了大半江山。\n"

	//img_pattern := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	//matches := img_pattern.FindAllStringSubmatch(text, -1)
	//for _, match := range matches {
	//	url := match[1]
	//	fmt.Println(url)
	//}
	images := FindImageUrls(text)
	fmt.Println(images)
}

//.image#?w=1486&h=1193&s=53343&e=png&a=1&b=546ec6
