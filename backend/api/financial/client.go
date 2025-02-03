package financial

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-helper-svelte/backend/api/types"
)

// Client 财务数据客户端
type Client struct {
	baseURL string
	licence string
	request func(ctx context.Context, endpoint string, freq types.KLineFreq) ([]byte, error)
}

// NewClient 创建新的财务数据客户端
func NewClient(baseURL, licence string, requestFunc func(ctx context.Context, endpoint string, freq types.KLineFreq) ([]byte, error)) *Client {
	return &Client{
		baseURL: baseURL,
		licence: licence,
		request: requestFunc,
	}
}

// GetProfitability 获取盈利能力数据
func (c *Client) GetProfitability(ctx context.Context, year int, quarter int) ([]types.ProfitabilityData, error) {
	endpoint := fmt.Sprintf("hicw/yl/%d/%d", year, quarter)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取盈利能力数据失败: %v", err)
	}

	var data []types.ProfitabilityData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析盈利能力数据失败: %v", err)
	}

	return data, nil
}

// GetOperation 获取运营能力数据
func (c *Client) GetOperation(ctx context.Context, year int, quarter int) ([]types.OperationData, error) {
	endpoint := fmt.Sprintf("hicw/yy/%d/%d", year, quarter)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取运营能力数据失败: %v", err)
	}

	var data []types.OperationData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析运营能力数据失败: %v", err)
	}

	return data, nil
}

// GetGrowth 获取成长能力数据
func (c *Client) GetGrowth(ctx context.Context, year int, quarter int) ([]types.GrowthData, error) {
	endpoint := fmt.Sprintf("hicw/cz/%d/%d", year, quarter)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取成长能力数据失败: %v", err)
	}

	var data []types.GrowthData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析成长能力数据失败: %v", err)
	}

	return data, nil
}

// GetSolvency 获取偿债能力数据
func (c *Client) GetSolvency(ctx context.Context, year int, quarter int) ([]types.SolvencyData, error) {
	endpoint := fmt.Sprintf("hicw/cznl/%d/%d", year, quarter)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取偿债能力数据失败: %v", err)
	}

	var data []types.SolvencyData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析偿债能力数据失败: %v", err)
	}

	return data, nil
}

// GetCashFlow 获取现金流量数据
func (c *Client) GetCashFlow(ctx context.Context, year int, quarter int) ([]types.CashFlowData, error) {
	endpoint := fmt.Sprintf("hicw/xj/%d/%d", year, quarter)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取现金流量数据失败: %v", err)
	}

	var data []types.CashFlowData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析现金流量数据失败: %v", err)
	}

	return data, nil
}

// GetPerformance 获取业绩报表数据
func (c *Client) GetPerformance(ctx context.Context, year int, quarter int) ([]types.PerformanceData, error) {
	endpoint := fmt.Sprintf("hicw/yjbb/%d/%d", year, quarter)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取业绩报表数据失败: %v", err)
	}

	var data []types.PerformanceData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析业绩报表数据失败: %v", err)
	}

	return data, nil
}

// GetStockProfitability 获取指定股票的盈利能力数据
func (c *Client) GetStockProfitability(ctx context.Context, code string, year int, quarter int) (*types.ProfitabilityData, error) {
	data, err := c.GetProfitability(ctx, year, quarter)
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.Dm == code {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("未找到股票代码 %s 的盈利能力数据", code)
}

// GetStockOperation 获取指定股票的运营能力数据
func (c *Client) GetStockOperation(ctx context.Context, code string, year int, quarter int) (*types.OperationData, error) {
	data, err := c.GetOperation(ctx, year, quarter)
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.Dm == code {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("未找到股票代码 %s 的运营能力数据", code)
}

// GetStockGrowth 获取指定股票的成长能力数据
func (c *Client) GetStockGrowth(ctx context.Context, code string, year int, quarter int) (*types.GrowthData, error) {
	data, err := c.GetGrowth(ctx, year, quarter)
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.Dm == code {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("未找到股票代码 %s 的成长能力数据", code)
}

// GetStockSolvency 获取指定股票的偿债能力数据
func (c *Client) GetStockSolvency(ctx context.Context, code string, year int, quarter int) (*types.SolvencyData, error) {
	data, err := c.GetSolvency(ctx, year, quarter)
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.Dm == code {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("未找到股票代码 %s 的偿债能力数据", code)
}

// GetStockCashFlow 获取指定股票的现金流量数据
func (c *Client) GetStockCashFlow(ctx context.Context, code string, year int, quarter int) (*types.CashFlowData, error) {
	data, err := c.GetCashFlow(ctx, year, quarter)
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.Dm == code {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("未找到股票代码 %s 的现金流量数据", code)
}

// GetStockPerformance 获取指定股票的业绩报表数据
func (c *Client) GetStockPerformance(ctx context.Context, code string, year int, quarter int) (*types.PerformanceData, error) {
	data, err := c.GetPerformance(ctx, year, quarter)
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.Dm == code {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("未找到股票代码 %s 的业绩报表数据", code)
}

// GetStockFinancialData 获取指定股票的所有财务数据
func (c *Client) GetStockFinancialData(ctx context.Context, code string, year int, quarter int) (*types.StockFinancialData, error) {
	profitability, _ := c.GetStockProfitability(ctx, code, year, quarter)
	operation, _ := c.GetStockOperation(ctx, code, year, quarter)
	growth, _ := c.GetStockGrowth(ctx, code, year, quarter)
	solvency, _ := c.GetStockSolvency(ctx, code, year, quarter)
	cashFlow, _ := c.GetStockCashFlow(ctx, code, year, quarter)
	performance, _ := c.GetStockPerformance(ctx, code, year, quarter)

	return &types.StockFinancialData{
		Code:          code,
		Year:          year,
		Quarter:       quarter,
		Profitability: profitability,
		Operation:     operation,
		Growth:        growth,
		Solvency:      solvency,
		CashFlow:      cashFlow,
		Performance:   performance,
	}, nil
}
