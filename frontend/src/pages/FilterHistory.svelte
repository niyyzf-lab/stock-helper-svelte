<script lang="ts">
  import { History, RefreshCw, ChevronRight, Trash2, ChevronLeft, ChevronUp } from 'lucide-svelte'
  import { fade, fly } from 'svelte/transition'
  import EmptyState from '../components/EmptyState.svelte'
  import { toastStore } from '../stores/toast'
  import ExecutionRecordDetail from '../components/ExecutionRecordDetail.svelte'
  import Modal from '../components/Modal.svelte'
  import { onMount, onDestroy } from 'svelte'

  // 历史记录接口
  interface ExecutionRecord {
    fileName: string        // 文件名
    strategyId: number     // 策略ID
    strategyName: string   // 策略名称
    executionTime: string  // 执行时间
    signalCount: number    // 信号数量
    processedCount: number // 处理数量
    totalStocks: number    // 总股票数
  }

  let records: ExecutionRecord[] = []
  let loading = false
  let selectedRecord: ExecutionRecord | null = null
  let showConfirmModal = false
  let recordToDelete: ExecutionRecord | null = null
  let showingDetail = false
  let currentIndex = -1

  // 加载历史记录
  async function loadRecords() {
    loading = true
    try {
      // 添加最短加载时间
      const startTime = Date.now()
      const result = await (window as any).go.main.App.GetExecutionRecords()
      // 确保至少显示500ms的加载状态
      const elapsed = Date.now() - startTime
      if (elapsed < 500) {
        await new Promise(resolve => setTimeout(resolve, 500 - elapsed))
      }
      
      console.log('加载的记录列表:', result)
      records = result || []
      console.log('记录数量:', records.length)
    } catch (err) {
      toastStore.error('加载历史记录失败')
      console.error('加载历史记录失败:', err)
      records = []
    } finally {
      loading = false
    }
  }

  // 格式化日期
  function formatDate(dateStr: string | undefined): string {
    if (!dateStr) return '未知时间'
    try {
      const date = new Date(dateStr)
      if (isNaN(date.getTime())) {
        return '无效时间'
      }
      
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      })
    } catch (err) {
      return '无效时间'
    }
  }

  // 格式化数字
  function formatNumber(num: number | undefined): string {
    if (num === undefined || isNaN(num)) {
      return '0'
    }
    return num.toString()
  }

  // 格式化进度
  function formatProgress(processed: number | undefined, total: number | undefined): string {
    if (processed === undefined || total === undefined || 
        typeof processed !== 'number' || typeof total !== 'number') {
      return '0/0'
    }
    return `${processed}/${total}`
  }

  // 刷新列表
  function refreshList() {
    loadRecords()
  }

  // 查看记录详情
  function viewRecord(record: ExecutionRecord) {
    console.log('查看记录:', record)
    console.log('当前记录列表:', records)
    selectedRecord = record
    showingDetail = true
    currentIndex = getCurrentIndex()
    console.log('当前记录索引:', currentIndex)
  }

  // 返回列表
  function backToList() {
    showingDetail = false
    selectedRecord = null
    currentIndex = -1
  }

  // 导航到上一条记录
  function navigatePrev() {
    if (!selectedRecord || !records.length) return
    const index = getCurrentIndex()
    console.log('导航到上一条:', { currentIndex: index, canNavigate: index > 0 })
    if (index > 0) {
      selectedRecord = records[index - 1]
      currentIndex = index - 1
      console.log('更新后的索引:', currentIndex)
    }
  }

  // 导航到下一条记录
  function navigateNext() {
    if (!selectedRecord || !records.length) return
    const index = getCurrentIndex()
    console.log('导航到下一条:', { currentIndex: index, canNavigate: index < records.length - 1 })
    if (index < records.length - 1) {
      selectedRecord = records[index + 1]
      currentIndex = index + 1
      console.log('更新后的索引:', currentIndex)
    }
  }

  // 获取当前记录索引
  function getCurrentIndex(): number {
    if (!selectedRecord || !records.length) return -1
    const index = records.findIndex(r => r.fileName === selectedRecord!.fileName)
    console.log('计算导航索引:', {
      selectedFileName: selectedRecord?.fileName,
      allFileNames: records.map(r => r.fileName),
      calculatedIndex: index
    })
    return index
  }

  // 键盘导航处理
  function handleKeydown(event: KeyboardEvent) {
    if (!showingDetail) return
    
    if (event.key === 'ArrowLeft') {
      navigatePrev()
    } else if (event.key === 'ArrowRight') {
      navigateNext()
    } else if (event.key === 'Escape') {
      backToList()
    }
  }

  async function deleteRecord(record: ExecutionRecord) {
    try {
      await (window as any).go.main.App.DeleteExecutionRecord(record.fileName)
      toastStore.success('删除成功')
      showConfirmModal = false
      recordToDelete = null
      // 先重置状态，再重新加载记录
      records = []
      await loadRecords()
    } catch (err) {
      toastStore.error('删除失败')
      console.error('删除记录失败:', err)
      showConfirmModal = false
      recordToDelete = null
    }
  }

  function confirmDelete(record: ExecutionRecord, e: Event) {
    e.stopPropagation()
    recordToDelete = record
    showConfirmModal = true
  }

  // 组件加载时获取历史记录
  loadRecords()

  // 添加键盘事件监听
  onMount(() => {
    window.addEventListener('keydown', handleKeydown)
  })

  onDestroy(() => {
    window.removeEventListener('keydown', handleKeydown)
  })
</script>

<div class="page-container">
  <div class="main-container">
    <header class="header" in:fade={{duration: 300}}>
      <div class="header-content">
        <div class="header-left" in:fly={{x: -20, duration: 400, delay: 100}}>
          <div class="title-group">
            {#if showingDetail}
              <button class="back-btn" on:click={backToList} in:fly={{x: -20, duration: 400}}>
                <ChevronLeft size={16} strokeWidth={1.5} />
                <span>返回列表</span>
              </button>
              <div class="divider" in:fade={{duration: 300, delay: 150}}></div>
              <span class="strategy-name" in:fly={{x: -10, duration: 400, delay: 200}}>{selectedRecord?.strategyName || ''}</span>
              <div class="record-index" in:fly={{x: -10, duration: 400, delay: 250}}>
                {currentIndex + 1}/{records.length}
              </div>
            {:else}
              <h1 in:fly={{x: -20, duration: 400}}>筛选历史</h1>
              <div class="divider" in:fade={{duration: 300, delay: 150}}></div>
              <span class="version" in:fly={{x: -10, duration: 400, delay: 200}}>Beta</span>
            {/if}
          </div>
        </div>
        <div class="header-right" in:fly={{x: 20, duration: 400, delay: 100}}>
          {#if !showingDetail}
            <button 
              class="btn solid refresh-btn" 
              class:loading 
              on:click={refreshList} 
              disabled={loading} 
              in:fly={{x: 20, duration: 300}}
            >
              <div class="icon-wrapper" class:spin={loading}>
                <RefreshCw size={16} />
              </div>
              <span class="btn-text" style="min-width: 4em;">
                {loading ? '刷新中' : '刷新列表'}
              </span>
            </button>
          {:else}
            <div class="navigation-controls" in:fly={{x: 20, duration: 400}}>
              <button 
                class="btn outline" 
                on:click={navigatePrev}
                disabled={currentIndex <= 0}
                in:fly={{x: 20, duration: 400, delay: 150}}
              >
                <ChevronLeft size={16} strokeWidth={1.5} />
                上一个
              </button>
              <button 
                class="btn outline" 
                on:click={navigateNext}
                disabled={currentIndex >= records.length - 1}
                in:fly={{x: 20, duration: 400, delay: 200}}
              >
                下一个
                <ChevronRight size={16} strokeWidth={1.5} />
              </button>
            </div>
          {/if}
        </div>
      </div>
    </header>

    <main class="main">
      <div class="content-wrapper">
        {#if loading}
          <div class="loading-state" in:fade={{duration: 200}}>
            <RefreshCw size={24} class="spin" />
            <span>加载中...</span>
          </div>
        {:else if showingDetail && selectedRecord}
          <div class="detail-view" 
            in:fade={{duration: 200}}
            out:fade={{duration: 150}}
          >
            {#key selectedRecord.fileName}
              <ExecutionRecordDetail 
                fileName={selectedRecord.fileName}
                onClose={backToList}
                isModal={false}
              />
            {/key}
          </div>
        {:else if records.length === 0}
          <section class="empty-state" 
            in:fade={{duration: 200}}
            out:fade={{duration: 150}}
          >
            <div class="empty-content">
              <div class="empty-icon">
                <History size={24} strokeWidth={1.5} />
              </div>
              <div class="empty-text">
                <h3>暂无筛选记录</h3>
                <p>执行策略筛选后的结果会保存在这里</p>
              </div>
            </div>
          </section>
        {:else}
          <div class="records-grid">
            {#each records as record, i}
              <div 
                class="record-card" 
                role="button"
                tabindex="0"
                on:click={() => viewRecord(record)}
                on:keydown={(e) => e.key === 'Enter' && viewRecord(record)}
                in:fade={{
                  duration: 200,
                  delay: i * 50
                }}
              >
                <div class="card-header" in:fade={{duration: 200, delay: i * 50 + 50}}>
                  <div class="strategy-info">
                    <div class="strategy-icon" in:fade={{duration: 200, delay: i * 50 + 100}}>
                      <History size={16} strokeWidth={1.5} />
                    </div>
                    <span class="strategy-name" in:fade={{duration: 200, delay: i * 50 + 150}}>{record.strategyName || '未知策略'}</span>
                  </div>
                  <ChevronRight size={16} class="arrow-icon"  />
                </div>
                <div class="card-content" in:fade={{duration: 200, delay: i * 50 + 250}}>
                  <div class="stats-row">
                    <div class="stat" in:fade={{duration: 200, delay: i * 50 + 300}}>
                      <span class="label">信号数量</span>
                      <span class="value">{formatNumber(record.signalCount)}</span>
                    </div>
                    <div class="stat" in:fade={{duration: 200, delay: i * 50 + 350}}>
                      <span class="label">处理进度</span>
                      <span class="value">
                        {#if typeof record.processedCount === 'number' && typeof record.totalStocks === 'number'}
                          {record.processedCount}/{record.totalStocks}
                        {:else}
                          0/0
                        {/if}
                      </span>
                    </div>
                  </div>
                  <div class="time-info" in:fade={{duration: 200, delay: i * 50 + 400}}>
                    执行时间：{formatDate(record.executionTime)}
                  </div>
                  <button 
                    class="delete-btn"
                    on:click={(e) => confirmDelete(record, e)}
                    on:keydown={(e) => e.key === 'Enter' && confirmDelete(record, e)}
                    in:fade={{duration: 200, delay: i * 50 + 450}}
                  >
                    <Trash2 size={16} strokeWidth={1.5} />
                  </button>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </main>
  </div>
</div>

{#if showConfirmModal && recordToDelete}
  <Modal
    title="确认删除"
    show={showConfirmModal}
    on:close={() => {
      showConfirmModal = false
      recordToDelete = null
    }}
  >
    <div class="confirm-content">
      <p>确定要删除策略 "{recordToDelete.strategyName}" 的执行记录吗？</p>
      <p class="text-sm text-gray-500">此操作无法撤销</p>
      <div class="confirm-actions">
        <button 
          class="btn outline"
          on:click={() => {
            showConfirmModal = false
            recordToDelete = null
          }}
        >
          取消
        </button>
        <button 
          class="btn solid danger"
          on:click={() => recordToDelete && deleteRecord(recordToDelete)}
        >
          删除
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
    padding-top: 0;
    overflow: hidden;
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

  .back-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    color: #374151;
    font-size: 16px;
    font-weight: 500;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .back-btn:hover {
    color: #1f2937;
    transform: translateX(-2px);
  }

  .back-btn:active {
    transform: translateX(-1px) scale(0.98);
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
  }

  .header-right {
    display: flex;
    align-items: center;
    animation: slideInRight 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .navigation-controls {
    display: flex;
    gap: 8px;
  }

  .btn {
    height: 34px;
    padding: 0 16px;
    font-size: 14px;
    font-weight: 500;
    border-radius: 6px;
    display: flex;
    align-items: center;
    gap: 6px;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .btn:active:not(:disabled) {
    transform: scale(0.95);
  }

  .btn.solid {
    color: white;
    background: #2563eb;
    border: none;
  }

  .btn.outline {
    color: #374151;
    background: white;
    border: 1px solid #e5e7eb;
  }

  .btn.solid:hover:not(:disabled) {
    background: #1d4ed8;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.2);
  }

  .btn.outline:hover:not(:disabled) {
    background: #f9fafb;
    border-color: #d1d5db;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  }

  .btn:disabled {
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
    overflow: auto;
    min-height: 0; /* 重要：允许内容区域收缩 */
  }

  .content-wrapper {
    width: 100%;
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0; /* 重要：允许内容区域收缩 */
  }

  /* 空状态 */
  .empty-state {
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

  .empty-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 32px;
    text-align: center;
    width: 100%;
    max-width: 480px;
    margin: 0 auto;
    transform-origin: center;
  }

  .empty-icon {
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
  }

  .empty-text h3 {
    font-size: 16px;
    font-weight: 600;
    color: #111827;
    margin: 0 0 8px;
  }

  .empty-text p {
    font-size: 14px;
    color: #6b7280;
    margin: 0;
  }

  /* 加载状态 */
  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    height: 400px;
    color: #6b7280;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  /* 记录网格 */
  .records-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;
    padding: 4px;
    align-items: start;
    height: auto;
  }

  /* 记录卡片 */
  .record-card {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    padding: 16px;
    cursor: pointer;
    position: relative;
    height: auto;
    min-height: 140px;
    display: flex;
    flex-direction: column;
    width: 100%;
    transform-origin: center center;
    transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
    animation: scaleIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) backwards;
    animation-delay: calc(var(--index, 0) * 50ms + 200ms);
  }

  .record-card:hover {
    transform: translateY(-2px) scale(1.02);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
    border-color: #d1d5db;
  }

  .record-card:active {
    transform: translateY(0) scale(0.98);
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 12px;
  }

  .strategy-info {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    min-width: 0;
  }

  .strategy-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: #f3f4f6;
    border-radius: 6px;
    color: #2563eb;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    animation: slideUp 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) backwards;
    animation-delay: calc(var(--index, 0) * 50ms + 250ms);
  }

  .strategy-name {
    font-weight: 500;
    color: #111827;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    animation: slideUp 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) backwards;
    animation-delay: calc(var(--index, 0) * 50ms + 300ms);
  }

  .card-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    gap: 12px;
  }

  .stats-row {
    display: flex;
    gap: 24px;
    margin-bottom: auto;
    animation: slideUp 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) backwards;
    animation-delay: calc(var(--index, 0) * 50ms + 350ms);
  }

  .stat {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .stat .label {
    font-size: 12px;
    color: #6b7280;
  }

  .stat .value {
    font-size: 14px;
    font-weight: 500;
    color: #111827;
  }

  .time-info {
    font-size: 12px;
    color: #6b7280;
    margin-top: auto;
    animation: slideUp 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) backwards;
    animation-delay: calc(var(--index, 0) * 50ms + 400ms);
  }

  .delete-btn {
    position: absolute;
    top: 16px;
    right: 16px;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    color: #ef4444;
    opacity: 0;
    transform: scale(0.9);
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .record-card:hover .delete-btn {
    opacity: 1;
    transform: scale(1);
  }

  .delete-btn:hover {
    background: #fef2f2;
    border-color: #fca5a5;
    transform: scale(1.1);
  }

  .delete-btn:active {
    transform: scale(0.95);
  }

  .record-card:hover .strategy-icon {
    transform: scale(1.1);
    background: #e0e7ff;
  }

  /* 详情视图 */
  .detail-view {
    flex: 1;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    overflow: hidden;
    height: 100%;
  }

  .record-index {
    margin-left: 12px;
    font-size: 13px;
    color: #6b7280;
    background: #f3f4f6;
    padding: 2px 8px;
    border-radius: 4px;
  }

  .confirm-content {
    padding: 16px 0;
  }

  .confirm-content p {
    margin: 0;
  }

  .confirm-content p + p {
    margin-top: 8px;
  }

  .text-sm {
    font-size: 14px;
  }

  .text-gray-500 {
    color: #6b7280;
  }

  .confirm-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
  }

  .btn.danger {
    background: #ef4444;
  }

  .btn.danger:hover:not(:disabled) {
    background: #dc2626;
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
  }

  :global(.strategy-modal) {
    min-width: 600px !important;
    max-width: 80vw !important;
  }

  /* 优化滚动条样式 */
  .main::-webkit-scrollbar {
    width: 6px;
  }

  .main::-webkit-scrollbar-track {
    background: transparent;
    margin: 4px 0;
  }

  .main::-webkit-scrollbar-thumb {
    background: var(--neutral-200);
    border-radius: 6px;
    border: 2px solid transparent;
    background-clip: padding-box;
  }

  .main::-webkit-scrollbar-thumb:hover {
    background: var(--neutral-300);
    border: 2px solid transparent;
    background-clip: padding-box;
  }

  @keyframes scaleIn {
    from {
      opacity: 0;
      transform: scale(0.8);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  @keyframes fadeInScale {
    from {
      opacity: 0;
      transform: scale(0.9);
    }
    to {
      opacity: 1;
      transform: scale(1);
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

  .refresh-btn {
    position: relative;
    overflow: hidden;
  }

  .refresh-btn::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      to right,
      transparent,
      rgba(255, 255, 255, 0.1),
      transparent
    );
    transform: translateX(-100%);
  }

  .refresh-btn.loading::before {
    animation: shimmer 1.5s infinite;
  }

  .btn-text {
    display: inline-block;
    transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  @keyframes shimmer {
    100% {
      transform: translateX(100%);
    }
  }

  .spin {
    animation: spin 1s linear infinite;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

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
</style> 