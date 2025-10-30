// Package cron 定时任务
package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/investool/datacenter"
	"github.com/axiaoxin-com/investool/models"
	"github.com/sirupsen/logrus"
)

// SyncIndustryList 同步行业列表
func SyncIndustryList() {
	ctx := context.Background()
	indlist, err := datacenter.EastMoney.QueryIndustryList(ctx)
	if err != nil {
		logrus.Errorf("SyncIndustryList QueryIndustryList error: %v", err)
		promSyncError.WithLabelValues("SyncIndustryList").Inc()
		return
	}
	if len(indlist) != 0 {
		models.StockIndustryList = indlist
	}

	// 保存数据到数据库
	if models.DB != nil {
		// 清空旧的行业数据
		if err := models.DB.Exec("TRUNCATE TABLE industries").Error; err != nil {
			logrus.Errorf("SyncIndustryList truncate error: %v", err)
			promSyncError.WithLabelValues("SyncIndustryList").Inc()
		}

		// 批量插入新的行业数据
		industries := make([]models.IndustryDB, 0, len(indlist))
		for _, industry := range indlist {
			industries = append(industries, models.IndustryDB{
				Name:      industry,
				UpdatedAt: time.Now(),
			})
		}

		if len(industries) > 0 {
			if err := models.DB.CreateInBatches(industries, 100).Error; err != nil {
				logrus.Errorf("SyncIndustryList save to database error: %v", err)
				promSyncError.WithLabelValues("SyncIndustryList").Inc()
			} else {
				logrus.Info(fmt.Sprintf("SyncIndustryList saved %d industries to database successfully", len(industries)))
			}
		}
	}

}
