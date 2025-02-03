<script lang="ts">
  import type { AnalysisProgress } from '../../types/analysis'
  import Modal from '../Modal.svelte'
  import { fade, slide } from 'svelte/transition'

  import { RefreshCwIcon, TrendingUpIcon, TrendingDownIcon, MinusIcon, AlertTriangleIcon, BarChart3Icon, LineChartIcon, ActivityIcon, InfoIcon } from 'lucide-svelte'
  import { EventsOn } from '../../../wailsjs/runtime/runtime'
  import { onMount, onDestroy } from 'svelte'
  import AIAnalysisDetailModal from './AIAnalysisDetailModal.svelte'
  import { AnalyzeStock } from '../../../wailsjs/go/main/App'

  export let code: string = ''

  let analysisData: any = null
  let isAnalyzing: boolean = false
  let analysisError: string = ''
  let showDetailModal = false
  let quickAnalysis: StockAnalysis | null = null

  // 添加进度状态
  let progress: AnalysisProgress = {
    step: 0,
    total: 0,
    message: '',
    phase: ''
  }

  let activeTab = 'fundamental' // 'fundamental' | 'technical'

  // 添加分析阶段描述
  const phaseDescriptions = {
    data: {
      title: '数据收集',
      steps: [
        '获取公司基本信息',
        '获取财务数据',
        '获取市场行情',
        '获取资金流向',
        '获取股东信息',
        '获取行业数据',
        '获取技术指标',
        '获取风险信息',
        '整合分析数据',
        '准备AI分析'
      ]
    },
    analysis: {
      title: 'AI分析',
      steps: [
        '快速策略分析',
        '深度价值分析',
        '技术形态分析',
        '风险评估',
        '综合建议生成'
      ]
    }
  };

  // 添加指标说明
  const metricDescriptions = {
    performanceScore: "业绩评分基于收入增长率(30%)、净利率(40%)、ROE(30%)加权计算，分数越高表示业绩越好",
    revenueGrowth: "营收增长率反映公司业务扩张速度，↑表示加速增长，↓表示增速放缓，!表示异常波动",
    roe: "净资产收益率(ROE)反映公司盈利能力，15%以上为优秀，8-15%为良好，低于8%需关注",
    pe: "市盈率(PE)反映估值水平，需要结合行业平均和公司成长性综合判断",
    pb: "市净率(PB)反映资产溢价程度，需要结合行业特征和ROE水平判断",
    risk: "风险评分基于财务风险(40%)、经营风险(30%)、市场风险(30%)计算，分数越高风险越大",
    targetPrice: "基于估值和成长性分析的目标价格，仅供参考，不构成投资建议",
    stopLoss: "建议的风险控制价格，跌破此价格应当考虑止损出场",
    mainForce: "主力资金动向反映大资金流向，方向和强度越大影响越明显",
    netMargin: "净利率反映公司盈利能力，数值越高表示盈利能力越强"
  };

  // 获取当前阶段的进度描述
  function getProgressDescription(phase: string, step: number): string {
    const phaseInfo = phaseDescriptions[phase as keyof typeof phaseDescriptions];
    if (!phaseInfo) return '';
    
    const currentStep = phaseInfo.steps[step - 1] || phaseInfo.steps[0];
    return currentStep;
  }

  // 获取总体进度
  function getTotalProgress(): number {
    if (!progress.phase) return 0;
    if (progress.phase === 'data') {
      return (progress.step / progress.total) * 50;
    } else if (progress.phase === 'analysis') {
      return 50 + (quickAnalysis ? 75 : (progress.step / progress.total) * 50);
    }
    return 100;
  }

  // 获取进度状态类
  function getProgressStatusClass(): string {
    if (quickAnalysis) return 'progress-quick';
    if (progress.phase === 'analysis') return 'progress-analysis';
    return 'progress-data';
  }

  // 修改 handleRefresh 函数
  async function handleRefresh(event?: MouseEvent | number) {
    const retryCount = typeof event === 'number' ? event : 0;
    console.log('开始分析, code:', code); // 添加日志
    if (isAnalyzing && !analysisError) return;
    
    try {
      isAnalyzing = true;
      analysisError = '';
      analysisData = null;
      progress = {
        step: 0,
        total: 0,
        message: '',
        phase: 'data'
      };
      
      console.log('调用 AnalyzeStock...'); // 添加日志
      const result = await AnalyzeStock(code);
      console.log('AnalyzeStock 返回结果:', result); // 添加日志
      
      if (result) {
        analysisData = result;
        console.log('分析数据已更新:', analysisData); // 添加日志
      }
    } catch (error) {
      console.error('分析请求失败:', error);
      analysisError = error instanceof Error ? error.message : String(error);
    } finally {
      if (analysisError) {
        isAnalyzing = false;
      }
    }
  }

  // 监听代码变化自动开始分析
  $: if (code) {
    handleRefresh();
  }

  onMount(() => {
    console.log('组件挂载，设置事件监听...');
    
    // 修改事件名称
    EventsOn("analysis_progress", (data: AnalysisProgress) => {
        console.log('收到进度更新:', data);
        if (!analysisError) {
            progress = data;
        }
    });

    EventsOn("analysis_error", (error: string) => {
        console.log('收到错误:', error);
        analysisError = error;
        isAnalyzing = false;
    });

    EventsOn("analysis_update", (response: AnalysisResponse) => {
        console.log('收到分析更新:', response);
        if (response.data) {
            analysisData = response.data;
        }
    });

    if (code) {
        console.log('初始化时存在code，开始分析:', code);
        handleRefresh();
    }
  });

  onDestroy(() => {
    // 清理状态
    analysisData = null;
    quickAnalysis = null;
    isAnalyzing = false;
    analysisError = '';
  });

  // 获取当前显示的数据
  $: displayData = quickAnalysis || analysisData;

  // 修改 safeGet 函数
  function safeGet(obj: any, path: string): string {
    if (!obj) return 'error';
    
    const value = path.split('.').reduce((acc: any, part: string) => {
      if (acc === null || acc === undefined) return null;
      return acc[part];
    }, obj);

    // 增加数据有效性检查
    if (value === undefined || value === null || value === '') return 'miss';
    if (typeof value === 'number' && isNaN(value)) return 'error';
    if (typeof value === 'string' && value.toLowerCase() === 'nan') return 'error';
    if (typeof value === 'string' && value.toLowerCase() === 'null') return 'miss';
    if (typeof value === 'string' && value.toLowerCase() === 'undefined') return 'miss';
    
    return value;
  }

  // 添加判断是否为快速分析数据的函数
  function isQuickAnalysis(data: any): boolean {
    return false;
  }

  // 修改显示逻辑
  $: showDetailButton = !isAnalyzing && analysisData;
  $: showQuickNotice = false;

  // 获取操作建议的颜色类
  function getActionColorClass(action: string): string {
    if (action.includes('买入')) return 'action-buy';
    if (action.includes('卖出') || action.includes('清仓')) return 'action-sell';
    return 'action-hold';
  }

  // 修改 formatScore 函数
  function formatScore(score: string, dependencies?: string[]): string {
    if (score === 'miss' || score === 'error' || score === '?') return '缺失';
    
    // 如果提供了依赖项，检查所有依赖项是否都有有效数据
    if (dependencies && dependencies.length > 0) {
      const hasMissingData = dependencies.some(dep => {
        const value = safeGet(displayData, dep);
        return value === 'miss' || value === 'error' || value === '?';
      });
      if (hasMissingData) return '缺失';
    }
    
    const num = parseInt(score);
    if (isNaN(num)) return '错误';
    return `${num}`; 
  }

  function formatValue(value: string): { value: string; class: string } {
    if (value === 'miss') return { value: '缺失', class: 'miss' };
    if (value === 'error') return { value: '错误', class: 'error' };
    if (value === '?') return { value: '待定', class: 'pending' };
    
    let formattedValue = value;
    let trendClass = '';
    
    // 处理趋势
    if (value.includes('↑')) {
      formattedValue = value.replace('↑', '');
      trendClass = 'trend-up';
    } else if (value.includes('↓')) {
      formattedValue = value.replace('↓', '');
      trendClass = 'trend-down';
    } else if (value.includes('→')) {
      formattedValue = value.replace('→', '');
      trendClass = 'trend-flat';
    }

    // 处理感叹号
    if (value.includes('!')) {
      formattedValue = formattedValue.replace('!', '');
      trendClass += ' trend-alert';
    }

    // 如果有百分比，确保正确显示
    if (formattedValue.includes('%')) {
      const [num, ...rest] = formattedValue.split('%');
      formattedValue = `${num}%${rest.join('')}`;
    }

    return { value: formattedValue, class: trendClass };
  }

  // 修改风险等级样式获取函数
  function getRiskClass(value: string): string {
    if (value === 'miss' || value === 'error' || value === '?') return '';
    const num = parseInt(value);
    if (isNaN(num)) return '';
    return num >= 70 ? 'risk-high' : num <= 30 ? 'risk-low' : 'risk-medium';
  }

  // 修改响应式声明
  $: revenueGrowthValue = formatValue(safeGet(analysisData, 'financialAnalysis.growth.revenueGrowth'));
  $: roeValue = formatValue(safeGet(analysisData, 'financialAnalysis.profitability.roe'));
  $: netMarginValue = formatValue(safeGet(analysisData, 'financialAnalysis.profitability.netMargin'));
  $: performanceScore = formatScore(safeGet(analysisData, 'scores.growthScore'));
  $: profitabilityScore = formatScore(safeGet(analysisData, 'scores.profitabilityScore'));
  $: healthScore = formatScore(safeGet(analysisData, 'scores.healthScore'));
  $: institutionalScore = formatScore(safeGet(analysisData, 'scores.institutionalScore'));

  // 修改建议相关的响应式声明
  $: action = safeGet(analysisData, 'suggestion.action');
  $: actionReason = safeGet(analysisData, 'suggestion.reason');

  let tooltipContent = '';
  let tooltipVisible = false;

  // 更新 tooltipMap
  const tooltipMap: Record<string, string> = {
    performanceScore: '综合评估公司的业绩表现',
    revenueGrowth: '营业收入的同比增长率',
    roe: '净资产收益率，反映股东权益的收益水平',
    netMargin: '净利润率，反映公司的盈利能力',
    profitabilityScore: '公司盈利能力的综合评分',
    healthScore: '公司财务健康状况评分',
    institutionalScore: '机构投资者对公司的认可度评分',
    risk: '综合风险评估分数'
  };

  // 修改 showTooltip 函数的类型
  function showTooltip(metric: keyof typeof tooltipMap) {
    tooltipContent = tooltipMap[metric] || '';
  }

  // 修改 hideTooltip 函数
  function hideTooltip() {
    tooltipContent = '';
  }

  // 修改 StockAnalysis 接口定义
  interface StockAnalysis {
    companyProfile: {
      industry: string;
      businessModel: string;
      coreBusiness: string;
      marketPosition: string;
      advantages: string[];
      challenges: string[];
    };
    financialAnalysis: {
      growth: {
        revenueGrowth: string;
        profitGrowth: string;
        growthStability: string;
        growthScore: string;
      };
      profitability: {
        grossMargin: string;
        netMargin: string;
        roe: string;
        profitabilityScore: string;
      };
      health: {
        assetLiabilityRatio: string;
        operatingCashFlow: string;
        investingCashFlow: string;
        financingCashFlow: string;
        cashFlowScore: string;
        healthScore: string;
      };
      institutionalRecognition: {
        fundHoldingRatio: string;
        fundCount: string;
        majorShareholderStability: string;
        institutionalScore: string;
      };
    };
    scores: {
      growthScore: string;
      profitabilityScore: string;
      healthScore: string;
      institutionalScore: string;
      totalScore: string;
    };
    suggestion: {
      summary: string;
      advantages: string[];
      risks: string[];
      action: string;
      reason: string;
    };
  }

  interface AnalysisResponse {
    type: 'quick' | 'full';
    data: StockAnalysis;
  }
</script>

{#if isAnalyzing && !analysisData}
  <!-- 加载状态 -->
  <div class="ai-analysis" in:fade={{duration: 200}}>
    <div class="header">
      <span>AI 智能分析</span>
      <button class="icon-button" disabled>
        <RefreshCwIcon size={14} class="spinning" />
      </button>
    </div>

    <!-- 添加错误提示 -->
    {#if analysisError}
      <div class="error-message" in:fade>
        <AlertTriangleIcon size={16} />
        <span>{analysisError}</span>
      </div>
    {/if}

    <div class="progress-view">
      <div class="progress-status">
        <span class="phase-text">{progress.phase === 'analysis' ? 'AI分析' : '数据收集'}</span>
        <span class="percent-text">{Math.round(getTotalProgress())}%</span>
      </div>
      
      <div class="progress-track">
        <div class="progress-bar {analysisError ? 'error' : ''}" style="width: {getTotalProgress()}%"></div>
      </div>
      
      <div class="step-info">
        <span class="step-text">
          {#if analysisError}
            分析已终止
          {:else}
            {getProgressDescription(progress.phase, progress.step)}
          {/if}
        </span>
        {#if progress.phase === 'data' && !analysisError}
          <span class="step-count">{progress.step}/{progress.total}</span>
        {/if}
      </div>

      <!-- 添加重试按钮 -->
      {#if analysisError}
        <button class="retry-button" on:click={() => handleRefresh()}>
          <RefreshCwIcon size={14} />
          重试分析
        </button>
      {/if}
    </div>
  </div>

{:else if analysisData}
  <div class="ai-analysis" in:fade={{duration: 200}}>
    <div class="header">
      <span>AI 智能分析</span>
      <button class="icon-button" on:click={handleRefresh}>
        <RefreshCwIcon size={14} />
      </button>
    </div>

    <div class="suggestion-view">
      <!-- 投资建议 -->
      <div class="action {getActionColorClass(analysisData.investment.recommendation)}">
        {analysisData.investment.recommendation}
      </div>
      
      <!-- 核心指标 -->
      <div class="metrics-grid">
        <!-- 业绩评分 -->
        <div class="metric-card">
          <div class="metric-header">
            <span class="metric-label">业绩评分</span>
            <InfoIcon size={16} class="info-icon" on:mouseenter={() => showTooltip('performanceScore')} on:mouseleave={hideTooltip} />
          </div>
          <div class="metric-value">
            {analysisData.financial.performanceScore}
            <span class="metric-unit">分</span>
          </div>
        </div>

        <!-- 盈利质量 -->
        <div class="metric-card">
          <div class="metric-header">
            <span class="metric-label">盈利质量</span>
            <InfoIcon size={16} class="info-icon" on:mouseenter={() => showTooltip('profitQuality')} on:mouseleave={hideTooltip} />
          </div>
          <div class="metric-value">
            {analysisData.financial.profitQuality}
          </div>
        </div>

        <!-- 风险等级 -->
        <div class="metric-card">
          <div class="metric-header">
            <span class="metric-label">风险等级</span>
            <InfoIcon size={16} class="info-icon" on:mouseenter={() => showTooltip('riskLevel')} on:mouseleave={hideTooltip} />
          </div>
          <div class="metric-value {getRiskClass(analysisData.risk.riskLevel)}">
            {analysisData.risk.riskLevel}
            <span class="metric-unit">分</span>
          </div>
        </div>

        <!-- 目标价格 -->
        <div class="metric-card">
          <div class="metric-header">
            <span class="metric-label">目标价格</span>
            <InfoIcon size={16} class="info-icon" on:mouseenter={() => showTooltip('targetPrice')} on:mouseleave={hideTooltip} />
          </div>
          <div class="metric-value">
            {analysisData.investment.targetPrice}
          </div>
        </div>
      </div>

      <!-- 关键指标趋势 -->
      <div class="trends-section">
        <h3>关键指标</h3>
        <div class="trends-grid">
          {#each analysisData.financial.keyMetrics as metric}
            <div class="trend-item">
              <span class="trend-label">{metric}</span>
            </div>
          {/each}
        </div>
      </div>

      <!-- 投资要点 -->
      <div class="key-points">
        <h3>投资要点</h3>
        <ul>
          {#each analysisData.investment.keyPoints as point}
            <li>{point}</li>
          {/each}
        </ul>
      </div>

      <!-- 主要风险 -->
      <div class="key-points">
        <h3>主要风险</h3>
        <ul>
          {#each analysisData.risk.mainRisks as risk}
            <li>{risk}</li>
          {/each}
        </ul>
      </div>

      <!-- 查看详情按钮 -->
      <button class="detail-button" on:click={() => showDetailModal = true}>
        查看详细分析
      </button>
    </div>
  </div>
{/if}

{#if tooltipContent}
  <div class="tooltip">
    {tooltipContent}
  </div>
{/if}

{#if showDetailModal}
  <Modal on:close={() => showDetailModal = false}>
    <AIAnalysisDetailModal data={analysisData} />
  </Modal>
{/if}

<style>
  .ai-analysis {
    display: flex;
    flex-direction: column;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    border-radius: 12px;
    overflow: hidden;
  }

  .header {
    padding: 12px 16px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    font-weight: 500;
    color: var(--text-primary);
    border-bottom: 1px solid var(--border-color);
    background: var(--surface);
  }

  .icon-button {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: none;
    color: var(--text-secondary);
    cursor: pointer;
    border-radius: 6px;
    transition: var(--transition-colors);
  }

  .icon-button:hover:not(:disabled) {
    background: var(--hover-bg);
    color: var(--text-primary);
  }

  .icon-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .spinning {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* 进度视图 */
  .progress-view {
    padding: 16px;
    background: var(--surface);
  }

  .progress-status {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .phase-text {
    font-size: 14px;
    font-weight: 500;
    color: var(--text-primary);
  }

  .percent-text {
    font-size: 14px;
    font-weight: 600;
    color: var(--primary-500);
  }

  .progress-track {
    height: 4px;
    background: var(--hover-bg);
    border-radius: 2px;
    overflow: hidden;
  }

  .progress-bar {
    height: 100%;
    background: var(--primary-500);
    transition: width 0.3s ease-out;
    border-radius: 2px;
  }

  .step-info {
    margin-top: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .step-text {
    font-size: 13px;
    color: var(--text-secondary);
  }

  .step-count {
    font-size: 12px;
    color: var(--text-secondary);
  }

  /* 建议视图 */
  .suggestion-view {
    padding: 20px 16px;
    text-align: center;
    background: var(--surface);
  }

  .action {
    font-size: 28px;
    font-weight: 600;
    margin: 24px auto;
    padding: 16px 32px;
    border-radius: 16px;
    text-align: center;
    max-width: fit-content;
    transition: all 0.2s ease;
  }

  .action.action-buy {
    color: var(--success-text);
    background: var(--success-bg);
    border: 1px solid var(--success-text);
  }

  .action.action-sell {
    color: var(--danger-text);
    background: var(--danger-bg);
    border: 1px solid var(--danger-text);
  }

  .action.action-hold {
    color: var(--warning-text);
    background: var(--warning-bg);
    border: 1px solid var(--warning-text);
  }

  .metrics {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    padding: 16px;
    margin: 0 auto;
    max-width: 600px;
  }

  .metric {
    padding: 16px;
    border-radius: 12px;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    min-width: 0;
    transition: all 0.2s ease;
  }

  .metric:hover {
    border-color: var(--primary-300);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    transform: translateY(-1px);
  }

  .value {
    font-size: 20px;
    font-weight: 600;
    color: var(--text-primary);
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .label {
    font-size: 14px;
    color: var(--text-secondary);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 4px;
  }

  .reason {
    margin-top: 16px;
    font-size: 14px;
    color: var(--text-secondary);
  }

  /* 按钮和提示 */
  .action-button {
    margin: 0 16px 16px;
    width: calc(100% - 32px);
    padding: 12px;
    background: var(--primary-500);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: var(--transition-colors);
  }

  .action-button:hover {
    background: var(--primary-600);
  }

  .notice {
    margin: 0 16px 16px;
    padding: 12px;
    background: var(--surface-variant);
    border-radius: 8px;
    text-align: center;
    font-size: 14px;
    color: var(--text-secondary);
    border: 1px solid var(--border-color);
  }

  .dots {
    display: inline-block;
    width: 12px;
    height: 12px;
    margin-left: 4px;
  }

  .dots::after {
    content: '...';
    animation: dots 1.5s steps(4, end) infinite;
  }

  @keyframes dots {
    0%, 20% { content: '.'; }
    40% { content: '..'; }
    60% { content: '...'; }
    80%, 100% { content: ''; }
  }

  /* 空状态 */
  .empty {
    padding: 32px 16px;
    text-align: center;
    background: var(--surface);
  }

  .empty span {
    display: block;
    margin-bottom: 16px;
    color: var(--text-secondary);
  }

  /* 风险样式 */
  .risk-high { color: var(--danger-text, #dc2626) !important; }
  .risk-medium { color: var(--warning-text, #d97706) !important; }
  .risk-low { color: var(--success-text, #059669) !important; }

  /* 添加新的 CSS 变量 */
  :root {
    --success-text: #059669;
    --danger-text: #dc2626;
    --warning-text: #d97706;
    --success-bg: #ecfdf5;
    --danger-bg: #fef2f2;
    --warning-bg: #fffbeb;
  }

  /* 暗色模式变量覆盖 */
  :global(.dark) {
    --success-text: #34d399;
    --danger-text: #f87171;
    --warning-text: #fbbf24;
    --success-bg: rgba(52, 211, 153, 0.1);
    --danger-bg: rgba(248, 113, 113, 0.1);
    --warning-bg: rgba(251, 191, 36, 0.1);
  }

  /* 趋势样式 */
  .trend-up {
    color: var(--success-text);
  }

  .trend-up::after {
    content: '▲';
    font-size: 12px;
    margin-left: 4px;
  }

  .trend-down {
    color: var(--danger-text);
  }

  .trend-down::after {
    content: '▼';
    font-size: 12px;
    margin-left: 4px;
  }

  .trend-flat {
    color: var(--warning-text);
  }

  .trend-flat::after {
    content: '─';
    font-size: 12px;
    margin-left: 4px;
  }

  .trend-alert {
    font-weight: 700;
  }

  .trend-alert::before {
    content: '!';
    color: var(--warning-text);
    font-size: 14px;
    margin-right: 4px;
  }

  /* 状态样式 */
  .miss, .error, .pending {
    font-size: 16px;
    opacity: 0.8;
  }

  .miss { color: var(--warning-text); }
  .error { color: var(--danger-text); }
  .pending { color: var(--info-text); }

  /* 标签样式 */
  .label {
    font-size: 13px;
    color: var(--text-secondary);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 4px;
  }

  .tooltip {
    position: fixed;
    background: rgba(0, 0, 0, 0.8);
    color: white;
    padding: 8px 12px;
    border-radius: 4px;
    font-size: 14px;
    max-width: 300px;
    z-index: 1000;
    pointer-events: none;
    transition: opacity 0.2s;
  }

  /* 添加新的样式 */
  .metrics-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    margin: 20px 0;
  }

  .metric-card {
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 12px;
    padding: 16px;
    transition: all 0.2s ease;
  }

  .metric-card:hover {
    border-color: var(--primary-300);
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .metric-header {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 8px;
  }

  .metric-label {
    font-size: 14px;
    color: var(--text-secondary);
  }

  .metric-value {
    font-size: 24px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .metric-unit {
    font-size: 14px;
    color: var(--text-secondary);
    margin-left: 4px;
  }

  .trends-section {
    margin: 24px 0;
  }

  .trends-section h3 {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 16px;
    color: var(--text-primary);
  }

  .trends-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .trend-item {
    padding: 12px;
    background: var(--surface-variant);
    border-radius: 8px;
    font-size: 14px;
  }

  .key-points {
    margin: 24px 0;
    text-align: left;
  }

  .key-points h3 {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 16px;
    color: var(--text-primary);
  }

  .key-points ul {
    list-style: none;
    padding: 0;
  }

  .key-points li {
    margin: 8px 0;
    padding-left: 20px;
    position: relative;
    font-size: 14px;
    color: var(--text-secondary);
  }

  .key-points li::before {
    content: "•";
    position: absolute;
    left: 0;
    color: var(--primary-500);
  }

  .detail-button {
    width: 100%;
    padding: 12px;
    background: var(--primary-500);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s ease;
  }

  .detail-button:hover {
    background: var(--primary-600);
  }

  .error-message {
    margin: 12px 16px;
    padding: 12px;
    background: var(--danger-bg);
    color: var(--danger-text);
    border-radius: 8px;
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .progress-bar.error {
    background: var(--danger-text);
  }

  .retry-button {
    margin-top: 16px;
    padding: 8px 16px;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 6px;
    color: var(--text-primary);
    font-size: 14px;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 8px;
    transition: all 0.2s ease;
  }

  .retry-button:hover {
    background: var(--hover-bg);
    border-color: var(--primary-300);
  }
</style> 