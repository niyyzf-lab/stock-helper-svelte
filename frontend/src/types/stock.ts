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
  o: number
  c: number
  h: number
  l: number
  d: string
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
