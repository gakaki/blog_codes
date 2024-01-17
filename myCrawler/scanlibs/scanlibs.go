package scanlibs

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"myCrawler/utils"
	"strconv"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		fmt.Printf("%d", err)
	}
}

func getEveryPageItems(doc *goquery.Document) []*ScanLibItem {
	scanlibItems := make([]*ScanLibItem, 0)
	doc.Find("article").Each(func(index int, ele *goquery.Selection) {
		var item = &ScanLibItem{}
		item.ID, _ = ele.Attr("id")
		item.URL, _ = ele.Find("a").Attr("href")
		item.CreatedAt, _ = ele.Find("time").Attr("datetime")
		item.Thumbnil, _ = ele.Find("img.aligncenter").Attr("src")
		item.Title = ele.Find("b").Text()
		item.SubTitle = ele.Find("p").Text()
		item.IsVideo = strings.Index(item.URL, "video") >= 0
		scanlibItems = append(scanlibItems, item)
	})
	return scanlibItems
}
func makePageUrl(pageIndex int) string {
	url := fmt.Sprintf("https://scanlibs.com/page/%d/", pageIndex)
	return url
}

func setEveryPage(page *ScanLibPage) {
	fmt.Println("request page ", page.Index)
	doc, err := utils.RequestGetDocument(page.URL)
	if err != nil {
		CheckError(err)
		return
	}
	items := getEveryPageItems(doc)
	page.Items = items
	utils.WriteJSON(items, fmt.Sprintf("json/scanlibs/page-%d.json", page.Index))
}

func parseThanGetLastPageIndex(doc *goquery.Document) int {
	t := doc.Find(".page-numbers:nth-last-child(2)").Text()
	t = strings.Replace(t, ",", "", 2)
	d, _ := strconv.Atoi(t)
	return d
}

func getTotalPages() []*ScanLibPage {
	url := "https://scanlibs.com/page/1"
	doc, err := utils.RequestGetDocument(url)
	if err != nil {
		CheckError(err)
		panic(err)
	}

	totalPageIndex := parseThanGetLastPageIndex(doc)
	scanLibPages := make([]*ScanLibPage, 0)
	for i := 1; i <= totalPageIndex; i++ {
		page := &ScanLibPage{}
		page.Index = i
		page.URL = makePageUrl(i)
		scanLibPages = append(scanLibPages, page)
	}
	fmt.Println("总共有多少页", totalPageIndex)
	return scanLibPages
}
