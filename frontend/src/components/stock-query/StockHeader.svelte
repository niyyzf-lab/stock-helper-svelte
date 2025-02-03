<!-- 股票头部组件 -->
<script lang="ts">
  import { Share2, ChevronRight, ChevronLeft } from 'lucide-svelte';
  import { createEventDispatcher } from 'svelte';
  import { fade } from 'svelte/transition';

  // 导出属性
  export let code: string = '';
  export let name: string = '';
  export let freq: string = 'dh';

  // 创建事件分发器
  const dispatch = createEventDispatcher();

  // 时间周期选项
  const timeFrames = [
    { value: '5m', label: '5M', shortLabel: '5' },
    { value: '15m', label: '15M', shortLabel: '15' },
    { value: '30m', label: '30M', shortLabel: '30' },
    { value: '60m', label: '60M', shortLabel: '60' },
    { value: 'dh', label: '日线', shortLabel: '日' },
    { value: 'wh', label: '周线', shortLabel: '周' },
    { value: 'mh', label: '月线', shortLabel: '月' },
    { value: 'yh', label: '年线', shortLabel: '年' }
  ];


  // 处理时间周期变化
  function handleFreqChange(newFreq: string) {
    freq = newFreq;
    dispatch('freqChange', { freq: newFreq });
  }

  // 处理分享点击
  function handleShareClick() {
    dispatch('share');
  }

  let timeFramesRef: HTMLDivElement;

  function scrollTimeFrames(direction: 'left' | 'right') {
    if (!timeFramesRef) return;
    const scrollAmount = direction === 'left' ? -100 : 100;
    timeFramesRef.scrollBy({ left: scrollAmount, behavior: 'smooth' });
  }
</script>

<div class="stock-header">
  <!-- 左侧股票信息 -->
  <div class="stock-info">
    <div class="stock-name">
      <span class="stock-code">{code}</span>
      <span class="stock-badge">A股</span>
      <ChevronRight size={16} class="chevron-icon" />
      <span class="stock-title">{name}</span>
    </div>
  </div>

  <!-- 右侧操作区 -->
  <div class="operations">
    <!-- 时间周期选择器 -->
    <div class="time-selector">
      <button class="scroll-btn" on:click={() => scrollTimeFrames('left')}>
        <ChevronLeft size={16} />
      </button>
      
      <div class="time-frames" bind:this={timeFramesRef}>
        {#each timeFrames as frame}
          <button
            class="time-frame-btn"
            class:active={freq === frame.value}
            on:click={() => handleFreqChange(frame.value)}
          >
            {frame.label}
          </button>
        {/each}
      </div>

      <button class="scroll-btn" on:click={() => scrollTimeFrames('right')}>
        <ChevronRight size={16} />
      </button>
    </div>

    <!-- 分享按钮 -->
    <button class="share-btn" on:click={handleShareClick}>
      <Share2 size={16} />
    </button>
  </div>
</div>

<style>
  .stock-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    position: relative;
    z-index: 10;
    background: var(--surface);
    transition: all 0.2s ease;
  }

  .stock-info {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 2px 0;
  }

  .stock-name {
    display: flex;
    align-items: center;
    gap: 6px;
    color: var(--text-primary);
    font-size: 14px;
    letter-spacing: 0.2px;
  }

  .stock-code {
    font-weight: 600;
    font-size: 15px;
    min-width: 60px;
    display: inline-block;
  }

  .stock-badge {
    font-size: 11px;
    padding: 1px 6px;
    border-radius: 3px;
    background: var(--hover-bg);
    color: var(--text-secondary);
    font-weight: 500;
    letter-spacing: 0.3px;
  }

  .chevron-icon {
    color: var(--text-tertiary);
    margin: 0 -4px;
  }

  .stock-title {
    color: var(--text-secondary);
    font-weight: 500;
  }

  .operations {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .time-selector {
    display: flex;
    align-items: center;
    gap: 2px;
    position: relative;
    z-index: 1000;
    background: var(--surface);
    border-radius: 6px;
    padding: 2px;
    border: 1px solid var(--border-color);
  }

  .time-frames {
    display: flex;
    gap: 2px;
    overflow-x: auto;
    scrollbar-width: none;
    -ms-overflow-style: none;
    max-width: 300px;
    padding: 0 2px;
  }

  .time-frames::-webkit-scrollbar {
    display: none;
  }

  .scroll-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    border-radius: 4px;
    color: var(--text-secondary);
    background: transparent;
    border: none;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    height: 24px;
    width: 24px;
    line-height: 1;
  }

  .scroll-btn:hover {
    background: var(--hover-bg);
    color: var(--text-primary);
  }

  .scroll-btn:active {
    transform: scale(0.95);
  }

  .time-frame-btn {
    padding: 0 10px;
    height: 24px;
    min-width: 42px;
    border-radius: 4px;
    background: transparent;
    border: none;
    color: var(--text-secondary);
    cursor: pointer;
    white-space: nowrap;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    font-weight: 500;
    font-size: 12px;
    position: relative;
    overflow: hidden;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    line-height: 1;
  }

  .time-frame-btn:hover {
    color: var(--text-primary);
    background: var(--hover-bg);
  }

  .time-frame-btn.active {
    background: var(--primary-600, #2563eb);
    color: white;
    font-weight: 500;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .time-frame-btn.active:hover {
    background: var(--primary-700, #1d4ed8);
  }

  .share-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    height: 28px;
    width: 28px;
    border-radius: 6px;
    color: var(--text-secondary);
    background: transparent;
    border: 1px solid var(--border-color);
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    line-height: 1;
  }

  .share-btn:hover {
    background: var(--hover-bg);
    color: var(--text-primary);
    border-color: var(--text-secondary);
  }

  .share-btn:active {
    transform: scale(0.95);
  }

  @media (max-width: 640px) {
    .stock-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
      padding: 12px;
    }

    .operations {
      width: 100%;
    }

    .time-selector {
      flex: 1;
    }

    .time-frames {
      max-width: none;
      width: 100%;
    }

    .stock-name {
      font-size: 13px;
    }

    .stock-code {
      font-size: 14px;
    }
  }
</style>