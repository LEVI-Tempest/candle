package risk

import (
	"math"
	"time"
)

// RiskManager 风险管理器
type RiskManager struct {
	portfolio   *Portfolio
	riskProfile RiskProfile
	config      RiskConfig
	metrics     *RiskMetrics
	alerts      []RiskAlert
}

// NewRiskManager 创建风险管理器
func NewRiskManager(initialCapital float64, config RiskConfig) *RiskManager {
	return &RiskManager{
		portfolio: &Portfolio{
			TotalValue:  initialCapital,
			Cash:        initialCapital,
			Positions:   make(map[string]*Position),
			RiskBudget:  initialCapital * config.MaxTotalRisk,
			PeakValue:   initialCapital,
		},
		riskProfile: RiskProfile{
			Conservative: 0.2,
			Moderate:     0.5,
			Aggressive:   0.8,
		},
		config:  config,
		metrics: &RiskMetrics{},
		alerts:  make([]RiskAlert, 0),
	}
}

// CalculateVolatility 计算波动率
func (rm *RiskManager) CalculateVolatility(prices []float64) float64 {
	if len(prices) < 2 {
		return 0
	}

	// 计算日收益率
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		returns[i-1] = (prices[i] - prices[i-1]) / prices[i-1]
	}

	// 计算收益率标准差
	mean := 0.0
	for _, r := range returns {
		mean += r
	}
	mean /= float64(len(returns))

	variance := 0.0
	for _, r := range returns {
		variance += math.Pow(r-mean, 2)
	}
	variance /= float64(len(returns)-1)

	// 年化波动率 (假设252个交易日)
	return math.Sqrt(variance) * math.Sqrt(252)
}

// CalculateVaR 计算风险价值
func (rm *RiskManager) CalculateVaR(prices []float64, confidence float64) float64 {
	if len(prices) < 2 {
		return 0
	}

	// 计算日收益率
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		returns[i-1] = (prices[i] - prices[i-1]) / prices[i-1]
	}

	// 计算收益率标准差
	mean := 0.0
	for _, r := range returns {
		mean += r
	}
	mean /= float64(len(returns))

	variance := 0.0
	for _, r := range returns {
		variance += math.Pow(r-mean, 2)
	}
	variance /= float64(len(returns)-1)
	stdDev := math.Sqrt(variance)

	// 根据置信度确定分位数
	var zScore float64
	switch confidence {
	case 0.90:
		zScore = 1.282
	case 0.95:
		zScore = 1.645
	case 0.99:
		zScore = 2.326
	default:
		zScore = 1.645
	}

	// VaR = 投资组合价值 × 标准差 × 分位数
	return rm.portfolio.TotalValue * stdDev * zScore
}

// CalculateMaxDrawdown 计算最大回撤
func (rm *RiskManager) CalculateMaxDrawdown(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}

	maxDrawdown := 0.0
	currentDrawdown := 0.0
	peak := values[0]

	for _, value := range values {
		if value > peak {
			peak = value
		}

		drawdown := (peak - value) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
		if drawdown > currentDrawdown {
			currentDrawdown = drawdown
		}
	}

	return maxDrawdown, currentDrawdown
}

// CalculatePositionSize 计算仓位大小 (凯利公式)
func (rm *RiskManager) CalculatePositionSize(winRate, avgWin, avgLoss float64) float64 {
	if avgWin <= 0 || avgLoss <= 0 {
		return 0
	}

	// 凯利公式: f = (bp - q) / b
	// 其中: b = 盈亏比, p = 胜率, q = 败率
	b := avgWin / avgLoss
	p := winRate
	q := 1 - winRate

	kellyFraction := (b*p - q) / b

	// 限制在合理范围内
	if kellyFraction > 0.25 {
		kellyFraction = 0.25 // 最大25%
	}
	if kellyFraction < 0 {
		kellyFraction = 0
	}

	return kellyFraction
}

// CalculateStopLoss 计算止损价格
func (rm *RiskManager) CalculateStopLoss(entryPrice float64, stopLossType string, params map[string]float64) float64 {
	switch stopLossType {
	case "fixed":
		ratio := params["ratio"]
		return entryPrice * (1 - ratio)
	case "atr":
		atr := params["atr"]
		multiplier := params["multiplier"]
		return entryPrice - (atr * multiplier)
	case "trailing":
		ratio := params["ratio"]
		highestPrice := params["highest"]
		return math.Max(entryPrice*(1-ratio), highestPrice*(1-ratio))
	default:
		return entryPrice * (1 - rm.config.StopLossRatio)
	}
}

// CheckRiskLimits 检查风险限制
func (rm *RiskManager) CheckRiskLimits() []RiskAlert {
	var alerts []RiskAlert

	// 检查总风险
	totalRisk := rm.calculateTotalRisk()
	if totalRisk > rm.config.MaxTotalRisk {
		alerts = append(alerts, RiskAlert{
			Type:      "total_risk",
			Message:   "总风险超过限制",
			Level:     "high",
			Timestamp: time.Now(),
			Value:     totalRisk,
			Threshold: rm.config.MaxTotalRisk,
		})
	}

	// 检查最大回撤
	if rm.portfolio.CurrentDrawdown > rm.config.MaxDrawdownLimit {
		alerts = append(alerts, RiskAlert{
			Type:      "max_drawdown",
			Message:   "当前回撤超过限制",
			Level:     "critical",
			Timestamp: time.Now(),
			Value:     rm.portfolio.CurrentDrawdown,
			Threshold: rm.config.MaxDrawdownLimit,
		})
	}

	// 检查单笔仓位
	for symbol, position := range rm.portfolio.Positions {
		positionValue := float64(position.Quantity) * position.CurrentPrice
		positionRatio := positionValue / rm.portfolio.TotalValue

		if positionRatio > rm.config.MaxPositionSize {
			alerts = append(alerts, RiskAlert{
				Type:      "position_size",
				Message:   "单笔仓位过大: " + symbol,
				Level:     "medium",
				Timestamp: time.Now(),
				Value:     positionRatio,
				Threshold: rm.config.MaxPositionSize,
			})
		}
	}

	return alerts
}

// calculateTotalRisk 计算总风险
func (rm *RiskManager) calculateTotalRisk() float64 {
	totalRisk := 0.0

	for _, position := range rm.portfolio.Positions {
		positionValue := float64(position.Quantity) * position.CurrentPrice
		stopLossRisk := (position.CurrentPrice - position.StopLoss) / position.CurrentPrice
		totalRisk += positionValue * stopLossRisk
	}

	return totalRisk / rm.portfolio.TotalValue
}

// UpdatePortfolio 更新投资组合
func (rm *RiskManager) UpdatePortfolio(symbol string, currentPrice float64) {
	if position, exists := rm.portfolio.Positions[symbol]; exists {
		position.CurrentPrice = currentPrice

		// 更新总价值
		rm.updateTotalValue()

		// 检查止损止盈
		rm.checkStopLossTakeProfit(symbol)
	}
}

// updateTotalValue 更新总价值
func (rm *RiskManager) updateTotalValue() {
	totalValue := rm.portfolio.Cash

	for _, position := range rm.portfolio.Positions {
		positionValue := float64(position.Quantity) * position.CurrentPrice
		totalValue += positionValue
	}

	rm.portfolio.TotalValue = totalValue

	// 更新峰值和回撤
	if totalValue > rm.portfolio.PeakValue {
		rm.portfolio.PeakValue = totalValue
	}

	rm.portfolio.CurrentDrawdown = (rm.portfolio.PeakValue - totalValue) / rm.portfolio.PeakValue

	if rm.portfolio.CurrentDrawdown > rm.portfolio.MaxDrawdown {
		rm.portfolio.MaxDrawdown = rm.portfolio.CurrentDrawdown
	}
}

// checkStopLossTakeProfit 检查止损止盈
func (rm *RiskManager) checkStopLossTakeProfit(symbol string) {
	position := rm.portfolio.Positions[symbol]

	// 检查止损
	if position.CurrentPrice <= position.StopLoss {
		rm.addAlert(RiskAlert{
			Type:      "stop_loss",
			Message:   "触发止损: " + symbol,
			Level:     "high",
			Timestamp: time.Now(),
			Value:     position.CurrentPrice,
			Threshold: position.StopLoss,
		})
	}

	// 检查止盈
	if position.CurrentPrice >= position.TakeProfit {
		rm.addAlert(RiskAlert{
			Type:      "take_profit",
			Message:   "触发止盈: " + symbol,
			Level:     "medium",
			Timestamp: time.Now(),
			Value:     position.CurrentPrice,
			Threshold: position.TakeProfit,
		})
	}
}

// addAlert 添加预警
func (rm *RiskManager) addAlert(alert RiskAlert) {
	rm.alerts = append(rm.alerts, alert)
}

// GetAlerts 获取预警
func (rm *RiskManager) GetAlerts() []RiskAlert {
	return rm.alerts
}

// ClearAlerts 清除预警
func (rm *RiskManager) ClearAlerts() {
	rm.alerts = make([]RiskAlert, 0)
}

// GetPortfolio 获取投资组合
func (rm *RiskManager) GetPortfolio() *Portfolio {
	return rm.portfolio
}

// GetMetrics 获取风险指标
func (rm *RiskManager) GetMetrics() *RiskMetrics {
	return rm.metrics
} 