<script lang="ts">
  import { fade, fly } from 'svelte/transition'
  import { onMount } from 'svelte'
  import { toastStore } from '../stores/toast'
  import Modal from '../components/Modal.svelte'
  import EmptyState from '../components/stock-test/EmptyState.svelte'
  import StockListModal from '../components/stock-test/StockListModal.svelte'
  import TestView from '../components/stock-test/TestView.svelte'
  import type { StockQuiz } from '../types/stock'

  // 添加动画控制变量
  let mounted = false
  let showStockList = false
  let loading = false
  let stocks: any[] = []
  let selectedCount = 15// 固定15只股票
  let retryCount = 0
  const MAX_RETRIES = 3
  
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

  // 添加难度等级类型
  type DifficultyLevel = 1 | 2 | 3 | 4 | 5;
  
  // 添加难度选择相关状态
  let showDifficultyModal = false;
  let selectedDifficulty: DifficultyLevel | null = null;

  // 记录得分列表
  type ScoreRecord = {
    stockCode: string;
    stockName: string;
    score: number;
    difficulty: DifficultyLevel;
  };
  
  let scoreList: ScoreRecord[] = [];

  // 添加新的状态管理
  let quizCache: { [key: string]: StockQuiz } = {} // 题目缓存
  let loadingQuizzes = false // 加载状态
  let quizQueue: string[] = [] // 待加载的股票代码队列
  const PRELOAD_THRESHOLD = 5 // 预加载阈值
  const MAX_CONCURRENT_REQUESTS = 2 // 最大并发请求数
  let activeRequests = 0 // 当前活跃请求数

  let loadingFirstQuiz = false // 添加首个题目加载状态
  let hasFirstQuiz = false // 添加首个题目状态检查

  onMount(() => {
    requestAnimationFrame(() => {
      mounted = true
    })
    selectedDifficulty = getStoredDifficulty();
  })

  // 从 localStorage 获取上次选择的难度
  function getStoredDifficulty(): DifficultyLevel {
    try {
      const stored = localStorage.getItem('stockTestDifficulty');
      if (stored) {
        const level = parseInt(stored);
        if (level >= 1 && level <= 5) {
          return level as DifficultyLevel;
        }
      }
    } catch (e) {
      console.warn('Failed to read difficulty from localStorage:', e);
    }
    return 3; // 默认中等难度
  }

  // 保存难度到 localStorage
  function storeDifficulty(level: DifficultyLevel) {
    try {
      localStorage.setItem('stockTestDifficulty', level.toString());
    } catch (e) {
      console.warn('Failed to save difficulty to localStorage:', e);
    }
  }

  // 开始测试
  function startTest() {
    if (loading) return;
    showDifficultyModal = true;
  }

  // 重置测试
  function resetTest() {
    showChart = false;
    currentStock = null;
    currentStockIndex = 0;
    showActions = false;
    scoreList = [];
    stocks = [];
  }

  // 处理用户答案
  async function handleAnswer(direction: 'up' | 'down' | 'shock') {
    showActions = false;
    moveToNextStock();
  }

  // 移动到下一只股票
  function moveToNextStock() {
    currentStockIndex++
    if (currentStockIndex < stocks.length) {
      currentStock = stocks[currentStockIndex]
      showActions = false
      
      // 检查是否需要触发新的预加载
      const remainingStocks = stocks.length - currentStockIndex
      if (remainingStocks <= PRELOAD_THRESHOLD) {
        preloadQuizzes()
      }
      
      setTimeout(() => {
        showActions = true
      }, 5000)
    } else {
      showChart = false
    }
  }

  // 难度选择处理
  function handleDifficultySelect(level: DifficultyLevel) {
    selectedDifficulty = level;
    storeDifficulty(level);
    showDifficultyModal = false;
    getRandomStocks();
  }

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
      const allStocks = await (window as any).go.main.App.GetIndexList()
      
      // 过滤掉 ST 和退市的股票
      const validStocks = allStocks.filter((stock: { mc: string; dm: string }) => {
        const name = stock.mc || ''
        return !name.includes('ST') && 
               !name.includes('退') &&
               !name.includes('退市') &&
               !name.includes('摘牌')
      })
      
      // 随机选择股票
      const initialSelection = shuffleArray(validStocks).slice(0, selectedCount * 2)
      
      const validStockData = []
      for (let i = 0; i < initialSelection.length; i++) {
        const stock = initialSelection[i]
        progress.current = i + 1
        progress.currentStock = `${stock.mc} (${stock.dm})`
        progress.stage = '获取日线数据'
        
        try {
          // 获取日线数据
          const dayData = await (window as any).go.main.App.GetKLineData(stock.dm, 'dh')
          
          if (dayData && dayData.length >= 46) {
            progress.stage = '处理历史数据'
            
            // 获取最近一年的数据
            const oneYearData = dayData.slice(-252)
            
            // 确保有足够的数据进行测试
            if (oneYearData.length >= 46) {
              // 找出涨停日期
              const limitUpDays = []
              for (let j = 30; j < oneYearData.length - 15; j++) {
                const day = oneYearData[j]
                const prevDay = oneYearData[j - 1]
                // 使用新的涨停判断函数
                if (prevDay && isLimitUp(stock, prevDay.c, day.c)) {
                  limitUpDays.push(j)
                }
              }

              let selectedData
              let selectedIndex
              
              // 70%的概率选择涨停后的日期
              if (limitUpDays.length > 0 && Math.random() < 0.7) {
                // 随机选择一个涨停日
                const randomLimitUpIndex = limitUpDays[Math.floor(Math.random() * limitUpDays.length)]
                // 在涨停后3-15天内随机选择一天
                const daysAfterLimitUp = Math.floor(Math.random() * 13) + 3
                selectedIndex = randomLimitUpIndex + daysAfterLimitUp
                selectedData = oneYearData[selectedIndex]
              } else {
                // 30%概率随机选择
                const availableData = oneYearData.slice(-90, -15)
                const randomIndex = Math.floor(Math.random() * availableData.length)
                selectedData = availableData[randomIndex]
                selectedIndex = oneYearData.indexOf(selectedData)
              }
              
              // 获取历史数据(前30天)和未来数据(15天)
              if (selectedIndex >= 30 && selectedIndex + 15 < oneYearData.length) {
                const historyData = oneYearData.slice(selectedIndex - 30, selectedIndex)
                const futureData = oneYearData.slice(selectedIndex + 1, selectedIndex + 16)
                
                if (historyData.length === 30 && futureData.length === 15) {
                  stock.testDate = selectedData.d
                  stock.price = selectedData.c
                  stock.historyData = historyData
                  stock.futureData = futureData
                  validStockData.push(stock)
                }
              }

              // 如果已经获取到足够的股票，就提前结束
              if (validStockData.length >= selectedCount) {
                break
              }
            }
          }
        } catch (err) {
          console.error(`获取股票 ${stock.dm} 数据失败:`, err)
          continue
        }
      }
      
      // 更新股票列表
      stocks = validStockData.slice(0, selectedCount)
      
      // 如果没有足够的有效股票，重试或提示错误
      if (stocks.length < selectedCount) {
        retryCount++
        if (retryCount < MAX_RETRIES) {
          toastStore.warning(`有效数据不足，正在进行第 ${retryCount} 次重试...`)
          return getRandomStocks()
        } else {
          toastStore.error('无法获取足够的有效数据，请稍后重试')
          return
        }
      }
      
      showStockList = true
      retryCount = 0 // 重置重试计数
      
      // 加载首个题目
      loadFirstQuiz()
    } catch (err) {
      console.error('获取股票数据失败:', err)
      toastStore.error('获取股票数据失败')
    } finally {
      loading = false
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

  // 添加判断涨停的函数
  function isLimitUp(stock: { dm: string }, prevPrice: number, currentPrice: number): boolean {
    // 判断是否是科创板或创业板
    const isSTAR = stock.dm.startsWith('688') // 科创板
    const isGEM = stock.dm.startsWith('300') || stock.dm.startsWith('301') // 创业板
    
    const priceChange = (currentPrice - prevPrice) / prevPrice
    
    if (isSTAR || isGEM) {
      return priceChange >= 19.5 // 20%涨停
    } else {
      return priceChange >= 9.5 // 10%涨停
    }
  }

  // 开始实际测试
  async function startStockTest() {
    if (!stocks || !stocks.length) {
      console.error('没有可用的股票数据')
      return
    }
    
    showStockList = false
    currentStockIndex = 0
    currentStock = stocks[0]
    
    // 确保当前股票的题目已经在缓存中
    if (!quizCache[currentStock.dm]) {
      console.error('当前股票题目未加载')
      return
    }
    
    showChart = true
    setTimeout(() => {
      showActions = true
    }, 5000)
  }

  // 重置列表状态
  function resetListState() {
    showStockList = false
    stocks = []
    loading = false
    quizCache = {} // 清空题目缓存
    hasFirstQuiz = false
    quizQueue = [] // 清空队列
    activeRequests = 0 // 重置活跃请求数
    loadingQuizzes = false // 重置加载状态
    progress = {
      current: 0,
      total: selectedCount,
      currentStock: '',
      stage: ''
    }
  }

  // 请求单个股票的题目
  async function requestQuizForStock(stockCode: string, testDate: string): Promise<void> {
    // 检查股票是否在当前列表中
    const stockExists = stocks.some(s => s.dm === stockCode)
    if (!stockExists) {
      console.warn('请求的股票不在当前列表中:', stockCode)
      return
    }

    if (quizCache[stockCode] || 
        (stocks.length > 0 && stockCode === stocks[0].dm) || 
        activeRequests >= MAX_CONCURRENT_REQUESTS) return
    
    activeRequests++
    try {
      const quiz = await (window as any).go.main.App.GenerateStockQuizWithOptions(
        stockCode,
        testDate,
        selectedDifficulty || 3,
        0
      )
      quizCache[stockCode] = quiz
      console.log('已缓存题目:', stockCode)
    } catch (err) {
      console.error('获取AI题目失败:', stockCode, err)
    } finally {
      activeRequests--
      processQuizQueue() // 处理队列中的下一个请求
    }
  }

  // 处理题目请求队列
  async function processQuizQueue() {
    if (!quizQueue.length || activeRequests >= MAX_CONCURRENT_REQUESTS) return
    
    const nextStock = quizQueue[0]
    // 检查股票是否仍在当前列表中
    const stockData = stocks.find(s => s.dm === nextStock)
    if (stockData) {
      quizQueue.shift() // 移除队列中的第一个
      requestQuizForStock(stockData.dm, stockData.testDate)
    } else {
      // 如果股票不在列表中，直接移除
      console.warn('队列中的股票不在当前列表中，跳过:', nextStock)
      quizQueue.shift()
      processQuizQueue() // 继续处理下一个
    }
  }

  // 预加载题目
  async function preloadQuizzes() {
    if (loadingQuizzes) return
    loadingQuizzes = true
    
    try {
      // 检查stocks是否有效
      if (!stocks || !stocks.length) {
        console.warn('没有可用的股票列表')
        return
      }

      // 将未缓存的股票添加到队列
      quizQueue = stocks
        .slice(1)  // 跳过第一个股票
        .filter(stock => !quizCache[stock.dm])
        .map(stock => stock.dm)
      
      console.log('预加载队列:', quizQueue)
      
      // 开始处理队列
      for (let i = 0; i < MAX_CONCURRENT_REQUESTS; i++) {
        processQuizQueue()
      }
    } finally {
      loadingQuizzes = false
    }
  }

  // 修改首个题目加载函数
  async function loadFirstQuiz() {
    if (!stocks || !stocks.length) return
    
    loadingFirstQuiz = true
    hasFirstQuiz = false
    
    try {
      const firstStock = stocks[0]
      // 直接请求第一个题目，不通过缓存和队列机制
      activeRequests++
      const quiz = await (window as any).go.main.App.GenerateStockQuizWithOptions(
        firstStock.dm,
        firstStock.testDate,
        selectedDifficulty || 3,
        0
      )
      quizCache[firstStock.dm] = quiz
      console.log('已加载首个题目:', firstStock.dm)
      hasFirstQuiz = true
      
      // 开始预加载其他题目
      preloadQuizzes()
    } catch (err) {
      console.error('加载首个题目失败:', err)
      toastStore.error('加载题目失败，请重试')
    } finally {
      loadingFirstQuiz = false
      activeRequests--
    }
  }
</script>

{#if mounted}
<div class="page-container" in:fade={{duration: 300}}>
  <div class="main-container">
    <main class="main" in:fly={{y: 20, duration: 400, delay: 300}}>
      {#if !showChart && scoreList.length === 0}
        <div class="content-wrapper">
          <EmptyState
            {loading}
            {progress}
            onStart={startTest}
          />
        </div>
      {:else if !showChart && scoreList.length > 0}
        <div class="content-wrapper">
          <div class="score-panel">
            <h2>测试结果</h2>
            <div class="score-list">
              {#each scoreList as score}
                <div class="score-item">
                  <span>{score.stockName} ({score.stockCode})</span>
                  <span>得分: {score.score}</span>
                </div>
              {/each}
            </div>
            <button
              on:reset={() => {
                resetTest()
              }}
            >重新开始</button>
          </div>
        </div>
      {:else}
        <TestView
          {currentStock}
          {currentStockIndex}
          {stocks}
          {showActions}
          difficulty={selectedDifficulty || 3}
          currentQuiz={quizCache[currentStock?.dm]}
        />
      {/if}
    </main>
  </div>
</div>
{/if}

<!-- 难度选择模态框 -->
{#if showDifficultyModal}
  <Modal
    show={true}
    title="选择难度等级"
    on:close={() => showDifficultyModal = false}
    class_="simple-modal"
  >
    <div class="difficulty-modal">
      <div class="difficulty-options">
        <button
          class="difficulty-button"
          class:selected={selectedDifficulty === 1}
          on:click={() => selectedDifficulty = 1}
        >
          <div class="level-tag">L1</div>
          <div class="difficulty-content">
            <div class="difficulty-title">初学者</div>
            <div class="difficulty-desc">基础价格和成交量分析</div>
          </div>
        </button>
        
        <button
          class="difficulty-button"
          class:selected={selectedDifficulty === 2}
          on:click={() => selectedDifficulty = 2}
        >
          <div class="level-tag">L2</div>
          <div class="difficulty-content">
            <div class="difficulty-title">简单</div>
            <div class="difficulty-desc">基本技术指标分析</div>
          </div>
        </button>
        
        <button
          class="difficulty-button"
          class:selected={selectedDifficulty === 3}
          on:click={() => selectedDifficulty = 3}
        >
          <div class="level-tag">L3</div>
          <div class="difficulty-content">
            <div class="difficulty-title">中等</div>
            <div class="difficulty-desc">综合技术指标分析</div>
          </div>
        </button>
        
        <button
          class="difficulty-button"
          class:selected={selectedDifficulty === 4}
          on:click={() => selectedDifficulty = 4}
        >
          <div class="level-tag">L4</div>
          <div class="difficulty-content">
            <div class="difficulty-title">困难</div>
            <div class="difficulty-desc">市场形态与趋势分析</div>
          </div>
        </button>
        
        <button
          class="difficulty-button"
          class:selected={selectedDifficulty === 5}
          on:click={() => selectedDifficulty = 5}
        >
          <div class="level-tag">L5</div>
          <div class="difficulty-content">
            <div class="difficulty-title">专家</div>
            <div class="difficulty-desc">深度技术分析与市场研判</div>
          </div>
        </button>
      </div>
      <div class="modal-footer">
        <button 
          class="confirm-button" 
          disabled={!selectedDifficulty}
          on:click={() => {
            if (selectedDifficulty) {
              handleDifficultySelect(selectedDifficulty);
            }
          }}
        >
          开始测试
        </button>
      </div>
    </div>
  </Modal>
{/if}

<!-- 使用股票列表模态框 -->
<StockListModal
  show={showStockList}
  {stocks}
  onClose={() => {
    resetListState()
  }}
  onStart={() => {
    showStockList = false
    startStockTest()
  }}
  {loadingFirstQuiz}
  {hasFirstQuiz}
/>

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

  .score-panel {
    padding: 2rem;
    background: white;
    border-radius: 0.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }
  
  .score-list {
    margin: 1rem 0;
  }
  
  .score-item {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem;
    border-bottom: 1px solid #e5e7eb;
  }


  .difficulty-modal {
    padding: 1.5rem;
    min-width: 800px;
    background: transparent;
    width: 100%;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 2rem;
    max-width: 1000px;
  }
  
  .difficulty-options {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1rem;
    width: 100%;
  }
  
  .difficulty-button {
    display: flex;
    flex-direction: row;
    align-items: center;
    padding: 1.25rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    background: #ffffff;
    width: 100%;
    transition: all 0.2s;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    text-align: left;
    height: 100%;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
  }
  
  .level-tag {
    font-size: 0.875rem;
    font-weight: 600;
    margin-right: 1.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    width: 2.25rem;
    height: 2.25rem;
    background: #f3f4f6;
    border-radius: 0.5rem;
    color: #6b7280;
  }
  
  .difficulty-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 0.5rem;
    min-height: 3.5rem;
  }
  
  .difficulty-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: #111827;
    line-height: 1.2;
  }
  
  .difficulty-desc {
    font-size: 0.875rem;
    color: #6b7280;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  
  .difficulty-button:hover {
    border-color: #3b82f6;
    background: #f9fafb;
    transform: translateY(-1px);
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    z-index: 1;
  }
  
  .difficulty-button.selected {
    border-color: #3b82f6;
    border-width: 1px;
    background: #eff6ff;
    box-shadow: 0 0 0 1px #3b82f6, 0 2px 4px -1px rgba(59, 130, 246, 0.1);
    z-index: 2;
  }
  
  .difficulty-button.selected .level-tag {
    background: #3b82f6;
    color: #ffffff;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    padding-top: 1rem;
    border-top: 1px solid #e5e7eb;
  }
  
  .confirm-button {
    padding: 0.75rem 2rem;
    background-color: #3b82f6;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-weight: 600;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .confirm-button:hover:not(:disabled) {
    background-color: #2563eb;
    transform: translateY(-1px);
  }
  
  .confirm-button:disabled {
    background-color: #9ca3af;
    cursor: not-allowed;
    opacity: 0.7;
  }

  @media (max-width: 640px) {
    .difficulty-options {
      grid-template-columns: 1fr;
    }
    
    .difficulty-button {
      padding: 1rem;
    }
    
    .level-tag {
      width: 1.75rem;
      height: 1.75rem;
    }
    
    .confirm-button {
      width: 100%;
    }
  }
</style> 