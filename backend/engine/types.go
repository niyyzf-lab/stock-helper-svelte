package engine

import (
	"context"
	"time"

	"stock-helper-svelte/backend/api"
)

// Strategy 策略定义
type Strategy struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FilePath    string `json:"filePath"`
}

// StrategyMeta 策略元数据
type StrategyMeta struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ExecutionConfig 执行引擎配置
type ExecutionConfig struct {
	WorkerPoolSize   int             // 工作池大小
	BatchSize        int             // 批处理大小
	RetryAttempts    int             // 重试次数
	RetryDelay       time.Duration   // 重试延迟
	ExecutionTimeout time.Duration   // 执行超时时间
	APIClient        *api.Client     // API客户端
	Context          context.Context // 上下文
}

// ExecutionStats 执行统计信息
type ExecutionStats struct {
	StartTime       time.Time
	ProcessedStocks int32
	TotalStocks     int32
	ErrorCount      int32
	CurrentSpeed    float64 // 每秒处理数量
	MemoryUsage     uint64  // 内存使用量
	GoroutineCount  int     // goroutine数量
}

// StatusUpdater 状态更新接口
type StatusUpdater interface {
	UpdateStatus(ExecutionStatus)
	UpdateProgress(processedStocks int, currentStock string)
	AddSignal(StockSignal)
}

// ExecutionEngine 执行引擎接口
type ExecutionEngine interface {
	// 核心执行方法
	Execute(strategy *Strategy) error

	// 控制方法
	Pause()
	Resume()
	Stop()

	// 状态查询
	IsPaused() bool
	IsRunning() bool
	GetExecutionStats() ExecutionStats

	// 资源清理
	Close() error
}

// StockSignal 股票信号
type StockSignal struct {
	Code     string  `json:"code"`     // 股票代码
	Name     string  `json:"name"`     // 股票名称
	Price    float64 `json:"price"`    // 当前价格
	Turnover float64 `json:"turnover"` // 换手率
	Change   float64 `json:"change"`   // 涨跌幅
	Reason   string  `json:"reason"`   // 信号原因
}
