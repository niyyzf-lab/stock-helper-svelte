package types

// CompanyInfo 公司基本信息
type CompanyInfo struct {
	Code                string `json:"dm"`        // 股票代码
	Name                string `json:"mc"`        // 公司名称
	Industry            string `json:"hy"`        // 所属行业
	MainBusiness        string `json:"zygw"`      // 主营业务
	RegisteredCapital   string `json:"zczb"`      // 注册资本(万元)
	TotalShares         string `json:"zgb"`       // 总股本(万股)
	FloatShares         string `json:"ltgb"`      // 流通股本(万股)
	Area                string `json:"qy"`        // 地区
	ListDate            string `json:"ssrq"`      // 上市日期
	EstablishDate       string `json:"clrq"`      // 成立日期
	Exchange            string `json:"jys"`       // 交易所
	Chairman            string `json:"dsz"`       // 董事长
	LegalRepresentative string `json:"frdb"`      // 法人代表
	Secretary           string `json:"dm1"`       // 董秘
	Website             string `json:"gw"`        // 公司网站
	Email               string `json:"dzyj"`      // 电子邮件
	RegisteredAddress   string `json:"zcdz"`      // 注册地址
	OfficeAddress       string `json:"bgdz"`      // 办公地址
	Profile             string `json:"gsjj"`      // 公司简介
	EmployeeCount       string `json:"ygrs"`      // 员工人数
	EnglishName         string `json:"ename"`     // 公司英文名称
	Market              string `json:"market"`    // 上市市场
	Concepts            string `json:"idea"`      // 概念及板块
	IssuePrice          string `json:"sprice"`    // 发行价格(元)
	Underwriter         string `json:"principal"` // 主承销商
	InstitutionType     string `json:"instype"`   // 机构类型7
	OrganizationType    string `json:"organ"`     // 组织形式
	SecretaryPhone      string `json:"sphone"`    // 董秘电话
	SecretaryFax        string `json:"sfax"`      // 董秘传真
	SecretaryEmail      string `json:"semail"`    // 董秘电子邮箱
	Phone               string `json:"phone"`     // 公司电话
	Fax                 string `json:"fax"`       // 公司传真
	PostalCode          string `json:"post"`      // 邮政编码
	InfoDisclosureUrl   string `json:"infosite"`  // 信息披露网址
	NameHistory         string `json:"oname"`     // 证券简称更名历史
	BusinessScope       string `json:"bscope"`    // 经营范围
}

// FinancialIndicator 主要财务指标
type FinancialIndicator struct {
	Code            string `json:"dm"`      // 股票代码
	Period          string `json:"date"`    // 报告期
	Revenue         string `json:"yysr"`    // 营业收入(元)
	RevenueYOY      string `json:"yysrzzl"` // 营业收入同比增长率(%)
	NetProfit       string `json:"jlr"`     // 净利润(元)
	NetProfitYOY    string `json:"jlrzzl"`  // 净利润同比增长率(%)
	GrossMargin     string `json:"mlr"`     // 毛利率(%)
	NetProfitMargin string `json:"jll"`     // 净利率(%)
	ROE             string `json:"jzcsyl"`  // 净资产收益率(%)
	DebtRatio       string `json:"zcfzl"`   // 资产负债率(%)
	EPS             string `json:"mgsy"`    // 每股收益(元)
	NetAssetPS      string `json:"mgjzc"`   // 每股净资产(元)
}

// ShareholderInfo 股东信息
type ShareholderInfo struct {
	Code        string `json:"dm"`     // 股票代码
	Date        string `json:"date"`   // 截止日期
	Name        string `json:"gdmc"`   // 股东名称
	SharesHeld  string `json:"cgs"`    // 持股数(万股)
	Percentage  string `json:"zltgbl"` // 占流通股比例(%)
	Nature      string `json:"gdxz"`   // 股东性质
	ShareChange string `json:"zde"`    // 较上期变动(万股)
	ChangeRatio string `json:"zdbl"`   // 变动比例(%)
}

// ExecutiveInfo 高管信息
type ExecutiveInfo struct {
	Code      string `json:"dm"`   // 股票代码
	Name      string `json:"xm"`   // 姓名
	Position  string `json:"zw"`   // 职务
	StartDate string `json:"rzsj"` // 任职时间
	Education string `json:"xl"`   // 学历
	Birth     string `json:"csny"` // 出生年月
	Salary    string `json:"xz"`   // 年薪(万元)
	Shares    string `json:"cgs"`  // 持股数(万股)
}

// CompanyNews 公司新闻
type CompanyNews struct {
	Code    string `json:"dm"`   // 股票代码
	Title   string `json:"bt"`   // 标题
	Content string `json:"nr"`   // 内容
	Date    string `json:"date"` // 发布日期
	Source  string `json:"ly"`   // 来源
	URL     string `json:"url"`  // 新闻链接
}

// CompanyAnnouncement 公司公告
type CompanyAnnouncement struct {
	Code  string `json:"dm"`   // 股票代码
	Title string `json:"bt"`   // 公告标题
	Type  string `json:"lx"`   // 公告类型
	Date  string `json:"date"` // 公告日期
	URL   string `json:"url"`  // 公告链接
}

// DividendInfo 分红信息
type DividendInfo struct {
	Code           string `json:"dm"`     // 股票代码
	AnnouncedDate  string `json:"sdate"`  // 公告日期
	GiveStocks     string `json:"give"`   // 每10股送股
	TransferStocks string `json:"change"` // 每10股转增
	Dividend       string `json:"send"`   // 每10股派息(税前)
	Progress       string `json:"line"`   // 进度
	ExRightDate    string `json:"cdate"`  // 除权除息日
	RegisterDate   string `json:"edate"`  // 股权登记日
	ListingDate    string `json:"hdate"`  // 红股上市日
}

// ShareUnlockInfo 解禁限售信息
type ShareUnlockInfo struct {
	Code          string `json:"dm"`      // 股票代码
	UnlockDate    string `json:"rdate"`   // 解禁日期
	UnlockShares  string `json:"ramount"` // 解禁数量(万股)
	MarketValue   string `json:"rprice"`  // 解禁股流通市值(亿元)
	Batch         string `json:"batch"`   // 上市批次
	AnnouncedDate string `json:"pdate"`   // 公告日期
}

// QuarterlyProfit 季度利润
type QuarterlyProfit struct {
	Code            string `json:"dm"`       // 股票代码
	EndDate         string `json:"date"`     // 截止日期
	Revenue         string `json:"income"`   // 营业收入
	Expenses        string `json:"expend"`   // 营业支出
	OperatingProfit string `json:"profit"`   // 营业利润
	TotalProfit     string `json:"totalp"`   // 利润总额
	NetProfit       string `json:"reprofit"` // 净利润
	BasicEPS        string `json:"basege"`   // 基本每股收益
	DilutedEPS      string `json:"ettege"`   // 稀释每股收益
}

// CompanyProfile 公司详细信息
type CompanyProfile struct {
	Name          string `json:"name"`      // 公司名称
	EnglishName   string `json:"ename"`     // 公司英文名称
	Market        string `json:"market"`    // 上市市场
	Concepts      string `json:"idea"`      // 概念及板块
	Description   string `json:"desc"`      // 公司简介
	BusinessScope string `json:"bscope"`    // 经营范围
	ListDate      string `json:"ldate"`     // 上市日期
	IPOPrice      string `json:"sprice"`    // 发行价格
	Underwriter   string `json:"principal"` // 主承销商
}

// BelongingIndex 所属指数信息
type BelongingIndex struct {
	Name      string `json:"mc"`   // 指数名称
	Code      string `json:"dm"`   // 指数代码
	JoinDate  string `json:"ind"`  // 进入日期
	LeaveDate string `json:"outd"` // 退出日期
}

// ShareholderTrend 股东变化趋势
type ShareholderTrend struct {
	Date   string `json:"jzrq"` // 截止日期
	Count  string `json:"gdhs"` // 股东户数
	Change string `json:"bh"`   // 比上期变化情况
}

// FundHolding 基金持股
type FundHolding struct {
	Date          string  `json:"jzrq"` // 截止日期
	FundName      string  `json:"jjmc"` // 基金名称
	FundCode      string  `json:"jjdm"` // 基金代码
	SharesHeld    float64 `json:"ccsl"` // 持仓数量
	FloatRatio    float64 `json:"ltbl"` // 占流通股比例
	MarketValue   float64 `json:"cgsz"` // 持股市值
	NetWorthRatio float64 `json:"jzbl"` // 占净值比例
}

// CashFlow 季度现金流
type CashFlow struct {
	Date           string `json:"date"`    // 截止日期
	OperateInflow  string `json:"jyin"`    // 经营活动现金流入
	OperateOutflow string `json:"jyout"`   // 经营活动现金流出
	OperateNet     string `json:"jyfinal"` // 经营活动净流量
	InvestInflow   string `json:"tzin"`    // 投资活动现金流入
	InvestOutflow  string `json:"tzout"`   // 投资活动现金流出
	InvestNet      string `json:"tzfinal"` // 投资活动净流量
	CashIncrease   string `json:"cashinc"` // 现金净增加额
}

// PerformanceForecast 业绩预告
type PerformanceForecast struct {
	AnnounceDate string `json:"pdate"` // 公告日期
	ReportDate   string `json:"rdate"` // 报告期
	Type         string `json:"type"`  // 类型
	Summary      string `json:"abs"`   // 业绩预告摘要
	LastYearEPS  string `json:"old"`   // 上年同期每股收益
}

// ShareUnlock 解禁信息
type ShareUnlock struct {
	UnlockDate   string  `json:"rdate"`   // 解禁日期
	UnlockShares float64 `json:"ramount"` // 解禁数量
	MarketValue  float64 `json:"rprice"`  // 解禁市值
	Batch        float64 `json:"batch"`   // 上市批次
	AnnounceDate string  `json:"pdate"`   // 公告日期
}

// Dividend 分红信息
type Dividend struct {
	AnnounceDate   string `json:"sdate"`  // 公告日期
	GiveStocks     string `json:"give"`   // 每10股送股
	TransferStocks string `json:"change"` // 每10股转增
	DividendAmount string `json:"send"`   // 每10股派息
	Progress       string `json:"line"`   // 进度
	ExRightDate    string `json:"cdate"`  // 除权除息日
}

// Executive 高管成员
type Executive struct {
	Name      string `json:"name"`  // 姓名
	Title     string `json:"title"` // 职务
	StartDate string `json:"sdate"` // 起始日期
	EndDate   string `json:"edate"` // 终止日期
}

// Director 董事会成员
type Director struct {
	Name      string `json:"name"`  // 姓名
	Title     string `json:"title"` // 职务
	StartDate string `json:"sdate"` // 起始日期
	EndDate   string `json:"edate"` // 终止日期
}

// Supervisor 监事会成员
type Supervisor struct {
	Name      string `json:"name"`  // 姓名
	Title     string `json:"title"` // 职务
	StartDate string `json:"sdate"` // 起始日期
	EndDate   string `json:"edate"` // 终止日期
}

// AdditionalIssue 增发信息
type AdditionalIssue struct {
	AnnounceDate string `json:"sdate"`  // 公告日期
	IssueType    string `json:"type"`   // 发行方式
	IssuePrice   string `json:"price"`  // 发行价格
	TotalAmount  string `json:"tprice"` // 实际募集资金总额
	IssueFee     string `json:"fprice"` // 发行费用总额
	IssueShares  string `json:"amount"` // 实际发行数量
}

// QuarterlyCashFlow 季度现金流
type QuarterlyCashFlow struct {
	Date           string `json:"date"`    // 截止日期
	OperateInflow  string `json:"jyin"`    // 经营活动现金流入小计
	OperateOutflow string `json:"jyout"`   // 经营活动现金流出小计
	OperateNet     string `json:"jyfinal"` // 经营活动产生的现金流量净额
	InvestInflow   string `json:"tzin"`    // 投资活动现金流入小计
	InvestOutflow  string `json:"tzout"`   // 投资活动现金流出小计
	InvestNet      string `json:"tzfinal"` // 投资活动产生的现金流量净额
	FinanceInflow  string `json:"czin"`    // 筹资活动现金流入小计
	FinanceOutflow string `json:"czout"`   // 筹资活动现金流出小计
	FinanceNet     string `json:"czfinal"` // 筹资活动产生的现金流量净额
	ExchangeEffect string `json:"hl"`      // 汇率变动对现金的影响
	CashIncrease   string `json:"cashinc"` // 现金及现金等价物净增加额
	BeginningCash  string `json:"cashs"`   // 期初现金及现金等价物余额
	EndingCash     string `json:"cashe"`   // 期末现金及现金等价物余额
}

// TopTenShareholder 十大股东
type TopTenShareholder struct {
	EndDate      string `json:"jzrq"` // 截止日期
	AnnounceDate string `json:"ggrq"` // 公告日期
	Description  string `json:"gdsm"` // 股东说明
	TotalCount   int    `json:"gdzs"` // 股东总数
	AvgShares    int    `json:"pjcg"` // 平均持股
}

// TopTenFloatShareholder 十大流通股东
type TopTenFloatShareholder struct {
	EndDate      string `json:"jzrq"` // 截止日期
	AnnounceDate string `json:"ggrq"` // 公告日期
}
