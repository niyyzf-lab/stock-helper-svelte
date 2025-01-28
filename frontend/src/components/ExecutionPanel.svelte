<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import { fade } from 'svelte/transition'
  import toast from 'svelte-french-toast'
  import { onMount, onDestroy } from 'svelte'

  interface ExecutionStatus {
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

  export let executionState: ExecutionStatus

  const dispatch = createEventDispatcher()
  let toastId: string;
  let lastProgress = 0;
  let progressInterval: any;
  let panelElement: HTMLElement;
  let isVisible = false;
  let mainContainer: HTMLElement;

  // 检查元素是否在视口中
  function checkVisibility() {
    if (!panelElement) return false;
    const rect = panelElement.getBoundingClientRect();
    const windowHeight = window.innerHeight || document.documentElement.clientHeight;
    return rect.top >= 0 && rect.bottom <= windowHeight;
  }

  // 滚动到结果面板
  function scrollToResults() {
    const mainContent = document.querySelector('.main-container');
    if (mainContent) {
      mainContent.scrollTo({
        top: mainContent.scrollHeight,
        behavior: 'smooth'
      });
    }
  }

  // 滚动到执行面板
  function scrollToExecution() {
    const mainContent = document.querySelector('.main-container');
    if (mainContent) {
      mainContent.scrollTo({
        top: 0,
        behavior: 'smooth'
      });
    }
    if (toastId) {
      toast.dismiss(toastId);
      toastId = '';
    }
  }

  // 格式化时间
  function formatTime(seconds: number): string {
    if (!seconds || seconds < 0) return '0秒'
    if (seconds < 60) return `${Math.floor(seconds)}秒`
    if (seconds < 3600) return `${Math.floor(seconds / 60)}分${Math.floor(seconds % 60)}秒`
    return `${Math.floor(seconds / 3600)}小时${Math.floor((seconds % 3600) / 60)}分`
  }

  // 格式化速度
  function formatSpeed(speed: number): string {
    if (!speed || speed < 0) return '0.0'
    return speed.toFixed(1)
  }

  // 监听滚动事件
  function handleScroll() {
    if (status === 'running') {
      updateToastProgress();
    }
  }

  // 计算已运行时间
  $: elapsedTime = executionState?.startTime ? 
    Math.floor((Date.now() - new Date(executionState.startTime).getTime()) / 1000) : 0

  // 安全获取进度
  $: progress = typeof executionState?.progress === 'number' ? executionState.progress : 0
  $: processedCount = typeof executionState?.processedCount === 'number' ? executionState.processedCount : 0
  $: totalStocks = typeof executionState?.totalStocks === 'number' ? executionState.totalStocks : 0
  $: speed = typeof executionState?.speed === 'number' ? executionState.speed : 0
  $: estimateTime = typeof executionState?.estimateTime === 'number' ? executionState.estimateTime : 0
  $: currentStock = executionState?.currentStock || ''
  $: status = executionState?.status || 'idle'

  // 监听状态变化
  $: {
    if (status === 'running' && progress > lastProgress) {
      updateToastProgress();
    } else if (status !== 'running' && toastId) {
      // 处理非运行状态
      if (status === 'completed' || status === 'error' || status === 'idle') {
        updateToastProgress();
      }
    }
  }

  onMount(() => {
    // 启动定时更新
    progressInterval = setInterval(() => {
      if (status === 'running') {
        updateToastProgress();
      }
    }, 1000);

    // 添加滚动监听
    window.addEventListener('scroll', handleScroll, { passive: true });
    
    // 添加 toast 点击监听
    document.addEventListener('click', (e) => {
      if ((e.target as HTMLElement).closest('.toast-clickable')) {
        scrollToExecution();
      }
    });
  });

  onDestroy(() => {
    // 清理定时器
    if (progressInterval) {
      clearInterval(progressInterval);
    }
    // 清理 toast
    if (toastId) {
      toast.dismiss(toastId);
    }
    // 移除滚动监听
    window.removeEventListener('scroll', handleScroll);
  });

  // 调试输出
  $: console.log('ExecutionPanel state:', { 
    progress, 
    processedCount, 
    totalStocks, 
    speed, 
    estimateTime, 
    currentStock, 
    status 
  })

  // 更新 toast 配置
  const toastConfig = {
    position: 'bottom-right' as const,
    duration: Infinity,
    style: 'min-width: 300px; cursor: pointer;',
    className: 'toast-clickable',
    ariaProps: {
      role: 'alert' as const,
      'aria-live': 'polite' as const
    },
    onClick: () => {
      scrollToExecution();
    }
  };

  // 更新 Toast 进度
  function updateToastProgress() {
    // 检查面板是否在视口中
    isVisible = checkVisibility();
    
    // 如果面板可见，不显示 toast
    if (isVisible) {
      if (toastId) {
        toast.dismiss(toastId);
        toastId = '';
      }
      return;
    }

    if (status === 'running' && progress > lastProgress) {
      if (!toastId) {
        toastId = toast.loading(
          `正在处理: ${currentStock || '准备中...'}`, 
          toastConfig
        );
      } else {
        toast.loading(
          `正在处理: ${currentStock || '准备中...'} (${progress.toFixed(1)}%)`,
          { ...toastConfig, id: toastId }
        );
      }
      lastProgress = progress;
    }

    // 处理完成或出错时
    if (status !== 'running' && toastId) {
      if (status === 'completed') {
        toast.success(
          `处理完成: 共处理 ${processedCount}/${totalStocks} 只股票`,
          { ...toastConfig, duration: 3000 }
        );
      } else if (status === 'error') {
        toast.error(
          `处理出错: ${executionState.error || '未知错误'}`,
          { ...toastConfig, duration: 5000 }
        );
      } else if (status === 'idle') {
        toast.success(
          `已停止: 共处理 ${processedCount}/${totalStocks} 只股票`,
          { ...toastConfig, duration: 3000 }
        );
      }
      toastId = '';
      lastProgress = 0;
    }
  }
</script>

{#if status === 'running' || status === 'paused'}
  <section class="execution-panel" bind:this={panelElement} in:fade={{duration: 200}}>
    <div class="status-bar">
      <div class="status-info">
        <div class="status-badge" class:paused={status === 'paused'}>
          {#if status === 'running'}
            <div class="pulse"></div>
            <span>执行中</span>
          {:else if status === 'paused'}
            <svg viewBox="0 0 24 24" width="14" height="14" stroke="currentColor" fill="none">
              <rect x="6" y="4" width="4" height="16" stroke="none" fill="currentColor"/>
              <rect x="14" y="4" width="4" height="16" stroke="none" fill="currentColor"/>
            </svg>
            <span>已暂停</span>
          {/if}
        </div>
        <div class="divider"></div>
        <div class="metrics">
          <div class="metric">
            <span class="label">已处理</span>
            <span class="value">{processedCount}/{totalStocks}</span>
          </div>
          <div class="metric">
            <span class="label">速度</span>
            <span class="value">{formatSpeed(speed)}个/秒</span>
          </div>
          <div class="metric">
            <span class="label">已用时间</span>
            <span class="value">{formatTime(elapsedTime)}</span>
          </div>
          {#if estimateTime > 0}
            <div class="metric">
              <span class="label">预计剩余</span>
              <span class="value">{formatTime(estimateTime)}</span>
            </div>
          {/if}
        </div>
      </div>
      <div class="controls">
        {#if status === 'running'}
          <button class="btn" on:click={() => dispatch('pause')}>
            <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" fill="none">
              <rect x="6" y="4" width="4" height="16" stroke="none" fill="currentColor"/>
              <rect x="14" y="4" width="4" height="16" stroke="none" fill="currentColor"/>
            </svg>
            暂停
          </button>
        {:else if status === 'paused'}
          <button class="btn" on:click={() => dispatch('resume')}>
            <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" fill="none">
              <path d="M5 3l14 9-14 9V3z" stroke="none" fill="currentColor"/>
            </svg>
            继续
          </button>
        {/if}
        <button class="btn stop" on:click={() => dispatch('stop')}>
          <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" fill="none">
            <rect x="4" y="4" width="16" height="16" stroke="none" fill="currentColor"/>
          </svg>
          停止
        </button>
      </div>
    </div>

    <div class="progress-bar">
      <div class="progress" style="width: {progress}%"></div>
    </div>

    <div class="current-stock">
      <span class="label">当前处理：</span>
      <span class="stock">{currentStock || '准备中...'}</span>
    </div>

    <button class="scroll-btn" on:click={scrollToResults}>
      <div class="scroll-icon">
        <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" fill="none" stroke-width="2">
          <path d="M12 5v14M5 12l7 7 7-7" />
        </svg>
        <span>查看结果</span>
      </div>
    </button>
  </section>
{/if}

<style>
  .execution-panel {
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    padding: 16px;
    background: white;
    position: relative;
  }

  .status-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .status-info {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .status-badge {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 10px;
    background: #ecfdf5;
    color: #059669;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 500;
  }

  .status-badge.paused {
    background: #f3f4f6;
    color: #6b7280;
  }

  .pulse {
    width: 8px;
    height: 8px;
    background: currentColor;
    border-radius: 50%;
    position: relative;
  }

  .pulse::after {
    content: '';
    position: absolute;
    width: 100%;
    height: 100%;
    background: currentColor;
    border-radius: 50%;
    animation: pulse 2s infinite;
  }

  .divider {
    width: 1px;
    height: 24px;
    background: #e5e7eb;
  }

  .metrics {
    display: flex;
    gap: 24px;
  }

  .metric {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .label {
    font-size: 13px;
    color: #6b7280;
  }

  .value {
    font-size: 13px;
    color: #111827;
    font-weight: 500;
  }

  .controls {
    display: flex;
    gap: 8px;
  }

  .btn {
    height: 32px;
    padding: 0 12px;
    display: flex;
    align-items: center;
    gap: 6px;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    background: white;
    color: #374151;
    font-size: 13px;
    font-weight: 500;
    transition: all 0.15s;
  }

  .btn:hover {
    background: #f9fafb;
    border-color: #d1d5db;
  }

  .btn.stop {
    color: #dc2626;
  }

  .btn.stop:hover {
    background: #fef2f2;
    border-color: #fecaca;
  }

  .progress-bar {
    height: 4px;
    background: #f3f4f6;
    border-radius: 2px;
    overflow: hidden;
    margin-bottom: 12px;
  }

  .progress {
    height: 100%;
    background: #2563eb;
    border-radius: 2px;
    transition: width 0.3s ease;
  }

  .current-stock {
    font-size: 13px;
  }

  .current-stock .stock {
    color: #111827;
    font-weight: 500;
  }

  @keyframes pulse {
    0% {
      transform: scale(1);
      opacity: 1;
    }
    100% {
      transform: scale(2.5);
      opacity: 0;
    }
  }

  .scroll-btn {
    position: absolute;
    bottom: 12px;
    right: 12px;
    border: none;
    background: transparent;
    cursor: pointer;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
  }

  .scroll-icon {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    border-radius: 16px;
    color: #2563eb;
    font-size: 13px;
    font-weight: 500;
    background: rgba(37, 99, 235, 0.1);
    transition: all 0.3s ease;
  }

  .scroll-icon span {
    opacity: 0;
    transform: translateX(-10px);
    transition: all 0.3s ease;
  }

  .scroll-btn:hover .scroll-icon {
    background: rgba(37, 99, 235, 0.15);
    padding-right: 16px;
  }

  .scroll-btn:hover .scroll-icon span {
    opacity: 1;
    transform: translateX(0);
  }

  .scroll-btn:hover svg {
    transform: translateY(2px);
  }

  .scroll-btn svg {
    transition: transform 0.3s ease;
  }

  .scroll-btn:hover .scroll-icon {
    animation: none;
  }
</style> 