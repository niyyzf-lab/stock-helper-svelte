package indicators

// KDJResult KDJ 计算结果
type KDJResult struct {
	K []float64
	D []float64
	J []float64
}

// CalculateKDJ 计算KDJ指标
// prices: 收盘价数组
// n: RSV周期（默认9）
// m1: K值平滑系数（默认3）
// m2: D值平滑系数（默认3）
func CalculateKDJ(prices []float64, n, m1, m2 int) (*KDJResult, error) {
	// 输入参数验证
	if prices == nil {
		return &KDJResult{
			K: make([]float64, 0),
			D: make([]float64, 0),
			J: make([]float64, 0),
		}, nil
	}
	if len(prices) == 0 {
		return &KDJResult{
			K: make([]float64, 0),
			D: make([]float64, 0),
			J: make([]float64, 0),
		}, nil
	}

	// 使用默认参数
	if n <= 0 {
		n = 9
	}
	if m1 <= 0 {
		m1 = 3
	}
	if m2 <= 0 {
		m2 = 3
	}

	// 检查价格数组中是否存在无效值
	for i, price := range prices {
		if !IsValidFloat(price) {
			prices[i] = 0
		}
	}

	// 计算RSV
	rsv := make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		if i < n-1 {
			rsv[i] = 50 // 将初始值设为50
			continue
		}
		high := prices[i]
		low := prices[i]
		for j := 0; j < n; j++ {
			if prices[i-j] > high {
				high = prices[i-j]
			}
			if prices[i-j] < low {
				low = prices[i-j]
			}
		}
		if high == low {
			rsv[i] = 50
		} else {
			rsv[i] = (prices[i] - low) / (high - low) * 100
		}
	}

	// 计算KDJ
	k := make([]float64, len(prices))
	d := make([]float64, len(prices))
	j := make([]float64, len(prices))

	// 初始值
	k[0] = 50
	d[0] = 50
	j[0] = 50

	// 计算其他值
	for i := 1; i < len(prices); i++ {
		// 使用SMA计算K值和D值
		k[i] = (float64(m1-1)*k[i-1] + rsv[i]) / float64(m1)
		d[i] = (float64(m2-1)*d[i-1] + k[i]) / float64(m2)
		j[i] = 3*k[i] - 2*d[i]

		// 确保值在0-100范围内
		k[i] = Clamp(k[i], 0, 100)
		d[i] = Clamp(d[i], 0, 100)
		j[i] = Clamp(j[i], 0, 100)
	}

	return &KDJResult{K: k, D: d, J: j}, nil
}
