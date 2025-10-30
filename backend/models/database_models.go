// Package models 数据库模型
package models

import (
	"encoding/json"
	"time"

	"github.com/axiaoxin-com/investool/datacenter/eastmoney"
	"gorm.io/gorm"
)

// FundDB 基金数据库模型
type FundDB struct {
	Code                  string  `gorm:"primaryKey;column:code" json:"code"`
	Name                  string  `gorm:"column:name" json:"name"`
	Type                  string  `gorm:"column:type" json:"type"`
	EstablishedDate       string  `gorm:"column:established_date" json:"established_date"`
	NetAssetsScale        float64 `gorm:"column:net_assets_scale" json:"net_assets_scale"`
	IndexCode             string  `gorm:"column:index_code" json:"index_code"`
	IndexName             string  `gorm:"column:index_name" json:"index_name"`
	Rate                  string  `gorm:"column:rate" json:"rate"`
	FixedInvestmentStatus string  `gorm:"column:fixed_investment_status" json:"fixed_investment_status"`

	// JSONB 字段存储复杂数据
	Stddev         string `gorm:"column:stddev;type:jsonb" json:"stddev"`
	MaxRetracement string `gorm:"column:max_retracement;type:jsonb" json:"max_retracement"`
	Sharp          string `gorm:"column:sharp;type:jsonb" json:"sharp"`
	Performance    string `gorm:"column:performance;type:jsonb" json:"performance"`

	// 4433法则标记
	Is4433 bool `gorm:"column:is_4433;default:false;index" json:"is_4433"`

	// 更新时间追踪
	LastSyncTime   time.Time `gorm:"column:last_sync_time;index" json:"last_sync_time"`
	LastUpdateTime time.Time `gorm:"column:last_update_time" json:"last_update_time"`
	SyncVersion    int       `gorm:"column:sync_version;default:0;index" json:"sync_version"`

	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName 指定表名
func (FundDB) TableName() string {
	return "funds"
}

// FundStockDB 基金持仓股票数据库模型
type FundStockDB struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FundCode    string    `gorm:"column:fund_code;index" json:"fund_code"`
	StockCode   string    `gorm:"column:stock_code" json:"stock_code"`
	StockName   string    `gorm:"column:stock_name" json:"stock_name"`
	Industry    string    `gorm:"column:industry" json:"industry"`
	ExCode      string    `gorm:"column:ex_code" json:"ex_code"`
	HoldRatio   float64   `gorm:"column:hold_ratio" json:"hold_ratio"`
	AdjustRatio float64   `gorm:"column:adjust_ratio" json:"adjust_ratio"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (FundStockDB) TableName() string {
	return "fund_stocks"
}

// FundManagerRelationDB 基金经理关联数据库模型
type FundManagerRelationDB struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FundCode      string    `gorm:"column:fund_code;index;uniqueIndex:idx_fund_manager" json:"fund_code"`
	ManagerID     string    `gorm:"column:manager_id;uniqueIndex:idx_fund_manager" json:"manager_id"`
	ManagerName   string    `gorm:"column:manager_name" json:"manager_name"`
	WorkingDays   float64   `gorm:"column:working_days" json:"working_days"`
	ManageDays    float64   `gorm:"column:manage_days" json:"manage_days"`
	ManageRepay   float64   `gorm:"column:manage_repay" json:"manage_repay"`
	YearsAvgRepay float64   `gorm:"column:years_avg_repay" json:"years_avg_repay"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (FundManagerRelationDB) TableName() string {
	return "fund_manager_relations"
}

// UpdateSyncTime 更新同步时间
func (f *FundDB) UpdateSyncTime(db *gorm.DB) error {
	return db.Model(f).Updates(map[string]interface{}{
		"last_sync_time": time.Now(),
		"sync_version":   gorm.Expr("sync_version + 1"),
		"updated_at":     time.Now(),
	}).Error
}

// ToFund 将 FundDB 转换为 Fund
func (f *FundDB) ToFund() *Fund {
	fund := &Fund{
		Code:                  f.Code,
		Name:                  f.Name,
		Type:                  f.Type,
		EstablishedDate:       f.EstablishedDate,
		NetAssetsScale:        f.NetAssetsScale,
		IndexCode:             f.IndexCode,
		IndexName:             f.IndexName,
		Rate:                  f.Rate,
		FixedInvestmentStatus: f.FixedInvestmentStatus,
	}

	// 解析 JSONB 字段
	var stddev fundStddev
	json.Unmarshal([]byte(f.Stddev), &stddev)
	fund.Stddev = stddev

	var maxRetracement fundMaxRetracement
	json.Unmarshal([]byte(f.MaxRetracement), &maxRetracement)
	fund.MaxRetracement = maxRetracement

	var sharp fundSharp
	json.Unmarshal([]byte(f.Sharp), &sharp)
	fund.Sharp = sharp

	var performance fundPerformance
	json.Unmarshal([]byte(f.Performance), &performance)
	fund.Performance = performance

	// 加载持仓股票（如果数据库已初始化）
	if DB != nil {
		var stocks []FundStockDB
		DB.Where("fund_code = ?", f.Code).Find(&stocks)
		fund.Stocks = make([]fundStock, len(stocks))
		for i, s := range stocks {
			fund.Stocks[i] = fundStock{
				Code:        s.StockCode,
				Name:        s.StockName,
				Industry:    s.Industry,
				ExCode:      s.ExCode,
				HoldRatio:   s.HoldRatio,
				AdjustRatio: s.AdjustRatio,
			}
		}

		// 加载基金经理
		var manager FundManagerRelationDB
		if err := DB.Where("fund_code = ?", f.Code).First(&manager).Error; err == nil {
			fund.Manager = fundManager{
				ID:            manager.ManagerID,
				Name:          manager.ManagerName,
				WorkingDays:   manager.WorkingDays,
				ManageDays:    manager.ManageDays,
				ManageRepay:   manager.ManageRepay,
				YearsAvgRepay: manager.YearsAvgRepay,
			}
		}
	}

	return fund
}

// IsNeedUpdate 判断是否需要更新
// maxAge: 最大更新时间（天）
func (f *FundDB) IsNeedUpdate(maxAge int) bool {
	if f.SyncVersion == 0 {
		return true // 新基金，需要更新
	}

	age := time.Since(f.LastSyncTime).Hours() / 24
	return age > float64(maxAge)
}

// ToFundDB 将 Fund 转换为 FundDB
func (f *Fund) ToFundDB() *FundDB {
	stddevJSON, _ := json.Marshal(f.Stddev)
	maxRetracementJSON, _ := json.Marshal(f.MaxRetracement)
	sharpJSON, _ := json.Marshal(f.Sharp)
	performanceJSON, _ := json.Marshal(f.Performance)

	return &FundDB{
		Code:                  f.Code,
		Name:                  f.Name,
		Type:                  f.Type,
		EstablishedDate:       f.EstablishedDate,
		NetAssetsScale:        f.NetAssetsScale,
		IndexCode:             f.IndexCode,
		IndexName:             f.IndexName,
		Rate:                  f.Rate,
		FixedInvestmentStatus: f.FixedInvestmentStatus,
		Stddev:                string(stddevJSON),
		MaxRetracement:        string(maxRetracementJSON),
		Sharp:                 string(sharpJSON),
		Performance:           string(performanceJSON),
		LastSyncTime:          time.Now(),
		LastUpdateTime:        time.Now(),
		Is4433:                false,
	}
}

// ToFundStocks 将 Fund.Stocks 转换为 FundStockDB 列表
func (f *Fund) ToFundStocks() []FundStockDB {
	stocks := make([]FundStockDB, 0, len(f.Stocks))
	for _, stock := range f.Stocks {
		stocks = append(stocks, FundStockDB{
			FundCode:    f.Code,
			StockCode:   stock.Code,
			StockName:   stock.Name,
			Industry:    stock.Industry,
			ExCode:      stock.ExCode,
			HoldRatio:   stock.HoldRatio,
			AdjustRatio: stock.AdjustRatio,
			UpdatedAt:   time.Now(),
		})
	}
	return stocks
}

// ToFundManagerRelation 将 Fund.Manager 转换为 FundManagerRelationDB
func (f *Fund) ToFundManagerRelation() *FundManagerRelationDB {
	if f.Manager.ID == "" {
		return nil
	}
	return &FundManagerRelationDB{
		FundCode:      f.Code,
		ManagerID:     f.Manager.ID,
		ManagerName:   f.Manager.Name,
		WorkingDays:   f.Manager.WorkingDays,
		ManageDays:    f.Manager.ManageDays,
		ManageRepay:   f.Manager.ManageRepay,
		YearsAvgRepay: f.Manager.YearsAvgRepay,
		UpdatedAt:     time.Now(),
	}
}

// IndustryDB 行业列表数据库模型
type IndustryDB struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;uniqueIndex" json:"name"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (IndustryDB) TableName() string {
	return "industries"
}

// FundManagerDB 基金经理数据库模型
type FundManagerDB struct {
	ID                  string         `gorm:"primaryKey;column:id" json:"id"`
	Name                string         `gorm:"column:name" json:"name"`
	FundCompanyID       string         `gorm:"column:fund_company_id" json:"fund_company_id"`
	FundCompanyName     string         `gorm:"column:fund_company_name" json:"fund_company_name"`
	WorkingYears        float64        `gorm:"column:working_years" json:"working_years"`
	CurrentBestReturn   float64        `gorm:"column:current_best_return" json:"current_best_return"`
	CurrentBestFundCode string         `gorm:"column:current_best_fund_code" json:"current_best_fund_code"`
	CurrentBestFundName string         `gorm:"column:current_best_fund_name" json:"current_best_fund_name"`
	CurrentFundScale    float64        `gorm:"column:current_fund_scale" json:"current_fund_scale"`
	WorkingBestReturn   float64        `gorm:"column:working_best_return" json:"working_best_return"`
	Yieldse             float64        `gorm:"column:yieldse" json:"yieldse"`
	CurrentBestFundType string         `gorm:"column:current_best_fund_type" json:"current_best_fund_type"`
	Score               float64        `gorm:"column:score" json:"score"`
	Resume              string         `gorm:"column:resume;type:text" json:"resume"`
	AwardNum            int            `gorm:"column:award_num" json:"award_num"`
	CreatedAt           time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName 指定表名
func (FundManagerDB) TableName() string {
	return "fund_managers"
}

// FundManagerFundsDB 基金经理管理的基金关联表
type FundManagerFundsDB struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ManagerID string    `gorm:"column:manager_id;index" json:"manager_id"`
	FundCode  string    `gorm:"column:fund_code;index;uniqueIndex:idx_manager_fund" json:"fund_code"`
	FundName  string    `gorm:"column:fund_name" json:"fund_name"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (FundManagerFundsDB) TableName() string {
	return "fund_manager_funds"
}

// ToFundManagerDB 将 eastmoney.FundManagerInfo 转换为 FundManagerDB
func ToFundManagerDB(manager *eastmoney.FundManagerInfo) *FundManagerDB {
	return &FundManagerDB{
		ID:                  manager.ID,
		Name:                manager.Name,
		FundCompanyID:       manager.FundCompanyID,
		FundCompanyName:     manager.FundCompanyName,
		WorkingYears:        manager.WorkingYears,
		CurrentBestReturn:   manager.CurrentBestReturn,
		CurrentBestFundCode: manager.CurrentBestFundCode,
		CurrentBestFundName: manager.CurrentBestFundName,
		CurrentFundScale:    manager.CurrentFundScale,
		WorkingBestReturn:   manager.WorkingBestReturn,
		Yieldse:             manager.Yieldse,
		CurrentBestFundType: manager.CurrentBestFundType,
		Score:               manager.Score,
		Resume:              manager.Resume,
		AwardNum:            manager.AwardNum,
		UpdatedAt:           time.Now(),
	}
}

// ToFundManagerFundsDB 将基金经理管理的基金列表转换为 FundManagerFundsDB
func ToFundManagerFundsDB(managerID string, fundCodes, fundNames []string) []FundManagerFundsDB {
	result := make([]FundManagerFundsDB, len(fundCodes))
	for i := range fundCodes {
		if i < len(fundNames) {
			result[i] = FundManagerFundsDB{
				ManagerID: managerID,
				FundCode:  fundCodes[i],
				FundName:  fundNames[i],
				UpdatedAt: time.Now(),
			}
		}
	}
	return result
}

// FundDividendDB 基金分红数据库模型
type FundDividendDB struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FundCode   string    `gorm:"column:fund_code;index;uniqueIndex:idx_fund_div" json:"fund_code"`
	RegDate    string    `gorm:"column:reg_date;uniqueIndex:idx_fund_div" json:"reg_date"`
	Value      float64   `gorm:"column:value" json:"value"`
	RationDate string    `gorm:"column:ration_date" json:"ration_date"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (FundDividendDB) TableName() string {
	return "fund_dividends"
}

// FundAssetsProportionDB 基金资产占比数据库模型
type FundAssetsProportionDB struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FundCode  string    `gorm:"column:fund_code;index;uniqueIndex:idx_fund_assets" json:"fund_code"`
	PubDate   string    `gorm:"column:pub_date;uniqueIndex:idx_fund_assets" json:"pub_date"`
	Stock     string    `gorm:"column:stock" json:"stock"`
	Bond      string    `gorm:"column:bond" json:"bond"`
	Cash      string    `gorm:"column:cash" json:"cash"`
	Other     string    `gorm:"column:other" json:"other"`
	NetAssets string    `gorm:"column:net_assets" json:"net_assets"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (FundAssetsProportionDB) TableName() string {
	return "fund_assets_proportion"
}

// FundIndustryProportionDB 基金行业占比数据库模型
type FundIndustryProportionDB struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FundCode  string    `gorm:"column:fund_code;index" json:"fund_code"`
	PubDate   string    `gorm:"column:pub_date;index" json:"pub_date"`
	Industry  string    `gorm:"column:industry" json:"industry"`
	Prop      string    `gorm:"column:prop" json:"prop"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (FundIndustryProportionDB) TableName() string {
	return "fund_industry_proportions"
}

// ToFundDividends 将 Fund.HistoricalDividends 转换为 FundDividendDB 列表
func (f *Fund) ToFundDividends() []FundDividendDB {
	dividends := make([]FundDividendDB, 0, len(f.HistoricalDividends))
	for _, div := range f.HistoricalDividends {
		dividends = append(dividends, FundDividendDB{
			FundCode:   f.Code,
			RegDate:    div.RegDate,
			Value:      div.Value,
			RationDate: div.RationDate,
			UpdatedAt:  time.Now(),
		})
	}
	return dividends
}

// ToFundAssetsProportion 将 Fund.AssetsProportion 转换为 FundAssetsProportionDB
func (f *Fund) ToFundAssetsProportion() *FundAssetsProportionDB {
	if f.AssetsProportion.PubDate == "" {
		return nil
	}
	return &FundAssetsProportionDB{
		FundCode:  f.Code,
		PubDate:   f.AssetsProportion.PubDate,
		Stock:     f.AssetsProportion.Stock,
		Bond:      f.AssetsProportion.Bond,
		Cash:      f.AssetsProportion.Cash,
		Other:     f.AssetsProportion.Other,
		NetAssets: f.AssetsProportion.NetAssets,
		UpdatedAt: time.Now(),
	}
}

// ToFundIndustryProportions 将 Fund.IndustryProportions 转换为 FundIndustryProportionDB 列表
func (f *Fund) ToFundIndustryProportions() []FundIndustryProportionDB {
	proportions := make([]FundIndustryProportionDB, 0, len(f.IndustryProportions))
	for _, prop := range f.IndustryProportions {
		proportions = append(proportions, FundIndustryProportionDB{
			FundCode:  f.Code,
			PubDate:   prop.PubDate,
			Industry:  prop.Industry,
			Prop:      prop.Prop,
			UpdatedAt: time.Now(),
		})
	}
	return proportions
}
