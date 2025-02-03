-- @id: 2
-- @name: 十字星触底反转策略
-- @description: 监测最近3天内的十字星形态，结合MACD和RSI确认触底信号，验证后续上涨趋势

-- 策略参数
local STRATEGY_PARAMS = {
    DOJI_DAYS = 3,        -- 检查十字星的天数范围
    DOJI_THRESHOLD = 0.008, -- 开盘收盘价差阈值（放宽到0.8%）
    RSI_PERIOD = 14,      -- RSI周期
    RSI_OVERSOLD = 40,    -- RSI超卖阈值(放宽到40)
    MACD_FAST = 12,       -- MACD快线
    MACD_SLOW = 26,       -- MACD慢线
    MACD_SIGNAL = 9,      -- MACD信号线
    TREND_MA = 20,        -- 趋势均线
    TREND_MA_LONG = 60,   -- 长期趋势均线
    MIN_PERIODS = 60,     -- 最小所需数据周期
    VOLUME_RATIO = 1.0,   -- 成交量放大倍数(降低到1.0)
    PRICE_INCREASE = 0.015, -- 十字星后上涨确认幅度（降低到1.5%）
    BOTTOM_COMPARE_DAYS = 15, -- 相对低点比较天数（缩短到15天）
    TREND_UP_THRESHOLD = 0.01, -- 上升趋势判断阈值（降低到1%）
    MIN_AMOUNT = 2000000  -- 最小成交额要求（降低到200万）
}

-- 日志级别
local LOG_LEVEL = {
    DEBUG = "DEBUG",
    INFO = "INFO",
    WARN = "WARN",
    ERROR = "ERROR"
}

local CURRENT_LOG_LEVEL = LOG_LEVEL.DEBUG  -- 改为DEBUG级别

-- 日志输出函数
function print_log(level, stock_code, message, ...)
    -- 根据日志级别过滤输出
    if level == LOG_LEVEL.DEBUG and CURRENT_LOG_LEVEL ~= LOG_LEVEL.DEBUG then
        return
    end
    
    local formatted_msg = string.format(message, ...)
    if stock_code then
        log(string.format("[%s] [%s] %s", level, stock_code, formatted_msg))
    else
        log(string.format("[%s] %s", level, formatted_msg))
    end
end

-- 检查基本面条件
function check_fundamental(stock)
    -- 检查成交额
    if not stock.volume then
        log(string.format("[DEBUG] [%s] 成交量数据不存在", stock.code))
        return false
    end

    if not stock.close then
        log(string.format("[DEBUG] [%s] 收盘价数据不存在", stock.code))
        return false
    end

    -- 尝试转换数据为数值类型
    local volume = tonumber(stock.volume)
    local close = tonumber(stock.close)

    if not volume or not close then
        log(string.format("[DEBUG] [%s] 数据类型转换失败 - volume: %s, close: %s", 
            stock.code, tostring(stock.volume), tostring(stock.close)))
        return false
    end

    -- 检查数值是否有效
    if not is_valid_number(volume) or not is_valid_number(close) then
        log(string.format("[DEBUG] [%s] 数据无效 - volume: %s, close: %s", 
            stock.code, tostring(volume), tostring(close)))
        return false
    end

    -- 计算成交额
    local amount = volume * close

    -- 检查成交额是否满足最小要求
    if amount < STRATEGY_PARAMS.MIN_AMOUNT then
        log(string.format("[DEBUG] [%s] 成交额不足: %.2f万", stock.code, amount/10000))
        return false
    end
    
    return true
end

-- 处理单个股票
function process_stock(stock)
    -- 立即输出开始处理的日志
    log("[DEBUG] 开始处理股票")
    
    -- 快速检查基本数据是否存在
    if not stock then
        log("[ERROR] 股票数据为空")
        return
    end
    
    if not stock.code then
        log("[ERROR] 股票代码为空")
        return
    end
    
    if not stock.name then
        log("[ERROR] 股票名称为空")
        return
    end

    log(string.format("[DEBUG] [%s] 开始处理: %s", stock.code, stock.name))

    -- 检查市盈率
    if stock.pe_ttm then
        local pe = tonumber(stock.pe_ttm)
        if pe and is_valid_number(pe) and pe > STRATEGY_PARAMS.MAX_PE then
            log(string.format("[DEBUG] [%s] 市盈率过高: %.2f", stock.code, pe))
            return
        end
    end

    -- 获取K线数据
    log(string.format("[DEBUG] [%s] 开始获取K线数据", stock.code))
    local kdata = api.getKLineData(stock.code, "dh")
    
    if not kdata then
        log(string.format("[ERROR] [%s] 获取K线数据失败", stock.code))
        return
    end
    
    if type(kdata) ~= "table" then
        log(string.format("[ERROR] [%s] K线数据类型错误: %s", stock.code, type(kdata)))
        return
    end
    
    if #kdata < STRATEGY_PARAMS.MIN_PERIODS then
        log(string.format("[DEBUG] [%s] K线数据不足: %d天", stock.code, #kdata))
        return
    end

    -- 处理最近的数据
    local process_days = math.min(120, #kdata)
    local start_idx = #kdata - process_days + 1
    local recent_kdata = {}
    
    log(string.format("[DEBUG] [%s] 开始处理最近%d天数据", stock.code, process_days))
    
    -- 验证数据有效性
    for i = start_idx, #kdata do
        local k = kdata[i]
        if not k then
            log(string.format("[ERROR] [%s] 第%d天数据为空", stock.code, i))
            return
        end
        
        -- 检查每个必需字段
        local required_fields = {"open", "close", "high", "low", "volume", "time"}
        for _, field in ipairs(required_fields) do
            if not k[field] then
                log(string.format("[ERROR] [%s] 第%d天缺少%s数据", stock.code, i, field))
                return
            end
            if field ~= "time" then
                local value = tonumber(k[field])
                if not value then
                    log(string.format("[ERROR] [%s] 第%d天%s数据类型错误: %s", stock.code, i, field, tostring(k[field])))
                    return
                end
                if not is_valid_number(value) then
                    log(string.format("[ERROR] [%s] 第%d天%s数据无效: %s", stock.code, i, field, tostring(value)))
                    return
                end
                if value <= 0 then
                    log(string.format("[ERROR] [%s] 第%d天%s数据异常: %.4f", stock.code, i, field, value))
                    return
                end
                k[field] = value  -- 更新为数值类型
            end
        end
        
        table.insert(recent_kdata, k)
    end

    -- 提取价格数据
    log(string.format("[DEBUG] [%s] 开始提取价格数据", stock.code))
    local prices = {
        open = {},
        high = {},
        low = {},
        close = {},
        volume = {}
    }
    
    for _, k in ipairs(recent_kdata) do
        table.insert(prices.open, k.open)
        table.insert(prices.high, k.high)
        table.insert(prices.low, k.low)
        table.insert(prices.close, k.close)
        table.insert(prices.volume, k.volume)
    end

    -- 计算技术指标
    log(string.format("[DEBUG] [%s] 开始计算技术指标", stock.code))
    local indicators = calculate_indicators(prices)
    if not indicators then
        log(string.format("[ERROR] [%s] 计算技术指标失败", stock.code))
        return
    end

    -- 分析交易信号
    log(string.format("[DEBUG] [%s] 开始分析交易信号", stock.code))
    local signal = analyze_trading_signal(stock.code, recent_kdata, indicators)
    
    -- 处理信号
    if signal and signal.valid then
        local last_data = recent_kdata[#recent_kdata]
        if not last_data then
            log(string.format("[ERROR] [%s] 最新K线数据无效", stock.code))
            return
        end
        
        -- 确保所有必需的数据都是有效的
        if not is_valid_number(last_data.close) or
           not is_valid_number(last_data.turnover) or
           not is_valid_number(last_data.change) then
            log(string.format("[ERROR] [%s] 信号数据无效", stock.code))
            return
        end
        
        local message = generate_signal_message(signal, stock.code, recent_kdata)
        log(string.format("[INFO] [%s] 发现有效信号：%s", stock.code, message))
        
        -- 发送信号前再次验证数据
        api.sendSignal(
            stock.code,
            stock.name,
            last_data.close,
            last_data.turnover,
            last_data.change,
            message
        )
    else
        log(string.format("[DEBUG] [%s] 未发现有效信号", stock.code))
    end
    
    log(string.format("[DEBUG] [%s] 处理完成", stock.code))
end

-- 检查数组中的数值是否都有效
function check_valid_array(arr)
    if type(arr) ~= "table" then return false end
    for _, v in ipairs(arr) do
        if not is_valid_number(v) then
            return false
        end
    end
    return true
end

-- 计算技术指标
function calculate_indicators(prices)
    if not prices then
        log("[ERROR] 输入价格数据为空")
        return nil
    end
    
    if not prices.close or #prices.close == 0 then
        log("[ERROR] 收盘价数据为空")
        return nil
    end

    -- 检查输入数据的有效性
    if not check_valid_array(prices.close) then
        log("[ERROR] 收盘价数据包含无效值")
        return nil
    end
    
    if not check_valid_array(prices.volume) then
        log("[ERROR] 成交量数据包含无效值")
        return nil
    end

    local indicators = {}
    
    -- 计算RSI
    log("[DEBUG] 开始计算RSI")
    local rsi, err = api.indicator.calculateRSI(prices.close, STRATEGY_PARAMS.RSI_PERIOD)
    if err then 
        log(string.format("[ERROR] RSI计算失败: %s", tostring(err)))
        return nil 
    end
    if not check_valid_array(rsi) then
        log("[ERROR] RSI结果包含无效值")
        return nil
    end
    indicators.rsi = rsi
    
    -- 计算MACD
    log("[DEBUG] 开始计算MACD")
    local macd, err = api.indicator.calculateMACD(prices.close, 
        STRATEGY_PARAMS.MACD_FAST, 
        STRATEGY_PARAMS.MACD_SLOW, 
        STRATEGY_PARAMS.MACD_SIGNAL)
    if err then 
        log(string.format("[ERROR] MACD计算失败: %s", tostring(err)))
        return nil 
    end
    if not macd or type(macd) ~= "table" or not macd.dif then
        log("[ERROR] MACD结果结构无效")
        return nil
    end
    if not check_valid_array(macd.dif) then
        log("[ERROR] MACD.DIF包含无效值")
        return nil
    end
    indicators.macd = macd.dif
    
    -- 计算MA
    log("[DEBUG] 开始计算MA")
    local ma, err = api.indicator.calculateMA(prices.close, "sma", STRATEGY_PARAMS.TREND_MA)
    if err then 
        log(string.format("[ERROR] MA计算失败: %s", tostring(err)))
        return nil 
    end
    if not check_valid_array(ma) then
        log("[ERROR] MA结果包含无效值")
        return nil
    end
    indicators.ma = ma
    
    -- 计算成交量MA
    log("[DEBUG] 开始计算成交量MA")
    local volume_ma, err = api.indicator.calculateMA(prices.volume, "sma", 5)
    if err then 
        log(string.format("[ERROR] 成交量MA计算失败: %s", tostring(err)))
        return nil 
    end
    if not check_valid_array(volume_ma) then
        log("[ERROR] 成交量MA结果包含无效值")
        return nil
    end
    indicators.volume_ma = volume_ma
    
    -- 计算长期MA
    log("[DEBUG] 开始计算长期MA")
    local ma_long, err = api.indicator.calculateMA(prices.close, "sma", STRATEGY_PARAMS.TREND_MA_LONG)
    if err then 
        log(string.format("[ERROR] 长期MA计算失败: %s", tostring(err)))
        return nil 
    end
    if not check_valid_array(ma_long) then
        log("[ERROR] 长期MA结果包含无效值")
        return nil
    end
    indicators.ma_long = ma_long
    
    log("[DEBUG] 所有技术指标计算完成")
    return indicators
end

-- 检查是否形成十字星
function is_doji(open, close, high, low)
    -- 计算实体与影线比例
    local body = math.abs(close - open)
    local body_ratio = body / (high - low)
    local price = (open + close) / 2
    
    local is_doji = (body / price) <= STRATEGY_PARAMS.DOJI_THRESHOLD and body_ratio <= 0.3
    if is_doji then
        print_log(LOG_LEVEL.DEBUG, nil, "发现十字星形态 - 实体比例: %.2f%%, 影线比例: %.2f%%", 
            (body / price) * 100, body_ratio * 100)
    end
    
    return is_doji
end

-- 检查数值是否有效
function is_valid_number(n)
    return type(n) == "number" and n == n and n ~= math.huge and n ~= -math.huge
end

-- 检查是否处于上升趋势
function is_uptrend(kdata, indicators, idx)
    -- 检查长期均线趋势
    local ma_long = indicators.ma_long
    if not ma_long or not ma_long[idx] or not ma_long[idx-5] then
        return false
    end
    
    -- 防止除以0
    if ma_long[idx-5] <= 0 then
        return false
    end
    
    -- 计算长期均线斜率
    local ma_slope = (ma_long[idx] - ma_long[idx-5]) / ma_long[idx-5]
    
    -- 检查计算结果是否有效
    if not is_valid_number(ma_slope) then
        return false
    end
    
    -- 检查当前价格是否在长期均线上方
    local price_above_ma = kdata[idx].close > ma_long[idx]
    
    -- 检查短期均线趋势
    local ma = indicators.ma
    if not ma or not ma[idx] or not ma[idx-5] then
        return false
    end
    
    -- 防止除以0
    if ma[idx-5] <= 0 then
        return false
    end
    
    -- 计算短期均线斜率
    local short_ma_slope = (ma[idx] - ma[idx-5]) / ma[idx-5]
    
    -- 检查计算结果是否有效
    if not is_valid_number(short_ma_slope) then
        return false
    end
    
    -- 判断是否处于上升趋势
    local is_trend_up = ma_slope > STRATEGY_PARAMS.TREND_UP_THRESHOLD and 
                       short_ma_slope > STRATEGY_PARAMS.TREND_UP_THRESHOLD
    
    if is_trend_up then
        print_log(LOG_LEVEL.DEBUG, nil, "处于上升趋势 - 长期均线斜率: %.2f%%, 短期均线斜率: %.2f%%", 
            ma_slope * 100, short_ma_slope * 100)
    end
    
    return is_trend_up and price_above_ma
end

-- 检查是否为相对低点
function is_relative_bottom(kdata, idx)
    local compare_days = STRATEGY_PARAMS.BOTTOM_COMPARE_DAYS
    local start_idx = math.max(1, idx - compare_days)
    local current_low = kdata[idx].low
    
    -- 计算前期最高价和最低价
    local prev_high = kdata[start_idx].high
    local prev_low = kdata[start_idx].low
    
    for i = start_idx + 1, idx - 1 do
        prev_high = math.max(prev_high, kdata[i].high)
        prev_low = math.min(prev_low, kdata[i].low)
    end
    
    -- 计算前期价格区间
    local price_range = prev_high - prev_low
    -- 防止除以0
    if price_range <= 0 then
        return false
    end
    
    -- 当前价格在前期价格区间的相对位置（0-1）
    local relative_position = (current_low - prev_low) / price_range
    
    -- 检查计算结果是否有效
    if not is_valid_number(relative_position) then
        return false
    end
    
    -- 判断是否为相对低点（在前期价格区间的下二分之一）
    local is_bottom = relative_position <= 0.5
    
    if is_bottom then
        print_log(LOG_LEVEL.DEBUG, nil, "形成相对低点 - 当前价格在前期价格区间的位置: %.2f%%", 
            relative_position * 100)
    end
    
    return is_bottom
end

-- 检查是否跳空高开或涨停
function is_gap_up_or_limit_up(kdata, idx, stock_code)
    -- 只检查涨停，不再检查跳空高开
    -- 根据股票代码判断涨停幅度
    local limit_up_rate = 0.098  -- 默认A股涨停9.8%
    
    -- 科创板(688开头)和创业板(300开头)涨停20%
    if string.match(stock_code, "^688") or string.match(stock_code, "^300") then
        limit_up_rate = 0.198
    end
    
    -- 检查是否涨停
    local prev_close = kdata[idx-1].close
    local current_close = kdata[idx].close
    
    -- 防止除以0
    if prev_close <= 0 then
        return false
    end
    
    local price_change = (current_close - prev_close) / prev_close
    
    -- 检查计算结果是否有效
    if not is_valid_number(price_change) then
        return false
    end
    
    local is_limit_up = price_change >= limit_up_rate
    
    if is_limit_up then
        print_log(LOG_LEVEL.DEBUG, nil, "涨停板 - 涨幅: %.2f%%, 涨停幅度: %.2f%%", 
            price_change * 100, limit_up_rate * 100)
    end
    
    return is_limit_up
end

-- 检查触底信号
function check_bottom_signal(kdata, indicators, doji_idx, stock_code)
    if not kdata or not indicators then
        log(string.format("[ERROR] [%s] 输入数据无效", stock_code))
        return nil
    end

    if not stock_code then
        log("[ERROR] 股票代码为空")
        return nil
    end

    -- 检查索引是否有效
    if doji_idx < 1 or doji_idx > #kdata then
        log(string.format("[ERROR] [%s] 十字星索引无效: %d", stock_code, doji_idx))
        return nil
    end

    local result = {
        is_bottom = false,
        reasons = {}
    }
    
    -- 检查是否跳空高开或涨停
    if is_gap_up_or_limit_up(kdata, doji_idx, stock_code) then
        log(string.format("[DEBUG] [%s] 检测到跳空高开或涨停，不符合触底特征", stock_code))
        return result
    end
    
    -- 检查均线数据有效性
    if not indicators.ma_long or not indicators.ma_long[doji_idx] then
        log(string.format("[ERROR] [%s] 长期均线数据无效", stock_code))
        return nil
    end
    
    -- 检查是否处于上升趋势或价格在均线上方
    local trend_result = is_uptrend(kdata, indicators, doji_idx)
    local price_above_ma = kdata[doji_idx].close > indicators.ma_long[doji_idx]
    
    if not trend_result and not price_above_ma then
        log(string.format("[DEBUG] [%s] 不满足趋势条件", stock_code))
        return result
    end
    
    -- 检查是否为相对低点
    if not is_relative_bottom(kdata, doji_idx) then
        log(string.format("[DEBUG] [%s] 不满足相对低点条件", stock_code))
        return result
    end
    
    -- 统计满足的技术指标条件
    local condition_count = 0
    
    -- 检查RSI超卖
    if indicators.rsi and indicators.rsi[doji_idx] then
        if indicators.rsi[doji_idx] <= STRATEGY_PARAMS.RSI_OVERSOLD then
            table.insert(result.reasons, "RSI超卖")
            condition_count = condition_count + 1
            log(string.format("[DEBUG] [%s] RSI超卖条件满足: %.2f", stock_code, indicators.rsi[doji_idx]))
        end
    else
        log(string.format("[ERROR] [%s] RSI数据无效", stock_code))
        return nil
    end
    
    -- 检查MACD底背离
    if indicators.macd and indicators.macd[doji_idx] and indicators.macd[doji_idx-1] then
        local macd_trend = indicators.macd[doji_idx] > indicators.macd[doji_idx-1]
        local price_trend = kdata[doji_idx].low <= kdata[doji_idx-1].low
        
        if macd_trend and price_trend then
            table.insert(result.reasons, "MACD底背离")
            condition_count = condition_count + 1
            log(string.format("[DEBUG] [%s] MACD底背离条件满足", stock_code))
        end
    else
        log(string.format("[ERROR] [%s] MACD数据无效", stock_code))
        return nil
    end
    
    -- 检查量能放大
    if indicators.volume_ma and indicators.volume_ma[doji_idx] and kdata[doji_idx].volume then
        if kdata[doji_idx].volume >= indicators.volume_ma[doji_idx] * STRATEGY_PARAMS.VOLUME_RATIO then
            table.insert(result.reasons, "量能放大")
            condition_count = condition_count + 1
            log(string.format("[DEBUG] [%s] 量能放大条件满足: %.2f倍", stock_code, 
                kdata[doji_idx].volume / indicators.volume_ma[doji_idx]))
        end
    else
        log(string.format("[ERROR] [%s] 成交量数据无效", stock_code))
        return nil
    end
    
    -- 只要满足一个条件就可以
    result.is_bottom = condition_count >= 1
    
    if result.is_bottom then
        log(string.format("[INFO] [%s] 满足触底条件，原因: %s", stock_code, table.concat(result.reasons, "、")))
    else
        log(string.format("[DEBUG] [%s] 不满足触底条件", stock_code))
    end
    
    return result
end

-- 验证上涨趋势
function verify_uptrend(kdata, doji_idx, last_idx)
    -- 计算十字星后的价格涨幅
    local doji_price = kdata[doji_idx].close
    local current_price = kdata[last_idx].close
    local price_change = (current_price - doji_price) / doji_price
    
    print_log(LOG_LEVEL.DEBUG, nil, "验证上涨趋势 - 涨幅: %.2f%%", price_change * 100)
    
    return price_change >= STRATEGY_PARAMS.PRICE_INCREASE
end

-- 分析交易信号
function analyze_trading_signal(stock_code, kdata, indicators)
    if not stock_code then
        log("[ERROR] 股票代码为空")
        return { valid = false }
    end
    
    if not kdata or type(kdata) ~= "table" or #kdata == 0 then
        log(string.format("[ERROR] [%s] K线数据无效", stock_code))
        return { valid = false }
    end
    
    if not indicators or type(indicators) ~= "table" then
        log(string.format("[ERROR] [%s] 技术指标数据无效", stock_code))
        return { valid = false }
    end

    local last_idx = #kdata
    local signal = {
        valid = false,
        doji_idx = 0,
        bottom_info = nil
    }
    
    log(string.format("[DEBUG] [%s] 开始在最近%d天内寻找十字星形态", stock_code, STRATEGY_PARAMS.DOJI_DAYS))
    
    -- 确保有足够的数据进行分析
    if last_idx < STRATEGY_PARAMS.DOJI_DAYS + 1 then
        log(string.format("[ERROR] [%s] 数据不足以进行分析", stock_code))
        return signal
    end
    
    -- 在最近DOJI_DAYS天内寻找十字星（不包括最后一天）
    for i = last_idx - 1, math.max(last_idx - STRATEGY_PARAMS.DOJI_DAYS, 2), -1 do
        -- 检查数据有效性
        if not kdata[i] or not kdata[i].open or not kdata[i].close or 
           not kdata[i].high or not kdata[i].low or not kdata[i].time then
            log(string.format("[ERROR] [%s] 第%d天数据缺失", stock_code, i))
            goto continue
        end
        
        -- 检查数值有效性
        if not is_valid_number(kdata[i].open) or not is_valid_number(kdata[i].close) or
           not is_valid_number(kdata[i].high) or not is_valid_number(kdata[i].low) then
            log(string.format("[ERROR] [%s] 第%d天数据包含无效值", stock_code, i))
            goto continue
        end

        if is_doji(kdata[i].open, kdata[i].close, kdata[i].high, kdata[i].low) then
            log(string.format("[DEBUG] [%s] 在%s发现十字星形态", stock_code, kdata[i].time))
            
            -- 检查触底信号
            local bottom_info = check_bottom_signal(kdata, indicators, i, stock_code)
            if not bottom_info then
                log(string.format("[ERROR] [%s] 检查触底信号时发生错误", stock_code))
                goto continue
            end
            
            -- 验证上涨趋势
            if bottom_info.is_bottom then
                local uptrend_valid = verify_uptrend(kdata, i, last_idx)
                if uptrend_valid then
                    signal.valid = true
                    signal.doji_idx = i
                    signal.bottom_info = bottom_info
                    log(string.format("[INFO] [%s] 确认有效的触底反转信号", stock_code))
                    break
                end
            end
        end
        
        ::continue::
    end
    
    if not signal.valid then
        log(string.format("[DEBUG] [%s] 未找到有效的交易信号", stock_code))
    end
    
    return signal
end

-- 生成信号消息
function generate_signal_message(signal, stock_code, kdata)
    local doji_date = kdata[signal.doji_idx].time
    local reasons = table.concat(signal.bottom_info.reasons, "、")
    local price_change = (kdata[#kdata].close - kdata[signal.doji_idx].close) / kdata[signal.doji_idx].close * 100
    
    return string.format(
        "%s出现十字星触底信号，底部特征：%s，后续上涨%.2f%%",
        doji_date,
        reasons,
        price_change
    )
end 