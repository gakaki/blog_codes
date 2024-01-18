package juejinBook

import (
	"encoding/json"
	"fmt"
	"log"
	"myCrawler/utils"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultSaveDir default save dir.
	DefaultSaveDir = "book"
	// GetSectionURL get section url.
	GetSectionURL = "https://api.juejin.cn/booklet_api/v1/section/get"
	// GetBookInfoURL get book info url.
	GetBookInfoURL = "https://api.juejin.cn/booklet_api/v1/booklet/get"
)

type JuejinListRequest struct {
	CategoryID int
	Cursor     string
	Sort       int
	IsVIP      int
	Limit      int
}

func GetAllBookListSortLatestSaveToJSON() {
	// 先从chrome 链接里爬取 , 然后找到 body里的参数进行 修改参数
	url := "https://api.juejin.cn/booklet_api/v1/booklet/listbycategory?aid=2608&uuid=7220793504238650912&spider=0"

	juejinListRequest := JuejinListRequest{
		CategoryID: 0,
		Cursor:     "0",
		Sort:       10,
		IsVIP:      0,
		Limit:      10000,
	}
	payloadByte, err := json.Marshal(juejinListRequest)
	if err != nil {
		return
	}

	response, err := utils.PostToStructInputStruct[JuejinResponse](url, payloadByte, "")
	if err != nil {
		panic(err)
	}
	fmt.Println("掘金一共有", len(response.Data), "本册子")
	utils.WriteJSON(response, "juejin_book.json")
}

func NewConfig(config Config) (*Juejinxiaoce2Markdown, error) {
	if config.Sessionid == "" {
		return nil, fmt.Errorf("sessionid is empty")
	}
	if len(config.BookIDs) == 0 {
		return nil, fmt.Errorf("bookIDs is empty")
	}
	pwd := os.Getenv("PWD")
	if pwd == "" {
		return nil, fmt.Errorf("PWD is empty")
	}
	defaultSaveDir := filepath.Join(pwd, DefaultSaveDir)

	j := &Juejinxiaoce2Markdown{
		ImgPattern:        regexp.MustCompile(`!\[.*?\]\((.*?)\)`),
		Sessionid:         config.Sessionid,
		BookIDs:           config.BookIDs,
		SaveDir:           config.SaveDir,
		RequestHeaders:    map[string]string{"cookie": fmt.Sprintf("sessionid=%s;", config.Sessionid)},
		MarkdownSavePaths: make(map[string]string),
	}

	if j.SaveDir == "" {
		j.SaveDir = defaultSaveDir
	}
	if err := os.MkdirAll(j.SaveDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create save dir failed: %v", err)
	}
	return j, nil
}

func (j *Juejinxiaoce2Markdown) GetSectionRes(sectionID string) (JuejinSectionContent, error) {
	data := map[string]string{
		//"section_id": strconv.FormatInt(sectionID, 10),
		"section_id": sectionID,
	}
	return utils.PostToStructInputStruct[JuejinSectionContent](GetSectionURL, data, j.Sessionid)
}

func (j *Juejinxiaoce2Markdown) GetBookInfoRes(bookID string) (JuejinSection, error) {
	data := map[string]string{
		"booklet_id": bookID,
	}
	return utils.PostToStructInputStruct[JuejinSection](GetBookInfoURL, data, j.Sessionid)
}

func dealBookAndSectionTitle(s string) string {
	tmp := strings.ReplaceAll(s, "\\", "")
	tmp = strings.ReplaceAll(tmp, "/", "")
	tmp = strings.ReplaceAll(tmp, "|", "")
	return tmp
}

func (j *Juejinxiaoce2Markdown) Download() {
	// 并发 下载
	fmt.Println("books need to download count ", len(j.BookIDs))

	maxWorkerCount := 20
	queue := make(chan string, maxWorkerCount)
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{}

	for i := 0; i < maxWorkerCount; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)
			for bookId := range queue {
				err := j.DownloadOneBook(bookId)
				if err != nil {
					fmt.Println("Error", err)
				}
				time.Sleep(time.Second * 1)
			}
		}()
	}
	for _, bookID := range j.BookIDs {
		queue <- bookID
	}
	close(queue)
	wg.Wait()
	//utils.WriteJSON(salttigerItems, "salttigerItems.json")
}

func (j *Juejinxiaoce2Markdown) DownloadOneBook(bookID string) error {
	log.Printf("开始处理小册")

	juejinSection, err := j.GetBookInfoRes(bookID)
	if err != nil {
		return fmt.Errorf("GetBookInfoRes failed: %v", err)
	}

	bookTitle := dealBookAndSectionTitle(juejinSection.Data.Booklet.BaseInfo.Title)
	log.Printf("book_title: %s", bookTitle)
	bookSavePath := filepath.Join(j.SaveDir, bookTitle)
	if err := os.MkdirAll(bookSavePath, os.ModePerm); err != nil {
		return fmt.Errorf("create book save path failed: %v", err)
	}
	imgDir := filepath.Join(bookSavePath, "img")
	if err := os.MkdirAll(imgDir, os.ModePerm); err != nil {
		return fmt.Errorf("create img dir failed: %v", err)
	}

	sectionIDList := make([]string, 0, len(juejinSection.Data.Sections))
	for _, section := range juejinSection.Data.Sections {
		sectionIDList = append(sectionIDList, section.SectionID)
	}

	sectionTotalLength := len(sectionIDList)
	for sectionIndex, sectionID := range sectionIDList {
		sectionOrder := sectionIndex + 1

		juejinSectionContent, err := j.GetSectionRes(sectionID)
		if err != nil {
			return fmt.Errorf("GetSectionRes failed: %v", err)
		}

		sectionTitle := dealBookAndSectionTitle(juejinSectionContent.Data.Section.Title)
		markdownStr := juejinSectionContent.Data.Section.MarkdownShow
		markdownFilePath := filepath.Join(bookSavePath, fmt.Sprintf("%d-%s.md", sectionOrder, sectionTitle))
		sectionImgDir := filepath.Join(imgDir, strconv.Itoa(sectionOrder))

		log.Printf("进度: %d/%d, 处理 section >> %s", sectionOrder, sectionTotalLength, sectionTitle)

		if err := os.MkdirAll(sectionImgDir, os.ModePerm); err != nil {
			return fmt.Errorf("create section img dir failed: %v", err)
		}
		markdownRelativeImgDir := filepath.Join("img", strconv.Itoa(sectionOrder))
		j.MarkdownSavePaths[sectionID] = markdownFilePath

		j.SaveMarkdown(sectionIndex, markdownFilePath, sectionImgDir, markdownRelativeImgDir, markdownStr)
	}
	log.Printf("处理完成")
	return nil
}

func FindImageUrls(htmls string) []string {
	//imgPattern := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imgs := imgRE.FindAllStringSubmatch(htmls, -1)
	out := make([]string, 0)
	for _, img := range imgs {
		if strings.Contains(img[1], "https") {
			out = append(out, img[1])
		}
	}
	return out
}

func (j *Juejinxiaoce2Markdown) SaveMarkdown(sectionIndex int, markdownFilePath string, sectionImgDir string, markdownRelativeImgDir string, markdownStr string) {

	imgUrls := FindImageUrls(markdownStr)
	for imgIndex, imgUrl := range imgUrls {
		newImgUrl := strings.TrimSpace(imgUrl) // Remove newlines and extra spaces
		if strings.HasPrefix(newImgUrl, "//") {
			newImgUrl = "https:" + newImgUrl // Add https:// if missing
		}

		suffix := filepath.Ext(newImgUrl)
		suffix = ".png"                                                         // Get file extension
		imgFileName := fmt.Sprintf("%d%s", imgIndex+1, suffix)                  // Generate filename
		mdRelativeImgPath := filepath.Join(markdownRelativeImgDir, imgFileName) // Relative path for Markdown
		imgSavePath := filepath.Join(sectionImgDir, imgFileName)                // Full path to save image

		// Download image
		err := utils.RequestThanSaveImage(newImgUrl, imgSavePath)
		if err != nil {
			fmt.Println("Error downloading image:", err)
		}
		// Replace URL in Markdown string with relative path
		markdownStr = strings.ReplaceAll(markdownStr, imgUrl, mdRelativeImgPath)
	}

	err := os.WriteFile(markdownFilePath, []byte(markdownStr), 0644)
	if err != nil {
		fmt.Println("Error saving Markdown file:", err)
	}
}
