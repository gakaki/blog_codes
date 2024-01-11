package mgstage

type Link struct {
	Text string `json:"text"`
	Url string `json:"link"`
}
type Video struct {
	Url  string `json:"url"`
	FullUrl string `json:"url"`
	Price string `json:"price"`
	ID   string `json:"string"`
	Title string `json:"title"`

	CountFavorite string `json:"countFavorite"`
	CountPlay string `json:"countPlay"`

	Actor   Link `json:"actor"`
	Manufacturer Link `json:"manufacturer"`

	TimeLong  string `json:"timeLong"`
	StartDate string `json:"startDate"`
	SaleDate string `json:"saleDate"`

	Series    Link   `json:"series"`
	Company Link `json:"company"`

	Tags      []Link `json:"tags"`

	Rate string `json:"rate"`

	Description string `json:"description"`

	ImageHead string `json:"imageHead"`
	Images []string `json:"images"`

	Pid string `json:"pid"`
	VideoUrl string `json:"videoUrl"`

	//可以把评论也抓下来
}

