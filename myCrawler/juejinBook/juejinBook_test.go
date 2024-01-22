package juejinBook

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetAllXiaoces(t *testing.T) {
	GetAllBookListSortLatestSaveToJSON()
}

func TestDownload2Markdown(t *testing.T) {

	c := Config{
		Sessionid: "43680406158c713253e8bfe966d70f80",
		BookIDs: []string{
			//"7302990019642261567",
			//"6918979822425210891", // 0 打造通用型低
			//"7202598408815640631", //前端依赖治理
			//"7269673629964173331", // 前端可视化入门与实战
			"7288940354408022074", //web动画之旅
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
	images := FindImageUrls(1, text)
	fmt.Println(images)
}

func TestRenderPDF(t *testing.T) {
	//brew install pandoc
	//brew install --cask basictex

	// 遍历文件夹
	filepath.Walk("./book", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() == false && filepath.Ext(info.Name()) == ".md" {
			// 获取文件名
			fileName := filepath.Base(info.Name())
			// 切换到当前文件夹
			os.Chdir(filepath.Dir(path))

			newFileName := strings.Replace(fileName, filepath.Ext(fileName), "", 1)
			// 打印转换开始信息
			fmt.Println("转换开始：" + "pandoc " + fileName + " -o " + newFileName + ".pdf")

			// 调用 pandoc 进行格式转换
			cmd := exec.Command(fmt.Sprintf("pandoc %s -o %s", fileName, newFileName))
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("combined out:\n%s\n", string(out))
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}

			// 打印转换完成信息
			fmt.Println("转换完成...")
		}
		return err
	})

}
