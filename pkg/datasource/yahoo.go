// Yahoo Finance data source - 雅虎财经数据源
// API: https://query1.finance.yahoo.com/v8/finance/chart/{TICKER}
// 免费、无需 token，港股格式如 0700.HK
package datasource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

// YahooClient 雅虎财经 K 线客户端
type YahooClient struct {
	BaseURL string
	Client  *http.Client
}

// YahooChartResponse API 响应
type YahooChartResponse struct {
	Chart *struct {
		Result []*struct {
			Timestamp []int64 `json:"timestamp"`
			Indicators *struct {
				Quote []*struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []float64 `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

// NewYahooClient 创建 Yahoo Finance 客户端
func NewYahooClient() *YahooClient {
	return &YahooClient{
		BaseURL: "https://query1.finance.yahoo.com",
		Client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// TickerHK 港股代码转为 Yahoo 格式，如 00700 -> 0700.HK
func TickerHK(code string) string {
	code = strings.TrimSpace(code)
	code = strings.TrimLeft(code, "0")
	if code == "" {
		code = "0"
	}
	for len(code) < 4 {
		code = "0" + code
	}
	return code + ".HK"
}

// FetchHK 获取港股日线数据
func (c *YahooClient) FetchHK(code string, limit int) ([]*v1.Candlestick, error) {
	if limit <= 0 {
		limit = 30
	}
	ticker := TickerHK(code)

	end := time.Now()
	start := end.AddDate(0, 0, -limit-10)
	period1 := start.Unix()
	period2 := end.Unix()

	u := fmt.Sprintf("%s/v8/finance/chart/%s", c.BaseURL, url.PathEscape(ticker))
	reqURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	q := reqURL.Query()
	q.Set("period1", fmt.Sprintf("%d", period1))
	q.Set("period2", fmt.Sprintf("%d", period2))
	q.Set("interval", "1d")
	q.Set("events", "history")
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Candle/1.0)")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	var yr YahooChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&yr); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}
	if yr.Chart == nil || len(yr.Chart.Result) == 0 || yr.Chart.Result[0].Indicators == nil {
		return nil, nil
	}

	r := yr.Chart.Result[0]
	q0 := r.Indicators.Quote
	if len(q0) == 0 {
		return nil, nil
	}
	qt := q0[0]

	result := make([]*v1.Candlestick, 0, len(r.Timestamp))
	for i, ts := range r.Timestamp {
		if i >= len(qt.Close) || qt.Close[i] == 0 {
			continue
		}
		o := qt.Close[i]
		if i < len(qt.Open) {
			o = qt.Open[i]
		}
		h := qt.Close[i]
		if i < len(qt.High) {
			h = qt.High[i]
		}
		l := qt.Close[i]
		if i < len(qt.Low) {
			l = qt.Low[i]
		}
		vol := 0.0
		if i < len(qt.Volume) {
			vol = qt.Volume[i]
		}
		result = append(result, &v1.Candlestick{
			Timestamp: ts,
			Open:      o,
			High:      h,
			Low:       l,
			Close:     qt.Close[i],
			Volume:    vol,
		})
	}

	if len(result) > limit {
		result = result[len(result)-limit:]
	}
	return result, nil
}
