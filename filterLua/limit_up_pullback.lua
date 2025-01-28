-- @id: 4
-- @name: 涨停回拉策略
-- @description: 识别最后一个涨停板后跌破回拉的股票，要求最后一个涨停板后至少3天，期间出现跌破涨停收盘价，并且只在最后几天内首次突破回拉

-- ===================== 策略参数配置 =====================
-- 涨停阈值
local LIMIT_UP_THRESHOLD_NORMAL = 9.5     -- 普通股票涨停阈值
local LIMIT_UP_THRESHOLD_STAR = 19.5      -- 科创板、创业板涨停阈值

-- 回拉分析参数
local MIN_DAYS_AFTER_LIMIT = 3            -- 涨停板后最少需要的天数
local BREAK_WINDOW = 3                    -- 允许突破的时间窗口（最后几天）

-- 数据分析范围
local ANALYSIS_DAYS = 15                  -- 分析最近15天的数据

-- 处理单个股票
function process_stock(stock)
    -- 获取该股票的日线数据
    local kdata = api.getKLineData(stock.code, "dn")
    if not kdata then
        log("获取K线数据失败: " .. stock.code)
        return
    end

    -- 如果数据点不足，直接返回
    if #kdata < ANALYSIS_DAYS then
        return
    end

    -- 为每个K线数据添加股票代码
    for _, k in ipairs(kdata) do
        k.code = stock.code
    end

    -- 查找最后一个涨停并分析走势
    local lastLimitUp, analysis = findLastLimitUpAndAnalyze(kdata)

    -- 判断条件
    if lastLimitUp and analysis and
       analysis.daysAfterLimit >= MIN_DAYS_AFTER_LIMIT and
       analysis.hasBroken and
       analysis.hasRecovered and
       analysis.firstBreakIndex and
       analysis.hasBreakthrough then
        
        -- 获取跌破和突破的价格线名称
        local lineName = 
            analysis.brokenLine == "high" and "涨停最高" or
            analysis.brokenLine == "close" and "涨停收盘" or
            analysis.brokenLine == "low" and "涨停最低" or "未知"
        
        -- 发送股票信号到前端
        local lastKData = kdata[#kdata]
        local message = string.format(
            "最后涨停[%s]收盘价%.2f，跌破%s价%.2f后最低至%.2f，最后%d天内突破并回拉至%.2f",
            lastLimitUp.date,
            lastLimitUp.price,
            lineName,
            analysis.brokenPrice,
            analysis.lowestPrice,
            BREAK_WINDOW,
            analysis.recoveryPrice
        )
        
        api.sendSignal(
            stock.code,
            stock.name,
            lastKData.close,
            lastKData.turnover,
            lastKData.change,
            message
        )
    end
end

-- 判断是否涨停（根据不同股票类型判断）
function isLimitUp(kdata)
    if not kdata.code then
        return false
    end
    
    -- 科创板、创业板涨停20%
    if string.match(kdata.code, "^68") or string.match(kdata.code, "^30") then
        return kdata.change >= LIMIT_UP_THRESHOLD_STAR
    end
    -- 其他板块涨停10%
    return kdata.change >= LIMIT_UP_THRESHOLD_NORMAL
end

-- 分析涨停后的走势
function analyzePullback(kdata_list, limitUpIndex, limitUpClose)
    local result = {
        hasBroken = false,           -- 是否出现跌破
        hasRecovered = false,        -- 是否回拉
        daysAfterLimit = 0,          -- 涨停后的天数
        lowestPrice = nil,           -- 最低价
        recoveryPrice = nil,         -- 回拉价格
        firstBreakIndex = nil,       -- 首次突破的位置
        -- 突破线判断结果
        brokenLine = nil,            -- 跌破的价格线（high/close/low）
        brokenPrice = nil,           -- 跌破的价格值
        hasBreakthrough = false      -- 是否突破
    }
    
    -- 确保有足够的后续数据
    if limitUpIndex + MIN_DAYS_AFTER_LIMIT > #kdata_list then
        return result
    end
    
    -- 获取涨停当天的三个价格线
    local limitUpDay = kdata_list[limitUpIndex]
    local limitUpHigh = limitUpDay.high
    local limitUpLow = limitUpDay.low
    
    -- 计算涨停后的天数
    result.daysAfterLimit = #kdata_list - limitUpIndex
    
    -- 分三段分析：
    -- 1. 从涨停后一天到突破窗口前，寻找首次跌破的价格线
    local windowStartIndex = #kdata_list - BREAK_WINDOW + 1
    result.lowestPrice = kdata_list[limitUpIndex + 1].close  -- 初始化最低价
    
    for i = limitUpIndex + 1, windowStartIndex - 1 do
        -- 更新最低价
        if kdata_list[i].close < result.lowestPrice then
            result.lowestPrice = kdata_list[i].close
        end
        
        -- 检查是否跌破各个价格线，优先判断收盘价
        if not result.brokenLine then
            if kdata_list[i].close < limitUpClose then
                result.brokenLine = "close"
                result.brokenPrice = limitUpClose
                result.hasBroken = true
            elseif kdata_list[i].close < limitUpLow then
                result.brokenLine = "low"
                result.brokenPrice = limitUpLow
                result.hasBroken = true
            elseif kdata_list[i].close < limitUpHigh then
                result.brokenLine = "high"
                result.brokenPrice = limitUpHigh
                result.hasBroken = true
            end
        end
        
        -- 如果这期间有突破对应的价格线，直接返回（不符合条件）
        if result.brokenLine and kdata_list[i].close > result.brokenPrice then
            return result
        end
    end
    
    -- 2. 在突破窗口内（最后几天）寻找突破
    if result.hasBroken then
        for i = windowStartIndex, #kdata_list do
            -- 检查是否突破对应的价格线
            if kdata_list[i].close > result.brokenPrice then
                result.hasBreakthrough = true
                result.firstBreakIndex = i
                result.recoveryPrice = kdata_list[i].high
                break
            end
        end
    end
    
    -- 3. 如果找到突破，检查后续是否维持在该价格线之上
    if result.firstBreakIndex then
        result.hasRecovered = true
        result.recoveryPrice = kdata_list[result.firstBreakIndex].high
        
        for i = result.firstBreakIndex, #kdata_list do
            if kdata_list[i].close < result.brokenPrice then
                result.hasRecovered = false
                break
            end
            if kdata_list[i].high > result.recoveryPrice then
                result.recoveryPrice = kdata_list[i].high
            end
        end
    end
    
    return result
end

-- 查找最后一个涨停并分析走势
function findLastLimitUpAndAnalyze(kdata_list)
    local lastLimitUp = {
        date = nil,
        price = nil,
        index = nil
    }
    
    -- 只取最近ANALYSIS_DAYS天的数据
    local recent_data = {}
    local start_idx = math.max(1, #kdata_list - ANALYSIS_DAYS + 1)
    for i = start_idx, #kdata_list do
        table.insert(recent_data, kdata_list[i])
    end
    
    -- 查找最后一个涨停
    for i = #recent_data, 1, -1 do
        local kdata = recent_data[i]
        if isLimitUp(kdata) then
            lastLimitUp.date = kdata.time
            lastLimitUp.price = kdata.close
            lastLimitUp.index = i
            break
        end
    end
    
    -- 如果找到涨停，分析后续走势
    if lastLimitUp.index then
        local analysis = analyzePullback(recent_data, lastLimitUp.index, lastLimitUp.price)
        return lastLimitUp, analysis
    end
    
    return nil, nil
end 