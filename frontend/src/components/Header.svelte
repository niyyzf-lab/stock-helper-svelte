<script lang="ts">
  import { Search, X, Command, HelpCircle, ChevronDown, ChevronUp, FileSearch } from 'lucide-svelte'
  import WindowControls from './WindowControls.svelte'
  import Modal from './Modal.svelte'
  import { push } from 'svelte-spa-router'
  import { pinyin } from 'pinyin-pro'
  import { fade, fly, scale } from 'svelte/transition'
  import { elasticOut, quintOut } from 'svelte/easing'
    import { onMount } from 'svelte';

  interface Stock {
    code: string    // 对应后端的 dm
    name: string    // 对应后端的 mc
    exchange: string // 对应后端的 jys
    fullPinyin: string // 完整拼音
    firstPinyin: string // 首字母
  }

  interface StockParams {
    code: string
    name: string
    exchange: string
  }

  interface StockResponse {
    dm: string;    // 代码
    mc: string;    // 名称
    jys: string;   // 交易所
  }

  let mounted = false
  let showSearchModal = false
  let searchInput = ''
  let loading = false
  let stocks: Stock[] = []

  onMount(() => {
    // 使用 RAF 确保在下一帧执行，避免闪烁
    requestAnimationFrame(() => {
      mounted = true
    })
  })

  // 获取拼音信息
  function getPinyinInfo(text: string): { full: string, first: string } {
    // 获取完整拼音，移除声调
    const full = pinyin(text, { toneType: 'none', type: 'array' }).join('')
    // 获取首字母
    const first = pinyin(text, { pattern: 'first', toneType: 'none', type: 'array' }).join('')
    return { full, first }
  }

  async function loadStocks() {
    try {
      loading = true
      const result = await (window as any).go.main.App.GetIndexList()
      stocks = (result || []).map((stock: StockResponse) => {
        const pinyinInfo = getPinyinInfo(stock.mc || '')
        return {
          code: stock.dm || '',
          name: stock.mc || '',
          exchange: stock.jys || '',
          fullPinyin: pinyinInfo.full,
          firstPinyin: pinyinInfo.first
        }
      }).filter((stock: Stock) => stock.code && stock.name)
    } catch (err) {
      console.error('加载股票列表失败:', err)
      stocks = []
    } finally {
      loading = false
    }
  }

  function openSearchModal() {
    showSearchModal = true
    searchInput = ''  // 清除搜索输入
    loadStocks()
  }

  // 处理股票选择
  function handleStockSelect(stock: Stock) {
    showSearchModal = false
    searchInput = ''
    const params = new URLSearchParams({
      code: stock.code,
      name: stock.name,
      exchange: stock.exchange
    })
    push(`/stock?${params.toString()}`)
  }

  function clearSearch() {
    searchInput = ''
  }

  // 模糊匹配函数
  function fuzzyMatch(text: string, query: string): boolean {
    let i = 0
    const queryLen = query.length
    const textLen = text.length
    
    for (let j = 0; j < textLen && i < queryLen; j++) {
      if (query[i].toLowerCase() === text[j].toLowerCase()) {
        i++
      }
    }
    
    return i === queryLen
  }

  $: filteredStocks = (() => {
    const input = searchInput.toLowerCase().trim()
    if (!input) return []
    
    return stocks.filter(stock => {
      // 股票代码匹配
      if (stock.code.toLowerCase().includes(input)) return true
      
      // 股票名称匹配
      if (stock.name.toLowerCase().includes(input)) return true
      
      // 完整拼音匹配
      if (stock.fullPinyin.toLowerCase().includes(input)) return true
      
      // 首字母匹配
      if (stock.firstPinyin.toLowerCase().includes(input)) return true
      
      // 模糊匹配 - 名称
      if (fuzzyMatch(stock.name.toLowerCase(), input)) return true
      
      // 模糊匹配 - 拼音
      if (fuzzyMatch(stock.fullPinyin.toLowerCase(), input)) return true
      if (fuzzyMatch(stock.firstPinyin.toLowerCase(), input)) return true
      
      return false
    }).slice(0, 100)
  })()

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter' && filteredStocks.length > 0) {
      handleStockSelect(filteredStocks[0])
    }
  }

  // 添加快捷键处理
  function handleGlobalKeydown(event: KeyboardEvent) {
    if ((event.metaKey || event.ctrlKey) && event.key === 'k') {
      event.preventDefault()
      openSearchModal()
    }
  }
</script>

<svelte:window on:keydown={handleGlobalKeydown}/>

{#if mounted}
<header style="--wails-draggable:drag" in:fly={{y: -30, duration: 800, delay: 200, easing: quintOut}}>
  <div class="header-content">
    <div class="left" in:fly={{x: -30, duration: 800, delay: 400, easing: quintOut}}>
      <div class="search-box" 
        on:click={openSearchModal} 
        style="--wails-draggable:no-drag"
        in:scale={{duration: 600, delay: 600, easing: elasticOut}}
      >
        <Search size={16} strokeWidth={1.5} />
        <input 
          type="text" 
          placeholder="搜索股票代码、名称..." 
          readonly
        />
        <div class="shortcut-hint" in:fade={{duration: 400, delay: 800}}>
          <Command size={14} strokeWidth={1.5} />
          <span>K</span>
        </div>
      </div>
    </div>
    <div class="right" in:fly={{x: 30, duration: 800, delay: 400, easing: quintOut}}>
      <WindowControls />
    </div>
  </div>
</header>
{/if}

<Modal 
  bind:show={showSearchModal}
  title="搜索股票"
>
  <div class="search-modal">
    <div class="search-input" in:fly={{y: -20, duration: 400, delay: 100, easing: quintOut}}>
      <div class="search-icon">
        <Search size={16} strokeWidth={1.5} />
      </div>
      <input
        type="text"
        placeholder="输入股票代码、名称、拼音搜索..."
        bind:value={searchInput}
        on:keydown={handleKeydown}
        autofocus
      />
      {#if searchInput}
        <button class="clear-btn" on:click={clearSearch} in:scale={{duration: 400, easing: elasticOut}}>
          <X size={14} />
        </button>
      {:else}
        <div class="shortcut-hint" in:fade={{duration: 300, delay: 200}}>
          <span>ESC</span>
          <span class="hint-text">关闭</span>
        </div>
      {/if}
    </div>

    <div class="search-content">
      {#if loading}
        <div class="state-message" in:fade={{duration: 400}}>
          <div class="spinner" />
          <span>正在加载股票列表...</span>
        </div>
      {:else if stocks.length === 0}
        <div class="state-message" in:fade={{duration: 400}}>
          <span class="state-title">暂无股票数据</span>
          <span class="state-desc">请检查网络连接</span>
        </div>
      {:else if !searchInput}
        <div class="state-message" in:fly|local={{y: 30, duration: 600, delay: 200, easing: quintOut}}>
          <div class="empty-content" in:fly|local={{y: 20, duration: 600, delay: 400, easing: quintOut}}>
            <div class="empty-icon" in:fly|local={{y: -20, duration: 800, delay: 600, easing: elasticOut}}>
              <FileSearch size={28} strokeWidth={1.5} />
            </div>
            <div class="empty-text" in:fly|local={{y: 20, duration: 600, delay: 800}}>
              <h3>搜索股票</h3>
              <p>输入股票代码、名称或拼音进行搜索</p>
            </div>
          </div>
          <div class="search-tips" in:fly|local={{y: 30, duration: 600, delay: 1000, easing: quintOut}}>
            <div class="tips-row">
              <div class="tips-col">
                <div class="tips-label">代码：</div>
                <div class="tips-value">600519</div>
              </div>
              <div class="tips-col">
                <div class="tips-label">名称：</div>
                <div class="tips-value">贵州茅台</div>
              </div>
              <div class="tips-col">
                <div class="tips-label">拼音：</div>
                <div class="tips-value">guizhou</div>
              </div>
              <div class="tips-col">
                <div class="tips-label">首字母：</div>
                <div class="tips-value">gzmt</div>
              </div>
              <div class="tips-col">
                <div class="tips-label">模糊：</div>
                <div class="tips-value">贵茅</div>
              </div>
            </div>
          </div>
        </div>
      {:else if filteredStocks.length === 0}
        <div class="state-message" in:fly={{y: 30, duration: 600, delay: 200, easing: quintOut}}>
          <div class="empty-content" in:fly={{y: 20, duration: 600, delay: 400, easing: quintOut}}>
            <div class="empty-icon" in:fly={{y: -20, duration: 800, delay: 600, easing: elasticOut}}>
              <FileSearch size={28} strokeWidth={1.5} />
            </div>
            <div class="empty-text" in:fly={{y: 20, duration: 600, delay: 800}}>
              <h3>未找到匹配的股票</h3>
              <p>试试以下搜索方式</p>
            </div>
          </div>
          <div class="search-tips" in:fly={{y: 30, duration: 600, delay: 1000, easing: quintOut}}>
            <div class="tips-row">
              <div class="tips-col">
                <div class="tips-label">代码：</div>
                <div class="tips-value">600519</div>
              </div>
              <div class="tips-col">
                <div class="tips-label">名称：</div>
                <div class="tips-value">贵州茅台</div>
              </div>
              <div class="tips-col">
                <div class="tips-label">拼音：</div>
                <div class="tips-value">guizhou</div>
              </div>
              <div class="tips-col">
                <div class="tips-label">首字母：</div>
                <div class="tips-value">gzmt</div>
              </div>
              <div class="tips-col">
                <div class="tips-label">模糊：</div>
                <div class="tips-value">贵茅</div>
              </div>
            </div>
          </div>
        </div>
      {:else}
        <div in:fly={{y: 20, duration: 400, easing: quintOut}}>
          <div class="results-header">
            找到 {filteredStocks.length} 个结果
          </div>
          <div class="stock-list">
            {#each filteredStocks as stock, i}
              <button 
                class="stock-item"
                on:click={() => handleStockSelect(stock)}
                in:fly={{y: 20, duration: 400, delay: i * 50, easing: quintOut}}
              >
                <div class="stock-info">
                  <div class="stock-main">
                    <span class="stock-code">{stock.code}</span>
                    <span class="stock-name">{stock.name}</span>
                  </div>
                  <span class="stock-exchange">{stock.exchange}</span>
                </div>
              </button>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  </div>
</Modal>

<style>
  header {
    height: 60px;
    border-bottom: 1px solid var(--border-color);
    padding-left: 14px; /* 增加左边距 */
    transform-origin: top;
    opacity: 0;
    animation: fadeIn 0.8s cubic-bezier(0.4, 0, 0.2, 1) forwards;
    background: var(--surface);
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .header-content {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
  }

  .left {
    transform-origin: left;
  }

  .right {
    transform-origin: right;
  }

  .search-box {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 0 12px;
    height: 36px;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    color: var(--text-secondary);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    cursor: pointer;
    justify-content: space-between;
    transform-origin: center;
  }

  .search-box:hover {
    background: var(--hover-bg);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px var(--shadow-color);
  }

  input {
    width: 240px;
    border: none;
    background: none;
    font-size: var(--text-sm);
    color: var(--text-primary);
    outline: none;
    font-family: inherit;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  input::placeholder {
    color: var(--text-secondary);
  }

  .right {
    display: flex;
    align-items: center;
  }

  /* 搜索模态框样式 */
  .search-modal {
    width: 100%;
    min-width: 640px;
    max-width: 800px;
    display: flex;
    flex-direction: column;
    padding: 24px;
    padding-top: 4px;
    overflow: hidden;
    animation: modalShow 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes modalShow {
    from {
      opacity: 0;
      transform: scale(0.95) translateY(-10px);
    }
    to {
      opacity: 1;
      transform: scale(1) translateY(0);
    }
  }

  .search-input {
    position: relative;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 0 16px;
    height: 48px;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    color: var(--text-secondary);
    margin-bottom: 20px;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    transform-origin: top;
  }

  .search-input:focus-within {
    border-color: var(--primary-500);
    background: var(--surface);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
  }

  .search-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-tertiary);
    transition: all 0.2s ease;
  }

  .search-input:focus-within .search-icon {
    color: var(--primary-500);
  }

  .search-input input {
    flex: 1;
    height: 100%;
    width: auto;
    font-size: 15px;
    background: transparent;
    border: none;
    outline: none;
    cursor: text;
    caret-color: var(--primary-500);
    padding: 0;
  }

  .clear-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 6px;
    color: var(--text-secondary);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    background: none;
    border: none;
    cursor: pointer;
    opacity: 0.8;
    border-radius: 4px;
    transform-origin: center;
  }

  .clear-btn:hover {
    opacity: 1;
    color: var(--text-primary);
    background: var(--neutral-100);
    transform: scale(1.2);
  }

  .clear-btn:active {
    transform: scale(0.9);
  }

  .search-content {
    flex: 1;
    overflow-y: auto;
    margin: 0 -24px;
    padding: 0 24px;
    max-height: 500px;
    position: relative;
  }

  .state-message {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-direction: column;
    padding: 64px 24px 24px;
    color: var(--text-secondary);
    font-size: 14px;
    background: var(--surface-variant);
    border-radius: 12px;
    margin: 0;
    min-height: 400px;
    overflow: hidden;
    transform-origin: center;
  }

  .empty-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 24px;
    text-align: center;
    margin-bottom: auto;
    will-change: transform;
    transform-origin: center;
  }

  .empty-icon {
    width: 64px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--primary-500);
    background: var(--surface);
    border: 1px solid var(--border-color);
    border-radius: 16px;
    box-shadow: 0 4px 12px var(--shadow-color);
    animation: float 8s ease-in-out infinite;
    will-change: transform;
    transform-origin: center;
  }

  @keyframes float {
    0%, 100% {
      transform: translateY(0) rotate(0);
    }
    25% {
      transform: translateY(-8px) rotate(-3deg);
    }
    75% {
      transform: translateY(6px) rotate(3deg);
    }
  }

  .empty-text {
    display: flex;
    flex-direction: column;
    gap: 8px;
    transform-origin: center;
  }

  .empty-text h3 {
    font-size: 16px;
    font-weight: 500;
    color: var(--text-primary);
    margin: 0;
  }

  .empty-text p {
    font-size: 14px;
    color: var(--text-secondary);
    margin: 0;
  }

  .search-tips {
    width: 100%;
    padding: 16px;
    background: var(--surface);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    margin-top: auto;
    backdrop-filter: blur(8px);
    will-change: transform, opacity;
    transform-origin: bottom;
  }

  .tips-row {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    gap: 16px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .tips-col {
    display: flex;
    align-items: center;
    gap: 4px;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .tips-col:hover {
    transform: translateY(-2px);
  }

  .tips-label {
    font-size: 13px;
    color: var(--text-tertiary);
    white-space: nowrap;
  }

  .tips-value {
    font-size: 13px;
    color: var(--primary-500);
    background: var(--surface-variant);
    padding: 2px 8px;
    border-radius: 4px;
    font-family: var(--font-mono);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .tips-col:hover .tips-value {
    background: var(--hover-bg);
    color: var(--primary-400);
    transform: scale(1.1);
  }

  .stock-list {
    display: flex;
    flex-direction: column;
    gap: 1px;
    padding: 4px 0;
    position: relative;
    margin-bottom: 40px;
  }

  .stock-item {
    width: 100%;
    padding: 12px 16px;
    background: none;
    border: none;
    cursor: pointer;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    text-align: left;
    position: relative;
    border-radius: 8px;
    margin: 0;
    display: block;
    transform-origin: center;
  }

  .stock-item:hover {
    background: var(--hover-bg);
    transform: translateX(6px) scale(1.01);
  }

  .stock-item:active {
    transform: scale(0.98);
  }

  .stock-item:not(:last-child) {
    border-bottom: 1px solid var(--border-color);
  }

  .stock-info {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
  }

  .stock-main {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .stock-code {
    font-size: 15px;
    color: var(--text-primary);
    font-family: var(--font-mono);
    font-weight: 500;
    min-width: 80px;
    transition: all 0.2s ease;
  }

  .stock-name {
    font-size: 15px;
    color: var(--text-secondary);
    transition: all 0.2s ease;
  }

  .stock-exchange {
    font-size: 13px;
    color: var(--text-tertiary);
    padding: 2px 8px;
    background: var(--surface-variant);
    border-radius: 4px;
    font-weight: 500;
    letter-spacing: 0.02em;
    text-transform: uppercase;
    transition: all 0.2s ease;
  }

  .stock-item:hover .stock-code {
    color: var(--primary-600);
  }

  .stock-item:hover .stock-name {
    color: var(--primary-500);
  }

  .stock-item:hover .stock-exchange {
    color: var(--primary-500);
    background: var(--hover-bg);
  }

  .shortcut-hint {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 8px;
    background: var(--surface-variant);
    border-radius: 4px;
    color: var(--text-tertiary);
    font-size: 12px;
    font-family: var(--font-mono);
    user-select: none;
    transition: all 0.2s ease;
  }

  .hint-text {
    font-size: 13px;
    color: var(--text-tertiary);
    margin-left: 2px;
  }

  .results-header {
    padding: 12px 0;
    color: var(--text-secondary);
    font-size: 13px;
    border-bottom: 1px solid var(--border-color);
    margin-bottom: 4px;
    background: var(--surface);
    position: sticky;
    top: 0;
    z-index: 1;
    animation: slideDown 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .spinner {
    width: 20px;
    height: 20px;
    border: 2px solid var(--border-color);
    border-top-color: var(--primary-500);
    border-radius: 50%;
    animation: spin 1.2s cubic-bezier(0.4, 0, 0.2, 1) infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* 优化滚动条样式 */
  .search-content::-webkit-scrollbar {
    width: 6px;
  }

  .search-content::-webkit-scrollbar-track {
    background: transparent;
    margin: 4px 0;
  }

  .search-content::-webkit-scrollbar-thumb {
    background: var(--border-color);
    border-radius: 6px;
    border: 2px solid transparent;
    background-clip: padding-box;
  }

  .search-content::-webkit-scrollbar-thumb:hover {
    background: var(--text-secondary);
    border: 2px solid transparent;
    background-clip: padding-box;
  }
</style> 