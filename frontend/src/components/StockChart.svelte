<script lang="ts">
  import { onMount } from 'svelte';
  import { createChart, AreaSeries, CandlestickSeries, HistogramSeries, ColorType, LineStyle, LineSeries } from 'lightweight-charts';
  import type { IChartApi, DeepPartial, ChartOptions, LineSeriesOptions, SeriesOptionsCommon } from 'lightweight-charts';
  import type { KLineData } from '../types/stock';
  import { themeStore } from '../stores/theme';

  export let code: string;
  export let freq: string = 'dh';  // K线频率，默认日线
  export let defaultVisibleRange: number = 60;  // 默认显示点数
  export let rightOffset: number = 30;  // 默认右侧偏移点数
  
  let chartContainer: HTMLElement;
  let chart: IChartApi;
  let loading = false;
  let error: string | null = null;
  let containerHeight = 0;
  
  // 获取主题状态
  let isDark = false;  // 默认为亮色主题
  let currentTheme = 'light';
  // 增加缓存变量
  let dataMap = new Map<number, KLineData>(); // 时间戳到数据的映射
  let maLegendValues: Record<string, { value: string; color: string }> = {};
  let macdLegendValues: Record<string, { value: string; color: string }> = {};
  let kdjLegendValues: Record<string, { value: string; color: string }> = {};


  themeStore.subscribe(theme => {
    isDark = theme === 'dark';
    currentTheme = theme;
    requestAnimationFrame(() => {
      if (chart) {
        // 清理旧图表
        chart.remove();
        // 重新初始化图表
        initChart();
      }
    });
  });

  // 更新图表主题
  function updateChartTheme() {
    if (!chart) return;
    
    chart.applyOptions({
      layout: {
        textColor: isDark ? '#d1d5db' : '#1f2937',
        background: { type: ColorType.Solid, color: 'transparent' },
      },
      grid: {
        vertLines: {
          color: isDark ? 'rgba(197, 203, 206, 0.1)' : 'rgba(70, 130, 180, 0.1)',
        },
        horzLines: {
          color: isDark ? 'rgba(197, 203, 206, 0.1)' : 'rgba(70, 130, 180, 0.1)',
        },
      },
    });
  }
  
  // MA 配置 - 简化为固定值
  let maType: string = 'SMA';  // 默认使用 SMA
  let maPeriods: number[] = [6, 18, 98];  // 默认周期
  let maColors: string[] = ['#60a5fa', '#f472b6', '#a78bfa'];  // 修改为蓝色、粉色、紫色
  
  // MA 显示控制 - 默认全部显示
  let maVisibility: { [key: string]: boolean } = {
    'SMA6': true,
    'SMA18': true,
    'SMA98': true,
  };

  // 格式化数字
  function formatNumber(num: number | undefined): string {
    if (num === undefined || isNaN(num)) return '-';
    return num.toFixed(2);
  }

  // 获取K线数据
  async function fetchKLineData() {
    try {
      const data = await (window as any).go.main.App.GetKLineData(code, freq);
      if (!data || data.length === 0) throw new Error('没有数据');
      return data;
    } catch (err) {
      console.error('获取K线数据失败:', err);
      throw err;
    }
  }

  // 计算 MACD 指标
  async function calculateMACD(prices: number[]) {
    try {
      return await (window as any).go.main.App.CalculateMACD(prices);
    } catch (err) {
      console.error('计算MACD失败:', err);
      return null;
    }
  }

  // 计算 KDJ 指标
  async function calculateKDJ(prices: number[]) {
    try {
      if (!prices || prices.length === 0) {
        throw new Error('价格数据为空');
      }
      
      const result = await (window as any).go.main.App.CalculateKDJ(prices);
      
      if (result.Error) {
        throw new Error(`KDJ计算错误: ${result.Error}`);
      }
      
      if (!result.K || !result.D || !result.J) {
        throw new Error('KDJ数据无效');
      }
      
      if (!Array.isArray(result.K) || !Array.isArray(result.D) || !Array.isArray(result.J)) {
        throw new Error('KDJ数据格式错误');
      }
      
      if (result.K.length === 0 || result.D.length === 0 || result.J.length === 0) {
        throw new Error('KDJ数据为空');
      }
      
      return [result.K, result.D, result.J];
    } catch (err) {
      console.error('计算KDJ失败:', err);
      return null;
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
        maLines.push({
          period,
          data: maData
        });
      } catch (err) {
        console.error(`计算${maType}${period}失败:`, err);
      }
    }
    return maLines;
  }

  // 添加时间格式转换函数
  function convertTime(timeStr: string, freq: string) {
    try {
      // 将时间字符串转换为时间戳（秒级）
      const timestamp = Math.floor(new Date(timeStr.replace(/-/g, '/')).getTime() / 1000);
      return timestamp;
    } catch (err) {
      console.error('时间转换失败:', timeStr, err);
      return 0;
    }
  }

  // 添加时间格式转换函数
  function formatDateTime(timestamp: any): string {
    if (typeof timestamp === 'string') {
      return timestamp;
    }
    if (typeof timestamp === 'object' && timestamp.year && timestamp.month && timestamp.day) {
      // 处理 BusinessDay 类型
      const month = String(timestamp.month).padStart(2, '0');
      const day = String(timestamp.day).padStart(2, '0');
      return `${timestamp.year}-${month}-${day}`;
    }
    if (typeof timestamp === 'number') {
      const date = new Date(timestamp * 1000);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    }
    return 'Invalid Date';
  }

  const chartOptions: DeepPartial<ChartOptions> = {
    layout: {
      textColor: isDark ? '#d1d5db' : '#1f2937',
      background: { type: ColorType.Solid, color: 'transparent' },
      panes: {
        separatorColor: '#f22c3d',
        separatorHoverColor: 'rgba(255, 0, 0, 0.1)',
      },
    },
    grid: {
      vertLines: {
        visible: true,
        color: isDark ? 'rgba(197, 203, 206, 0.1)' : 'rgba(70, 130, 180, 0.1)',
      },
      horzLines: {
        visible: true,
        color: isDark ? 'rgba(197, 203, 206, 0.1)' : 'rgba(70, 130, 180, 0.1)',
      },
    },
    timeScale: {
      rightOffset: rightOffset,
      barSpacing: 12,
      borderVisible: true,  // 显示时间轴边框
    },
  };

  async function initChart() {
    if (!chartContainer || !code) return;
    
    try {
      loading = true;
      error = null;

      // 获取父容器的高度
      containerHeight = chartContainer.parentElement?.clientHeight ?? 0;

      // 创建工具提示元素
      const toolTip = document.createElement('div');
      toolTip.setAttribute('style', `
        position: absolute;
        display: none;
        padding: 12px;
        box-sizing: border-box;
        font-size: 12px;
        text-align: left;
        z-index: 1000;
        pointer-events: none;
        border-radius: 6px;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
        background: ${isDark ? 'rgba(0, 0, 0, 0.85)' : 'rgba(255, 255, 255, 0.95)'};
        color: ${isDark ? '#fff' : '#1f2937'};
        border: 1px solid ${isDark ? 'rgba(255, 255, 255, 0.2)' : 'rgba(0, 0, 0, 0.1)'};
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        backdrop-filter: blur(4px);
        -webkit-backdrop-filter: blur(4px);
      `);
      chartContainer.appendChild(toolTip);

      // 创建固定图例
      const maLegendFixed = document.createElement('div');
      maLegendFixed.setAttribute('style', `
        position: absolute;
        top: 4px;
        left: 4px;
        z-index: 3;
        padding: 4px 8px;
        background: ${isDark ? 'rgba(0, 0, 0, 0.5)' : 'rgba(255, 255, 255, 0.9)'};
        color: ${isDark ? '#fff' : '#1f2937'};
        border-radius: 4px;
        font-size: 12px;
        pointer-events: none;
        display: flex;
        flex-direction: row;
        gap: 12px;
        border: 1px solid ${isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)'};
        box-shadow: 0 2px 4px ${isDark ? 'rgba(0, 0, 0, 0.2)' : 'rgba(0, 0, 0, 0.05)'};
      `);

      const macdLegendFixed = document.createElement('div');
      macdLegendFixed.setAttribute('style', `
        position: absolute;
        top: 4px;
        left: 4px;
        z-index: 3;
        padding: 4px 8px;
        background: ${isDark ? 'rgba(0, 0, 0, 0.5)' : 'rgba(255, 255, 255, 0.9)'};
        color: ${isDark ? '#fff' : '#1f2937'};
        border-radius: 4px;
        font-size: 12px;
        pointer-events: none;
        display: flex;
        flex-direction: row;
        gap: 12px;
        border: 1px solid ${isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)'};
        box-shadow: 0 2px 4px ${isDark ? 'rgba(0, 0, 0, 0.2)' : 'rgba(0, 0, 0, 0.05)'};
      `);

      const kdjLegendFixed = document.createElement('div');
      kdjLegendFixed.setAttribute('style', `
        position: absolute;
        top: 4px;
        left: 4px;
        z-index: 3;
        padding: 4px 8px;
        background: ${isDark ? 'rgba(0, 0, 0, 0.5)' : 'rgba(255, 255, 255, 0.9)'};
        color: ${isDark ? '#fff' : '#1f2937'};
        border-radius: 4px;
        font-size: 12px;
        pointer-events: none;
        display: flex;
        flex-direction: row;
        gap: 12px;
        border: 1px solid ${isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)'};
        box-shadow: 0 2px 4px ${isDark ? 'rgba(0, 0, 0, 0.2)' : 'rgba(0, 0, 0, 0.05)'};
      `);

      const klineData = await fetchKLineData();
      const prices = klineData.map((item: KLineData) => item.c);
      // 创建数据映射
      dataMap.clear();
      klineData.forEach((item: KLineData) => {
        dataMap.set(convertTime(item.d, freq), item);
      });
      // 计算技术指标
      const [macdResult, kdjResult] = await Promise.all([
        calculateMACD(prices),
        calculateKDJ(prices)
      ]);
      
      chart = createChart(chartContainer, {
        ...chartOptions,
        width: chartContainer.clientWidth,
        height: containerHeight
      });

      // 等待一下让图表完全初始化
      await new Promise(resolve => setTimeout(resolve, 0));

      // 设置各个窗格的高度比例 70:15:15
      const panes = chart.panes();
      if (panes.length >= 3) {
        const klineHeight = Math.floor(containerHeight * 0.70);
        const volumeHeight = Math.floor(containerHeight * 0.15);
        const macdHeight = Math.floor(containerHeight * 0.15);

        panes[2].setHeight(macdHeight);  // MACD 15%
        panes[1].setHeight(volumeHeight);  // 成交量 15%
        panes[0].setHeight(klineHeight);  // K线图 70%

        // 等待一下让图表完全渲染
        await new Promise(resolve => setTimeout(resolve, 100));

        // 使用更精确的选择器定位窗格
        const chartElement = document.querySelector('.stock-chart-container-content .tv-lightweight-charts');
        
        if (chartElement) {
          const klinePane = chartElement.querySelector('table tr:nth-child(1) td:nth-child(2)') as HTMLElement;
          const macdPane = chartElement.querySelector('table tr:nth-child(5) td:nth-child(2)') as HTMLElement;
             
          if (klinePane && macdPane) {
            maLegendFixed.style.position = 'absolute';
            maLegendFixed.style.top = '4px';
            maLegendFixed.style.left = '4px';
            
            macdLegendFixed.style.position = 'absolute';
            macdLegendFixed.style.top = '4px';
            macdLegendFixed.style.left = '4px';
            
            kdjLegendFixed.style.position = 'absolute';
            kdjLegendFixed.style.top = '4px';
            kdjLegendFixed.style.left = '4px';

            // 将图例添加到对应的窗格
            klinePane.appendChild(maLegendFixed);
            const macdPaneParent = macdPane.closest('tr');
            if (macdPaneParent) {
              const prevRow = macdPaneParent.previousElementSibling;
              if (prevRow) {
                const prevPane = prevRow.querySelector('td:nth-child(2)') as HTMLElement;
                if (prevPane) {
                  prevPane.appendChild(macdLegendFixed);
                }
              }
            }
            macdPane.appendChild(kdjLegendFixed);

            // 确保图例容器的定位是相对于窗格
            klinePane.style.position = 'relative';
            macdPane.style.position = 'relative';
            
        
          } 
        } else {
          console.log('未找到图表容器元素');
          // 输出所有可能的容器
          console.log('所有.tv-lightweight-charts元素:', document.querySelectorAll('.tv-lightweight-charts'));
          console.log('所有.stock-chart-container-content元素:', document.querySelectorAll('.stock-chart-container-content'));
        }
      }

      // 先添加所有的图表
      const candlestickSeries = chart.addSeries(CandlestickSeries, {
        upColor: 'transparent',  // 涨为空心
        downColor: '#2563eb',   // 跌为更深的蓝色
        borderUpColor: '#dc2626',   // 涨边框为更深的红色
        borderDownColor: '#2563eb',  // 跌边框为蓝色
        wickUpColor: '#dc2626',      // 上影线为更深的红色
        wickDownColor: '#2563eb',    // 下影线为蓝色
        borderVisible: true,         // 显示边框
      });

      // 设置K线的位置和样式
      candlestickSeries.applyOptions({
        upColor: 'transparent',  // 涨为空心
        downColor: '#2563eb',   // 跌为更深的蓝色
        borderUpColor: '#dc2626',   // 涨边框为更深的红色
        borderDownColor: '#2563eb',  // 跌边框为蓝色
        wickUpColor: '#dc2626',      // 上影线为更深的红色
        wickDownColor: '#2563eb',    // 下影线为蓝色
        borderVisible: true,         // 显示边框
      });

      candlestickSeries.priceScale().applyOptions({
        scaleMargins: {
          top: 0.1,    // 距离顶部10%
          bottom: 0.2, // 距离底部40%
        },
        borderColor: isDark ? 'rgba(197, 203, 206, 0.3)' : 'rgba(0, 0, 0, 0.1)',
      });

      // 将成交量作为叠加层添加到主图表
      const volumeSeries = chart.addSeries(HistogramSeries, {
        color: isDark ? 'rgba(75, 85, 99, 0.3)' : 'rgba(107, 114, 128, 0.3)', // 降低成交量透明度
        priceFormat: {
          type: 'volume',
        },
        priceScaleId: '', // 设置为空字符串，作为叠加层
      });

      // 设置成交量的位置
      volumeSeries.priceScale().applyOptions({
        scaleMargins: {
          top: 0.8,    // 距离顶部70%
          bottom: 0,   // 紧贴底部
        },
      });

      // 等待所有数据设置完成
      await new Promise(resolve => setTimeout(resolve, 500));

      // 尝试多次查找图表容器，直到找到为止
      let attempts = 0;
      const maxAttempts = 10;
      const findChartContainer = async () => {
        while (attempts < maxAttempts) {
          const chartElement = document.querySelector('.stock-chart-container-content');
          
          if (chartElement) {
            const tvCharts = chartElement.querySelectorAll('.tv-lightweight-charts');
          
            
            if (tvCharts.length > 0) {
              const lastChart = tvCharts[tvCharts.length - 1];
  
              
              // 查找所有的 tr 元素
              const rows = lastChart.querySelectorAll('tr');
             
              
              if (rows.length >= 5) {
                const klinePane = rows[0].querySelector('td:nth-child(2)') as HTMLElement;
                const macdPane = rows[4].querySelector('td:nth-child(2)') as HTMLElement;
                
              
                
                if (klinePane && macdPane) {
                
                  
                  // 确保图例容器的定位是相对于窗格
                  klinePane.style.position = 'relative';
                  macdPane.style.position = 'relative';
                  
                  // 将图例添加到对应的窗格
                  klinePane.appendChild(maLegendFixed);
                  
                  // 修改MACD图例的位置到第3行
                  const macdRow = rows[2].querySelector('td:nth-child(2)') as HTMLElement;
                  if (macdRow) {
                    macdRow.style.position = 'relative';
                    macdRow.appendChild(macdLegendFixed);
                  }
                  
                  macdPane.appendChild(kdjLegendFixed);
                  
                 
                  return true;
                }
              }
            }
          }
          
          attempts++;
          await new Promise(resolve => setTimeout(resolve, 100));
        }
       
        return false;
      };

      findChartContainer();

      if (macdResult) {
        // MACD 柱状图 - 最重要
        const macdHistogram = chart.addSeries(HistogramSeries, {
          color: isDark ? 'rgba(99, 102, 241, 0.8)' : 'rgba(79, 70, 229, 0.8)',
          priceFormat: {
            type: 'price',
            precision: 2,
          },
          priceScaleId: 'macd'
        },1);

        // DIF 线 - 次重要
        const difLine = chart.addSeries(LineSeries, {
          color: '#60a5fa',
          lineWidth: 1,
          lineStyle: LineStyle.Solid,
          title: 'DIF',
          priceScaleId: 'macd',
          crosshairMarkerVisible: true,
          priceLineVisible: false,  // 禁用价格线
        },1);

        // DEA 线 - 第三重要
        const deaLine = chart.addSeries(LineSeries, {
          color: '#f472b6',
          lineWidth: 1,
          lineStyle: LineStyle.Dashed,
          title: 'DEA',
          priceScaleId: 'macd',
          crosshairMarkerVisible: true,
          priceLineVisible: false,  // 禁用价格线
        },1);

        // 设置 MACD 数据
        const macdData = macdResult.MACD.map((value: number, i: number) => ({
          time: convertTime(klineData[i].d, freq),
          value: value ,
          color: value >= 0 ? '#dc2626' : '#2563eb'  // 保持红蓝对比
        }));

        const difData = macdResult.DIF.map((value: number, i: number) => ({
          time: convertTime(klineData[i].d, freq),
          value: value
        }));

        const deaData = macdResult.DEA.map((value: number, i: number) => ({
          time: convertTime(klineData[i].d, freq),
          value: value
        }));

        macdHistogram.setData(macdData);
        difLine.setData(difData);
        deaLine.setData(deaData);
      }

      if (kdjResult) {
        // K 线 - 最重要
        const kLine = chart.addSeries(LineSeries, {
          color: '#f43f5e',
          lineWidth: 1,
          lineStyle: LineStyle.Solid,
          title: 'K',
          priceScaleId: 'kdj',
          crosshairMarkerVisible: true,
          priceLineVisible: false,  // 禁用价格线
        }, 2);

        // D 线 - 次重要
        const dLine = chart.addSeries(LineSeries, {
          color: '#3b82f6',
          lineWidth: 1,
          lineStyle: LineStyle.Dashed,
          title: 'D',
          priceScaleId: 'kdj',
          crosshairMarkerVisible: true,
          priceLineVisible: false,  // 禁用价格线
        }, 2);

        // J 线 - 第三重要
        const jLine = chart.addSeries(LineSeries, {
          color: '#a855f7',
          lineWidth: 1,
          lineStyle: LineStyle.LargeDashed,
          title: 'J',
          priceScaleId: 'kdj',
          crosshairMarkerVisible: true,
          priceLineVisible: false,  // 禁用价格线
        }, 2);

        // 设置 KDJ 数据
        const kData = kdjResult[0].map((value: number, i: number) => ({
          time: convertTime(klineData[i].d, freq),
          value: value
        }));

        const dData = kdjResult[1].map((value: number, i: number) => ({
          time: convertTime(klineData[i].d, freq),
          value: value
        }));

        const jData = kdjResult[2].map((value: number, i: number) => ({
          time: convertTime(klineData[i].d, freq),
          value: value
        }));

        kLine.setData(kData);
        dLine.setData(dData);
        jLine.setData(jData);
      }

      // 设置价格轴
      chart.priceScale('macd').applyOptions({
        scaleMargins: {
          top: 0.2,
          bottom: 0.2,
        },
        borderColor: 'rgba(197, 203, 206, 0.3)',
        borderVisible: true,  // 显示价格轴边框
      });

      // 适配暗色模式
      chart.applyOptions({
        layout: {
          textColor: isDark ? '#d1d5db' : '#1f2937',
        },
        grid: {
          vertLines: {
            visible: true,
            color: isDark ? 'rgba(197, 203, 206, 0.1)' : 'rgba(70, 130, 180, 0.1)',
          },
          horzLines: {
            visible: true,
            color: isDark ? 'rgba(197, 203, 206, 0.1)' : 'rgba(70, 130, 180, 0.1)',
          },
        },
      });

      // 转换K线数据格式
      const candleData = klineData.map((item: KLineData) => ({
        time: convertTime(item.d, freq),
        open: item.o,
        high: item.h,
        low: item.l,
        close: item.c
      }));

      // 转换成交量数据格式
      const volumeData = klineData.map((item: KLineData) => ({
        time: convertTime(item.d, freq),
        value: item.v || 0,
        turnover: item.hs || 0,  // 添加换手率数据
        color: item.c >= item.o ? '#dc2626' : '#2563eb' // 使用更深的红蓝色
      }));

      // 添加均线
      const maLines = await calculateMAs(prices);
      maLines.forEach((maLine, index) => {
        const opacity = 1 - (index * 0.12);
        const lineStyle = index === 0 ? LineStyle.Solid : 
                         index === 1 ? LineStyle.Dashed : 
                         LineStyle.LargeDashed;
        
        const lineOptions: DeepPartial<LineSeriesOptions & SeriesOptionsCommon> = {
          color: isDark 
            ? `rgba(${index === 0 ? '96, 165, 250' : index === 1 ? '244, 114, 182' : '167, 139, 250'}, ${opacity})` 
            : `rgba(${index === 0 ? '96, 165, 250' : index === 1 ? '244, 114, 182' : '167, 139, 250'}, ${opacity})`,
          lineWidth: index === 0 ? 2 : 1,
          title: `${maType}${maLine.period}`,
          priceLineVisible: false,  // 禁用价格线
          lastValueVisible: true,
          crosshairMarkerVisible: true,
          lastPriceAnimation: 0,  // 禁用最后价格动画
          lineStyle: lineStyle,
        };

        const lineSeries = chart.addSeries(LineSeries, lineOptions);
        const maData = maLine.data.map((value: number, i: number) => ({
          time: convertTime(klineData[i].d, freq),
          value: value !== 0 ? value : undefined
        })).filter((data: { value: undefined; }) => data.value !== undefined);

        lineSeries.setData(maData);
      });

      candlestickSeries.setData(candleData);
      volumeSeries.setData(volumeData);
      
      // 设置默认显示范围
      const timeScale = chart.timeScale();
      timeScale.applyOptions({
        rightOffset: rightOffset,
        barSpacing: 12,
      });

      const dataLength = candleData.length;
      if (dataLength > defaultVisibleRange) {
        const from = dataLength - defaultVisibleRange;
        const to = dataLength - 1;
        timeScale.setVisibleLogicalRange({
          from,
          to
        });
      } else {
        timeScale.fitContent();
      }

      // 添加窗口大小变化监听
      const resizeObserver = new ResizeObserver(entries => {
        if (chart && entries[0]) {
          const { width, height } = entries[0].contentRect;
          chart.resize(width, height);
        }
      });

      resizeObserver.observe(chartContainer);

      // 添加更新图例的辅助函数
      const updateLegend = (
        container: HTMLElement, 
        newValues: Record<string, { value: string; color: string }>, 
        cache: Record<string, { value: string; color: string }>
      ) => {
        let shouldUpdate = false;
        
        // 检查是否有变化
        for (const [key, value] of Object.entries(newValues)) {
          if (!cache[key] || cache[key].value !== value.value || cache[key].color !== value.color) {
            shouldUpdate = true;
            break;
          }
        }
        
        if (!shouldUpdate) return;
        
        // 批量更新 DOM
        container.innerHTML = Object.entries(newValues)
          .map(([title, data]) => `
            <div>
              <span style="color: ${isDark ? '#ffffff' : '#1f2937'}">${title}</span>
              <span style="color: ${data.color}">${data.value}</span>
            </div>
          `)
          .join('');
        
        // 更新缓存
        Object.assign(cache, newValues);
      };

      // 添加节流处理
      let lastUpdate = 0;
      const throttledCrosshairMove = (param: any) => {
        if (
          param.point === undefined ||
          !param.time ||
          param.point.x < 0 ||
          param.point.x > chartContainer.clientWidth ||
          param.point.y < 0 ||
          param.point.y > chartContainer.clientHeight
        ) {
          toolTip.style.display = 'none';
          return;
        }

        const now = Date.now();
        if (now - lastUpdate < 50) return; // 50ms 节流
        lastUpdate = now;

        // 获取当前K线数据
        const currentKLineData = dataMap.get(param.time as number);
        const candleDataPoint = param.seriesData.get(candlestickSeries) as any;
        const volumeDataPoint = param.seriesData.get(volumeSeries) as any;

        if (candleDataPoint && typeof candleDataPoint === 'object') {
          const { open, high, low, close } = candleDataPoint;
          const volume = volumeDataPoint?.value || 0;
          const change = ((close - open) / open * 100).toFixed(2);
          const changeColor = close >= open ? '#dc2626' : '#2563eb';
          const changeSign = close >= open ? '+' : '';

          // 更新 MA 图例
          if (open !== undefined && high !== undefined && low !== undefined && close !== undefined) {
            const newMaLegend = {
              '开盘': { value: formatNumber(open), color: changeColor },
              '最高': { value: formatNumber(high), color: changeColor },
              '最低': { value: formatNumber(low), color: changeColor },
              '收盘': { value: formatNumber(close), color: changeColor },
              '涨幅': { value: `${changeSign}${change}%`, color: changeColor },
              '成交量': { value: formatNumber(volume), color: changeColor },
              '换手率': { value: `${currentKLineData?.hs || 0}%`, color: changeColor }
            };
            updateLegend(maLegendFixed, newMaLegend, maLegendValues);
          }

          // 更新MACD图例值
          let macdHtml = '';
          const macdValues = new Map();
          param.seriesData.forEach((value: any, series: any) => {
            const options = series.options();
            if (options.priceScaleId === 'macd') {
              const title = options.title || 'MACD';
              const seriesColor = (options as any).color || '#ffffff';
              macdValues.set(title, { value: value.value, color: seriesColor });
            }
          });
          
          // 按固定顺序显示MACD值
          const macdOrder = ['DIF', 'DEA', 'MACD'];
          const newMacdLegend: Record<string, { value: string; color: string }> = {};
          macdOrder.forEach(title => {
            const data = macdValues.get(title);
            if (data) {
              newMacdLegend[title] = {
                value: formatNumber(data.value),
                color: data.color
              };
            }
          });
          updateLegend(macdLegendFixed, newMacdLegend, macdLegendValues);

          // 更新KDJ图例值
          let kdjHtml = '';
          const kdjValues = new Map();
          param.seriesData.forEach((value: any, series: any) => {
            const options = series.options();
            if (options.priceScaleId === 'kdj') {
              const title = options.title;
              const seriesColor = (options as any).color || '#ffffff';
              kdjValues.set(title, { value: value.value, color: seriesColor });
            }
          });
          
          // 按固定顺序显示KDJ值
          const kdjOrder = ['K', 'D', 'J'];
          const newKdjLegend: Record<string, { value: string; color: string }> = {};
          kdjOrder.forEach(title => {
            const data = kdjValues.get(title);
            if (data) {
              newKdjLegend[title] = {
                value: formatNumber(data.value),
                color: data.color
              };
            }
          });
          updateLegend(kdjLegendFixed, newKdjLegend, kdjLegendValues);

          // 更新悬浮提示
          let tooltipHtml = '';
          if (candleDataPoint && volumeDataPoint) {
            const price = candleDataPoint.close;
            const volume = volumeDataPoint.value;
            const change = ((candleDataPoint.close - candleDataPoint.open) / candleDataPoint.open * 100).toFixed(2);
            const changeSign = candleDataPoint.close >= candleDataPoint.open ? '+' : '';
            const changeColor = candleDataPoint.close >= candleDataPoint.open ? '#dc2626' : '#2563eb';
            
            // 根据时间找到对应的K线数据
            const currentKLineData = klineData.find((item: KLineData) => {
              const itemTime = convertTime(item.d, freq);
              return itemTime === param.time;
            });
            const turnover = currentKLineData?.hs || 0;

            tooltipHtml = `
              <div style="margin-bottom: 12px;">
                <div style="color: ${isDark ? '#9ca3af' : '#4b5563'}; margin-bottom: 8px; font-weight: 500;">${formatDateTime(param.time)}</div>
                <div style="display: grid; grid-template-columns: auto 1fr; gap: 8px 16px; min-width: 160px;">
                  <span style="color: ${isDark ? '#9ca3af' : '#4b5563'}; font-weight: 500;">价格</span>
                  <span style="color: ${changeColor}; text-align: right;">${formatNumber(price)}</span>
                  <span style="color: ${isDark ? '#9ca3af' : '#4b5563'}; font-weight: 500;">涨幅</span>
                  <span style="color: ${changeColor}; text-align: right;">${changeSign}${change}%</span>
                  <span style="color: ${isDark ? '#9ca3af' : '#4b5563'}; font-weight: 500;">成交量</span>
                  <span style="text-align: right;">${formatNumber(volume)}</span>
                  <span style="color: ${isDark ? '#9ca3af' : '#4b5563'}; font-weight: 500;">换手率</span>
                  <span style="text-align: right;">${turnover.toFixed(2)}%</span>
                </div>
              </div>
            `;
          }

          if (tooltipHtml) {
            toolTip.innerHTML = tooltipHtml;
            toolTip.style.display = 'block';

            const toolTipWidth = toolTip.clientWidth;
            const toolTipHeight = toolTip.clientHeight;
            const toolTipMargin = 15;

            let left = param.point.x + toolTipMargin;
            if (left > chartContainer.clientWidth - toolTipWidth) {
              left = param.point.x - toolTipMargin - toolTipWidth;
            }

            let top = param.point.y + toolTipMargin;
            if (top > chartContainer.clientHeight - toolTipHeight) {
              top = param.point.y - toolTipHeight - toolTipMargin;
            }

            toolTip.style.left = left + 'px';
            toolTip.style.top = top + 'px';
          } else {
            toolTip.style.display = 'none';
          }
        }
      };

      chart.subscribeCrosshairMove(throttledCrosshairMove);

      // 更新清理函数
      return () => {
        resizeObserver.disconnect();
        if (chart) {
          chart.remove();
        }
        if (toolTip && toolTip.parentNode) {
          toolTip.parentNode.removeChild(toolTip);
        }
        if (maLegendFixed && maLegendFixed.parentNode) {
          maLegendFixed.parentNode.removeChild(maLegendFixed);
        }
        if (macdLegendFixed && macdLegendFixed.parentNode) {
          macdLegendFixed.parentNode.removeChild(macdLegendFixed);
        }
        if (kdjLegendFixed && kdjLegendFixed.parentNode) {
          kdjLegendFixed.parentNode.removeChild(kdjLegendFixed);
        }
        // 清理缓存
        dataMap.clear();
        maLegendValues = {};
        macdLegendValues = {};
        kdjLegendValues = {};
      };
    } catch (err) {
      console.error('初始化图表失败:', err);
      error = err instanceof Error ? err.message : '初始化失败';
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    if (code) {
      initChart();
    }
  });

  // 响应式更新：监听 code 和 freq 的变化
  let currentCode = code;
  let currentFreq = freq;
  $: if ((code !== currentCode || freq !== currentFreq) && chartContainer) {
    currentCode = code;
    currentFreq = freq;
    if (chart) {
      chart.remove();
    }
    initChart();
  }
</script>

{#if !code}
<div class="chart-container">
  <div class="chart-overlay">
    <div class="notice">请选择股票</div>
  </div>
</div>
{:else}
<div bind:this={chartContainer} class="chart-container">
  {#if loading}
    <div class="chart-overlay">
      <div class="loading">加载中...</div>
    </div>
  {/if}
  {#if error}
    <div class="chart-overlay">
      <div class="error">
        <div class="error-message">{error}</div>
        <button class="refresh-btn" on:click={() => initChart()}>
          重新加载
        </button>
      </div>
    </div>
  {/if}
</div>
{/if}

<style>
  .chart-container {
    width: 100%;
    height: 100%;  /* 改为100%填充父级高度 */
    min-height: 300px;  /* 添加最小高度防止太小 */
    position: relative;
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
    background: transparent;
  }

  .loading {
    padding: 12px 20px;
    border-radius: 8px;
    background: var(--surface);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .error {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: 16px 24px;
    border-radius: 8px;
    background: var(--surface);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .error-message {
    color: var(--error-500);
  }

  .refresh-btn {
    padding: 8px 16px;
    border: none;
    border-radius: 6px;
    background: var(--primary-500);
    color: white;
    cursor: pointer;
  }

  .refresh-btn:hover {
    background: var(--primary-600);
  }

  .notice {
    color: var(--text-secondary);
    background: var(--surface-variant);
    padding: 8px 16px;
    border-radius: 4px;
  }

  .legend-container {
    position: absolute;
    z-index: 2;
    top: 12px;
    left: 12px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    font-size: 12px;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
    padding: 8px;
    border-radius: 4px;
    background: rgba(0, 0, 0, 0.7);
    color: #fff;
  }

  .legend {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .legend-line {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .legend-title {
    font-weight: 500;
    min-width: 50px;
  }

  .legend-value {
    opacity: 0.9;
  }

  :global(.legend-container) {
    pointer-events: none;
  }
</style>