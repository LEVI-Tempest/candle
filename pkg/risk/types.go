package risk

import (
	"time"
)

// RiskProfile 风险偏好配置
type RiskProfile struct {
	Conservative float64 // 保守型 (0.1-0.3)
	Moderate     float64 // 中性型 (0.3-0.6)
	Aggressive   float64 // 激进型 (0.6-0.9)
}

// Position 持仓信息
type Position struct {
	Symbol       string    // 股票代码
	Quantity     int       // 数量
	EntryPrice   float64   // 入场价格
	CurrentPrice float64   // 当前价格
	EntryTime    time.Time // 入场时间
	StopLoss     float64   // 止损价格
	TakeProfit   float64   // 止盈价格
}

// Portfolio 投资组合
type Portfolio struct {
	TotalValue       float64              // 总价值
	Cash             float64              // 现金
	Positions        map[string]*Position // 持仓
	RiskBudget       float64              // 风险预算
	MaxDrawdown      float64              // 最大回撤
	CurrentDrawdown  float64              // 当前回撤
	PeakValue        float64              // 历史峰值
}

// RiskMetrics 风险指标
type RiskMetrics struct {
	Volatility   float64 // 波动率
	VaR          float64 // 风险价值
	CVaR         float64 // 条件风险价值
	MaxDrawdown  float64 // 最大回撤
	SharpeRatio  float64 // 夏普比率
	SortinoRatio float64 // 索提诺比率
	CalmarRatio  float64 // 卡玛比率
	Correlation  float64 // 相关性
}

// RiskConfig 风险配置
type RiskConfig struct {
	MaxPositionSize   float64 // 单笔最大仓位比例
	MaxTotalRisk      float64 // 总风险上限
	StopLossRatio     float64 // 止损比例
	TakeProfitRatio   float64 // 止盈比例
	TrailingStopRatio float64 // 移动止损比例
	MaxDrawdownLimit  float64 // 最大回撤限制
	CorrelationLimit  float64 // 相关性限制
	VolatilityPeriod  int     // 波动率计算周期
	VaRConfidence     float64 // VaR置信度
	VaRTimeWindow     int     // VaR时间窗口(天)
}

// RiskAlert 风险预警
type RiskAlert struct {
	Type      string    // 预警类型
	Message   string    // 预警信息
	Level     string    // 预警级别 (low/medium/high/critical)
	Timestamp time.Time // 预警时间
	Value     float64   // 触发值
	Threshold float64   // 阈值
}

// DefaultRiskConfig 默认风险配置
func DefaultRiskConfig() RiskConfig {
	return RiskConfig{
		MaxPositionSize:   0.05, // 单笔最大5%
		MaxTotalRisk:      0.10, // 总风险最大10%
		StopLossRatio:     0.02, // 止损2%
		TakeProfitRatio:   0.06, // 止盈6%
		TrailingStopRatio: 0.03, // 移动止损3%
		MaxDrawdownLimit:  0.15, // 最大回撤15%
		CorrelationLimit:  0.7,  // 相关性限制0.7
		VolatilityPeriod:  20,   // 20天波动率
		VaRConfidence:     0.95, // 95%置信度
		VaRTimeWindow:     1,    // 1天VaR
	}
} 