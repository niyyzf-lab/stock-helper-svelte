<script lang="ts">
  import logo from '../assets/images/logo-universal.png'
  import { 
    LayoutDashboard, BarChart2, LineChart, 
    Search, Filter, Star, History, Settings, RefreshCw,
    ChevronLeft, ChevronRight, Brain, Play
  } from 'lucide-svelte'
  import { location, link } from 'svelte-spa-router'
  import { onDestroy, afterUpdate, onMount } from 'svelte'
  import tippy from 'tippy.js'
  import 'tippy.js/dist/tippy.css'
  import 'tippy.js/themes/light.css'
  import { EventsOn, EventsOff } from '../../wailsjs/runtime'
  import { sidebarStore } from '../stores/sidebar'
  import { toastStore } from '../stores/toast'
  import { fade, fly, slide } from 'svelte/transition'

  export let updateStatus: any = null
  let progressBarElement: HTMLElement
  let tippyInstance: any

  // 移除本地状态，使用 store
  $: isCollapsed = $sidebarStore.isCollapsed

  const mainNavItems = [
    { icon: LayoutDashboard, label: '概览', href: '/' },
    { icon: BarChart2, label: '股票分析', href: '/analysis' },
    { icon: LineChart, label: '股票追踪', href: '/track' },
    { icon: Filter, label: '股票筛选', href: '/filter' },
    { icon: Play, label: '策略测试', href: '/test' },
    { icon: History, label: '筛选历史', href: '/filter-history' },
  ]

  const myStockItems = [
    { icon: Star, label: '自选股票' },
  ]

  function formatProgress(status: any) {
    if (!status) return ''
    return `
      <div class="tippy-content">
        <div class="detail-item">
          <span class="label">进度</span>
          <span class="value">${status.progress.toFixed(1)}% (${status.completed}/${status.total})</span>
        </div>
        <div class="detail-item">
          <span class="label">速度</span>
          <span class="value">${status.speed.toFixed(1)} 个/秒</span>
        </div>
        <div class="detail-item">
          <span class="label">剩余时间</span>
          <span class="value">${status.estimateTime}秒</span>
        </div>
        ${status.errorCount > 0 ? `
          <div class="detail-item error">
            <span class="label">错误数</span>
            <span class="value">${status.errorCount}</span>
          </div>
        ` : ''}
      </div>
    `
  }

  afterUpdate(() => {
    if (progressBarElement && updateStatus) {
      if (tippyInstance) {
        tippyInstance.setContent(formatProgress(updateStatus))
      } else {
        tippyInstance = tippy(progressBarElement, {
          content: formatProgress(updateStatus),
          allowHTML: true,
          theme: 'light',
          placement: 'right',
          interactive: true,
          animation: 'scale',
          duration: [200, 0],
          appendTo: () => document.body
        })
      }
    }
  })

  onDestroy(() => {
    if (tippyInstance) {
      tippyInstance.destroy()
    }
  })

  onMount(() => {
    // 监听更新状态
    EventsOn('update:status', (status: any) => {
      console.log('收到更新状态:', status)
      updateStatus = status
    })

    // 返回清理函数
    return () => {
      EventsOff('update:status')
    }
  })

  // 添加调试日志
  $: console.log('Current hash:', $location)

  // 添加更新数据函数
  async function updateData() {
    if (updateStatus?.isUpdating) return
    try {
      toastStore.promise(
        (window as any).go.main.App.UpdateStockData(),
        {
          loading: '正在更新数据...',
          success: '数据更新成功！',
          error: '数据更新失败，请重试'
        }
      )
    } catch (err) {
      toastStore.error('更新数据时发生错误')
    }
  }



  // 修改切换函数使用 store
  function toggleSidebar() {
    sidebarStore.toggle()
  }

  // 修改动画控制
  let mounted = false
  onMount(() => {
    // 使用 RAF 确保在下一帧执行，避免闪烁
    requestAnimationFrame(() => {
      mounted = true
    })
  })
</script>

<aside class="sidebar" class:collapsed={$sidebarStore.isCollapsed} class:mounted={mounted}>
  <header style="--wails-draggable:drag">
    <div class="logo">
      <div class="logo-icon">
        <img src={logo} alt="Logo" />
      </div>
      {#if !isCollapsed}
        <span>股票助手</span>
      {/if}
    </div>
  </header>

  <button 
    class="collapse-btn" 
    on:click={toggleSidebar}
  >
    <svelte:component this={$sidebarStore.isCollapsed ? ChevronRight : ChevronLeft} size={16} strokeWidth={1.5} />
  </button>

  <nav>
    <div class="nav-group main-nav">
      {#each mainNavItems as item, i}
        <a 
          href={'#' + item.href}
          class="nav-item"
          class:active={$location === item.href}
          use:link
          style="--delay: {i * 50 + 200}ms"
        >
          <div class="icon-wrapper">
            <svelte:component this={item.icon} size={16} strokeWidth={1.5} />
          </div>
          {#if !isCollapsed}
            <span>{item.label}</span>
          {/if}
        </a>
      {/each}

      {#if updateStatus?.isUpdating}
        <div class="nav-item update-status">
          <div class="progress-section">
            <div class="progress-bar-wrapper" bind:this={progressBarElement}>
              <div class="progress-bar" style="width: {updateStatus?.progress || 0}%">
                {#if updateStatus?.errorCount > 0}
                  <div class="error-indicator" style="width: {(updateStatus.errorCount / updateStatus.total) * 100}%"></div>
                {/if}
              </div>
            </div>
            {#if !isCollapsed}
              <div class="current-info">
                <span class="current">{updateStatus?.current || '-'}</span>
                <span class="progress">{updateStatus?.progress?.toFixed(1) || 0}%</span>
              </div>
            {/if}
          </div>
        </div>
      {/if}
    </div>

    <div class="nav-group">
      <div class="group-header">
        {#if !isCollapsed}
          <span class="group-title">我的股票</span>
          <span class="badge">2</span>
        {/if}
      </div>
      {#each myStockItems as item, i}
        <button 
          class="nav-item"
          style="--delay: {i * 50 + 500}ms"
        >
          <div class="icon-wrapper">
            <svelte:component this={item.icon} size={16} strokeWidth={1.5} />
          </div>
          {#if !isCollapsed}
            <span>{item.label}</span>
            <div class="arrow" />
          {/if}
        </button>
      {/each}
    </div>
  </nav>

  <div class="sync-wrapper">
    <div class="button-group">
      <button 
        class="nav-item sync-btn" 
        on:click={updateData} 
        disabled={updateStatus?.isUpdating}
        style="--delay: 600ms"
      >
        <div class="icon-wrapper">
          <RefreshCw size={16} strokeWidth={1.5} />
        </div>
        {#if !isCollapsed}
          <span>同步数据</span>
        {/if}
      </button>
    </div>
  </div>

  <footer>
    <button class="nav-item" style="--delay: 700ms">
      <div class="icon-wrapper">
        <Settings size={16} strokeWidth={1.5} />
      </div>
      {#if !isCollapsed}
        <span>设置</span>
      {/if}
    </button>
  </footer>
</aside>

<style>
  .sidebar {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background: #fff;
    font-family: var(--font-sans);
    --text-primary: var(--neutral-800);
    --text-secondary: var(--neutral-400);
    --hover-bg: var(--neutral-50);
    --active-bg: var(--neutral-100);
    --border-color: var(--neutral-100);
    --accent-color: var(--primary-500);
    --spacing-unit: 12px;
    position: relative;
    border-right: 1px solid var(--border-color);
    width: 240px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    opacity: 0;
    transform: translateX(-20px);
  }

  .sidebar.mounted {
    opacity: 1;
    transform: translateX(0);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .sidebar.collapsed {
    width: 72px;
  }

  .sidebar.collapsed .nav-item {
    width: 48px;
    height: 48px;
    padding: 8px;
    margin: 0 auto;
    justify-content: center;
  }

  .collapse-btn {
    position: absolute;
    right: -12px;
    top: 20px;
    width: 24px;
    height: 24px;
    background: #fff;
    border: 1px solid var(--border-color);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: var(--text-secondary);
    z-index: 10;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
    opacity: 0;
    transform: translateX(20px);
  }

  .mounted .collapse-btn {
    opacity: 1;
    transform: translateX(0);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    transition-delay: 300ms;
  }

  .collapse-btn:hover {
    color: var(--text-primary);
    background: var(--hover-bg);
    transform: scale(1.1);
  }

  .collapse-btn:active {
    transform: scale(0.9);
  }

  .sidebar.collapsed nav {
    padding: 20px 12px;
  }

  /* 添加微妙的侧边阴影 */
  .sidebar::after {
    content: '';
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    width: 1px;
    background: linear-gradient(to bottom,
      rgba(0,0,0,0.02),
      rgba(0,0,0,0.04) 20%,
      rgba(0,0,0,0.04) 80%,
      rgba(0,0,0,0.02)
    );
  }

  header {
    height: 60px;
    display: flex;
    align-items: center;
    padding: 0 var(--spacing-unit);
    position: relative;
    user-select: none;
    background: rgba(255, 255, 255, 0.7);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border-bottom: 1px solid var(--border-color);
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 0 12px;
    width: 100%;
    opacity: 0;
    transform: translateY(-20px);
    white-space: nowrap;
    min-width: 0;
  }

  .mounted .logo {
    opacity: 1;
    transform: translateY(0);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .logo-icon {
    width: 32px;
    height: 32px;
    min-width: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #f8fafc 0%, #fff 100%);
    border-radius: 8px;
    box-shadow: 
      0 2px 4px rgba(0,0,0,0.02),
      0 0 0 1px rgba(0,0,0,0.02);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .logo-icon:hover {
    transform: scale(1.05);
  }

  .logo img {
    width: 24px;
    height: 24px;
  }

  .logo span {
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
    letter-spacing: -0.1px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .sidebar.collapsed .logo {
    justify-content: center;
    padding: 0;
  }

  nav {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 28px;
    padding: 20px var(--spacing-unit);
  }

  .nav-group {
    display: flex;
    flex-direction: column;
    gap: 1px;
  }

  .group-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 12px;
    margin-bottom: 8px;
    min-height: 20px;
  }

  .sidebar.collapsed .group-header {
    height: 8px;
    margin-bottom: 16px;
    opacity: 0.5;
    background: var(--border-color);
    margin: 8px auto;
    width: 24px;
    padding: 0;
    border-radius: 4px;
  }

  .group-title {
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
    color: var(--text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.5px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .badge {
    font-size: 10px;
    font-weight: 500;
    color: var(--text-secondary);
    background: var(--hover-bg);
    padding: 2px 8px;
    border-radius: 10px;
    border: 1px solid var(--border-color);
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    min-width: 24px;
    text-align: center;
    flex-shrink: 0;
  }

  .badge:hover {
    transform: scale(1.1);
  }

  .nav-item {
    position: relative;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 12px;
    color: var(--text-primary);
    text-decoration: none;
    border-radius: 8px;
    font-size: 13px;
    font-weight: 450;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    white-space: nowrap;
    overflow: hidden;
    background: none;
    border: none;
    width: 100%;
    text-align: left;
    cursor: pointer;
    opacity: 0;
    transform: translateX(-20px);
  }

  .mounted .nav-item {
    opacity: 1;
    transform: translateX(0);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    transition-delay: var(--delay);
  }

  .nav-item:hover {
    background: var(--hover-bg);
    transform: translateX(2px);
  }

  .nav-item.active {
    background: var(--active-bg);
    color: var(--accent-color);
    font-weight: 500;
    transform: none;
  }

  .nav-item:active {
    transform: scale(0.98);
  }

  .nav-item:hover .icon-wrapper {
    transform: scale(1.1);
  }

  .nav-item.active .icon-wrapper {
    transform: none;
    color: var(--accent-color);
  }

  .icon-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    min-width: 32px;
    border-radius: 8px;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .nav-item :global(svg) {
    opacity: 0.8;
    transition: all 0.15s ease;
  }

  .nav-item:hover :global(svg) {
    opacity: 1;
  }

  .arrow {
    position: absolute;
    right: 12px;
    width: 5px;
    height: 5px;
    border-right: 1.5px solid var(--text-secondary);
    border-bottom: 1.5px solid var(--text-secondary);
    transform: rotate(-45deg);
    opacity: 0;
    transition: all 0.2s ease;
  }

  .nav-item:hover .arrow {
    opacity: 0.5;
    right: 10px;
  }

  footer {
    padding: var(--spacing-unit);
    border-top: 1px solid var(--border-color);
  }

  footer .nav-item {
    color: var(--text-secondary);
  }

  footer .nav-item:hover {
    color: var(--text-primary);
  }

  /* 主导航组的特殊样式 */
  .main-nav .nav-item {
    padding: 9px 12px;
  }

  /* 优化活跃状态 */
  .nav-item.active :global(svg) {
    opacity: 1;
  }

  /* 更新状态样式 */
  .update-status {
    padding: 12px !important;
    margin: 8px 0;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
  }

  .progress-section {
    width: 100%;
  }

  .progress-bar-wrapper {
    position: relative;
    height: 4px;
    background: var(--neutral-100);
    border-radius: 2px;
    overflow: hidden;
    cursor: pointer;
    margin-bottom: 8px;
  }

  .progress-bar {
    position: relative;
    height: 100%;
    background: linear-gradient(90deg, 
      var(--primary-200), 
      var(--primary-500), 
      var(--primary-200));
    background-size: 200% 100%;
    border-radius: 2px;
    transition: width 0.3s ease;
    animation: shimmer 2s infinite;
  }

  .error-indicator {
    position: absolute;
    top: 0;
    right: 0;
    height: 100%;
    background: rgba(239, 68, 68, 0.5);
    transition: width 0.3s ease;
  }

  .current-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 11px;
  }

  .current {
    color: var(--text-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    flex: 1;
    margin-right: 8px;
  }

  .progress {
    color: var(--primary-500);
    font-weight: 600;
  }

  @keyframes shimmer {
    0% {
      background-position: 200% 0;
    }
    100% {
      background-position: -200% 0;
    }
  }

  .button-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .train-btn {
    width: 100%;
    color: var(--text-secondary);
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 6px;
    font-weight: 450;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .train-btn:hover {
    color: var(--text-primary);
    background: var(--hover-bg);
    border-color: var(--neutral-200);
    transform: translateY(-1px);
  }

  .train-btn:active {
    transform: scale(0.98);
  }

  .train-btn .icon-wrapper {
    color: var(--primary-500);
  }

  /* 修改 sync-wrapper 样式以适应两个按钮 */
  .sync-wrapper {
    padding: 0 var(--spacing-unit);
    margin-bottom: 8px;
  }

  .sync-btn {
    width: 100%;
    color: var(--text-secondary);
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 6px;
    font-weight: 450;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .sync-btn:hover:not(:disabled) {
    color: var(--text-primary);
    background: var(--hover-bg);
    border-color: var(--neutral-200);
    transform: translateY(-1px);
  }

  .sync-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .sync-btn:active:not(:disabled) {
    transform: scale(0.98);
  }

  /* ... rest of the styles ... */
</style> 
