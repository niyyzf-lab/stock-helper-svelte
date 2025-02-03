package indicators

// MACDResult MACD 计算结果
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
	// 输入参数验证
	if prices == nil {
		return &MACDResult{
			DIF:  make([]float64, 0),
			DEA:  make([]float64, 0),
			MACD: make([]float64, 0),
		}, nil
	}
	if len(prices) == 0 {
		return &MACDResult{
			DIF:  make([]float64, 0),
			DEA:  make([]float64, 0),
			MACD: make([]float64, 0),
		}, nil
	}

	// 使用默认参数
	if shortPeriod <= 0 {
		shortPeriod = 12
	}
	if longPeriod <= 0 {
		longPeriod = 26
	}
	if signalPeriod <= 0 {
		signalPeriod = 9
	}

	// 如果周期参数不合理，返回全0数组
	if shortPeriod >= longPeriod || longPeriod > len(prices) {
		result := &MACDResult{
			DIF:  make([]float64, len(prices)),
			DEA:  make([]float64, len(prices)),
			MACD: make([]float64, len(prices)),
		}
		return result, nil
	}

	// 检查价格数组中是否存在无效值
	for i, price := range prices {
		if !IsValidFloat(price) {
			prices[i] = 0
		}
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
		if i < longPeriod+signalPeriod-2 {
			macd[i] = 0
		} else {
			macd[i] = (dif[i] - dea[i]) * 2
		}
	}

	return &MACDResult{
		DIF:  dif,
		DEA:  dea,
		MACD: macd,
	}, nil
}
