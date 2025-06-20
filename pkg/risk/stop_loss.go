package risk

import "math"

// StopLossStrategy 止损策略接口
type StopLossStrategy interface {
	CalculateStopLoss(entryPrice, currentPrice float64, params map[string]float64) float64
}

// FixedStopLoss 固定止损
type FixedStopLoss struct {
	ratio float64
}

func NewFixedStopLoss(ratio float64) *FixedStopLoss {
	return &FixedStopLoss{ratio: ratio}
}

func (f *FixedStopLoss) CalculateStopLoss(entryPrice, currentPrice float64, params map[string]float64) float64 {
	return entryPrice * (1 - f.ratio)
}

// TrailingStopLoss 移动止损
type TrailingStopLoss struct {
	ratio float64
}

func NewTrailingStopLoss(ratio float64) *TrailingStopLoss {
	return &TrailingStopLoss{ratio: ratio}
}

func (t *TrailingStopLoss) CalculateStopLoss(entryPrice, currentPrice float64, params map[string]float64) float64 {
	highestPrice := params["highest"]
	if highestPrice == 0 {
		highestPrice = entryPrice
	}
	
	if currentPrice > highestPrice {
		highestPrice = currentPrice
	}
	
	return math.Max(entryPrice*(1-t.ratio), highestPrice*(1-t.ratio))
}

// ATRStopLoss ATR止损
type ATRStopLoss struct {
	multiplier float64
}

func NewATRStopLoss(multiplier float64) *ATRStopLoss {
	return &ATRStopLoss{multiplier: multiplier}
}

func (a *ATRStopLoss) CalculateStopLoss(entryPrice, currentPrice float64, params map[string]float64) float64 {
	atr := params["atr"]
	if atr == 0 {
		return entryPrice * 0.98 // 默认2%止损
	}
	return entryPrice - (atr * a.multiplier)
} 