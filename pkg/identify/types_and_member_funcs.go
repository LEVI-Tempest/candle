package identify

import (
	"math"

	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
	"github.com/LEVI-Tempest/Candle/pkg/utils"
)

type CandlestickWrapper struct {
	*v1.Candlestick
}

// Body returns the body of the candlestick
// always >= 0
func (c *CandlestickWrapper) Body() float64 {
	return utils.Abs(c.Close - c.Open)
}

func (c *CandlestickWrapper) UpperShadow() float64 {
	if c.Close >= c.Open {
		return c.High - c.Close
	}
	return c.High - c.Open
}

func (c *CandlestickWrapper) LowerShadow() float64 {
	if c.Close >= c.Open {
		return c.Open - c.Low
	}
	return c.Close - c.Low
}

func (c *CandlestickWrapper) Yang() bool {
	return c.Close >= c.Open
}

func (c *CandlestickWrapper) Yin() bool {
	return !c.Yang()
}

// IsBullish Yang
func (c *CandlestickWrapper) IsBullish() bool {
	return c.Yang()
}

// IsBearish Yin
func (c *CandlestickWrapper) IsBearish() bool {
	return c.Yin()
}

// DetermineTrend returns the trend based on simple rising/falling days count
// 基于简单的上涨/下跌天数计数来确定趋势
func DetermineTrend(cs []CandlestickWrapper, days int) Trend {
	if len(cs) <= days {
		return TrendUnknown
	}

	// Count the number of days when closing price is rising in the last days[]
	// 计算最近几天收盘价上涨的天数
	risingCount := 0
	for i := len(cs) - days; i < len(cs)-1; i++ {
		if cs[i+1].Close > cs[i].Close {
			risingCount++
		}
	}

	// Determine the trend based on the majority
	// 基于多数情况确定趋势
	if risingCount > days*3/5 {
		return TrendYang
	}

	return TrendYin
}

// AnalyzeLongTermTrend provides a comprehensive analysis of the long-term trend
// 提供全面的长期趋势分析
// Parameters:
// - cs: Candlestick data series (蜡烛图数据系列)
// - period: Number of days to analyze (分析的天数)
// - volatilityThreshold: Threshold for determining if the market is in a sideways/consolidation pattern (判断市场是否处于震荡/盘整模式的阈值)
//
// Returns:
// - Trend: Yang (上升), Yin (下降), or Middle (震荡)
// - Additional metrics about the trend strength and volatility
func AnalyzeLongTermTrend(cs []CandlestickWrapper, period int, volatilityThreshold float64) (Trend, map[string]float64) {
	if len(cs) < period {
		return TrendUnknown, nil
	}

	// Extract the relevant portion of the data for analysis
	// 提取相关部分的数据进行分析
	analysisData := cs[len(cs)-period:]

	// Calculate key metrics
	// 计算关键指标
	metrics := make(map[string]float64)

	// Starting and ending prices
	// 起始和结束价格
	startPrice := analysisData[0].Close
	endPrice := analysisData[len(analysisData)-1].Close

	// Overall price change and percentage
	// 总体价格变化和百分比
	priceChange := endPrice - startPrice
	priceChangePercent := (priceChange / startPrice) * 100
	metrics["priceChange"] = priceChange
	metrics["priceChangePercent"] = priceChangePercent

	// Calculate high and low during the period
	// 计算期间的最高价和最低价
	highPrice := analysisData[0].High
	lowPrice := analysisData[0].Low

	for _, candle := range analysisData {
		if candle.High > highPrice {
			highPrice = candle.High
		}
		if candle.Low < lowPrice {
			lowPrice = candle.Low
		}
	}

	metrics["highPrice"] = highPrice
	metrics["lowPrice"] = lowPrice

	// Calculate volatility (using high-low range as a percentage of average price)
	// 计算波动性（使用收盘价标准差占均值百分比，降低极值日对震荡识别的干扰）
	closeSum := 0.0
	for _, candle := range analysisData {
		closeSum += candle.Close
	}
	avgClose := closeSum / float64(len(analysisData))
	varianceSum := 0.0
	for _, candle := range analysisData {
		diff := candle.Close - avgClose
		varianceSum += diff * diff
	}
	volatility := 0.0
	if avgClose != 0 {
		volatility = (math.Sqrt(varianceSum/float64(len(analysisData))) / avgClose) * 100
	}
	metrics["volatility"] = volatility

	// Count rising and falling days
	// 计算上涨和下跌的天数
	risingDays := 0
	fallingDays := 0

	for i := 1; i < len(analysisData); i++ {
		if analysisData[i].Close > analysisData[i-1].Close {
			risingDays++
		} else if analysisData[i].Close < analysisData[i-1].Close {
			fallingDays++
		}
	}

	metrics["risingDays"] = float64(risingDays)
	metrics["fallingDays"] = float64(fallingDays)

	// Calculate trend strength
	// 计算趋势强度
	trendStrength := utils.Abs(float64(risingDays-fallingDays)) / float64(period-1) * 100
	metrics["trendStrength"] = trendStrength

	// Determine trend direction
	// 确定趋势方向
	var trend Trend

	// If volatility is low and price change is minimal, it's a sideways/consolidation pattern
	// 如果波动性低且价格变化很小，则为横盘/盘整模式
	if volatility < volatilityThreshold && utils.Abs(priceChangePercent) < volatilityThreshold {
		trend = TrendMiddle
	} else if trendStrength < 20 && utils.Abs(priceChangePercent) < volatilityThreshold*1.5 {
		// Weak directional strength with limited net change is treated as sideways.
		// 方向强度较弱且净涨跌有限时，视为震荡。
		trend = TrendMiddle
	} else if priceChange > 0 && risingDays > fallingDays {
		trend = TrendYang
	} else if priceChange < 0 && fallingDays > risingDays {
		trend = TrendYin
	} else {
		// If the signals are mixed, look at the overall price change
		// 如果信号混合，查看整体价格变化
		if priceChange > 0 {
			trend = TrendYang
		} else if priceChange < 0 {
			trend = TrendYin
		} else {
			trend = TrendMiddle
		}
	}

	return trend, metrics
}
