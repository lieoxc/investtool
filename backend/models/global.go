// Package models
// 全局变量

package models

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	SSLMode         string `yaml:"sslmode"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

var (
	// DBConfig 数据库配置
	DBConfig *DatabaseConfig
)

var (
	// StockIndustryList 东方财富股票行业列表
	StockIndustryList []string
	SyncFundTime      = time.Now()
	// AAACompanyBondSyl AAA公司债当期收益率
	AAACompanyBondSyl = -1.0 // datacenter.ChinaBond.QueryAAACompanyBondSyl(context.Background())
	// DB 全局数据库连接
	DB *gorm.DB
)

// LoadDatabaseConfig 从配置文件加载数据库配置
func LoadDatabaseConfig(configFile string) error {
	if configFile == "" {
		logrus.Warn("config file not specified, skipping database configuration")
		return nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		logrus.Warn("failed to read config file:" + err.Error())
		return nil
	}

	// 简单的 YAML 解析，只解析 database 部分
	var config struct {
		Database DatabaseConfig `yaml:"database"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		logrus.Warn("failed to parse config file:" + err.Error())
		return nil
	}

	DBConfig = &config.Database

	// 设置默认值
	if DBConfig.SSLMode == "" {
		DBConfig.SSLMode = "disable"
	}
	if DBConfig.MaxOpenConns == 0 {
		DBConfig.MaxOpenConns = 25
	}
	if DBConfig.MaxIdleConns == 0 {
		DBConfig.MaxIdleConns = 5
	}
	if DBConfig.ConnMaxLifetime == 0 {
		DBConfig.ConnMaxLifetime = 300
	}

	return nil
}

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	// 如果没有数据库配置，跳过初始化
	if DBConfig == nil || DBConfig.Host == "" || DBConfig.DBName == "" {
		logrus.Warn("database configuration not found, skipping database initialization")
		return nil
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		DBConfig.Host, DBConfig.User, DBConfig.Password, DBConfig.DBName, DBConfig.Port, DBConfig.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印 SQL 语句
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(DBConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(DBConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(DBConfig.ConnMaxLifetime) * time.Second)

	DB = db
	logrus.Info("database connected successfully")

	// 自动迁移数据库表结构
	if err := DB.AutoMigrate(&FundDB{}, &FundStockDB{}, &FundManagerRelationDB{}, &IndustryDB{}, &FundManagerDB{}, &FundManagerFundsDB{}, &FundDividendDB{}, &FundAssetsProportionDB{}, &FundIndustryProportionDB{}); err != nil {
		return fmt.Errorf("failed to auto migrate database: %w", err)
	}
	logrus.Info("database tables migrated successfully")

	return nil
}
