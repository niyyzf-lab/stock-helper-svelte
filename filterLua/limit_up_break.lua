-- @id: 2
-- @name: 涨停突破策略
-- @description: 识别涨停突破的股票，结合量价关系进行分析

-- ===================== 策略参数配置 =====================
-- 涨停阈值
local LIMIT_UP_THRESHOLD_NORMAL = 9.5     -- 普通股票涨停阈值
local LIMIT_UP_THRESHOLD_STAR = 19.5      -- 科创板、创业板涨停阈值

-- 涨停次数要求
local MIN_LIMIT_UP_COUNT = 3              -- 最少涨停次数
local MAX_LIMIT_UP_COUNT = 5              -- 最多涨停次数

-- 回调分析参数
local MIN_RETRACEMENT_DAYS = 3            -- 最少回调天数
local MAX_RETRACEMENT_DAYS = 10           -- 最多回调天数
local MIN_VOLUME_RATIO = 0.8              -- 回调期间最小量能比

-- 换手率要求
local MIN_TURNOVER_RATE = 4.0             -- 最小换手率(%)

-- 数据分析范围
local ANALYSIS_DAYS = 90                  -- 分析最近90天的数据

-- 更新进度并检查是否应该停止
local function updateProgressAndCheck(code, processed)
    local shouldContinue = api.updateProgress(code, processed)
    if not shouldContinue then
        return false
    end
    return true
end

-- 判断是否涨停（根据不同股票类型判断）
function isLimitUp(kdata)
    if not kdata.code then
        return false
    end
    
    log(string.format("检查涨停 - 代码: %s, 涨幅: %.2f%%", kdata.code, kdata.change))
    -- 科创板、创业板涨停20%
    if string.match(kdata.code, "^68") or string.match(kdata.code, "^30") then
        log(string.format("科创板/创业板股票，需要20%%涨停"))
        return kdata.change >= LIMIT_UP_THRESHOLD_STAR
    end
    -- 其他板块涨停10%
    log(string.format("普通股票，需要10%%涨停"))
    return kdata.change >= LIMIT_UP_THRESHOLD_NORMAL
end

-- 检查换手率条件
function checkTurnoverRate(kdata)
    log(string.format("检查换手率 - 换手率: %.2f%%", kdata.turnover))
    return kdata.turnover >= MIN_TURNOVER_RATE
end

-- 获取90天内最高价
function get90DaysHigh(kdata_list)
    local high = 0
    for _, kdata in ipairs(kdata_list) do
        if kdata.high > high then
            high = kdata.high
        end
    end
    return high
end

-- 获取指定日期之前的最高价
function getHighestPriceBeforeDate(kdata_list, date)
    local highest = 0
    for _, kdata in ipairs(kdata_list) do
        if kdata.time < date then
            if kdata.high > highest then
                highest = kdata.high
            end
        end
    end
    return highest
end

-- 分析最后涨停后的回调
function analyzeRetracement(kdata_list, lastLimitUpIndex, openPrice)
    local retracement = {
        days = 0,          -- 回调天数
        lowestPrice = nil, -- 最低价
        breakSupport = false, -- 是否跌破支撑
        avgVolume = 0,      -- 回调期间平均成交量
        volumeRatio = 0     -- 回调期间量能比（相对涨停日）
    }
    
    -- 确保有足够的后续数据进行分析
    if lastLimitUpIndex + MAX_RETRACEMENT_DAYS > #kdata_list then
        return false, retracement
    end
    
    -- 获取涨停日成交量作为基准
    local limitUpVolume = kdata_list[lastLimitUpIndex].volume
    local totalVolume = 0
    
    -- 检查后续3-10天
    for i = lastLimitUpIndex + 1, math.min(lastLimitUpIndex + MAX_RETRACEMENT_DAYS, #kdata_list) do
        retracement.days = retracement.days + 1
        local currentLow = kdata_list[i].low
        local currentVolume = kdata_list[i].volume
        totalVolume = totalVolume + currentVolume
        
        -- 更新最低价
        if not retracement.lowestPrice or currentLow < retracement.lowestPrice then
            retracement.lowestPrice = currentLow
        end
        
        -- 检查是否跌破支撑价
        if currentLow < openPrice then
            retracement.breakSupport = true
            break
        end
        
        -- 计算量能指标
        retracement.avgVolume = totalVolume / retracement.days
        retracement.volumeRatio = retracement.avgVolume / limitUpVolume
        
        -- 如果已经回调3天以上且没有跌破支撑，符合条件
        if retracement.days >= MIN_RETRACEMENT_DAYS 
           and not retracement.breakSupport 
           and retracement.volumeRatio >= MIN_VOLUME_RATIO then
            return true, retracement
        end
    end
    
    return false, retracement
end

-- 比较两个日期字符串
function compareDate(date1, date2)
    -- 日期格式：YYYY-MM-DD
    return date1 < date2
end

-- 按日期排序K线数据
function sortKLineByDate(kdata_list)
    local sorted = {}
    for _, kdata in ipairs(kdata_list) do
        table.insert(sorted, kdata)
    end
    
    table.sort(sorted, function(a, b)
        return compareDate(a.time, b.time)
    end)
    
    return sorted
end

-- 查找最近的涨停数量并判断突破
function findRecentLimitUps(kdata_list)
    local limitUps = {}
    local high90 = get90DaysHigh(kdata_list)
    local lastLimitUpIndex = nil
    
    log("\n开始数据排序和处理...")
    log(string.format("原始数据: 第一条: %s, 最后一条: %s", 
        kdata_list[1].time, 
        kdata_list[#kdata_list].time))
    
    -- 按日期排序，并只取最近90个交易日
    local sorted_kdata = sortKLineByDate(kdata_list)
    log(string.format("排序后数据: 第一条: %s, 最后一条: %s", 
        sorted_kdata[1].time, 
        sorted_kdata[#sorted_kdata].time))
    
    -- 只取最近90个交易日的数据
    local recent_kdata = {}
    local start_idx = math.max(1, #sorted_kdata - ANALYSIS_DAYS)
    for i = start_idx, #sorted_kdata do
        table.insert(recent_kdata, sorted_kdata[i])
    end
    
    -- 添加调试信息
    for i, kdata in ipairs(recent_kdata) do
        if kdata.change >= 9.5 then  -- 临时添加，检查是否有涨停数据
            log(string.format("发现可能的涨停: 日期=%s, 代码=%s, 涨幅=%.2f%%", 
                kdata.time, kdata.code, kdata.change))
        end
    end
    
    -- 查找涨停记录
    for i, kdata in ipairs(recent_kdata) do
        -- 确保code字段存在
        kdata.code = kdata.code or recent_kdata[1].code
        
        if isLimitUp(kdata) then
            -- 获取该涨停日之前的最高价
            local previousHigh = getHighestPriceBeforeDate(recent_kdata, kdata.time)
            
            log(string.format("\n发现涨停: %s", kdata.time))
            log(string.format("当日价格: 开盘:%.2f 最高:%.2f 最低:%.2f 收盘:%.2f", 
                kdata.open, kdata.high, kdata.low, kdata.close))
            log(string.format("涨幅: %.2f%%", kdata.change))
            log(string.format("之前高点: %.2f %s", 
                previousHigh,
                kdata.high > previousHigh and "【突破】" or ""))
            
            local isBreakthrough = kdata.high > previousHigh
            
            table.insert(limitUps, 1, {
                date = kdata.time,
                change = kdata.change,
                volume = kdata.volume,
                amount = kdata.amount,
                high = kdata.high,
                open = kdata.open,
                isBreakthrough = isBreakthrough,
                index = i
            })
            lastLimitUpIndex = i
            
            log(string.format("当前涨停列表: "))
            for idx, record in ipairs(limitUps) do
                log(string.format("%d. %s (涨幅: %.2f%%)", idx, record.date, record.change))
            end
            
            -- 如果已经找到5个涨停，提前结束
            if #limitUps >= MAX_LIMIT_UP_COUNT then
                break
            end
        end
    end
    
    if #limitUps > 0 then
        log("\n涨停统计:")
        log(string.format("找到涨停数量: %d", #limitUps))
        log(string.format("第一个涨停: %s", limitUps[1].date))
        log(string.format("最后涨停: %s", limitUps[#limitUps].date))
        if lastLimitUpIndex then
            log(string.format("最后涨停位置: %d/%d", lastLimitUpIndex, #recent_kdata))
        end
    end
    
    -- 分析最后一个涨停后的回调
    local validRetracement = false
    local retracementInfo = {
        days = 0,
        lowestPrice = nil,
        breakSupport = false,
        avgVolume = 0,
        volumeRatio = 0
    }
    
    if lastLimitUpIndex and #limitUps > 0 then
        log("\n开始回调分析...")
        log(string.format("分析位置: %d", lastLimitUpIndex))
        log(string.format("使用开盘价: %.2f", limitUps[1].open))
        
        validRetracement, retracementInfo = analyzeRetracement(
            recent_kdata,
            lastLimitUpIndex,
            limitUps[1].open
        )
        
        if validRetracement then
            log("回调分析结果: 符合条件")
            log(string.format("回调天数: %d", retracementInfo.days))
            log(string.format("最低价: %.2f", retracementInfo.lowestPrice))
            log(string.format("回调期间日均量: %.0f", retracementInfo.avgVolume))
            log(string.format("量能比: %.2f", retracementInfo.volumeRatio))
        else
            log("回调分析结果: 不符合条件")
            if retracementInfo.breakSupport then
                log("原因: 跌破支撑")
            elseif retracementInfo.volumeRatio < 0.8 then
                log(string.format("原因: 量能不足 (量能比: %.2f)", retracementInfo.volumeRatio))
            elseif retracementInfo.days < 3 then
                log("原因: 回调天数不足")
            end
        end
    end
    
    return limitUps, high90, validRetracement, retracementInfo
end

-- 并发处理股票数据
function processStocksParallel(indices, batchSize)
    local processed = 0
    local total = #indices
    
    -- 将股票分组
    for i, index in ipairs(indices) do
        -- 声明局部变量
        local limitUps = {}
        local high90 = 0
        local validRetracement = false
        local retracementInfo = {
            days = 0,
            lowestPrice = nil,
            breakSupport = false,
            avgVolume = 0,
            volumeRatio = 0
        }
        
        -- 更新进度并检查是否应该停止
        if not updateProgressAndCheck(index.code, processed) then
            log("策略执行停止")
            return
        end
        
        -- 获取90天K线数据
        local kdata, err = api.getKLineData(index.code, "dq")
        if err then
            log(string.format("获取%s数据失败: %s", index.name, err))
            goto continue
        end
        
        -- 为每个K线数据添加股票代码
        for _, k in ipairs(kdata) do
            k.code = index.code
        end
        
        -- 查找涨停记录和分析回调
        limitUps, high90, validRetracement, retracementInfo = findRecentLimitUps(kdata)
        
        -- 判断条件：
        log("\n开始检查筛选条件:")
        log(string.format("1. 涨停次数: %d (要求3-5次)", #limitUps))
        
        -- 检查 limitUps 是否为空
        if #limitUps == 0 then
            goto continue
        end
        
        log(string.format("2. 最新涨停突破: %s", limitUps[1].isBreakthrough and "是" or "否"))
        log(string.format("3. 回调企稳: %s", validRetracement and "是" or "否"))
        log(string.format("4. 换手率: %.2f%% (要求>4%%)", kdata[#kdata].turnover))
        
        if #limitUps >= MIN_LIMIT_UP_COUNT and #limitUps <= MAX_LIMIT_UP_COUNT 
           and limitUps[1].isBreakthrough 
           and validRetracement
           and checkTurnoverRate(kdata[#kdata]) then
            
            -- 发送股票信号到前端
            local lastKData = kdata[#kdata] -- 最新一天的数据
            api.sendSignal(
                index.code,
                index.name,
                lastKData.close,
                lastKData.turnover,
                lastKData.change,
                lastKData.volume,
                lastKData.amount,
                string.format("近90天内出现%d次涨停，最新涨停突破前期高点且回调企稳，换手率%.1f%%", #limitUps, lastKData.turnover)
            )
            
            log(string.format("\n发现目标股票: %s (%s)", index.name, index.code))
            log(string.format("最新涨停突破前期高点: %.2f", limitUps[1].high))
            log(string.format("最新换手率: %.2f%%", lastKData.turnover))
            log(string.format("回调分析:"))
            log(string.format("回调天数: %d", retracementInfo.days))
            log(string.format("最低价: %.2f", retracementInfo.lowestPrice))
            log(string.format("支撑位: %.2f (涨停开盘价)", limitUps[1].open))
            log("涨停记录:")
            log("----------------------------------------")
            for _, record in ipairs(limitUps) do
                log(string.format("日期: %s  涨幅: %.2f%%  最高: %.2f  成交量: %.0f  成交额: %.2f %s",
                    record.date,
                    record.change,
                    record.high,
                    record.volume,
                    record.amount,
                    record.isBreakthrough and "【突破】" or ""
                ))
            end
        end
        
        ::continue::
        processed = processed + 1
    end
end

-- 策略主函数
function main()
    log("开始执行涨停突破策略...")
    
    -- 获取指数列表
    local indices, err = api.getIndexList()
    if err then
        error(string.format("获取指数列表失败: %s", err))
    end
    
    -- 并发处理股票数据，每批10个
    processStocksParallel(indices, 10)
    
    log("\n策略执行完成")
end

