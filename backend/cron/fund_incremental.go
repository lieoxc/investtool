// Package cron 增量更新基金数据
package cron

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/axiaoxin-com/investool/datacenter"
	"github.com/axiaoxin-com/investool/datacenter/eastmoney"
	"github.com/axiaoxin-com/investool/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SyncFundIncremental 增量更新基金数据
// maxAge: 最大更新间隔（天），超过这个时间的基金会被更新
// batchSize: 每批更新的基金数量
// maxConcurrency: 最大并发数
func SyncFundIncremental(db *gorm.DB, maxAge int, batchSize int, maxConcurrency int) error {
	ctx := context.Background()
	logrus.Info("SyncFundIncremental start...")

	// 1. 获取基金列表
	efundlist, err := getFundListFromAPI(ctx)
	if err != nil {
		return err
	}

	// 2. 获取需要更新的基金
	fundCodesToUpdate := getFundsToUpdate(ctx, db, efundlist, maxAge, batchSize)

	if len(fundCodesToUpdate) == 0 {
		logrus.Info("No funds need to update")
		return nil
	}

	logrus.Infof("Found %d funds need to update", len(fundCodesToUpdate))

	// 3. 并发更新基金
	return updateFundsBatch(ctx, db, fundCodesToUpdate, maxConcurrency)
}

// getFundListFromAPI 从API获取基金列表
func getFundListFromAPI(ctx context.Context) (eastmoney.FundList, error) {
	fundTypes := []eastmoney.FundType{
		eastmoney.FundTypeBond,
	}

	efundlist := eastmoney.FundList{}
	for _, fundType := range fundTypes {
		efunds, err := datacenter.EastMoney.QueryAllFundList(ctx, fundType)
		if err != nil {
			return nil, err
		}
		efundlist = append(efundlist, efunds...)
	}
	return efundlist, nil
}

// getFundsToUpdate 获取需要更新的基金列表
func getFundsToUpdate(ctx context.Context, db *gorm.DB, efundlist eastmoney.FundList, maxAge int, batchSize int) []string {
	// 从数据库获取现有基金
	var existingFunds []models.FundDB
	db.Where("sync_version > 0").Find(&existingFunds)

	// 构建现有基金代码映射
	existingMap := make(map[string]bool)
	for _, f := range existingFunds {
		existingMap[f.Code] = true
	}

	// 找出需要更新的基金
	fundsToUpdate := []string{}
	cutoffTime := time.Now().AddDate(0, 0, -maxAge)

	// 新基金
	for _, efund := range efundlist {
		if !existingMap[efund.Fcode] {
			fundsToUpdate = append(fundsToUpdate, efund.Fcode)
		}
	}

	// 旧基金（超过maxAge天未更新的）
	var oldFunds []models.FundDB
	db.Where("last_sync_time < ? AND sync_version > 0", cutoffTime).
		Limit(batchSize - len(fundsToUpdate)).
		Find(&oldFunds)

	for _, fund := range oldFunds {
		fundsToUpdate = append(fundsToUpdate, fund.Code)
	}

	return fundsToUpdate
}

// updateFundsBatch 批量更新基金
func updateFundsBatch(ctx context.Context, db *gorm.DB, fundCodes []string, maxConcurrency int) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrency)
	results := make(chan error, len(fundCodes))

	// 预编译正则表达式
	fundCodeRegex := regexp.MustCompile(`\d{6}`)

	// 并发更新
	for _, fundCode := range fundCodes {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if !fundCodeRegex.MatchString(code) {
				results <- nil
				return
			}

			err := updateSingleFund(ctx, db, code)
			results <- err
		}(fundCode)
	}

	wg.Wait()
	close(results)

	// 检查错误
	successCount := 0
	errorCount := 0
	for err := range results {
		if err != nil {
			errorCount++
			logrus.Error(err.Error())
		} else {
			successCount++
		}
	}

	logrus.Infof("SyncFundIncremental completed: success=%d, error=%d", successCount, errorCount)

	return nil
}

// updateSingleFund 更新单个基金
func updateSingleFund(ctx context.Context, db *gorm.DB, fundCode string) error {
	// 从API获取基金详情
	fundresp := &eastmoney.RespFundInfo{}
	err := retry.Do(
		func() error {
			var err error
			fundresp, err = datacenter.EastMoney.QueryFundInfo(ctx, fundCode)
			return err
		},
		retry.OnRetry(func(n uint, err error) {
			logrus.Debugf("retry#%d: fundCode:%v %v", n, fundCode, err)
		}),
		retry.Attempts(3),
		retry.Delay(500*time.Millisecond),
	)

	if err != nil {
		return err
	}

	// 转换为模型
	fund := models.NewFund(ctx, fundresp)

	// 更新数据库
	return db.Transaction(func(tx *gorm.DB) error {
		// 转换并保存基金数据
		fundDB := convertFundToDB(fund)

		// UPSERT操作
		if err := tx.Where("code = ?", fundDB.Code).
			Assign(fundDB).
			FirstOrCreate(&fundDB).Error; err != nil {
			return errors.New("updateSingleFund FirstOrCreate error: " + err.Error())
		}

		// 删除旧持仓
		tx.Where("fund_code = ?", fund.Code).Delete(&models.FundStockDB{})

		// 插入新持仓
		for _, stock := range fund.Stocks {
			stockDB := models.FundStockDB{
				FundCode:    fund.Code,
				StockCode:   stock.Code,
				StockName:   stock.Name,
				Industry:    stock.Industry,
				ExCode:      stock.ExCode,
				HoldRatio:   stock.HoldRatio,
				AdjustRatio: stock.AdjustRatio,
				UpdatedAt:   time.Now(),
			}
			tx.Create(&stockDB)
		}

		// 更新同步时间
		return fundDB.UpdateSyncTime(tx)
	})
}

// convertFundToDB 转换Fund为FundDB
func convertFundToDB(fund *models.Fund) models.FundDB {
	now := time.Now()

	// 序列化JSONB字段
	stddevJSON, _ := json.Marshal(fund.Stddev)
	maxRetracementJSON, _ := json.Marshal(fund.MaxRetracement)
	sharpJSON, _ := json.Marshal(fund.Sharp)
	performanceJSON, _ := json.Marshal(fund.Performance)

	fundDB := models.FundDB{
		Code:                  fund.Code,
		Name:                  fund.Name,
		Type:                  fund.Type,
		EstablishedDate:       fund.EstablishedDate,
		NetAssetsScale:        fund.NetAssetsScale,
		IndexCode:             fund.IndexCode,
		IndexName:             fund.IndexName,
		Rate:                  fund.Rate,
		FixedInvestmentStatus: fund.FixedInvestmentStatus,
		Stddev:                string(stddevJSON),
		MaxRetracement:        string(maxRetracementJSON),
		Sharp:                 string(sharpJSON),
		Performance:           string(performanceJSON),
		LastSyncTime:          now,
		LastUpdateTime:        now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	return fundDB
}
