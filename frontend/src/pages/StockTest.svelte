<script lang="ts">
  import { fade, fly } from 'svelte/transition'
  import { quintOut } from 'svelte/easing'
  import { onMount, onDestroy } from 'svelte'
  import { Target } from 'lucide-svelte'
  import { toastStore } from '../stores/toast'
  import Modal from '../components/Modal.svelte'
  import EmptyState from '../components/stock-test/EmptyState.svelte'
  import StockListModal from '../components/stock-test/StockListModal.svelte'
  import ResultModal from '../components/stock-test/ResultModal.svelte'
  import TestView from '../components/stock-test/TestView.svelte'
  import type { TestResult } from 'src/types/stock';
  import StockChart from '../components/StockChart.svelte'
  import TestResultPanel from '../components/stock-test/TestResultPanel.svelte'

  // 添加动画控制变量
  let mounted = false
  let showStockList = false
  let loading = false
  let stocks: any[] = []
  let selectedCount = 1// 固定15只股票
  
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

  let retryCount = 0
  const MAX_RETRIES = 3

  // 在 script 标签中添加新的状态变量
  let showResultModal = false
  let currentResult:TestResult| null = null

  // 添加全屏状态变量
  let isFullscreen = false

  // 添加新的状态变量
  let showFullChartModal = false

  // 添加动画控制变量
  let currentStep = 0
  let animationTimer: number

  // 添加测试结果数组
  let testResults: TestResult[] = []

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
            const oneYearData = dayData.slice(-252) // 约等于一年的交易日
            
            // 确保有足够的数据进行测试
            if (oneYearData.length >= 46) {
              // 从最近3个月的数据中随机选择一个日期作为测试日期
              const availableData = oneYearData.slice(-90, -15) // 最近90天，排除最后15天
              const randomIndex = Math.floor(Math.random() * availableData.length)
              const selectedData = availableData[randomIndex]
              const selectedIndex = oneYearData.indexOf(selectedData)
              
              // 获取历史数据(前30天)和未来数据(15天)
              const historyData = oneYearData.slice(selectedIndex - 30, selectedIndex)
              const futureData = oneYearData.slice(selectedIndex + 1, selectedIndex + 16)
              
              if (historyData.length === 30 && futureData.length === 15) {
                stock.testDate = selectedData.d
                stock.price = selectedData.c
                stock.historyData = historyData
                stock.futureData = futureData
                validStockData.push(stock)
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
    // 重置所有状态
    showChart = false
    currentStock = null
    currentStockIndex = 0
    currentResult = null
    showActions = false
    testResults = []
    // 获取新的股票数据
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
  }

  // 处理用户答案
  async function handleAnswer(direction: 'up' | 'down' | 'shock') {
    showActions = false
    
    try {
      const futureData = currentStock.futureData
      
      if (futureData && futureData.length > 0) {
        const startPrice = futureData[0].c
        const endPrice = futureData[futureData.length - 1].c
        
        // 计算区间最高最低价
        const maxPrice = Math.max(...futureData.map((d: { h: number }) => d.h))
        const minPrice = Math.min(...futureData.map((d: { l: number }) => d.l))
        
        // 计算涨跌幅
        const priceChange = (endPrice - startPrice) / startPrice
        
        // 判断实际走势
        const isShock = Math.abs(priceChange) <= 0.03
        const isUp = priceChange > 0.03
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
          priceChange,
          actualDirection,
          daysCount: futureData.length,
          klineData: currentStock.historyData,
          futureData: futureData,
          prices: futureData.map((d: { c: number }) => d.c) // 添加收盘价数组
        }
        
        // 保存当前股票的结果
        currentStock.result = currentResult
        
        // 显示结果模态框
        showResultModal = true
      }
    } catch (err) {
      console.error('处理结果数据失败:', err)
      toastStore.error('处理结果数据失败')
      moveToNextStock()
    }
  }

  // 抽取移动到下一只股票的逻辑
  function moveToNextStock() {
    currentStockIndex++
    if (currentStockIndex < stocks.length) {
      currentStock = stocks[currentStockIndex]
      currentResult = null
      showActions = false
      setTimeout(() => {
        showActions = true
      }, 5000)
    } else {
      // 测试结束
      showChart = false
      testResults = stocks.map((stock, index) => ({
        ...stock.result,
        stockCode: stock.dm,
        stockName: stock.mc
      }))
    }
  }

  // 关闭结果模态框并继续
  function continueTest() {
    if (animationTimer) {
      clearInterval(animationTimer)
    }
    currentStep = 0
    showResultModal = false
    moveToNextStock()
  }

  // 添加全屏切换函数
  function toggleFullscreen() {
    isFullscreen = !isFullscreen
  }

  // 添加显示完整图表的函数
  function showFullChart() {
    showFullChartModal = true
  }

  $: if (showResultModal && currentResult) {
    startAnimation()
  }

  function startAnimation() {
    currentStep = 0
    if (animationTimer) clearInterval(animationTimer)
    
    animationTimer = window.setInterval(() => {
      if (currentStep >= (currentResult?.prices?.length || 0) - 1) {
        clearInterval(animationTimer)
        return
      }
      currentStep++
    }, 300)
  }

  // 在组件销毁时清理定时器
  onDestroy(() => {
    if (animationTimer) {
      clearInterval(animationTimer)
    }
  })

  $: chartOptions = {
    grid: {
      top: 40,
      right: 20,
      bottom: 40,
      left: 50,
    },
    xAxis: {
      type: 'category',
      data: Array.from({ length: currentResult?.prices?.length || 0 }, (_, i) => `Day ${i + 1}`),
      axisLine: { lineStyle: { color: '#e2e8f0' } },
      axisTick: { show: false },
      axisLabel: { color: '#64748b' }
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: '#e2e8f0', type: 'dashed' } },
      axisLabel: { 
        color: '#64748b',
        formatter: (value: number) => `${value.toFixed(2)}%`
      }
    },
    series: [{
      type: 'line',
      data: currentResult?.prices?.slice(0, currentStep + 1)?.map(price => (
        (price - (currentResult?.prices?.[0] || 0)) / (currentResult?.prices?.[0] || 1) * 100
      )) || [],
      showSymbol: false,
      lineStyle: {
        width: 3,
        color: '#3b82f6'
      },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [{
            offset: 0,
            color: 'rgba(59, 130, 246, 0.2)'
          }, {
            offset: 1,
            color: 'rgba(59, 130, 246, 0)'
          }]
        }
      },
      emphasis: {
        focus: 'series',
        itemStyle: {
          color: '#2563eb'
        }
      },
      animation: true,
      animationDuration: 300,
      animationEasing: 'cubicOut'
    }],
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const value = params[0].value
        return `Day ${params[0].dataIndex + 1}<br/>涨跌幅：${value?.toFixed(2)}%`
      },
      backgroundColor: 'rgba(0, 0, 0, 0.75)',
      borderWidth: 0,
      textStyle: {
        color: '#fff'
      }
    }
  }
</script>

{#if mounted}
<div class="page-container" in:fade={{duration: 300}}>
  <div class="main-container">
    <main class="main" in:fly={{y: 20, duration: 400, delay: 300}}>
      {#if !showChart && testResults.length === 0}
        <div class="content-wrapper">
          <EmptyState
            {loading}
            {progress}
            onStart={startTest}
          />
        </div>
      {:else if !showChart && testResults.length > 0}
        <div class="content-wrapper">
          <TestResultPanel 
            results={testResults}
            on:click={() => {
              testResults = []
              startTest()
            }}
          />
        </div>
      {:else}
        <TestView
          {currentStock}
          {currentStockIndex}
          {stocks}
          {showActions}
          onAnswer={handleAnswer}
        />
      {/if}
    </main>
  </div>
</div>
{/if}

<!-- 使用股票列表模态框 -->
<StockListModal
  show={showStockList}
  stocks={stocks}
  onClose={() => showStockList = false}
  onStart={() => {
    showStockList = false
    startStockTest()
  }}
/>

<!-- 修改结果模态框部分 -->
{#if showResultModal}
  <ResultModal
    show={true}
    result={currentResult}
    onClose={continueTest}
    onViewChart={() => {
      showFullChartModal = true
    }}
  />
{/if}

<!-- 修改完整图表模态框 -->
{#if showFullChartModal}
  <Modal 
    show={showFullChartModal}
    title={`${currentStock.mc} (${currentStock.dm}) 走势图`}
    on:close={() => showFullChartModal = false}
    class_="chart-modal"
  >
    <div class="chart-modal-content">
      <StockChart 
        code={currentStock.dm}
        endTime={(() => {
          // 计算结束日期：测试日期后15天
          const testDate = new Date(currentStock.testDate)
          testDate.setDate(testDate.getDate() + 15)
          return testDate.toISOString().split('T')[0]
        })()}
        freq="dh"
      />
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


  /* 修改图表模态框相关样式 */
  .chart-modal-content {
    width: 100%;
    height: 100%;
    min-height: 600px;
    padding: 0;  /* 移除内边距 */
  }
</style> 