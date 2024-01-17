package leetcodeCN

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"myCrawler/utils"
	"net/http"
	"os"
	"reflect"
	"strings"
)

func request(url string, query string) ([]byte, error) {
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	ShareHeader(req)
	return baseRequest(req)
}
func baseRequest(req *http.Request) ([]byte, error) {
	ShareHeader(req)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//fmt.Println(string(body))
	return body, nil
}
func graphql(postQL string) ([]byte, error) {
	url := "https://leetcode-cn.com/graphql"
	method := "POST"
	payload := strings.NewReader(postQL)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return baseRequest(req)
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("%d", err)
	}
}

func writeToJSON(body []byte, json interface{}) {
	typeOfA := reflect.TypeOf(json)
	//fmt.Println(typeOfA.Name(), typeOfA.Kind())
	fileName := fmt.Sprintf("json/%s.json", typeOfA.Name())
	_ = os.Mkdir("json", os.ModePerm)

	//jsonByte, _ := json.([]byte)
	_ = ioutil.WriteFile(fileName, body, os.ModePerm)
}

func getAllQuestionsTranslation() (QuestionCHNJSON, error) {
	body, _ := graphql(`{"operationName":"getQuestionTranslation","variables":{},"query":"query getQuestionTranslation($lang: String) {\n  translations: allAppliedQuestionTranslations(lang: $lang) {\n    title\n    questionId\n    __typename\n  }\n}\n"}`)
	var questionCHNJSON QuestionCHNJSON
	err := json.Unmarshal(body, &questionCHNJSON)
	checkError(err)
	writeToJSON(body, questionCHNJSON)
	return questionCHNJSON, err
}

func getAllQuestions() (AllQuestionJSON, error) {
	body, _ := request("https://leetcode-cn.com/api/problems/all/", "")
	var questionJSON AllQuestionJSON
	err := json.Unmarshal(body, &questionJSON)
	checkError(err)
	writeToJSON(body, questionJSON)
	return questionJSON, err
}
func getFavorites() (FavoriteJSON, error) {
	body, _ := request("https://leetcode-cn.com/problems/api/favorites/", "")
	var favoriteJSON FavoriteJSON
	err := json.Unmarshal(body, &favoriteJSON)
	checkError(err)
	writeToJSON(body, favoriteJSON)
	return favoriteJSON, err
}
func getTags() (TagJSON, error) {
	body, _ := request("https://leetcode-cn.com/problems/api/tags/", "")
	var tagJSON TagJSON
	err := json.Unmarshal(body, &tagJSON)
	checkError(err)
	writeToJSON(body, tagJSON)
	return tagJSON, err
}

func getCompanys() (CompanyJSON, error) {
	body, _ := graphql(`{"operationName":"getHotCompanies","variables":{"onlyCompanyCards":true},"query":"query getHotCompanies($onlyCompanyCards: Boolean!) {\n  interviewHotCards(onlyCompanyCards: $onlyCompanyCards) {\n    id\n    numQuestions\n    company {\n      name\n      slug\n      imgUrl\n      __typename\n    }\n    __typename\n  }\n}\n"}`)
	var companyJSON CompanyJSON
	err := json.Unmarshal(body, &companyJSON)
	checkError(err)
	writeToJSON(body, companyJSON)
	return companyJSON, err
}

func getCompnayQuestions() (CompanyQuestionJSON, error) {
	body, _ := graphql(`{"operationName":"companyTag","variables":{"slug":"bytedance"},"query":"query companyTag($slug: String!) {\n  interviewCard(companySlug: $slug) {\n    id\n    isFavorite\n    isPremiumOnly\n    privilegeExpiresAt\n    jobsCompany {\n      name\n      jobPostingNum\n      isVerified\n      description\n      logo\n      logoPath\n      postingTypeCounts {\n        count\n        postingType\n        __typename\n      }\n      industryDisplay\n      scaleDisplay\n      financingStageDisplay\n      website\n      legalName\n      __typename\n    }\n    __typename\n  }\n  interviewCompanyOptions(query: $slug) {\n    id\n    __typename\n  }\n  companyTag(slug: $slug) {\n    name\n    id\n    imgUrl\n    translatedName\n    frequencies\n    questions {\n      title\n      translatedTitle\n      titleSlug\n      questionId\n      stats\n      status\n      questionFrontendId\n      difficulty\n      frequencyTimePeriod\n      topicTags {\n        id\n        name\n        slug\n        translatedName\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  jobsCompany(companySlug: $slug) {\n    name\n    legalName\n    logo\n    description\n    website\n    industryDisplay\n    scaleDisplay\n    financingStageDisplay\n    isVerified\n    __typename\n  }\n}\n"}`)
	var companyQuestionJSON CompanyQuestionJSON
	//fmt.Println(string(body))
	err := json.Unmarshal(body, &companyQuestionJSON)
	checkError(err)
	writeToJSON(body, companyQuestionJSON)
	return companyQuestionJSON, err
}

func getQuestionDetail(titleSlug string) (QuestionJSON, error) {
	queryStatement := fmt.Sprintf("{\"operationName\":\"questionData\",\"variables\":{\"titleSlug\":\"%s\"},\"query\":\"query questionData($titleSlug: String!) {\\n  question(titleSlug: $titleSlug) {\\n    questionId\\n    questionFrontendId\\n    boundTopicId\\n    title\\n    titleSlug\\n    content\\n    translatedTitle\\n    translatedContent\\n    isPaidOnly\\n    difficulty\\n    likes\\n    dislikes\\n    isLiked\\n    similarQuestions\\n    contributors {\\n      username\\n      profileUrl\\n      avatarUrl\\n      __typename\\n    }\\n    langToValidPlayground\\n    topicTags {\\n      name\\n      slug\\n      translatedName\\n      __typename\\n    }\\n    companyTagStats\\n    codeSnippets {\\n      lang\\n      langSlug\\n      code\\n      __typename\\n    }\\n    stats\\n    hints\\n    solution {\\n      id\\n      canSeeDetail\\n      __typename\\n    }\\n    status\\n    sampleTestCase\\n    metaData\\n    judgerAvailable\\n    judgeType\\n    mysqlSchemas\\n    enableRunCode\\n    envInfo\\n    book {\\n      id\\n      bookName\\n      pressName\\n      source\\n      shortDescription\\n      fullDescription\\n      bookImgUrl\\n      pressImgUrl\\n      productUrl\\n      __typename\\n    }\\n    isSubscribed\\n    isDailyQuestion\\n    dailyRecordStatus\\n    editorType\\n    ugcQuestionId\\n    style\\n    __typename\\n  }\\n}\\n\"}", titleSlug)
	body, _ := graphql(queryStatement)
	fmt.Println(string(body))
	var questionJSON QuestionJSON
	err := json.Unmarshal(body, &questionJSON)
	checkError(err)
	utils.WriteToJSONByFileName(body, fmt.Sprintf("%s.json", titleSlug))
	return questionJSON, err
}

func ShareHeader(req *http.Request) {

	req.Header.Add("authority", "leetcode-cn.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("accept", "*/*")
	req.Header.Add("origin", "https://leetcode-cn.com")
	req.Header.Add("x-csrftoken", "EuuNpdOIWUDy1PwTBI5KseiWLl6NxQjpoWXTscj2b5MRlZaM86zTlemE09G1rxty")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4033.2 Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://leetcode-cn.com/problems/valid-permutations-for-di-sequence/")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Add("cookie", "_ga=GA1.2.1965944468.1574254532; __auc=fedaa78216e88e1f18109f26189; gr_user_id=4983e725-f1e4-4757-8252-11842415420b; _uab_collina=157425454000517146858749; grwng_uid=b9655aee-3902-4271-b016-d7bbb84ec0c2; a2873925c34ecbd2_gr_last_sent_cs1=gakaki; csrftoken=EuuNpdOIWUDy1PwTBI5KseiWLl6NxQjpoWXTscj2b5MRlZaM86zTlemE09G1rxty; LEETCODE_SESSION=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJfYXV0aF91c2VyX2lkIjoiMTAyNjczNiIsIl9hdXRoX3VzZXJfYmFja2VuZCI6ImF1dGhlbnRpY2F0aW9uLmF1dGhfYmFja2VuZHMuUGhvbmVBdXRoZW50aWNhdGlvbkJhY2tlbmQiLCJfYXV0aF91c2VyX2hhc2giOiI3ZWQyZTczMjFhNjk5NjNhMmI3MzY2MGMwN2JhNmU1NTczYWM2ZmMzIiwiaWQiOjEwMjY3MzYsImVtYWlsIjoiIiwidXNlcm5hbWUiOiJnYWtha2kiLCJ1c2VyX3NsdWciOiJnYWtha2kiLCJhdmF0YXIiOiJodHRwczovL2Fzc2V0cy5sZWV0Y29kZS1jbi5jb20vYWxpeXVuLWxjLXVwbG9hZC9kZWZhdWx0X2F2YXRhci5wbmciLCJwaG9uZV92ZXJpZmllZCI6dHJ1ZSwidGltZXN0YW1wIjoiMjAyMC0wMS0xNCAxNDo1MDoyMi42Nzg2NTcrMDA6MDAiLCJSRU1PVEVfQUREUiI6IjE3Mi4yMS41LjE5NSIsIklERU5USVRZIjoiNDhlYzEwNTIzMDZiYTkwOTRkZmRlOGU4NTVjN2QwZjgiLCJfc2Vzc2lvbl9leHBpcnkiOjEzODI0MDB9.mgw5KxqaNt5Oo_E3DQ-wYtajnUsTnjjVveV0-nNLIiw; _gid=GA1.2.122172361.1579752216; Hm_lvt_fa218a3ff7179639febdb15e372f411c=1580047689,1580053847,1580054570,1580102807; __asc=5967496f16fe577701dc6fea7ee; a2873925c34ecbd2_gr_session_id=5c56bbb2-baad-421b-ba73-7cffdccd35cc; a2873925c34ecbd2_gr_last_sent_sid_with_cs1=5c56bbb2-baad-421b-ba73-7cffdccd35cc; a2873925c34ecbd2_gr_session_id_5c56bbb2-baad-421b-ba73-7cffdccd35cc=true; Hm_lpvt_fa218a3ff7179639febdb15e372f411c=1580104888; a2873925c34ecbd2_gr_cs1=gakaki; _gat_gtag_UA_131851415_1=1; csrftoken=EuuNpdOIWUDy1PwTBI5KseiWLl6NxQjpoWXTscj2b5MRlZaM86zTlemE09G1rxty")

}
