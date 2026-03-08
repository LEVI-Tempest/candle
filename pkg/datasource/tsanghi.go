// Tsanghi data source - 曾海数据源
// API: https://tsanghi.com/api/fin/stock/{exchange}/daily
package datasource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

// TsanghiExchange 交易所代码 | Exchange code
// XSHE = 深圳, XSHG = 上海
const (
	TsanghiXSHE = "XSHE" // 深圳 | Shenzhen
	TsanghiXSHG = "XSHG" // 上海 | Shanghai
)

// TsanghiClient fetches historical daily OHLCV from tsanghi.com
// TsanghiClient 从 tsanghi.com 获取历史日线 OHLCV
type TsanghiClient struct {
	BaseURL string // default: https://tsanghi.com
	Token   string // API token; "demo" for demo
	Client  *http.Client
}

// TsanghiDailyItem API 返回的单条日线
type TsanghiDailyItem struct {
	Ticker string  `json:"ticker"`
	Date   string  `json:"date"`   // YYYY-MM-DD
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// TsanghiResponse API 响应
type TsanghiResponse struct {
	Code int                `json:"code"`
	Msg  string             `json:"msg"`
	Data []TsanghiDailyItem `json:"data"`
}

// NewTsanghiClient creates a Tsanghi API client
// 创建 Tsanghi API 客户端
func NewTsanghiClient(token string) *TsanghiClient {
	if token == "" {
		token = "demo"
	}
	return &TsanghiClient{
		BaseURL: "https://tsanghi.com",
		Token:   token,
		Client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// Fetch fetches daily candlestick data for the given ticker and exchange
// Fetch 获取指定交易所和代码的日线数据
func (c *TsanghiClient) Fetch(exchange, ticker string, opts *FetchOptions) ([]*v1.Candlestick, error) {
	if opts == nil {
		opts = &FetchOptions{Limit: 60, Order: 2}
	}
	limit := opts.Limit
	if limit <= 0 {
		limit = 60
	}
	order := opts.Order
	if order != 1 && order != 2 {
		order = 2
	}

	u := fmt.Sprintf("%s/api/fin/stock/%s/daily", c.BaseURL, exchange)
	reqURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	q := reqURL.Query()
	q.Set("token", c.Token)
	q.Set("ticker", ticker)
	q.Set("order", fmt.Sprintf("%d", order))
	q.Set("limit", fmt.Sprintf("%d", limit))
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

	var tr TsanghiResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}
	if tr.Code != 200 {
		return nil, fmt.Errorf("api error: code=%d msg=%s", tr.Code, tr.Msg)
	}

	// Convert to proto Candlestick
	// 转换为 proto Candlestick
	result := make([]*v1.Candlestick, 0, len(tr.Data))
	for _, item := range tr.Data {
		ts, err := parseDate(item.Date)
		if err != nil {
			continue
		}
		result = append(result, &v1.Candlestick{
			Timestamp: ts.Unix(),
			Open:      item.Open,
			High:      item.High,
			Low:       item.Low,
			Close:     item.Close,
			Volume:    item.Volume,
		})
	}
	return result, nil
}

func parseDate(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", s, time.Local)
}
