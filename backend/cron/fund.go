// Package cron 定时任务
package cron

import (
	"context"
	"time"

	"github.com/axiaoxin-com/investool/core"
	"github.com/axiaoxin-com/investool/datacenter"
	"github.com/axiaoxin-com/investool/datacenter/eastmoney"
	"github.com/axiaoxin-com/investool/models"
	"github.com/sirupsen/logrus"
)

// SyncFund 同步基金数据
func SyncFund() {
	ctx := context.Background()
	logrus.Info("SyncFund request start...")
	fundTypes := []eastmoney.FundType{
		eastmoney.FundTypeBond,
	}
	efundlist := eastmoney.FundList{}
	for _, fundType := range fundTypes {
		efunds, err := datacenter.EastMoney.QueryAllFundList(ctx, fundType)
		if err != nil {
			logrus.Errorf("SyncFund QueryAllFundList error:%v", err)
			promSyncError.WithLabelValues("SyncFund").Inc()
			return
		}
		efundlist = append(efundlist, efunds...)
	}

	fundCodes := []string{}
	for _, efund := range efundlist {
		fundCodes = append(fundCodes, efund.Fcode)
	}
	s := core.NewSearcher(ctx)
	data, err := s.SearchFunds(ctx, fundCodes)
	if err != nil {
		logrus.Errorf("SyncFund SearchFunds error:%v", err)
		promSyncError.WithLabelValues("SyncFund").Inc()
		return
	}
	fundlist := models.FundList{}
	typeMap := map[string]struct{}{}
	for _, fund := range data {
		fundlist = append(fundlist, fund)
		typeMap[fund.Type] = struct{}{}
	}
	// 更新同步时间
	models.SyncFundTime = time.Now()

	// 保存数据到数据库
	if models.DB == nil {
		logrus.Warn("SyncFund database not initialized, skipping database update")
		return
	}
	// 批量保存基金数据
	for _, fund := range fundlist {
		fundDB := fund.ToFundDB()
		// 使用 CreateOrUpdate 模式保存基金基本信息
		if err := models.DB.Save(fundDB).Error; err != nil {
			logrus.Errorf("SyncFund Save fund error: code=%s, error=%v", fund.Code, err)
			promSyncError.WithLabelValues("SyncFund").Inc()
			continue
		}

		// 保存基金持仓股票
		stocks := fund.ToFundStocks()
		if len(stocks) > 0 {
			// 删除旧的持仓记录
			models.DB.Where("fund_code = ?", fund.Code).Delete(&models.FundStockDB{})
			// 插入新的持仓记录
			if err := models.DB.CreateInBatches(stocks, 100).Error; err != nil {
				logrus.Errorf("SyncFund Save stocks error: code=%s, error=%v", fund.Code, err)
				promSyncError.WithLabelValues("SyncFund").Inc()
			}
		}

		// 保存基金经理关联
		managerRel := fund.ToFundManagerRelation()
		if managerRel != nil {
			// 使用 ON CONFLICT 处理重复
			if err := models.DB.Where("fund_code = ? AND manager_id = ?", managerRel.FundCode, managerRel.ManagerID).
				Assign(*managerRel).
				FirstOrCreate(managerRel).Error; err != nil {
				logrus.Errorf("SyncFund Save manager error: code=%s, error=%v", fund.Code, err)
				promSyncError.WithLabelValues("SyncFund").Inc()
			}
		}

		// 保存基金分红记录
		dividends := fund.ToFundDividends()
		if len(dividends) > 0 {
			// 删除旧的分红记录
			models.DB.Where("fund_code = ?", fund.Code).Delete(&models.FundDividendDB{})
			// 插入新的分红记录
			if err := models.DB.CreateInBatches(dividends, 100).Error; err != nil {
				logrus.Errorf("SyncFund Save dividends error: code=%s, error=%v", fund.Code, err)
				promSyncError.WithLabelValues("SyncFund").Inc()
			}
		}

		// 保存基金资产占比
		assetsProp := fund.ToFundAssetsProportion()
		if assetsProp != nil {
			// 使用 ON CONFLICT 处理重复
			if err := models.DB.Where("fund_code = ? AND pub_date = ?", assetsProp.FundCode, assetsProp.PubDate).
				Assign(*assetsProp).
				FirstOrCreate(assetsProp).Error; err != nil {
				logrus.Errorf("SyncFund Save assets proportion error: code=%s, error=%v", fund.Code, err)
				promSyncError.WithLabelValues("SyncFund").Inc()
			}
		}

		// 保存基金行业占比
		industryProps := fund.ToFundIndustryProportions()
		if len(industryProps) > 0 {
			// 删除旧的行业占比记录
			models.DB.Where("fund_code = ?", fund.Code).Delete(&models.FundIndustryProportionDB{})
			// 插入新的行业占比记录
			if err := models.DB.CreateInBatches(industryProps, 100).Error; err != nil {
				logrus.Errorf("SyncFund Save industry proportions error: code=%s, error=%v", fund.Code, err)
				promSyncError.WithLabelValues("SyncFund").Inc()
			}
		}
	}
	// 更新4433列表
	Update4433(fundlist)

	logrus.Info("SyncFund saved to database successfully")
}

// Update4433 更新4433检测结果
func Update4433(allFundlist models.FundList) {
	ctx := context.Background()
	fundlist := models.FundList{}
	for _, fund := range allFundlist {
		if fund.Is4433(ctx) {
			fundlist = append(fundlist, fund)
		}
	}
	// 更新 models 变量
	fundlist.Sort(models.FundSortTypeWeek)

	// 更新数据库中4433基金的is_4433字段
	if models.DB == nil {
		logrus.Warn("Update4433 database not initialized, skipping database update")
		return
	}
	// 获取所有4433基金的代码
	codes := make([]string, len(fundlist))
	for i, fund := range fundlist {
		codes[i] = fund.Code
	}

	// 先标记所有4433基金为true
	if len(codes) > 0 {
		if err := models.DB.Model(&models.FundDB{}).
			Where("code IN ?", codes).
			Update("is_4433", true).Error; err != nil {
			logrus.Errorf("Update4433 update is_4433=true error: %v", err)
		}
	}

	// 将所有非4433基金标记为false
	if len(codes) > 0 {
		if err := models.DB.Model(&models.FundDB{}).
			Where("code NOT IN ?", codes).
			Update("is_4433", false).Error; err != nil {
			logrus.Errorf("Update4433 update is_4433=false error: %v", err)
		}
	}

	logrus.Infof("Update4433 found %d funds that meet 4433 criteria", len(fundlist))

}
