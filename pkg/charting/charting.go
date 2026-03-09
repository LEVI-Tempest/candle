package charting

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/LEVI-Tempest/Candle/pkg/identify"
	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

/***
 * @author Tempest
 * @description identify
 * @date 2025/05/18
 * @ref: go-echarts Examples：https://github.com/go-echarts/examples
 * @version 1.0.0
 */

type klineData struct {
	date string
	data [4]float32
}

var kd = [...]klineData{
	{date: "2018/1/24", data: [4]float32{2320.26, 2320.26, 2287.3, 2362.94}},
	{date: "2018/1/25", data: [4]float32{2300, 2291.3, 2288.26, 2308.38}},
	{date: "2018/1/28", data: [4]float32{2295.35, 2346.5, 2295.35, 2346.92}},
	{date: "2018/1/29", data: [4]float32{2347.22, 2358.98, 2337.35, 2363.8}},
	{date: "2018/1/30", data: [4]float32{2360.75, 2382.48, 2347.89, 2383.76}},
	{date: "2018/1/31", data: [4]float32{2383.43, 2385.42, 2371.23, 2391.82}},
	{date: "2018/2/1", data: [4]float32{2377.41, 2419.02, 2369.57, 2421.15}},
	{date: "2018/2/4", data: [4]float32{2425.92, 2428.15, 2417.58, 2440.38}},
	{date: "2018/2/5", data: [4]float32{2411, 2433.13, 2403.3, 2437.42}},
	{date: "2018/2/6", data: [4]float32{2432.68, 2434.48, 2427.7, 2441.73}},
	{date: "2018/2/7", data: [4]float32{2430.69, 2418.53, 2394.22, 2433.89}},
	{date: "2018/2/8", data: [4]float32{2416.62, 2432.4, 2414.4, 2443.03}},
	{date: "2018/2/18", data: [4]float32{2441.91, 2421.56, 2415.43, 2444.8}},
	{date: "2018/2/19", data: [4]float32{2420.26, 2382.91, 2373.53, 2427.07}},
	{date: "2018/2/20", data: [4]float32{2383.49, 2397.18, 2370.61, 2397.94}},
	{date: "2018/2/21", data: [4]float32{2378.82, 2325.95, 2309.17, 2378.82}},
	{date: "2018/2/22", data: [4]float32{2322.94, 2314.16, 2308.76, 2330.88}},
	{date: "2018/2/25", data: [4]float32{2320.62, 2325.82, 2315.01, 2338.78}},
	{date: "2018/2/26", data: [4]float32{2313.74, 2293.34, 2289.89, 2340.71}},
	{date: "2018/2/27", data: [4]float32{2297.77, 2313.22, 2292.03, 2324.63}},
	{date: "2018/2/28", data: [4]float32{2322.32, 2365.59, 2308.92, 2366.16}},
	{date: "2018/3/1", data: [4]float32{2364.54, 2359.51, 2330.86, 2369.65}},
	{date: "2018/3/4", data: [4]float32{2332.08, 2273.4, 2259.25, 2333.54}},
	{date: "2018/3/5", data: [4]float32{2274.81, 2326.31, 2270.1, 2328.14}},
	{date: "2018/3/6", data: [4]float32{2333.61, 2347.18, 2321.6, 2351.44}},
	{date: "2018/3/7", data: [4]float32{2340.44, 2324.29, 2304.27, 2352.02}},
	{date: "2018/3/8", data: [4]float32{2326.42, 2318.61, 2314.59, 2333.67}},
	{date: "2018/3/11", data: [4]float32{2314.68, 2310.59, 2296.58, 2320.96}},
	{date: "2018/3/12", data: [4]float32{2309.16, 2286.6, 2264.83, 2333.29}},
	{date: "2018/3/13", data: [4]float32{2282.17, 2263.97, 2253.25, 2286.33}},
	{date: "2018/3/14", data: [4]float32{2255.77, 2270.28, 2253.31, 2276.22}},
	{date: "2018/3/15", data: [4]float32{2269.31, 2278.4, 2250, 2312.08}},
	{date: "2018/3/18", data: [4]float32{2267.29, 2240.02, 2239.21, 2276.05}},
	{date: "2018/3/19", data: [4]float32{2244.26, 2257.43, 2232.02, 2261.31}},
	{date: "2018/3/20", data: [4]float32{2257.74, 2317.37, 2257.42, 2317.86}},
	{date: "2018/3/21", data: [4]float32{2318.21, 2324.24, 2311.6, 2330.81}},
	{date: "2018/3/22", data: [4]float32{2321.4, 2328.28, 2314.97, 2332}},
	{date: "2018/3/25", data: [4]float32{2334.74, 2326.72, 2319.91, 2344.89}},
	{date: "2018/3/26", data: [4]float32{2318.58, 2297.67, 2281.12, 2319.99}},
	{date: "2018/3/27", data: [4]float32{2299.38, 2301.26, 2289, 2323.48}},
	{date: "2018/3/28", data: [4]float32{2273.55, 2236.3, 2232.91, 2273.55}},
	{date: "2018/3/29", data: [4]float32{2238.49, 2236.62, 2228.81, 2246.87}},
	{date: "2018/4/1", data: [4]float32{2229.46, 2234.4, 2227.31, 2243.95}},
	{date: "2018/4/2", data: [4]float32{2234.9, 2227.74, 2220.44, 2253.42}},
	{date: "2018/4/3", data: [4]float32{2232.69, 2225.29, 2217.25, 2241.34}},
	{date: "2018/4/8", data: [4]float32{2196.24, 2211.59, 2180.67, 2212.59}},
	{date: "2018/4/9", data: [4]float32{2215.47, 2225.77, 2215.47, 2234.73}},
	{date: "2018/4/10", data: [4]float32{2224.93, 2226.13, 2212.56, 2233.04}},
	{date: "2018/4/11", data: [4]float32{2236.98, 2219.55, 2217.26, 2242.48}},
	{date: "2018/4/12", data: [4]float32{2218.09, 2206.78, 2204.44, 2226.26}},
	{date: "2018/4/15", data: [4]float32{2199.91, 2181.94, 2177.39, 2204.99}},
	{date: "2018/4/16", data: [4]float32{2169.63, 2194.85, 2165.78, 2196.43}},
	{date: "2018/4/17", data: [4]float32{2195.03, 2193.8, 2178.47, 2197.51}},
	{date: "2018/4/18", data: [4]float32{2181.82, 2197.6, 2175.44, 2206.03}},
	{date: "2018/4/19", data: [4]float32{2201.12, 2244.64, 2200.58, 2250.11}},
	{date: "2018/4/22", data: [4]float32{2236.4, 2242.17, 2232.26, 2245.12}},
	{date: "2018/4/23", data: [4]float32{2242.62, 2184.54, 2182.81, 2242.62}},
	{date: "2018/4/24", data: [4]float32{2187.35, 2218.32, 2184.11, 2226.12}},
	{date: "2018/4/25", data: [4]float32{2213.19, 2199.31, 2191.85, 2224.63}},
	{date: "2018/4/26", data: [4]float32{2203.89, 2177.91, 2173.86, 2210.58}},
	{date: "2018/5/2", data: [4]float32{2170.78, 2174.12, 2161.14, 2179.65}},
	{date: "2018/5/3", data: [4]float32{2179.05, 2205.5, 2179.05, 2222.81}},
	{date: "2018/5/6", data: [4]float32{2212.5, 2231.17, 2212.5, 2236.07}},
	{date: "2018/5/7", data: [4]float32{2227.86, 2235.57, 2219.44, 2240.26}},
	{date: "2018/5/8", data: [4]float32{2242.39, 2246.3, 2235.42, 2255.21}},
	{date: "2018/5/9", data: [4]float32{2246.96, 2232.97, 2221.38, 2247.86}},
	{date: "2018/5/10", data: [4]float32{2228.82, 2246.83, 2225.81, 2247.67}},
	{date: "2018/5/13", data: [4]float32{2247.68, 2241.92, 2231.36, 2250.85}},
	{date: "2018/5/14", data: [4]float32{2238.9, 2217.01, 2205.87, 2239.93}},
	{date: "2018/5/15", data: [4]float32{2217.09, 2224.8, 2213.58, 2225.19}},
	{date: "2018/5/16", data: [4]float32{2221.34, 2251.81, 2210.77, 2252.87}},
	{date: "2018/5/17", data: [4]float32{2249.81, 2282.87, 2248.41, 2288.09}},
	{date: "2018/5/20", data: [4]float32{2286.33, 2299.99, 2281.9, 2309.39}},
	{date: "2018/5/21", data: [4]float32{2297.11, 2305.11, 2290.12, 2305.3}},
	{date: "2018/5/22", data: [4]float32{2303.75, 2302.4, 2292.43, 2314.18}},
	{date: "2018/5/23", data: [4]float32{2293.81, 2275.67, 2274.1, 2304.95}},
	{date: "2018/5/24", data: [4]float32{2281.45, 2288.53, 2270.25, 2292.59}},
	{date: "2018/5/27", data: [4]float32{2286.66, 2293.08, 2283.94, 2301.7}},
	{date: "2018/5/28", data: [4]float32{2293.4, 2321.32, 2281.47, 2322.1}},
	{date: "2018/5/29", data: [4]float32{2323.54, 2324.02, 2321.17, 2334.33}},
	{date: "2018/5/30", data: [4]float32{2316.25, 2317.75, 2310.49, 2325.72}},
	{date: "2018/5/31", data: [4]float32{2320.74, 2300.59, 2299.37, 2325.53}},
	{date: "2018/6/3", data: [4]float32{2300.21, 2299.25, 2294.11, 2313.43}},
	{date: "2018/6/4", data: [4]float32{2297.1, 2272.42, 2264.76, 2297.1}},
	{date: "2018/6/5", data: [4]float32{2270.71, 2270.93, 2260.87, 2276.86}},
	{date: "2018/6/6", data: [4]float32{2264.43, 2242.11, 2240.07, 2266.69}},
	{date: "2018/6/7", data: [4]float32{2242.26, 2210.9, 2205.07, 2250.63}},
	{date: "2018/6/13", data: [4]float32{2190.1, 2148.35, 2126.22, 2190.1}},
}

func klineBase() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Kline-example",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y)
	return kline
}

func klineDataZoomInside() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "DataZoom(inside)",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y)
	return kline
}

func klineDataZoomBoth() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "DataZoom(inside&slider)",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y)
	return kline
}

func klineDataZoomYAxis() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "DataZoom(yAxis)",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			YAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y)
	return kline
}

func klineStyle() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "different style",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y).
		SetSeriesOptions(
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "highest value",
				Type:     "max",
				ValueDim: "highest",
			}),
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "lowest value",
				Type:     "min",
				ValueDim: "lowest",
			}),
			charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
				Label: &opts.Label{
					Show: opts.Bool(true),
				},
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:        "#ec0000",
				Color0:       "#00da3c",
				BorderColor:  "#8A0000",
				BorderColor0: "#008F28",
			}),
		)
	return kline
}

type KlineExamples struct{}

func (KlineExamples) Examples() {
	page := components.NewPage()
	page.AddCharts(
		klineBase(),
		klineDataZoomInside(),
		klineDataZoomBoth(),
		klineDataZoomYAxis(),
		klineStyle(),
	)

	f, err := os.Create("kline.html")
	if err != nil {
		panic(err)

	}
	page.Render(io.MultiWriter(f))
}

// Enhanced Kline Chart with Pattern Recognition
// 带有形态识别的增强K线图

// TimeFrame represents different time periods for candlestick charts
// 支持不同时间周期的切换
type TimeFrame string

const (
	TimeFrame1Min  TimeFrame = "1m"
	TimeFrame5Min  TimeFrame = "5m"
	TimeFrame15Min TimeFrame = "15m"
	TimeFrame1Hour TimeFrame = "1h"
	TimeFrame1Day  TimeFrame = "1d"
	TimeFrame1Week TimeFrame = "1w"
)

// Kline represents a candlestick chart with additional functionality
// K线图结构，包含额外功能
type Kline struct {
	chart      *charts.Kline
	data       []identify.CandlestickWrapper
	currentPos int
	timeFrame  TimeFrame
}

// TrendType represents the type of trend
// 趋势类型
type TrendType string

const (
	TrendTypeUp   TrendType = "up"
	TrendTypeDown TrendType = "down"
)

// LevelType represents support or resistance level
// 支撑阻力位类型
type LevelType string

const (
	LevelTypeSupport    LevelType = "support"
	LevelTypeResistance LevelType = "resistance"
)

// Point represents a coordinate point
// 坐标点
type Point struct {
	X float64
	Y float64
}

// EnhancedKline represents an enhanced candlestick chart with additional features
// 优化后的K线数据结构
type EnhancedKline struct {
	*charts.Kline
	Patterns          []Pattern                     // Detected patterns (识别出的形态)
	VolumeSignals     []identify.VolumePriceSignal  // Volume-price signals (量价信号)
	Evidences         []identify.PatternEvidence    // Structured pattern evidences (结构化证据)
	Indicators        map[string][]float64          // Technical indicators (技术指标)
	TrendLines        []TrendLine                   // Trend lines (趋势线)
	SupportResistance []Level                       // Support and resistance levels (支撑阻力位)
	TimeFrame         TimeFrame                     // Current time frame (当前时间周期)
	Data              []identify.CandlestickWrapper // Candlestick data (蜡烛图数据)
}

// Pattern represents information about a detected candlestick pattern
// 形态信息
type Pattern struct {
	Type     string  // Pattern type (形态类型)
	Position int     // Position in the data (位置)
	Strength float64 // Pattern strength (形态强度)
	Risk     float64 // Risk level (风险等级)
	Price    float64 // Price at which pattern was detected (检测到形态的价格)
	Time     string  // Time when pattern was detected (检测到形态的时间)
}

// TrendLine represents a trend line on the chart
// 趋势线
type TrendLine struct {
	StartPoint Point     // Starting point (起始点)
	EndPoint   Point     // Ending point (结束点)
	Type       TrendType // Trend type: up/down (上升/下降)
}

// Level represents a support or resistance level
// 支撑阻力位
type Level struct {
	Price    float64   // Price level (价格水平)
	Type     LevelType // Level type: support/resistance (支撑/阻力)
	Strength float64   // Level strength (强度)
}

// NewKline creates a new Kline instance
// 创建新的K线实例
func NewKline() *Kline {
	return &Kline{
		chart:     charts.NewKLine(),
		data:      make([]identify.CandlestickWrapper, 0),
		timeFrame: TimeFrame1Day,
	}
}

// NewEnhancedKline creates a new EnhancedKline instance
// 创建新的增强K线实例
func NewEnhancedKline() *EnhancedKline {
	return &EnhancedKline{
		Kline:             charts.NewKLine(),
		Patterns:          make([]Pattern, 0),
		VolumeSignals:     make([]identify.VolumePriceSignal, 0),
		Evidences:         make([]identify.PatternEvidence, 0),
		Indicators:        make(map[string][]float64),
		TrendLines:        make([]TrendLine, 0),
		SupportResistance: make([]Level, 0),
		TimeFrame:         TimeFrame1Day,
		Data:              make([]identify.CandlestickWrapper, 0),
	}
}

// LoadData loads candlestick data into the enhanced kline chart
// 加载蜡烛图数据到增强K线图
func (ek *EnhancedKline) LoadData(candles []*v1.Candlestick) {
	ek.Data = make([]identify.CandlestickWrapper, len(candles))
	for i, candle := range candles {
		ek.Data[i] = identify.NewCandlestickWrapper(candle)
	}
}

// AutoDetectPatterns automatically detects candlestick patterns in the data
// 自动识别K线形态
func (ek *EnhancedKline) AutoDetectPatterns() {
	patterns := make([]Pattern, 0)

	// Single candlestick patterns (单根K线形态)
	for i := 0; i < len(ek.Data); i++ {
		candle := []identify.CandlestickWrapper{ek.Data[i]}
		timestamp := time.Unix(ek.Data[i].Timestamp, 0).Format("2006-01-02 15:04:05")

		// Check single candle patterns
		if identify.Doji(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Doji",
				Position: i,
				Strength: 0.7,
				Risk:     0.5,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.LongLeggedDoji(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Long-Legged Doji",
				Position: i,
				Strength: 0.85,
				Risk:     0.5,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.Hammer(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Hammer",
				Position: i,
				Strength: 0.8,
				Risk:     0.3,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.HangingMan(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Hanging Man",
				Position: i,
				Strength: 0.8,
				Risk:     0.7,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.InvertedHammer(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Inverted Hammer",
				Position: i,
				Strength: 0.7,
				Risk:     0.4,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.ShootingStar(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Shooting Star",
				Position: i,
				Strength: 0.8,
				Risk:     0.6,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.Marubozu(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Marubozu",
				Position: i,
				Strength: 0.9,
				Risk:     0.2,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}
		if identify.WhiteMarubozu(candle) {
			patterns = append(patterns, Pattern{
				Type:     "White Marubozu",
				Position: i,
				Strength: 0.95,
				Risk:     0.2,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}
		if identify.BlackMarubozu(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Black Marubozu",
				Position: i,
				Strength: 0.95,
				Risk:     0.2,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.SpinningTop(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Spinning Top",
				Position: i,
				Strength: 0.5,
				Risk:     0.8,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.Umbrella(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Umbrella",
				Position: i,
				Strength: 0.7,
				Risk:     0.4,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.DragonflyDoji(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Dragonfly Doji",
				Position: i,
				Strength: 0.8,
				Risk:     0.4,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.GravestoneDoji(candle) {
			patterns = append(patterns, Pattern{
				Type:     "Gravestone Doji",
				Position: i,
				Strength: 0.8,
				Risk:     0.5,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}
	}

	// Two candlestick patterns (双根K线形态)
	for i := 1; i < len(ek.Data); i++ {
		candles := []identify.CandlestickWrapper{ek.Data[i], ek.Data[i-1]}
		timestamp := time.Unix(ek.Data[i].Timestamp, 0).Format("2006-01-02 15:04:05")

		if identify.BullishEngulfing(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Bullish Engulfing",
				Position: i,
				Strength: 0.9,
				Risk:     0.2,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.BearishEngulfing(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Bearish Engulfing",
				Position: i,
				Strength: 0.9,
				Risk:     0.2,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.PiercingLine(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Piercing Line",
				Position: i,
				Strength: 0.8,
				Risk:     0.3,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.DarkCloudCover(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Dark Cloud Cover",
				Position: i,
				Strength: 0.8,
				Risk:     0.3,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.TweezerBottoms(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Tweezer Bottoms",
				Position: i,
				Strength: 0.7,
				Risk:     0.4,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.TweezerTops(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Tweezer Tops",
				Position: i,
				Strength: 0.7,
				Risk:     0.4,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.FallingWindow(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Falling Window",
				Position: i,
				Strength: 0.6,
				Risk:     0.5,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.RisingWindow(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Rising Window",
				Position: i,
				Strength: 0.6,
				Risk:     0.5,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.BullishHarami(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Bullish Harami",
				Position: i,
				Strength: 0.8,
				Risk:     0.3,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.BearishHarami(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Bearish Harami",
				Position: i,
				Strength: 0.8,
				Risk:     0.3,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}
	}

	// Three candlestick patterns (三根K线形态)
	for i := 2; i < len(ek.Data); i++ {
		candles := []identify.CandlestickWrapper{ek.Data[i], ek.Data[i-1], ek.Data[i-2]}
		timestamp := time.Unix(ek.Data[i].Timestamp, 0).Format("2006-01-02 15:04:05")

		if identify.MorningStar(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Morning Star",
				Position: i,
				Strength: 0.9,
				Risk:     0.1,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.EveningStar(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Evening Star",
				Position: i,
				Strength: 0.9,
				Risk:     0.1,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}
		if identify.MorningDojiStar(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Morning Doji Star",
				Position: i,
				Strength: 0.98,
				Risk:     0.1,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}
		if identify.EveningDojiStar(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Evening Doji Star",
				Position: i,
				Strength: 0.98,
				Risk:     0.1,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.ThreeWhiteSoldiers(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Three White Soldiers",
				Position: i,
				Strength: 0.95,
				Risk:     0.1,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.ThreeBlackCrows(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Three Black Crows",
				Position: i,
				Strength: 0.95,
				Risk:     0.1,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}
	}

	// Five candlestick patterns (五根K线形态)
	for i := 4; i < len(ek.Data); i++ {
		candles := []identify.CandlestickWrapper{
			ek.Data[i], ek.Data[i-1], ek.Data[i-2], ek.Data[i-3], ek.Data[i-4],
		}
		timestamp := time.Unix(ek.Data[i].Timestamp, 0).Format("2006-01-02 15:04:05")

		if identify.RisingThreeMethods(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Rising Three Methods",
				Position: i,
				Strength: 0.9,
				Risk:     0.2,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}

		if identify.FallingThreeMethods(candles) {
			patterns = append(patterns, Pattern{
				Type:     "Falling Three Methods",
				Position: i,
				Strength: 0.9,
				Risk:     0.2,
				Price:    ek.Data[i].Close,
				Time:     timestamp,
			})
		}
	}

	ek.Patterns = patterns
	// Analyze volume-price signals after pattern detection
	// 形态识别后补充量价信号分析
	ek.VolumeSignals = identify.AnalyzeVolumePriceSignals(ek.Data, 5)
	ek.Evidences = identify.BuildPatternEvidence(
		toPatternSignals(ek.Patterns),
		ek.Data,
		identify.DefaultEvidenceConfig(),
	)
}

func toPatternSignals(patterns []Pattern) []identify.PatternSignal {
	out := make([]identify.PatternSignal, 0, len(patterns))
	for _, p := range patterns {
		out = append(out, identify.PatternSignal{
			Type:      p.Type,
			Direction: patternDirection(p.Type),
			Position:  p.Position,
			Strength:  p.Strength,
			Risk:      p.Risk,
			Price:     p.Price,
			Time:      p.Time,
		})
	}
	return out
}

func patternDirection(patternType string) string {
	switch patternType {
	case "Hammer",
		"Inverted Hammer",
		"Bullish Engulfing",
		"Piercing Line",
		"Morning Star",
		"Three White Soldiers",
		"Tweezer Bottoms",
		"Rising Window",
		"Rising Three Methods":
		return "bullish"
	case "Hanging Man",
		"Shooting Star",
		"Bearish Engulfing",
		"Dark Cloud Cover",
		"Evening Star",
		"Three Black Crows",
		"Tweezer Tops",
		"Falling Window",
		"Falling Three Methods":
		return "bearish"
	default:
		return "neutral"
	}
}

// MarkPatterns marks detected patterns on the chart
// 在图表上标记检测到的形态
func (ek *EnhancedKline) MarkPatterns() {
	if len(ek.Patterns) == 0 {
		return
	}

	// Create precise coordinate mark points (limit to top 10 for readability)
	// 为检测到的形态创建精确坐标标记，最多 10 个以减轻拥挤
	displayPatterns := getDisplayPatterns(ek.Patterns, 0.85)
	if len(displayPatterns) > 10 {
		displayPatterns = displayPatterns[:10]
	}
	markPointData := make([]opts.MarkPointNameCoordItem, 0, len(displayPatterns))

	for _, pattern := range displayPatterns {
		if pattern.Position < 0 || pattern.Position >= len(ek.Data) {
			continue
		}

		date := time.Unix(ek.Data[pattern.Position].Timestamp, 0).Format("2006-01-02")
		offsetFactor := 1.004

		markPointData = append(markPointData, opts.MarkPointNameCoordItem{
			Name:       fmt.Sprintf("%s %s", getPatternDirectionTag(pattern.Type), getPatternShortName(pattern.Type)),
			Coordinate: []interface{}{date, pattern.Price * offsetFactor},
			Value:      fmt.Sprintf("%.2f", pattern.Price),
			Symbol:     getDirectionSymbol(pattern.Type),
			SymbolSize: 10,
			ItemStyle: &opts.ItemStyle{
				Color:       getPatternColor(pattern.Type),
				BorderColor: getPatternColor(pattern.Type),
				BorderWidth: 1,
				Opacity:     0.9,
			},
		})
	}

	// Add all pattern marks into the candlestick series
	// 将所有形态标记点添加到蜡烛图序列
	if len(markPointData) > 0 {
		ek.Kline.SetSeriesOptions(
			charts.WithMarkPointNameCoordItemOpts(markPointData...),
			charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
				Label: &opts.Label{
					Show:      opts.Bool(true),
					Position:  "top",
					FontSize:  9,
					Color:     "#222222",
					Formatter: "{b}",
				},
			}),
		)
	}
}

// getPatternColor returns the color for a specific pattern type
// 获取特定形态类型的颜色
func getPatternColor(patternType string) string {
	switch patternType {
	// Bullish patterns (看涨形态)
	case "Hammer", "Inverted Hammer", "Bullish Engulfing", "Piercing Line", "Morning Star", "Morning Doji Star",
		"Three White Soldiers", "Dragonfly Doji", "Bullish Harami", "Rising Three Methods", "White Marubozu":
		return "#00da3c" // Green
	// Bearish patterns (看跌形态)
	case "Hanging Man", "Shooting Star", "Bearish Engulfing", "Dark Cloud Cover", "Evening Star", "Evening Doji Star",
		"Three Black Crows", "Gravestone Doji", "Bearish Harami", "Falling Three Methods", "Black Marubozu":
		return "#ec0000" // Red
	// Neutral/Reversal patterns (中性/反转形态)
	case "Doji", "Spinning Top", "Tweezer Tops", "Tweezer Bottoms":
		return "#ffaa00" // Orange
	case "Long-Legged Doji":
		return "#ff8800" // Deep Orange
	// Gap patterns (缺口形态)
	case "Rising Window", "Falling Window":
		return "#0066cc" // Blue
	// Strong patterns (强势形态)
	case "Marubozu":
		return "#9900cc" // Purple
	default:
		return "#666666" // Gray
	}
}

// getPatternSymbol returns the symbol for a specific pattern type
// 获取特定形态类型的符号
func getPatternSymbol(patternType string) string {
	switch patternType {
	// Bullish patterns (看涨形态)
	case "Hammer", "Inverted Hammer", "Bullish Engulfing", "Piercing Line", "Morning Star", "Morning Doji Star",
		"Three White Soldiers", "Dragonfly Doji", "Bullish Harami", "Rising Three Methods", "White Marubozu":
		return "triangle" // Triangle pointing up
	// Bearish patterns (看跌形态)
	case "Hanging Man", "Shooting Star", "Bearish Engulfing", "Dark Cloud Cover", "Evening Star", "Evening Doji Star",
		"Three Black Crows", "Gravestone Doji", "Bearish Harami", "Falling Three Methods", "Black Marubozu":
		return "triangleDown" // Triangle pointing down
	// Neutral/Reversal patterns (中性/反转形态)
	case "Doji", "Spinning Top":
		return "diamond" // Diamond
	case "Long-Legged Doji":
		return "circle"
	// Support/Resistance patterns (支撑/阻力形态)
	case "Tweezer Tops", "Tweezer Bottoms":
		return "rect" // Rectangle
	// Gap patterns (缺口形态)
	case "Rising Window", "Falling Window":
		return "arrow" // Arrow
	// Strong patterns (强势形态)
	case "Marubozu":
		return "star" // Star
	default:
		return "circle" // Circle
	}
}

// getPatternDirectionTag returns a clear direction marker for each pattern
// getPatternDirectionTag 返回每个形态的方向标记（看涨/看跌/中性）
func getPatternDirectionTag(patternType string) string {
	switch patternType {
	// Bullish patterns (看涨形态)
	case "Hammer", "Inverted Hammer", "Bullish Engulfing", "Piercing Line", "Morning Star", "Morning Doji Star",
		"Three White Soldiers", "Dragonfly Doji", "Bullish Harami", "Rising Three Methods", "Rising Window", "White Marubozu":
		return "↑看涨"
	// Bearish patterns (看跌形态)
	case "Hanging Man", "Shooting Star", "Bearish Engulfing", "Dark Cloud Cover", "Evening Star", "Evening Doji Star",
		"Three Black Crows", "Gravestone Doji", "Bearish Harami", "Falling Three Methods", "Falling Window", "Black Marubozu":
		return "↓看跌"
	default:
		return "→中性"
	}
}

func getDirectionSymbol(patternType string) string {
	switch getPatternDirectionTag(patternType) {
	case "↑看涨":
		return "triangle"
	case "↓看跌":
		return "triangleDown"
	default:
		return "circle"
	}
}

func getPatternShortName(patternType string) string {
	switch patternType {
	case "Doji":
		return "十字星"
	case "Dragonfly Doji":
		return "蜻蜓十字"
	case "Gravestone Doji":
		return "墓碑十字"
	case "Spinning Top":
		return "纺锤线"
	case "Morning Star":
		return "晨星"
	case "Evening Star":
		return "暮星"
	case "Morning Doji Star":
		return "晨十字星"
	case "Evening Doji Star":
		return "暮十字星"
	case "Hammer":
		return "锤头"
	case "Hanging Man":
		return "上吊线"
	case "Inverted Hammer":
		return "倒锤头"
	case "Shooting Star":
		return "流星"
	case "Three White Soldiers":
		return "三白兵"
	case "Three Black Crows":
		return "三黑鸦"
	case "Bullish Engulfing":
		return "看涨吞没"
	case "Bearish Engulfing":
		return "看跌吞没"
	case "Piercing Line":
		return "刺透"
	case "Dark Cloud Cover":
		return "乌云盖顶"
	case "Bullish Harami":
		return "看涨孕线"
	case "Bearish Harami":
		return "看跌孕线"
	case "Rising Three Methods":
		return "上升三法"
	case "Falling Three Methods":
		return "下降三法"
	case "Rising Window":
		return "上升窗口"
	case "Falling Window":
		return "下降窗口"
	case "White Marubozu":
		return "光头阳线"
	case "Black Marubozu":
		return "光头阴线"
	case "Long-Legged Doji":
		return "长腿十字"
	default:
		return patternType
	}
}

// getDisplayPatterns keeps one strongest signal per candle for readability.
// getDisplayPatterns 为每根K线只保留一个最强信号，提升可读性。
func getDisplayPatterns(patterns []Pattern, minStrength float64) []Pattern {
	bestByPos := make(map[int]Pattern)
	for _, p := range patterns {
		if p.Strength < minStrength {
			continue
		}
		old, ok := bestByPos[p.Position]
		if !ok || p.Strength > old.Strength {
			bestByPos[p.Position] = p
		}
	}
	out := make([]Pattern, 0, len(bestByPos))
	for _, p := range bestByPos {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Position < out[j].Position })
	return out
}

// CreateChart creates and configures the candlestick chart with patterns
// 创建并配置带有形态的蜡烛图
func (ek *EnhancedKline) CreateChart(title string) {
	if len(ek.Data) == 0 {
		return
	}

	// Prepare data for chart
	// 为图表准备数据
	x := make([]string, len(ek.Data))
	y := make([]opts.KlineData, len(ek.Data))
	volume := make([]opts.BarData, len(ek.Data))
	volMA5 := make([]opts.LineData, len(ek.Data))
	volMA10 := make([]opts.LineData, len(ek.Data))

	for i, candle := range ek.Data {
		x[i] = time.Unix(candle.Timestamp, 0).Format("2006-01-02")
		y[i] = opts.KlineData{
			Value: []interface{}{candle.Open, candle.Close, candle.Low, candle.High},
		}
		volumeColor := "#ec0000"
		if candle.Close >= candle.Open {
			volumeColor = "#00da3c"
		}
		volume[i] = opts.BarData{
			Value: candle.Volume,
			ItemStyle: &opts.ItemStyle{
				Color: volumeColor,
			},
		}
		volMA5[i] = opts.LineData{Value: movingAvgVolumeAt(ek.Data, i, 5)}
		volMA10[i] = opts.LineData{Value: movingAvgVolumeAt(ek.Data, i, 10)}
	}

	// Create dedicated lower panel for volume (xAxis[1], yAxis[1])
	// 为成交量创建独立下方面板（xAxis[1], yAxis[1]）
	ek.Kline.ExtendXAxis(opts.XAxis{
		Type:      "category",
		GridIndex: 1,
		AxisLabel: &opts.AxisLabel{
			Show: opts.Bool(false),
		},
		AxisTick: &opts.AxisTick{
			Show: opts.Bool(false),
		},
	})
	ek.Kline.ExtendYAxis(opts.YAxis{
		Scale:       opts.Bool(true),
		GridIndex:   1,
		Name:        "Volume",
		SplitNumber: 2,
		AxisLabel: &opts.AxisLabel{
			Show: opts.Bool(true),
		},
	})

	// Configure chart options
	// 配置图表选项
	displayPatterns := getDisplayPatterns(ek.Patterns, 0.85)
	if len(displayPatterns) > 10 {
		displayPatterns = displayPatterns[:10]
	}
	ek.Kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: fmt.Sprintf("Detected %d patterns, showing %d | 点击图例可切换 Volume/VOL MA5/VOL MA10", len(ek.Patterns), len(displayPatterns)),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
			GridIndex:   0,
		}, 0),
		charts.WithYAxisOpts(opts.YAxis{
			Scale:     opts.Bool(true),
			GridIndex: 0,
		}, 0),
		charts.WithGridOpts(
			opts.Grid{
				Left:         "8%",
				Right:        "8%",
				Top:          "10%",
				Height:       "58%",
				ContainLabel: opts.Bool(true),
			},
			opts.Grid{
				Left:         "8%",
				Right:        "8%",
				Top:          "74%",
				Height:       "16%",
				ContainLabel: opts.Bool(true),
			},
		),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0, 1},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0, 1},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    opts.Bool(true),
			Trigger: "axis",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:         opts.Bool(true),
			SelectedMode: "multiple",
			Data:         []string{"Volume", "VOL MA5", "VOL MA10"},
			Top:          "3%",
			Right:        "3%",
		}),
	)

	// Add data to chart
	// 将数据添加到图表
	ek.Kline.SetXAxis(x).AddSeries("Candlestick", y)

	// Mark patterns on candlestick series first, so marks won't leak into volume series
	// 先在蜡烛图序列打标，避免标记污染成交量序列
	ek.MarkPatterns()

	// Overlay volume bars on secondary axis
	// 在副轴叠加成交量柱状图
	volumeBar := charts.NewBar()
	volumeBar.SetXAxis(x).AddSeries("Volume", volume,
		charts.WithBarChartOpts(opts.BarChart{
			XAxisIndex: 1,
			YAxisIndex: 1,
		}),
	)
	ek.Kline.Overlap(volumeBar)

	// Overlay volume moving averages on volume panel
	// 在成交量面板叠加均量线
	volLine := charts.NewLine()
	volLine.SetXAxis(x).
		AddSeries("VOL MA5", volMA5,
			charts.WithLineChartOpts(opts.LineChart{
				XAxisIndex: 1,
				YAxisIndex: 1,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{Width: 1}),
		).
		AddSeries("VOL MA10", volMA10,
			charts.WithLineChartOpts(opts.LineChart{
				XAxisIndex: 1,
				YAxisIndex: 1,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{Width: 1}),
		)
	ek.Kline.Overlap(volLine)
}

// RenderToFile renders the chart to an HTML file
// 将图表渲染到HTML文件
func (ek *EnhancedKline) RenderToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer f.Close()

	return ek.Kline.Render(io.MultiWriter(f))
}

// GetPatternSummary returns a summary of detected patterns
// 获取检测到的形态摘要
func (ek *EnhancedKline) GetPatternSummary() map[string]int {
	summary := make(map[string]int)
	for _, pattern := range ek.Patterns {
		summary[pattern.Type]++
	}
	return summary
}

func movingAvgVolumeAt(data []identify.CandlestickWrapper, idx, n int) float64 {
	if n <= 0 || idx < 0 || len(data) == 0 {
		return 0
	}
	start := idx - n + 1
	if start < 0 {
		start = 0
	}
	sum := 0.0
	count := 0
	for i := start; i <= idx; i++ {
		sum += data[i].Volume
		count++
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}
