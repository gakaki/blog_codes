package Satiger

// {
// "id": "1",
// "url": "https://scanlibs.com/mastering-opencv-python-practical-processing/",
// "yearmonth": "2020年十二月",
// "title": "Accelerating Cloud Adoption",
// "thumbnil": "",
// "pubDate": "2019",
// "officalUrl": "2019",
// "officalPress": "orreilly",
// "baiduUrl": "orreilly",
// "baiduCode": "orreilly",
// "description": "balabala",
// "createdAt": "balabala",
// "tags": [
// "books"
// ]
// }
type Tag struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type SatigerItem struct {
	ID            string   `json:"id"`
	URL           string   `json:"url"`
	Yearmonth     string   `json:"yearmonth"`
	Title         string   `json:"title"`
	Thumbnil      string   `json:"thumbnil"`
	PubDate       string   `json:"pubDate"`
	OfficalURL    string   `json:"officalUrl"`
	OfficalPress  string   `json:"officalPress"`
	BaiduURL      string   `json:"baiduUrl"`
	OtherLinks    []string `json:"otherLinks"`
	BaiduCode     string   `json:"baiduCode"`
	Description   string   `json:"description"`
	CreatedAt     string   `json:"createdAt"`
	Tags          []Tag    `json:"tags"`
	ZlibSearchUrl string   `json:"zlibSearchUrl"`
}
