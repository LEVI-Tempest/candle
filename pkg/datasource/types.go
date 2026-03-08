// Package datasource provides data sources for candlestick charts.
// 数据源包 - 为蜡烛图提供数据获取接口
package datasource

import v1 "github.com/LEVI-Tempest/Candle/pkg/proto"

// Fetcher fetches candlestick data from external APIs
// Fetcher 从外部 API 获取蜡烛图数据
type Fetcher interface {
	// Fetch returns candlestick data for the given symbol
	// Fetch 返回指定标的的蜡烛图数据
	Fetch(symbol string, opts *FetchOptions) ([]*v1.Candlestick, error)
}

// FetchOptions configures fetch behavior
// FetchOptions 配置获取行为
type FetchOptions struct {
	// Limit max number of candles to return; 0 = default (e.g. 60)
	// Limit 最大返回蜡烛数量；0 表示使用默认值（如 60）
	Limit int
	// Order: 1=asc, 2=desc (oldest first / newest first)
	// Order: 1=升序, 2=降序（最旧在前/最新在前）
	Order int
}
