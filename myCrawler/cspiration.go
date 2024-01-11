package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	Cspiration "myCrawler/cspiration"
	"myCrawler/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var cspirationItems = make([]*Cspiration.Item, 0)

//var queueFethItems = make(chan *scanlibs.ScanLibItem, 3)
//var queueFethItems = make(chan int, 3)
//var wg sync.WaitGroup

func RequestCspirationGetDocument(url string) (*goquery.Document, error) {
	timeout := time.Duration(45 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "coc=1; PHPSESSID=q62v530bnrtfsqcfmnj2beuhp0; uuid=33e76d09fc86e676de2e0e8b887c8ee5; adc=1; __ulfpc=202012061038023125; _ga=GA1.2.98736236.1607222284; _gid=GA1.2.293888026.1607222284; bWdzdGFnZS5jb20%3D-_lr_uf_-r2icil=962de34a-89e7-413c-81cb-4d47111a2ac7; _gat_UA-58252858-1=1; _gat_UA-158726521-1=1; bWdzdGFnZS5jb20%3D-_lr_hb_-r2icil%2Fmgs={%22heartbeat%22:1607224305081}; bWdzdGFnZS5jb20%3D-_lr_tabs_-r2icil%2Fmgs={%22sessionID%22:0%2C%22recordingID%22:%224-e33ab609-ebc7-4322-a136-500d9641fe5f%22%2C%22lastActivity%22:1607224310786}//Cache-Control: no-cache; PHPSESSID=97n84d8jcorg6t45epdvbjvl02; uuid=33e76d09fc86e676de2e0e8b887c8ee5")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error", err)
		return nil, err
	}
	//函数结束后关闭相关链接
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		utils.WriteToFile(bodyBytes, "cspiration.html")
		//doc, err := goquery.NewDocumentFromReader(resp.Body)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyString))
		if err != nil {
			fmt.Println("error", err)
			log.Fatal(err)
			return nil, err
		}
		return doc, nil
	} else {
		fmt.Println("错误号码:")
		fmt.Println(resp.StatusCode)
		panic(fmt.Sprintf("status code is %2d", resp.StatusCode))
		return nil, err
	}
}

func GetFileBytes(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path) //read config file
	if err != nil {
		fmt.Println("read json file error")
		return nil, err
	}
	return data, nil
}
func dealPage() {
	//url := "https://cspiration.com/leetcodeClassification#nav-jieshao"
	//doc, err := RequestCspirationGetDocument(url)
	fileBytes, err := GetFileBytes("json/cspiration.html")
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(fileBytes)))
	if err != nil {
		fmt.Println("error", err)
		log.Fatal(err)
	}
	//fmt.Println(doc.Html())
	get400Questions(doc)
	get250CategoryQuestions(doc)
	getDSCategoryQuestions(doc)

}

var questionDS = make([]*Cspiration.Item, 0)
var question250 = make([]*Cspiration.Item, 0)
var question400 = make([]*Cspiration.Item, 0)
var questionMaps = make(map[string]*Cspiration.Item)

func getCommonCategoryQuestion(selector string, categoryName string, doc *goquery.Document, itemArray []*Cspiration.Item) []*Cspiration.Item {

	doc.Find(selector).Each(func(index int, ele *goquery.Selection) {

		questionNo := ele.Find("td:nth-child(1)").Text()
		questionNode := ele.Find("td:nth-child(2) > a")

		if questionNo != "" {

			item := questionMaps[questionNo]
			if item == nil {
				item = &Cspiration.Item{}
			}

			item.BigCategory = categoryName
			item.QuestionNo = questionNo
			item.Title = questionNode.Text()
			item.LeetCodeLink, _ = questionNode.Attr("href")
			questionMaps[questionNo] = item

			itemArray = append(itemArray, item)
		}
	})

	return itemArray
}

// Leetcode 前 400 重点 250 题
func get250CategoryQuestions(doc *goquery.Document) {
	selector := "#nav-zhongdiantimu > div > div.col-body > table > tbody > tr"
	categoryName := "Leetcode 前 400 重点 250 题"
	question250 = getCommonCategoryQuestion(selector, categoryName, doc, question250)
}

// Data Science Leetcode 精简版
func getDSCategoryQuestions(doc *goquery.Document) {
	selector := "#nav-DSzhongdiantimu > div > div.col-body > table > tbody > tr"
	categoryName := "Data Science Leetcode 精简版"
	questionDS = getCommonCategoryQuestion(selector, categoryName, doc, questionDS)
}

func get400Questions(doc *goquery.Document) {
	//print(doc)
	bigCategory := "LeetCode分类顺序类 400题"
	doc.Find("div.tab-pane").Each(func(index int, divEle *goquery.Selection) {
		divEleId, _ := divEle.Attr("id")
		if divEleId == "nav-question1" {
			return
		}
		if divEleId != "nav-jieshao" {
			category := divEle.Find("div > div.col-header").Text()
			subCategory := ""

			divEle.Find("div > div.col-body > table > tbody > tr").Each(func(index int, trEle *goquery.Selection) {
				firstText := trEle.Find("td:nth-child(1)").Text()
				questionNo, err := strconv.Atoi(firstText)
				if err != nil || questionNo == 0 {
					subCategory = firstText
				} else {
					//说明是questionNo
					questionNode := trEle.Find("td:nth-child(2) > a")

					item := questionMaps[firstText]
					if item == nil {
						item = &Cspiration.Item{}
						item.BigCategory = bigCategory
						item.Category = category
						item.SubCategory = subCategory
						item.QuestionNo = strconv.Itoa(questionNo)
						item.Title = questionNode.Text()
						item.LeetCodeLink, _ = questionNode.Attr("href")
						item.Desc = trEle.Find("td:nth-child(4)").Text()
						questionMaps[item.QuestionNo] = item
						question400 = append(question400, item)
					}

				}
			})
		}
	})
}
