package juejinBook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	// DefaultSaveDir default save dir.
	DefaultSaveDir = "book"
	// GetSectionURL get section url.
	GetSectionURL = "https://api.juejin.cn/booklet_api/v1/section/get"
	// GetBookInfoURL get book info url.
	GetBookInfoURL = "https://api.juejin.cn/booklet_api/v1/booklet/get"
)

// NewJuejinxiaoce2Markdown new juejinxiaoce2Markdown.
func NewJuejinxiaoce2Markdown(config Config) (*Juejinxiaoce2Markdown, error) {
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
		MarkdownSavePaths: make(map[int64]string),
	}

	if j.SaveDir == "" {
		j.SaveDir = defaultSaveDir
	}
	if err := os.MkdirAll(j.SaveDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create save dir failed: %v", err)
	}
	return j, nil
}

// Config config.
type Config struct {
	Sessionid string  `yaml:"sessionid"`
	BookIDs   []int64 `yaml:"book_ids"`
	SaveDir   string  `yaml:"save_dir"`
}

// GetSectionRes get section res.
func (j *Juejinxiaoce2Markdown) GetSectionRes(sectionID int64) (*http.Response, error) {
	data := map[string]string{
		"section_id": strconv.FormatInt(sectionID, 10),
	}
	return j.PostJSON(GetSectionURL, data)
}

// GetBookInfoRes get book info res.
func (j *Juejinxiaoce2Markdown) GetBookInfoRes(bookID int64) (*http.Response, error) {
	data := map[string]string{
		"booklet_id": strconv.FormatInt(bookID, 10),
	}
	return j.PostJSON(GetBookInfoURL, data)
}

// PostJSON post json.
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

// ClearSlash clear slash.
func ClearSlash(s string) string {
	return strings.ReplaceAll(s, "\\", "").ReplaceAll("/", "").ReplaceAll("|", "")
}

// DealABook deal a book.
func (j *Juejinxiaoce2Markdown) DealABook(bookID int64) error {
	log.Printf("开始处理小册")

	res, err := j.GetBookInfoRes(bookID)
	if err != nil {
		return fmt.Errorf("GetBookInfoRes failed: %v", err)
	}
	var resData struct {
		Data struct {
			Booklet struct {
				BaseInfo struct {
					Title string `json:"title"`
				} `json:"base_info"`
			} `json:"booklet"`
			Sections []struct {
				SectionID int64 `json:"section_id"`
			} `json:"sections"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return fmt.Errorf("json.NewDecoder.Decode failed: %v", err)
	}

	bookTitle := ClearSlash(resData.Data.Booklet.BaseInfo.Title)
	log.Printf("book_title: %s", bookTitle)
	bookSavePath := filepath.Join(j.saveDir, bookTitle)
	if err := os.MkdirAll(bookSavePath, os.ModePerm); err != nil {
		return fmt.Errorf("create book save path failed: %v", err)
	}
	imgDir := filepath.Join(bookSavePath, "img")
	if err := os.MkdirAll(imgDir, os.ModePerm); err != nil {
		return fmt.Errorf("create img dir failed: %v", err)
	}

	sectionIDList := make([]int64, 0, len(resData.Data.Sections))
	for _, section := range resData.Data.Sections {
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
		var resData struct {
			Data struct {
				Section struct {
					Title        string `json:"title"`
					MarkdownShow string `json:"markdown_show"`
				} `json:"section"`
			} `json:"data"`
		}
		if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
			return fmt.Errorf("json.NewDecoder.Decode failed: %v", err)
		}

		sectionTitle := ClearSlash(resData.Data.Section.Title)
		markdownStr := resData.Data.Section.MarkdownShow
		markdownFilePath := filepath.Join(bookSavePath, fmt.Sprintf("%d-%s.md", sectionOrder, sectionTitle))
		sectionImgDir := filepath.Join(imgDir, strconv.Itoa(sectionOrder))
		if err := os.MkdirAll(sectionImgDir, os.ModePerm); err != nil {
			return fmt.Errorf("create section img dir failed: %v", err)
		}
		markdownRelativeImgDir := filepath.Join("img", strconv.Itoa(sectionOrder))
		j.MarkdownSavePaths[sectionID] = markdownFilePath
		j.SaveMarkdown(markdownFilePath, sectionImgDir, markdownRelativeImgDir, markdownStr)

		j.DealABook(bookID)
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
	}
}
