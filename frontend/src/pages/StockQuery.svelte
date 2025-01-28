<script lang="ts">
  import { Search, Clock3, LineChart, CandlestickChart, Star, Share2, Bell, ChevronDown } from 'lucide-svelte'
  import { fade, fly, scale } from 'svelte/transition'
  import EmptyState from '../components/EmptyState.svelte'
  import StockChart from '../components/StockChart.svelte'
  import { location } from 'svelte-spa-router'
  import { onMount } from 'svelte'
  import { clickOutside } from '../directives/clickOutside'

  // 接收路由参数
  export let params: { code?: string } = {}
  
  // 组件状态
  let fullCode = ''
  let stockCode = ''
  let stockName = ''
  let stockExchange = ''

  // 监听参数变化
  $: {
    if (params.code) {
      fullCode = params.code
      const parts = fullCode.split('-')
      stockCode = parts[0] || ''
      stockName = parts[1] || ''
      stockExchange = parts[2] || ''
    }
  }

  // 图表配置
  let maType: 'SMA' | 'EMA' = 'SMA'
  let chartType: 'kline' | 'time' = 'kline'
  
  // 时间周期选项
  const timeFrames = [
    { value: '5m', label: '5分钟' },
    { value: '15m', label: '15分钟' },
    { value: '30m', label: '30分钟' },
    { value: '60m', label: '60分钟' },
    { value: 'dh', label: '日线' },
  ]
  let selectedTimeFrame = 'dh'
  let useDropdown = false

  // 监听容器宽度变化
  let chartControlsEl: HTMLElement
  
  function checkWidth() {
    if (chartControlsEl) {
      const controlsWidth = chartControlsEl.offsetWidth
      useDropdown = controlsWidth < 400
    }
  }

  onMount(() => {
    checkWidth()
    window.addEventListener('resize', checkWidth)
    return () => window.removeEventListener('resize', checkWidth)
  })

  // 收藏状态
  let isFavorite = false

  // 切换图表类型
  function toggleChartType() {
    chartType = chartType === 'kline' ? 'time' : 'kline'
  }

  // 切换收藏
  function toggleFavorite() {
    isFavorite = !isFavorite
  }

  // 下拉框状态
  let isDropdownOpen = false
  
  // 切换下拉框
  function toggleDropdown() {
    isDropdownOpen = !isDropdownOpen
  }
  
  // 选择时间周期
  function selectTimeFrame(value: string) {
    selectedTimeFrame = value
    isDropdownOpen = false
  }

  // 股票信息提示
  let showStockTip = false

  // 持股数据
  let shareData = {
    time: '',
    superShareVolume: 0,
    superShareRatio: 0,
    largeShareVolume: 0,
    largeShareRatio: 0,
    smallShareVolume: 0,
    smallShareRatio: 0,
    retailShareVolume: 0,
    retailShareRatio: 0
  };

  // 当日持股数据
  let dailyShareData = {
    superNetInflow: 0,
    largeNetInflow: 0,
    smallNetInflow: 0,
    retailNetInflow: 0,
    superInflowRatio: 0,
    largeInflowRatio: 0,
    smallInflowRatio: 0,
    retailInflowRatio: 0,
    superOutflowRatio: 0,
    largeOutflowRatio: 0,
    smallOutflowRatio: 0,
    retailOutflowRatio: 0,
    superNetRatio: 0,
    largeNetRatio: 0,
    smallNetRatio: 0,
    retailNetRatio: 0,
    superYesterdayNetRatio: 0,
    largeYesterdayNetRatio: 0,
    smallYesterdayNetRatio: 0,
    retailYesterdayNetRatio: 0
  };

  // 是否显示累计持股
  let showCumulative = true;

  // 处理持股数据更新
  function handleShareUpdate(event: CustomEvent) {
    const data = event.detail;
    console.log('Received share data:', data);
    
    // 计算总流入和总流出
    const totalInflow = data.superInflow + data.largeInflow + data.smallInflow + data.retailInflow;
    const totalOutflow = Math.abs(data.superInflow - data.superNetInflow) + 
                        Math.abs(data.largeInflow - data.largeNetInflow) + 
                        Math.abs(data.smallInflow - data.smallNetInflow) + 
                        Math.abs(data.retailInflow - data.retailNetInflow);
    
    // 计算总净流入（只计算正值）
    const positiveNetInflows = [
      data.superNetInflow,
      data.largeNetInflow,
      data.smallNetInflow,
      data.retailNetInflow
    ].filter(flow => flow > 0);
    
    const totalPositiveNetInflow = positiveNetInflows.reduce((sum, flow) => sum + flow, 0);
    
    // 计算总持股量（只计算正值）
    const totalPositiveShareVolume = [
      Math.max(0, data.superShareVolume),
      Math.max(0, data.smallShareVolume),
      Math.max(0, data.retailShareVolume)
    ].reduce((sum, volume) => sum + volume, 0);
    
    // 累计持股数据
    shareData = {
      time: data.time,
      superShareVolume: Math.max(0, data.superShareVolume),
      superShareRatio: totalPositiveShareVolume > 0 ? (Math.max(0, data.superShareVolume) / totalPositiveShareVolume * 100) : 0,
      largeShareVolume: Math.max(0, data.largeShareVolume),
      largeShareRatio: totalPositiveShareVolume > 0 ? (Math.max(0, data.largeShareVolume) / totalPositiveShareVolume * 100) : 0,
      smallShareVolume: Math.max(0, data.smallShareVolume),
      smallShareRatio: totalPositiveShareVolume > 0 ? (Math.max(0, data.smallShareVolume) / totalPositiveShareVolume * 100) : 0,
      retailShareVolume: Math.max(0, data.retailShareVolume),
      retailShareRatio: totalPositiveShareVolume > 0 ? (Math.max(0, data.retailShareVolume) / totalPositiveShareVolume * 100) : 0
    };
    
    // 当日资金流向数据
    dailyShareData = {
      // 净流入金额
      superNetInflow: data.superNetInflow || 0,
      largeNetInflow: data.largeNetInflow || 0,
      smallNetInflow: data.smallNetInflow || 0,
      retailNetInflow: data.retailNetInflow || 0,
      
      // 流入占比
      superInflowRatio: totalInflow ? (data.superInflow / totalInflow * 100) : 0,
      largeInflowRatio: totalInflow ? (data.largeInflow / totalInflow * 100) : 0,
      smallInflowRatio: totalInflow ? (data.smallInflow / totalInflow * 100) : 0,
      retailInflowRatio: totalInflow ? (data.retailInflow / totalInflow * 100) : 0,
      
      // 流出占比
      superOutflowRatio: totalOutflow ? (Math.abs(data.superInflow - data.superNetInflow) / totalOutflow * 100) : 0,
      largeOutflowRatio: totalOutflow ? (Math.abs(data.largeInflow - data.largeNetInflow) / totalOutflow * 100) : 0,
      smallOutflowRatio: totalOutflow ? (Math.abs(data.smallInflow - data.smallNetInflow) / totalOutflow * 100) : 0,
      retailOutflowRatio: totalOutflow ? (Math.abs(data.retailInflow - data.retailNetInflow) / totalOutflow * 100) : 0,
      
      // 净流入占比
      superNetRatio: totalPositiveNetInflow && data.superNetInflow > 0 ? (data.superNetInflow / totalPositiveNetInflow * 100) : 0,
      largeNetRatio: totalPositiveNetInflow && data.largeNetInflow > 0 ? (data.largeNetInflow / totalPositiveNetInflow * 100) : 0,
      smallNetRatio: totalPositiveNetInflow && data.smallNetInflow > 0 ? (data.smallNetInflow / totalPositiveNetInflow * 100) : 0,
      retailNetRatio: totalPositiveNetInflow && data.retailNetInflow > 0 ? (data.retailNetInflow / totalPositiveNetInflow * 100) : 0,
      
      // 昨日净占比
      superYesterdayNetRatio: data.superYesterdayNetRatio || 0,
      largeYesterdayNetRatio: data.largeYesterdayNetRatio || 0,
      smallYesterdayNetRatio: data.smallYesterdayNetRatio || 0,
      retailYesterdayNetRatio: data.retailYesterdayNetRatio || 0
    };

    console.log('Updated share data:', shareData);
    console.log('Updated daily share data:', dailyShareData);
  }

  // 切换显示模式
  function toggleDisplayMode() {
    showCumulative = !showCumulative;
  }

  // 格式化数字
  function formatVolume(volume: number): string {
    if (typeof volume !== 'number' || isNaN(volume)) return '-';
    const absVolume = Math.abs(volume);
    if (absVolume >= 100000000) {
      return `${(volume / 100000000).toFixed(2)}亿`;
    }
    if (absVolume >= 10000) {
      return `${(volume / 10000).toFixed(2)}万`;
    }
    return volume.toFixed(2);
  }

  // 复制股票代码
  function copyStockCode() {
    navigator.clipboard.writeText(stockCode)
      .then(() => {
        showStockTip = false
        // 可以添加一个复制成功的提示
      })
  }
</script>

<div class="page-container">
  <div class="main-container">
    <div class="layout-container" in:fade={{duration: 300}}>
      <!-- 左侧主内容区 -->
      <div class="left-section">
        <!-- 上方卡片 -->
        <div class="content-card stock-header-card" 
          in:fly={{y: -20, duration: 400, delay: 200}}>
          <!-- 左侧股票信息 -->
          <div class="stock-info" in:fly={{x: -20, duration: 400, delay: 400}}>
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div class="stock-code" on:click={copyStockCode} on:mouseenter={() => showStockTip = true} on:mouseleave={() => showStockTip = false}>
              {stockCode}
              {#if showStockTip}
                <div class="stock-tip" transition:scale={{duration: 150, start: 0.95}}>
                  点击复制股票代码
                </div>
              {/if}
            </div>
            <div class="stock-exchange" style="transform: translateY({showStockTip ? '2px' : '0'})">
              {stockExchange || 'SH'}
            </div>
            <div class="stock-name" style="transform: translateY({showStockTip ? '2px' : '0'})">
              {stockName}
            </div>
          </div>

          <!-- 分隔线 -->
          <div class="divider" in:scale={{duration: 300, delay: 300}}></div>

          <!-- 图表类型切换 -->
          <div class="chart-controls" bind:this={chartControlsEl} in:fly={{x: 20, duration: 400, delay: 500}}>
            {#if useDropdown}
              <!-- 自定义下拉选择 -->
              <div 
                class="custom-select" 
                use:clickOutside
                
              >
                <button 
                  class="select-button" 
                  on:click={toggleDropdown}
                  aria-haspopup="listbox"
                  aria-expanded={isDropdownOpen}
                >
                  <span>{timeFrames.find(f => f.value === selectedTimeFrame)?.label}</span>
                  <ChevronDown size={14} class="select-icon" style="transform: rotate({isDropdownOpen ? 180 : 0}deg)" />
                </button>
                
                {#if isDropdownOpen}
                  <div class="select-dropdown" transition:fade={{duration: 100}}>
                    {#each timeFrames as frame}
                      <button
                        class="dropdown-item"
                        class:active={selectedTimeFrame === frame.value}
                        on:click={() => selectTimeFrame(frame.value)}
                      >
                        {frame.label}
                      </button>
                    {/each}
                  </div>
                {/if}
              </div>
            {:else}
              <!-- 时间周期按钮组 -->
              <div class="time-buttons">
                {#each timeFrames as frame}
                  <button 
                    class="time-button" 
                    class:active={selectedTimeFrame === frame.value}
                    on:click={() => selectedTimeFrame = frame.value}
                  >
                    {frame.label}
                  </button>
                {/each}
              </div>
            {/if}

            <!-- K线切换按钮 -->
            <button class="icon-button" class:active={chartType === 'kline'} on:click={toggleChartType}>
              <CandlestickChart size={18} />
              <span>K线</span>
            </button>
          </div>

          <!-- 分隔线 -->
          <div class="divider" in:scale={{duration: 300, delay: 300}}></div>

          <!-- 右侧功能按钮组 -->
          <div class="action-buttons" in:fly={{x: 20, duration: 400, delay: 600}}>
            <button class="icon-button" class:active={isFavorite} on:click={toggleFavorite}>
              <Star size={18} />
            </button>
            <button class="icon-button">
              <Bell size={18} />
            </button>
            <button class="icon-button">
              <Share2 size={18} />
            </button>
          </div>
        </div>
        
        <!-- 下方卡片 -->
        <div class="content-card stock-detail-card"
          in:fly={{y: 20, duration: 400, delay: 300}}>
          <div class="card-content">
            <!-- <StockChart code={stockCode} maType={maType} on:shareUpdate={handleShareUpdate} /> -->
            <StockChart code={stockCode} freq={selectedTimeFrame}  />
          </div>
        </div>
      </div>
      
      <!-- 右侧布局 -->
      <div class="right-section">
        <!-- 右上区域 -->
        <div class="content-card share-distribution-card" in:fly={{x: 20, duration: 400, delay: 400}}>
          <div class="card-header">
            <div class="header-left">
              <h3>持股分布</h3>
              <span class="update-time">{shareData.time}</span>
            </div>
          </div>
          
          <div class="share-distribution">
            <!-- 累计持股 -->
            <div class="share-section">
              <div class="section-header">
                <span>累计持股</span>
              </div>
              <div class="share-bar-section">
                <div class="share-bar">
                  {#if shareData.superShareRatio > 0}
                    <div class="bar super positive" 
                      style="width: {shareData.superShareRatio}%" 
                      title="超大单: {shareData.superShareRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if shareData.largeShareRatio > 0}
                    <div class="bar large positive" 
                      style="width: {shareData.largeShareRatio}%" 
                      title="大单: {shareData.largeShareRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if shareData.smallShareRatio > 0}
                    <div class="bar small positive" 
                      style="width: {shareData.smallShareRatio}%" 
                      title="小单: {shareData.smallShareRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if shareData.retailShareRatio > 0}
                    <div class="bar retail positive" 
                      style="width: {shareData.retailShareRatio}%" 
                      title="散单: {shareData.retailShareRatio.toFixed(1)}%">
                    </div>
                  {/if}
                </div>
              </div>
              
              <div class="share-details">
                <div class="share-item super">
                  <span class="label">超大单</span>
                  <span class="volume">{formatVolume(shareData.superShareVolume)}元</span>
                  <span class="ratio">{shareData.superShareRatio.toFixed(1)}%</span>
                </div>
                <div class="share-item large">
                  <span class="label">大单</span>
                  <span class="volume">{formatVolume(shareData.largeShareVolume)}元</span>
                  <span class="ratio">{shareData.largeShareRatio.toFixed(1)}%</span>
                </div>
                <div class="share-item small">
                  <span class="label">小单</span>
                  <span class="volume">{formatVolume(shareData.smallShareVolume)}元</span>
                  <span class="ratio">{shareData.smallShareRatio.toFixed(1)}%</span>
                </div>
                <div class="share-item retail">
                  <span class="label">散单</span>
                  <span class="volume">{formatVolume(shareData.retailShareVolume)}元</span>
                  <span class="ratio">{shareData.retailShareRatio.toFixed(1)}%</span>
                </div>
              </div>
            </div>

            <!-- 分隔线 -->
            <div class="section-divider"></div>

            <!-- 当日资金流向 -->
            <div class="share-section">
              <div class="section-header">
                <span>当日资金流向</span>
              </div>
              
              <!-- 资金流入 -->
              <div class="share-bar-section">
                <span class="section-label">流入</span>
                <div class="share-bar">
                  {#if dailyShareData.superInflowRatio > 0}
                    <div class="bar super positive" 
                      style="width: {dailyShareData.superInflowRatio}%" 
                      title="超大单流入: {dailyShareData.superInflowRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.largeInflowRatio > 0}
                    <div class="bar large positive" 
                      style="width: {dailyShareData.largeInflowRatio}%" 
                      title="大单流入: {dailyShareData.largeInflowRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.smallInflowRatio > 0}
                    <div class="bar small positive" 
                      style="width: {dailyShareData.smallInflowRatio}%" 
                      title="小单流入: {dailyShareData.smallInflowRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.retailInflowRatio > 0}
                    <div class="bar retail positive" 
                      style="width: {dailyShareData.retailInflowRatio}%" 
                      title="散单流入: {dailyShareData.retailInflowRatio.toFixed(1)}%">
                    </div>
                  {/if}
                </div>
              </div>

              <!-- 资金流出 -->
              <div class="share-bar-section">
                <span class="section-label">流出</span>
                <div class="share-bar">
                  {#if dailyShareData.superOutflowRatio > 0}
                    <div class="bar super negative" 
                      style="width: {dailyShareData.superOutflowRatio}%" 
                      title="超大单流出: {dailyShareData.superOutflowRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.largeOutflowRatio > 0}
                    <div class="bar large negative" 
                      style="width: {dailyShareData.largeOutflowRatio}%" 
                      title="大单流出: {dailyShareData.largeOutflowRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.smallOutflowRatio > 0}
                    <div class="bar small negative" 
                      style="width: {dailyShareData.smallOutflowRatio}%" 
                      title="小单流出: {dailyShareData.smallOutflowRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.retailOutflowRatio > 0}
                    <div class="bar retail negative" 
                      style="width: {dailyShareData.retailOutflowRatio}%" 
                      title="散单流出: {dailyShareData.retailOutflowRatio.toFixed(1)}%">
                    </div>
                  {/if}
                </div>
              </div>
              
              <!-- 净流入 -->
              <div class="share-bar-section">
                <span class="section-label">净流入</span>
                <div class="share-bar">
                  {#if dailyShareData.superNetInflow > 0}
                    <div class="bar super positive" 
                      style="width: {dailyShareData.superNetRatio}%" 
                      title="超大单净流入占比: {dailyShareData.superNetRatio.toFixed(1)}%">
                    </div>
                  {:else if dailyShareData.superNetInflow < 0}
                    <div class="bar super negative" 
                      style="width: {Math.abs(dailyShareData.superNetRatio)}%" 
                      title="超大单净流出占比: {Math.abs(dailyShareData.superNetRatio).toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.largeNetInflow > 0}
                    <div class="bar large positive" 
                      style="width: {dailyShareData.largeNetRatio}%" 
                      title="大单净流入占比: {dailyShareData.largeNetRatio.toFixed(1)}%">
                    </div>
                  {:else if dailyShareData.largeNetInflow < 0}
                    <div class="bar large negative" 
                      style="width: {Math.abs(dailyShareData.largeNetRatio)}%" 
                      title="大单净流出占比: {Math.abs(dailyShareData.largeNetRatio).toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.smallNetInflow > 0}
                    <div class="bar small positive" 
                      style="width: {dailyShareData.smallNetRatio}%" 
                      title="小单净流入占比: {dailyShareData.smallNetRatio.toFixed(1)}%">
                    </div>
                  {:else if dailyShareData.smallNetInflow < 0}
                    <div class="bar small negative" 
                      style="width: {Math.abs(dailyShareData.smallNetRatio)}%" 
                      title="小单净流出占比: {Math.abs(dailyShareData.smallNetRatio).toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.retailNetInflow > 0}
                    <div class="bar retail positive" 
                      style="width: {dailyShareData.retailNetRatio}%" 
                      title="散单净流入占比: {dailyShareData.retailNetRatio.toFixed(1)}%">
                    </div>
                  {:else if dailyShareData.retailNetInflow < 0}
                    <div class="bar retail negative" 
                      style="width: {Math.abs(dailyShareData.retailNetRatio)}%" 
                      title="散单净流出占比: {Math.abs(dailyShareData.retailNetRatio).toFixed(1)}%">
                    </div>
                  {/if}
                </div>
              </div>
              
              <!-- 昨日净占比 -->
              <div class="share-bar-section">
                <span class="section-label">昨日净</span>
                <div class="share-bar">
                  {#if dailyShareData.superYesterdayNetRatio > 0}
                    <div class="bar super positive" 
                      style="width: {dailyShareData.superYesterdayNetRatio}%" 
                      title="超大单昨日净占比: {dailyShareData.superYesterdayNetRatio.toFixed(1)}%">
                    </div>
                  {:else if dailyShareData.superYesterdayNetRatio < 0}
                    <div class="bar super negative" 
                      style="width: {Math.abs(dailyShareData.superYesterdayNetRatio)}%" 
                      title="超大单昨日净占比: {dailyShareData.superYesterdayNetRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.largeYesterdayNetRatio > 0}
                    <div class="bar large positive" 
                      style="width: {dailyShareData.largeYesterdayNetRatio}%" 
                      title="大单昨日净占比: {dailyShareData.largeYesterdayNetRatio.toFixed(1)}%">
                    </div>
                  {:else if dailyShareData.largeYesterdayNetRatio < 0}
                    <div class="bar large negative" 
                      style="width: {Math.abs(dailyShareData.largeYesterdayNetRatio)}%" 
                      title="大单昨日净占比: {dailyShareData.largeYesterdayNetRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.smallYesterdayNetRatio > 0}
                    <div class="bar small positive" 
                      style="width: {dailyShareData.smallYesterdayNetRatio}%" 
                      title="小单昨日净占比: {dailyShareData.smallYesterdayNetRatio.toFixed(1)}%">
                    </div>
                  {:else if dailyShareData.smallYesterdayNetRatio < 0}
                    <div class="bar small negative" 
                      style="width: {Math.abs(dailyShareData.smallYesterdayNetRatio)}%" 
                      title="小单昨日净占比: {dailyShareData.smallYesterdayNetRatio.toFixed(1)}%">
                    </div>
                  {/if}
                  {#if dailyShareData.retailYesterdayNetRatio > 0}
                    <div class="bar retail positive" 
                      style="width: {dailyShareData.retailYesterdayNetRatio}%" 
                      title="散单昨日净占比: {dailyShareData.retailYesterdayNetRatio.toFixed(1)}%">
                    </div>
                  {:else if dailyShareData.retailYesterdayNetRatio < 0}
                    <div class="bar retail negative" 
                      style="width: {Math.abs(dailyShareData.retailYesterdayNetRatio)}%" 
                      title="散单昨日净占比: {dailyShareData.retailYesterdayNetRatio.toFixed(1)}%">
                    </div>
                  {/if}
                </div>
              </div>

              <div class="share-details">
                <div class="share-item super">
                  <span class="label">超大单</span>
                  <span class="volume" class:positive={dailyShareData.superNetInflow > 0} 
                        class:negative={dailyShareData.superNetInflow < 0}>
                    {dailyShareData.superNetInflow > 0 ? '+' : ''}{formatVolume(dailyShareData.superNetInflow)}元
                  </span>
                  <span class="ratio">
                    {dailyShareData.superNetInflow > 0 ? dailyShareData.superNetRatio.toFixed(1) : '0.0'}%
                  </span>
                </div>
                <div class="share-item large">
                  <span class="label">大单</span>
                  <span class="volume" class:positive={dailyShareData.largeNetInflow > 0} 
                        class:negative={dailyShareData.largeNetInflow < 0}>
                    {dailyShareData.largeNetInflow > 0 ? '+' : ''}{formatVolume(dailyShareData.largeNetInflow)}元
                  </span>
                  <span class="ratio">
                    {dailyShareData.largeNetInflow > 0 ? dailyShareData.largeNetRatio.toFixed(1) : '0.0'}%
                  </span>
                </div>
                <div class="share-item small">
                  <span class="label">小单</span>
                  <span class="volume" class:positive={dailyShareData.smallNetInflow > 0} 
                        class:negative={dailyShareData.smallNetInflow < 0}>
                    {dailyShareData.smallNetInflow > 0 ? '+' : ''}{formatVolume(dailyShareData.smallNetInflow)}元
                  </span>
                  <span class="ratio">
                    {dailyShareData.smallNetInflow > 0 ? dailyShareData.smallNetRatio.toFixed(1) : '0.0'}%
                  </span>
                </div>
                <div class="share-item retail">
                  <span class="label">散单</span>
                  <span class="volume" class:positive={dailyShareData.retailNetInflow > 0} 
                        class:negative={dailyShareData.retailNetInflow < 0}>
                    {dailyShareData.retailNetInflow > 0 ? '+' : ''}{formatVolume(dailyShareData.retailNetInflow)}元
                  </span>
                  <span class="ratio">
                    {dailyShareData.retailNetInflow > 0 ? dailyShareData.retailNetRatio.toFixed(1) : '0.0'}%
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 右下区域 -->
        <div class="content-card" in:fly={{x: 20, duration: 400, delay: 500}}>
          <p>右下内容</p>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  /* 全局容器 */
  .page-container {
    display: flex;
    height: 100%;
    color: #1f2937;
  }

  /* 主内容区域容器 */
  .main-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    padding: 12px;
    padding-top: 0;
    height: 100%;
  }

  /* 三分布局容器 */
  .layout-container {
    display: flex;
    gap: 12px;
    height: 100%;
    padding-top: 12px;
  }

  /* 左侧区域 */
  .left-section {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  /* 右侧区域 */
  .right-section {
    width: 280px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    height: 100%;
  }

  .right-section .content-card {
    flex: 1;
    padding: 16px;
  }

  /* 内容卡片通用样式 */
  .content-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
  }

  .content-card:hover {
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    transform: translateY(-1px);
  }

  /* 股票头部卡片 */
  .stock-header-card {
    height: 48px;
    padding: 0 16px;
    display: flex;
    align-items: center;
    gap: 16px;
  }

  /* 股票信息 */
  .stock-info {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 160px;
    position: relative;
  }

  .stock-code {
    font-size: 18px;
    font-weight: 600;
    color: #111827;
    line-height: 1;
    cursor: pointer;
    position: relative;
    padding: 4px 0;
    transition: all 0.2s;
  }

  .stock-code:hover {
    color: #2563eb;
  }

  .stock-tip {
    position: absolute;
    top: -30px;
    left: 50%;
    transform: translateX(-50%);
    background: #1f2937;
    color: white;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: normal;
    white-space: nowrap;
    pointer-events: none;
  }

  .stock-tip::after {
    content: '';
    position: absolute;
    bottom: -4px;
    left: 50%;
    transform: translateX(-50%);
    border-left: 4px solid transparent;
    border-right: 4px solid transparent;
    border-top: 4px solid #1f2937;
  }

  .stock-exchange {
    font-size: 12px;
    color: #6b7280;
    background: #f3f4f6;
    padding: 2px 6px;
    border-radius: 4px;
    font-weight: 500;
    transition: transform 0.2s ease;
  }

  .stock-name {
    font-size: 14px;
    color: #6b7280;
    line-height: 1;
    transition: transform 0.2s ease;
  }

  /* 分割线 */
  .divider {
    width: 1px;
    height: 24px;
    background-color: #e5e7eb;
    margin: 0;
    transform-origin: center;
  }

  /* 图表控制区 */
  .chart-controls {
    display: flex;
    gap: 12px;
    align-items: center;
    flex: 1;
    min-width: 0;
  }

  /* 时间按钮组 */
  .time-buttons {
    display: flex;
    gap: 1px;
    background: #f1f5f9;
    padding: 2px;
    border-radius: 6px;
  }

  .time-button {
    border: none;
    background: transparent;
    padding: 6px 12px;
    font-size: 13px;
    color: #64748b;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.2s;
    white-space: nowrap;
  }

  .time-button:hover {
    color: #334155;
  }

  .time-button.active {
    background: white;
    color: #2563eb;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  /* 自定义下拉选择框 */
  .custom-select {
    position: relative;
    min-width: 100px;
    max-width: 140px;
  }

  .select-button {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    padding: 6px 12px;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    background: white;
    color: #64748b;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .select-button:hover {
    border-color: #94a3b8;
  }

  .select-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    margin-top: 4px;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    z-index: 10;
  }

  .dropdown-item {
    display: block;
    width: 100%;
    padding: 8px 12px;
    text-align: left;
    border: none;
    background: transparent;
    color: #64748b;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .dropdown-item:hover {
    background: #f8fafc;
    color: #334155;
  }

  .dropdown-item.active {
    color: #2563eb;
    background: #eff6ff;
  }

  /* 功能按钮组 */
  .action-buttons {
    display: flex;
    gap: 8px;
    margin-left: auto;
  }

  /* 图标按钮 */
  .icon-button {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 6px 12px;
    border: none;
    background: transparent;
    color: #64748b;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .icon-button:hover {
    background: #f1f5f9;
    color: #334155;
  }

  .icon-button.active {
    color: #2563eb;
    background: #eff6ff;
  }

  .icon-button span {
    font-size: 14px;
  }

  /* 股票详情卡片 */
  .stock-detail-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 500px;
  }

  .card-content {
    padding: 0;
    flex: 1;
    height: 100%;
    overflow: hidden;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* 修复点击事件类型错误 */
  :global(.custom-select) {
    position: relative;
  }

  /* 持股分布卡片样式优化 */
  .share-distribution-card {
    padding: 12px !important;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 8px;
    border-bottom: 1px solid #e5e7eb;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .header-left h3 {
    font-size: 15px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }

  .update-time {
    font-size: 12px;
    color: #64748b;
  }

  .share-distribution {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .share-section {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 4px;
  }

  .section-header span {
    font-size: 13px;
    font-weight: 600;
    color: #334155;
  }

  .section-divider {
    height: 1px;
    background: #e5e7eb;
    margin: 8px 0;
    transition: opacity 0.3s ease;
  }

  .share-bar-section {
    display: flex;
    align-items: center;
    gap: 8px;
    margin: 4px 0;
  }

  .section-label {
    width: 42px;
    font-size: 12px;
    color: #475569;
    font-weight: 500;
  }

  .share-bar {
    flex: 1;
    display: flex;
    height: 16px;
    background: #f1f5f9;
    border-radius: 3px;
    overflow: hidden;
  }

  .share-bar .bar {
    height: 100%;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .share-bar .bar.positive {
    opacity: 0.9;
    transform-origin: left;
  }

  .share-bar .bar.negative {
    opacity: 0.75;
    transform-origin: left;
  }

  .share-details {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-top: 8px;
  }

  .share-item {
    display: flex;
    align-items: center;
    padding: 2px 0;
    font-size: 12px;
    line-height: 1.3;
    transition: all 0.3s ease;
  }

  .share-item:hover {
    background: rgba(0, 0, 0, 0.02);
  }

  .share-item .label {
    width: 42px;
    font-weight: 600;
    transition: color 0.3s ease;
  }

  .share-item .volume {
    flex: 1;
    padding: 0 8px;
    text-align: right;
    font-family: Monaco, monospace;
    transition: all 0.3s ease;
  }

  .share-item .ratio {
    min-width: 48px;
    text-align: right;
    font-weight: 500;
    font-family: Monaco, monospace;
    transition: all 0.3s ease;
  }

  .volume.positive { 
    color: #ef4444 !important; 
    font-weight: 500;
    transition: color 0.3s ease;
  }
  
  .volume.negative { 
    color: #10b981 !important; 
    font-weight: 500;
    transition: color 0.3s ease;
  }

  .share-item.super .label { color: #f97316; }
  .share-item.large .label { color: #3b82f6; }
  .share-item.small .label { color: #a855f7; }
  .share-item.retail .label { color: #14b8a6; }

  .share-bar .bar.super { 
    background: #f97316; 
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }
  .share-bar .bar.large { 
    background: #3b82f6; 
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }
  .share-bar .bar.small { 
    background: #a855f7; 
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }
  .share-bar .bar.retail { 
    background: #14b8a6; 
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .section-divider {
    height: 1px;
    background: #e5e7eb;
    margin: 8px 0;
    transition: opacity 0.3s ease;
  }

  .share-section:hover .section-divider {
    opacity: 0.8;
  }
</style> 