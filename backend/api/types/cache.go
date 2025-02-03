package types

import "time"

// KLineFreq K线周期类型
type KLineFreq string

const (
	FREQ_5MIN        KLineFreq = "5m"  // 5分钟
	FREQ_15MIN       KLineFreq = "15m" // 15分钟
	FREQ_30MIN       KLineFreq = "30m" // 30分钟
	FREQ_60MIN       KLineFreq = "60m" // 60分钟
	FREQ_DAILY       KLineFreq = "dn"  // 日线(未复权)
	FREQ_DAILY_QFQ   KLineFreq = "dq"  // 日线(前复权)
	FREQ_DAILY_HFQ   KLineFreq = "dh"  // 日线(后复权)
	FREQ_WEEKLY      KLineFreq = "wn"  // 周线(未复权)
	FREQ_WEEKLY_QFQ  KLineFreq = "wq"  // 周线(前复权)
	FREQ_WEEKLY_HFQ  KLineFreq = "wh"  // 周线(后复权)
	FREQ_MONTHLY     KLineFreq = "mn"  // 月线(未复权)
	FREQ_MONTHLY_QFQ KLineFreq = "mq"  // 月线(前复权)
	FREQ_MONTHLY_HFQ KLineFreq = "mh"  // 月线(后复权)
	FREQ_YEARLY      KLineFreq = "yn"  // 年线(未复权)
	FREQ_YEARLY_QFQ  KLineFreq = "yq"  // 年线(前复权)
	FREQ_YEARLY_HFQ  KLineFreq = "yh"  // 年线(后复权)
)

// GetNextUpdateTime 获取下次更新时间(针对日线及以上级别)
func GetNextUpdateTime() time.Time {
	now := time.Now()
	today4PM := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, now.Location())

	if now.Before(today4PM) {
		return today4PM
	}

	tomorrow4PM := today4PM.Add(24 * time.Hour)
	return tomorrow4PM
}

// GetTTL 获取到下次更新时间的剩余秒数
func GetTTL() int64 {
	nextUpdate := GetNextUpdateTime()
	return int64(time.Until(nextUpdate).Seconds())
}

// 缓存时间常量(秒)
const (
	// 分钟级别K线缓存时间
	CacheTime_KLine_5Min  = 300  // 5分钟K线缓存5分钟
	CacheTime_KLine_15Min = 900  // 15分钟K线缓存15分钟
	CacheTime_KLine_30Min = 1800 // 30分钟K线缓存30分钟
	CacheTime_KLine_60Min = 3600 // 60分钟K线缓存1小时

	// 实时数据缓存时间
	CacheTime_RealtimeData = 60 // 实时数据缓存1分钟

	// 日线及以上级别数据缓存到下一个交易日16:00
	CacheTime_Dynamic = -1 // 标记需要动态计算TTL的缓存类型
)

// GetKLineCacheTime 根据K线频率获取缓存时间
func GetKLineCacheTime(freq KLineFreq) int {
	switch freq {
	case FREQ_5MIN:
		return CacheTime_KLine_5Min
	case FREQ_15MIN:
		return CacheTime_KLine_15Min
	case FREQ_30MIN:
		return CacheTime_KLine_30Min
	case FREQ_60MIN:
		return CacheTime_KLine_60Min
	case FREQ_DAILY, FREQ_DAILY_QFQ, FREQ_DAILY_HFQ,
		FREQ_WEEKLY, FREQ_WEEKLY_QFQ, FREQ_WEEKLY_HFQ,
		FREQ_MONTHLY, FREQ_MONTHLY_QFQ, FREQ_MONTHLY_HFQ,
		FREQ_YEARLY, FREQ_YEARLY_QFQ, FREQ_YEARLY_HFQ:
		return CacheTime_Dynamic
	default:
		return CacheTime_Dynamic
	}
}

// IsDynamicCache 判断是否是需要动态计算TTL的缓存类型
func IsDynamicCache(cacheTime int) bool {
	return cacheTime == CacheTime_Dynamic
}
