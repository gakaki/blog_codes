package scanlibs

/*
{
  "id": 1,
  "url": "https://scanlibs.com/mastering-opencv-python-practical-processing/",
  "title": "Mastering OpenCV 4 with Python: A practical guide covering topics from image processing, augmented reality to deep learning with OpenCV 4 and Python 3.7",
  "createdAt": "December 13, 2020",
  "author": "",
  "pubDate": "2019",
  "isVideo": false,
  "isbn": "978-123-145",
  "language": "english",
  "format": "pdf/ebub",
  "formats": [
    "pdf",
    "epub"
  ],
  "size": "533Mb",
  "thumbnil": "",
  "tags": [
    "books"
  ],
  "review_url": "http://amazon",
  "pages": "532",
  "scanlibs_url": "http://amazon",
  "description": "balabala"
}
*/

type ScanLibItem struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	SubTitle    string   `json:"subTitle"`
	CreatedAt   string   `json:"createdAt"`
	IsVideo     bool     `json:"isVideo"`
	Author      string   `json:"author"`
	PubDate     string   `json:"pubDate"`
	Isbn        string   `json:"isbn"`
	Language    string   `json:"language"`
	Format      string   `json:"format"`
	Formats     []string `json:"formats"`
	Size        string   `json:"size"`
	Thumbnil    string   `json:"thumbnil"`
	Tags        []string `json:"tags"`
	ReviewURL   string   `json:"review_url"`
	Pages       string   `json:"pages"`
	ScanlibsURL string   `json:"scanlibs_url"`
	Description string   `json:"description"`
}
type ScanLibPage struct {
	Index int            `json:"index"`
	URL   string         `json:"url"`
	Items []*ScanLibItem `json:"scanlib_items"`
}
