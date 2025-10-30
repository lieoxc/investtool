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
	"github.com/sirupsen/logrus"
)

// FundService 基金服务
type FundService struct{}

// NewFundService 创建基金服务实例
func NewFundService() *FundService {
	return &FundService{}
}

// GetFundIndex 获取4433基金列表
func (s *FundService) GetFundIndex(ctx context.Context, params FundIndexParams) (*FundIndexResponse, error) {
	logrus.Info("GetFundIndex params:", params)
	// 检查数据库是否初始化
	if models.DB == nil {
		logging.Error(ctx, "database not initialized")
		return &FundIndexResponse{}, nil
	}
	var allFound int64
	if err := models.DB.Model(&models.FundDB{}).Count(&allFound).Error; err != nil {
		logrus.Error("GetFundIndex get all found error:" + err.Error())
		return &FundIndexResponse{}, nil
	}

	// 从数据库获取4433基金
	var fundDBs []models.FundDB
	query := models.DB.Model(&models.FundDB{}).Where("is_4433 = ?", true)

	// 如果指定了基金类型，添加类型过滤
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	// 获取总数
	var fourFourThreeCount int64
	if err := query.Count(&fourFourThreeCount).Error; err != nil {
		logrus.Error("GetFundIndex get four four three count error:" + err.Error())
		return &FundIndexResponse{}, nil
	}
	// 分页
	pagi := goutils.PaginateByPageNumSize(int(fourFourThreeCount), params.PageNum, params.PageSize)

	// 重新构建查询（因为 Count 已经修改了 query）
	query = models.DB.Model(&models.FundDB{}).Where("is_4433 = ?", true)
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	// 这里需要根据sort参数进行排序
	orderClause := getOrderClause(models.FundSortType(params.Sort))
	if err := query.Order(orderClause).Offset(pagi.StartIndex).Limit(pagi.PageSize).Find(&fundDBs).Error; err != nil {
		return nil, err
	}

	// 转换为基金列表
	fundList := make([]*models.Fund, len(fundDBs))
	for i, fd := range fundDBs {
		fundList[i] = fd.ToFund()
	}

	// 获取所有基金类型
	var fundTypes []string
	if err := models.DB.Model(&models.FundDB{}).
		Where("is_4433 = ?", true).
		Distinct("type").
		Pluck("type", &fundTypes).Error; err != nil {
		logrus.Error("GetFundIndex get fund types error:" + err.Error())
	}

	// 计算总页数，避免除零
	totalPages := 0
	if pagi.PageSize > 0 {
		totalPages = int(fourFourThreeCount) / pagi.PageSize
		if int(fourFourThreeCount)%pagi.PageSize > 0 {
			totalPages++
		}
	}

	return &FundIndexResponse{
		FundList: fundList,
		Pagination: PaginationResponse{
			PageNum:    pagi.PageNum,
			PageSize:   pagi.PageSize,
			Total:      int(fourFourThreeCount),
			TotalPages: totalPages,
			StartIndex: pagi.StartIndex,
			EndIndex:   pagi.EndIndex,
		},
		UpdatedAt:     models.SyncFundTime.Format("2006-01-02 15:04:05"),
		AllFundCount:  int(allFound),
		Fund4433Count: int(fourFourThreeCount),
		FundTypes:     fundTypes,
	}, nil
}

// getOrderClause 根据排序类型返回 ORDER BY 子句
func getOrderClause(sort models.FundSortType) string {
	switch sort {
	case models.FundSortTypeWeek:
		return "performance->>'week_profit_ratio' DESC"
	case models.FundSortTypeMonth1:
		return "performance->>'month_1_profit_ratio' DESC"
	case models.FundSortTypeMonth3:
		return "performance->>'month_3_profit_ratio' DESC"
	case models.FundSortTypeMonth6:
		return "performance->>'month_6_profit_ratio' DESC"
	case models.FundSortTypeYear1:
		return "performance->>'year_1_profit_ratio' DESC"
	case models.FundSortTypeYear2:
		return "performance->>'year_2_profit_ratio' DESC"
	case models.FundSortTypeYear3:
		return "performance->>'year_3_profit_ratio' DESC"
	case models.FundSortTypeYear5:
		return "performance->>'year_5_profit_ratio' DESC"
	case models.FundSortTypeThisYear:
		return "performance->>'this_year_profit_ratio' DESC"
	case models.FundSortTypeHistorical:
		return "performance->>'historical_profit_ratio' DESC"
	default:
		return "updated_at DESC"
	}
}

// GetFundFilter 基金筛选
func (s *FundService) GetFundFilter(ctx context.Context, params FundFilterParams) (*FundIndexResponse, error) {
	// 从数据库构建查询
	var fundDBs []models.FundDB
	query := models.DB.Model(&models.FundDB{})

	// 应用筛选条件
	filter := params.ParamFundListFilter

	// 基金类型筛选
	if len(filter.Types) > 0 {
		query = query.Where("type IN ?", filter.Types)
	}

	// 基金规模筛选
	if filter.MinScale > 0 {
		query = query.Where("net_assets_scale >= ?", filter.MinScale*100000000)
	}
	if filter.MaxScale > 0 {
		query = query.Where("net_assets_scale <= ?", filter.MaxScale*100000000)
	}

	// 成立年限筛选（需要从成立日期计算）
	// 这里简化处理，后续可以优化

	// 绩效排名筛选（使用 JSONB 查询）
	if filter.Year1RankRatio > 0 {
		query = query.Where("(performance->>'year_1_rank_ratio')::float <= ?", filter.Year1RankRatio)
	}
	if filter.Month6RankRatio > 0 {
		query = query.Where("(performance->>'month_6_rank_ratio')::float <= ?", filter.Month6RankRatio)
	}
	if filter.Month3RankRatio > 0 {
		query = query.Where("(performance->>'month_3_rank_ratio')::float <= ?", filter.Month3RankRatio)
	}

	// 基金经理年限筛选（通过经理关联表）
	if filter.MinManagerYears > 0 {
		subQuery := models.DB.Model(&models.FundManagerRelationDB{}).
			Where("manage_days >= ?", filter.MinManagerYears*365).
			Select("fund_code")
		query = query.Where("code IN (?)", subQuery)
	}

	// 波动率筛选
	if filter.Max135AvgStddev > 0 {
		query = query.Where("(stddev->>'avg_135')::float <= ?", filter.Max135AvgStddev)
	}

	// 夏普比率筛选
	if filter.Min135AvgSharp > 0 {
		query = query.Where("(sharp->>'avg_135')::float >= ?", filter.Min135AvgSharp)
	}

	// 最大回撤筛选
	if filter.Max135AvgRetr > 0 {
		query = query.Where("(max_retracement->>'avg_135')::float <= ?", filter.Max135AvgRetr)
	}

	// 基金类型额外过滤
	if params.ParamFundIndex.Type != "" {
		query = query.Where("type = ?", params.ParamFundIndex.Type)
	}

	// 获取总数
	var totalCount int64
	query.Count(&totalCount)

	// 分页
	pagi := goutils.PaginateByPageNumSize(int(totalCount), params.ParamFundIndex.PageNum, params.ParamFundIndex.PageSize)

	// 排序
	orderClause := getOrderClause(models.FundSortType(params.ParamFundIndex.Sort))
	if err := query.Order(orderClause).Offset(pagi.StartIndex).Limit(pagi.PageSize).Find(&fundDBs).Error; err != nil {
		return nil, err
	}

	// 转换为基金列表
	fundList := make([]*models.Fund, len(fundDBs))
	for i, fd := range fundDBs {
		fundList[i] = fd.ToFund()
	}

	// 获取基金类型
	var fundTypes []string
	if err := models.DB.Model(&models.FundDB{}).
		Distinct("type").
		Pluck("type", &fundTypes).Error; err != nil {
		logging.Error(ctx, "GetFundFilter get fund types error:"+err.Error())
	}

	// 计算总页数，避免除零
	totalPages := 0
	if pagi.PageSize > 0 {
		totalPages = int(totalCount) / pagi.PageSize
		if int(totalCount)%pagi.PageSize > 0 {
			totalPages++
		}
	}

	return &FundIndexResponse{
		FundList: fundList,
		Pagination: PaginationResponse{
			PageNum:    pagi.PageNum,
			PageSize:   pagi.PageSize,
			Total:      int(totalCount),
			TotalPages: totalPages,
			StartIndex: pagi.StartIndex,
			EndIndex:   pagi.EndIndex,
		},
		UpdatedAt:     models.SyncFundTime.Format("2006-01-02 15:04:05"),
		AllFundCount:  int(totalCount),
		Fund4433Count: int(totalCount),
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
	managers := eastmoney.FundManagerInfoList{}.Filter(ctx, eastmoney.ParamFundManagerFilter{
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

	// 计算总页数，避免除零
	totalPages := 0
	if pagi.PageSize > 0 {
		totalPages = len(managers) / pagi.PageSize
		if len(managers)%pagi.PageSize > 0 {
			totalPages++
		}
	}

	return &FundManagerResponse{
		Managers: []FundManagerInfo{},
		Pagination: PaginationResponse{
			PageNum:    pagi.PageNum,
			PageSize:   pagi.PageSize,
			Total:      len(managers),
			TotalPages: totalPages,
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
