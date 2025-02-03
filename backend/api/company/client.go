package company

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-helper-svelte/backend/api/types"
)

// Client 公司信息客户端
type Client struct {
	baseURL string
	licence string
	request func(ctx context.Context, endpoint string, freq types.KLineFreq) ([]byte, error)
}

// NewClient 创建新的公司信息客户端
func NewClient(baseURL, licence string, requestFunc func(ctx context.Context, endpoint string, freq types.KLineFreq) ([]byte, error)) *Client {
	return &Client{
		baseURL: baseURL,
		licence: licence,
		request: requestFunc,
	}
}

// GetCompanyInfo 获取公司基本信息
func (c *Client) GetCompanyInfo(ctx context.Context, code string) (*types.CompanyInfo, error) {
	endpoint := fmt.Sprintf("hscp/gsjj/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取公司信息失败: %v", err)
	}

	var info types.CompanyInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("解析公司信息失败: %v", err)
	}

	return &info, nil
}

// GetFinancialIndicators 获取财务指标
func (c *Client) GetFinancialIndicators(ctx context.Context, code string) ([]types.FinancialIndicator, error) {
	endpoint := fmt.Sprintf("hscp/cwzb/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取财务指标失败: %v", err)
	}

	var indicators []types.FinancialIndicator
	if err := json.Unmarshal(body, &indicators); err != nil {
		return nil, fmt.Errorf("解析财务指标失败: %v", err)
	}

	return indicators, nil
}

// GetShareholders 获取股东信息
func (c *Client) GetShareholders(ctx context.Context, code string) ([]types.ShareholderInfo, error) {
	endpoint := fmt.Sprintf("hscp/sdgd/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取股东信息失败: %v", err)
	}

	var shareholders []types.ShareholderInfo
	if err := json.Unmarshal(body, &shareholders); err != nil {
		return nil, fmt.Errorf("解析股东信息失败: %v", err)
	}

	return shareholders, nil
}

// GetNews 获取公司新闻
func (c *Client) GetNews(ctx context.Context, code string) ([]types.CompanyNews, error) {
	endpoint := fmt.Sprintf("hscp/xwxx/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取公司新闻失败: %v", err)
	}

	var news []types.CompanyNews
	if err := json.Unmarshal(body, &news); err != nil {
		return nil, fmt.Errorf("解析公司新闻失败: %v", err)
	}

	return news, nil
}

// GetAnnouncements 获取公司公告
func (c *Client) GetAnnouncements(ctx context.Context, code string) ([]types.CompanyAnnouncement, error) {
	endpoint := fmt.Sprintf("hscp/ggxx/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取公司公告失败: %v", err)
	}

	var announcements []types.CompanyAnnouncement
	if err := json.Unmarshal(body, &announcements); err != nil {
		return nil, fmt.Errorf("解析公司公告失败: %v", err)
	}

	return announcements, nil
}

// GetCompanyProfile 获取公司详细信息
func (c *Client) GetCompanyProfile(ctx context.Context, code string) (*types.CompanyProfile, error) {
	endpoint := fmt.Sprintf("hscp/gsjj/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取公司详细信息失败: %v", err)
	}

	var profile types.CompanyProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("解析公司详细信息失败: %v", err)
	}

	return &profile, nil
}

// GetBelongingIndices 获取所属指数信息
func (c *Client) GetBelongingIndices(ctx context.Context, code string) ([]types.BelongingIndex, error) {
	endpoint := fmt.Sprintf("hscp/sszs/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取所属指数信息失败: %v", err)
	}

	var indices []types.BelongingIndex
	if err := json.Unmarshal(body, &indices); err != nil {
		return nil, fmt.Errorf("解析所属指数信息失败: %v", err)
	}

	return indices, nil
}

// GetExecutives 获取高管成员信息
func (c *Client) GetExecutives(ctx context.Context, code string) ([]types.Executive, error) {
	endpoint := fmt.Sprintf("hscp/ljgg/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取高管成员信息失败: %v", err)
	}

	var executives []types.Executive
	if err := json.Unmarshal(body, &executives); err != nil {
		return nil, fmt.Errorf("解析高管成员信息失败: %v", err)
	}

	return executives, nil
}

// GetDirectors 获取董事会成员信息
func (c *Client) GetDirectors(ctx context.Context, code string) ([]types.Director, error) {
	endpoint := fmt.Sprintf("hscp/ljds/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取董事会成员信息失败: %v", err)
	}

	var directors []types.Director
	if err := json.Unmarshal(body, &directors); err != nil {
		return nil, fmt.Errorf("解析董事会成员信息失败: %v", err)
	}

	return directors, nil
}

// GetSupervisors 获取监事会成员信息
func (c *Client) GetSupervisors(ctx context.Context, code string) ([]types.Supervisor, error) {
	endpoint := fmt.Sprintf("hscp/ljjj/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取监事会成员信息失败: %v", err)
	}

	var supervisors []types.Supervisor
	if err := json.Unmarshal(body, &supervisors); err != nil {
		return nil, fmt.Errorf("解析监事会成员信息失败: %v", err)
	}

	return supervisors, nil
}

// GetDividendHistory 获取分红历史
func (c *Client) GetDividendHistory(ctx context.Context, code string) ([]types.Dividend, error) {
	endpoint := fmt.Sprintf("hscp/jnfh/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取分红历史失败: %v", err)
	}

	var dividends []types.Dividend
	if err := json.Unmarshal(body, &dividends); err != nil {
		return nil, fmt.Errorf("解析分红历史失败: %v", err)
	}

	return dividends, nil
}

// GetAdditionalIssues 获取增发信息
func (c *Client) GetAdditionalIssues(ctx context.Context, code string) ([]types.AdditionalIssue, error) {
	endpoint := fmt.Sprintf("hscp/jnzf/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取增发信息失败: %v", err)
	}

	var issues []types.AdditionalIssue
	if err := json.Unmarshal(body, &issues); err != nil {
		return nil, fmt.Errorf("解析增发信息失败: %v", err)
	}

	return issues, nil
}

// GetShareUnlocks 获取解禁信息
func (c *Client) GetShareUnlocks(ctx context.Context, code string) ([]types.ShareUnlock, error) {
	endpoint := fmt.Sprintf("hscp/jjxs/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取解禁信息失败: %v", err)
	}

	var unlocks []types.ShareUnlock
	if err := json.Unmarshal(body, &unlocks); err != nil {
		return nil, fmt.Errorf("解析解禁信息失败: %v", err)
	}

	return unlocks, nil
}

// GetQuarterlyCashFlow 获取季度现金流
func (c *Client) GetQuarterlyCashFlow(ctx context.Context, code string) ([]types.QuarterlyCashFlow, error) {
	endpoint := fmt.Sprintf("hscp/jdxj/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取季度现金流失败: %v", err)
	}

	var cashflows []types.QuarterlyCashFlow
	if err := json.Unmarshal(body, &cashflows); err != nil {
		return nil, fmt.Errorf("解析季度现金流失败: %v", err)
	}

	return cashflows, nil
}

// GetTopTenShareholders 获取十大股东信息
func (c *Client) GetTopTenShareholders(ctx context.Context, code string) ([]types.TopTenShareholder, error) {
	endpoint := fmt.Sprintf("hscp/sdgd/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取十大股东信息失败: %v", err)
	}

	var shareholders []types.TopTenShareholder
	if err := json.Unmarshal(body, &shareholders); err != nil {
		return nil, fmt.Errorf("解析十大股东信息失败: %v", err)
	}

	return shareholders, nil
}

// GetTopTenFloatShareholders 获取十大流通股东信息
func (c *Client) GetTopTenFloatShareholders(ctx context.Context, code string) ([]types.TopTenFloatShareholder, error) {
	endpoint := fmt.Sprintf("hscp/ltgd/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取十大流通股东信息失败: %v", err)
	}

	var shareholders []types.TopTenFloatShareholder
	if err := json.Unmarshal(body, &shareholders); err != nil {
		return nil, fmt.Errorf("解析十大流通股东信息失败: %v", err)
	}

	return shareholders, nil
}

// GetFundHoldings 获取基金持股信息
func (c *Client) GetFundHoldings(ctx context.Context, code string) ([]types.FundHolding, error) {
	endpoint := fmt.Sprintf("hscp/jjcg/%s", code)
	body, err := c.request(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("获取基金持股信息失败: %v", err)
	}

	var holdings []types.FundHolding
	if err := json.Unmarshal(body, &holdings); err != nil {
		return nil, fmt.Errorf("解析基金持股信息失败: %v", err)
	}

	return holdings, nil
}
