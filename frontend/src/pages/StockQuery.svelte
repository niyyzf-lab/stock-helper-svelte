<script lang="ts">
  import AiAnalysisCard from '../components/stock-query/AIAnalysisCard.svelte';
  import StockHeader from '../components/stock-query/StockHeader.svelte';
  import StockChart from '../components/StockChart.svelte';
  import { onMount } from 'svelte';
  import { location, querystring } from 'svelte-spa-router'
  import { Splitpanes, Pane } from 'svelte-splitpanes'
  import type { IPaneSizingEvent } from 'svelte-splitpanes'

  // 状态变量
  let code = '';
  let name = '';
  let broker = '';
  let freq = 'dh';

  // 从查询参数解析股票信息
  $: {
    const params = new URLSearchParams($querystring)
    code = params.get('code') || ''
    name = params.get('name') || ''
    broker = params.get('exchange') || ''
  }

  // 处理频率变化
  function handleFreqChange(event: CustomEvent<{freq: string}>) {
    freq = event.detail.freq;
  }

  // 处理分享
  function handleShare() {
    const params = new URLSearchParams({
      code,
      name,
      exchange: broker
    }).toString()
    
    const shareUrl = `${window.location.origin}/#/stock?${params}`
    if (navigator.share) {
      navigator.share({
        title: `${name}(${code})`,
        text: `查看 ${name}(${code}) 的股票信息`,
        url: shareUrl
      }).catch(err => {
        console.warn('分享失败:', err);
      });
    } else {
      navigator.clipboard.writeText(shareUrl).then(() => {
        alert('链接已复制到剪贴板');
      }).catch(err => {
        console.error('复制失败:', err);
      });
    }
  }

  // 保存布局
  function handleSplitChange(event: CustomEvent<IPaneSizingEvent[]>) {
    try {
      const sizes = event.detail.map(pane => pane.size);
      localStorage.setItem('stock-query-layout', JSON.stringify(sizes));
    } catch (e) {
      console.warn('Failed to save layout:', e);
    }
  }

  // 获取初始布局
  let initialSizes = [70, 30];
  onMount(() => {
    try {
      const saved = localStorage.getItem('stock-query-layout');
      if (saved) {
        initialSizes = JSON.parse(saved);
      }
    } catch (e) {
      console.warn('Failed to load layout:', e);
    }
  });
</script>

<div class="stock-query-container">
  <Splitpanes on:resize={handleSplitChange}>
    <Pane minSize={60} maxSize={75}  size={initialSizes[0]}>
      <div class="stock-query-container-left">
        <div class="stock-chart-container">
          <div class="stock-chart-container-header card">
            <StockHeader 
              {code}
              {name}
              {freq}
              on:freqChange={handleFreqChange}
              on:share={handleShare}
            />
          </div>
          <div class="stock-chart-container-content card">
            <StockChart {code} {freq} />
          </div>
        </div>
      </div>
    </Pane>
    <Pane minSize={25} size={initialSizes[1]}>
      <div class="stock-query-container-right">
        <div class="ai-card card"> 
          <AiAnalysisCard {code} />
        </div>
        <div class="other-placeholder card"> </div>
      </div>
    </Pane>
  </Splitpanes>
</div>

<style>
  /* 引入 splitpanes 样式 */
  :global(.splitpanes) {
    height: 100%;
  }
  
  /* 定义飞入动画关键帧 */
  @keyframes flyIn {
    0% {
      opacity: 0;
      transform: translateX(30px);
    }
    70% {
      transform: translateX(-5px);
    }
    85% {
      transform: translateX(2px);
    }
    100% {
      opacity: 1;
      transform: translateX(0);
    }
  }

  /* 为左侧和右侧容器中的卡片设置不同的动画延迟 */
  .stock-query-container-left .card {
    animation: flyIn 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
  }

  .stock-query-container-right .card {
    animation: flyIn 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
    animation-delay: 0.1s;
  }
  
  /* 覆盖默认主题的边框样式 */
  :global(.default-theme.splitpanes--vertical > .splitpanes__splitter),
  :global(.default-theme .splitpanes--vertical > .splitpanes__splitter) {
    border: none !important;
    border-left: none !important;
    border-right: none !important;
    width: 4px !important;
    margin: 0 !important;
    background: var(--border-color) !important;
  }
  
  :global(.splitpanes__pane) {
    height: 100%;
    overflow: hidden;
  }
  
  :global(.splitpanes__splitter) {
    background: var(--border-color) !important;
    position: relative;
    width: 4px !important;
    border: none !important;
    margin: 0 !important;
    transition: var(--transition-colors);
  }
  
  :global(.splitpanes__splitter:hover) {
    background: var(--neutral-400) !important;
  }
  
  :global(.splitpanes__splitter::before) {
    content: '';
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 2px;
    height: 40px;
    background: var(--neutral-300);
    border-radius: 2px;
    transition: var(--transition-colors);
  }
  
  :global(.splitpanes__splitter:hover::before) {
    height: 60px;
    background: var(--neutral-500);
    box-shadow: var(--shadow-sm);
  }

  .card {
    background-color: var(--surface);
    border-radius: 12px;
    padding: 5px;
    min-height: 20px;
    box-shadow: var(--shadow-sm);
    border: 1px solid var(--border-color);
    width: 100%;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    transition: var(--transition-colors);
  }

  .card:hover {
    box-shadow: var(--shadow-md);
    border-color: var(--border-color);
  }

  .stock-query-container {
    height: 100%;
    overflow: hidden;
    position: relative;
    box-sizing: border-box;
    background: var(--background);
  }

  .stock-query-container-left, .stock-query-container-right {
    display: flex;
    flex-direction: column;
    gap: 16px;
    overflow-y: auto;
    padding: 20px;
    height: 100%;
    position: relative;
    overflow-x: hidden;
    background: var(--background);
  }
  .stock-query-container-left{
    padding-right: 10px;
  }
  .stock-query-container-right{
    padding-left: 10px;
  }
  .stock-query-container-left::-webkit-scrollbar,
  .stock-query-container-right::-webkit-scrollbar {
    width: 6px;
  }

  .stock-query-container-left::-webkit-scrollbar-thumb,
  .stock-query-container-right::-webkit-scrollbar-thumb {
    background-color: var(--neutral-300);
    border-radius: 3px;
  }

  .stock-query-container-left::-webkit-scrollbar-track,
  .stock-query-container-right::-webkit-scrollbar-track {
    background: transparent;
  }

  .stock-chart-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
    height: 100%;
    min-height: 0;
    flex: 1;
    transition: all 0.3s ease;
  }

  .stock-chart-container-content {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    background: #fff;
    border-radius: 12px;
    overflow: hidden;
  }

  .ai-card, .other-placeholder {
    min-height: 100px;
    height: auto;
    flex: 0 0 auto;
  }

  /* 暗色模式特殊处理 */
  :global(.dark) .card {
    background: var(--surface);
    border-color: var(--border-color);
  }

  :global(.dark) .card:hover {
    box-shadow: 0 3px 12px rgba(0, 0, 0, 0.2);
  }

  :global(.dark) .stock-query-container-left::-webkit-scrollbar-thumb,
  :global(.dark) .stock-query-container-right::-webkit-scrollbar-thumb {
    background-color: var(--neutral-600);
  }

  :global(.dark) :global(.splitpanes__splitter::before) {
    background: var(--neutral-600);
  }

  :global(.dark) :global(.splitpanes__splitter:hover::before) {
    background: var(--neutral-500);
    box-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
  }
</style>
