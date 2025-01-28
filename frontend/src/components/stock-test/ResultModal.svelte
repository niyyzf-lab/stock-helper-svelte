<script lang="ts">
  import Modal from '../Modal.svelte'
  import type { TestResult } from '../../types/stock'
  import { Chart } from 'svelte-echarts'
  import { init, use } from 'echarts/core'
  import type { ComposeOption } from 'echarts/core'

  import { LineChart, CandlestickChart } from 'echarts/charts'
  import {
    GridComponent,
    TooltipComponent,
    ToolboxComponent,
    TitleComponent,
    DataZoomComponent
  } from 'echarts/components'
  import { CanvasRenderer } from 'echarts/renderers'
  import { onDestroy, onMount } from 'svelte'



  use([
    LineChart,
    CandlestickChart,
    GridComponent,
    TooltipComponent,
    ToolboxComponent,
    TitleComponent,
    DataZoomComponent,
    CanvasRenderer
  ])

  export let show = false
  export let result: TestResult | null = null
  export let onClose = () => {}
  export let onViewChart = () => {}

  let currentStep = -1
  let animationTimer: number | undefined
  let showResult = false
  let showBadge = false
  let showGrid = false
  let showButtons = false
  let chartInstance: any = null
  let mounted = false

  $: priceChangeClass = result?.priceChange 
    ? (result.priceChange > 0.03 ? 'up' : result.priceChange < -0.03 ? 'down' : 'shock')
    : ''

  $: isDataValid = result && 
    Array.isArray(result.klineData) && 
    Array.isArray(result.futureData) &&
    result.klineData.length > 0 &&
    currentStep >= 0

  $: chartOptions = {
    grid: {
      top: 0,
      right: 0,
      bottom: 0,
      left: 0,
      containLabel: false
    },
    tooltip: {
      trigger: 'axis' as const,
      axisPointer: {
        type: 'none' as const
      },
      position: function (pos: number[], params: any, dom: any, rect: any, size: any) {
        const obj: { [key: string]: number } = { top: 10 }
        obj[['left', 'right'][+(pos[0] < size.viewSize[0] / 2)]] = 5
        return obj
      },
      backgroundColor: 'rgba(255, 255, 255, 0.9)',
      borderWidth: 0,
      padding: [8, 12],
      textStyle: {
        color: '#1f2937'
      },
      formatter: function(params: any) {
        if (!Array.isArray(params)) {
          params = [params]
        }
        const item = params[0]
        if (!item) return ''
        
        const dataIndex = item.dataIndex
        const data = item.data
        const isHistory = dataIndex < (result?.klineData?.length || 0)
        
        let html = `<div style="padding: 0">`
        html += `<div style="margin-bottom: 4px;color:#6b7280;font-size:12px;">
          ${isHistory ? '历史数据' : '未来走势'}</div>`
        const date = new Date(item.name)
        html += `<div style="font-size:13px;font-weight:500;margin-bottom:8px;">
          ${date.getMonth() + 1}月${date.getDate()}日</div>`
        html += `<div style="margin-top:8px;font-size:12px;line-height:1.6;">
          <div>开盘：${data[0].toFixed(2)}</div>
          <div>收盘：${data[1].toFixed(2)}</div>
          <div>最低：${data[2].toFixed(2)}</div>
          <div>最高：${data[3].toFixed(2)}</div>
        </div>`
        html += `</div>`
        return html
      }
    },
    xAxis: {
      type: 'category' as const,
      data: isDataValid ? [
        ...result!.klineData.map(d => d.d),
        ...result!.futureData.slice(0, currentStep + 1).map(d => d.d)
      ] : [],
      show: false
    },
    yAxis: {
      type: 'value' as const,
      scale: true,
      show: false
    },
    dataZoom: [{
      type: 'inside' as const,
      start: Math.max(0, ((currentStep + (result?.klineData?.length ?? 0)) - 20) / ((result?.klineData?.length ?? 1) + currentStep) * 100),
      end: 100,
      zoomLock: true,
      moveOnMouseMove: true
    }, {
      show: false
    }],
    series: [
      {
        type: 'candlestick' as const,
        name: '历史数据',
        z: 10,
        data: isDataValid ? result!.klineData.map(item => [
          item.o,
          item.c,
          item.l,
          item.h
        ]) : [],
        itemStyle: {
          color: '#bfbfbf',
          color0: '#e0e0e0',
          borderColor: '#bfbfbf',
          borderColor0: '#e0e0e0'
        }
      },
      {
        type: 'candlestick' as const,
        name: '未来走势',
        z: 20,
        data: isDataValid ? [
          ...Array(result!.klineData.length).fill(['-', '-', '-', '-']),
          ...result!.futureData.slice(0, currentStep + 1).map(item => [
            item.o,
            item.c,
            item.l,
            item.h
          ])
        ] : [],
        itemStyle: {
          color: 'transparent',
          color0: '#3b82f6',
          borderColor: '#ef4444',
          borderColor0: '#3b82f6'
        }
      }
    ],
    animation: true,
    animationDuration: 400,
    animationEasing: 'cubicOut' as const,
    animationThreshold: 2000
  }

  $: if (show && result) {
    startAnimation()
  }

  onMount(() => {
    mounted = true
  })

  // 在 chartOptions 更新时手动更新图表
  $: if (chartInstance && isDataValid) {
    try {
      chartInstance.setOption({
        xAxis: {
          data: [
            ...result!.klineData.map(d => d.d),
            ...result!.futureData.slice(0, currentStep + 1).map(d => d.d)
          ]
        },
        series: [
          {
            data: result!.klineData.map(item => [
              item.o,
              item.c,
              item.l,
              item.h
            ])
          },
          {
            data: [
              ...Array(result!.klineData.length).fill(['-', '-', '-', '-']),
              ...result!.futureData.slice(0, currentStep + 1).map(item => [
                item.o,
                item.c,
                item.l,
                item.h
              ])
            ]
          }
        ]
      }, {
        replaceMerge: ['xAxis', 'series']
      })
    } catch (error) {
      console.warn('Failed to update chart:', error)
    }
  }

  function onChartInit(e: any) {
    chartInstance = e.detail
    if (chartInstance) {
      chartInstance.setOption(chartOptions)
    }
  }

  function updateChart() {
    if (!chartInstance || !result || !mounted) return
    
    try {
      const historyData = result.klineData.map(item => [
        item.o,
        item.c,
        item.l,
        item.h
      ])
      
      const futureData = [
        ...Array(result.klineData.length).fill(['-', '-', '-', '-']),
        ...result.futureData.slice(0, currentStep + 1).map(item => [
          item.o,
          item.c,
          item.l,
          item.h
        ])
      ]

      const dates = [
        ...result.klineData.map(d => d.d),
        ...result.futureData.slice(0, currentStep + 1).map(d => d.d)
      ]

      chartInstance.setOption({
        xAxis: {
          data: dates
        },
        series: [
          { data: historyData },
          { data: futureData }
        ]
      }, {
        replaceMerge: ['xAxis', 'series']
      })
    } catch (error) {
      console.warn('Failed to update chart:', error)
    }
  }

  function startAnimation() {
    currentStep = -1
    showResult = false
    showBadge = false
    showGrid = false
    showButtons = false
    
    if (animationTimer) {
      clearInterval(animationTimer)
      animationTimer = undefined
    }
    
    setTimeout(() => {
      if (!mounted) return
      showResult = true
      setTimeout(() => {
        if (!mounted) return
        showBadge = true
      }, 400)
    }, 2000)

    setTimeout(() => {
      if (!mounted) return
      currentStep = 0
      animationTimer = window.setInterval(() => {
        if (!result?.futureData || !mounted) {
          if (animationTimer) {
            clearInterval(animationTimer)
            animationTimer = undefined
          }
          return
        }
        
        if (currentStep === Math.floor(result.futureData.length * 0.7)) {
          showGrid = true
        }
        
        if (currentStep >= result.futureData.length - 1) {
          if (animationTimer) {
            clearInterval(animationTimer)
            animationTimer = undefined
          }
          
          if (mounted && chartInstance) {
            try {
              chartInstance.setOption({
                dataZoom: [{
                  start: 0,
                  end: 100,
                  zoomLock: false,
                  moveOnMouseMove: false
                }],
                grid: {
                  top: 40,
                  right: 10,
                  bottom: 60,
                  left: 50,
                  containLabel: true
                },
                xAxis: {
                  show: true,
                  axisLine: { show: true },
                  axisTick: { show: true },
                  axisLabel: { show: true }
                },
                yAxis: {
                  show: true,
                  axisLine: { show: true },
                  axisTick: { show: true },
                  axisLabel: { show: true }
                },
                series: [
                  {
                    type: 'candlestick',
                    name: '历史数据',
                    data: result.klineData.map(item => [
                      item.o,
                      item.c,
                      item.l,
                      item.h
                    ]),
                    itemStyle: {
                      color: '#bfbfbf',
                      color0: '#e0e0e0',
                      borderColor: '#bfbfbf',
                      borderColor0: '#e0e0e0'
                    }
                  },
                  {
                    type: 'candlestick',
                    name: '未来走势',
                    data: result.futureData.map(item => [
                      item.o,
                      item.c,
                      item.l,
                      item.h
                    ]),
                    itemStyle: {
                      color: 'transparent',
                      color0: '#3b82f6',
                      borderColor: '#ef4444',
                      borderColor0: '#3b82f6'
                    }
                  }
                ]
              }, {
                replaceMerge: ['dataZoom', 'grid', 'xAxis', 'yAxis', 'series']
              })
            } catch (error) {
              console.warn('Failed to update chart layout:', error)
            }
          }
          
          setTimeout(() => {
            if (mounted) showButtons = true
          }, 400)
          
          return
        }
        currentStep++
        updateChart()
      }, 400)
    }, 800)
  }

  onDestroy(() => {
    mounted = false
    if (animationTimer) {
      clearInterval(animationTimer)
      animationTimer = undefined
    }
    // 清理图表实例
    if (chartInstance) {
      try {
        chartInstance.dispose()
      } catch (error) {
        console.warn('Failed to dispose chart:', error)
      }
      chartInstance = null
    }
  })
</script>

<Modal 
  show={show}
  title="预测结果"
  on:close={onClose}
  class_="result-modal"
>
  <div class="result-container">
    <div class="content-wrapper">
      <div class="chart-section">
        <Chart 
          {init} 
          options={chartOptions} 
          on:init={onChartInit}
        />
      </div>

      <div class="info-section">
        <div class="result-content">
          <div class="result-badge {result?.correct ? 'correct' : 'wrong'}" 
            class:slide-down={showBadge} 
            class:minimized={showGrid}>
            <div class="badge-icon">
              {#if result?.correct}
                <svg class="check" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M20 6L9 17L4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              {:else}
                <svg class="cross" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M18 6L6 18M6 6L18 18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              {/if}
            </div>
            <div class="badge-text">
              <div class="badge-title">{result?.correct ? '预测正确' : '预测错误'}</div>
              <div class="badge-desc">{result?.correct ? '恭喜您，判断准确！' : '继续加油，下次会更好！'}</div>
            </div>
          </div>

          <div class="result-grid" class:fade-in={showGrid}>
            <div class="result-item">
              <div class="item-label">您的预测</div>
              <div class="item-value prediction {result?.direction}">
                {#if result?.direction === 'up'}上涨
                {:else if result?.direction === 'down'}下跌
                {:else}震荡{/if}
              </div>
            </div>

            <div class="result-item">
              <div class="item-label">实际走势</div>
              <div class="item-value prediction {result?.actualDirection}">
                {#if result?.actualDirection === 'up'}上涨
                {:else if result?.actualDirection === 'down'}下跌
                {:else}震荡{/if}
              </div>
            </div>

            <div class="result-item highlight">
              <div class="item-label">涨跌幅</div>
              <div class="item-value change {priceChangeClass}">
                {result?.priceChange ? (result.priceChange * 100).toFixed(2) : 0}%
              </div>
            </div>

            <div class="result-item">
              <div class="item-label">观察天数</div>
              <div class="item-value">{result?.daysCount || 0}天</div>
            </div>
          </div>
        </div>

        <div class="button-group" class:fade-in={showButtons}>
          <button class="next-btn" on:click={onClose}>
            下一题
          </button>
          <button class="view-chart-btn" on:click={onViewChart}>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M15 3h6v6M9 21H3v-6M21 3l-7 7M3 21l7-7"/>
            </svg>
            <span>查看完整图表</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</Modal>

<style>
  .result-container {
    width: 1080px;
    min-height: 680px;
    height: 100%;
    padding: 0;
    background: #fff;
    display: flex;
    flex-direction: column;
    border-radius: var(--radius-lg);
    overflow: hidden;
  }

  .content-wrapper {
    display: flex;
    flex-direction: row;
    height: 100%;
  }

  /* 图表区域 */
  .chart-section {
    flex: 1;
    height: 100%;
    padding: 40px;
    position: relative;
  }

  /* 添加分隔线 */
  .chart-section::after {
    content: '';
    position: absolute;
    right: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 1px;
    height: 80%; /* 不完全等长 */
    background: linear-gradient(
      to bottom,
      transparent,
      var(--border-color) 20%,
      var(--border-color) 80%,
      transparent
    );
  }

  /* 结果区域 */
  .info-section {
    width: 360px;
    padding: 40px 32px;
    display: flex;
    flex-direction: column;
    justify-content: space-between; /* 确保内容分布到整个高度 */
  }

  /* 结果内容容器 */
  .result-content {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  /* 徽章样式 */
  .result-badge {
    text-align: center;
    padding: 0 0 32px;
    margin-bottom: 32px;
    border-bottom: 1px solid var(--border-color);
    opacity: 0;
    transform: translateY(-20px);
    transition: all 0.6s cubic-bezier(0.16, 1, 0.3, 1);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 16px;
  }

  .result-badge.minimized {
    padding: 0 0 24px;
    margin-bottom: 24px;
    justify-content: flex-start;
  }

  .result-badge.slide-down {
    opacity: 1;
    transform: translateY(0);
  }

  .badge-icon {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .result-badge.correct .badge-icon {
    color: #22c55e;
  }

  .result-badge.wrong .badge-icon {
    color: #ef4444;
  }

  .badge-title {
    font-size: 1.125rem;
    font-weight: 600;
    margin-bottom: 4px;
    color: var(--text-primary);
  }

  .badge-desc {
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .badge-text {
    text-align: left;
  }

  .result-badge:not(.minimized) .badge-text {
    text-align: center;
  }

  /* 结果网格 */
  .result-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 24px;
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.6s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .result-grid.fade-in {
    opacity: 1;
    transform: translateY(0);
  }

  .result-item {
    position: relative;
  }

  .result-item.highlight {
    padding-left: 16px;
  }

  .result-item.highlight::before {
    content: '';
    position: absolute;
    left: 0;
    top: 4px;
    bottom: 4px;
    width: 3px;
    background: var(--primary-500);
    border-radius: 3px;
  }

  .item-label {
    font-size: 0.75rem;
    font-weight: 500;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--text-secondary);
    margin-bottom: 8px;
  }

  .item-value {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--text-primary);
    line-height: 1.2;
  }

  /* 按钮组 */
  .button-group {
    padding-top: 40px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-top: 0; /* 移除 margin-top: auto */
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .button-group.fade-in {
    opacity: 1;
    transform: translateY(0);
  }

  .next-btn {
    padding: 14px 24px;
    background: var(--primary-600);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 0.9375rem;
    font-weight: 600;
    letter-spacing: 0.3px;
    transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
    transform: translateY(0);
  }

  .next-btn:hover {
    background: var(--primary-700);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(var(--primary-600-rgb), 0.2);
  }

  .view-chart-btn {
    padding: 14px 24px;
    border: 1px solid var(--border-color);
    background: transparent;
    border-radius: 8px;
    color: var(--text-secondary);
    font-size: 0.9375rem;
    font-weight: 500;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
    transform: translateY(0);
  }

  .view-chart-btn:hover {
    border-color: var(--primary-200);
    color: var(--primary-600);
    transform: translateY(-1px);
    background: var(--primary-50);
  }

  /* 动画相关 */
  .check, .cross {
    stroke-dasharray: 60;
    stroke-dashoffset: 60;
    transition: stroke-dashoffset 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .result-badge.slide-down .check,
  .result-badge.slide-down .cross {
    stroke-dashoffset: 0;
  }

  /* 涨跌颜色 */
  .item-value.prediction.up,
  .item-value.change.up {
    color: #ef4444;
  }

  .item-value.prediction.down,
  .item-value.change.down {
    color: #22c55e;
  }

  .item-value.prediction.shock {
    color: #eab308;
  }

  :global(.modal-content.result-modal) {
    background: var(--neutral-50) !important;
  }
</style> 