# 增量更新方案设计

## 问题分析

当前系统的数据更新问题：
1. **全量更新**：每次更新都要获取所有基金数据并重新查询
2. **性能问题**：API 调用频繁，内存占用大
3. **更新慢**：每次更新耗时很长
4. **资源浪费**：大部分基金数据没有变化

## 解决方案

### 方案一：基于数据库的增量更新（推荐）

#### 核心思路
1. 使用 PostgreSQL 存储基金数据
2. 记录每个基金的上次更新时间
3. 只更新需要更新的基金
4. 按优先级和上次更新时间智能判断

#### 数据库设计

```sql
-- 基金主表
CREATE TABLE funds (
    code VARCHAR(20) PRIMARY KEY,
    name VARCHAR(200),
    type VARCHAR(50),
    established_date DATE,
    net_assets_scale DECIMAL(20,2),
    index_code VARCHAR(20),
    index_name VARCHAR(200),
    rate VARCHAR(50),
    fixed_investment_status VARCHAR(20),
    
    -- 性能指标（JSONB存储复杂数据）
    stddev JSONB,
    max_retracement JSONB,
    sharp JSONB,
    performance JSONB,
    
    -- 更新时间追踪
    last_sync_time TIMESTAMP DEFAULT NOW(),
    last_update_time TIMESTAMP DEFAULT NOW(),
    sync_version INTEGER DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 基金持仓股票
CREATE TABLE fund_stocks (
    id SERIAL PRIMARY KEY,
    fund_code VARCHAR(20) REFERENCES funds(code) ON DELETE CASCADE,
    stock_code VARCHAR(20),
    stock_name VARCHAR(200),
    industry VARCHAR(100),
    ex_code VARCHAR(10),
    hold_ratio DECIMAL(10,4),
    adjust_ratio DECIMAL(10,4),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 基金经理关联
CREATE TABLE fund_manager_relations (
    id SERIAL PRIMARY KEY,
    fund_code VARCHAR(20) REFERENCES funds(code) ON DELETE CASCADE,
    manager_id VARCHAR(50),
    manager_name VARCHAR(100),
    working_days INTEGER,
    manage_days INTEGER,
    manage_repay DECIMAL(10,4),
    years_avg_repay DECIMAL(10,4),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(fund_code, manager_id)
);

-- 基金分红历史
CREATE TABLE fund_dividends (
    id SERIAL PRIMARY KEY,
    fund_code VARCHAR(20) REFERENCES funds(code) ON DELETE CASCADE,
    reg_date VARCHAR(20),
    value DECIMAL(10,4),
    ration_date VARCHAR(20),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 基金资产占比
CREATE TABLE fund_assets_proportion (
    id SERIAL PRIMARY KEY,
    fund_code VARCHAR(20) REFERENCES funds(code) ON DELETE CASCADE,
    pub_date VARCHAR(20),
    stock VARCHAR(50),
    bond VARCHAR(50),
    cash VARCHAR(50),
    other VARCHAR(50),
    net_assets VARCHAR(50),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(fund_code, pub_date)
);

-- 基金行业占比
CREATE TABLE fund_industry_proportions (
    id SERIAL PRIMARY KEY,
    fund_code VARCHAR(20) REFERENCES funds(code) ON DELETE CASCADE,
    pub_date VARCHAR(20),
    industry VARCHAR(100),
    prop VARCHAR(50),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引优化查询
CREATE INDEX idx_funds_type ON funds(type);
CREATE INDEX idx_funds_last_sync ON funds(last_sync_time);
CREATE INDEX idx_funds_sync_version ON funds(sync_version);
CREATE INDEX idx_fund_stocks_fund_code ON fund_stocks(fund_code);
CREATE INDEX idx_fund_dividends_fund_code ON fund_dividends(fund_code);
```

### 更新策略

#### 1. 初次全量导入
- 第一次运行时，全量导入所有基金数据
- 设置 `last_sync_time` 和 `sync_version`

#### 2. 增量更新策略

**策略A：按时间增量**
```go
// 每天更新的基金（热基金）
fundsToUpdate := db.Where("last_sync_time < ?", time.Now().Add(-24*time.Hour))
    .Or("sync_version = ?", 0)
    .Find(&funds)

// 每周更新的基金（温基金）
fundsToUpdate := db.Where("last_sync_time < ?", time.Now().Add(-7*24*time.Hour))
    .Find(&funds)
```

**策略B：按优先级更新**
```go
// 优先级策略
// 1. 新基金（sync_version = 0）
// 2. 久未更新的基金（超过7天）
// 3. 热门基金（访问次数多，需要定期更新）

// 查询需要更新的基金
fundsToUpdate := []Fund{}
db.Where("sync_version = 0 OR last_sync_time < ?", 
    time.Now().Add(-7*24*time.Hour)).Find(&fundsToUpdate)
```

**策略C：分批更新**
```go
// 每天更新一部分
// 例如：1000个基金，每天更新200个
limit := 200
offset := 0
db.Where("sync_version = 0 OR last_sync_time < ?", 
    time.Now().Add(-7*24*time.Hour))
    .Limit(limit)
    .Offset(offset)
    .Find(&fundsToUpdate)
```

#### 3. 实现逻辑

```go
// SyncFundIncremental 增量更新基金
func SyncFundIncremental() {
    ctx := context.Background()
    
    // 1. 获取需要更新的基金
    fundsToUpdate := getFundsToUpdate(ctx)
    
    // 2. 并发更新（控制并发数）
    for _, fundCode := range fundsToUpdate {
        updateSingleFund(ctx, fundCode)
    }
    
    // 3. 更新同步时间
    updateLastSyncTime(ctx)
}

// getFundsToUpdate 获取需要更新的基金列表
func getFundsToUpdate(ctx context.Context) []string {
    // 策略1：按时间判断
    oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour)
    
    var funds []Fund
    db.Where("last_sync_time < ? OR sync_version = 0", oneWeekAgo)
       .Limit(100) // 每次更新100个
       .Find(&funds)
    
    codes := []string{}
    for _, f := range funds {
        codes = append(codes, f.Code)
    }
    return codes
}

// updateSingleFund 更新单个基金
func updateSingleFund(ctx context.Context, fundCode string) error {
    // 1. 从API获取最新数据
    fundresp, err := datacenter.EastMoney.QueryFundInfo(ctx, fundCode)
    if err != nil {
        return err
    }
    
    // 2. 转换数据
    fund := models.NewFund(ctx, fundresp)
    
    // 3. 更新数据库（使用事务）
    return db.Transaction(func(tx *gorm.DB) error {
        // 更新基金主表
        tx.Model(&Fund{}).Where("code = ?", fund.Code).Updates(fund)
        
        // 删除旧持仓数据
        tx.Where("fund_code = ?", fund.Code).Delete(&FundStock{})
        
        // 插入新持仓数据
        for _, stock := range fund.Stocks {
            tx.Create(&FundStock{
                FundCode: fund.Code,
                StockCode: stock.Code,
                StockName: stock.Name,
                HoldRatio: stock.HoldRatio,
                AdjustRatio: stock.AdjustRatio,
            })
        }
        
        // 更新同步时间
        tx.Model(&Fund{}).Where("code = ?", fund.Code).Updates(map[string]interface{}{
            "last_sync_time": time.Now(),
            "sync_version": gorm.Expr("sync_version + 1"),
            "updated_at": time.Now(),
        })
        
        return nil
    })
}
```

### 优势

1. **性能提升**
   - 只更新需要更新的基金
   - 减少API调用次数

2. **内存优化**
   - 不一次性加载所有数据到内存
   - 分批处理

3. **数据一致性**
   - 记录更新时间
   - 支持数据版本追踪

4. **灵活性强**
   - 可配置更新策略
   - 支持优先级

### 方案二：基于消息队列的更新

使用 RabbitMQ 或 Kafka 实现异步更新：

```go
// 1. 生产者：发现需要更新的基金
func produceUpdateMessages() {
    fundsToUpdate := getFundsToUpdate()
    for _, fundCode := range fundsToUpdate {
        queue.Publish("fund.update", fundCode)
    }
}

// 2. 消费者：更新基金
func consumeUpdateMessages() {
    queue.Consume("fund.update", func(msg string) {
        updateSingleFund(context.Background(), msg)
    })
}
```

## 推荐方案

**优先使用方案一：基于数据库的增量更新**

理由：
1. ✅ 实现简单，不增加系统复杂度
2. ✅ 性能提升明显
3. ✅ 数据一致性好
4. ✅ 便于维护和扩展
