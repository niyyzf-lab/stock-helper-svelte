<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { createChart } from 'lightweight-charts';
    import type { LineWidth } from 'lightweight-charts';  // 修改导入方式
  
    export let code: string;
    export let freq: string = 'dh';  // K线频率，默认日线
    export let maType: string = 'SMA';  // 默认使用 SMA
    export let maPeriods: number[] = [6, 18, 98];  // 默认周期
    export let endTime: string = '';  // 格式应该是 "2024-01-01" 或 "2024-01-01 12:34:56"

    let chartContainer: HTMLElement;
    let macdContainer: HTMLElement;  // 新增MACD容器引用
    let chartInstance: any = null;
    let macdChartInstance: any = null;  // 新增MACD图表实例
    let resizeObserver: ResizeObserver | null = null;
    let loading = false;
    let error: string | null = null;
    let series = {
      candlestick: null as any,
      volume: null as any,
      ma: [] as any[],  // 添加均线数组
      macd: {
        dif: null as any,
        dea: null as any,
        histogram: null as any
      }
    };

    // 添加价格信息状态
    let priceInfo = {
        open: 0,
        high: 0,
        low: 0,
        close: 0,
        change: 0,
        volume: 0
    };

    // 在 script 部分添加工具提示状态
    let tooltip = {
      visible: false,
      x: 0,
      y: 0,
      data: null as any
    };

    // 修改图表基础配置
    const chartOptions = {
      layout: {
        background: { color: '#ffffff' },
        textColor: '#4b5563',
        fontSize: 12,
        fontFamily: 'system-ui, -apple-system, sans-serif'
      },
      grid: {
        vertLines: { color: 'rgba(229, 231, 235, 0.5)' },
        horzLines: { color: 'rgba(229, 231, 235, 0.5)' }
      },
      rightPriceScale: {
        borderVisible: false,
        autoScale: true,      // 启用自动缩放
        scaleMargins: {
          top: 0.1,           // 调整上边距为10%
          bottom: 0.1         // 调整下边距为10%
        },
        borderColor: 'rgba(229, 231, 235, 0.8)',
        textColor: '#4b5563'
      },
      timeScale: {
        borderVisible: false,
        timeVisible: true,
        secondsVisible: false,
        rightOffset: 12,
        rightBarStaysOnScroll: true,  // 保持右侧空白
        fixLeftEdge: false,
        fixRightEdge: false,
        lockVisibleTimeRangeOnResize: false,
        barSpacing: 12,
        minBarSpacing: 8
      },
      crosshair: {
        mode: 1,
        vertLine: {
          color: 'rgba(75, 85, 99, 0.3)',
          width: 1 as LineWidth,
          style: 2
        },
        horzLine: {
          color: 'rgba(75, 85, 99, 0.3)',
          width: 1 as LineWidth,
          style: 2
        }
      },
      handleScroll: {
        mouseWheel: true,
        pressedMouseMove: true,     // 启用鼠标拖动
        horzTouchDrag: true,        // 启用触摸拖动
        vertTouchDrag: false
      },
      handleScale: {
        axisPressedMouseMove: false,
        mouseWheel: true,
        pinch: true,
        minBarSpacing: 0.1
      }
    };

    // 添加时间转换函数
    function convertTime(timeStr: string, freq: string): number {
      try {
        if (freq === 'dh') {  // 日线
          return Math.floor(new Date(`${timeStr} 00:00:00 GMT+0800`).getTime() / 1000);
        } else {  // 分时
          // 分时数据格式可能是 "2024-01-01 09:30:00" 这样的格式
          return Math.floor(new Date(`${timeStr} GMT+0800`).getTime() / 1000);
        }
      } catch (err) {
        console.error('时间转换失败:', timeStr, err);
        return 0;
      }
    }

    // 获取K线数据
    async function fetchKLineData() {
      try {
        const data = await (window as any).go.main.App.GetKLineData(code, freq);
        console.log('原始K线数据:', data);
        
        if (!data || data.length === 0) throw new Error('没有数据');
        return data;
      } catch (err) {
        console.error('获取K线数据失败:', err);
        throw err;
      }
    }

    // 计算均线
    async function calculateMAs(prices: number[]) {
      const maLines = [];
      for (const period of maPeriods) {
        try {
          const maData = await (window as any).go.main.App.CalculateMA(
            prices, 
            maType.toLowerCase(), 
            period
          );
          maLines.push(maData);
        } catch (err) {
          console.error(`计算${maType}${period}失败:`, err);
        }
      }
      return maLines;
    }

    // 添加 MACD 计算函数
    async function calculateMACD(prices: number[]) {
      try {
        return await (window as any).go.main.App.CalculateMACD(prices);
      } catch (err) {
        console.error('计算MACD失败:', err);
        return null;
      }
    }

    // 修改同步函数，减少等待时间
    function syncTimeScales(sourceScale: any, targetScale: any, range: any) {
        if (!range || !targetScale) return;
        
        // 直接更新，不使用 requestAnimationFrame
        const currentRange = targetScale.getVisibleLogicalRange();
        if (!currentRange || 
            currentRange.from !== range.from || 
            currentRange.to !== range.to) {
            targetScale.setVisibleLogicalRange(range);
        }
    }

    // 修改时间轴同步部分
    let syncInProgress = false;
    const SYNC_DELAY = 16;  // 减少延迟到16ms (一帧)

    // 修改初始化函数，将十字线移动事件订阅移到MACD代码块外
    async function initChart() {
      if (!chartContainer || !macdContainer || !code) return;
      
      try {
        loading = true;
        error = null;

        await cleanup();
        let klineData = await fetchKLineData();
        
        // 如果设置了截至时间，过滤数据
        if (endTime) {
          const endTimestamp = new Date(endTime).getTime();
          klineData = klineData.filter(item => {
            const itemTime = freq === 'dh' ? 
              new Date(`${item.d} 00:00:00`).getTime() :  // 日线
              new Date(item.d).getTime();                 // 分时
            return itemTime <= endTimestamp;
          });

          // 如果过滤后没有数据，抛出错误
          if (klineData.length === 0) {
            throw new Error('截至时间之前没有数据');
          }
        }

        // 创建主图表实例
        const containerHeight = chartContainer.parentElement.clientHeight;
        chartInstance = createChart(chartContainer, {
          ...chartOptions,
          width: chartContainer.clientWidth,
          height: Math.floor(containerHeight * 0.7)
        });

        // 创建MACD图表实例
        macdChartInstance = createChart(macdContainer, {
          ...chartOptions,
          width: macdContainer.clientWidth,
          height: Math.floor(containerHeight * 0.3)  // 使用父容器高度计算
        });

        // 修改MACD图表的价格轴配置
        macdChartInstance.applyOptions({
          rightPriceScale: {
            scaleMargins: {
              top: 0.1,     // MACD图表保持适当边距
              bottom: 0.1
            },
            borderVisible: false,
            textColor: '#4b5563'
          },
          timeScale: {
            visible: false,    // 隐藏时间轴
            borderVisible: false
          }
        });

        // 设置K线图
        series.candlestick = chartInstance.addCandlestickSeries({
          upColor: 'transparent',      // 上涨时内部透明
          downColor: '#3b82f6',       // 下跌时使用蓝色实心
          borderVisible: true,
          borderUpColor: '#ef4444',   // 上涨时边框为红色
          borderDownColor: '#3b82f6', // 下跌时边框为蓝色
          wickUpColor: '#ef4444',     // 上影线为红色
          wickDownColor: '#3b82f6'    // 下影线为蓝色
        });

        // 设置成交量图
        series.volume = chartInstance.addHistogramSeries({
          priceFormat: {
            type: 'volume'
          },
          priceScaleId: 'volume',
          scaleMargins: {
            top: 0.85,    // 调整成交量位置，占用底部15%空间
            bottom: 0.02  // 留出小边距
          }
        });

        // 设置成交量坐标轴
        chartInstance.priceScale('volume').applyOptions({
          scaleMargins: {
            top: 0.85,
            bottom: 0.02
          },
          visible: false
        });

        // 转换并设置K线数据
        const candleData = klineData.map(item => ({
          time: convertTime(item.d, freq),
          open: item.o,
          high: item.h,
          low: item.l,
          close: item.c
        }));
        
        // 转换并设置成交量数据
        const volumeData = klineData.map(item => ({
          time: convertTime(item.d, freq),
          value: item.v,
          color: item.c >= item.o ? 
            'rgba(239, 68, 68, 0.3)' :  // 上涨红色 30% 透明度
            'rgba(59, 130, 246, 0.3)'   // 下跌蓝色 30% 透明度
        }));

        series.candlestick.setData(candleData, {
          lastValueVisible: true,
          priceLineVisible: true,
          priceLineWidth: 1,
          priceLineColor: 'rgba(0, 0, 0, 0.3)',
          priceLineStyle: 2,  // 虚线
          lastPriceAnimation: 1  // 启用价格动画
        });
        series.volume.setData(volumeData);

        // 计算均线
        const prices = klineData.map(item => item.c);
        const maLines = await calculateMAs(prices);

        // 添加均线
        const maColors = [
          'rgba(59, 130, 246, 0.8)',  // 蓝色 80% 透明度
          'rgba(236, 72, 153, 0.8)',  // 粉色 80% 透明度
          'rgba(139, 92, 246, 0.8)'   // 紫色 80% 透明度
        ];
        series.ma = maLines.map((data, index) => {
          const lineSeries = chartInstance.addLineSeries({
            color: maColors[index % maColors.length],
            lineWidth: 1.5,                   
            title: `${maType}${maPeriods[index]}`,
            priceLineVisible: false,
            lastValueVisible: true,
            crosshairMarkerVisible: false,    // 移除选中原点
            lineStyle: 0,
            disableApproximation: false,
            priceFormat: {
              type: 'price',
              precision: 2,
              minMove: 0.01
            }
          });

          // 过滤掉无效的数据点
          const validData = data
            .map((value: number, i: number) => ({
              time: convertTime(klineData[i].d, freq),
              value: value === 0 ? undefined : value
            }))
            .filter(item => item.value !== undefined);

          lineSeries.setData(validData);
          return lineSeries;
        });

        // 添加十字线移动事件订阅 (移到这里)
        chartInstance.subscribeCrosshairMove(param => {
          // 同步MACD图表的十字线
          if (macdChartInstance && param.point) {
            macdChartInstance.setCrosshairPosition(
              param.time,
              param.point.y / chartContainer.clientHeight * macdContainer.clientHeight,
              null  // 移除 param.seriesPrices.get() 调用
            );
          } else if (macdChartInstance) {
            macdChartInstance.clearCrosshairPosition();
          }

          // 处理工具提示
          if (param.point && param.time) {  // 确保有时间和位置
            const point = klineData.find(d => convertTime(d.d, freq) === param.time);
            if (point) {
              // 更新价格信息面板
              priceInfo = {
                open: point.o,
                high: point.h,
                low: point.l,
                close: point.c,
                change: point.zd || ((point.c - point.o) / point.o * 100),
                volume: point.v
              };

              // 获取图表容器的位置
              const rect = chartContainer.getBoundingClientRect();
              
              // 计算工具提示的位置（相对于视口）
              let tooltipX = rect.left + param.point.x + 12;
              let tooltipY = rect.top + param.point.y - 8;

              // 设置固定的工具提示尺寸
              const tooltipWidth = 160;
              const tooltipHeight = 120;
              
              // 调整位置，避免超出视口
              const windowWidth = window.innerWidth;
              const windowHeight = window.innerHeight;

              if (tooltipX + tooltipWidth > windowWidth - 5) {
                tooltipX = rect.left + param.point.x - tooltipWidth - 12;
              }

              if (tooltipY < 5) {
                tooltipY = rect.top + param.point.y + 8;
              }

              if (tooltipY + tooltipHeight > windowHeight - 5) {
                tooltipY = windowHeight - tooltipHeight - 5;
              }

              tooltipX = Math.max(5, tooltipX);

              // 更新工具提示
              tooltip = {
                visible: true,
                x: tooltipX,
                y: tooltipY,
                data: {
                  time: point.d,
                  price: point.c.toFixed(2),
                  change: (point.zd || ((point.c - point.o) / point.o * 100)).toFixed(2),
                  volume: (point.v / 10000).toFixed(2),
                  amount: (point.e / 100000000).toFixed(2)
                }
              };
            }
          } else {
            tooltip.visible = false;
          }
        });

        // 计算并添加MACD
        const macdResult = await calculateMACD(prices);
        if (macdResult) {
          // 添加MACD柱状图
          series.macd.histogram = macdChartInstance.addHistogramSeries({
            color: '#6b7280',
            priceFormat: {
              type: 'price',
              precision: 4
            }
          });

          // 添加DIF线
          series.macd.dif = macdChartInstance.addLineSeries({
            color: 'rgba(37, 99, 235, 0.8)',
            lineWidth: 1.5,
            title: 'DIF'
          });

          // 添加DEA线
          series.macd.dea = macdChartInstance.addLineSeries({
            color: 'rgba(236, 72, 153, 0.8)',
            lineWidth: 1.5,
            title: 'DEA'
          });

          // 设置数据
          const macdData = macdResult.MACD.map((value, i) => ({
            time: convertTime(klineData[i].d, freq),
            value: value,
            color: value >= 0 ? 'rgba(239, 68, 68, 0.6)' : 'rgba(59, 130, 246, 0.6)'
          }));

          const difData = macdResult.DIF.map((value, i) => ({
            time: convertTime(klineData[i].d, freq),
            value: value
          }));

          const deaData = macdResult.DEA.map((value, i) => ({
            time: convertTime(klineData[i].d, freq),
            value: value
          }));

          series.macd.histogram.setData(macdData);
          series.macd.dif.setData(difData);
          series.macd.dea.setData(deaData);

          // 修改时间轴同步
          const mainTimeScale = chartInstance.timeScale();
          const macdTimeScale = macdChartInstance.timeScale();

          // 修改时间轴同步部分
          mainTimeScale.subscribeVisibleLogicalRangeChange(range => {
            if (!syncInProgress && range) {
                syncInProgress = true;
                syncTimeScales(mainTimeScale, macdTimeScale, range);
                setTimeout(() => {
                    syncInProgress = false;
                }, SYNC_DELAY);
            }
          });

          macdTimeScale.subscribeVisibleLogicalRangeChange(range => {
            if (!syncInProgress && range) {
                syncInProgress = true;
                syncTimeScales(macdTimeScale, mainTimeScale, range);
                setTimeout(() => {
                    syncInProgress = false;
                }, SYNC_DELAY);
            }
          });

          // 设置初始范围
          const totalBars = klineData.length;
          const visibleBars = Math.min(100, totalBars);
          const rightPadding = freq === 'dh' ? 12 : 8;  // 分时图不需要右侧空白

          const initialRange = {
            from: Math.max(0, totalBars - visibleBars),
            to: totalBars + rightPadding
          };

          // 同步设置两个图表的范围
          mainTimeScale.setVisibleLogicalRange(initialRange);
          macdTimeScale.setVisibleLogicalRange(initialRange);

          // 根据频率设置图表配置
          if (freq !== 'dh') {  // 分时图
              chartInstance.applyOptions({
                  timeScale: {
                      rightOffset: 0,
                      fixRightEdge: true
                  }
              });
              macdChartInstance.applyOptions({
                  timeScale: {
                      rightOffset: 0,
                      fixRightEdge: true
                  }
              });
          } else {  // K线图
              chartInstance.applyOptions({
                  timeScale: {
                      rightOffset: 12,
                      fixRightEdge: false
                  }
              });
              macdChartInstance.applyOptions({
                  timeScale: {
                      rightOffset: 12,
                      fixRightEdge: false
                  }
              });
          }
        }

        // 更新最新价格信息
        const lastData = klineData[klineData.length - 1];
        priceInfo = {
            open: lastData.o,
            high: lastData.h,
            low: lastData.l,
            close: lastData.c,
            change: lastData.zd || ((lastData.c - lastData.o) / lastData.o * 100),
            volume: lastData.v
        };

        // 修改自动调整大小的逻辑
        if (resizeObserver) resizeObserver.disconnect();
        resizeObserver = new ResizeObserver(entries => {
          const containerHeight = entries[0].contentRect.height;
          const width = entries[0].contentRect.width;
          
          if (chartInstance) {
            chartInstance.resize(
              width,
              Math.floor(containerHeight * 0.7)
            );
          }
          if (macdChartInstance) {
            macdChartInstance.resize(
              width,
              Math.floor(containerHeight * 0.3)
            );
          }
        });
        resizeObserver.observe(chartContainer.parentElement);

      } catch (err) {
        console.error('初始化图表失败:', err);
        error = err instanceof Error ? err.message : '初始化失败';
      } finally {
        loading = false;
      }
    }

    // 修改 cleanup 函数
    async function cleanup() {
      try {
        if (resizeObserver) {
          resizeObserver.disconnect();
          resizeObserver = null;
        }

        // 清理主图表
        if (chartInstance) {
          // 清理均线
          if (series.ma.length) {
            series.ma.forEach(s => {
              if (s) chartInstance.removeSeries(s);
            });
            series.ma = [];
          }

          if (series.volume) {
            chartInstance.removeSeries(series.volume);
            series.volume = null;
          }

          if (series.candlestick) {
            chartInstance.removeSeries(series.candlestick);
            series.candlestick = null;
          }

          chartInstance.remove();
          chartInstance = null;
        }

        // 清理MACD图表
        if (macdChartInstance) {
          // 清理MACD相关的系列
          if (series.macd.histogram) {
            macdChartInstance.removeSeries(series.macd.histogram);
            series.macd.histogram = null;
          }
          if (series.macd.dif) {
            macdChartInstance.removeSeries(series.macd.dif);
            series.macd.dif = null;
          }
          if (series.macd.dea) {
            macdChartInstance.removeSeries(series.macd.dea);
            series.macd.dea = null;
          }

          macdChartInstance.remove();
          macdChartInstance = null;
        }

        // 重置系列对象
        series = {
          candlestick: null,
          volume: null,
          ma: [],
          macd: {
            dif: null,
            dea: null,
            histogram: null
          }
        };

        await new Promise(resolve => requestAnimationFrame(resolve));
      } catch (err) {
        console.error('清理图表失败:', err);
      }
    }

    // 生命周期钩子
    onMount(async () => {
      if (code && chartContainer) {
        await initChart();
      }
    });

    onDestroy(cleanup);

    // 响应式更新：同时监听 code 和 freq 的变化
    let currentCode = code;
    let currentFreq = freq;
    let currentEndTime = endTime;
    $: if ((code !== currentCode || freq !== currentFreq || endTime !== currentEndTime) && chartContainer) {
      currentCode = code;
      currentFreq = freq;
      currentEndTime = endTime;
      (async () => {
        await cleanup();
        await initChart();
      })();
    }
</script>

{#if !code}
<div class="chart-container">
  <div class="chart-overlay">
    <div class="notice">请选择股票</div>
  </div>
</div>
{:else}
<div class="chart-container">
  <!-- 价格信息面板 -->
  <div class="price-panel">
    <div class="price-info">
      <span class="price-label">开</span>
      <span class="price-value" class:up={priceInfo.close >= priceInfo.open} class:down={priceInfo.close < priceInfo.open}>
        {priceInfo.open.toFixed(2)}
      </span>
    </div>
    <div class="price-info">
      <span class="price-label">高</span>
      <span class="price-value up">{priceInfo.high.toFixed(2)}</span>
    </div>
    <div class="price-info">
      <span class="price-label">低</span>
      <span class="price-value down">{priceInfo.low.toFixed(2)}</span>
    </div>
    <div class="price-info">
      <span class="price-label">收</span>
      <span class="price-value" class:up={priceInfo.close >= priceInfo.open} class:down={priceInfo.close < priceInfo.open}>
        {priceInfo.close.toFixed(2)}
      </span>
    </div>
    <div class="price-info">
      <span class="price-label">涨跌</span>
      <span class="price-value" class:up={priceInfo.change >= 0} class:down={priceInfo.change < 0}>
        {priceInfo.change >= 0 ? '+' : ''}{priceInfo.change.toFixed(2)}%
      </span>
    </div>
    <div class="price-info">
      <span class="price-label">量</span>
      <span class="price-value">{(priceInfo.volume / 10000).toFixed(2)}万</span>
    </div>
  </div>

  <div bind:this={chartContainer} class="main-chart"></div>
  <div bind:this={macdContainer} class="macd-chart"></div>

  <!-- 添加工具提示 -->
  {#if tooltip.visible}
    <div 
      class="tooltip"
      style="left: {tooltip.x}px; top: {tooltip.y}px;"
    >
      <div class="tooltip-time">{tooltip.data.time}</div>
      <div class="tooltip-row">
        <span>价格</span>
        <span class:up={tooltip.data.change >= 0} class:down={tooltip.data.change < 0}>
          {tooltip.data.price}
        </span>
      </div>
      <div class="tooltip-row">
        <span>涨跌</span>
        <span class:up={tooltip.data.change >= 0} class:down={tooltip.data.change < 0}>
          {tooltip.data.change}%
        </span>
      </div>
      <div class="tooltip-row">
        <span>成交量</span>
        <span>{tooltip.data.volume}万</span>
      </div>
      <div class="tooltip-row">
        <span>成交额</span>
        <span>{tooltip.data.amount}亿</span>
      </div>
    </div>
  {/if}

  {#if loading}
    <div class="chart-overlay">
      <div class="loading">
        <svg class="spin" viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" stroke-width="2" fill="none">
          <path d="M12 3v3m6.364 2.636l-2.121 2.121M21 12h-3m-2.636 6.364l-2.121-2.121M12 21v-3m-6.364-2.636l2.121-2.121M3 12h3m2.636-6.364l2.121 2.121" />
        </svg>
        <span>加载中...</span>
      </div>
    </div>
  {/if}
  {#if error}
    <div class="chart-overlay">
      <div class="error">
        <div class="error-message">{error}</div>
        <button class="refresh-btn" on:click={() => initChart()}>
          <svg class="refresh-icon" viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none">
            <path d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c-4.97 0-9-4.03-9-9m9 9a9 9 0 009-9" />
          </svg>
          <span>重新加载</span>
        </button>
      </div>
    </div>
  {/if}
</div>
{/if}

<style>
  .chart-container {
    width: 100%;
    height: 100%;
    min-height: 500px;
    display: flex;
    flex-direction: column;
    position: relative;
    overflow: hidden;  /* 防止溢出 */
  }

  .main-chart {
    position: absolute;  /* 使用绝对定位 */
    top: 0;
    left: 0;
    right: 0;
    height: 70%;  /* 使用百分比高度 */
    min-height: 350px;
  }

  .macd-chart {
    position: absolute;  /* 使用绝对定位 */
    left: 0;
    right: 0;
    bottom: 0;
    height: 30%;  /* 使用百分比高度 */
    min-height: 150px;
    border-top: 1px solid rgba(229, 231, 235, 0.5);
  }

  .chart-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.9);
    z-index: 1;
  }

  .loading {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 20px;
    border-radius: 8px;
    color: #6b7280;
    background: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .error {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: 16px 24px;
    border-radius: 8px;
    background: white;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .error-message {
    color: #ef4444;
    font-size: 14px;
    font-weight: 500;
  }

  .refresh-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 16px;
    border: none;
    border-radius: 6px;
    background: #2563eb;
    color: white;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .refresh-btn:hover {
    background: #1d4ed8;
  }

  .notice {
    color: #6b7280;
    background: #f3f4f6;
    padding: 8px 16px;
    border-radius: 4px;
    font-size: 14px;
  }

  .spin {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .price-panel {
    position: absolute;
    top: 8px;
    left: 8px;
    background: rgba(255, 255, 255, 0.85);  /* 增加透明度 */
    padding: 8px 12px;  /* 减小内边距 */
    border-radius: 6px;  /* 减小圆角 */
    border: 1px solid rgba(229, 231, 235, 0.4);  /* 边框更透明 */
    font-size: 11px;  /* 缩小字体 */
    z-index: 2;
    backdrop-filter: blur(8px);
    display: flex;
    gap: 12px;  /* 减小间距 */
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);  /* 减小阴影 */
    animation: fadeIn 0.3s ease;
  }

  .price-info {
    display: flex;
    align-items: center;
    gap: 4px;  /* 减小标签和值之间的间距 */
  }

  .price-label {
    color: #4b5563;  /* 稍微淡化标签颜色 */
    font-weight: 500;
    font-family: system-ui, -apple-system, sans-serif;
    font-size: 11px;  /* 缩小字体 */
  }

  .price-value {
    font-family: 'SF Mono', SFMono-Regular, ui-monospace, Menlo, monospace;
    font-weight: 500;  /* 减小字重 */
    min-width: 48px;  /* 减小最小宽度 */
    text-align: right;
    letter-spacing: -0.2px;
    font-size: 11px;  /* 缩小字体 */
  }

  .up {
    color: #dc2626;
    font-weight: 600;
  }

  .down {
    color: #2563eb;
    font-weight: 600;
  }

  /* 修改悬停效果 */
  .price-info:hover {
    background: rgba(243, 244, 246, 0.3);  /* 减小悬停背景不透明度 */
    border-radius: 4px;  /* 减小圆角 */
    padding: 1px 4px;  /* 减小内边距 */
    margin: -1px -4px;  /* 调整外边距 */
    transition: all 0.15s ease;
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

  .tooltip {
    position: fixed;
    z-index: 9999;
    background: rgba(255, 255, 255, 0.85);  /* 降低不透明度到0.85 */
    border: 1px solid rgba(229, 231, 235, 0.3);  /* 边框更透明 */
    border-radius: 6px;
    padding: 10px;
    font-size: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);  /* 减小阴影 */
    backdrop-filter: blur(8px);
    pointer-events: none;
    min-width: 140px;
  }

  .tooltip-time {
    color: #4b5563;  /* 淡化颜色 */
    font-weight: 500;  /* 减小字重 */
    margin-bottom: 6px;
    font-size: 11px;
    font-family: 'SF Mono', monospace;
  }

  .tooltip-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 3px;  /* 减小行间距 */
    font-size: 11px;  /* 缩小字体 */
  }

  .tooltip-row span:first-child {
    color: #6b7280;
    font-weight: 400;  /* 减小字重 */
  }

  .tooltip-row span:last-child {
    font-family: 'SF Mono', monospace;
    font-weight: 500;
  }
</style>

