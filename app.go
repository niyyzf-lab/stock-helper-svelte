package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"stock-helper-svelte/backend/api"
	"stock-helper-svelte/backend/data"
	"stock-helper-svelte/backend/engine"
	"stock-helper-svelte/backend/indicators"
	"stock-helper-svelte/backend/scheduler"
	"stock-helper-svelte/backend/strategy"
	"stock-helper-svelte/backend/types"

	"github.com/tidwall/buntdb"
)

// App struct
type App struct {
	ctx             context.Context
	apiClient       *api.Client
	strategyManager *strategy.Manager
	updater         *data.Updater
	dataManager     *data.Manager
	scheduler       *scheduler.Scheduler
	db              *buntdb.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// 获取用户的 AppData 目录
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("无法获取用户配置目录:", err)
	}

	// 创建应用数据目录
	appDir := filepath.Join(appDataDir, "stock-helper-svelte.exe")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		log.Fatal("无法创建应用数据目录:", err)
	}

	// 初始化数据库
	dbPath := filepath.Join(appDir, "cache.db")
	db, err := buntdb.Open(dbPath)
	if err != nil {
		log.Fatal("无法打开数据库:", err)
	}
	a.db = db

	// 设置数据库配置
	err = db.SetConfig(buntdb.Config{
		SyncPolicy: buntdb.Always, // 文件模式下使用同步写入
	})
	if err != nil {
		log.Fatal("无法设置数据库配置:", err)
	}

	// 初始化 API 客户端
	APIBase := "http://b.biyingapi.com"
	Licence := "546DE618-D8DA-40C6-9274-D28FFF9E1130"
	apiClient, err := api.NewClient(APIBase, Licence, db)
	if err != nil {
		log.Fatal("无法初始化 API 客户端:", err)
	}
	a.apiClient = apiClient

	// 初始化数据管理器
	a.dataManager = data.NewManager(a.apiClient)

	// 初始化策略管理器
	a.strategyManager = strategy.NewManager("filterLua", a.apiClient, ctx)

	// 初始化数据更新器
	a.updater = data.NewUpdater(a.apiClient, ctx)

	// 初始化调度器
	a.scheduler = scheduler.NewScheduler(ctx, a.dataManager)

	// 启动调度器
	if err := a.scheduler.Start(); err != nil {
		log.Printf("警告: 启动定时任务调度器失败: %v\n", err)
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetStrategies returns all available strategies
func (a *App) GetStrategies() []engine.Strategy {
	return a.strategyManager.GetStrategies()
}

// GetStrategyByID returns a strategy by its ID
func (a *App) GetStrategyByID(id int) (*engine.Strategy, error) {
	return a.strategyManager.GetStrategyByID(id)
}

// ExecuteStrategy 执行策略
func (a *App) ExecuteStrategy(strategyID int) error {
	// 获取策略信息
	strategy, err := a.strategyManager.GetStrategyByID(strategyID)
	if err != nil {
		return fmt.Errorf("failed to get strategy: %v", err)
	}

	// 执行策略
	if err := a.strategyManager.ExecuteStrategy(strategy); err != nil {
		return fmt.Errorf("failed to start execution: %v", err)
	}

	return nil
}

// GetDataUpdateStatus 获取数据更新状态
func (a *App) GetDataUpdateStatus() data.UpdateStatus {
	return a.updater.GetStatus()
}

// GetIndexList 获取指数列表
func (a *App) GetIndexList() ([]api.Index, error) {
	return a.apiClient.GetIndexList()
}

// GetKLineData 获取K线数据
func (a *App) GetKLineData(code string, freq api.KLineFreq) ([]api.KLineData, error) {
	return a.apiClient.GetKLineData(code, freq)
}

// UpdateStockData 更新股票数据
func (a *App) UpdateStockData() error {
	return a.updater.UpdateData()
}

// GetExecutionState 获取执行状态
func (a *App) GetExecutionState() engine.ExecutionStatus {
	return a.strategyManager.GetCurrentStatus()
}

// PauseExecution 暂停策略执行
func (a *App) PauseExecution() {
	a.strategyManager.Pause()
}

// ResumeExecution 恢复策略执行
func (a *App) ResumeExecution() {
	a.strategyManager.Resume()
}

// StopExecution 停止策略执行
func (a *App) StopExecution() {
	a.strategyManager.Stop()
}

// GetExecutionResults 获取执行结果
type ExecutionResults struct {
	Signals     []types.StockSignal    `json:"signals"`     // 选股信号
	TotalStocks int                    `json:"totalStocks"` // 总股票数
	Status      engine.ExecutionStatus `json:"status"`      // 执行状态
}

func (a *App) GetExecutionResults() ExecutionResults {
	status := a.strategyManager.GetCurrentStatus()
	signals := a.strategyManager.GetCurrentSignals()
	return ExecutionResults{
		Signals:     signals,
		TotalStocks: status.TotalStocks,
		Status:      status,
	}
}

// GetRealtimeData 获取实时交易数据
func (a *App) GetRealtimeData(code string) (api.RealtimeData, error) {
	return a.apiClient.GetRealtimeData(code)
}

// GetHistoricalTransactions 获取历史成交分布数据
func (a *App) GetHistoricalTransactions(code string) ([]api.HistoricalTransaction, error) {
	return a.apiClient.GetHistoricalTransactions(code)
}

// beforeClose is called when the app is about to quit
func (a *App) beforeClose(_ context.Context) {
	// 停止调度器
	if a.scheduler != nil {
		a.scheduler.Stop()
	}

	// 关闭数据库
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Println("关闭数据库时发生错误:", err)
		}
	}
}

// GetSchedulerStatus 获取调度器状态
func (a *App) GetSchedulerStatus() map[string]interface{} {
	if a.scheduler != nil {
		return a.scheduler.GetStatus()
	}
	return map[string]interface{}{
		"isRunning": false,
		"lastRun":   "",
	}
}

// RunSchedulerNow 立即执行一次数据更新
func (a *App) RunSchedulerNow() error {
	if a.scheduler != nil {
		return a.scheduler.RunNow()
	}
	return fmt.Errorf("scheduler not initialized")
}

// GetExecutionRecords 获取执行记录列表
func (a *App) GetExecutionRecords() ([]engine.ExecutionRecord, error) {
	return a.strategyManager.GetExecutionRecords()
}

// GetExecutionRecord 获取具体执行记录内容
func (a *App) GetExecutionRecord(fileName string) (*engine.ExecutionResult, error) {
	return a.strategyManager.GetExecutionRecord(fileName)
}

// DeleteExecutionRecord 删除执行记录
func (a *App) DeleteExecutionRecord(fileName string) error {
	return a.strategyManager.DeleteExecutionRecord(fileName)
}

// CalculateMA 计算移动平均线
func (a *App) CalculateMA(prices []float64, maType string, period int) ([]float64, error) {
	return indicators.CalculateMA(prices, indicators.MAType(maType), period)
}

// CalculateMACD 计算MACD指标
func (a *App) CalculateMACD(prices []float64) (*indicators.MACDResult, error) {
	return indicators.CalculateMACD(prices, 12, 26, 9) // 使用默认参数
}
