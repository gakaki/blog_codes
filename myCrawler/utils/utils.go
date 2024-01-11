package utils

import (
	"crypto/tls"
	"encoding/json"
	fmt "fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func WriteToFile(body []byte, name string) {
	fileName := fmt.Sprintf("json/%s", name)
	_ = os.Mkdir("json", os.ModePerm)
	_ = ioutil.WriteFile(fileName, body, os.ModePerm)
}
func WriteToJSONByFileName(body []byte, name string) {
	fileName := fmt.Sprintf("json/%s.json", name)
	_ = os.Mkdir("json", os.ModePerm)
	_ = ioutil.WriteFile(fileName, body, os.ModePerm)
}
func WriteJSON(j interface{}, fileName string) {
	f, err := json.MarshalIndent(&j, "", " ")

	if err != nil {
		fmt.Println(err)
	} else {
		os.WriteFile(fmt.Sprintf("%s", fileName), f, 0777)
	}
}
func ReadJSONBytes(fileName string) []byte {
	fileRead, _ := os.ReadFile(fmt.Sprintf("%s", fileName))
	return fileRead
}
func ReadJSON[T interface{}](fileName string) (t T) {
	json.Unmarshal(ReadJSONBytes(fileName), &t)
	return t
}
func GetCommonHeaders() map[string]string {
	return map[string]string{
		//"pragma":        "no-cache",
		//"cache-control": "no-cache",
		"User-Agent": `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36`,
		//"Authorization": "Bearer NjM0MzQ5MDIxMzg4OpTbOJlBQ8pQNsF078ecJkGtwTCb",
		"Cookie": "uuid=67b2a73cbeeac80f3cc1f7538299bc98; _ga=GA1.2.1416779059.1704272338; _ga_XGPRLSR61S=GS1.2.1704982317.8.1.1704982452.13.0.0; _ga_92ER0V7HV2=GS1.2.1704982317.8.1.1704982452.13.0.0; __ulfpc=202401031658581756; adc=1; bWdzdGFnZS5jb20%3D-_lr_tabs_-r2icil%2Fmgs={%22sessionID%22:0%2C%22recordingID%22:%225-b671bf91-b2d2-4e1d-b0d9-236f2c560758%22%2C%22webViewID%22:null%2C%22lastActivity%22:1704982452191}; bWdzdGFnZS5jb20%3D-_lr_hb_-r2icil%2Fmgs={%22heartbeat%22:1704982462612}; _gid=GA1.2.498118497.1704895794; PHPSESSID=hal6k2uoeqv6102li0teekjn51; coc=1; bWdzdGFnZS5jb20%3D-_lr_uf_-r2icil=ab59e25f-29fd-4ae0-9c1d-dc67b430d3d9",

		//"accept-encoding": "gzip, deflate, br",
		//"Content-Type":    "application/json",
	}
}

func GetHttpClient() *resty.Client {

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(10 * time.Second)

	client.SetHeaders(GetCommonHeaders())
	client.SetContentLength(true)

	client.SetProxy("socks5://127.0.0.1:7890")

	client.
		SetRetryCount(5).
		SetRetryWaitTime(4 * time.Second).
		SetDebug(false)

	return client
}

func RequestString(url string) (string, error) {
	resp, err := GetHttpClient().R().
		//EnableTrace().
		Get(url)

	if err != nil {
		return "", err
	}
	//fmt.Println("Response Info:")
	//fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode())

	if resp.StatusCode() == http.StatusOK {
		bodyString := string(resp.Body())
		return bodyString, nil
	} else {
		fmt.Println("错误号码:")
		panic(fmt.Sprintf("url is %s ,status code is %d %s", url, resp.StatusCode(), string(resp.Body())))
		return "", err
	}
}

func RequestGetDocument(url string) (*goquery.Document, error) {
	body, err := RequestString(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		fmt.Println("error", err)
		return nil, err
	}
	return doc, nil
}
