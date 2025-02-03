package api

import (
	"github.com/yankeguo/zhipu"
)

// GLM4Response 智谱AI响应结构
type GLM4Response struct {
	ID      string                    `json:"id"`
	Created int64                     `json:"created"`
	Model   string                    `json:"model"`
	Usage   zhipu.ChatCompletionUsage `json:"usage"`
	Choices []struct {
		Index   int         `json:"index"`
		Message GLM4Message `json:"message"`
	} `json:"choices"`
}

// GLM4Message 消息结构
type GLM4Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StockAnalysis 股票分析结果
type StockAnalysis struct {
	CompanyProfile CompanyProfile `json:"companyProfile"` // 公司概况
	Financial      Financial      `json:"financial"`      // 财务分析
	Operation      Operation      `json:"operation"`      // 经营分析
	Risk           Risk           `json:"risk"`           // 风险分析
	Investment     Investment     `json:"investment"`     // 投资建议
}

// CompanyProfile 公司概况
type CompanyProfile struct {
	Industry       string   `json:"industry"`       // 所属行业
	MarketPosition string   `json:"marketPosition"` // 市场地位
	BusinessModel  string   `json:"businessModel"`  // 商业模式
	CoreBusiness   string   `json:"coreBusiness"`   // 核心业务
	Advantages     []string `json:"advantages"`     // 竞争优势
	Challenges     []string `json:"challenges"`     // 面临挑战
}

// Financial 财务分析
type Financial struct {
	PerformanceScore string   `json:"performanceScore"` // 整体业绩评分(1-100)
	GrowthTrend      string   `json:"growthTrend"`      // 增长趋势
	ProfitQuality    string   `json:"profitQuality"`    // 盈利质量
	KeyMetrics       []string `json:"keyMetrics"`       // 关键指标
	Concerns         []string `json:"concerns"`         // 需关注问题
}

// Operation 经营分析
type Operation struct {
	MarketShare     string   `json:"marketShare"`     // 市场份额
	CompetitiveEdge string   `json:"competitiveEdge"` // 竞争优势
	Efficiency      string   `json:"efficiency"`      // 运营效率
	Strengths       []string `json:"strengths"`       // 优势
	Weaknesses      []string `json:"weaknesses"`      // 劣势
}

// Risk 风险分析
type Risk struct {
	RiskLevel    string   `json:"riskLevel"`    // 风险等级(1-100)
	RiskTrend    string   `json:"riskTrend"`    // 风险趋势
	MainRisks    []string `json:"mainRisks"`    // 主要风险
	SpecialNotes []string `json:"specialNotes"` // 特别说明
}

// Investment 投资建议
type Investment struct {
	Recommendation string   `json:"recommendation"` // 投资建议
	TargetPrice    string   `json:"targetPrice"`    // 目标价格
	StopLoss       string   `json:"stopLoss"`       // 止损价格
	TimeHorizon    string   `json:"timeHorizon"`    // 投资期限
	KeyPoints      []string `json:"keyPoints"`      // 关键要点
}
