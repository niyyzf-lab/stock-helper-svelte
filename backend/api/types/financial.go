package types

import (
	"fmt"
	"strings"
)

// ProfitabilityData 盈利能力数据
type ProfitabilityData struct {
	Dm     string  `json:"dm"`     // 代码
	Mc     string  `json:"mc"`     // 名称
	Jzcsy  float64 `json:"jzcsy"`  // 净资产收益率(%)
	Jll    float64 `json:"jll"`    // 净利率(%)
	Mll    float64 `json:"mll"`    // 毛利率(%)
	Jlr    float64 `json:"jlr"`    // 净利润(百万元)
	Mgsy   float64 `json:"mgsy"`   // 每股收益(元)
	Yysr   float64 `json:"yysr"`   // 营业收入(百万元)
	Mgzysr float64 `json:"mgzysr"` // 每股主营业务收入(元)
	Y      int     `json:"y"`      // 报告年份
	Q      int     `json:"q"`      // 报告季度
	Yq     string  `json:"yq"`     // 报告期描述
}

// OperationData 运营能力数据
type OperationData struct {
	Dm       string  `json:"dm"`       // 代码
	Mc       string  `json:"mc"`       // 名称
	Yszzzl   float64 `json:"yszzzl"`   // 应收账款周转率(次)
	Yszts    float64 `json:"yszts"`    // 应收账款周转天数(天)
	Chzzl    float64 `json:"chzzl"`    // 存货周转率(次)
	Chzzts   float64 `json:"chzzts"`   // 存货周转天数(天)
	Ldzczzl  float64 `json:"ldzczzl"`  // 流动资产周转率(次)
	Ldzczzts float64 `json:"ldzczzts"` // 流动资产周转天数(天)
	Y        int     `json:"y"`        // 报告年份
	Q        int     `json:"q"`        // 报告季度
	Yq       string  `json:"yq"`       // 报告期描述
}

// GrowthData 成长能力数据
type GrowthData struct {
	Dm     string  `json:"dm"`     // 代码
	Mc     string  `json:"mc"`     // 名称
	Zyzzl  float64 `json:"zyzzl"`  // 主营业务收入增长率(%)
	Jlrzzl float64 `json:"jlrzzl"` // 净利润增长率(%)
	Jzczzl float64 `json:"jzczzl"` // 净资产增长率(%)
	Zzczzl float64 `json:"zzczzl"` // 总资产增长率(%)
	Mgzzl  float64 `json:"mgzzl"`  // 每股收益增长率(%)
	Gdzzl  float64 `json:"gdzzl"`  // 股东权益增长率(%)
	Y      int     `json:"y"`      // 报告年份
	Q      int     `json:"q"`      // 报告季度
	Yq     string  `json:"yq"`     // 报告期描述
}

// SolvencyData 偿债能力数据
type SolvencyData struct {
	Dm    string  `json:"dm"`    // 代码
	Mc    string  `json:"mc"`    // 名称
	Ldbl  float64 `json:"ldbl"`  // 流动比率(%)
	Sdbl  float64 `json:"sdbl"`  // 速动比率(%)
	Xjbl  float64 `json:"xjbl"`  // 现金比率(%)
	Lxbs  float64 `json:"lxbs"`  // 利息支付倍数
	Gdbl  float64 `json:"gdbl"`  // 股东权益比率(%)
	Zcfzl float64 `json:"zcfzl"` // 资产负债率(%)
	Y     int     `json:"y"`     // 报告年份
	Q     int     `json:"q"`     // 报告季度
	Yq    string  `json:"yq"`    // 报告期描述
}

// CashFlowData 现金流量数据
type CashFlowData struct {
	Dm       string  `json:"dm"`       // 代码
	Mc       string  `json:"mc"`       // 名称
	Xjlxxbl  float64 `json:"xjlxxbl"`  // 经营现金净流量对销售收入比率(%)
	Xjlhbl   float64 `json:"xjlhbl"`   // 资产的经营现金流量回报率(%)
	Xjljlrbl float64 `json:"xjljlrbl"` // 经营现金净流量与净利润的比率(%)
	Xjlfzbl  float64 `json:"xjlfzbl"`  // 经营现金净流量对负债比率(%)
	Xjllbl   float64 `json:"xjllbl"`   // 现金流量比率(%)
	Y        int     `json:"y"`        // 报告年份
	Q        int     `json:"q"`        // 报告季度
	Yq       string  `json:"yq"`       // 报告期描述
}

// PerformanceData 业绩报表数据
type PerformanceData struct {
	Dm     string  `json:"dm"`     // 代码
	Mc     string  `json:"mc"`     // 名称
	Mgsy   float64 `json:"mgsy"`   // 每股收益(元)
	Mgsytb float64 `json:"mgsytb"` // 每股收益同比(%)
	Mgjz   float64 `json:"mgjz"`   // 每股净资产(元)
	Jzsy   float64 `json:"jzsy"`   // 净资产收益率(%)
	Mgxjl  float64 `json:"mgxjl"`  // 每股现金流量(元)
	Jlr    float64 `json:"jlr"`    // 净利润(万元)
	Jlrtb  float64 `json:"jlrtb"`  // 净利润同比(%)
	Fpfa   string  `json:"fpfa"`   // 分配方案
	Rdate  string  `json:"rdate"`  // 发布日期
	Y      int     `json:"y"`      // 报告年份
	Q      int     `json:"q"`      // 报告季度
	Yq     string  `json:"yq"`     // 报告期描述
}

// StockFinancialData 单个股票的完整财务数据
type StockFinancialData struct {
	Code          string             `json:"code"`          // 股票代码
	Year          int                `json:"year"`          // 年份
	Quarter       int                `json:"quarter"`       // 季度
	Profitability *ProfitabilityData `json:"profitability"` // 盈利能力数据
	Operation     *OperationData     `json:"operation"`     // 运营能力数据
	Growth        *GrowthData        `json:"growth"`        // 成长能力数据
	Solvency      *SolvencyData      `json:"solvency"`      // 偿债能力数据
	CashFlow      *CashFlowData      `json:"cashFlow"`      // 现金流量数据
	Performance   *PerformanceData   `json:"performance"`   // 业绩报表数据
}

// FormatFinancialData 格式化财务数据为字符串
func (s *StockFinancialData) FormatFinancialData() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("财务数据分析 (%d年%d季度):\n\n", s.Year, s.Quarter))

	if s.Profitability != nil {
		result.WriteString(fmt.Sprintf("盈利能力:\n"+
			"净资产收益率: %.2f%%\n"+
			"净利率: %.2f%%\n"+
			"毛利率: %.2f%%\n"+
			"净利润: %.2f百万元\n"+
			"每股收益: %.2f元\n\n",
			s.Profitability.Jzcsy,
			s.Profitability.Jll,
			s.Profitability.Mll,
			s.Profitability.Jlr,
			s.Profitability.Mgsy))
	}

	if s.Growth != nil {
		result.WriteString(fmt.Sprintf("成长能力:\n"+
			"营收增长率: %.2f%%\n"+
			"净利润增长率: %.2f%%\n"+
			"净资产增长率: %.2f%%\n\n",
			s.Growth.Zyzzl,
			s.Growth.Jlrzzl,
			s.Growth.Jzczzl))
	}

	if s.Operation != nil {
		result.WriteString(fmt.Sprintf("运营能力:\n"+
			"应收账款周转率: %.2f次\n"+
			"存货周转率: %.2f次\n"+
			"流动资产周转率: %.2f次\n\n",
			s.Operation.Yszzzl,
			s.Operation.Chzzl,
			s.Operation.Ldzczzl))
	}

	if s.Solvency != nil {
		result.WriteString(fmt.Sprintf("偿债能力:\n"+
			"流动比率: %.2f%%\n"+
			"速动比率: %.2f%%\n"+
			"资产负债率: %.2f%%\n\n",
			s.Solvency.Ldbl,
			s.Solvency.Sdbl,
			s.Solvency.Zcfzl))
	}

	if s.CashFlow != nil {
		result.WriteString(fmt.Sprintf("现金流量:\n"+
			"经营现金流量比率: %.2f%%\n"+
			"现金流量比率: %.2f%%\n\n",
			s.CashFlow.Xjlxxbl,
			s.CashFlow.Xjllbl))
	}

	if s.Performance != nil {
		result.WriteString(fmt.Sprintf("业绩表现:\n"+
			"每股收益: %.2f元 (同比%.2f%%)\n"+
			"每股净资产: %.2f元\n"+
			"净利润: %.2f万元 (同比%.2f%%)\n",
			s.Performance.Mgsy,
			s.Performance.Mgsytb,
			s.Performance.Mgjz,
			s.Performance.Jlr,
			s.Performance.Jlrtb))
	}

	return result.String()
}
