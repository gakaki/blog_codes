package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"myCrawler/mgstage"
	"myCrawler/utils"
	"os"
	"strings"
	"sync"

	"net/url"
	"time"
)

func check(err error) {
	if err != nil {
		fmt.Println("出错了 请查看", err)
		panic(err)
	}
}

func getMgstageTodayURl() string {
	yyyymmdd := getTimeYYYYMMDD()
	url := "https://www.mgstage.com/search/cSearch.php?search_word=&sale_start_range=%s-%s&sort=new&list_cnt=120&type=top"
	return fmt.Sprintf(url, yyyymmdd, yyyymmdd)
}
func getTimeYYYYMMDD() string {
	t := time.Now() //.Add(-3600*24*time.Second)
	return fmt.Sprintf("%d.%02d.%02d",
		t.Year(), t.Month(), t.Day()-1)
}
func getTimeYYYYMMDDHHMMSS() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func getSavePath(url string, id string) string {
	urlSplits := strings.Split(url, "/")

	fileDir := fmt.Sprintf("assets/%s", id)
	err := os.Mkdir(fileDir, os.ModePerm)
	checkError(err)
	savePath := fmt.Sprintf("assets/%s/%s", id, urlSplits[len(urlSplits)-1])
	return savePath
}

func downlaodThanSave(url string, path string) {
	if url == "" {
		return
	}
	responseStr, e := utils.RequestString(url)
	if e != nil {
		log.Fatal(e)
	}

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(responseStr))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("download success", url, path)
}

func downloadImages(v *mgstage.Video) {
	for _, imageUrl := range v.Images {
		downlaodThanSave(string(imageUrl), getSavePath(string(imageUrl), v.ID))
	}
}

var queue = make(chan *mgstage.Video, 3)
var queueVideos = make(chan *mgstage.Video, 3)
var queueImages = make(chan *mgstage.Video, 10)
var wg sync.WaitGroup

var videoMap = map[string]*mgstage.Video{}
var videoArray = make([]*mgstage.Video, 0)

func DetailPageToImagesVideo() {
	for v := range queue {
		//go func(v *mgstage.Video) {
		fmt.Println("=============")

		fmt.Println("queue url  full", v.FullUrl)
		detailDoc, err := utils.RequestGetDocument(v.FullUrl)
		if err != nil {
			log.Fatal(err)
		}
		getVideoModel(v, detailDoc)

		if err != nil {
			log.Fatal(err)
		}
		GetVideoMP4(v)
		fmt.Println("queue url video", v.VideoUrl)
		//为什么这里不需要使用go func呢{} 因为buffered channel吗
		queueVideos <- v
		queueImages <- v
		//}(Video)

	}
}
func getList() {

	mgstageTodayURl := getMgstageTodayURl()
	fmt.Println("start mgstage", mgstageTodayURl)

	doc, err := utils.RequestGetDocument(mgstageTodayURl)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("div.search_list > div > ul > li").Each(func(index int, ele *goquery.Selection) {
		v := mgstage.Video{}
		v.Url, _ = getVideoUrl(ele)
		v.FullUrl = getFullUrl(v.Url)
		v.ID = getVideoId(v.Url)
		v.Price = getPrice(ele)
		videoMap[v.ID] = &v
		videoArray = append(videoArray, &v)
		//fmt.Println(v.FullUrl)
		queue <- &v
	})

}
func downloadVideo(v *mgstage.Video) {
	if v.VideoUrl != "" {
		fmt.Printf("视频文件 %s ", v.VideoUrl)
		downlaodThanSave(string(v.VideoUrl), getSavePath(string(v.VideoUrl), v.ID))
	}
	//if(v.VideoUrl==""){
	//	fmt.Printf("没有找到视频文件 %s %s",v.ID,v.FullUrl)
	//	return
	//}
}
func GetVideoMP4(v *mgstage.Video) {
	toGetVideoUrl := fmt.Sprintf("https://www.mgstage.com/sampleplayer/sampleRespons.php?pid=%s", v.Pid)
	jsonString, err := utils.RequestString(toGetVideoUrl)
	if err == nil {
		fmt.Printf("%v", err)
	}
	var mgstageVideo JSONMgstageVideo
	err = json.Unmarshal([]byte(jsonString), &mgstageVideo)
	if err == nil {
		//fmt.Printf("%#v \n %#v \n", err, mgstageVideo.Url)
	}
	if mgstageVideo.Url != "" {
		v.VideoUrl = getMp4UrlFromISMURL(mgstageVideo.Url)
	} else {
		fmt.Println("mgstage Video url 为空", v.FullUrl)
	}
}

type JSONMgstageVideo struct {
	Url string `json:"url"`
}

func getLinkByItemDoc(selectorTd string, doc *goquery.Document) mgstage.Link {
	tdEle := doc.Find(selectorTd)
	linkElem := tdEle.Find("a")
	href, _ := linkElem.Attr("href")
	text := tdEle.Text()

	text = cleanText(text)
	href = cleanText(href)
	if href != "" {
		href = getFullUrl(href)
	}
	return mgstage.Link{text, href}
}
func getLinkByItemSelection(selectorTd string, doc *goquery.Selection) mgstage.Link {
	linkElem := doc.Find(selectorTd)
	href, _ := linkElem.Attr("href")
	text := linkElem.Text()
	text = cleanText(text)
	href = cleanText(href)
	if href != "" {
		href = getFullUrl(href)
	}
	return mgstage.Link{text, href}
}
func cleanText(str string) string {
	s := strings.TrimLeft(str, " ")
	s = strings.Trim(str, "\n")
	s = strings.Trim(str, "\r")
	s = strings.TrimSpace(str)
	return s
}

func getMp4UrlFromISMURL(url string) string {
	//url:="https://sample.mgstage.com/sample/nanpatv/200gana/2396/200gana-2396_20201202T125901.ism/request?uid=10000000-0000-0000-0000-00000000000a&amp;pid=16b6ae62-e6d6-412c-afda-c8b4709c86eb"
	end := strings.Index(url, "/request?")
	lastStr := url[:end]
	lastStr = strings.Replace(lastStr, ".ism", ".mp4", 1)
	//fmt.Println(lastStr)
	return lastStr
}

func getDocText(selector string, doc *goquery.Document) string {
	t := doc.Find(selector).Text()
	return cleanText(t)
}
func getVideoModel(v *mgstage.Video, detailDoc *goquery.Document) (*mgstage.Video, error) {

	v.Title = getDocText("#center_column > div.common_detail_cover > h1", detailDoc)

	v.CountFavorite = getDocText("#playing > dl.detail_fav_cnt", detailDoc)
	v.CountPlay = detailDoc.Find("#playing > dl.playing").Text()

	v.Actor = getLinkByItemDoc("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(1) > td", detailDoc)
	v.Manufacturer = getLinkByItemDoc("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(2) > td", detailDoc)

	v.TimeLong = getDocText("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(3) > td", detailDoc)
	v.StartDate = getDocText("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(5) > td", detailDoc)
	v.SaleDate = getDocText("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(6) > td", detailDoc)

	v.Series = getLinkByItemDoc("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(7) > td", detailDoc)
	v.Company = getLinkByItemDoc("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(7) > td", detailDoc)

	detailDoc.Find("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(9) > td").Each(func(index int, ele *goquery.Selection) {
		link := getLinkByItemSelection("a", ele)
		v.Tags = append(v.Tags, link)
	})

	v.Rate = getDocText("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(11) > td > span", detailDoc)

	v.ImageHead, _ = detailDoc.Find("#center_column > div.common_detail_cover > div.detail_left > div > div > h2 > img").Attr("src")
	v.Description = getDocText("#introduction > dd > p.txt.introduction", detailDoc)

	//photos
	detailDoc.Find("#sample-photo > dd > ul > li > a").Each(func(index int, ele *goquery.Selection) {
		link, _ := ele.Attr("href")
		v.Images = append(v.Images, link)
	})

	pidUrl, _ := detailDoc.Find(".button_sample").Attr("href")
	pidUrlSplits := strings.Split(pidUrl, "/")
	v.Pid = pidUrlSplits[len(pidUrlSplits)-1]
	return v, nil
}

func getVideoUrl(ele *goquery.Selection) (string, bool) {
	return ele.Find("a").Eq(0).Attr("href")
}

func getPrice(ele *goquery.Selection) string {
	return ele.Find(".price").Text()
}

func getVideoId(avUrl string) string {
	return strings.Split(avUrl, "/")[3]
}
func getFullUrl(avUrl string) string {
	return "https://www.mgstage.com" + avUrl
}
func getFullUrlJoin(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return " "
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return " "
	}
	return baseUrl.ResolveReference(uri).String()
}
