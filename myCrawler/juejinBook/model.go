package juejinBook

import "regexp"

type Config struct {
	Sessionid string   `yaml:"sessionid"`
	BookIDs   []string `yaml:"book_ids"`
	SaveDir   string   `yaml:"save_dir"`
}

type Juejinxiaoce2Markdown struct {
	ImgPattern        *regexp.Regexp
	Sessionid         string
	BookIDs           []string
	SaveDir           string
	RequestHeaders    map[string]string
	MarkdownSavePaths map[string]string
}
type JuejinBook struct {
	BookID int64
	Title  string
	Cover  string
}

type JuejinBookDownloadConfig struct {
	SessionID         string
	Books             string
	BookIDs           []string
	SaveDir           string
	RequestHeaders    map[string]string
	MarkdownSavePaths map[int64]string
}
type JuejinSection struct {
	Data struct {
		Booklet struct {
			BaseInfo struct {
				Title string `json:"title"`
			} `json:"base_info"`
		} `json:"booklet"`
		Sections []struct {
			SectionID string `json:"section_id"`
		} `json:"sections"`
	} `json:"data"`
}

type JuejinSectionContent struct {
	Data struct {
		Section struct {
			Title        string `json:"title"`
			MarkdownShow string `json:"markdown_show"`
		} `json:"section"`
	} `json:"data"`
}

type JuejinResponseBook struct {
	BookletID string `json:"booklet_id"`
	BaseInfo  struct {
		ID                int    `json:"id"`
		BookletID         string `json:"booklet_id"`
		Title             string `json:"title"`
		Price             int    `json:"price"`
		CategoryID        string `json:"category_id"`
		Status            int    `json:"status"`
		UserID            string `json:"user_id"`
		VerifyStatus      int    `json:"verify_status"`
		Summary           string `json:"summary"`
		CoverImg          string `json:"cover_img"`
		SectionCount      int    `json:"section_count"`
		SectionIds        string `json:"section_ids"`
		IsFinished        int    `json:"is_finished"`
		Ctime             int    `json:"ctime"`
		Mtime             int    `json:"mtime"`
		PutOnTime         int    `json:"put_on_time"`
		PullOffTime       int64  `json:"pull_off_time"`
		FinishedTime      int64  `json:"finished_time"`
		RecycleBinTime    int64  `json:"recycle_bin_time"`
		VerifyTime        int64  `json:"verify_time"`
		SubmitTime        int    `json:"submit_time"`
		TopTime           int    `json:"top_time"`
		WechatGroupImg    string `json:"wechat_group_img"`
		WechatGroupDesc   string `json:"wechat_group_desc"`
		WechatGroupSignal string `json:"wechat_group_signal"`
		ReadTime          int    `json:"read_time"`
		BuyCount          int    `json:"buy_count"`
		CourseType        int    `json:"course_type"`
		BackgroundImg     string `json:"background_img"`
		IsDistribution    int    `json:"is_distribution"`
		DistributionImg   string `json:"distribution_img"`
		Commission        int    `json:"commission"`
		CanVipBorrow      bool   `json:"can_vip_borrow"`
	} `json:"base_info"`
	UserInfo struct {
		UserID            string `json:"user_id"`
		UserName          string `json:"user_name"`
		Company           string `json:"company"`
		JobTitle          string `json:"job_title"`
		AvatarLarge       string `json:"avatar_large"`
		Level             int    `json:"level"`
		Description       string `json:"description"`
		FolloweeCount     int    `json:"followee_count"`
		FollowerCount     int    `json:"follower_count"`
		PostArticleCount  int    `json:"post_article_count"`
		DiggArticleCount  int    `json:"digg_article_count"`
		GotDiggCount      int    `json:"got_digg_count"`
		GotViewCount      int    `json:"got_view_count"`
		PostShortmsgCount int    `json:"post_shortmsg_count"`
		DiggShortmsgCount int    `json:"digg_shortmsg_count"`
		Isfollowed        bool   `json:"isfollowed"`
		FavorableAuthor   int    `json:"favorable_author"`
		Power             int    `json:"power"`
		StudyPoint        int    `json:"study_point"`
		University        struct {
			UniversityID string `json:"university_id"`
			Name         string `json:"name"`
			Logo         string `json:"logo"`
		} `json:"university"`
		Major struct {
			MajorID  string `json:"major_id"`
			ParentID string `json:"parent_id"`
			Name     string `json:"name"`
		} `json:"major"`
		StudentStatus           int  `json:"student_status"`
		SelectEventCount        int  `json:"select_event_count"`
		SelectOnlineCourseCount int  `json:"select_online_course_count"`
		Identity                int  `json:"identity"`
		IsSelectAnnual          bool `json:"is_select_annual"`
		SelectAnnualRank        int  `json:"select_annual_rank"`
		AnnualListType          int  `json:"annual_list_type"`
		ExtraMap                struct {
		} `json:"extraMap"`
		IsLogout       int   `json:"is_logout"`
		AnnualInfo     []any `json:"annual_info"`
		AccountAmount  int   `json:"account_amount"`
		UserGrowthInfo struct {
			UserID                   int64   `json:"user_id"`
			Jpower                   int     `json:"jpower"`
			Jscore                   float64 `json:"jscore"`
			JpowerLevel              int     `json:"jpower_level"`
			JscoreLevel              int     `json:"jscore_level"`
			JscoreTitle              string  `json:"jscore_title"`
			AuthorAchievementList    []int   `json:"author_achievement_list"`
			VipLevel                 int     `json:"vip_level"`
			VipTitle                 string  `json:"vip_title"`
			JscoreNextLevelScore     int     `json:"jscore_next_level_score"`
			JscoreThisLevelMiniScore int     `json:"jscore_this_level_mini_score"`
			VipScore                 int     `json:"vip_score"`
		} `json:"user_growth_info"`
		IsVip bool `json:"is_vip"`
	} `json:"user_info"`
	IsBuy               bool `json:"is_buy"`
	SectionUpdatedCount int  `json:"section_updated_count"`
	IsNew               bool `json:"is_new"`
	MaxDiscount         struct {
		DiscountType  int    `json:"discount_type"`
		DiscountID    string `json:"discount_id"`
		CouponID      string `json:"coupon_id"`
		CouponBasicID string `json:"coupon_basic_id"`
		Name          string `json:"name"`
		Desc          string `json:"desc"`
		DiscountRate  int    `json:"discount_rate"`
		Price         int    `json:"price"`
		DiscountMoney int    `json:"discount_money"`
		PayMoney      int    `json:"pay_money"`
		IsLimitedTime int    `json:"is_limited_time"`
		StartTime     int    `json:"start_time"`
		EndTime       int    `json:"end_time"`
	} `json:"max_discount"`
}

type JuejinResponse struct {
	ErrNo   int                  `json:"err_no"`
	ErrMsg  string               `json:"err_msg"`
	Data    []JuejinResponseBook `json:"data"`
	Cursor  string               `json:"cursor"`
	Count   int                  `json:"count"`
	HasMore bool                 `json:"has_more"`
}
