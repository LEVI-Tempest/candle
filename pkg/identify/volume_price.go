package identify

// VolumePriceSignal 量价信号
type VolumePriceSignal struct {
	Type      string  // 信号类型 | Signal type
	Direction string  // 方向: bullish/bearish/neutral | Direction
	Position  int     // K线位置索引 | Candle index
	Strength  float64 // 信号强度 0-1 | Strength score
	Price     float64 // 触发价格 | Trigger price
	Volume    float64 // 触发成交量 | Trigger volume
	Reason    string  // 触发原因 | Trigger reason
}

// AnalyzeVolumePriceSignals analyzes volume-price behavior with lightweight rules.
// AnalyzeVolumePriceSignals 使用轻量规则分析量价关系。
func AnalyzeVolumePriceSignals(cs []CandlestickWrapper, lookback int) []VolumePriceSignal {
	if len(cs) == 0 {
		return nil
	}
	if lookback < 3 {
		lookback = 5
	}

	signals := make([]VolumePriceSignal, 0)
	for i := range cs {
		c := cs[i]

		avgVol := averageVolumeBefore(cs, i, lookback)
		if avgVol <= 0 {
			continue
		}
		volRatio := c.Volume / avgVol

		// 放量阳线确认 | Bullish confirmation with high volume
		if c.IsBullish() && volRatio >= 1.8 {
			signals = append(signals, VolumePriceSignal{
				Type:      "Volume Breakout Bullish",
				Direction: "bullish",
				Position:  i,
				Strength:  0.85,
				Price:     c.Close,
				Volume:    c.Volume,
				Reason:    "bullish candle with >=1.8x average volume",
			})
		}

		// 放量阴线确认 | Bearish confirmation with high volume
		if c.IsBearish() && volRatio >= 1.8 {
			signals = append(signals, VolumePriceSignal{
				Type:      "Volume Breakdown Bearish",
				Direction: "bearish",
				Position:  i,
				Strength:  0.85,
				Price:     c.Close,
				Volume:    c.Volume,
				Reason:    "bearish candle with >=1.8x average volume",
			})
		}

		// 长腿十字+放量：拐点预警 | Long-legged doji with high volume: turning-point warning
		if LongLeggedDoji([]CandlestickWrapper{c}) && volRatio >= 1.5 {
			signals = append(signals, VolumePriceSignal{
				Type:      "High-Volume Long-Legged Doji",
				Direction: "neutral",
				Position:  i,
				Strength:  0.9,
				Price:     c.Close,
				Volume:    c.Volume,
				Reason:    "long-legged doji with >=1.5x average volume",
			})
		}

		// 简单双向背离预警（连续3根） | Simple divergence warning (3-bar sequence)
		if i >= 2 {
			c0, c1, c2 := cs[i-2], cs[i-1], cs[i]
			priceUp := c0.Close < c1.Close && c1.Close < c2.Close
			priceDown := c0.Close > c1.Close && c1.Close > c2.Close
			volDown := c0.Volume > c1.Volume && c1.Volume > c2.Volume

			if priceUp && volDown {
				signals = append(signals, VolumePriceSignal{
					Type:      "Bearish Volume Divergence",
					Direction: "bearish",
					Position:  i,
					Strength:  0.8,
					Price:     c2.Close,
					Volume:    c2.Volume,
					Reason:    "price up 3 bars while volume down 3 bars",
				})
			}
			if priceDown && volDown {
				signals = append(signals, VolumePriceSignal{
					Type:      "Downtrend Volume Exhaustion",
					Direction: "neutral",
					Position:  i,
					Strength:  0.75,
					Price:     c2.Close,
					Volume:    c2.Volume,
					Reason:    "price down 3 bars with shrinking volume",
				})
			}
		}
	}
	return signals
}

func averageVolumeBefore(cs []CandlestickWrapper, i, lookback int) float64 {
	if i <= 0 {
		return 0
	}
	start := i - lookback
	if start < 0 {
		start = 0
	}
	if start == i {
		return 0
	}
	sum := 0.0
	for j := start; j < i; j++ {
		sum += cs[j].Volume
	}
	return sum / float64(i-start)
}
