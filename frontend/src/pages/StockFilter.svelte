<script lang="ts">
  import { link } from 'svelte-spa-router'
  import { ListFilter, RefreshCw } from 'lucide-svelte'
  import Modal from '../components/Modal.svelte'
  import { onDestroy, onMount } from 'svelte'
  import { EventsOn, EventsOff } from '../../wailsjs/runtime'
  import { fade, fly, slide, scale } from 'svelte/transition'
  import { elasticOut, quintOut } from 'svelte/easing'
  import StrategySelector from '../components/StrategySelector.svelte'
  import ExecutionPanel from '../components/ExecutionPanel.svelte'
  import ResultsPanel from '../components/ResultsPanel.svelte'
  import EmptyState from '../components/EmptyState.svelte'
  import { toastStore } from '../stores/toast'

  type Strategy = {
    id: number
    name: string
    description: string
    filePath: string
  }

  type ExecutionStatus = {
    status: 'running' | 'paused' | 'completed' | 'error' | 'idle';
    startTime: string;
    totalStocks: number;
    processedCount: number;
    currentStock: string;
    progress: number;
    speed: number;
    estimateTime: number;
    error?: string;
  }

  type StockSignal = {
    code: string;
    name: string;
    price: number;
    turnover: number;
    change: number;
    volume: number;
    amount: number;
    reason: string;
  }

  let showStrategyModal = false
  let selectedStrategyId: number | null = null
  let strategies: Strategy[] = []
  let loading = false
  let error: string | null = null
  let selectedStrategy: Strategy | null = null
  let executing = false
  let signals: StockSignal[] = []
  let totalStocks = 0
  let matchedStocks = 0
  let executionState: ExecutionStatus = {
    status: 'idle',
    startTime: new Date().toISOString(),
    totalStocks: 0,
    processedCount: 0,
    currentStock: '',
    progress: 0,
    speed: 0,
    estimateTime: 0
  }
  let statusMessage: string | null = null
  let statusTimer: number | null = null

  // 添加动画控制变量
  let mounted = false

  // 加载策略列表
  async function loadStrategies() {
    loading = true
    error = null
    try {
      strategies = await (window as any).go.main.App.GetStrategies()
    } catch (err) {
      error = '加载策略列表失败'
    } finally {
      loading = false
    }
  }

  // 组件加载时初始化
  onMount(() => {
    // 使用 RAF 确保在下一帧执行，避免闪烁
    requestAnimationFrame(() => {
      mounted = true
    })
    
    console.log('组件挂载 - 初始状态:', { executing, selectedStrategy })
    
    // 初始化状态
    initializeExecutionState().then(() => {
      console.log('初始化完成后的状态:', { 
        executing, 
        selectedStrategy,
        executionState: executionState.status,
        hasSignals: signals.length > 0
      })
      // 如果没有正在执行的策略，才加载策略列表
      if (!executing) {
        loadStrategies()
      }
    })

    // 监听引擎错误事件
    EventsOn("engine:error", (event) => {
      const { level, message, details, component } = event;
      const fullMessage = details ? `${message}\n${details}` : message;
      
      switch (level) {
        case "info":
          toastStore.info(`[${component}] ${fullMessage}`);
          break;
        case "warning":
          toastStore.warning(`[${component}] ${fullMessage}`);
          break;
        case "error":
          toastStore.error(`[${component}] ${fullMessage}`);
          break;
        case "fatal":
          toastStore.error(`[${component}] ${fullMessage}`, {
            duration: 5000,
            icon: '💀'
          });
          break;
      }
    });

    // 监听执行状态更新
    EventsOn('engine:status', (state: ExecutionStatus) => {
      console.log('Received execution status:', state)
      if (!state) return;
      
      executionState = {
        ...executionState,
        ...state,
        startTime: state.startTime || executionState.startTime,
        totalStocks: state.totalStocks || 0,
        processedCount: state.processedCount || 0,
        progress: state.progress || 0,
        speed: state.speed || 0,
        estimateTime: state.estimateTime || 0,
        currentStock: state.currentStock || ''
      }
      // 更新总数
      totalStocks = state.totalStocks || 0
      
      // 更新执行状态
      executing = state.status === 'running' || state.status === 'paused'
      
      // 如果状态是 completed 或 error，确保清理执行状态
      if (state.status === 'completed' || state.status === 'error' || state.status === 'idle') {
        executing = false
      }
    })
    
    // 监听股票信号
    EventsOn('engine:signal', (signal: StockSignal) => {
      // 使用不可变更新方式
      signals = [...signals, signal].slice(-1000)
      // 更新匹配数
      matchedStocks = signals.length
    })

    // 返回清理函数
    return () => {
      EventsOff('execution:status')
      EventsOff('stock:signal')
      EventsOff('engine:error')
    }
  })

  // 初始化执行状态
  async function initializeExecutionState() {
    console.log('开始初始化执行状态')
    try {
      const state = await (window as any).go.main.App.GetExecutionState()
      console.log('获取到执行状态 (完整内容):', JSON.stringify(state, null, 2))
      
      if (state) {
        executionState = state
        console.log('更新执行状态:', executionState)
        
        // 如果有正在执行的策略
        if (state.status === 'running' || state.status === 'paused') {
          console.log('检测到正在执行的策略, status:', state.status)
          executing = true
          
          // 获取当前执行的策略信息
          if (state.strategyId) {
            console.log('尝试获取策略信息, strategyId:', state.strategyId)
            try {
              const strategy = await (window as any).go.main.App.GetStrategyByID(state.strategyId)
              console.log('获取到策略信息:', strategy)
              
              if (strategy) {
                selectedStrategy = strategy
                console.log('设置当前策略:', selectedStrategy)
                
                // 获取已有的执行结果
                try {
                  const results = await (window as any).go.main.App.GetExecutionResults()
                  console.log('获取到执行结果:', results)
                  
                  if (results) {
                    signals = results.signals || []
                    totalStocks = results.totalStocks || 0
                    matchedStocks = signals.length
                    console.log('更新执行结果:', { 
                      signalsCount: signals.length, 
                      totalStocks, 
                      matchedStocks 
                    })
                  }
                } catch (err) {
                  console.error('获取执行结果失败:', err)
                }
              }
            } catch (err) {
              console.error('获取策略信息失败:', err)
              executing = false // 如果获取策略失败，重置执行状态
              console.log('重置执行状态 (策略获取失败)')
            }
          } else {
            console.error('执行状态中缺少 strategyId:', state)
            executing = false // 如果没有策略ID，重置执行状态
            console.log('重置执行状态 (无策略ID)')
          }
        } else {
          console.log('没有正在执行的策略, status:', state.status)
        }
      } else {
        console.log('没有获取到执行状态')
      }
    } catch (error) {
      console.error('初始化执行状态失败:', error)
      executing = false
      console.log('重置执行状态 (初始化失败)')
    }
  }

  function showStatus(message: string) {
    statusMessage = message
    if (statusTimer) clearTimeout(statusTimer)
    statusTimer = setTimeout(() => {
      statusMessage = null
    }, 3000) as unknown as number
  }

  async function executeStrategy() {
    if (!selectedStrategy) return
    if (executing) {
      showStatus('策略执行中,请先停止策略再执行')
      return
    }
    
    // 重置状态
    signals = []
    matchedStocks = 0
    totalStocks = 0
    executing = true
    
    try {
      await (window as any).go.main.App.ExecuteStrategy(selectedStrategy.id)
    } catch (err) {
      console.error('执行策略失败:', err)
      executing = false
      showStatus('执行策略失败')
    }
  }

  // 暂停执行
  async function pauseExecution() {
    try {
      await (window as any).go.main.App.PauseExecution()
    } catch (err) {
      console.error('暂停执行失败:', err)
    }
  }

  // 恢复执行
  async function resumeExecution() {
    try {
      await (window as any).go.main.App.ResumeExecution()
    } catch (err) {
      console.error('恢复执行失败:', err)
    }
  }

  // 停止执行
  async function stopExecution() {
    try {
      await (window as any).go.main.App.StopExecution()
      executing = false
      // 不需要手动设置状态，让后端的状态更新事件来处理
    } catch (err) {
      console.error('停止执行失败:', err)
      showStatus('停止执行失败')
    }
  }

  function handleStrategyConfirm(event: CustomEvent<Strategy>) {
    selectedStrategy = event.detail
    showStrategyModal = false
    selectedStrategyId = null
  }
</script>

{#if mounted}
<div class="page-container" in:fade={{duration: 300}}>
  {#if statusMessage}
    <div class="status-message" transition:fade={{duration: 200}}>
      <svg viewBox="0 0 24 24" width="16" height="16">
        <circle cx="12" cy="12" r="10" stroke="currentColor" fill="none" stroke-width="2"/>
        <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2"/>
        <line x1="12" y1="16" x2="12.01" y2="16" stroke="currentColor" stroke-width="2"/>
      </svg>
      <span>{statusMessage}</span>
    </div>
  {/if}

  <div class="main-container">
    <header class="header" in:fade={{duration: 200}}>
      <div class="header-content">
        <div class="header-left" in:fly={{x: -20, duration: 300}}>
          <div class="title-group">
            <h1 in:fly={{x: -20, duration: 300}}>股票筛选</h1>
            <div class="divider" in:fade={{duration: 200, delay: 100}}></div>
            <span class="version" in:fly={{x: -10, duration: 300, delay: 150}}>Beta</span>
          </div>
          {#if selectedStrategy}
            <div class="divider vertical" in:scale={{duration: 200, delay: 400}}></div>
            <div class="strategy-name" in:fly={{x: -20, duration: 400, delay: 500}}>
              <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" fill="none">
                <path d="M3 12h4l3-9 4 18 3-9h4" stroke-width="1.5"/>
              </svg>
              {selectedStrategy.name}
            </div>
          {/if}
        </div>
        <div class="header-right" in:fly={{x: 20, duration: 400, delay: 400}}>
          {#if selectedStrategy}
            <button class="btn text" on:click={() => showStrategyModal = true}>
              切换策略
            </button>
            <div class="divider vertical" in:scale={{duration: 200}}></div>
            <button class="btn solid" on:click={executeStrategy} disabled={executing}>
              开始执行
            </button>
          {:else}
            <button class="btn solid" on:click={() => showStrategyModal = true}>
              选择策略
            </button>
          {/if}
        </div>
      </div>
    </header>

    <main class="main" in:fly={{y: 20, duration: 400, delay: 300}}>
      <div class="content-wrapper">
        {#if selectedStrategy}
          {#if executionState && (executionState.status === 'running' || executionState.status === 'paused')}
            <ExecutionPanel 
              {executionState}
              on:pause={pauseExecution}
              on:resume={resumeExecution}
              on:stop={stopExecution}
            />
          {:else}
            <section class="ready-state" in:fly={{y: 20, duration: 400, easing: quintOut}}>
              <div class="ready-content">
                <div class="strategy-preview">
                  <div class="strategy-icon pulse" in:scale={{duration: 400, delay: 200, easing: elasticOut}}>
                    <svg viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" fill="none">
                      {#if selectedStrategy.id === 1}
                        <path d="M3 12h4l3-9 4 18 3-9h4" stroke-width="1.5"/>
                      {:else if selectedStrategy.id === 2}
                        <path d="M3 3v18h18M7 12v5M11 8v9M15 11v6M19 8v9" stroke-width="1.5"/>
                      {:else}
                        <path d="M3 10h18M7 15c4 0 8-10 12 0" stroke-width="1.5"/>
                      {/if}
                    </svg>
                  </div>
                  <div class="strategy-info" in:fly={{y: 20, duration: 400, delay: 300}}>
                    <h3>{selectedStrategy.name}</h3>
                    <p>{selectedStrategy.description}</p>
                  </div>
                </div>
                <button class="btn solid large" 
                  on:click={executeStrategy}
                  in:scale={{duration: 400, delay: 400, easing: elasticOut}}
                >
                  <svg class="start-icon" viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" fill="none">
                    <path d="M5 3l14 9-14 9V3z" stroke-width="1.5"/>
                  </svg>
                  开始执行
                </button>
              </div>
            </section>
          {/if}

          {#if signals.length > 0}
            <ResultsPanel 
              stockSignals={signals}
              {totalStocks}
              {matchedStocks}
            />
          {/if}
        {:else}
          <EmptyState on:select={() => showStrategyModal = true} />
        {/if}
      </div>
    </main>

    <Modal 
      bind:show={showStrategyModal}
      title="选择策略"
    >
      <StrategySelector
        {strategies}
        {selectedStrategyId}
        {loading}
        {error}
        on:confirm={handleStrategyConfirm}
        on:cancel={() => showStrategyModal = false}
        on:retry={loadStrategies}
      />
    </Modal>
  </div>
</div>
{/if}

<style>
  /* 全局容器 */
  .page-container {
    display: flex;
    height: 100%;
    color: #1f2937;
    /* overflow: hidden; */
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  /* 主内容区域容器 */
  .main-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: 24px;
    padding-top: 0;
    /* overflow: hidden; */
  }

  /* 顶部导航栏 */
  .header {
    height: 60px;
    border-bottom: 1px solid #f3f4f6;
    flex-shrink: 0;
    animation: fadeIn 0.2s ease-out;
  }

  .header-content {
    height: 100%;
    max-width: 1440px;
    margin: 0 auto;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .header-left {
    display: flex;
    align-items: center;
    animation: slideInLeft 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .title-group {
    display: flex;
    align-items: center;
    transform-origin: left;
    animation: slideInLeft 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) 0.1s backwards;
  }

  h1 {
    font-size: 18px;
    font-weight: 600;
    color: #111827;
  }

  .version {
    font-size: 13px;
    font-weight: 500;
    color: #2563eb;
  }

  .divider {
    width: 1px;
    height: 12px;
    background: #e5e7eb;
    margin: 0 12px;
    transform-origin: center;
  }

  .divider.vertical {
    height: 16px;
    margin: 0 20px;
  }

  .strategy-name {
    font-size: 14px;
    color: #4b5563;
    display: flex;
    align-items: center;
    gap: 8px;
    transform-origin: left;
  }

  .header-right {
    display: flex;
    align-items: center;
    transform-origin: right;
  }

  .btn {
    height: 34px;
    padding: 0 16px;
    font-size: 14px;
    font-weight: 500;
    border-radius: 6px;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .text {
    color: #4b5563;
    background: transparent;
    border: none;
  }

  .text:hover {
    color: #111827;
    transform: translateY(-1px);
  }

  .solid {
    color: white;
    background: #2563eb;
    border: none;
  }

  .solid:hover:not(:disabled) {
    background: #1d4ed8;
    transform: translateY(-1px);
  }

  .solid:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* 主内容区 */
  .main {
    flex: 1;
    width: 100%;
    display: flex;
    flex-direction: column;
    margin-top: 24px;
  }

  .content-wrapper {
    width: 100%;
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  /* 准备状态 */
  .ready-state {
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    padding: 32px;
    flex: 1;
    background: linear-gradient(to right, #f8fafc, #f0f7ff);
    min-height: 400px;
    display: flex;
    align-items: center;
    justify-content: center;
    transform-origin: center;
  }

  .ready-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 32px;
    text-align: center;
    width: 100%;
    max-width: 480px;
    margin: 0 auto;
  }

  .strategy-preview {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
  }

  .strategy-icon {
    width: 64px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #2563eb;
    background: white;
    border: 1px solid #93c5fd;
    border-radius: 16px;
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.1);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    transform-origin: center;
  }

  .strategy-icon.pulse {
    animation: pulse 2s infinite;
  }

  .strategy-info {
    transform-origin: center;
  }

  .strategy-info h3 {
    font-size: 18px;
    font-weight: 600;
    color: #111827;
    margin: 0 0 8px;
  }

  .strategy-info p {
    font-size: 14px;
    color: #6b7280;
    margin: 0;
    max-width: 480px;
  }

  .btn.solid.large {
    display: flex;
    align-items: center;
    gap: 8px;
    height: 44px;
    padding: 0 32px;
    font-size: 15px;
    font-weight: 600;
    background: linear-gradient(135deg, #2563eb, #1d4ed8);
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(37, 99, 235, 0.25);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    transform-origin: center;
  }

  .btn.solid.large:hover {
    transform: translateY(-2px) scale(1.02);
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
  }

  .btn.solid.large:active {
    transform: translateY(0) scale(0.98);
  }

  .start-icon {
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .btn.solid.large:hover .start-icon {
    transform: translateX(4px);
  }

  @keyframes pulse {
    0% {
      box-shadow: 0 4px 12px rgba(37, 99, 235, 0.1);
      transform: scale(1);
    }
    50% {
      box-shadow: 0 4px 24px rgba(37, 99, 235, 0.2);
      transform: scale(1.05);
    }
    100% {
      box-shadow: 0 4px 12px rgba(37, 99, 235, 0.1);
      transform: scale(1);
    }
  }

  /* 状态消息 */
  .status-message {
    position: fixed;
    top: 24px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 20px;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    z-index: 1000;
    color: #4b5563;
    animation: slideDown 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translate(-50%, -20px);
    }
    to {
      opacity: 1;
      transform: translate(-50%, 0);
    }
  }

  :global(.strategy-modal) {
    min-width: 600px !important;
    max-width: 80vw !important;
    animation: modalShow 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes modalShow {
    from {
      opacity: 0;
      transform: scale(0.95);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  /* 动画关键帧 */
  @keyframes slideInLeft {
    from {
      opacity: 0;
      transform: translateX(-20px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  @keyframes slideInRight {
    from {
      opacity: 0;
      transform: translateX(20px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  /* 策略卡片动画 */
  .strategy-card {
    transform-origin: center;
    transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .strategy-card:hover {
    transform: translateY(-2px) scale(1.02);
  }

  .strategy-card:active {
    transform: scale(0.98);
  }

  .strategy-icon {
    transform-origin: center;
    transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .strategy-card:hover .strategy-icon {
    transform: scale(1.1);
  }
</style> 

