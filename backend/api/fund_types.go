// API 请求参数和响应结构体定义
package api

import (
	"github.com/axiaoxin-com/investool/core"
	"github.com/axiaoxin-com/investool/datacenter/eastmoney"
	"github.com/axiaoxin-com/investool/models"
)

// FundIndexParams 基金首页请求参数
type FundIndexParams struct {
	PageNum  int    `json:"page_num"  form:"page_num" binding:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
	Sort     int    `json:"sort"      form:"sort"`
	Type     string `json:"type"      form:"type"`
}

// FundIndexResponse 基金首页响应
type FundIndexResponse struct {
	FundList      []*models.Fund     `json:"fund_list"`
	Pagination    PaginationResponse `json:"pagination"`
	UpdatedAt     string             `json:"updated_at"`
	AllFundCount  int                `json:"all_fund_count"`
	Fund4433Count int                `json:"fund_4433_count"`
	FundTypes     []string           `json:"fund_types"`
}

// FundFilterParams 基金筛选请求参数
type FundFilterParams struct {
	ParamFundListFilter models.ParamFundListFilter `json:"filter"`
	ParamFundIndex      FundIndexParams            `json:"index"`
}

// FundCheckParams 基金检测请求参数
type FundCheckParams struct {
	Code                 string              `json:"fundcode" binding:"required"`
	MinScale             float64             `json:"min_scale"`
	MaxScale             float64             `json:"max_scale"`
	MinManagerYears      float64             `json:"min_manager_years"`
	Year1RankRatio       float64             `json:"year_1_rank_ratio"`
	ThisYear235RankRatio float64             `json:"this_year_235_rank_ratio"`
	Month6RankRatio      float64             `json:"month_6_rank_ratio"`
	Month3RankRatio      float64             `json:"month_3_rank_ratio"`
	Max135AvgStddev      float64             `json:"max_135_avg_stddev"`
	Min135AvgSharp       float64             `json:"min_135_avg_sharp"`
	Max135AvgRetr        float64             `json:"max_135_avg_retr"`
	CheckStocks          bool                `json:"check_stocks"`
	StockCheckerOptions  core.CheckerOptions `json:"stock_checker_options"`
}

// FundCheckResponse 基金检测响应
type FundCheckResponse struct {
	Funds             []*models.Fund                        `json:"funds"`
	Param             FundCheckParams                       `json:"param"`
	StockCheckResults map[string]core.FundStocksCheckResult `json:"stock_check_results,omitempty"`
}

// FundManagerParams 基金经理筛选参数
type FundManagerParams struct {
	Name                string  `json:"name" form:"name"`
	MinWorkingYears     int     `json:"min_working_years" form:"min_working_years"`
	MinYieldse          float64 `json:"min_yieldse" form:"min_yieldse"`
	MaxCurrentFundCount int     `json:"max_current_fund_count" form:"max_current_fund_count"`
	MinScale            float64 `json:"min_scale" form:"min_scale"`
	PageNum             int     `json:"page_num" form:"page_num" binding:"min=1"`
	PageSize            int     `json:"page_size" form:"page_size" binding:"min=1,max=100"`
	Sort                string  `json:"sort" form:"sort"`
	FundType            string  `json:"fund_type" form:"fund_type"`
}

// FundManagerResponse 基金经理响应
type FundManagerResponse struct {
	Managers   []FundManagerInfo  `json:"managers"`
	Pagination PaginationResponse `json:"pagination"`
}

// FundManagerInfo 基金经理信息
type FundManagerInfo struct {
	eastmoney.FundManagerInfo
	BestFundIs4433 bool `json:"best_fund_is_4433"`
}

// FundSimilarityParams 基金相似度请求参数
type FundSimilarityParams struct {
	Codes string `json:"codes" form:"codes" binding:"required"`
}

// QueryByStockParams 股票选基请求参数
type QueryByStockParams struct {
	Keywords string `json:"keywords" form:"keywords" binding:"required"`
}
