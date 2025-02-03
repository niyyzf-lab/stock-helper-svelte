package types

// Index 股票指数信息
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

// MainForceMinute 主力资金分钟走势
type MainForceMinute struct {
	Time             string  `json:"t"`     // 时间yyyy-MM-dd HH:mm:ss
	ChangePercent    float64 `json:"zdf"`   // 涨跌幅(%)
	MainForceInflow  float64 `json:"lrzj"`  // 主力流入(元)
	InflowRate       float64 `json:"lrl"`   // 主力流入率(%)
	MainForceOutflow float64 `json:"lczj"`  // 主力流出(元)
	NetInflow        float64 `json:"jlr"`   // 主力净流入(元)
	NetInflowRate    float64 `json:"jlrl"`  // 主力净流入率(%)
	RetailInflowRate float64 `json:"shlrl"` // 散户流入率(%)
}

// CapitalFlow 资金流向趋势
type CapitalFlow struct {
	Time                string  `json:"t"`      // 时间yyyy-MM-dd
	ChangePercent       float64 `json:"zdf"`    // 涨跌幅(%)
	NetInflow           float64 `json:"jlr"`    // 净流入(元)
	NetInflowRate       float64 `json:"jlrl"`   // 净流入率(%)
	MainForceNetInflow  float64 `json:"zljlr"`  // 主力净流入(元)
	MainForceInflowRate float64 `json:"zljlrl"` // 主力净流入率(%)
	IndustryNetInflow   float64 `json:"hyjlr"`  // 行业净流入(元)
	IndustryInflowRate  float64 `json:"hyjlrl"` // 行业净流入率(%)
	TurnoverRate        float64 `json:"hsl"`    // 换手率(%)
}

// MainForcePhase 阶段主力动向
type MainForcePhase struct {
	Time           string  `json:"t"`      // 时间yyyy-MM-dd
	NetInflow3Day  float64 `json:"jlr3"`   // 近3日主力净流入(元)
	NetInflow5Day  float64 `json:"jlr5"`   // 近5日主力净流入(元)
	NetInflow10Day float64 `json:"jlr10"`  // 近10日主力净流入(元)
	Rate3Day       float64 `json:"jlrl3"`  // 近3日主力净流入率(%)
	Rate5Day       float64 `json:"jlrl5"`  // 近5日主力净流入率(%)
	Rate10Day      float64 `json:"jlrl10"` // 近10日主力净流入率(%)
}
