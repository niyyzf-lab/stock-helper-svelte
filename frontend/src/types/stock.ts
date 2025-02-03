export interface StockSignal {
  code: string
  name: string
  price: number
  change: number
  turnover: number
  reason: string
}

export interface Stock {
  dm: string
  mc: string
  testDate: string
}

export interface KLineData {
  d: string;  // 交易时间
  o: number;  // 开盘价
  h: number;  // 最高价
  l: number;  // 最低价
  c: number;  // 收盘价
  v: number;  // 成交量
  e: number;  // 成交额
  zf: number;  // 振幅
  hs: number;  // 换手率
  zd: number;  // 涨跌幅
  zde: number;  // 涨跌额
}

export interface TestResult {
  direction: 'up' | 'down' | 'shock'
  correct: boolean
  nextPrice: number
  currentPrice: number
  maxPrice: number
  minPrice: number
  priceChange: number
  actualDirection: 'up' | 'down' | 'shock'
  shortTermTrend: 'up' | 'down' | 'shock'
  shortTermChange: number
  longTermTrend: 'up' | 'down' | 'shock'
  longTermChange: number
  daysCount: number
  klineData: KLineData[]  // 历史K线数据(前30天)
  futureData: KLineData[] // 未来K线数据(15天)
  prices: number[]        // 用于图表展示的价格数组
}

export interface Progress {
  current: number
  total: number
  currentStock: string
  stage: string
}

export interface QuizOption {
  id: string;
  content: string;
}

export type StockQuiz = {
  quizzes: Array<{
    title: string
    question: string
    options: Array<{
      id: string
      content: string
    }>
    correctAnswer: string
    explanation: string
  }>
  level: number
}
