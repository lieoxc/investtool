// Package cron 定时任务
package cron

import (
	"context"

	"github.com/axiaoxin-com/investool/datacenter"
	"github.com/axiaoxin-com/investool/models"
	"github.com/sirupsen/logrus"
)

// SyncFundManagers 同步基金经理
func SyncFundManagers() {
	ctx := context.Background()

	// 检查数据库是否初始化
	if models.DB == nil {
		logrus.Error("SyncFundManagers: database not initialized")
		return
	}

	managers, err := datacenter.EastMoney.FundMangers(ctx, "zq", "penavgrowth", "desc")
	if err != nil {
		logrus.Errorf("SyncFundManagers error: %v", err)
		return
	}
	managers.SortByYieldse()

	// 保存数据到数据库
	for _, manager := range managers {
		// 保存基金经理基本信息
		managerDB := models.ToFundManagerDB(manager)
		if err := models.DB.Save(managerDB).Error; err != nil {
			logrus.Errorf("SyncFundManagers Save manager error: id=%s, error=%v", manager.ID, err)
			continue
		}

		// 保存基金经理管理的基金列表
		if len(manager.FundCodes) > 0 {
			// 删除旧的基金记录
			models.DB.Where("manager_id = ?", manager.ID).Delete(&models.FundManagerFundsDB{})
			// 插入新的基金记录
			funds := models.ToFundManagerFundsDB(manager.ID, manager.FundCodes, manager.FundNames)
			if len(funds) > 0 {
				if err := models.DB.CreateInBatches(funds, 100).Error; err != nil {
					logrus.Errorf("SyncFundManagers Save manager funds error: id=%s, error=%v", manager.ID, err)
				}
			}
		}
	}

	logrus.Info("SyncFundManagers saved to database successfully")
}
