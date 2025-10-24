// API 控制器层 - 处理 HTTP 请求和响应
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FundController 基金控制器
type FundController struct {
	service *FundService
}

// NewFundController 创建基金控制器
func NewFundController() *FundController {
	return &FundController{
		service: NewFundService(),
	}
}

// GetFundIndex 获取4433基金列表
func (c *FundController) GetFundIndex(ctx *gin.Context) {
	var params FundIndexParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, BadRequestResponse("参数绑定失败", err))
		return
	}

	// 设置默认值
	if params.PageNum == 0 {
		params.PageNum = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 20
	}
	if params.Sort == 0 {
		params.Sort = 1 // FundSortTypeWeek
	}

	result, err := c.service.GetFundIndex(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalErrorResponse("获取基金列表失败", err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result))
}

// GetFundFilter 基金筛选
func (c *FundController) GetFundFilter(ctx *gin.Context) {
	var params FundFilterParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, BadRequestResponse("参数绑定失败", err))
		return
	}

	// 设置默认值
	if params.ParamFundListFilter.MinScale == 0 {
		params.ParamFundListFilter.MinScale = 2.0
	}
	if params.ParamFundListFilter.MaxScale == 0 {
		params.ParamFundListFilter.MaxScale = 50.0
	}
	if params.ParamFundListFilter.MinEstabYears == 0 {
		params.ParamFundListFilter.MinEstabYears = 5.0
	}
	if params.ParamFundListFilter.MinManagerYears == 0 {
		params.ParamFundListFilter.MinManagerYears = 5.0
	}
	if params.ParamFundListFilter.Year1RankRatio == 0 {
		params.ParamFundListFilter.Year1RankRatio = 25.0
	}
	if params.ParamFundListFilter.ThisYear235RankRatio == 0 {
		params.ParamFundListFilter.ThisYear235RankRatio = 25.0
	}
	if params.ParamFundListFilter.Month6RankRatio == 0 {
		params.ParamFundListFilter.Month6RankRatio = 33.33
	}
	if params.ParamFundListFilter.Month3RankRatio == 0 {
		params.ParamFundListFilter.Month3RankRatio = 33.33
	}

	if params.ParamFundIndex.PageNum == 0 {
		params.ParamFundIndex.PageNum = 1
	}
	if params.ParamFundIndex.PageSize == 0 {
		params.ParamFundIndex.PageSize = 20
	}

	result, err := c.service.GetFundFilter(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalErrorResponse("基金筛选失败", err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result))
}

// CheckFund 基金检测
func (c *FundController) CheckFund(ctx *gin.Context) {
	var params FundCheckParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, BadRequestResponse("参数绑定失败", err))
		return
	}

	// 设置默认值
	if params.MinScale == 0 {
		params.MinScale = 2.0
	}
	if params.MaxScale == 0 {
		params.MaxScale = 50.0
	}
	if params.MinManagerYears == 0 {
		params.MinManagerYears = 5.0
	}
	if params.Year1RankRatio == 0 {
		params.Year1RankRatio = 25.0
	}
	if params.ThisYear235RankRatio == 0 {
		params.ThisYear235RankRatio = 25.0
	}
	if params.Month6RankRatio == 0 {
		params.Month6RankRatio = 33.33
	}
	if params.Month3RankRatio == 0 {
		params.Month3RankRatio = 33.33
	}
	if params.Max135AvgStddev == 0 {
		params.Max135AvgStddev = 25.0
	}
	if params.Min135AvgSharp == 0 {
		params.Min135AvgSharp = 1.0
	}
	if params.Max135AvgRetr == 0 {
		params.Max135AvgRetr = 25.0
	}

	result, err := c.service.CheckFund(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalErrorResponse("基金检测失败", err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result))
}

// GetFundManagers 基金经理筛选
func (c *FundController) GetFundManagers(ctx *gin.Context) {
	var params FundManagerParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, BadRequestResponse("参数绑定失败", err))
		return
	}

	// 设置默认值
	if params.MinWorkingYears == 0 {
		params.MinWorkingYears = 8
	}
	if params.MinYieldse == 0 {
		params.MinYieldse = 15.0
	}
	if params.MaxCurrentFundCount == 0 {
		params.MaxCurrentFundCount = 10
	}
	if params.MinScale == 0 {
		params.MinScale = 60.0
	}
	if params.PageNum == 0 {
		params.PageNum = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 20
	}
	if params.Sort == "" {
		params.Sort = "yieldse"
	}

	result, err := c.service.GetFundManagers(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalErrorResponse("获取基金经理列表失败", err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result))
}

// GetFundSimilarity 基金持仓相似度
func (c *FundController) GetFundSimilarity(ctx *gin.Context) {
	var params FundSimilarityParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, BadRequestResponse("参数绑定失败", err))
		return
	}

	result, err := c.service.GetFundSimilarity(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalErrorResponse("基金相似度检测失败", err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result))
}

// QueryByStock 股票选基
func (c *FundController) QueryByStock(ctx *gin.Context) {
	var params QueryByStockParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, BadRequestResponse("参数绑定失败", err))
		return
	}

	result, err := c.service.QueryByStock(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, InternalErrorResponse("股票选基失败", err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result))
}
