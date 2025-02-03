<script lang="ts">
  import type { StockAnalysis } from '../../types/analysis'
  import Modal from '../Modal.svelte'
  import { AlertTriangleIcon, TrendingUpIcon, TrendingDownIcon, MinusIcon } from 'lucide-svelte'

  export let show = false
  export let data: StockAnalysis | null = null

  // 获取操作建议的颜色类
  function getActionColorClass(action: string): string {
    if (action.includes('买入')) return 'action-buy'
    if (action.includes('卖出') || action.includes('清仓')) return 'action-sell'
    return 'action-hold'
  }

  // 安全获取数据
  function safeGet(obj: any, path: string, defaultValue: any = '--'): any {
    return path.split('.').reduce((acc, part) => acc && acc[part], obj) ?? defaultValue
  }

  // 格式化趋势值
  function formatTrendValue(value: string): string {
    if (!value || value === '--') return '--'
    if (value.includes('↑')) return `${value} 📈`
    if (value.includes('↓')) return `${value} 📉`
    if (value.includes('→')) return `${value} ➡️`
    return value
  }

  // 格式化评分
  function formatScore(score: string): string {
    if (!score || score === '--') return '--'
    const num = parseInt(score)
    if (isNaN(num)) return score
    return `${num}/100`
  }

  // 获取评分颜色类
  function getScoreColorClass(score: string): string {
    const num = parseInt(score)
    if (isNaN(num)) return ''
    if (num >= 80) return 'score-high'
    if (num <= 40) return 'score-low'
    return 'score-medium'
  }
</script>

<Modal bind:show title="AI 智能分析详情">
  {#if data}
    <div class="analysis-detail">
      <!-- 公司概况 -->
      <section>
        <h3>公司概况</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="label">所属行业</span>
            <span class="value">{data.companyProfile.industry}</span>
          </div>
          <div class="info-item">
            <span class="label">商业模式</span>
            <span class="value">{data.companyProfile.businessModel}</span>
          </div>
          <div class="info-item">
            <span class="label">核心业务</span>
            <span class="value">{data.companyProfile.coreBusiness}</span>
          </div>
          <div class="info-item">
            <span class="label">市场地位</span>
            <span class="value">{data.companyProfile.marketPosition}</span>
          </div>
        </div>

        <div class="advantages-challenges">
          <div class="section">
            <h4>竞争优势</h4>
            <ul>
              {#each data.companyProfile.advantages as advantage}
                <li>{advantage}</li>
              {/each}
            </ul>
          </div>
          <div class="section">
            <h4>面临挑战</h4>
            <ul>
              {#each data.companyProfile.challenges as challenge}
                <li>{challenge}</li>
              {/each}
            </ul>
          </div>
        </div>
      </section>

      <!-- 财务分析 -->
      <section>
        <h3>财务分析</h3>
        <div class="financial-analysis">
          <!-- 成长性 -->
          <div class="analysis-section">
            <h4>成长性分析</h4>
            <div class="metrics">
              <div class="metric">
                <span class="label">营收增长</span>
                <span class="value">{data.financialAnalysis.growth.revenueGrowth}</span>
              </div>
              <div class="metric">
                <span class="label">利润增长</span>
                <span class="value">{data.financialAnalysis.growth.profitGrowth}</span>
              </div>
              <div class="metric">
                <span class="label">增长稳定性</span>
                <span class="value">{data.financialAnalysis.growth.growthStability}</span>
              </div>
              <div class="metric">
                <span class="label">成长性评分</span>
                <span class="value">{data.financialAnalysis.growth.growthScore}</span>
              </div>
            </div>
          </div>

          <!-- 盈利能力 -->
          <div class="analysis-section">
            <h4>盈利能力分析</h4>
            <div class="metrics">
              <div class="metric">
                <span class="label">毛利率</span>
                <span class="value">{data.financialAnalysis.profitability.grossMargin}</span>
              </div>
              <div class="metric">
                <span class="label">净利率</span>
                <span class="value">{data.financialAnalysis.profitability.netMargin}</span>
              </div>
              <div class="metric">
                <span class="label">ROE</span>
                <span class="value">{data.financialAnalysis.profitability.roe}</span>
              </div>
              <div class="metric">
                <span class="label">盈利能力评分</span>
                <span class="value">{data.financialAnalysis.profitability.profitabilityScore}</span>
              </div>
            </div>
          </div>

          <!-- 财务健康 -->
          <div class="analysis-section">
            <h4>财务健康分析</h4>
            <div class="metrics">
              <div class="metric">
                <span class="label">资产负债率</span>
                <span class="value">{data.financialAnalysis.health.assetLiabilityRatio}</span>
              </div>
              <div class="metric">
                <span class="label">经营现金流</span>
                <span class="value">{data.financialAnalysis.health.operatingCashFlow}</span>
              </div>
              <div class="metric">
                <span class="label">投资现金流</span>
                <span class="value">{data.financialAnalysis.health.investingCashFlow}</span>
              </div>
              <div class="metric">
                <span class="label">筹资现金流</span>
                <span class="value">{data.financialAnalysis.health.financingCashFlow}</span>
              </div>
              <div class="metric">
                <span class="label">现金流评分</span>
                <span class="value">{data.financialAnalysis.health.cashFlowScore}</span>
              </div>
              <div class="metric">
                <span class="label">健康评分</span>
                <span class="value">{data.financialAnalysis.health.healthScore}</span>
              </div>
            </div>
          </div>

          <!-- 机构认可度 -->
          <div class="analysis-section">
            <h4>机构认可度分析</h4>
            <div class="metrics">
              <div class="metric">
                <span class="label">基金持股比例</span>
                <span class="value">{data.financialAnalysis.institutionalRecognition.fundHoldingRatio}</span>
              </div>
              <div class="metric">
                <span class="label">基金持股数量</span>
                <span class="value">{data.financialAnalysis.institutionalRecognition.fundCount}</span>
              </div>
              <div class="metric">
                <span class="label">大股东稳定性</span>
                <span class="value">{data.financialAnalysis.institutionalRecognition.majorShareholderStability}</span>
              </div>
              <div class="metric">
                <span class="label">机构认可度评分</span>
                <span class="value">{data.financialAnalysis.institutionalRecognition.institutionalScore}</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 投资建议 -->
      <section>
        <h3>投资建议</h3>
        <div class="suggestion">
          <p class="summary">{data.suggestion.summary}</p>
          
          <div class="advantages-risks">
            <div class="section">
              <h4>投资优势</h4>
              <ul>
                {#each data.suggestion.advantages as advantage}
                  <li>{advantage}</li>
                {/each}
              </ul>
            </div>
            <div class="section">
              <h4>投资风险</h4>
              <ul>
                {#each data.suggestion.risks as risk}
                  <li>{risk}</li>
                {/each}
              </ul>
            </div>
          </div>

          <div class="action-section">
            <div class="action {data.suggestion.action.includes('买入') ? 'buy' : data.suggestion.action.includes('卖出') ? 'sell' : 'hold'}">
              {data.suggestion.action}
            </div>
            <p class="reason">{data.suggestion.reason}</p>
          </div>
        </div>
      </section>
    </div>
  {/if}
</Modal>

<style>
  .analysis-detail {
    padding: 24px;
    max-width: 800px;
    margin: 0 auto;
  }

  section {
    margin-bottom: 32px;
  }

  h3 {
    font-size: 20px;
    font-weight: 600;
    margin-bottom: 16px;
    color: var(--text-primary);
  }

  h4 {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 12px;
    color: var(--text-primary);
  }

  .info-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 16px;
    margin-bottom: 24px;
  }

  .info-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .label {
    font-size: 14px;
    color: var(--text-secondary);
  }

  .value {
    font-size: 16px;
    color: var(--text-primary);
  }

  .advantages-challenges, .advantages-risks {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 24px;
  }

  .section {
    background: var(--surface-variant);
    padding: 16px;
    border-radius: 8px;
  }

  ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  li {
    margin-bottom: 8px;
    padding-left: 20px;
    position: relative;
  }

  li:before {
    content: "•";
    position: absolute;
    left: 0;
    color: var(--primary-500);
  }

  .financial-analysis {
    display: grid;
    gap: 24px;
  }

  .analysis-section {
    background: var(--surface-variant);
    padding: 16px;
    border-radius: 8px;
  }

  .metrics {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 16px;
  }

  .metric {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .suggestion {
    background: var(--surface-variant);
    padding: 24px;
    border-radius: 8px;
  }

  .summary {
    font-size: 16px;
    line-height: 1.5;
    margin-bottom: 24px;
    color: var(--text-primary);
  }

  .action-section {
    margin-top: 24px;
    text-align: center;
  }

  .action {
    display: inline-block;
    padding: 8px 24px;
    border-radius: 20px;
    font-weight: 600;
    margin-bottom: 16px;
  }

  .action.buy {
    background: var(--success-container);
    color: var(--success-text);
  }

  .action.sell {
    background: var(--danger-container);
    color: var(--danger-text);
  }

  .action.hold {
    background: var(--warning-container);
    color: var(--warning-text);
  }

  .reason {
    font-size: 14px;
    color: var(--text-secondary);
    max-width: 600px;
    margin: 0 auto;
  }
</style> 