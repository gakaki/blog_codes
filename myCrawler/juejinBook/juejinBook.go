package juejinBook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"myCrawler/utils"
	"net/http"
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

func RequestPostJSON[T interface{}](url string) (t T, err error) {

	payload := strings.NewReader(`{
		"category_id": "0",
		"cursor": "0",
		"sort": 10,
		"is_vip": 0,
		"limit": 1000
	}`)

	resp, err := utils.GetHttpClient().R().
		//EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(url)

	if err != nil {
		return t, err
	}
	//fmt.Println("Response Info:")
	//fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode())

	if resp.StatusCode() == http.StatusOK {
		//bodyString := string(resp.Body())
		json.Unmarshal(resp.Body(), &t)
		return t, nil
	} else {
		fmt.Println("错误号码:")
		panic(fmt.Sprintf("url is %s ,status code is %d %s", url, resp.StatusCode(), string(resp.Body())))
		return t, err
	}
}

func GetAllBookListSortLatestSaveToJSON() {
	// 先从chrome 链接里爬取 , 然后找到 body里的参数进行 修改参数
	url := "https://api.juejin.cn/booklet_api/v1/booklet/listbycategory?aid=2608&uuid=7220793504238650912&spider=0"
	response, err := RequestPostJSON[JuejinResponse](url)
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

func (j *Juejinxiaoce2Markdown) GetSectionRes(sectionID string) (*http.Response, error) {
	data := map[string]string{
		//"section_id": strconv.FormatInt(sectionID, 10),
		"section_id": sectionID,
	}
	return j.PostJSON(GetSectionURL, data)
}

// GetBookInfoRes get book info res.
func (j *Juejinxiaoce2Markdown) GetBookInfoRes(bookID string) (*http.Response, error) {
	data := map[string]string{
		"booklet_id": bookID,
		//"booklet_id": strconv.FormatInt(bookID, 10),
	}
	return j.PostJSON(GetBookInfoURL, data)
}

func (j *Juejinxiaoce2Markdown) PostJSON(url string, data interface{}) (*http.Response, error) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("http.Post failed: %v", err)
	}
	return res, nil
}

func dealBookTitle(s string) string {
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
				err := j.DealABook(bookId)
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

// DealABook deal a book.
func (j *Juejinxiaoce2Markdown) DealABook(bookID string) error {
	log.Printf("开始处理小册")

	res, err := j.GetBookInfoRes(bookID)
	if err != nil {
		return fmt.Errorf("GetBookInfoRes failed: %v", err)
	}

	var juejinSection JuejinSection
	if err := json.NewDecoder(res.Body).Decode(&juejinSection); err != nil {
		return fmt.Errorf("json.NewDecoder.Decode failed: %v", err)
	}

	bookTitle := dealBookTitle(juejinSection.Data.Booklet.BaseInfo.Title)
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

	sectionCount := len(sectionIDList)
	for index, sectionID := range sectionIDList {
		sectionOrder := index + 1
		log.Printf("进度: %d/%d, 处理 section", sectionOrder, sectionCount)

		res, err := j.GetSectionRes(sectionID)
		if err != nil {
			return fmt.Errorf("GetSectionRes failed: %v", err)
		}
		var juejinSectionContent JuejinSectionContent
		if err := json.NewDecoder(res.Body).Decode(&juejinSectionContent); err != nil {
			return fmt.Errorf("json.NewDecoder.Decode failed: %v", err)
		}

		sectionTitle := dealBookTitle(juejinSectionContent.Data.Section.Title)
		markdownStr := juejinSectionContent.Data.Section.MarkdownShow
		markdownFilePath := filepath.Join(bookSavePath, fmt.Sprintf("%d-%s.md", sectionOrder, sectionTitle))
		sectionImgDir := filepath.Join(imgDir, strconv.Itoa(sectionOrder))
		if err := os.MkdirAll(sectionImgDir, os.ModePerm); err != nil {
			return fmt.Errorf("create section img dir failed: %v", err)
		}
		markdownRelativeImgDir := filepath.Join("img", strconv.Itoa(sectionOrder))
		j.MarkdownSavePaths[sectionID] = markdownFilePath
		j.SaveMarkdown(markdownFilePath, sectionImgDir, markdownRelativeImgDir, markdownStr)

	}
	log.Printf("处理完成")
	return nil
}

// SaveMarkdown save markdown.
func (j *Juejinxiaoce2Markdown) SaveMarkdown(markdownFilePath, sectionImgDir, markdownRelativeImgDir, markdownStr string) {
	imgURLs := j.ImgPattern.FindAllString(markdownStr, -1)
	for imgIndex, imgURL := range imgURLs {
		newImgURL := strings.ReplaceAll(imgURL, "\n", "")
		if strings.HasPrefix(newImgURL, "//") {
			newImgURL = "https:" + newImgURL
		}
		imgFileName := fmt.Sprintf("%d%s", imgIndex+1, filepath.Ext(newImgURL))
		fmt.Println(imgFileName)
	}
}
