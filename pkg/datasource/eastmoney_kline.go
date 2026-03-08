// East Money K-line data source - 东方财富 K 线数据源
// API: http://push2his.eastmoney.com/api/qt/stock/kline/get
// 支持 A 股、港股（secid 格式：124.00700 港股腾讯）
package datasource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

// EastMoneyKlineClient 东方财富 K 线客户端
type EastMoneyKlineClient struct {
	BaseURL string
	Client  *http.Client
}

// EastMoneyKlineItem API 返回的 K 线数据（f51-f61 字段）
type EastMoneyKlineItem struct {
	// 原始字符串格式: "2026-03-06,180.20,182.00,178.50,181.80,..."
	// f51:日期 f52:开盘 f53:收盘 f54:最高 f55:最低 f56:成交量 f57:成交额 f58:振幅 f59:涨跌幅 f60:涨跌额 f61:换手率
	Raw string `json:"-"` // 使用 klines 数组的字符串解析
}

// EastMoneyKlineResponse API 响应
type EastMoneyKlineResponse struct {
	Data *struct {
		Code   string   `json:"code"`
		Name   string   `json:"name"`
		Klines []string `json:"klines"`
	} `json:"data"`
}

// NewEastMoneyKlineClient 创建 East Money K 线客户端
func NewEastMoneyKlineClient() *EastMoneyKlineClient {
	return &EastMoneyKlineClient{
		BaseURL: "http://push2his.eastmoney.com",
		Client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// SecIDHK 构建港股 secid，例如 00700 -> 124.00700（东方财富港股市场代码 124）
func SecIDHK(code string) string {
	code = strings.TrimSpace(code)
	// 补齐为 5 位
	for len(code) < 5 {
		code = "0" + code
	}
	return "124." + code
}

// FetchHK fetches HK stock daily kline from East Money
// 从东方财富获取港股日线数据
func (c *EastMoneyKlineClient) FetchHK(code string, limit int) ([]*v1.Candlestick, error) {
	if limit <= 0 {
		limit = 30
	}
	secid := SecIDHK(code)

	u := fmt.Sprintf("%s/api/qt/stock/kline/get", c.BaseURL)
	reqURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	q := reqURL.Query()
	q.Set("secid", secid)
	q.Set("fields1", "f1,f2,f3,f4,f5,f6")
	q.Set("fields2", "f51,f52,f53,f54,f55,f56,f57")
	q.Set("klt", "101")
	q.Set("fqt", "0")
	q.Set("beg", "0")
	q.Set("end", "20500000")
	q.Set("lmt", strconv.Itoa(limit))
	q.Set("ut", "7eea3edcaed734bea9cbfc24409ed989")
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Referer", "https://quote.eastmoney.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Candle/1.0)")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	var er EastMoneyKlineResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}
	if er.Data == nil || len(er.Data.Klines) == 0 {
		return nil, nil
	}

	result := make([]*v1.Candlestick, 0, len(er.Data.Klines))
	for _, s := range er.Data.Klines {
		// 格式: "2026-03-06,180.20,182.00,178.50,181.80,1234567,..."
		parts := strings.Split(s, ",")
		if len(parts) < 6 {
			continue
		}
		t, err := time.ParseInLocation("2006-01-02", parts[0], time.Local)
		if err != nil {
			continue
		}
		open, _ := strconv.ParseFloat(parts[1], 64)
		closeVal, _ := strconv.ParseFloat(parts[2], 64)
		high, _ := strconv.ParseFloat(parts[3], 64)
		low, _ := strconv.ParseFloat(parts[4], 64)
		vol := 0.0
		if len(parts) > 5 {
			vol, _ = strconv.ParseFloat(parts[5], 64)
		}
		result = append(result, &v1.Candlestick{
			Timestamp: t.Unix(),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closeVal,
			Volume:    vol,
		})
	}
	return result, nil
}
