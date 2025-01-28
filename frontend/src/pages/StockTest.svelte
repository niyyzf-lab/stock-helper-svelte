<script lang="ts">
  import { fade, fly } from 'svelte/transition'
  import { quintOut } from 'svelte/easing'
  import { onMount } from 'svelte'
  import { Brain, ChartLine, Target, ArrowRight } from 'lucide-svelte'
  import { toastStore } from '../stores/toast'
  import Modal from '../components/Modal.svelte'
  import StockChart from '../components/StockChart.svelte'

  // 添加动画控制变量
  let mounted = false
  let showStockList = false
  let loading = false
  let stocks: any[] = []
  let selectedCount = 15 // 默认15只股票
  let selectedDays = 5 // 默认15天
  let showSettingsModal = false // 新增：控制设置弹窗显示
  
  // 添加进度状态
  let progress = {
    current: 0,
    total: selectedCount,
    currentStock: '',
    stage: '' // 用于显示当前阶段
  }

  let showChart = false
  let currentStock: any = null
  let currentStockIndex = 0

  let showActions = false

  let showTimeoutModal = false // 控制超时模态框

  let retryCount = 0
  const MAX_RETRIES = 3

  // 添加分时切换变量
  let timeframe = '5m'

  // 在 script 标签中添加新的状态变量
  let showResultModal = false
  let currentResult: {
    direction: 'up' | 'down' | 'shock'
    correct: boolean
    nextPrice?: number
    currentPrice?: number
    maxPrice?: number
    minPrice?: number
    priceChange?: number
  } | null = null

  onMount(() => {
    // 使用 RAF 确保在下一帧执行，避免闪烁
    requestAnimationFrame(() => {
      mounted = true
    })
  })

  // 获取随机股票列表
  async function getRandomStocks() {
    loading = true
    progress = {
      current: 0,
      total: selectedCount,
      currentStock: '',
      stage: '获取股票列表'
    }

    try {
      // 获取股票列表
      const allStocks = await (window as any).go.main.App.GetIndexList()
      
      // 过滤掉 ST 和退市的股票
      const validStocks = allStocks.filter(stock => {
        const name = stock.mc || ''
        return !name.includes('ST') && 
               !name.includes('退') &&
               !name.includes('退市') &&
               !name.includes('摘牌')
      })
      
      // 随机选择更多的股票以提高成功率
      const initialSelection = shuffleArray(validStocks).slice(0, selectedCount * 2)
      
      // 为每只股票获取一个随机的历史日期
      const validStockData = []
      for (let i = 0; i < initialSelection.length; i++) {
        const stock = initialSelection[i]
        progress.current = i + 1
        progress.currentStock = `${stock.mc} (${stock.dm})`
        progress.stage = '获取分时数据'
        
        try {
          // 获取15分钟级别数据
          const minuteData = await (window as any).go.main.App.GetKLineData(stock.dm, '5m')
          if (minuteData && minuteData.length > selectedDays * 48) {
            progress.stage = '处理历史数据'
            const randomIndex = Math.floor(Math.random() * (minuteData.length - selectedDays * 48))
            const selectedData = minuteData[randomIndex]
            
            stock.testDate = selectedData.d.split(' ')[0]
            stock.testTime = selectedData.d.split(' ')[1]
            stock.price = selectedData.c
            validStockData.push(stock)

            // 如果已经获取到足够的股票，就提前结束
            if (validStockData.length >= selectedCount) {
              break
            }
          }
        } catch (err) {
          console.error(`获取股票 ${stock.dm} 数据失败:`, err)
          continue
        }
      }
      
      // 更新股票列表，只包含有效的数据
      stocks = validStockData.slice(0, selectedCount)
      
      // 如果没有足够的有效股票，重试或提示错误
      if (stocks.length < selectedCount) {
        retryCount++
        if (retryCount < MAX_RETRIES) {
          toastStore.warning(`有效数据不足，正在进行第 ${retryCount} 次重试...`)
          return getRandomStocks()
        } else {
          toastStore.error('无法获取足够的有效数据，请尝试减少测试数量或缩短时间范围')
          return
        }
      }
      
      showStockList = true
      retryCount = 0 // 重置重试计数
    } catch (err) {
      console.error('获取股票数据失败:', err)
      toastStore.error('获取股票数据失败')
    } finally {
      loading = false
      // 重置进度
      progress = {
        current: 0,
        total: selectedCount,
        currentStock: '',
        stage: ''
      }
    }
  }

  // 数组随机排序
  function shuffleArray(array: any[]) {
    for (let i = array.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [array[i], array[j]] = [array[j], array[i]]
    }
    return array
  }

  // 格式化日期 - 增加更详细的格式
  function formatDate(dateStr: string, timeStr?: string) {
    if (!dateStr) return '未知时间'
    try {
      const date = new Date(dateStr)
      let result = date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
      if (timeStr) {
        result += ' ' + timeStr
      }
      return result
    } catch {
      return '无效日期'
    }
  }

  // 开始测试
  function startTest() {
    if (loading) return
    showSettingsModal = true // 显示设置弹窗而不是直接开始测试
  }

  // 新增：实际开始测试的函数
  function confirmAndStartTest() {
    showSettingsModal = false
    getRandomStocks()
  }

  // 开始实际测试
  function startStockTest() {
    showStockList = false
    currentStockIndex = 0
    currentStock = stocks[0]
    showChart = true
    
    // 5秒后显示操作按钮
    setTimeout(() => {
      showActions = true
    }, 5000)

    // 30秒后显示提示模态框
    setTimeout(() => {
      if (showActions) { // 如果用户还没有做出选择
        showTimeoutModal = true
      }
    }, 30000)
  }

  // 处理用户答案
  async function handleAnswer(direction: 'up' | 'down' | 'shock') {
    showTimeoutModal = false
    showActions = false  // 立即隐藏操作按钮
    
    if (direction !== 'skip') {
      // 获取后续数据来判断正确性
      try {
        // 获取日线数据来判断后续走势
        const nextData = await (window as any).go.main.App.GetKLineData(
          currentStock.dm,
          'dh',  // 使用日线数据
          currentStock.testDate,
          currentStock.testTime,
          5  // 获取后续5天数据
        )
        
        if (nextData && nextData.length > 0) {
          const startPrice = nextData[0].c
          const endPrice = nextData[nextData.length - 1].c
          
          // 计算区间最高最低价
          const maxPrice = Math.max(...nextData.map(d => d.h))
          const minPrice = Math.min(...nextData.map(d => d.l))
          
          // 计算涨跌幅
          const priceChange = (endPrice - startPrice) / startPrice
          
          // 判断实际走势
          // 震荡：最大涨跌幅在3%以内
          const isShock = Math.abs(priceChange) <= 0.03
          // 上涨：涨幅超过3%
          const isUp = priceChange > 0.03
          // 下跌：跌幅超过3%
          const isDown = priceChange < -0.03
          
          let actualDirection: 'up' | 'down' | 'shock'
          if (isShock) actualDirection = 'shock'
          else if (isUp) actualDirection = 'up'
          else actualDirection = 'down'
          
          // 记录结果
          currentResult = {
            direction,
            correct: direction === actualDirection,
            nextPrice: endPrice,
            currentPrice: startPrice,
            maxPrice,
            minPrice,
            priceChange
          }
          
          // 显示结果模态框
          showResultModal = true
        }
      } catch (err) {
        console.error('获取后续数据失败:', err)
        toastStore.error('获取结果数据失败')
        moveToNextStock()
      }
    } else {
      moveToNextStock()
    }
  }

  // 抽取移动到下一只股票的逻辑
  function moveToNextStock() {
    currentStockIndex++
    if (currentStockIndex < stocks.length) {
      currentStock = stocks[currentStockIndex]
      // 重置状态
      currentResult = null
      showActions = false
      // 5秒后显示操作按钮
      setTimeout(() => {
        showActions = true
      }, 5000)
    } else {
      // 测试结束
      showChart = false
      // TODO: 显示总体测试结果
    }
  }

  // 关闭结果模态框并继续
  function continueTest() {
    showResultModal = false
    moveToNextStock()
  }

  // 监听 timeframe 变化
  $: if (timeframe && currentStock) {
    // 重新加载图表数据
    // TODO: 更新 StockChart 组件，添加 timeframe 属性
  }
</script>

{#if mounted}
<div class="page-container" in:fade={{duration: 300}}>
  <div class="main-container">
    <main class="main" in:fly={{y: 20, duration: 400, delay: 300}}>
      {#if !showChart}
        <div class="content-wrapper">
          <section class="empty-state">
            <div class="empty-content">
              <div class="empty-icon">
                <Brain size={24} strokeWidth={1.5} />
              </div>
              <div class="empty-text">
                <h3>开始交易测试</h3>
                <p>通过历史数据测试您的交易判断能力</p>
              </div>
              <div class="feature-list">
                <div class="feature-item">
                  <div class="feature-icon-wrapper">
                    <Target size={20} strokeWidth={1.5} />
                  </div>
                  <div class="feature-info">
                    <h4>随机选股</h4>
                    <p>系统随机选择股票进行测试</p>
                  </div>
                </div>
                <div class="feature-item">
                  <div class="feature-icon-wrapper">
                    <ChartLine size={20} strokeWidth={1.5} />
                  </div>
                  <div class="feature-info">
                    <h4>历史数据</h4>
                    <p>基于真实历史数据进行模拟</p>
                  </div>
                </div>
              </div>
              <button class="btn solid large start-btn" on:click={startTest} disabled={loading}>
                {#if loading}
                  <div class="loading-spinner"></div>
                  <div class="loading-text">
                    <span class="stage">{progress.stage}</span>
                    <span class="progress-text">
                      {progress.currentStock ? `${progress.current}/${progress.total} - ${progress.currentStock}` : '准备中...'}
                    </span>
                  </div>
                {:else}
                  <span>开始测试</span>
                {/if}
              </button>
            </div>
          </section>
        </div>
      {:else}
        <div class="test-container">
          <div class="chart-wrapper">
            <StockChart 
              code={currentStock.dm}
              endDate={currentStock.testDate}
              defaultVisibleRange={120}
              freq={timeframe}
            />
          </div>
          
          <div class="test-footer">
            <div class="footer-content">
              <!-- 左侧股票信息 -->
              <div class="stock-info">
                <div class="stock-basic">
                  <span class="stock-code">{currentStock.dm}</span>
                  <span class="divider">·</span>
                  <span class="stock-name">{currentStock.mc}</span>
                </div>
                <span class="test-time">
                  {formatDate(currentStock.testDate, currentStock.testTime)}
                </span>
              </div>

              <!-- 中间分时切换 -->
              <div class="time-switch">
                <button 
                  class="switch-btn" 
                  class:active={timeframe === '5m'} 
                  on:click={() => timeframe = '5m'}
                >5m</button>
                <button 
                  class="switch-btn" 
                  class:active={timeframe === '15m'} 
                  on:click={() => timeframe = '15m'}
                >15m</button>
                <button 
                  class="switch-btn" 
                  class:active={timeframe === '30m'} 
                  on:click={() => timeframe = '30m'}
                >30m</button>
                <button 
                  class="switch-btn" 
                  class:active={timeframe === 'dh'} 
                  on:click={() => timeframe = 'dh'}
                >日线</button>
              </div>

              <!-- 右侧内容 -->
              <div class="footer-right">
                <!-- 进度信息 -->
                <div class="progress-info">
                  <Target size={14} strokeWidth={1.5} />
                  <span class="progress-text">
                    <span class="current">{currentStockIndex + 1}</span>
                    <span class="divider">/</span>
                    <span class="total">{stocks.length}</span>
                  </span>
                </div>

                <!-- 操作按钮 -->
                {#if showActions}
                  <div class="action-buttons" in:fly={{x: 20, duration: 400, easing: quintOut}}>
                    <button class="action-btn up" on:click={() => handleAnswer('up')}>
                      <span class="btn-icon">↑</span>
                      <span class="btn-label">看涨</span>
                    </button>
                    <button class="action-btn shock" on:click={() => handleAnswer('shock')}>
                      <span class="btn-icon">↔</span>
                      <span class="btn-label">震荡</span>
                    </button>
                    <button class="action-btn down" on:click={() => handleAnswer('down')}>
                      <span class="btn-icon">↓</span>
                      <span class="btn-label">看跌</span>
                    </button>
                  </div>
                {/if}
              </div>
            </div>
          </div>
        </div>
      {/if}
    </main>
  </div>
</div>
{/if}

<!-- 修改模态框部分 -->
{#if showStockList}
  <Modal 
    show={showStockList}
    title="测试股票列表"
    on:close={() => showStockList = false}
  >
    <div class="stock-list-container">
      <div class="info-text">
        <p>系统已为您随机选择以下股票进行测试：</p>
      </div>
      <div class="stock-list">
        {#each stocks as stock}
          <div class="stock-item">
            <div class="stock-info">
              <div class="stock-name">
                <span class="code">{stock.dm}</span>
                <span class="name">{stock.mc}</span>
              </div>
              <div class="stock-date">
                测试时间: {formatDate(stock.testDate, stock.testTime)}
              </div>
            </div>
          </div>
        {/each}
      </div>
      <div class="modal-footer">
        <p class="hint-text">点击确认后将开始测试，系统会依次展示每只股票的历史走势</p>
        <button class="btn solid large" on:click={() => {
          showStockList = false
          startStockTest() // 新增：开始实际测试的函数
        }}>
          开始测试
        </button>
      </div>
    </div>
  </Modal>
{/if}

<!-- 修改设置弹窗的布局 -->
{#if showSettingsModal}
  <Modal 
    show={showSettingsModal}
    title="测试设置"
    on:close={() => showSettingsModal = false}
  >
    <div class="settings-container">
      <div class="settings-group">
        <label class="setting-label">
          <div class="label-content">
            <div class="label-main">
              <span class="label-text">测试股票数量</span>
              <div class="input-wrapper">
                <button 
                  class="stepper-btn" 
                  on:click={() => selectedCount = Math.max(1, selectedCount - 1)}
                >-</button>
                <input 
                  type="number" 
                  bind:value={selectedCount}
                  min="1"
                  max="30"
                  class="number-input"
                />
                <button 
                  class="stepper-btn" 
                  on:click={() => selectedCount = Math.min(30, selectedCount + 1)}
                >+</button>
                <span class="unit">只</span>
              </div>
            </div>
            <span class="label-desc">建议选择 5-20 只股票</span>
          </div>
        </label>
      </div>

      <div class="settings-group">
        <label class="setting-label">
          <div class="label-content">
            <div class="label-main">
              <span class="label-text">历史数据范围</span>
              <div class="input-wrapper">
                <button 
                  class="stepper-btn" 
                  on:click={() => selectedDays = Math.max(5, selectedDays - 1)}
                >-</button>
                <input 
                  type="number" 
                  bind:value={selectedDays}
                  min="5"
                  max="60"
                  class="number-input"
                />
                <button 
                  class="stepper-btn" 
                  on:click={() => selectedDays = Math.min(60, selectedDays + 1)}
                >+</button>
                <span class="unit">天</span>
              </div>
            </div>
            <span class="label-desc warning">
              <Target size={14} strokeWidth={1.5} />
              显示超过15天可能导致部分分时数据无法完整显示
            </span>
          </div>
        </label>
      </div>

      <div class="modal-footer">
        <button 
          class="btn solid large" 
          on:click={confirmAndStartTest}
        >
          确认开始
        </button>
      </div>
    </div>
  </Modal>
{/if}

<!-- 添加超时提示模态框 -->
{#if showTimeoutModal}
  <Modal 
    show={showTimeoutModal}
    title="提示"
    on:close={() => showTimeoutModal = false}
  >
    <div class="timeout-modal">
      <p>您似乎在思考很久，需要更多时间吗？</p>
      <div class="modal-footer">
        <button 
          class="btn outline" 
          on:click={() => showTimeoutModal = false}
        >
          继续思考
        </button>
        <button 
          class="btn solid" 
          on:click={() => handleAnswer('skip')}
        >
          跳过此题
        </button>
      </div>
    </div>
  </Modal>
{/if}

<!-- 修改结果展示模态框 -->
{#if showResultModal && currentResult}
  <Modal 
    show={showResultModal}
    title="预测结果"
    on:close={continueTest}
  >
    <div class="result-modal">
      <!-- 预测结果标志 -->
      <div class="result-header">
        <div class="result-badge {currentResult.correct ? 'correct' : 'wrong'}">
          {#if currentResult.correct}
            <div class="badge-icon">✓</div>
            <div class="badge-text">
              <div class="badge-title">预测正确</div>
              <div class="badge-desc">恭喜您，判断准确！</div>
            </div>
          {:else}
            <div class="badge-icon">✕</div>
            <div class="badge-text">
              <div class="badge-title">预测错误</div>
              <div class="badge-desc">继续加油，下次会更好！</div>
            </div>
          {/if}
        </div>
      </div>

      <!-- 预测详情 -->
      <div class="result-details">
        <div class="detail-row">
          <div class="detail-item">
            <span class="label">您的预测</span>
            <span class="value prediction {currentResult.direction}">
              {#if currentResult.direction === 'up'}
                看涨
              {:else if currentResult.direction === 'down'}
                看跌
              {:else}
                震荡
              {/if}
            </span>
          </div>
          <div class="detail-item">
            <span class="label">实际涨跌</span>
            <span class="value change {currentResult.priceChange > 0.03 ? 'up' : currentResult.priceChange < -0.03 ? 'down' : 'shock'}">
              {#if currentResult.priceChange}
                {(currentResult.priceChange * 100).toFixed(2)}%
              {/if}
            </span>
          </div>
          <div class="detail-item">
            <span class="label">波动范围</span>
            <span class="value">
              {#if currentResult.maxPrice && currentResult.minPrice && currentResult.currentPrice}
                {(((currentResult.maxPrice - currentResult.minPrice) / currentResult.currentPrice) * 100).toFixed(2)}%
              {/if}
            </span>
          </div>
        </div>
      </div>

      <!-- 完整走势图 -->
      <div class="result-chart-wrapper">
        <div class="chart-title">完整走势</div>
        <div class="result-chart">
          <StockChart 
            code={currentStock.dm}
            endDate={currentStock.testDate}
            showFuture={true}
            freq={timeframe}
          />
        </div>
      </div>

      <!-- 底部按钮 -->
      <div class="modal-footer">
        <button class="btn solid large" on:click={continueTest}>
          继续测试
        </button>
      </div>
    </div>
  </Modal>
{/if}

<style>
  /* 全局容器 */
  .page-container {
    display: flex;
    height: 100%;
    color: #1f2937;
    overflow: hidden;
  }

  /* 主内容区域容器 */
  .main-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: 24px;
    overflow: hidden;
  }

  /* 主内容区 */
  .main {
    flex: 1;
    width: 100%;
    display: flex;
    flex-direction: column;
    overflow: auto;
    min-height: 0;
  }

  .content-wrapper {
    width: 100%;
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  /* 空状态 */
  .empty-state {
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: 0;
    flex: 1;
    background: linear-gradient(to bottom right, #f8fafc, #f0f7ff);
    min-height: 400px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .empty-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 32px;
    text-align: center;
    width: 100%;
    max-width: 640px;
    margin: 0 auto;
    padding: 48px 24px;
  }

  .empty-icon {
    width: 72px;
    height: 72px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #f0f7ff 0%, #e0f2fe 100%);
    color: var(--primary-500);
    border: 1px solid var(--primary-200);
    border-radius: 20px;
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.08);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .empty-icon:hover {
    transform: translateY(-2px) scale(1.02);
    box-shadow: 0 6px 16px rgba(37, 99, 235, 0.12);
  }

  .empty-text h3 {
    font-size: 18px;
    font-weight: 600;
    color: #111827;
    margin: 0 0 8px;
  }

  .empty-text p {
    font-size: 14px;
    color: #6b7280;
    margin: 0;
  }

  /* 特性列表 */
  .feature-list {
    display: flex;
    gap: 20px;
    width: 100%;
    margin: 8px 0;
  }

  .feature-item {
    flex: 1;
    display: flex;
    align-items: flex-start;
    gap: 16px;
    padding: 24px;
    background: white;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .feature-item:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.08);
    border-color: var(--primary-200);
  }

  .feature-icon-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 44px;
    height: 44px;
    background: var(--primary-50);
    color: var(--primary-500);
    border-radius: 12px;
    flex-shrink: 0;
  }

  .feature-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .feature-info h4 {
    font-size: 15px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }

  .feature-info p {
    font-size: 13px;
    color: var(--text-secondary);
    margin: 0;
    line-height: 1.5;
  }

  /* 按钮样式 */
  .btn {
    height: 34px;
    padding: 0 16px;
    font-size: 14px;
    font-weight: 500;
    border-radius: 6px;
    display: flex;
    align-items: center;
    gap: 8px;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .btn.solid {
    color: white;
    background: #2563eb;
    border: none;
  }

  .btn.solid:hover {
    background: #1d4ed8;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.2);
  }

  .btn.solid.large {
    min-height: 44px;
    height: auto;
    padding: 8px 32px;
    font-size: 15px;
    font-weight: 600;
    background: linear-gradient(135deg, #2563eb, #1d4ed8);
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(37, 99, 235, 0.25);
  }

  .btn.solid.large:hover {
    transform: translateY(-2px) scale(1.02);
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
  }

  .btn.solid.large:active {
    transform: translateY(0) scale(0.98);
  }

  /* 添加新样式 */
  .loading-spinner {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: white;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .stock-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px 0;
  }

  .stock-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .stock-item:hover {
    transform: translateX(4px);
    border-color: #93c5fd;
    background: #f8fafc;
  }

  .stock-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .stock-name {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .code {
    font-family: monospace;
    font-weight: 500;
    color: #2563eb;
  }

  .name {
    font-weight: 500;
    color: #111827;
  }

  .stock-date {
    font-size: 13px;
    color: #6b7280;
  }

  /* Modal 样式覆盖 */
  :global(.modal-content) {
    min-width: 480px !important;
  }

  /* 添加新的样式 */
  .stock-list-container {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .info-text {
    color: #4b5563;
    font-size: 14px;
    margin-bottom: 4px;
  }

  .stock-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px;
    background: #f9fafb;
    border-radius: 8px;
    max-height: 320px;
    overflow-y: auto;
  }

  .stock-item {
    display: flex;
    align-items: center;
    padding: 12px;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
  }

  .modal-footer {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
    padding-top: 16px;
    border-top: 1px solid #e5e7eb;
  }

  .hint-text {
    font-size: 13px;
    color: #6b7280;
    text-align: center;
    margin: 0;
  }

  /* 修改已有的样式 */
  .stock-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
    width: 100%;
  }

  .stock-name {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .code {
    font-family: monospace;
    font-weight: 500;
    color: #2563eb;
  }

  .name {
    font-weight: 500;
    color: #111827;
  }

  .stock-date {
    font-size: 13px;
    color: #6b7280;
  }

  /* Modal 样式覆盖 */
  :global(.modal-content) {
    min-width: 480px !important;
    max-width: 560px !important;
  }

  /* 修改按钮相关样式 */
  .start-btn {
    width: 100%;
    max-width: 320px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    background: linear-gradient(135deg, #2563eb, #1d4ed8);
    border-radius: 8px;
    font-size: 15px;
    font-weight: 600;
    color: white;
    border: none;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .start-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.25);
  }

  .start-btn:active:not(:disabled) {
    transform: translateY(0) scale(0.98);
  }

  .start-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .start-btn .loading-spinner {
    width: 20px;
    height: 20px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: white;
    animation: spin 1s linear infinite;
  }

  .start-btn .loading-text {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    min-width: 200px;
  }

  .start-btn .stage,
  .start-btn .progress-text {
    color: white;
    background: transparent;
  }

  .start-btn .progress-text {
    font-size: 12px;
    opacity: 0.8;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* 修改测试界面样式 */
  .test-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    gap: 1px;
    background: var(--neutral-100);
  }

  .chart-wrapper {
    flex: 1;
    min-height: 0;
    background: white;
    overflow: hidden;
  }

  .test-footer {
    background: white;
    padding: 12px 16px;
  }

  .footer-content {
    display: flex;
    align-items: center;
    gap: 24px;
  }

  .footer-right {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-left: auto;
  }

  .stock-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .stock-basic {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .stock-code {
    font-family: var(--font-mono);
    font-size: 15px;
    font-weight: 600;
    color: var(--primary-600);
  }

  .stock-name {
    font-size: 15px;
    font-weight: 500;
    color: var(--text-primary);
  }

  .test-time {
    font-size: 13px;
    color: var(--text-secondary);
  }

  .divider {
    color: var(--neutral-300);
  }

  .progress-info {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    background: var(--neutral-50);
    border: 1px solid var(--border-color);
    border-radius: 20px;
  }

  .progress-text {
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .progress-info .current {
    color: var(--primary-600);
    font-weight: 600;
  }

  .progress-info .total {
    color: var(--text-secondary);
  }

  .action-buttons {
    display: flex;
    gap: 8px;
    margin-left: 16px;
  }

  .action-btn {
    min-width: 90px;
    height: 32px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    border: none;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 500;
    color: white;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn .btn-icon {
    font-size: 16px;
    font-weight: 600;
  }

  .action-btn.up {
    background: linear-gradient(135deg, #ef4444, #dc2626);
  }

  .action-btn.down {
    background: linear-gradient(135deg, #10b981, #059669);
  }

  .action-btn.shock {
    background: linear-gradient(135deg, #6366f1, #4f46e5);
  }

  .action-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .action-btn:active {
    transform: translateY(0) scale(0.98);
  }

  /* 更新设置弹窗相关样式 */
  .settings-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 16px 8px;
  }

  .settings-group {
    display: flex;
    flex-direction: column;
    padding: 16px;
    background: var(--neutral-50);
    border-radius: 8px;
  }

  .setting-label {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .label-content {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .label-main {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
  }

  .label-text {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .label-desc {
    font-size: 13px;
    color: var(--text-secondary);
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .label-desc.warning {
    color: var(--warning-500);
  }

  .input-wrapper {
    display: flex;
    align-items: center;
    gap: 8px;
    background: white;
    padding: 4px;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    min-width: 160px;
  }

  .stepper-btn {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: var(--neutral-100);
    color: var(--text-primary);
    border-radius: 6px;
    font-size: 16px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .stepper-btn:hover {
    background: var(--neutral-200);
  }

  .stepper-btn:active {
    transform: scale(0.95);
  }

  .number-input {
    width: 60px;
    height: 28px;
    border: none;
    background: transparent;
    font-size: 14px;
    font-weight: 500;
    text-align: center;
    padding: 0;
  }

  .number-input:focus {
    outline: none;
  }

  .unit {
    font-size: 14px;
    color: var(--text-secondary);
    padding-right: 8px;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    padding-top: 16px;
  }

  .modal-footer .btn {
    min-width: 120px;
  }

  /* 添加超时模态框样式 */
  .timeout-modal {
    padding: 16px;
    text-align: center;
  }

  .timeout-modal p {
    margin: 0 0 24px;
    color: var(--text-primary);
  }

  .timeout-modal .modal-footer {
    display: flex;
    justify-content: center;
    gap: 16px;
  }

  .btn.outline {
    border: 1px solid var(--border-color);
    background: white;
    color: var(--text-primary);
  }

  /* 修改分时切换按钮样式 */
  .time-switch {
    display: flex;
    gap: 1px;
    padding: 2px;
    background: var(--neutral-100);
    border-radius: 6px;
  }

  .switch-btn {
    padding: 6px 10px;  /* 减小内边距使按钮更紧凑 */
    min-width: 56px;    /* 添加最小宽度保持对齐 */
    font-size: 13px;
    border: none;
    background: transparent;
    color: var(--text-secondary);
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.2s;
  }

  .switch-btn:hover {
    color: var(--text-primary);
    background: rgba(0, 0, 0, 0.04);
  }

  .switch-btn.active {
    background: white;
    color: var(--primary-600);
    font-weight: 500;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  /* 修改结果展示模态框样式 */
  .result-modal {
    display: flex;
    flex-direction: column;
    gap: 24px;
    padding: 20px;
  }

  .result-header {
    display: flex;
    justify-content: center;
  }

  .result-badge {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px 24px;
    border-radius: 12px;
    background: var(--neutral-50);
    border: 1px solid var(--border-color);
  }

  .result-badge.correct {
    background: var(--success-50);
    border-color: var(--success-200);
  }

  .result-badge.wrong {
    background: var(--error-50);
    border-color: var(--error-200);
  }

  .badge-icon {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    font-weight: bold;
    border-radius: 50%;
    background: white;
  }

  .result-badge.correct .badge-icon {
    color: var(--success-500);
  }

  .result-badge.wrong .badge-icon {
    color: var(--error-500);
  }

  .badge-text {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .badge-title {
    font-size: 16px;
    font-weight: 600;
  }

  .badge-desc {
    font-size: 13px;
    opacity: 0.8;
  }

  .result-badge.correct .badge-text {
    color: var(--success-700);
  }

  .result-badge.wrong .badge-text {
    color: var(--error-700);
  }

  .result-details {
    padding: 16px;
    background: var(--neutral-50);
    border-radius: 12px;
  }

  .detail-row {
    display: flex;
    justify-content: space-around;
    gap: 32px;
  }

  .detail-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
  }

  .detail-item .label {
    font-size: 13px;
    color: var(--text-secondary);
  }

  .detail-item .value {
    font-size: 16px;
    font-weight: 600;
  }

  .value.prediction.up,
  .value.change.up {
    color: var(--error-500);
  }

  .value.prediction.down,
  .value.change.down {
    color: var(--success-500);
  }

  .result-chart-wrapper {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .chart-title {
    font-size: 14px;
    font-weight: 500;
    color: var(--text-secondary);
  }

  .result-chart {
    height: 300px;
    background: white;
    border: 1px solid var(--border-color);
    border-radius: 12px;
    overflow: hidden;
  }

  .modal-footer {
    display: flex;
    justify-content: center;
    padding-top: 8px;
  }

  .modal-footer .btn {
    min-width: 140px;
  }

  /* 添加震荡相关样式 */
  .value.prediction.shock,
  .value.change.shock {
    color: var(--primary-500);
  }
</style> 