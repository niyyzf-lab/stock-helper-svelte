package engine

import (
	"context"
	"fmt"
	"sync"

	"stock-helper-svelte/backend/api"
	"stock-helper-svelte/backend/api/types" // 确保这个导入路径正确
	"stock-helper-svelte/backend/indicators"

	lua "github.com/yuin/gopher-lua"
)

// Worker 工作单元
type Worker struct {
	id            int
	luaState      *lua.LState
	apiClient     *api.Client // 修改为具体的API客户端类型
	metrics       *ExecutionMetrics
	strategy      *Strategy
	ctx           context.Context
	statusUpdater StatusUpdater
}

// NewWorker 创建新的工作单元
func NewWorker(id int, strategy *Strategy, metrics *ExecutionMetrics, ctx context.Context, apiClient *api.Client, statusUpdater StatusUpdater) (*Worker, error) {
	L := lua.NewState()
	if L == nil {
		return nil, NewEngineError(ErrLuaStateCreation, "failed to create Lua state", nil)
	}

	worker := &Worker{
		id:            id,
		luaState:      L,
		metrics:       metrics,
		strategy:      strategy,
		ctx:           ctx,
		apiClient:     apiClient,
		statusUpdater: statusUpdater,
	}

	// 注册Lua函数
	if err := worker.registerLuaFunctions(); err != nil {
		L.Close()
		return nil, ErrLuaFuncRegFailed(err)
	}

	// 加载策略文件
	if err := L.DoFile(strategy.FilePath); err != nil {
		L.Close()
		return nil, ErrLuaScriptFailed(fmt.Errorf("failed to load strategy file: %v", err))
	}

	return worker, nil
}

// ProcessStock 处理单个股票
func (w *Worker) ProcessStock(stock types.Index) error {
	select {
	case <-w.ctx.Done():
		return NewEngineError(ErrWorkerClosed, "worker context cancelled", w.ctx.Err())
	default:
	}

	// 创建股票数据表
	stockTable := w.luaState.NewTable()
	w.luaState.SetField(stockTable, "code", lua.LString(stock.Code))
	w.luaState.SetField(stockTable, "name", lua.LString(stock.Name))
	w.luaState.SetField(stockTable, "exchange", lua.LString(stock.Exchange))

	// 调用process_stock函数
	err := w.luaState.CallByParam(lua.P{
		Fn:      w.luaState.GetGlobal("process_stock"),
		NRet:    0,
		Protect: true,
	}, stockTable)

	if err != nil {
		w.metrics.IncrementErrors()
		return ErrLuaScriptFailed(fmt.Errorf("failed to process stock %s: %v", stock.Code, err))
	}

	w.metrics.IncrementProcessed()
	return nil
}

// registerLuaFunctions 注册Lua函数
func (w *Worker) registerLuaFunctions() error {
	// 注册全局函数
	w.registerGlobalFunctions()

	// 注册API函数
	w.registerAPIFunctions()

	return nil
}

// registerGlobalFunctions 注册全局函数
func (w *Worker) registerGlobalFunctions() {
	// 注册日志函数
	w.luaState.SetGlobal("log", w.luaState.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		fmt.Printf("[Worker-%d] %s\n", w.id, msg)
		return 0
	}))
}

// registerAPIFunctions 注册API函数
func (w *Worker) registerAPIFunctions() {
	// 创建API表
	apiTable := w.luaState.NewTable()

	// 注册获取指数列表函数
	w.luaState.SetField(apiTable, "getIndexList", w.luaState.NewFunction(w.luaGetIndexList))

	// 注册获取K线数据函数
	w.luaState.SetField(apiTable, "getKLineData", w.luaState.NewFunction(w.luaGetKLineData))

	// 注册发送股票信号函数
	w.luaState.SetField(apiTable, "sendSignal", w.luaState.NewFunction(w.luaSendStockSignal))

	// 注册更新进度函数
	w.luaState.SetField(apiTable, "updateProgress", w.luaState.NewFunction(w.luaUpdateProgress))

	// 注册技术指标相关函数
	indicatorTable := w.luaState.NewTable()

	// 注册RSI函数
	w.luaState.SetField(indicatorTable, "calculateRSI", w.luaState.NewFunction(w.luaCalculateRSI))

	// 注册MACD函数
	w.luaState.SetField(indicatorTable, "calculateMACD", w.luaState.NewFunction(w.luaCalculateMACD))

	// 注册MA函数
	w.luaState.SetField(indicatorTable, "calculateMA", w.luaState.NewFunction(w.luaCalculateMA))

	// 将indicator表设置为api表的一个字段
	w.luaState.SetField(apiTable, "indicator", indicatorTable)

	// 将API表设置为全局变量
	w.luaState.SetGlobal("api", apiTable)
}

// luaGetIndexList 获取指数列表的Lua包装函数
func (w *Worker) luaGetIndexList(L *lua.LState) int {
	indices, err := w.apiClient.Market.GetIndexList(context.Background())
	if err != nil {
		luaErr := NewAPIRequestError("getIndexList", err)
		L.Push(lua.LNil)
		L.Push(lua.LString(luaErr.Error()))
		return 2
	}

	// 创建返回表
	indexTable := L.NewTable()
	for _, index := range indices {
		indexData := L.NewTable()
		L.SetField(indexData, "code", lua.LString(index.Code))
		L.SetField(indexData, "name", lua.LString(index.Name))
		L.SetField(indexData, "exchange", lua.LString(index.Exchange))
		indexTable.Append(indexData)
	}

	L.Push(indexTable)
	return 1
}

// luaGetKLineData 获取K线数据的Lua包装函数
func (w *Worker) luaGetKLineData(L *lua.LState) int {
	code := L.ToString(1)
	freq := L.ToString(2)

	if code == "" || freq == "" {
		luaErr := NewInvalidStockDataError(code, fmt.Errorf("code and freq are required"))
		L.Push(lua.LNil)
		L.Push(lua.LString(luaErr.Error()))
		return 2
	}

	data, err := w.apiClient.Market.GetKLineData(context.Background(), code, types.KLineFreq(freq))
	if err != nil {
		luaErr := NewAPIRequestError("getKLineData", err)
		L.Push(lua.LNil)
		L.Push(lua.LString(luaErr.Error()))
		return 2
	}

	// 创建返回表
	dataTable := L.NewTable()
	for _, item := range data {
		itemTable := L.NewTable()
		L.SetField(itemTable, "time", lua.LString(item.Time))
		L.SetField(itemTable, "open", lua.LNumber(item.Open))
		L.SetField(itemTable, "high", lua.LNumber(item.High))
		L.SetField(itemTable, "low", lua.LNumber(item.Low))
		L.SetField(itemTable, "close", lua.LNumber(item.Close))
		L.SetField(itemTable, "volume", lua.LNumber(item.Volume))
		L.SetField(itemTable, "amount", lua.LNumber(item.Amount))
		L.SetField(itemTable, "amplitude", lua.LNumber(item.Amplitude))
		L.SetField(itemTable, "turnover", lua.LNumber(item.Turnover))
		L.SetField(itemTable, "change", lua.LNumber(item.Change))
		L.SetField(itemTable, "changeAmount", lua.LNumber(item.ChangeAmt))
		dataTable.Append(itemTable)
	}

	L.Push(dataTable)
	return 1
}

// luaSendStockSignal 发送股票信号的Lua包装函数
func (w *Worker) luaSendStockSignal(L *lua.LState) int {
	code := L.ToString(1)
	name := L.ToString(2)
	price := float64(L.ToNumber(3))
	turnover := float64(L.ToNumber(4))
	change := float64(L.ToNumber(5))
	reason := L.ToString(6)

	signal := StockSignal{
		Code:     code,
		Name:     name,
		Price:    price,
		Turnover: turnover,
		Change:   change,
		Reason:   reason,
	}

	w.statusUpdater.AddSignal(signal)
	return 0
}

// luaUpdateProgress 更新进度的Lua包装函数
func (w *Worker) luaUpdateProgress(L *lua.LState) int {
	currentStock := L.ToString(1)
	processedStocks := L.ToInt(2)

	w.statusUpdater.UpdateProgress(processedStocks, currentStock)
	L.Push(lua.LBool(true))
	return 1
}

// luaCalculateRSI RSI指标的Lua包装函数
func (w *Worker) luaCalculateRSI(L *lua.LState) int {
	// 检查参数
	prices := luaTableToFloat64Slice(L.CheckTable(1))
	period := L.CheckInt(2)

	// 计算RSI
	rsi, err := indicators.CalculateRSI(prices, period)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// 创建返回表
	rsiTable := L.NewTable()
	for i, v := range rsi {
		rsiTable.RawSetInt(i+1, lua.LNumber(v))
	}

	L.Push(rsiTable)
	L.Push(lua.LNil)
	return 2
}

// luaCalculateMACD MACD指标的Lua包装函数
func (w *Worker) luaCalculateMACD(L *lua.LState) int {
	// 检查参数
	prices := luaTableToFloat64Slice(L.CheckTable(1))
	shortPeriod := L.CheckInt(2)
	longPeriod := L.CheckInt(3)
	signalPeriod := L.CheckInt(4)

	// 计算MACD
	result, err := indicators.CalculateMACD(prices, shortPeriod, longPeriod, signalPeriod)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// 创建返回表
	macdTable := L.NewTable()

	// 设置DIF数组
	difTable := L.NewTable()
	for i, v := range result.DIF {
		difTable.RawSetInt(i+1, lua.LNumber(v))
	}
	L.SetField(macdTable, "dif", difTable)

	// 设置DEA数组
	deaTable := L.NewTable()
	for i, v := range result.DEA {
		deaTable.RawSetInt(i+1, lua.LNumber(v))
	}
	L.SetField(macdTable, "dea", deaTable)

	// 设置MACD数组
	macdArrayTable := L.NewTable()
	for i, v := range result.MACD {
		macdArrayTable.RawSetInt(i+1, lua.LNumber(v))
	}
	L.SetField(macdTable, "macd", macdArrayTable)

	L.Push(macdTable)
	L.Push(lua.LNil)
	return 2
}

// luaCalculateMA MA指标的Lua包装函数
func (w *Worker) luaCalculateMA(L *lua.LState) int {
	// 检查参数
	prices := luaTableToFloat64Slice(L.CheckTable(1))
	maType := indicators.MAType(L.CheckString(2))
	period := L.CheckInt(3)

	// 计算MA
	result, err := indicators.CalculateMA(prices, maType, period)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// 创建返回表
	maTable := L.NewTable()
	for i, v := range result {
		maTable.RawSetInt(i+1, lua.LNumber(v))
	}

	L.Push(maTable)
	L.Push(lua.LNil)
	return 2
}

// luaTableToFloat64Slice 辅助函数：将Lua表转换为float64切片
func luaTableToFloat64Slice(table *lua.LTable) []float64 {
	result := make([]float64, 0, table.Len())
	table.ForEach(func(_ lua.LValue, value lua.LValue) {
		if num, ok := value.(lua.LNumber); ok {
			result = append(result, float64(num))
		}
	})
	return result
}

// Close 关闭工作单元
func (w *Worker) Close() {
	if w.luaState != nil {
		w.luaState.Close()
	}
}

// WorkerPool 工作池
type WorkerPool struct {
	workers       []*Worker
	stockChan     chan types.Index
	errorChan     chan error
	metrics       *ExecutionMetrics
	apiClient     *api.Client
	statusUpdater StatusUpdater
	wg            sync.WaitGroup
	taskWg        sync.WaitGroup // 用于等待所有任务完成
	ctx           context.Context
	cancel        context.CancelFunc
	closed        bool       // 标记是否已关闭
	closeMutex    sync.Mutex // 保护 closed 标记
	closeOnce     sync.Once  // 确保只关闭一次
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool(size int, strategy *Strategy, metrics *ExecutionMetrics, ctx context.Context, apiClient *api.Client, statusUpdater StatusUpdater) (*WorkerPool, error) {
	if size <= 0 {
		return nil, NewInvalidConfigError("WorkerPoolSize", fmt.Errorf("worker pool size must be positive"))
	}

	if strategy == nil {
		return nil, NewInvalidConfigError("Strategy", fmt.Errorf("strategy is required"))
	}

	if metrics == nil {
		return nil, NewInvalidConfigError("Metrics", fmt.Errorf("metrics is required"))
	}

	if apiClient == nil {
		return nil, NewInvalidConfigError("APIClient", fmt.Errorf("API client is required"))
	}

	if statusUpdater == nil {
		return nil, NewInvalidConfigError("StatusUpdater", fmt.Errorf("status updater is required"))
	}

	poolCtx, cancel := context.WithCancel(ctx)

	pool := &WorkerPool{
		workers:       make([]*Worker, size),
		stockChan:     make(chan types.Index, size*2),
		errorChan:     make(chan error, size),
		metrics:       metrics,
		apiClient:     apiClient,
		statusUpdater: statusUpdater,
		ctx:           poolCtx,
		cancel:        cancel,
	}

	// 创建工作单元
	for i := 0; i < size; i++ {
		worker, err := NewWorker(i, strategy, metrics, poolCtx, apiClient, statusUpdater)
		if err != nil {
			pool.Close()
			return nil, ErrWorkerPoolFailed(fmt.Errorf("failed to create worker %d: %v", i, err))
		}
		pool.workers[i] = worker
	}

	return pool, nil
}

// Submit 提交股票到工作池
func (p *WorkerPool) Submit(stock types.Index) {
	p.closeMutex.Lock()
	if p.closed || p.stockChan == nil {
		p.closeMutex.Unlock()
		return
	}
	p.closeMutex.Unlock()

	p.taskWg.Add(1)
	select {
	case <-p.ctx.Done():
		p.taskWg.Done()
		return
	default:
		select {
		case <-p.ctx.Done():
			p.taskWg.Done()
			return
		case p.stockChan <- stock:
		}
	}
}

// Start 启动工作池
func (p *WorkerPool) Start(ctx context.Context) {
	for _, worker := range p.workers {
		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case stock, ok := <-p.stockChan:
					if !ok {
						return
					}
					if err := w.ProcessStock(stock); err != nil {
						p.errorChan <- err
					}
					p.taskWg.Done()
				}
			}
		}(worker)
	}
}

// Wait 等待所有任务完成
func (p *WorkerPool) Wait() {
	p.taskWg.Wait() // 等待所有任务完成
	p.closeOnce.Do(func() {
		if p.stockChan != nil {
			close(p.stockChan)
		}
	})
	p.wg.Wait() // 等待所有工作协程退出
}

// DrainErrors 清空错误通道
func (p *WorkerPool) DrainErrors() []error {
	var errors []error
	for {
		select {
		case err := <-p.errorChan:
			errors = append(errors, err)
		default:
			return errors
		}
	}
}

// Close 关闭工作池
func (p *WorkerPool) Close() {
	p.closeMutex.Lock()
	if p.closed {
		p.closeMutex.Unlock()
		return
	}
	p.closed = true
	p.closeMutex.Unlock()

	// 取消上下文
	if p.cancel != nil {
		p.cancel()
	}

	// 等待所有任务完成并关闭通道
	p.Wait()

	// 关闭所有工作单元
	for _, worker := range p.workers {
		if worker != nil {
			worker.Close()
		}
	}

	// 关闭错误通道
	p.closeMutex.Lock()
	if p.errorChan != nil {
		close(p.errorChan)
		p.errorChan = nil
	}
	p.closeMutex.Unlock()
}
