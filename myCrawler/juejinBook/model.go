package juejinBook

import "regexp"

// Juejinxiaoce2Markdown juejinxiaoce2Markdown.
type Juejinxiaoce2Markdown struct {
	ImgPattern        *regexp.Regexp
	Sessionid         string
	BookIDs           []int64
	SaveDir           string
	RequestHeaders    map[string]string
	MarkdownSavePaths map[int64]string
}
type JuejinBook struct {
	BookID int64
	Title  string
	Cover  string
}

// Juejinxiaoce2Markdown juejinxiaoce2Markdown.
type JuejinBookDownloadConfig struct {
	SessionID         string
	Books             string
	BookIDs           []int64
	SaveDir           string
	RequestHeaders    map[string]string
	MarkdownSavePaths map[int64]string
}
