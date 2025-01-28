package types

type ExecutionStatus string

const (
	ExecutionStatusIdle      ExecutionStatus = "idle"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusPaused    ExecutionStatus = "paused"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusError     ExecutionStatus = "error"
)

type ExecutionState struct {
	Status       ExecutionStatus `json:"status"`
	StrategyID   int             `json:"strategyId"`
	StrategyName string          `json:"strategyName"`
	StartTime    string          `json:"startTime"`
	Progress     Progress        `json:"progress"`
	Error        string          `json:"error,omitempty"`
	Speed        float64         `json:"speed"`        // 处理速度（个/秒）
	EstimateTime float64         `json:"estimateTime"` // 预计剩余时间（秒）
}

type Progress struct {
	TotalStocks     int     `json:"totalStocks"`
	ProcessedStocks int     `json:"processedStocks"`
	CurrentStock    string  `json:"currentStock"`
	Percentage      float64 `json:"percentage"`
}
