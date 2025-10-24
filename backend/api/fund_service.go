// API 服务层 - 处理业务逻辑
package api

import (
	"context"
	"sync"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/investool/core"
	"github.com/axiaoxin-com/investool/datacenter/eastmoney"
	"github.com/axiaoxin-com/investool/models"
	"github.com/axiaoxin-com/logging"
)

// FundService 基金服务
type FundService struct{}

// NewFundService 创建基金服务实例
func NewFundService() *FundService {
	return &FundService{}
}

// GetFundIndex 获取4433基金列表
func (s *FundService) GetFundIndex(ctx context.Context, params FundIndexParams) (*FundIndexResponse, error) {
	fundList := models.Fund4433List

	// 过滤
	if params.Type != "" {
		fundList = fundList.FilterByType(params.Type)
	}

	// 排序
	if params.Sort > 0 {
		fundList.Sort(models.FundSortType(params.Sort))
	}

	// 分页
	totalCount := len(fundList)
	pagi := goutils.PaginateByPageNumSize(totalCount, params.PageNum, params.PageSize)
	result := fundList[pagi.StartIndex:pagi.EndIndex]

	return &FundIndexResponse{
		FundList: result,
		Pagination: PaginationResponse{
			PageNum:    pagi.PageNum,
			PageSize:   pagi.PageSize,
			Total:      len(fundList),
			TotalPages: len(fundList) / pagi.PageSize,
			StartIndex: pagi.StartIndex,
			EndIndex:   pagi.EndIndex,
		},
		UpdatedAt:     models.SyncFundTime.Format("2006-01-02 15:04:05"),
		AllFundCount:  len(models.FundAllList),
		Fund4433Count: totalCount,
		FundTypes:     models.Fund4433TypeList,
	}, nil
}

// GetFundFilter 基金筛选
func (s *FundService) GetFundFilter(ctx context.Context, params FundFilterParams) (*FundIndexResponse, error) {
	fundList := models.FundAllList.Filter(ctx, params.ParamFundListFilter)
	fundTypes := fundList.Types()

	// 过滤
	if params.ParamFundIndex.Type != "" {
		fundList = fundList.FilterByType(params.ParamFundIndex.Type)
	}

	// 排序
	if params.ParamFundIndex.Sort > 0 {
		fundList.Sort(models.FundSortType(params.ParamFundIndex.Sort))
	}

	// 分页
	pagi := goutils.PaginateByPageNumSize(len(fundList), params.ParamFundIndex.PageNum, params.ParamFundIndex.PageSize)
	result := fundList[pagi.StartIndex:pagi.EndIndex]

	return &FundIndexResponse{
		FundList: result,
		Pagination: PaginationResponse{
			PageNum:    pagi.PageNum,
			PageSize:   pagi.PageSize,
			Total:      len(fundList),
			TotalPages: len(fundList) / pagi.PageSize,
			StartIndex: pagi.StartIndex,
			EndIndex:   pagi.EndIndex,
		},
		UpdatedAt:     models.SyncFundTime.Format("2006-01-02 15:04:05"),
		AllFundCount:  len(models.FundAllList),
		Fund4433Count: len(fundList),
		FundTypes:     fundTypes,
	}, nil
}

// CheckFund 基金检测
func (s *FundService) CheckFund(ctx context.Context, params FundCheckParams) (*FundCheckResponse, error) {
	if params.Code == "" {
		return nil, ErrFundCodeRequired
	}

	codes := goutils.SplitStringFields(params.Code)
	searcher := core.NewSearcher(ctx)
	fundsMap, err := searcher.SearchFunds(ctx, codes)
	if err != nil {
		return nil, err
	}

	funds := []*models.Fund{}
	for _, fund := range fundsMap {
		funds = append(funds, fund)
	}

	response := &FundCheckResponse{
		Funds: funds,
		Param: params,
	}

	// 如果需要检测持仓个股
	if params.CheckStocks {
		if len(funds) > 50 {
			return nil, ErrTooManyFunds
		}

		stockCheckResults := map[string]core.FundStocksCheckResult{}
		checker := core.NewChecker(ctx, params.StockCheckerOptions)
		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, fund := range funds {
			wg.Add(1)
			go func(fund *models.Fund) {
				defer wg.Done()
				checkResult, err := checker.CheckFundStocks(ctx, fund)
				if err != nil {
					logging.Errorf(ctx, "CheckFundStocks code:%s err:%v", fund.Code, err)
					return
				}
				mu.Lock()
				stockCheckResults[fund.Code] = checkResult
				mu.Unlock()
			}(fund)
		}
		wg.Wait()

		response.StockCheckResults = stockCheckResults
	}

	return response, nil
}

// GetFundManagers 基金经理筛选
func (s *FundService) GetFundManagers(ctx context.Context, params FundManagerParams) (*FundManagerResponse, error) {
	// 筛选
	managers := models.FundManagers.Filter(ctx, eastmoney.ParamFundManagerFilter{
		MinWorkingYears:     params.MinWorkingYears,
		MinYieldse:          params.MinYieldse,
		MaxCurrentFundCount: params.MaxCurrentFundCount,
		MinScale:            params.MinScale,
		FundType:            params.FundType,
	})

	// 排序
	switch params.Sort {
	case "yieldse":
		managers.SortByYieldse()
	case "scale":
		managers.SortByScale()
	case "score":
		managers.SortByScore()
	case "an":
		managers.SortByAwardNum()
	case "fc":
		managers.SortByFundCount()
	case "cbr":
		managers.SortByCurrentBestReturn()
	case "wbr":
		managers.SortByCurrentBestReturn()
	}

	// 分页
	pagi := goutils.PaginateByPageNumSize(len(managers), params.PageNum, params.PageSize)
	managers = managers[pagi.StartIndex:pagi.EndIndex]

	// 获取这批基金经理的代表基金是否是4433基金
	bestFundCodes := []string{}
	for _, m := range managers {
		bestFundCodes = append(bestFundCodes, m.CurrentBestFundCode)
	}
	searcher := core.NewSearcher(ctx)
	bestFundInfoMap, err := searcher.SearchFunds(ctx, bestFundCodes)
	if err != nil {
		logging.Error(ctx, "SearchFunds err:"+err.Error())
	}

	// 返回结果item
	type managerInfo struct {
		eastmoney.FundManagerInfo
		BestFundIs4433 bool
	}
	result := []managerInfo{}
	for _, m := range managers {
		i := bestFundInfoMap[m.CurrentBestFundCode]
		r := managerInfo{
			FundManagerInfo: *m,
			BestFundIs4433:  i.Is4433(ctx),
		}
		result = append(result, r)
	}

	return &FundManagerResponse{
		Managers: []FundManagerInfo{},
		Pagination: PaginationResponse{
			PageNum:    pagi.PageNum,
			PageSize:   pagi.PageSize,
			Total:      len(managers),
			TotalPages: len(managers) / pagi.PageSize,
			StartIndex: pagi.StartIndex,
			EndIndex:   pagi.EndIndex,
		},
	}, nil
}

// GetFundSimilarity 基金持仓相似度
func (s *FundService) GetFundSimilarity(ctx context.Context, params FundSimilarityParams) (interface{}, error) {
	if params.Codes == "" {
		return nil, ErrFundCodesRequired
	}

	codeList := goutils.SplitStringFields(params.Codes)
	checker := core.NewChecker(ctx, core.DefaultCheckerOptions)
	result, err := checker.GetFundStocksSimilarity(ctx, codeList)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// QueryByStock 股票选基
func (s *FundService) QueryByStock(ctx context.Context, params QueryByStockParams) (interface{}, error) {
	if params.Keywords == "" {
		return nil, ErrKeywordsRequired
	}

	// 这里需要实现股票选基的业务逻辑
	// 暂时返回模拟数据
	return map[string]interface{}{
		"fund_count":  0,
		"stock_count": 0,
		"message":     "功能开发中",
	}, nil
}
