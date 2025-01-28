package types

// StockSignal 股票信号
type StockSignal struct {
	Code     string  `json:"code"`     // 股票代码
	Name     string  `json:"name"`     // 股票名称
	Price    float64 `json:"price"`    // 当前价格
	Change   float64 `json:"change"`   // 涨跌幅
	Turnover float64 `json:"turnover"` // 换手率
	Reason   string  `json:"reason"`   // 选股原因
}
