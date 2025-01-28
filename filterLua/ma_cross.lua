-- @id: 1
-- @name: MA交叉策略
-- @description: 监测均线交叉信号，要求最近3-5天内出现金叉且之后无死叉，且价格在98日线之上，90天内有2-5个涨停，金叉时成交量放大，换手率大于5%

-- 策略参数
local STRATEGY_PARAMS = {
    SHORT_MA = 6,     -- 短期均线天数
    LONG_MA = 18,     -- 长期均线天数
    TREND_MA = 98,    -- 趋势均线天数
    MIN_DAYS = 3,     -- 最小检查天数
    MAX_DAYS = 5,     -- 最大检查天数
    MIN_PERIODS = 98, -- 最小所需数据周期
    CHECK_DAYS = 90,  -- 涨停检查天数
    MIN_LIMIT_UP = 2, -- 最小涨停数
    MAX_LIMIT_UP = 5, -- 最大涨停数
    VOLUME_RATIO = 2, -- 成交量放大倍数
    MIN_TURNOVER = 5  -- 最小换手率(%)
}

-- 处理单个股票
function process_stock(stock)
    -- 获取该股票的日线数据
    local kdata = api.getKLineData(stock.code, "dh")
    if not kdata then
        log("获取K线数据失败: " .. stock.code)
        return
    end

    -- 如果数据点不足，直接返回
    if #kdata < STRATEGY_PARAMS.MIN_PERIODS then
        return
    end

    -- 计算各均线
    local short_ma = calculate_ma(kdata, STRATEGY_PARAMS.SHORT_MA)
    local long_ma = calculate_ma(kdata, STRATEGY_PARAMS.LONG_MA)
    local trend_ma = calculate_ma(kdata, STRATEGY_PARAMS.TREND_MA)

    -- 检查最近的交叉信号
    local last_idx = #kdata
    local cross_info = find_recent_cross(short_ma, long_ma, last_idx)
    
    -- 检查是否在98日线之上
    local last_close = kdata[last_idx].close
    local is_above_trend = last_close > trend_ma[last_idx]
    
    -- 检查90天内涨停数量
    local limit_up_info = check_limit_up_count(kdata, last_idx)
    local limit_up_valid = limit_up_info.count >= STRATEGY_PARAMS.MIN_LIMIT_UP and 
                          limit_up_info.count <= STRATEGY_PARAMS.MAX_LIMIT_UP
    
    -- 检查金叉时的成交量和换手率
    local volume_valid = false
    local turnover_valid = false
    if cross_info.has_golden_cross then
        local cross_idx = last_idx - cross_info.days_since_cross
        volume_valid = check_volume_increase(kdata, cross_idx)
        turnover_valid = kdata[cross_idx].turnover >= STRATEGY_PARAMS.MIN_TURNOVER
    end
    
    if cross_info.has_golden_cross and 
       not cross_info.has_death_cross and 
       is_above_trend and 
       limit_up_valid and 
       volume_valid and
       turnover_valid then
        -- 发送金叉信号
        local last_data = kdata[last_idx]
        local cross_idx = last_idx - cross_info.days_since_cross
        api.sendSignal(
            stock.code,
            stock.name,
            last_data.close,
            last_data.turnover,
            last_data.change,
            string.format("MA%d上穿MA%d形成金叉，距今%d天，站上MA%d，%d天内%d个涨停，换手率%.2f%%", 
                STRATEGY_PARAMS.SHORT_MA,
                STRATEGY_PARAMS.LONG_MA,
                cross_info.days_since_cross,
                STRATEGY_PARAMS.TREND_MA,
                STRATEGY_PARAMS.CHECK_DAYS,
                limit_up_info.count,
                kdata[cross_idx].turnover
            )
        )
    end
end

-- 计算移动平均线
function calculate_ma(kdata, period)
    local ma = {}
    for i = period, #kdata do
        local sum = 0
        for j = 0, period - 1 do
            sum = sum + kdata[i-j].close
        end
        ma[i] = sum / period
    end
    return ma
end

-- 检查涨停数量
function check_limit_up_count(kdata, last_idx)
    local result = {
        count = 0,
        dates = {}
    }
    
    local start_idx = math.max(1, last_idx - STRATEGY_PARAMS.CHECK_DAYS + 1)
    for i = start_idx, last_idx do
        if i > 1 then  -- 确保有前一天的数据
            -- 涨停判断：收盘价涨幅超过9.8%
            if (kdata[i].close - kdata[i-1].close) / kdata[i-1].close >= 0.098 then
                result.count = result.count + 1
                table.insert(result.dates, kdata[i].date)
            end
        end
    end
    
    return result
end

-- 检查成交量放大
function check_volume_increase(kdata, cross_idx)
    if cross_idx < 6 then return false end  -- 确保有足够的数据计算平均成交量
    
    -- 计算前5天的平均成交量
    local avg_volume = 0
    for i = 1, 5 do
        avg_volume = avg_volume + kdata[cross_idx - i].volume
    end
    avg_volume = avg_volume / 5
    
    -- 检查金叉当天的成交量是否放大
    return kdata[cross_idx].volume >= avg_volume * STRATEGY_PARAMS.VOLUME_RATIO
end

-- 查找最近的交叉信号
function find_recent_cross(short_ma, long_ma, idx)
    local result = {
        has_golden_cross = false,
        has_death_cross = false,
        days_since_cross = 0
    }
    
    -- 检查指定天数范围
    for i = 0, STRATEGY_PARAMS.MAX_DAYS - 1 do
        local curr_idx = idx - i
        if curr_idx < 2 then break end
        
        -- 检查金叉
        if short_ma[curr_idx] > long_ma[curr_idx] and short_ma[curr_idx-1] < long_ma[curr_idx-1] then
            -- 只在指定的天数范围内记录金叉
            if i >= STRATEGY_PARAMS.MIN_DAYS - 1 then
                result.has_golden_cross = true
                result.days_since_cross = i
                -- 找到金叉后，检查之后是否有死叉
                for j = curr_idx + 1, idx do
                    if short_ma[j] < long_ma[j] and short_ma[j-1] > long_ma[j-1] then
                        result.has_death_cross = true
                        break
                    end
                end
                break
            end
        end
    end
    
    return result
end 