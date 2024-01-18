package utils

import (
	"crypto/tls"
	"encoding/json"
	fmt "fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
	"strings"
	"time"
)

func WriteToFile(body []byte, name string) {
	fileName := fmt.Sprintf("json/%s", name)
	_ = os.Mkdir("json", os.ModePerm)
	_ = os.WriteFile(fileName, body, os.ModePerm)
}
func WriteToJSONByFileName(body []byte, name string) {
	fileName := fmt.Sprintf("json/%s.json", name)
	_ = os.Mkdir("json", os.ModePerm)
	_ = os.WriteFile(fileName, body, os.ModePerm)
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
		"User-Agent": `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36`,
		//"Authorization": "Bearer NjM0MzQ5MDIxMzg4OpTbOJlBQ8pQNsF078ecJkGtwTCb",
		"Cookie": "PHPSESSID=ciek7ijvn0vgedjbq2vak4e114; uuid=ce12aa6efbf35ba1eac11b1e69b279b4; coc=1; _ga=GA1.2.1193284767.1703415248; __ulfpc=202312241854081801; bWdzdGFnZS5jb20%3D-_lr_uf_-r2icil=81b78c30-df88-44fd-88dd-d4729b88b373; adc=1; _gid=GA1.2.1824065489.1705508051; _gat_UA-158726521-1=1; _gat_UA-58252858-1=1; bWdzdGFnZS5jb20%3D-_lr_tabs_-r2icil%2Fmgs={%22sessionID%22:0%2C%22recordingID%22:%225-eb6008d2-2a33-4d56-b256-1c9025f69e6c%22%2C%22webViewID%22:null%2C%22lastActivity%22:1705587724057}; bWdzdGFnZS5jb20%3D-_lr_hb_-r2icil%2Fmgs={%22heartbeat%22:1705587724057}; _ga_XGPRLSR61S=GS1.2.1705585804.8.1.1705587754.25.0.0; _ga_92ER0V7HV2=GS1.2.1705587719.12.1.1705587754.25.0.0",

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
		SetRetryWaitTime(10 * time.Second).
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
func CheckError(err error) {
	if err != nil {
		fmt.Printf("%d", err)
	}
}
func RequestThanSaveImage(url string, saveImagePath string) error {

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(10 * time.Second)
	resp, err := client.
		SetRetryCount(5).
		SetRetryWaitTime(10 * time.Second).
		SetDebug(false).
		R().Get(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		err = os.WriteFile(saveImagePath, resp.Body(), 0644) // Save image
		if err != nil {
			fmt.Println("Error saving image:", saveImagePath, url, err)
			return err
		}
	} else {
		fmt.Println(fmt.Sprintf("url is %s ,status code is %d %s", url, resp.StatusCode(), string(resp.Body())))
		RequestThanSaveImage(url, saveImagePath)
	}
	return nil
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

type Request struct {
	CategoryID int
	Cursor     string
	Sort       int
	IsVIP      bool
	Limit      int
}

func PostToStructInputStruct[T interface{}](url string, data interface{}, sessionId string) (t T, err error) {
	payload, err := json.Marshal(data)
	return PostToStructInputBytes[T](url, payload, sessionId)
}
func PostToStructInputBytes[T interface{}](url string, payloadBytes []byte, sessionId string) (t T, err error) {
	resp, err := GetHttpClient().R().
		//EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", "sessionid="+sessionId).
		SetBody(payloadBytes).
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
