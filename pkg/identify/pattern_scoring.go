package identify

import (
	"math"
	"time"
)

/***
 * @author Tempest
 * @description Pattern Scoring System (形态评分系统)
 * @date 2025/01/27
 * @version 1.0.0
 */

// PatternResult 形态识别结果
type PatternResult struct {
	Pattern     string    // 形态名称
	Detected    bool      // 是否检测到
	Strength    float64   // 形态强度 (0-1)
	Confidence  float64   // 置信度 (0-1)
	Position    int       // 位置索引
	Direction   string    // 方向 (bullish/bearish/neutral)
	Description string    // 形态描述
	Timestamp   time.Time // 检测时间
}

// PatternConfig 形态识别配置
type PatternConfig struct {
	// 基础阈值配置
	BodyRatioThreshold     float64 // 实体比例阈值
	ShadowRatioThreshold   float64 // 影线比例阈值
	VolumeRatioThreshold   float64 // 成交量比例阈值
	
	// 趋势背景要求
	TrendContextRequired   bool    // 是否需要趋势背景
	TrendStrengthThreshold float64 // 趋势强度阈值
	
	// 评分权重
	ShapeWeight           float64 // 形状权重
	SizeWeight            float64 // 大小权重
	VolumeWeight          float64 // 成交量权重
	TrendWeight           float64 // 趋势权重
	ContextWeight         float64 // 上下文权重
}

// DefaultPatternConfig 默认配置
func DefaultPatternConfig() PatternConfig {
	return PatternConfig{
		BodyRatioThreshold:     0.1,
		ShadowRatioThreshold:   0.2,
		VolumeRatioThreshold:   1.5,
		TrendContextRequired:   true,
		TrendStrengthThreshold: 0.6,
		ShapeWeight:            0.4,
		SizeWeight:             0.3,
		VolumeWeight:           0.1,
		TrendWeight:            0.1,
		ContextWeight:          0.1,
	}
}

// PatternScorer 形态评分器
type PatternScorer struct {
	config PatternConfig
}

// NewPatternScorer 创建新的形态评分器
func NewPatternScorer(config PatternConfig) *PatternScorer {
	return &PatternScorer{
		config: config,
	}
}

// ScoreHammer 锤头线评分
func (ps *PatternScorer) ScoreHammer(cs []CandlestickWrapper) PatternResult {
	if len(cs) < 1 {
		return PatternResult{Pattern: "Hammer", Detected: false}
	}
	
	c := cs[0]
	body := c.Body()
	lowerShadow := c.LowerShadow()
	upperShadow := c.UpperShadow()
	totalRange := c.High - c.Low
	
	if body == 0 || totalRange == 0 {
		return PatternResult{Pattern: "Hammer", Detected: false}
	}
	
	// 基础条件检查
	bodyRatio := body / totalRange
	lowerShadowRatio := lowerShadow / totalRange
	upperShadowRatio := upperShadow / totalRange
	
	// 形状评分 (40%)
	shapeScore := 0.0
	if bodyRatio < ps.config.BodyRatioThreshold {
		shapeScore += 0.4
	}
	if lowerShadowRatio > ps.config.ShadowRatioThreshold {
		shapeScore += 0.4
	}
	if upperShadowRatio < 0.05 {
		shapeScore += 0.2
	}
	
	// 大小评分 (30%)
	sizeScore := 0.0
	if lowerShadow > 2*body {
		sizeScore += 0.5
	}
	if lowerShadow > 3*body {
		sizeScore += 0.3
	}
	if body > 0.01*totalRange {
		sizeScore += 0.2
	}
	
	// 成交量评分 (10%)
	volumeScore := 0.0
	if len(cs) > 1 {
		avgVolume := calculateAverageVolume(cs[1:min(len(cs), 6)])
		if c.Volume > avgVolume*ps.config.VolumeRatioThreshold {
			volumeScore = 1.0
		}
	}
	
	// 趋势背景评分 (10%)
	trendScore := 0.0
	if ps.config.TrendContextRequired && len(cs) > 5 {
		trend := analyzeTrendContext(cs[1:6])
		if trend == "downtrend" {
			trendScore = 1.0
		}
	}
	
	// 上下文评分 (10%)
	contextScore := 0.0
	if c.Close > c.Low && c.Close > c.Open {
		contextScore = 1.0
	}
	
	// 计算总分
	totalScore := shapeScore*ps.config.ShapeWeight +
		sizeScore*ps.config.SizeWeight +
		volumeScore*ps.config.VolumeWeight +
		trendScore*ps.config.TrendWeight +
		contextScore*ps.config.ContextWeight
	
	detected := totalScore > 0.7
	confidence := math.Min(totalScore, 1.0)
	
	return PatternResult{
		Pattern:     "Hammer",
		Detected:    detected,
		Strength:    totalScore,
		Confidence:  confidence,
		Position:    0,
		Direction:   "bullish",
		Description: "锤头线 - 看涨反转信号",
		Timestamp:   time.Now(),
	}
}

// ScoreDoji 十字星评分
func (ps *PatternScorer) ScoreDoji(cs []CandlestickWrapper) PatternResult {
	if len(cs) < 1 {
		return PatternResult{Pattern: "Doji", Detected: false}
	}
	
	c := cs[0]
	body := c.Body()
	upperShadow := c.UpperShadow()
	lowerShadow := c.LowerShadow()
	totalRange := c.High - c.Low
	
	if totalRange == 0 {
		return PatternResult{Pattern: "Doji", Detected: false}
	}
	
	// 形状评分 (50%)
	shapeScore := 0.0
	bodyRatio := body / totalRange
	if bodyRatio < ps.config.BodyRatioThreshold {
		shapeScore += 0.5
	}
	
	// 影线评分 (30%)
	shadowScore := 0.0
	upperShadowRatio := upperShadow / totalRange
	lowerShadowRatio := lowerShadow / totalRange
	if upperShadowRatio > ps.config.ShadowRatioThreshold && lowerShadowRatio > ps.config.ShadowRatioThreshold {
		shadowScore = 1.0
	}
	
	// 平衡性评分 (20%)
	balanceScore := 0.0
	shadowBalance := math.Abs(upperShadow - lowerShadow) / totalRange
	if shadowBalance < 0.1 {
		balanceScore = 1.0
	}
	
	totalScore := shapeScore*0.5 + shadowScore*0.3 + balanceScore*0.2
	detected := totalScore > 0.7
	confidence := math.Min(totalScore, 1.0)
	
	return PatternResult{
		Pattern:     "Doji",
		Detected:    detected,
		Strength:    totalScore,
		Confidence:  confidence,
		Position:    0,
		Direction:   "neutral",
		Description: "十字星 - 市场犹豫信号",
		Timestamp:   time.Now(),
	}
}

// ScoreMarubozu 光头光脚评分
func (ps *PatternScorer) ScoreMarubozu(cs []CandlestickWrapper) PatternResult {
	if len(cs) < 1 {
		return PatternResult{Pattern: "Marubozu", Detected: false}
	}
	
	c := cs[0]
	body := c.Body()
	upperShadow := c.UpperShadow()
	lowerShadow := c.LowerShadow()
	totalRange := c.High - c.Low
	
	if totalRange == 0 {
		return PatternResult{Pattern: "Marubozu", Detected: false}
	}
	
	// 实体大小评分 (40%)
	sizeScore := 0.0
	bodyRatio := body / totalRange
	if bodyRatio > 0.8 {
		sizeScore = 1.0
	} else if bodyRatio > 0.6 {
		sizeScore = 0.7
	}
	
	// 影线缺失评分 (40%)
	shadowScore := 0.0
	tolerance := 0.01 * body
	if upperShadow <= tolerance && lowerShadow <= tolerance {
		shadowScore = 1.0
	} else if upperShadow <= tolerance*2 && lowerShadow <= tolerance*2 {
		shadowScore = 0.7
	}
	
	// 方向评分 (20%)
	directionScore := 0.0
	if c.IsBullish() {
		directionScore = 1.0
	} else if c.IsBearish() {
		directionScore = 1.0
	}
	
	totalScore := sizeScore*0.4 + shadowScore*0.4 + directionScore*0.2
	detected := totalScore > 0.7
	confidence := math.Min(totalScore, 1.0)
	
	direction := "neutral"
	if c.IsBullish() {
		direction = "bullish"
	} else if c.IsBearish() {
		direction = "bearish"
	}
	
	return PatternResult{
		Pattern:     "Marubozu",
		Detected:    detected,
		Strength:    totalScore,
		Confidence:  confidence,
		Position:    0,
		Direction:   direction,
		Description: "光头光脚 - 强烈趋势信号",
		Timestamp:   time.Now(),
	}
}

// ScoreBullishEngulfing 看涨吞噬评分
func (ps *PatternScorer) ScoreBullishEngulfing(cs []CandlestickWrapper) PatternResult {
	if len(cs) < 2 {
		return PatternResult{Pattern: "BullishEngulfing", Detected: false}
	}
	
	first := cs[1]
	second := cs[0]
	
	// 基础条件检查
	if !first.IsBearish() || !second.IsBullish() {
		return PatternResult{Pattern: "BullishEngulfing", Detected: false}
	}
	
	// 吞噬程度评分 (50%)
	engulfingScore := 0.0
	if second.Open <= first.Close && second.Close >= first.Open {
		engulfingRatio := (second.Close - second.Open) / (first.Open - first.Close)
		if engulfingRatio > 1.5 {
			engulfingScore = 1.0
		} else if engulfingRatio > 1.2 {
			engulfingScore = 0.8
		} else {
			engulfingScore = 0.6
		}
	}
	
	// 实体大小评分 (30%)
	sizeScore := 0.0
	firstBodyRatio := first.Body() / (first.High - first.Low)
	secondBodyRatio := second.Body() / (second.High - second.Low)
	if firstBodyRatio > 0.6 && secondBodyRatio > 0.6 {
		sizeScore = 1.0
	} else if firstBodyRatio > 0.4 && secondBodyRatio > 0.4 {
		sizeScore = 0.7
	}
	
	// 成交量评分 (20%)
	volumeScore := 0.0
	if second.Volume > first.Volume*1.2 {
		volumeScore = 1.0
	} else if second.Volume > first.Volume {
		volumeScore = 0.7
	}
	
	totalScore := engulfingScore*0.5 + sizeScore*0.3 + volumeScore*0.2
	detected := totalScore > 0.6
	confidence := math.Min(totalScore, 1.0)
	
	return PatternResult{
		Pattern:     "BullishEngulfing",
		Detected:    detected,
		Strength:    totalScore,
		Confidence:  confidence,
		Position:    0,
		Direction:   "bullish",
		Description: "看涨吞噬 - 强烈反转信号",
		Timestamp:   time.Now(),
	}
}

// ScoreMorningStar 启明星评分
func (ps *PatternScorer) ScoreMorningStar(cs []CandlestickWrapper) PatternResult {
	if len(cs) < 3 {
		return PatternResult{Pattern: "MorningStar", Detected: false}
	}
	
	first := cs[2]  // 长阴线
	second := cs[1] // 小实体
	third := cs[0]  // 长阳线
	
	// 第一根蜡烛评分 (30%)
	firstScore := 0.0
	if first.IsBearish() {
		firstBodyRatio := first.Body() / (first.High - first.Low)
		if firstBodyRatio > 0.6 {
			firstScore = 1.0
		} else if firstBodyRatio > 0.4 {
			firstScore = 0.7
		}
	}
	
	// 第二根蜡烛评分 (20%)
	secondScore := 0.0
	secondBodyRatio := second.Body() / (second.High - second.Low)
	if secondBodyRatio < 0.3 {
		secondScore = 1.0
	} else if secondBodyRatio < 0.5 {
		secondScore = 0.7
	}
	
	// 第三根蜡烛评分 (30%)
	thirdScore := 0.0
	if third.IsBullish() {
		thirdBodyRatio := third.Body() / (third.High - third.Low)
		if thirdBodyRatio > 0.6 {
			thirdScore = 1.0
		} else if thirdBodyRatio > 0.4 {
			thirdScore = 0.7
		}
	}
	
	// 位置关系评分 (20%)
	positionScore := 0.0
	if second.High < first.Close && third.Close > (first.Open+first.Close)/2 {
		positionScore = 1.0
	} else if second.High < first.Close {
		positionScore = 0.7
	}
	
	totalScore := firstScore*0.3 + secondScore*0.2 + thirdScore*0.3 + positionScore*0.2
	detected := totalScore > 0.6
	confidence := math.Min(totalScore, 1.0)
	
	return PatternResult{
		Pattern:     "MorningStar",
		Detected:    detected,
		Strength:    totalScore,
		Confidence:  confidence,
		Position:    0,
		Direction:   "bullish",
		Description: "启明星 - 强烈看涨反转信号",
		Timestamp:   time.Now(),
	}
}

// 辅助函数

// calculateAverageVolume 计算平均成交量
func calculateAverageVolume(candles []CandlestickWrapper) float64 {
	if len(candles) == 0 {
		return 0
	}
	
	total := 0.0
	for _, c := range candles {
		total += c.Volume
	}
	return total / float64(len(candles))
}

// analyzeTrendContext 分析趋势背景
func analyzeTrendContext(candles []CandlestickWrapper) string {
	if len(candles) < 3 {
		return "unknown"
	}
	
	risingCount := 0
	fallingCount := 0
	
	for i := 1; i < len(candles); i++ {
		if candles[i].Close > candles[i-1].Close {
			risingCount++
		} else if candles[i].Close < candles[i-1].Close {
			fallingCount++
		}
	}
	
	if risingCount > fallingCount*2 {
		return "uptrend"
	} else if fallingCount > risingCount*2 {
		return "downtrend"
	}
	return "sideways"
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AnalyzeAllPatterns 分析所有形态
func (ps *PatternScorer) AnalyzeAllPatterns(cs []CandlestickWrapper) []PatternResult {
	var results []PatternResult
	
	// 单根蜡烛形态
	if len(cs) >= 1 {
		results = append(results, ps.ScoreHammer(cs))
		results = append(results, ps.ScoreDoji(cs))
		results = append(results, ps.ScoreMarubozu(cs))
	}
	
	// 双根蜡烛形态
	if len(cs) >= 2 {
		results = append(results, ps.ScoreBullishEngulfing(cs))
	}
	
	// 三根蜡烛形态
	if len(cs) >= 3 {
		results = append(results, ps.ScoreMorningStar(cs))
	}
	
	return results
}

// GetStrongestPatterns 获取最强的形态
func (ps *PatternScorer) GetStrongestPatterns(cs []CandlestickWrapper, minStrength float64) []PatternResult {
	allPatterns := ps.AnalyzeAllPatterns(cs)
	var strongPatterns []PatternResult
	
	for _, pattern := range allPatterns {
		if pattern.Detected && pattern.Strength >= minStrength {
			strongPatterns = append(strongPatterns, pattern)
		}
	}
	
	return strongPatterns
} 