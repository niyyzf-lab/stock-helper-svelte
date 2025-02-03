package indicators

// MAType 定义了移动平均线的类型
type MAType string

const (
	SMA MAType = "sma" // 简单移动平均线
	EMA MAType = "ema" // 指数移动平均线
	WMA MAType = "wma" // 加权移动平均线
	TMA MAType = "tma" // 三重移动平均线
)

// CalculateMA 计算移动平均线
func CalculateMA(prices []float64, maType MAType, period int) ([]float64, error) {
	if prices == nil {
		return make([]float64, 0), nil
	}
	if len(prices) == 0 {
		return make([]float64, 0), nil
	}
	if period <= 0 {
		return make([]float64, len(prices)), nil
	}
	if period > len(prices) {
		return make([]float64, len(prices)), nil
	}

	// 检查价格数组中是否存在无效值
	for i, price := range prices {
		if !IsValidFloat(price) {
			prices[i] = 0
		}
	}

	switch maType {
	case SMA:
		return calculateSMA(prices, period), nil
	case EMA:
		return calculateEMA(prices, period), nil
	case WMA:
		return calculateWMA(prices, period), nil
	case TMA:
		return calculateTMA(prices, period), nil
	default:
		return make([]float64, len(prices)), nil
	}
}

// calculateSMA 计算简单移动平均线
func calculateSMA(prices []float64, period int) []float64 {
	result := make([]float64, len(prices))

	// 填充无法计算的位置为0
	for i := 0; i < period-1; i++ {
		result[i] = 0
	}

	// 计算第一个SMA
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	result[period-1] = sum / float64(period)

	// 计算后续的SMA
	for i := period; i < len(prices); i++ {
		sum = sum - prices[i-period] + prices[i]
		result[i] = sum / float64(period)
	}

	return result
}

// calculateEMA 计算指数移动平均线
func calculateEMA(prices []float64, period int) []float64 {
	result := make([]float64, len(prices))

	// 填充无法计算的位置为0
	for i := 0; i < period-1; i++ {
		result[i] = 0
	}

	// 使用SMA作为第一个EMA值
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	result[period-1] = sum / float64(period)

	// 计算后续的EMA
	multiplier := 2.0 / float64(period+1)
	for i := period; i < len(prices); i++ {
		// 使用更稳定的计算方法
		result[i] = prices[i]*multiplier + result[i-1]*(1-multiplier)
	}

	return result
}

// calculateWMA 计算加权移动平均线
func calculateWMA(prices []float64, period int) []float64 {
	result := make([]float64, len(prices))

	// 填充无法计算的位置为0
	for i := 0; i < period-1; i++ {
		result[i] = 0
	}

	// 计算权重之和
	weightSum := float64(period * (period + 1) / 2)

	// 计算WMA
	for i := period - 1; i < len(prices); i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			weight := float64(period - j)
			sum += prices[i-j] * weight
		}
		result[i] = sum / weightSum
	}

	return result
}

// calculateTMA 计算三重移动平均线
func calculateTMA(prices []float64, period int) []float64 {
	result := make([]float64, len(prices))

	// 计算三重SMA
	sma1 := calculateSMA(prices, period)
	sma2 := calculateSMA(sma1, period)
	sma3 := calculateSMA(sma2, period)

	copy(result, sma3)
	return result
}
