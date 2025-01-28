<script lang="ts">
  import Router, { link, location } from 'svelte-spa-router'
  import MainLayout from './layouts/MainLayout.svelte'
  import Sidebar from './components/Sidebar.svelte'
  import Overview from './pages/Overview.svelte'
  import StockFilter from './pages/StockFilter.svelte'
  import FilterHistory from './pages/FilterHistory.svelte'
  import StockQuery from './pages/StockQuery.svelte'
  import StockTest from './pages/StockTest.svelte'
  import { Toaster } from 'svelte-french-toast'
  import { toastStore } from './stores/toast'
  import { onMount } from 'svelte'
  import { EventsOn } from '../wailsjs/runtime'



  // 路由配置
  const routes = {
    // 注意：这里的路由路径不需要带 #
    '/': Overview,
    '/filter': StockFilter,
    '/analysis': Overview,
    '/track': Overview,
    '/filter-history': FilterHistory,
    '/stock/:code': StockQuery,
    '/test': StockTest
  }

  // 调试日志
  $: console.log('Current route:', $location)

  onMount(() => {
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
  });
</script>

<MainLayout>
  <div slot="sidebar">
    <Sidebar />
  </div>
  
  <Toaster />
  <Router {routes} />
</MainLayout>

<style>
  :global(.toast) {
    --toast-background: var(--surface);
    --toast-color: var(--text-primary);
    --toast-border: 1px solid var(--border-color);
    --toast-success: var(--success-500);
    --toast-error: var(--error-500);
  }
</style>
