package cspiration

/*
{
  "id": "",
  "bigCategory":"leetcode400题",
  "category": "Array",
  "subCategory": "基础",
  "desc": "很少考",
  "questionNo": "27",
  "title": "find a number",
  "leetCodeLink": "http://leetcode",
}
*/

type Item struct {
	ID           string `json:"id"`
	BigCategory  string `json:"bigCategory"`
	Category     string `json:"category"`
	SubCategory  string `json:"subCategory"`
	Desc         string `json:"desc"`
	QuestionNo   string    `json:"questionNo"`
	Title        string `json:"title"`
	LeetCodeLink string `json:"leetCodeLink"`
}