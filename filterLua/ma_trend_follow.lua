-- @id: 1245125
-- @name: MA趋势跟踪策略
-- @description: 监测均线金叉后N天内的持续上涨趋势，要求趋势保持向上且波动幅度在限定范围内

-- 策略参数
local STRATEGY_PARAMS = {
    SHORT_MA = 5,      -- 短期均线天数
    LONG_MA = 10,      -- 长期均线天数
    TREND_MA = 20,     -- 趋势均线天数
    CHECK_DAYS = 5,    -- 金叉后检查天数
    MIN_PERIODS = 30,  -- 最小所需数据周期
    
    -- 趋势参数
    MIN_TREND_SLOPE = 0.3,    -- 最小趋势斜率
    MAX_TREND_SLOPE = 2.0,    -- 最大趋势斜率
    MAX_DEVIATION = 0.02,     -- 最大偏离度(单日)
    MIN_UP_DAYS_RATIO = 0.8,  -- 上涨天数比例要求
    
    -- 量价参数
    MIN_VOLUME_RATIO = 1.5,   -- 金叉时最小成交量放大倍数
    MIN_TURNOVER = 3.0,       -- 最小换手率(%)
    MAX_TURNOVER = 15.0       -- 最大换手率(%)
}

-- 处理单个股票
function process_stock(stock)
    -- 获取日线数据
    local kdata = api.getKLineData(stock.code, "dh")
    if not kdata or #kdata < STRATEGY_PARAMS.MIN_PERIODS then
        return
    end

    -- 提取价格和成交量数据
    local prices = {}
    local volumes = {}
    for _, k in ipairs(kdata) do
        table.insert(prices, k.close)
        table.insert(volumes, k.volume)
    end

    -- 计算技术指标
    local indicators = calculate_indicators(prices, volumes)
    
    -- 分析交易信号
    local signal = analyze_trading_signal(kdata, indicators)
    
    -- 发送交易信号
    if signal.valid then
        local last_data = kdata[#kdata]
        api.sendSignal(
            stock.code,
            stock.name,
            last_data.close,
            last_data.turnover,
            last_data.change,
            generate_signal_message(signal, kdata)
        )
    end
end

-- 计算技术指标
function calculate_indicators(prices, volumes)
    return {
        short_ma = api.indicator.calculateMA(prices, "sma", STRATEGY_PARAMS.SHORT_MA),
        long_ma = api.indicator.calculateMA(prices, "sma", STRATEGY_PARAMS.LONG_MA),
        trend_ma = api.indicator.calculateMA(prices, "sma", STRATEGY_PARAMS.TREND_MA),
        volume_ma = api.indicator.calculateMA(volumes, "sma", 5)
    }
end

-- 分析交易信号
function analyze_trading_signal(kdata, indicators)
    local last_idx = #kdata
    
    -- 查找最近的金叉
    local cross_info = find_golden_cross(indicators.short_ma, indicators.long_ma, last_idx)
    if not cross_info.found then
        return {valid = false}
    end
    
    -- 检查金叉后的趋势
    local trend_info = check_trend_after_cross(kdata, indicators, cross_info.index)
    if not trend_info.valid then
        return {valid = false}
    end
    
    -- 检查成交量和换手率
    if not check_volume_and_turnover(kdata, indicators, cross_info.index) then
        return {valid = false}
    end
    
    -- 返回有效信号
    return {
        valid = true,
        cross_info = cross_info,
        trend_info = trend_info,
        days_checked = trend_info.days_checked,
        up_days_ratio = trend_info.up_days_ratio,
        avg_slope = trend_info.avg_slope
    }
end

-- 查找金叉位置
function find_golden_cross(short_ma, long_ma, last_idx)
    for i = last_idx, last_idx - STRATEGY_PARAMS.CHECK_DAYS, -1 do
        if i < 2 then break end
        
        -- 检查金叉
        if short_ma[i] > long_ma[i] and short_ma[i-1] <= long_ma[i-1] then
            return {
                found = true,
                index = i
            }
        end
    end
    
    return {found = false}
end

-- 检查金叉后的趋势
function check_trend_after_cross(kdata, indicators, cross_idx)
    local result = {
        valid = false,
        up_days = 0,
        total_days = 0,
        max_deviation = 0,
        avg_slope = 0
    }
    
    -- 计算检查的天数
    local days_to_check = math.min(
        STRATEGY_PARAMS.CHECK_DAYS,
        #kdata - cross_idx
    )
    
    if days_to_check < 2 then
        return result
    end
    
    -- 计算趋势特征
    local slopes = {}
    local deviations = {}
    local up_days = 0
    
    for i = cross_idx, cross_idx + days_to_check - 1 do
        -- 计算日间斜率
        local daily_slope = (kdata[i+1].close - kdata[i].close) / kdata[i].close
        table.insert(slopes, daily_slope)
        
        -- 计算与趋势线的偏离度
        local deviation = math.abs(kdata[i].close - indicators.trend_ma[i]) / indicators.trend_ma[i]
        table.insert(deviations, deviation)
        
        -- 统计上涨天数
        if daily_slope > 0 then
            up_days = up_days + 1
        end
    end
    
    -- 计算平均斜率和最大偏离度
    local avg_slope = calculate_average(slopes)
    local max_deviation = math.max(table.unpack(deviations))
    local up_days_ratio = up_days / days_to_check
    
    -- 检查趋势是否符合要求
    result.valid = (
        avg_slope >= STRATEGY_PARAMS.MIN_TREND_SLOPE and
        avg_slope <= STRATEGY_PARAMS.MAX_TREND_SLOPE and
        max_deviation <= STRATEGY_PARAMS.MAX_DEVIATION and
        up_days_ratio >= STRATEGY_PARAMS.MIN_UP_DAYS_RATIO
    )
    
    result.days_checked = days_to_check
    result.up_days_ratio = up_days_ratio
    result.avg_slope = avg_slope
    result.max_deviation = max_deviation
    
    return result
end

-- 检查成交量和换手率
function check_volume_and_turnover(kdata, indicators, cross_idx)
    -- 检查成交量放大
    local volume_ratio = kdata[cross_idx].volume / indicators.volume_ma[cross_idx]
    if volume_ratio < STRATEGY_PARAMS.MIN_VOLUME_RATIO then
        return false
    end
    
    -- 检查换手率范围
    local turnover = kdata[cross_idx].turnover
    if turnover < STRATEGY_PARAMS.MIN_TURNOVER or turnover > STRATEGY_PARAMS.MAX_TURNOVER then
        return false
    end
    
    return true
end

-- 计算平均值
function calculate_average(t)
    local sum = 0
    for _, v in ipairs(t) do
        sum = sum + v
    end
    return sum / #t
end

-- 生成信号消息
function generate_signal_message(signal, kdata)
    return string.format(
        "MA%d上穿MA%d形成金叉，后续%d天内保持上涨趋势，上涨天数比例%.2f，平均斜率%.2f%%",
        STRATEGY_PARAMS.SHORT_MA,
        STRATEGY_PARAMS.LONG_MA,
        signal.days_checked,
        signal.up_days_ratio * 100,
        signal.avg_slope * 100
    )
end 