package indicators

// CalculateRSI 计算相对强弱指标(RSI)
// prices: 收盘价数组
// period: 计算周期（常用值：6、12、24）
func CalculateRSI(prices []float64, period int) ([]float64, error) {
	// 输入参数验证
	if prices == nil {
		return make([]float64, 0), nil
	}
	if len(prices) == 0 {
		return make([]float64, 0), nil
	}

	// 使用默认参数
	if period <= 0 {
		period = 14
	}

	// 检查价格数组中是否存在无效值
	for i, price := range prices {
		if !IsValidFloat(price) {
			prices[i] = 0
		}
	}

	size := len(prices)
	rsi := make([]float64, size)
	gains := make([]float64, size)
	losses := make([]float64, size)

	// 填充前期未计算的值为50
	for i := 0; i < period; i++ {
		rsi[i] = 50
	}

	// 计算涨跌幅
	for i := 1; i < size; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains[i] = change
			losses[i] = 0
		} else {
			gains[i] = 0
			losses[i] = -change
		}
	}

	// 计算首个RSI值
	var sumGain, sumLoss float64
	for i := 1; i <= period; i++ {
		sumGain += gains[i]
		sumLoss += losses[i]
	}

	// 使用Wilder's Smoothing Method
	avgGain := sumGain / float64(period)
	avgLoss := sumLoss / float64(period)

	if avgLoss == 0 {
		rsi[period] = 100
	} else {
		rs := avgGain / avgLoss
		rsi[period] = 100 - (100 / (1 + rs))
	}

	// 计算后续的RSI值
	for i := period + 1; i < size; i++ {
		// 更新平均涨跌幅
		avgGain = ((avgGain * float64(period-1)) + gains[i]) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + losses[i]) / float64(period)

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}

		// 确保RSI值在0-100范围内
		rsi[i] = Clamp(rsi[i], 0, 100)
	}

	return rsi, nil
}
