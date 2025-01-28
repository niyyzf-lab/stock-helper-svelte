package api

// StockBasic 股票基本信息
type StockBasic struct {
	Code   string `json:"code"`   // 股票代码
	Name   string `json:"name"`   // 股票名称
	Market string `json:"market"` // 市场类型（主板、创业板等）
}

// Index 指数信息
type Index struct {
	Code     string `json:"dm"`  // 指数代码
	Name     string `json:"mc"`  // 指数名称
	Exchange string `json:"jys"` // 交易所
}

// KLineData K线数据
type KLineData struct {
	Time      string  `json:"d"`   // 交易时间
	Open      float64 `json:"o"`   // 开盘价
	High      float64 `json:"h"`   // 最高价
	Low       float64 `json:"l"`   // 最低价
	Close     float64 `json:"c"`   // 收盘价
	Volume    float64 `json:"v"`   // 成交量(手)
	Amount    float64 `json:"e"`   // 成交额(元)
	Amplitude float64 `json:"zf"`  // 振幅(%)
	Turnover  float64 `json:"hs"`  // 换手率(%)
	Change    float64 `json:"zd"`  // 涨跌幅(%)
	ChangeAmt float64 `json:"zde"` // 涨跌额(元)
}

// KLineFreq K线周期
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

// RealtimeData 实时交易数据
type RealtimeData struct {
	FiveMinChange float64 `json:"fm"`    // 五分钟涨跌幅（%）
	High          float64 `json:"h"`     // 最高价（元）
	Turnover      float64 `json:"hs"`    // 换手（%）
	VolumeRatio   float64 `json:"lb"`    // 量比（%）
	Low           float64 `json:"l"`     // 最低价（元）
	FloatValue    float64 `json:"lt"`    // 流通市值（元）
	Open          float64 `json:"o"`     // 开盘价（元）
	PE            float64 `json:"pe"`    // 市盈率
	ChangePercent float64 `json:"pc"`    // 涨跌幅（%）
	Price         float64 `json:"p"`     // 当前价格（元）
	TotalValue    float64 `json:"sz"`    // 总市值（元）
	Amount        float64 `json:"cje"`   // 成交额（元）
	Change        float64 `json:"ud"`    // 涨跌额（元）
	Volume        float64 `json:"v"`     // 成交量（手）
	PreClose      float64 `json:"yc"`    // 昨日收盘价（元）
	Amplitude     float64 `json:"zf"`    // 振幅（%）
	Speed         float64 `json:"zs"`    // 涨速（%）
	PB            float64 `json:"sjl"`   // 市净率
	Change60Day   float64 `json:"zdf60"` // 60日涨跌幅（%）
	ChangeYTD     float64 `json:"zdfnc"` // 年初至今涨跌幅（%）
	Time          string  `json:"t"`     // 更新时间yyyy-MM-dd HH:mm:ss
}

// HistoricalTransaction 历史成交分布数据
type HistoricalTransaction struct {
	Time            string  `json:"t"`      // 时间yyyy-MM-dd
	Close           float64 `json:"c"`      // 收盘价(元)
	ChangePercent   float64 `json:"zdf"`    // 涨跌幅(%)
	NetInflowRate   float64 `json:"jlrl"`   // 净流入率(%)
	TurnoverRate    float64 `json:"hsl"`    // 换手率(%)
	TotalNetInflow  float64 `json:"qbjlr"`  // 全部净流入(元)
	SuperInflow     float64 `json:"cddlr"`  // 超大单流入(元)
	SuperNetInflow  float64 `json:"cddjlr"` // 超大单净流入(元)
	LargeInflow     float64 `json:"ddlr"`   // 大单流入(元)
	LargeNetInflow  float64 `json:"ddjlr"`  // 大单净流入(元)
	SmallInflow     float64 `json:"xdlr"`   // 小单流入(元)
	SmallNetInflow  float64 `json:"xdjlr"`  // 小单净流入(元)
	RetailInflow    float64 `json:"sdlr"`   // 散单流入(元)
	RetailNetInflow float64 `json:"sdjlr"`  // 散单净流入(元)
}
