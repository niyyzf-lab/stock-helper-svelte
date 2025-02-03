package indicators

// TrendStrength 计算数列的趋势强度
// values: 输入数列
// period: 分析周期 (建议值：14)
// 返回值范围: [-1, 1]
// -1 表示极强下跌趋势
// 1 表示极强上涨趋势
// 0 表示无明显趋势
func TrendStrength(values []float64, period int) ([]float64, error) {
	// 输入参数验证
	if values == nil {
		return make([]float64, 0), nil
	}
	if len(values) == 0 {
		return make([]float64, 0), nil
	}

	// 使用默认参数
	if period <= 0 {
		period = 14
	}

	// 检查数组中是否存在无效值
	for i, value := range values {
		if !IsValidFloat(value) {
			values[i] = 0
		}
	}

	result := make([]float64, len(values))

	// 填充前期未计算的值为0
	for i := 0; i < period-1; i++ {
		result[i] = 0
	}

	// 计算相关指标
	// 1. 计算移动平均线
	sma, _ := CalculateMA(values, SMA, period)

	// 2. 计算价格动量
	momentum := calculateMomentum(values, period)

	// 3. 计算线性回归斜率
	slopes := calculateLinearRegressionSlope(values, period)

	// 4. 计算波动率
	volatility := calculateVolatility(values, period)

	// 综合计算趋势强度
	for i := period - 1; i < len(values); i++ {
		// 1. 价格相对均线位置 (-0.25 到 0.25)
		pricePosition := (values[i] - sma[i]) / sma[i]
		priceScore := Clamp(pricePosition, -0.25, 0.25)

		// 2. 动量得分 (-0.25 到 0.25)
		momentumScore := Clamp(momentum[i]/100, -0.25, 0.25)

		// 3. 斜率得分 (-0.25 到 0.25)
		slopeScore := Clamp(slopes[i]*5, -0.25, 0.25)

		// 4. 波动率权重 (0.5 到 1.0)
		volatilityWeight := 0.5 + Clamp(volatility[i], 0, 0.5)

		// 综合计算
		trendScore := (priceScore + momentumScore + slopeScore) * volatilityWeight

		// 确保结果在 [-1, 1] 范围内
		result[i] = Clamp(trendScore, -1.0, 1.0)
	}

	return result, nil
}

// calculateMomentum 计算动量
func calculateMomentum(values []float64, period int) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < period-1; i++ {
		result[i] = 0
	}

	for i := period - 1; i < len(values); i++ {
		if values[i-period+1] != 0 {
			result[i] = ((values[i] / values[i-period+1]) - 1) * 100
		} else {
			result[i] = 0
		}
	}

	return result
}

// calculateLinearRegressionSlope 计算线性回归斜率
func calculateLinearRegressionSlope(values []float64, period int) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < period-1; i++ {
		result[i] = 0
	}

	for i := period - 1; i < len(values); i++ {
		sumX := 0.0
		sumY := 0.0
		sumXY := 0.0
		sumX2 := 0.0

		for j := 0; j < period; j++ {
			x := float64(j)
			y := values[i-period+1+j]

			sumX += x
			sumY += y
			sumXY += x * y
			sumX2 += x * x
		}

		// 计算斜率
		n := float64(period)
		denominator := n*sumX2 - sumX*sumX
		if denominator != 0 {
			slope := (n*sumXY - sumX*sumY) / denominator
			result[i] = slope
		} else {
			result[i] = 0
		}
	}

	return result
}

// calculateVolatility 计算波动率
func calculateVolatility(values []float64, period int) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < period-1; i++ {
		result[i] = 0
	}

	for i := period - 1; i < len(values); i++ {
		sum := 0.0
		mean := 0.0

		// 计算均值
		for j := 0; j < period; j++ {
			mean += values[i-j]
		}
		mean /= float64(period)

		if mean != 0 {
			// 计算标准差
			for j := 0; j < period; j++ {
				diff := values[i-j] - mean
				sum += diff * diff
			}

			// 计算波动率
			result[i] = Clamp(Sqrt(sum/float64(period))/mean, 0, 1)
		} else {
			result[i] = 0
		}
	}

	return result
}
