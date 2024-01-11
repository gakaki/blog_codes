package main

import (
	"fmt"
	"testing"
)

func TestLeetCode(t *testing.T) {
	////获得所有的题目
	//allQuestionJSON, _ := getAllQuestions()
	////获得所有题目的中文名字
	//questionJSONTranslation , _ := getAllQuestionsTranslation()
	////所有比赛题目合集
	//favoritesJSON,_ := getFavorites()
	////所有分类合集
	//tagsJSON,_ := getTags()
	//
	//companyJSON,_ := getCompanys()
	////model数据重新组织
	////写入elasticsearch 和 dgraphql数据库
	//companyQuestionsJSON,_ := getCompnayQuestions()

	titleSlug := "valid-permutations-for-di-sequence"
	questionDetail, _ := getQuestionDetail(titleSlug)

	fmt.Println(questionDetail)

	//fmt.Println(
	//	len(allQuestionJSON.StatStatusPairs),
	//	len(questionJSONTranslation.Data.Translations),
	//	len(favoritesJSON),
	//	len(tagsJSON.Topics),
	//	len(companyJSON.Data.InterviewHotCards),
	//	len(companyQuestionsJSON.Data.CompanyTag.Questions),
	//)
}
