<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import { fade } from 'svelte/transition'
  import toast from 'svelte-french-toast'

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
  let panelElement: HTMLElement;
  let isVisible = false;

  // 检查元素是否在视口中
  function checkVisibility() {
    if (!panelElement) return false;
    const rect = panelElement.getBoundingClientRect();
    const windowHeight = window.innerHeight || document.documentElement.clientHeight;
    return rect.top >= 0 && rect.bottom <= windowHeight;
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

  // 安全获取进度
  $: progress = typeof executionState?.progress === 'number' ? executionState.progress : 0
  $: processedCount = typeof executionState?.processedCount === 'number' ? executionState.processedCount : 0
  $: totalStocks = typeof executionState?.totalStocks === 'number' ? executionState.totalStocks : 0
  $: speed = typeof executionState?.speed === 'number' ? executionState.speed : 0
  $: estimateTime = typeof executionState?.estimateTime === 'number' ? executionState.estimateTime : 0
  $: currentStock = executionState?.currentStock || ''
  $: status = executionState?.status || 'idle'
  $: elapsedTime = executionState?.startTime ? 
    Math.floor((Date.now() - new Date(executionState.startTime).getTime()) / 1000) : 0

  // 更新 Toast 进度
  function updateToastProgress() {
    isVisible = checkVisibility();
    
    if (isVisible) {
      if (toastId) {
        toast.dismiss(toastId);
        toastId = '';
      }
      return;
    }

    if (status === 'running' && progress >= 0) {
      if (!toastId) {
        toastId = toast.loading(
          `正在处理: ${currentStock || '准备中...'}`, 
          {
            position: 'bottom-right',
            duration: Infinity,
            style: 'min-width: 300px; cursor: pointer;',
            className: 'toast-clickable'
          }
        );

        // 添加点击事件监听
        const toastElement = document.querySelector('.toast-clickable');
        if (toastElement) {
          toastElement.addEventListener('click', scrollToExecution);
        }
      } else {
        toast.loading(
          `正在处理: ${currentStock || '准备中...'} (${progress.toFixed(1)}%)`,
          { id: toastId }
        );
      }
    }

    if (status !== 'running' && toastId) {
      if (status === 'completed') {
        toast.success(
          `处理完成: 共处理 ${processedCount}/${totalStocks} 只股票`,
          { id: toastId, duration: 3000 }
        );
      } else if (status === 'error') {
        toast.error(
          `处理出错: ${executionState.error || '未知错误'}`,
          { id: toastId, duration: 5000 }
        );
      } else if (status === 'idle') {
        toast.success(
          `已停止: 共处理 ${processedCount}/${totalStocks} 只股票`,
          { id: toastId, duration: 3000 }
        );
      }
      toastId = '';
    }
  }

  $: {
    if (executionState) {
      updateToastProgress();
    }
  }
</script>

{#if status === 'running' || status === 'paused'}
  <section class="execution-panel" bind:this={panelElement} in:fade={{duration: 200}}>
    <div class="status-header">
      <div class="status-info">
        <div class="status-badge" class:running={status === 'running'} class:paused={status === 'paused'}>
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
          <button class="btn pause" on:click={() => dispatch('pause')}>
            <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" fill="none">
              <rect x="6" y="4" width="4" height="16" stroke="none" fill="currentColor"/>
              <rect x="14" y="4" width="4" height="16" stroke="none" fill="currentColor"/>
            </svg>
            暂停
          </button>
        {:else if status === 'paused'}
          <button class="btn resume" on:click={() => dispatch('resume')}>
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

    <button class="scroll-btn" on:click={scrollToExecution}>
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
    border: 1px solid var(--border-color);
    border-radius: var(--radius-xl);
    background: var(--surface);
    padding: 28px;
    margin-bottom: 24px;
    box-shadow: var(--shadow-sm);
    position: relative;
    overflow: hidden;
    touch-action: manipulation;
  }

  .execution-panel::before {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(180deg,
      color-mix(in srgb, var(--primary-500) 2%, transparent),
      transparent 70%
    );
    pointer-events: none;
    opacity: 0.3;
  }

  .status-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 32px;
    position: relative;
    z-index: 1;
  }

  .status-info {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .status-badge {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    border-radius: var(--radius-lg);
    font-size: 14px;
    font-weight: var(--font-medium);
    letter-spacing: -0.01em;
    min-width: 90px;
    justify-content: center;
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
  }

  .status-badge.running {
    color: var(--success-600);
    background: color-mix(in srgb, var(--success-500) 8%, var(--surface));
    border: 1px solid color-mix(in srgb, var(--success-500) 15%, transparent);
  }

  :global(.dark) .status-badge.running {
    background: color-mix(in srgb, var(--success-500) 15%, var(--surface));
    color: var(--success-400);
    border-color: color-mix(in srgb, var(--success-400) 20%, transparent);
  }

  .status-badge.paused {
    color: var(--warning-600);
    background: color-mix(in srgb, var(--warning-500) 8%, var(--surface));
    border: 1px solid color-mix(in srgb, var(--warning-500) 15%, transparent);
  }

  :global(.dark) .status-badge.paused {
    background: color-mix(in srgb, var(--warning-500) 15%, var(--surface));
    color: var(--warning-400);
    border-color: color-mix(in srgb, var(--warning-400) 20%, transparent);
  }

  .pulse {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: currentColor;
    position: relative;
  }

  .pulse::before {
    content: '';
    position: absolute;
    inset: -4px;
    border-radius: 50%;
    background: currentColor;
    animation: pulse 1.5s cubic-bezier(0.4, 0, 0.6, 1) infinite;
    opacity: 0.5;
  }

  @keyframes pulse {
    0% { transform: scale(1); opacity: 0.5; }
    50% { transform: scale(1.5); opacity: 0; }
    100% { transform: scale(1); opacity: 0.5; }
  }

  .metrics {
    display: grid;
    grid-template-columns: repeat(4, minmax(120px, 1fr));
    gap: 32px;
    margin-left: 32px;
    position: relative;
  }

  .metrics::before {
    content: '';
    position: absolute;
    left: -16px;
    top: 50%;
    transform: translateY(-50%);
    width: 1px;
    height: 32px;
    background: var(--border-color);
    opacity: 0.5;
  }

  .metric {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .label {
    font-size: 13px;
    color: var(--text-tertiary);
  }

  .value {
    font-size: 15px;
    font-weight: var(--font-medium);
    color: var(--text-primary);
    font-feature-settings: "tnum";
    font-variant-numeric: tabular-nums;
  }

  .controls {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-left: 24px;
    position: relative;
  }

  .controls::before {
    content: '';
    position: absolute;
    left: -24px;
    top: 50%;
    transform: translateY(-50%);
    width: 1px;
    height: 24px;
    background: var(--border-color);
    opacity: 0.5;
  }

  .btn {
    appearance: none;
    -webkit-appearance: none;
    border: none;
    outline: none;
    height: 40px;
    min-width: 100px;
    padding: 0 20px;
    font-size: var(--text-base);
    font-weight: var(--font-medium);
    border-radius: var(--radius-lg);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    touch-action: manipulation;
    user-select: none;
    -webkit-tap-highlight-color: transparent;
    position: relative;
    overflow: hidden;
    will-change: transform;
    transform: translateZ(0);
    cursor: pointer;
  }

  .btn.pause {
    color: var(--warning-600);
    background: color-mix(in srgb, var(--warning-500) 8%, var(--surface));
    box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--warning-500) 15%, transparent),
                0 1px 2px color-mix(in srgb, var(--warning-500) 10%, transparent);
  }

  .btn.resume {
    color: var(--success-600);
    background: color-mix(in srgb, var(--success-500) 8%, var(--surface));
    box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--success-500) 15%, transparent),
                0 1px 2px color-mix(in srgb, var(--success-500) 10%, transparent);
  }

  .btn.stop {
    color: var(--error-600);
    background: color-mix(in srgb, var(--error-500) 8%, var(--surface));
    box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--error-500) 15%, transparent),
                0 1px 2px color-mix(in srgb, var(--error-500) 10%, transparent);
  }

  .btn:hover {
    transform: translateY(-1px);
  }

  .btn:active {
    transform: translateY(1px) scale(0.98);
  }

  .progress-bar {
    height: 8px;
    background: var(--surface-variant);
    border-radius: 4px;
    overflow: hidden;
    position: relative;
    box-shadow: inset 0 1px 2px color-mix(in srgb, black 5%, transparent);
    margin-top: 32px;
  }

  .progress {
    height: 100%;
    background: var(--primary-500);
    border-radius: 4px;
    transition: width 0.3s ease;
    position: relative;
    min-width: 8px;
  }

  .progress::before {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(90deg,
      var(--primary-400),
      var(--primary-500)
    );
    opacity: 1;
  }

  .progress::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(
      90deg,
      transparent 0%,
      rgba(255, 255, 255, 0.2) 50%,
      transparent 100%
    );
    animation: shine 2s linear infinite;
  }

  @keyframes shine {
    from { transform: translateX(-100%); }
    to { transform: translateX(100%); }
  }

  .current-stock {
    margin-top: 16px;
    font-size: 14px;
    color: var(--text-secondary);
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 16px;
    background: color-mix(in srgb, var(--surface-variant) 50%, transparent);
    border-radius: var(--radius-lg);
    border: 1px solid color-mix(in srgb, var(--border-color) 50%, transparent);
  }

  .stock {
    color: var(--text-primary);
    font-weight: var(--font-medium);
    font-feature-settings: "tnum";
    font-variant-numeric: tabular-nums;
  }

  .scroll-btn {
    position: absolute;
    bottom: 16px;
    right: 16px;
    border: none;
    background: transparent;
    cursor: pointer;
    padding: 8px 16px;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--primary-500);
    font-size: 13px;
    font-weight: var(--font-medium);
    background: color-mix(in srgb, var(--primary-500) 8%, var(--surface));
    transition: all 0.3s ease;
    touch-action: manipulation;
    -webkit-tap-highlight-color: transparent;
    z-index: 2;
  }

  .scroll-btn:hover {
    background: color-mix(in srgb, var(--primary-500) 12%, var(--surface));
    transform: translateY(-1px);
  }

  .scroll-btn:active {
    transform: translateY(0) scale(0.98);
  }

  .scroll-icon {
    display: flex;
    align-items: center;
    gap: 6px;
    transition: all 0.3s ease;
  }

  .scroll-icon svg {
    width: 16px;
    height: 16px;
    transition: transform 0.3s ease;
  }

  .scroll-btn:hover .scroll-icon svg {
    transform: translateY(2px);
  }
</style> 