package engine

import (
	"time"
)

// ExecutionStatus 执行状态
type ExecutionStatus struct {
	Status         string    `json:"status"`         // 状态：running, paused, completed, error
	StartTime      time.Time `json:"startTime"`      // 开始时间
	TotalStocks    int       `json:"totalStocks"`    // 总股票数
	ProcessedCount int       `json:"processedCount"` // 已处理数量
	CurrentStock   string    `json:"currentStock"`   // 当前处理的股票
	Progress       float64   `json:"progress"`       // 进度百分比
	Speed          float64   `json:"speed"`          // 处理速度(个/秒)
	EstimateTime   int       `json:"estimateTime"`   // 预计剩余时间(秒)
	Error          string    `json:"error"`          // 错误信息
	StrategyId     int       `json:"strategyId"`     // 策略ID
}

// Status constants
const (
	StatusIdle      = "idle"      // 空闲状态
	StatusRunning   = "running"   // 运行中
	StatusPaused    = "paused"    // 已暂停
	StatusCompleted = "completed" // 已完成
	StatusError     = "error"     // 错误
	StatusStopped   = "stopped"   // 已停止
)

// engineState 引擎内部状态
type engineState struct {
	status         string    // 当前状态
	startTime      time.Time // 开始时间
	totalStocks    int32     // 总股票数
	processedCount int32     // 已处理数量
	currentStock   string    // 当前处理的股票
	speed          float64   // 处理速度
	error          string    // 错误信息
	strategyId     int       // 当前策略ID
	paused         bool      // 是否暂停
	shouldStop     bool      // 是否应该停止
}

// newEngineState 创建新的引擎状态
func newEngineState() *engineState {
	return &engineState{
		status: StatusIdle,
	}
}

// toExecutionStatus 转换为外部状态
func (s *engineState) toExecutionStatus() ExecutionStatus {
	progress := float64(0)
	if s.totalStocks > 0 {
		progress = float64(s.processedCount) / float64(s.totalStocks) * 100
	}

	estimateTime := 0
	if s.speed > 0 {
		remaining := s.totalStocks - s.processedCount
		estimateTime = int(float64(remaining) / s.speed)
	}

	return ExecutionStatus{
		Status:         s.status,
		StartTime:      s.startTime,
		TotalStocks:    int(s.totalStocks),
		ProcessedCount: int(s.processedCount),
		CurrentStock:   s.currentStock,
		Progress:       progress,
		Speed:          s.speed,
		EstimateTime:   estimateTime,
		Error:          s.error,
		StrategyId:     s.strategyId,
	}
}

// ExecutionResult 执行结果
type ExecutionResult struct {
	StrategyID      int           `json:"strategyId"`      // 策略ID
	StrategyName    string        `json:"strategyName"`    // 策略名称
	ExecutionTime   time.Time     `json:"executionTime"`   // 执行时间
	CompletionTime  time.Time     `json:"completionTime"`  // 完成时间
	TotalStocks     int           `json:"totalStocks"`     // 总股票数
	ProcessedStocks int           `json:"processedStocks"` // 已处理股票数
	Signals         []StockSignal `json:"signals"`         // 信号列表
}

// ExecutionRecord 执行记录
type ExecutionRecord struct {
	FileName       string    `json:"fileName"`       // 文件名
	StrategyID     int       `json:"strategyId"`     // 策略ID
	StrategyName   string    `json:"strategyName"`   // 策略名称
	ExecutionTime  time.Time `json:"executionTime"`  // 执行时间
	SignalCount    int       `json:"signalCount"`    // 信号数量
	ProcessedCount int       `json:"processedCount"` // 处理数量
	TotalStocks    int       `json:"totalStocks"`    // 总股票数
}
