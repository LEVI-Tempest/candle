// East Money data source - 东方财富数据源
// API: http://push2.eastmoney.com/api/qt/stock/get
// 适用于单日实时/当日数据，返回最新价、OHLC 等
package datasource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

// EastMoneyClient fetches real-time/single-day data from East Money
// EastMoneyClient 从东方财富获取实时/单日数据
type EastMoneyClient struct {
	BaseURL string
	Client  *http.Client
}

// EastMoneyStockItem API 返回的股票数据
type EastMoneyStockItem struct {
	F43 string `json:"f43"` // 最新价 latest
	F57 string `json:"f57"` // 涨跌幅 pctChange
	F58 string `json:"f58"` // 涨跌额 change
	F59 string `json:"f59"` // 开盘价 open
	F60 string `json:"f60"` // 最高价 high
	F61 string `json:"f61"` // 最低价 low
}

// EastMoneyResponse API 响应
type EastMoneyResponse struct {
	Data EastMoneyStockItem `json:"data"`
}

// NewEastMoneyClient creates an East Money API client
// 创建东方财富 API 客户端
func NewEastMoneyClient() *EastMoneyClient {
	return &EastMoneyClient{
		BaseURL: "http://push2.eastmoney.com",
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// SecID builds East Money secid from market and code
// SecID 从市场和代码构建东方财富 secid
// market: 1=上海, 0=深圳
func SecID(market int, code string) string {
	return fmt.Sprintf("%d.%s", market, code)
}

// FetchOne fetches real-time (single candle) data for the given secid
// FetchOne 获取指定 secid 的实时（单根蜡烛）数据
func (c *EastMoneyClient) FetchOne(secid string) (*v1.Candlestick, error) {
	u := fmt.Sprintf("%s/api/qt/stock/get", c.BaseURL)
	reqURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	q := reqURL.Query()
	q.Set("secid", secid)
	q.Set("fields", "f43,f57,f58,f59,f60,f61")
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	var er EastMoneyResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	d := er.Data
	open, _ := strconv.ParseFloat(d.F59, 64)
	high, _ := strconv.ParseFloat(d.F60, 64)
	low, _ := strconv.ParseFloat(d.F61, 64)
	closeVal, _ := strconv.ParseFloat(d.F43, 64)

	if open == 0 && high == 0 && low == 0 && closeVal == 0 {
		return nil, fmt.Errorf("no valid data for secid %s", secid)
	}

	ts := time.Now()
	return &v1.Candlestick{
		Timestamp: ts.Unix(),
		Open:      open,
		High:      high,
		Low:       low,
		Close:     closeVal,
		Volume:    0,
	}, nil
}
