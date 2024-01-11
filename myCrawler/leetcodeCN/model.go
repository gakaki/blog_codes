package leetcodeCN

type QuestionJSON struct {
	Data struct {
		Question struct {
			QuestionID            string        `json:"questionId"`
			QuestionFrontendID    string        `json:"questionFrontendId"`
			BoundTopicID          int           `json:"boundTopicId"`
			Title                 string        `json:"title"`
			TitleSlug             string        `json:"titleSlug"`
			Content               string        `json:"content"`
			TranslatedTitle       string        `json:"translatedTitle"`
			TranslatedContent     string        `json:"translatedContent"`
			IsPaidOnly            bool          `json:"isPaidOnly"`
			Difficulty            string        `json:"difficulty"`
			Likes                 int           `json:"likes"`
			Dislikes              int           `json:"dislikes"`
			IsLiked               interface{}   `json:"isLiked"`
			SimilarQuestions      string        `json:"similarQuestions"`
			Contributors          []interface{} `json:"contributors"`
			LangToValidPlayground string        `json:"langToValidPlayground"`
			TopicTags             []struct {
				Name           string `json:"name"`
				Slug           string `json:"slug"`
				TranslatedName string `json:"translatedName"`
				Typename       string `json:"__typename"`
			} `json:"topicTags"`
			CompanyTagStats interface{} `json:"companyTagStats"`
			CodeSnippets    []struct {
				Lang     string `json:"lang"`
				LangSlug string `json:"langSlug"`
				Code     string `json:"code"`
				Typename string `json:"__typename"`
			} `json:"codeSnippets"`
			Stats           string        `json:"stats"`
			Hints           []interface{} `json:"hints"`
			Solution        interface{}   `json:"solution"`
			Status          interface{}   `json:"status"`
			SampleTestCase  string        `json:"sampleTestCase"`
			MetaData        string        `json:"metaData"`
			JudgerAvailable bool          `json:"judgerAvailable"`
			JudgeType       string        `json:"judgeType"`
			MysqlSchemas    []interface{} `json:"mysqlSchemas"`
			EnableRunCode   bool          `json:"enableRunCode"`
			EnvInfo         string        `json:"envInfo"`
			Book            interface{}   `json:"book"`
			IsSubscribed    bool          `json:"isSubscribed"`
			Typename        string        `json:"__typename"`
		} `json:"question"`
	} `json:"data"`
}

type AllQuestionJSON struct {
	UserName        string `json:"user_name"`
	NumSolved       int    `json:"num_solved"`
	NumTotal        int    `json:"num_total"`
	AcEasy          int    `json:"ac_easy"`
	AcMedium        int    `json:"ac_medium"`
	AcHard          int    `json:"ac_hard"`
	StatStatusPairs []struct {
		Stat struct {
			QuestionID          int    `json:"question_id"`
			QuestionTitle       string `json:"question__title"`
			QuestionTitleSlug   string `json:"question__title_slug"`
			QuestionHide        bool   `json:"question__hide"`
			TotalAcs            int    `json:"total_acs"`
			TotalSubmitted      int    `json:"total_submitted"`
			TotalColumnArticles int    `json:"total_column_articles"`
			FrontendQuestionID  string `json:"frontend_question_id"`
			IsNewQuestion       bool   `json:"is_new_question"`
		} `json:"stat"`
		Status     interface{} `json:"status"`
		Difficulty struct {
			Level int `json:"level"`
		} `json:"difficulty"`
		PaidOnly  bool `json:"paid_only"`
		IsFavor   bool `json:"is_favor"`
		Frequency int  `json:"frequency"`
		Progress  int  `json:"progress"`
	} `json:"stat_status_pairs"`
	FrequencyHigh int    `json:"frequency_high"`
	FrequencyMid  int    `json:"frequency_mid"`
	CategorySlug  string `json:"category_slug"`
}

type QuestionCHNJSON struct {
	Data struct {
		Translations []struct {
			QuestionID string `json:"questionId"`
			Title      string `json:"title"`
			Typename   string `json:"__typename"`
		} `json:"translations"`
	} `json:"data"`
}

type FavoriteJSON []struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Questions []int  `json:"questions"`
	Type      string `json:"type"`
}

type TagJSON struct {
	Companies []interface{} `json:"companies"`
	Topics    []struct {
		Slug           string `json:"slug"`
		Name           string `json:"name"`
		Questions      []int  `json:"questions"`
		TranslatedName string `json:"translatedName"`
	} `json:"topics"`
}

type CompanyJSON struct {
	Data struct {
		InterviewHotCards []struct {
			ID           string `json:"id"`
			NumQuestions int    `json:"numQuestions"`
			Company      struct {
				Name     string `json:"name"`
				Slug     string `json:"slug"`
				ImgURL   string `json:"imgUrl"`
				Typename string `json:"__typename"`
			} `json:"company"`
			Typename string `json:"__typename"`
		} `json:"interviewHotCards"`
	} `json:"data"`
}

type CompanyQuestionJSON struct {
	Data struct {
		InterviewCard struct {
			ID                 string      `json:"id"`
			IsFavorite         bool        `json:"isFavorite"`
			IsPremiumOnly      bool        `json:"isPremiumOnly"`
			PrivilegeExpiresAt interface{} `json:"privilegeExpiresAt"`
			JobsCompany        struct {
				Name                  string        `json:"name"`
				JobPostingNum         int           `json:"jobPostingNum"`
				IsVerified            bool          `json:"isVerified"`
				Description           string        `json:"description"`
				Logo                  string        `json:"logo"`
				LogoPath              string        `json:"logoPath"`
				PostingTypeCounts     []interface{} `json:"postingTypeCounts"`
				IndustryDisplay       string        `json:"industryDisplay"`
				ScaleDisplay          string        `json:"scaleDisplay"`
				FinancingStageDisplay string        `json:"financingStageDisplay"`
				Website               string        `json:"website"`
				LegalName             string        `json:"legalName"`
				Typename              string        `json:"__typename"`
			} `json:"jobsCompany"`
			Typename string `json:"__typename"`
		} `json:"interviewCard"`
		InterviewCompanyOptions []struct {
			ID       int    `json:"id"`
			Typename string `json:"__typename"`
		} `json:"interviewCompanyOptions"`
		CompanyTag struct {
			Name           string      `json:"name"`
			ID             string      `json:"id"`
			ImgURL         string      `json:"imgUrl"`
			TranslatedName interface{} `json:"translatedName"`
			Frequencies    string      `json:"frequencies"`
			Questions      []struct {
				QuestionID          string      `json:"questionId"`
				TitleSlug           string      `json:"titleSlug"`
				QuestionFrontendID  string      `json:"questionFrontendId"`
				Status              interface{} `json:"status"`
				Title               string      `json:"title"`
				TranslatedTitle     string      `json:"translatedTitle"`
				Difficulty          string      `json:"difficulty"`
				Stats               string      `json:"stats"`
				IsPaidOnly          bool        `json:"isPaidOnly"`
				FrequencyTimePeriod interface{} `json:"frequencyTimePeriod"`
				TopicTags           []struct {
					ID             string `json:"id"`
					Name           string `json:"name"`
					TranslatedName string `json:"translatedName"`
					Slug           string `json:"slug"`
					Typename       string `json:"__typename"`
				} `json:"topicTags"`
				Typename string `json:"__typename"`
			} `json:"questions"`
			Typename string `json:"__typename"`
		} `json:"companyTag"`
		JobsCompany struct {
			Name                  string `json:"name"`
			LegalName             string `json:"legalName"`
			Logo                  string `json:"logo"`
			Description           string `json:"description"`
			Website               string `json:"website"`
			IndustryDisplay       string `json:"industryDisplay"`
			ScaleDisplay          string `json:"scaleDisplay"`
			FinancingStageDisplay string `json:"financingStageDisplay"`
			IsVerified            bool   `json:"isVerified"`
			Typename              string `json:"__typename"`
		} `json:"jobsCompany"`
	} `json:"data"`
}
