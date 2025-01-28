package indicators

import (
	"fmt"
)

// MACD 计算结果
type MACDResult struct {
	DIF  []float64 // 差离值
	DEA  []float64 // 信号线
	MACD []float64 // MACD柱
}

// CalculateMACD 计算MACD指标
// prices: 收盘价数组
// shortPeriod: 快线周期(默认12)
// longPeriod: 慢线周期(默认26)
// signalPeriod: 信号线周期(默认9)
func CalculateMACD(prices []float64, shortPeriod, longPeriod, signalPeriod int) (*MACDResult, error) {
	if len(prices) == 0 {
		return nil, fmt.Errorf("价格数组不能为空")
	}

	// 计算快线EMA
	shortEMA := calculateEMA(prices, shortPeriod)
	// 计算慢线EMA
	longEMA := calculateEMA(prices, longPeriod)

	// 计算DIF
	dif := make([]float64, len(prices))
	for i := range prices {
		if i < longPeriod-1 {
			dif[i] = 0
		} else {
			dif[i] = shortEMA[i] - longEMA[i]
		}
	}

	// 计算DEA(信号线)
	dea := calculateEMA(dif, signalPeriod)

	// 计算MACD柱
	macd := make([]float64, len(prices))
	for i := range prices {
		macd[i] = (dif[i] - dea[i]) * 2
	}

	return &MACDResult{
		DIF:  dif,
		DEA:  dea,
		MACD: macd,
	}, nil
}
