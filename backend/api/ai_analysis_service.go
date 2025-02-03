package api

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"stock-helper-svelte/backend/api/company"
	"stock-helper-svelte/backend/api/financial"
	"stock-helper-svelte/backend/api/market"
	"stock-helper-svelte/backend/api/types"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yankeguo/zhipu"
)

// EventAnalysisProgress 分析进度事件
const EventAnalysisProgress = "analysis:progress"

// ProgressMessage 进度消息
type ProgressMessage struct {
	Step    int    `json:"step"`
	Total   int    `json:"total"`
	Message string `json:"message"`
	Phase   string `json:"phase"`
}

// Service AI分析服务
type Service struct {
	zhipuClient *zhipu.Client
	company     *company.Client
	market      *market.Client
	financial   *financial.Client
	apiClient   *Client
	ctx         context.Context
}

// NewService 创建新的AI分析服务
func NewService(companyClient *company.Client, marketClient *market.Client, financialClient *financial.Client, zhipuApiKey string, apiClient *Client, ctx context.Context) *Service {
	client, _ := zhipu.NewClient(zhipu.WithAPIKey(zhipuApiKey))
	return &Service{
		zhipuClient: client,
		company:     companyClient,
		market:      marketClient,
		financial:   financialClient,
		apiClient:   apiClient,
		ctx:         ctx,
	}
}

// AnalyzeStock 分析股票
func (s *Service) AnalyzeStock(code string) (*StockAnalysis, error) {
	fmt.Printf("开始分析股票: %s\n", code)

	// 1. 获取股票数据
	stockData, err := s.getStockData(code)
	if err != nil {
		fmt.Printf("获取股票数据失败: %v\n", err)
		return nil, fmt.Errorf("获取股票数据失败: %w", err)
	}

	// 2. 验证数据完整性
	if strings.Contains(stockData, "数据为空") {
		fmt.Println("股票数据不完整")
		return nil, fmt.Errorf("股票数据不完整")
	}

	// 3. 发送 AI 分析状态
	fmt.Println("发送AI分析进度事件")
	runtime.EventsEmit(s.ctx, EventAnalysisProgress, ProgressMessage{
		Step:    0,
		Total:   0,
		Message: "AI 正在进行分析",
		Phase:   "analysis",
	})

	// 4. 进行完整分析
	fmt.Println("开始执行AI分析")
	analysis, err := s.performAnalysis(stockData)
	if err != nil {
		fmt.Printf("分析失败: %v\n", err)
		return nil, fmt.Errorf("分析失败: %w", err)
	}

	fmt.Println("分析完成")
	return analysis, nil
}

// getStockData 获取股票数据
func (s *Service) getStockData(code string) (string, error) {
	fmt.Printf("开始获取股票数据: %s\n", code)
	var result strings.Builder
	totalSteps := 7 // 总步骤数

	// 发送进度通知函数
	sendProgress := func(step int, message string) {
		fmt.Printf("发送进度更新: step=%d, message=%s\n", step, message)
		// 修改事件名称以匹配前端监听
		runtime.EventsEmit(s.ctx, "analysis_progress", ProgressMessage{
			Step:    step,
			Total:   totalSteps,
			Message: message,
			Phase:   "data",
		})
	}

	// 1. 获取公司详细信息
	sendProgress(1, "正在获取公司信息...")
	profile, err := s.company.GetCompanyProfile(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("获取公司详细信息失败: %v", err)
	}
	result.WriteString(formatCompanyProfile(profile))
	result.WriteString("\n\n")

	// 2. 获取财务数据（最近5年每个季度）
	sendProgress(2, "正在获取财务数据...")
	financialData, err := s.getFinancialDataForPeriods(code)
	if err != nil {
		return "", fmt.Errorf("获取财务数据失败: %v", err)
	}
	result.WriteString(formatFinancialDataHistory(financialData))
	result.WriteString("\n\n")

	// 3. 获取现金流
	sendProgress(3, "正在获取现金流数据...")
	cashflows, err := s.company.GetQuarterlyCashFlow(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("获取现金流失败: %v", err)
	}
	result.WriteString(formatCashFlow(cashflows))
	result.WriteString("\n\n")

	// 4. 获取股东信息
	sendProgress(4, "正在获取股东信息...")
	shareholders, err := s.company.GetShareholders(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("获取股东信息失败: %v", err)
	}
	result.WriteString(formatShareholders(shareholders))
	result.WriteString("\n\n")

	// 5. 获取基金持股
	sendProgress(5, "正在获取机构持股...")
	fundHoldings, err := s.company.GetFundHoldings(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("获取基金持股失败: %v", err)
	}
	result.WriteString(formatFundHoldings(fundHoldings))
	result.WriteString("\n\n")

	// 6. 获取行业数据
	sendProgress(6, "正在获取行业数据...")
	// TODO: 添加行业数据获取逻辑

	// 7. 整合分析数据
	sendProgress(7, "正在整合分析数据...")
	// 数据整合逻辑已经在前面完成

	// 完成
	runtime.EventsEmit(s.ctx, "analysis_progress", ProgressMessage{
		Step:    totalSteps,
		Total:   totalSteps,
		Message: "数据获取完成",
		Phase:   "data",
	})

	return result.String(), nil
}

// getFinancialDataForPeriods 获取最近5年每个季度的财务数据
func (s *Service) getFinancialDataForPeriods(code string) ([]*types.StockFinancialData, error) {
	var result []*types.StockFinancialData
	currentYear := time.Now().Year()
	currentQuarter := (int(time.Now().Month())-1)/3 + 1

	// 获取最近20个季度（5年）的数据
	for i := 0; i < 20; i++ {
		year := currentYear
		quarter := currentQuarter - i

		// 调整年份和季度
		for quarter <= 0 {
			year--
			quarter += 4
		}

		data, err := s.financial.GetStockFinancialData(context.Background(), code, year, quarter)
		if err != nil {
			continue // 跳过获取失败的季度
		}
		result = append(result, data)
	}

	// 按时间排序（从新到旧）
	sort.Slice(result, func(i, j int) bool {
		if result[i].Year != result[j].Year {
			return result[i].Year > result[j].Year
		}
		return result[i].Quarter > result[j].Quarter
	})

	return result, nil
}

// formatFinancialDataHistory 格式化历史财务数据
func formatFinancialDataHistory(data []*types.StockFinancialData) string {
	if len(data) == 0 {
		return "历史财务数据为空"
	}

	var result strings.Builder
	result.WriteString("历史财务数据分析:\n\n")

	// 按年度和季度组织数据
	for _, period := range data {
		result.WriteString(fmt.Sprintf("=== %d年第%d季度 ===\n", period.Year, period.Quarter))

		if period.Profitability != nil {
			result.WriteString(fmt.Sprintf("盈利能力指标:\n"+
				"营业收入: %.2f百万元\n"+
				"净利润: %.2f百万元\n"+
				"毛利率: %.2f%%\n"+
				"净利率: %.2f%%\n"+
				"ROE: %.2f%%\n",
				period.Profitability.Yysr,
				period.Profitability.Jlr,
				period.Profitability.Mll,
				period.Profitability.Jll,
				period.Profitability.Jzcsy))
		}

		if period.Growth != nil {
			result.WriteString(fmt.Sprintf("\n增长指标:\n"+
				"营收增长: %.2f%%\n"+
				"净利润增长: %.2f%%\n"+
				"净资产增长: %.2f%%\n",
				period.Growth.Zyzzl,
				period.Growth.Jlrzzl,
				period.Growth.Jzczzl))
		}

		if period.Solvency != nil {
			result.WriteString(fmt.Sprintf("\n偿债能力:\n"+
				"资产负债率: %.2f%%\n"+
				"流动比率: %.2f\n"+
				"速动比率: %.2f\n",
				period.Solvency.Zcfzl,
				period.Solvency.Ldbl,
				period.Solvency.Sdbl))
		}

		result.WriteString("\n")
	}

	// 添加同比环比分析
	if len(data) >= 4 {
		result.WriteString("\n同比环比分析:\n")
		latest := data[0]
		lastYear := data[3]    // 去年同期
		lastQuarter := data[1] // 上一季度

		if latest.Profitability != nil && lastYear.Profitability != nil {
			yoyRevenue := (latest.Profitability.Yysr - lastYear.Profitability.Yysr) / lastYear.Profitability.Yysr * 100
			yoyProfit := (latest.Profitability.Jlr - lastYear.Profitability.Jlr) / lastYear.Profitability.Jlr * 100
			result.WriteString(fmt.Sprintf("同比增长:\n营收: %.2f%%\n净利润: %.2f%%\n", yoyRevenue, yoyProfit))
		}

		if latest.Profitability != nil && lastQuarter.Profitability != nil {
			qoqRevenue := (latest.Profitability.Yysr - lastQuarter.Profitability.Yysr) / lastQuarter.Profitability.Yysr * 100
			qoqProfit := (latest.Profitability.Jlr - lastQuarter.Profitability.Jlr) / lastQuarter.Profitability.Jlr * 100
			result.WriteString(fmt.Sprintf("环比增长:\n营收: %.2f%%\n净利润: %.2f%%\n", qoqRevenue, qoqProfit))
		}
	}

	return result.String()
}

// performAnalysis 执行分析
func (s *Service) performAnalysis(stockData string) (*StockAnalysis, error) {
	systemPrompt := `你是一位专业的股票分析师，需要对公司进行深入的基本面分析。
请基于提供的数据，按照以下JSON结构体定义进行分析输出：

{
    "companyProfile": {
        "industry": "所属行业",
        "marketPosition": "市场地位",
        "businessModel": "商业模式",
        "coreBusiness": "核心业务",
        "advantages": ["竞争优势1", "竞争优势2", "..."],
        "challenges": ["面临挑战1", "面临挑战2", "..."]
    },
    "financial": {
        "performanceScore": "整体业绩评分(1-100)",
        "growthTrend": "增长趋势描述",
        "profitQuality": "盈利质量分析",
        "keyMetrics": ["关键指标1", "关键指标2", "..."],
        "concerns": ["需关注问题1", "需关注问题2", "..."]
    },
    "operation": {
        "marketShare": "市场份额情况",
        "competitiveEdge": "竞争优势分析",
        "efficiency": "运营效率评估",
        "strengths": ["优势1", "优势2", "..."],
        "weaknesses": ["劣势1", "劣势2", "..."]
    },
    "risk": {
        "riskLevel": "风险等级(1-100)",
        "riskTrend": "风险趋势分析",
        "mainRisks": ["主要风险1", "主要风险2", "..."],
        "specialNotes": ["特别说明1", "特别说明2", "..."]
    },
    "investment": {
        "recommendation": "投资建议",
        "targetPrice": "目标价格",
        "stopLoss": "止损价格",
        "timeHorizon": "投资期限",
        "keyPoints": ["关键要点1", "关键要点2", "..."]
    }
}

请确保：
1. 严格按照上述JSON结构输出
2. 所有数组至少包含3个元素
3. 评分类数值(performanceScore, riskLevel)使用1-100的范围
4. 价格相关数值需包含具体数字
5. 分析内容要客观、专业、具体`

	service := s.zhipuClient.ChatCompletion("glm-4-flash").
		AddMessage(zhipu.ChatCompletionMessage{
			Role:    "system",
			Content: systemPrompt,
		}).
		AddMessage(zhipu.ChatCompletionMessage{
			Role:    "user",
			Content: fmt.Sprintf("请对以下公司数据进行分析并以JSON格式输出：\n\n%s", stockData),
		})

	service.SetResponseFormat("json_object")

	response, err := service.Do(context.Background())
	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("分析返回结果为空")
	}

	var analysis StockAnalysis
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &analysis); err != nil {
		return nil, fmt.Errorf("解析分析结果失败: %w", err)
	}

	// 验证关键字段
	if analysis.CompanyProfile.Industry == "" ||
		analysis.Financial.PerformanceScore == "" ||
		analysis.Investment.Recommendation == "" {
		return nil, fmt.Errorf("分析结果关键字段缺失")
	}

	return &analysis, nil
}

// formatCompanyProfile 格式化公司概况
func formatCompanyProfile(profile *types.CompanyProfile) string {
	return fmt.Sprintf(`公司概况:
名称: %s
所属市场: %s
概念板块: %s
主营业务: %s
上市时间: %s
发行价: %s元`,
		profile.Name,
		profile.Market,
		profile.Concepts,
		profile.BusinessScope,
		profile.ListDate,
		profile.IPOPrice)
}

// formatCashFlow 格式化现金流
func formatCashFlow(cashflows []types.QuarterlyCashFlow) string {
	if len(cashflows) == 0 {
		return "现金流数据为空"
	}

	var result strings.Builder
	result.WriteString("现金流(最近5期):\n")

	// 只取最近5期数据
	end := len(cashflows)
	if end > 5 {
		end = 5
	}

	for i := 0; i < end; i++ {
		cf := cashflows[i]
		result.WriteString(fmt.Sprintf(`
%s:
经营活动:
- 流入: %s万元
- 流出: %s万元
- 净额: %s万元
投资活动:
- 流入: %s万元
- 流出: %s万元
- 净额: %s万元
筹资活动:
- 流入: %s万元
- 流出: %s万元
- 净额: %s万元
现金净增加: %s万元`,
			cf.Date,
			cf.OperateInflow,
			cf.OperateOutflow,
			cf.OperateNet,
			cf.InvestInflow,
			cf.InvestOutflow,
			cf.InvestNet,
			cf.FinanceInflow,
			cf.FinanceOutflow,
			cf.FinanceNet,
			cf.CashIncrease))
	}

	return result.String()
}

// formatShareholders 格式化股东信息
func formatShareholders(shareholders []types.ShareholderInfo) string {
	if len(shareholders) == 0 {
		return "股东信息为空"
	}

	var result strings.Builder
	result.WriteString("主要股东:\n")

	for _, sh := range shareholders {
		result.WriteString(fmt.Sprintf(`
%s:
持股数: %s万股
占比: %s%%
性质: %s
变动: %s万股 (%s)`,
			sh.Name,
			sh.SharesHeld,
			sh.Percentage,
			sh.Nature,
			sh.ShareChange,
			sh.ChangeRatio))
	}

	return result.String()
}

// formatFundHoldings 格式化基金持股
func formatFundHoldings(holdings []types.FundHolding) string {
	if len(holdings) == 0 {
		return "基金持股数据为空"
	}

	var result strings.Builder
	result.WriteString("基金持股:\n")

	for _, h := range holdings {
		result.WriteString(fmt.Sprintf(`
%s (%s):
持股数: %.2f万股
占流通股: %.2f%%
市值: %.2f万元
占净值: %.2f%%`,
			h.FundName,
			h.FundCode,
			h.SharesHeld/10000,
			h.FloatRatio,
			h.MarketValue/10000,
			h.NetWorthRatio))
	}

	return result.String()
}

// GetStockAnalysis 获取股票数据
func (s *Service) GetStockAnalysis(code string) (string, error) {
	return s.getStockData(code)
}

// Chat 通用对话接口
func (s *Service) Chat(messages []GLM4Message) (*GLM4Response, error) {
	// 转换消息格式
	var zhipuMessages []zhipu.ChatCompletionMessage
	for _, msg := range messages {
		zhipuMessages = append(zhipuMessages, zhipu.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 创建对话服务
	service := s.zhipuClient.ChatCompletion("glm-4-flash")
	for _, msg := range zhipuMessages {
		service.AddMessage(msg)
	}

	// 执行对话
	response, err := service.Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("对话失败: %w", err)
	}

	// 转换返回格式
	result := &GLM4Response{
		ID:      response.ID,
		Created: response.Created,
		Model:   response.Model,
		Usage:   response.Usage,
	}

	if len(response.Choices) > 0 {
		result.Choices = []struct {
			Index   int         `json:"index"`
			Message GLM4Message `json:"message"`
		}{
			{
				Index: response.Choices[0].Index,
				Message: GLM4Message{
					Role:    response.Choices[0].Message.Role,
					Content: response.Choices[0].Message.Content,
				},
			},
		}
	}

	return result, nil
}
