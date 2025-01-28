package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sort"

	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/tidwall/buntdb"
)

const (
	EndpointIndexList = "hslt/list"  // 指数列表接口
	EndpointKLine     = "hszbl/fsjy" // K线数据接口
	EndpointRealTime  = "hsrl/ssjy"  // 实时交易数据接口
	EndpointHistTrans = "hsmy/lscj"  // 历史成交分布接口

	// HTTP相关常量
	DefaultTimeout        = 30 * time.Second // 默认超时时间
	BackupTimeout         = 2 * time.Second  // 备用接口超时时间
	MaxIdleConns          = 500              // 最大空闲连接数
	MaxIdleConnsPerHost   = 100              // 每个host最大空闲连接
	IdleConnTimeout       = 90 * time.Second // 空闲连接超时时间
	TLSHandshakeTimeout   = 5 * time.Second  // TLS握手超时时间
	ExpectContinueTimeout = 1 * time.Second  // 100-continue等待时间
	ResponseHeaderTimeout = 5 * time.Second  // 响应头超时时间

	// 缓存时间常量
	CacheTime5Min  = time.Minute * 5
	CacheTime15Min = time.Minute * 15
	CacheTime30Min = time.Minute * 30
	CacheTime60Min = time.Hour
	CacheTimeDay   = time.Hour * 24

	// 工作池相关常量
	MaxWorkers = 200  // 增加工作协程数
	MaxQueue   = 5000 // 增加队列长度
	RateLimit  = 500  // 增加请求限制
	RateBurst  = 1000 // 增加突发限制

	// 重试相关常量
	MaxRetries    = 3           // 最大重试次数
	RetryInterval = time.Second // 重试间隔
)

var (
	// 对象池
	bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	gzipReaderPool = sync.Pool{
		New: func() interface{} {
			return new(gzip.Reader)
		},
	}
)

// WorkRequest 工作请求结构
type WorkRequest struct {
	ctx      context.Context
	url      string
	timeout  time.Duration
	resultCh chan<- WorkResult
}

// WorkResult 工作结果结构
type WorkResult struct {
	data []byte
	err  error
}

// WorkerPool 工作池结构
type WorkerPool struct {
	workerCount int
	queue       chan WorkRequest
	limiter     *rate.Limiter
	metrics     *PoolMetrics
}

// PoolMetrics 池指标结构
type PoolMetrics struct {
	activeWorkers  int64
	queuedRequests int64
	completedTasks int64
	failedTasks    int64
	mu             sync.RWMutex
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool(workerCount int) *WorkerPool {
	if workerCount <= 0 {
		workerCount = MaxWorkers
	}

	pool := &WorkerPool{
		workerCount: workerCount,
		queue:       make(chan WorkRequest, MaxQueue),
		limiter:     rate.NewLimiter(rate.Limit(RateLimit), RateBurst),
		metrics:     &PoolMetrics{},
	}

	// 启动工作协程
	for i := 0; i < workerCount; i++ {
		go pool.worker()
	}

	// 启动指标监控
	go pool.monitorMetrics()

	return pool
}

// monitorMetrics 监控并打印池指标
func (p *WorkerPool) monitorMetrics() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		p.metrics.mu.RLock()
		log.Printf("[指标] 活动工作协程: %d, 队列请求: %d, 完成任务: %d, 失败任务: %d",
			p.metrics.activeWorkers,
			p.metrics.queuedRequests,
			p.metrics.completedTasks,
			p.metrics.failedTasks,
		)
		p.metrics.mu.RUnlock()
	}
}

// Submit 提交工作请求
func (p *WorkerPool) Submit(ctx context.Context, url string, timeout time.Duration) ([]byte, error) {
	p.metrics.mu.Lock()
	p.metrics.queuedRequests++
	p.metrics.mu.Unlock()

	// 等待限流器许可
	if err := p.limiter.Wait(ctx); err != nil {
		p.metrics.mu.Lock()
		p.metrics.failedTasks++
		p.metrics.mu.Unlock()
		return nil, fmt.Errorf("rate limit exceeded: %v", err)
	}

	resultCh := make(chan WorkResult, 1)
	req := WorkRequest{
		ctx:      ctx,
		url:      url,
		timeout:  timeout,
		resultCh: resultCh,
	}

	// 提交请求到队列
	select {
	case p.queue <- req:
	case <-ctx.Done():
		p.metrics.mu.Lock()
		p.metrics.failedTasks++
		p.metrics.mu.Unlock()
		return nil, ctx.Err()
	}

	// 等待结果
	select {
	case result := <-resultCh:
		if result.err != nil {
			p.metrics.mu.Lock()
			p.metrics.failedTasks++
			p.metrics.mu.Unlock()
		} else {
			p.metrics.mu.Lock()
			p.metrics.completedTasks++
			p.metrics.mu.Unlock()
		}
		return result.data, result.err
	case <-ctx.Done():
		p.metrics.mu.Lock()
		p.metrics.failedTasks++
		p.metrics.mu.Unlock()
		return nil, ctx.Err()
	}
}

type Client struct {
	baseURL    string
	backupURLs []string
	licence    string
	httpClient *http.Client
	db         *buntdb.DB
	pool       *WorkerPool
	cacheMu    sync.RWMutex // 用于缓存操作的互斥锁
}

func NewClient(baseURL, licence string, db *buntdb.DB) (*Client, error) {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          MaxIdleConns,
		MaxIdleConnsPerHost:   MaxIdleConnsPerHost,
		IdleConnTimeout:       IdleConnTimeout,
		TLSHandshakeTimeout:   TLSHandshakeTimeout,
		ExpectContinueTimeout: ExpectContinueTimeout,
		ResponseHeaderTimeout: ResponseHeaderTimeout,
		// 启用HTTP/2
		ForceAttemptHTTP2: true,
		MaxConnsPerHost:   200,
		WriteBufferSize:   64 * 1024,
		ReadBufferSize:    64 * 1024,
	}

	httpClient := &http.Client{
		Timeout:   DefaultTimeout,
		Transport: transport,
	}

	client := &Client{
		baseURL: baseURL,
		backupURLs: []string{
			"http://api.biyingapi.com",
			"http://n.biyingapi.com",
			"http://b.biyingapi.com",
		},
		licence:    licence,
		httpClient: httpClient,
		db:         db,
		pool:       NewWorkerPool(MaxWorkers),
	}

	// 设置数据库配置
	if err := db.SetConfig(buntdb.Config{
		SyncPolicy:           buntdb.EverySecond,
		AutoShrinkPercentage: 100,
		AutoShrinkMinSize:    32 * 1024,
	}); err != nil {
		return nil, err
	}
	// client.ForceCleanCache()
	// 启动定期清理过期数据的goroutine
	go client.cleanupLoop()

	return client, nil
}

// cleanupLoop 定期清理过期数据
func (c *Client) cleanupLoop() {
	ticker := time.NewTicker(15 * time.Minute) // 增加清理频率到每15分钟
	defer ticker.Stop()

	for range ticker.C {
		if err := c.cleanup(); err != nil {
			log.Printf("[错误] 清理缓存失败: %v", err)
		}
	}
}

// cleanup 清理过期数据
func (c *Client) cleanup() error {
	var totalDeleted int
	const batchSize = 100 // 每批处理的key数量

	err := c.db.Update(func(tx *buntdb.Tx) error {
		var keysToDelete []string
		keysToDelete = make([]string, 0, batchSize)

		err := tx.AscendKeys("cache:*", func(key, value string) bool {
			ttl, err := tx.TTL(key)
			if err != nil {
				return true
			}

			if ttl < 0 {
				keysToDelete = append(keysToDelete, key)
				// 当达到批处理大小时，执行删除
				if len(keysToDelete) >= batchSize {
					for _, k := range keysToDelete {
						if _, err := tx.Delete(k); err != nil {
							log.Printf("[警告] 删除缓存失败: %s", k)
						} else {
							totalDeleted++
						}
					}
					keysToDelete = keysToDelete[:0] // 清空切片但保留容量
				}
			}
			return true
		})

		if err != nil {
			return fmt.Errorf("遍历缓存key失败: %v", err)
		}

		// 处理剩余的keys
		for _, key := range keysToDelete {
			if _, err := tx.Delete(key); err != nil {
				log.Printf("[警告] 删除缓存失败: %s", key)
			} else {
				totalDeleted++
			}
		}

		return nil
	})

	if err == nil && totalDeleted > 0 {
		log.Printf("[清理] 删除 %d 个过期缓存", totalDeleted)
	}

	// 只在删除了大量缓存时才进行收缩
	if totalDeleted > batchSize {
		if err := c.db.Shrink(); err != nil {
			log.Printf("[警告] 收缩数据库失败: %v", err)
		}
	}

	return err
}

// getCacheKey 生成缓存key
func (c *Client) getCacheKey(endpoint string) string {
	// 只使用endpoint作为缓存key，不包含baseURL和license
	return fmt.Sprintf("cache:%s", endpoint)
}

// getFromCache 从缓存获取数据
func (c *Client) getFromCache(key string) ([]byte, error) {
	var data []byte
	err := c.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		data = []byte(val)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// setToCache 设置缓存数据
func (c *Client) setToCache(key string, data []byte, ttl time.Duration) error {
	// 使用对象池中的buffer进行字符串转换
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	// 写入数据到buffer
	buf.Write(data)

	return c.db.Update(func(tx *buntdb.Tx) error {
		opts := &buntdb.SetOptions{
			Expires: true,
			TTL:     ttl,
		}
		_, _, err := tx.Set(key, buf.String(), opts)
		return err
	})
}

// sortKLineByDate 按日期排序K线数据(旧->新)
func sortKLineByDate(data []KLineData) {
	sort.Slice(data, func(i, j int) bool {
		// 将字符串日期转换为time.Time
		timeI, err := time.Parse("2006-01-02", data[i].Time)
		if err != nil {
			// 移除 log.Printf("解析日期失败[%s]: %v\n", data[i].Time, err)
			return false
		}
		timeJ, err := time.Parse("2006-01-02", data[j].Time)
		if err != nil {
			// 移除 log.Printf("解析日期失败[%s]: %v\n", data[j].Time, err)
			return false
		}
		return timeI.Before(timeJ)
	})
}

// getCacheTTL 根据数据类型和频率获取缓存时间
func (c *Client) getCacheTTL(endpoint string, freq KLineFreq) time.Duration {
	now := time.Now()
	today16 := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, now.Location())

	// 指数列表和日线数据共用相同的TTL逻辑
	if strings.Contains(endpoint, EndpointIndexList) ||
		(strings.Contains(endpoint, EndpointKLine) &&
			(freq == FREQ_DAILY_HFQ || freq == FREQ_WEEKLY_HFQ || freq == FREQ_MONTHLY_HFQ || freq == FREQ_YEARLY_HFQ)) {
		if now.After(today16) {
			return today16.Add(24 * time.Hour).Sub(now)
		}
		return today16.Sub(now)
	}

	// K线数据根据频率设置缓存时间
	if strings.Contains(endpoint, EndpointKLine) {
		switch freq {
		case FREQ_5MIN:
			return CacheTime5Min
		case FREQ_15MIN:
			return CacheTime15Min
		case FREQ_30MIN:
			return CacheTime30Min
		case FREQ_60MIN:
			return CacheTime60Min
		}
	}

	return CacheTimeDay
}

// Request 发送请求并返回响应数据（使用工作池）
func (c *Client) Request(ctx context.Context, path string, freq KLineFreq) ([]byte, error) {
	cacheKey := c.getCacheKey(path)

	// 使用读写锁保护缓存操作
	c.cacheMu.RLock()
	data, err := c.getFromCache(cacheKey)
	c.cacheMu.RUnlock()
	if err == nil {
		return data, nil
	}

	// 首先尝试主接口
	url := fmt.Sprintf("%s/%s/%s", c.baseURL, path, c.licence)
	data, err = c.pool.Submit(ctx, url, DefaultTimeout)
	if err == nil {
		c.cacheMu.Lock()
		if err := c.setToCache(cacheKey, data, c.getCacheTTL(path, freq)); err != nil {
			log.Printf("[警告] 缓存数据失败: %v", err)
		}
		c.cacheMu.Unlock()
		return data, nil
	}

	// 主接口失败，尝试备用接口
	for _, baseURL := range c.backupURLs {
		url := fmt.Sprintf("%s/%s/%s", baseURL, path, c.licence)
		data, err = c.pool.Submit(ctx, url, BackupTimeout)
		if err == nil {
			c.cacheMu.Lock()
			if err := c.setToCache(cacheKey, data, c.getCacheTTL(path, freq)); err != nil {
				log.Printf("[警告] 缓存数据失败: %v", err)
			}
			c.cacheMu.Unlock()
			return data, nil
		}
	}

	return nil, fmt.Errorf("所有接口请求失败，最后错误: %v", err)
}

// GetIndexList 获取指数列表
func (c *Client) GetIndexList() ([]Index, error) {
	body, err := c.Request(context.Background(), EndpointIndexList, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get index list: %v", err)
	}

	var indices []Index
	if err := json.Unmarshal(body, &indices); err != nil {
		return nil, fmt.Errorf("failed to parse index list: %v", err)
	}

	return indices, nil
}

// GetKLineData 获取K线数据
func (c *Client) GetKLineData(code string, freq KLineFreq) ([]KLineData, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", EndpointKLine, code, freq)

	// 获取缓存key
	cacheKey := c.getCacheKey(endpoint)
	// 尝试从缓存获取，如果存在就直接返回
	if data, err := c.getFromCache(cacheKey); err == nil {
		var klineData []KLineData
		if err := json.Unmarshal(data, &klineData); err == nil {
			// 验证数据时效性
			if len(klineData) > 0 {
				lastDate := klineData[len(klineData)-1].Time
				lastTime, err := time.Parse("2006-01-02", lastDate)
				if err == nil {
					today := time.Now()
					// 如果是交易日且已过16点，但数据不是今天的，强制更新
					if today.Hour() >= 16 && lastTime.Before(today) {
						goto FetchNew
					}
				}
			}
			return klineData, nil
		}
	}

FetchNew:
	body, err := c.Request(context.Background(), endpoint, freq)
	if err != nil {
		return nil, fmt.Errorf("failed to get kline data: %v", err)
	}

	var klineData []KLineData
	if err := json.Unmarshal(body, &klineData); err != nil {
		return nil, fmt.Errorf("failed to parse kline data: %v", err)
	}

	sortKLineByDate(klineData)
	return klineData, nil
}

// GetRealtimeData 获取实时交易数据
func (c *Client) GetRealtimeData(code string) (RealtimeData, error) {
	endpoint := fmt.Sprintf("%s/%s", EndpointRealTime, code)

	body, err := c.Request(context.Background(), endpoint, "")
	if err != nil {
		return RealtimeData{}, fmt.Errorf("failed to get realtime data: %v", err)
	}

	// API 返回的是一个对象,而不是数组
	var realtimeData RealtimeData
	if err := json.Unmarshal(body, &realtimeData); err != nil {
		return RealtimeData{}, fmt.Errorf("failed to parse realtime data: %v", err)
	}

	return realtimeData, nil
}

// GetHistoricalTransactions 获取历史成交分布数据
func (c *Client) GetHistoricalTransactions(code string) ([]HistoricalTransaction, error) {
	endpoint := fmt.Sprintf("%s/%s", EndpointHistTrans, code)

	body, err := c.Request(context.Background(), endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get historical transactions: %v", err)
	}

	var transactions []HistoricalTransaction
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, fmt.Errorf("failed to parse historical transactions: %v", err)
	}

	// 按时间倒序排序（新->旧）
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Time > transactions[j].Time
	})

	return transactions, nil
}

// cleanAllCache 清理所有缓存
func (c *Client) cleanAllCache() error {
	log.Printf("[清理] 开始清理所有缓存...")
	err := c.db.Update(func(tx *buntdb.Tx) error {
		var keysToDelete []string
		err := tx.AscendKeys("cache:*", func(key, value string) bool {
			keysToDelete = append(keysToDelete, key)
			return true
		})
		if err != nil {
			return fmt.Errorf("遍历缓存key失败: %v", err)
		}

		for _, key := range keysToDelete {
			if _, err := tx.Delete(key); err != nil {
				log.Printf("[警告] 删除缓存失败: %s, 错误: %v", key, err)
			} else {
				log.Printf("[清理] 删除缓存: %s", key)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("清理所有缓存失败: %v", err)
	}

	log.Printf("[清理] 缓存清理完成")
	return nil
}

// ForceCleanCache 强制清理所有缓存的公开方法
func (c *Client) ForceCleanCache() error {
	log.Printf("[清理] 开始强制清理所有缓存...")
	return c.cleanAllCache()
}

// worker 工作协程
func (p *WorkerPool) worker() {
	p.metrics.mu.Lock()
	p.metrics.activeWorkers++
	p.metrics.mu.Unlock()

	defer func() {
		p.metrics.mu.Lock()
		p.metrics.activeWorkers--
		p.metrics.mu.Unlock()
	}()

	for req := range p.queue {
		var result WorkResult
		var retries int

		for retries < MaxRetries {
			result = p.processRequest(req)
			if result.err == nil {
				break
			}

			// 如果是上下文取消或状态码错误，不重试
			if req.ctx.Err() != nil || strings.Contains(result.err.Error(), "unexpected status code") {
				break
			}

			retries++
			if retries < MaxRetries {
				// 使用指数退避策略
				backoff := time.Duration(retries*retries) * 100 * time.Millisecond
				if backoff > RetryInterval {
					backoff = RetryInterval
				}
				time.Sleep(backoff)
			}
		}

		// 发送最终结果
		select {
		case req.resultCh <- result:
		case <-req.ctx.Done():
		}
	}
}

// processRequest 处理单个请求
func (p *WorkerPool) processRequest(req WorkRequest) WorkResult {
	result := WorkResult{}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(req.ctx, "GET", req.url, nil)
	if err != nil {
		result.err = fmt.Errorf("创建请求失败: %w", err)
		return result
	}

	// 设置请求头
	httpReq.Header.Set("Accept-Encoding", "gzip, deflate")
	httpReq.Header.Set("Connection", "keep-alive")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0")
	httpReq.Header.Set("Accept", "application/json")

	// 执行请求
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		result.err = fmt.Errorf("请求失败: %w", err)
		return result
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		result.err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return result
	}

	// 获取buffer
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	// 处理gzip压缩
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader := gzipReaderPool.Get().(*gzip.Reader)
		if err := gzReader.Reset(resp.Body); err != nil {
			gzipReaderPool.Put(gzReader)
			result.err = fmt.Errorf("gzip重置失败: %w", err)
			return result
		}
		defer func() {
			gzReader.Close()
			gzipReaderPool.Put(gzReader)
		}()
		reader = gzReader
	}

	// 读取响应体
	if _, err := io.Copy(buf, reader); err != nil {
		result.err = fmt.Errorf("读取响应失败: %w", err)
		return result
	}

	// 检查响应内容
	if buf.Len() == 0 {
		result.err = fmt.Errorf("空响应")
		return result
	}

	// 返回结果
	result.data = append([]byte(nil), buf.Bytes()...)
	return result
}
