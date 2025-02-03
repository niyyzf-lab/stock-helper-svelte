package market

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"stock-helper-svelte/backend/api/types"
	"time"
)

// Client 市场数据客户端
type Client struct {
	baseURL string
	licence string
	request func(ctx context.Context, endpoint string, freq types.KLineFreq) ([]byte, error)
}

// NewClient 创建新的市场数据客户端
func NewClient(baseURL, licence string, requestFunc func(ctx context.Context, endpoint string, freq types.KLineFreq) ([]byte, error)) *Client {
	return &Client{
		baseURL: baseURL,
		licence: licence,
		request: requestFunc,
	}
}

// GetIndexList 获取指数列表
func (c *Client) GetIndexList(ctx context.Context) ([]types.Index, error) {
	endpoint := "hslt/list"
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取指数列表失败: %v", err)
	}

	var indices []types.Index
	if err := json.Unmarshal(body, &indices); err != nil {
		return nil, fmt.Errorf("解析指数列表失败: %v", err)
	}

	return indices, nil
}

// GetKLineData 获取K线数据
func (c *Client) GetKLineData(ctx context.Context, code string, freq types.KLineFreq) ([]types.KLineData, error) {
	endpoint := fmt.Sprintf("hszbl/fsjy/%s/%s", code, freq)
	body, err := c.request(ctx, endpoint, freq)
	if err != nil {
		return nil, fmt.Errorf("获取K线数据失败: %v", err)
	}

	var klineData []types.KLineData
	if err := json.Unmarshal(body, &klineData); err != nil {
		return nil, fmt.Errorf("解析K线数据失败: %v", err)
	}

	sortKLineByDate(klineData)
	return klineData, nil
}

// GetRealtimeData 获取实时交易数据
func (c *Client) GetRealtimeData(ctx context.Context, code string) (*types.RealtimeData, error) {
	endpoint := fmt.Sprintf("hsrl/ssjy/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取实时数据失败: %v", err)
	}

	var realtimeData types.RealtimeData
	if err := json.Unmarshal(body, &realtimeData); err != nil {
		return nil, fmt.Errorf("解析实时数据失败: %v", err)
	}

	return &realtimeData, nil
}

// GetHistoricalTransactions 获取历史成交分布数据
func (c *Client) GetHistoricalTransactions(ctx context.Context, code string) ([]types.HistoricalTransaction, error) {
	endpoint := fmt.Sprintf("hsmy/lscj/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取历史成交数据失败: %v", err)
	}

	var transactions []types.HistoricalTransaction
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, fmt.Errorf("解析历史成交数据失败: %v", err)
	}

	// 按时间倒序排序（新->旧）
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Time > transactions[j].Time
	})

	return transactions, nil
}

// sortKLineByDate 按日期排序K线数据(旧->新)
func sortKLineByDate(data []types.KLineData) {
	sort.Slice(data, func(i, j int) bool {
		timeI, err := time.Parse("2006-01-02", data[i].Time)
		if err != nil {
			return false
		}
		timeJ, err := time.Parse("2006-01-02", data[j].Time)
		if err != nil {
			return false
		}
		return timeI.Before(timeJ)
	})
}

// GetMainForceMinute 获取主力资金分钟走势
func (c *Client) GetMainForceMinute(ctx context.Context, code string) ([]types.MainForceMinute, error) {
	endpoint := fmt.Sprintf("hsmy/zlzj/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取主力资金分钟走势失败: %v", err)
	}

	var data []types.MainForceMinute
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析主力资金分钟走势失败: %v", err)
	}

	return data, nil
}

// GetCapitalFlow 获取资金流向趋势
func (c *Client) GetCapitalFlow(ctx context.Context, code string) ([]types.CapitalFlow, error) {
	endpoint := fmt.Sprintf("hsmy/zjlr/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取资金流向趋势失败: %v", err)
	}

	var data []types.CapitalFlow
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析资金流向趋势失败: %v", err)
	}

	return data, nil
}

// GetRecentCapitalFlow 获取最近10天资金流向趋势
func (c *Client) GetRecentCapitalFlow(ctx context.Context, code string) ([]types.CapitalFlow, error) {
	endpoint := fmt.Sprintf("hsmy/zhlrt/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取最近资金流向趋势失败: %v", err)
	}

	var data []types.CapitalFlow
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析最近资金流向趋势失败: %v", err)
	}

	return data, nil
}

// GetMainForcePhase 获取阶段主力动向
func (c *Client) GetMainForcePhase(ctx context.Context, code string) ([]types.MainForcePhase, error) {
	endpoint := fmt.Sprintf("hsmy/jddx/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取阶段主力动向失败: %v", err)
	}

	var data []types.MainForcePhase
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析阶段主力动向失败: %v", err)
	}

	return data, nil
}

// GetRecentMainForcePhase 获取最近10天阶段主力动向
func (c *Client) GetRecentMainForcePhase(ctx context.Context, code string) ([]types.MainForcePhase, error) {
	endpoint := fmt.Sprintf("hsmy/jddxt/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取最近阶段主力动向失败: %v", err)
	}

	var data []types.MainForcePhase
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析最近阶段主力动向失败: %v", err)
	}

	return data, nil
}
