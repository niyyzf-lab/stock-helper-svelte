-- @id: 1
-- @name: MA交叉策略
-- @description: 监测均线交叉信号，要求最近3-5天内出现金叉且之后无死叉，且价格在98日线之上，90天内有2-5个涨停，金叉时成交量放大，换手率大于5%，趋势向上

-- 策略参数
local STRATEGY_PARAMS = {
    SHORT_MA = 6,     -- 短期均线天数
    LONG_MA = 18,     -- 长期均线天数
    TREND_MA = 98,    -- 趋势均线天数
    MIN_DAYS = 3,     -- 金叉最少要满3天
    MAX_DAYS = 15,    -- 金叉最多不超过15天
    MIN_PERIODS = 98, -- 最小所需数据周期
    CHECK_DAYS = 90,  -- 涨停检查天数
    MIN_LIMIT_UP = 2, -- 最小涨停数
    MAX_LIMIT_UP = 5, -- 最大涨停数
    VOLUME_RATIO = 1.5, -- 成交量放大倍数
    MIN_TURNOVER = 5, -- 最小换手率(%)
    TREND_PERIOD = 20,-- 趋势分析周期 VOLUME
    MIN_TREND = 0.2,  -- 最小趋势强度要求

    -- 不同市场涨停幅度（稍微放宽一些，避免错过接近涨停的）
    LIMIT_UP_RATES = {
        ["00"] = 0.095, -- 深证主板 9.5%（放宽）
        ["30"] = 0.195, -- 创业板 19.5%（放宽）
        ["60"] = 0.095, -- 上证主板 9.5%（放宽）
        ["68"] = 0.195, -- 科创板 19.5%（放宽）
    }
}

-- 获取股票对应的涨停幅度
local function get_limit_up_rate(stock_code)
    local prefix = string.sub(stock_code, 1, 2)
    return STRATEGY_PARAMS.LIMIT_UP_RATES[prefix] or 0.098
end

-- 处理单个股票
function process_stock(stock)
    -- 获取该股票的日线数据
    local kdata = api.getKLineData(stock.code, "dh")
    if not kdata or #kdata < STRATEGY_PARAMS.MIN_PERIODS then
        return
    end

    -- 提取价格数据
    local prices = {}
    local volumes = {}
    for _, k in ipairs(kdata) do
        table.insert(prices, k.close)
        table.insert(volumes, k.volume)
    end

    -- 计算技术指标
    local indicators = calculate_indicators(prices, volumes)
    
    -- 分析交易信号
    local signal = analyze_trading_signal(stock.code, kdata, indicators)
    
    -- 发送交易信号
    if signal.valid then
        local last_data = kdata[#kdata]
        api.sendSignal(
            stock.code,
            stock.name,
            last_data.close,
            last_data.turnover,
            last_data.change,
            generate_signal_message(signal, stock.code, kdata)
        )
    end
end

-- 计算技术指标
function calculate_indicators(prices, volumes)
    return {
        short_ma = api.indicator.calculateMA(prices, "sma", STRATEGY_PARAMS.SHORT_MA),
        long_ma = api.indicator.calculateMA(prices, "sma", STRATEGY_PARAMS.LONG_MA),
        trend_ma = api.indicator.calculateMA(prices, "sma", STRATEGY_PARAMS.TREND_MA),
        trend = api.indicator.calculateTrend(prices, STRATEGY_PARAMS.TREND_PERIOD),
        volume_ma = api.indicator.calculateMA(volumes, "sma", 5)
    }
end

-- 检查涨停数量
function check_limit_up_count(stock_code, kdata, last_idx)
    local result = {
        count = 0,
        dates = {}
    }
    
    local limit_up_rate = get_limit_up_rate(stock_code)
    local start_idx = math.max(1, last_idx - STRATEGY_PARAMS.CHECK_DAYS + 1)
    
    for i = start_idx, last_idx do
        if i > 1 then
            local change_rate = (kdata[i].close - kdata[i-1].close) / kdata[i-1].close
            if change_rate >= limit_up_rate then
                result.count = result.count + 1
                table.insert(result.dates, kdata[i].time)
            end
        end
    end
    
    return result
end

-- 分析交易信号
function analyze_trading_signal(stock_code, kdata, indicators)
    local last_idx = #kdata
    
    -- 查找交叉信号
    local cross_info = find_recent_cross(indicators.short_ma, indicators.long_ma, last_idx)
    if not cross_info.has_golden_cross or cross_info.has_death_cross then
        return {valid = false}
    end
    
    -- 检查趋势
    local trend_valid = check_trend(kdata, indicators, last_idx)
    if not trend_valid then
        return {valid = false}
    end
    
    -- 检查涨停数量
    local limit_up_info = check_limit_up_count(stock_code, kdata, last_idx)
    if not (limit_up_info.count >= STRATEGY_PARAMS.MIN_LIMIT_UP and 
            limit_up_info.count <= STRATEGY_PARAMS.MAX_LIMIT_UP) then
        return {valid = false}
    end
    
    -- 检查金叉时的成交量和换手率
    local cross_idx = last_idx - cross_info.days_since_cross
    if not (check_volume_increase(kdata, indicators.volume_ma, cross_idx) and 
            kdata[cross_idx].turnover >= STRATEGY_PARAMS.MIN_TURNOVER) then
        return {valid = false}
    end
    
    -- 返回有效信号
    return {
        valid = true,
        cross_info = cross_info,
        limit_up_info = limit_up_info,
        trend_strength = indicators.trend[#indicators.trend],
        cross_turnover = kdata[cross_idx].turnover
    }
end

-- 检查趋势
function check_trend(kdata, indicators, idx)
    -- 检查价格是否在趋势线之上
    local is_above_trend = kdata[idx].close > indicators.trend_ma[idx]
    
    -- 检查趋势强度
    local trend_strength = indicators.trend[idx]
    local has_strong_trend = trend_strength >= STRATEGY_PARAMS.MIN_TREND
    
    -- 放宽条件：价格在趋势线上或趋势强度达标即可
    return is_above_trend or has_strong_trend
end

-- 检查成交量放大
function check_volume_increase(kdata, volume_ma, cross_idx)
    if cross_idx < 3 then 
        return false 
    end
    -- 检查金叉当天或后一天的成交量放大
    return kdata[cross_idx].volume >= volume_ma[cross_idx] * STRATEGY_PARAMS.VOLUME_RATIO or
           (cross_idx + 1 <= #kdata and kdata[cross_idx + 1].volume >= volume_ma[cross_idx + 1] * STRATEGY_PARAMS.VOLUME_RATIO)
end

-- 查找最近的交叉信号
function find_recent_cross(short_ma, long_ma, idx)
    local result = {
        has_golden_cross = false,
        has_death_cross = false,
        days_since_cross = 0
    }
    
    -- 修改遍历逻辑，先找到最近的金叉
    local found_cross = false
    for i = 0, STRATEGY_PARAMS.MAX_DAYS do
        local curr_idx = idx - i
        if curr_idx < 2 then break end
        
        -- 检查金叉
        if short_ma[curr_idx] > long_ma[curr_idx] and 
           short_ma[curr_idx-1] <= long_ma[curr_idx-1] then
            -- 金叉后必须保持至少3天的多头排列
            local valid_days = 0
            for j = curr_idx, math.min(curr_idx + STRATEGY_PARAMS.MIN_DAYS - 1, idx) do
                if short_ma[j] > long_ma[j] then
                    valid_days = valid_days + 1
                else
                    break
                end
            end
            
            if valid_days >= STRATEGY_PARAMS.MIN_DAYS then
                result.has_golden_cross = true
                result.days_since_cross = i
                -- 只检查金叉之后到现在是否有死叉
                for j = curr_idx + 1, idx do
                    if short_ma[j] < long_ma[j] and 
                       short_ma[j-1] >= long_ma[j-1] then
                        result.has_death_cross = true
                        break
                    end
                end
                found_cross = true
                break
            end
        end
    end
    
    return result
end

-- 生成信号消息
function generate_signal_message(signal, stock_code, kdata)
    local limit_up_rate = get_limit_up_rate(stock_code) * 100
    
    return string.format(
        "MA%d上穿MA%d形成金叉，距今%d天，站上MA%d，%d天内%d个涨停(涨停幅度%.1f%%)，趋势强度%.2f，换手率%.2f%%",
        STRATEGY_PARAMS.SHORT_MA,
        STRATEGY_PARAMS.LONG_MA,
        signal.cross_info.days_since_cross,
        STRATEGY_PARAMS.TREND_MA,
        STRATEGY_PARAMS.CHECK_DAYS,
        signal.limit_up_info.count,
        limit_up_rate,
        signal.trend_strength,
        signal.cross_turnover
    )
end 