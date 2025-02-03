package indicators

// IsValidFloat 检查浮点数是否有效
func IsValidFloat(value float64) bool {
	// 检查是否为无限大或无效值
	if value > 1.797693134862315708145274237317043567981e+308 ||
		value < -1.797693134862315708145274237317043567981e+308 {
		return false
	}
	return true
}

// Clamp 将值限制在指定范围内
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Sqrt 安全的平方根计算
func Sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	return x * x
}
