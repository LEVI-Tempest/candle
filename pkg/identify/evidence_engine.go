package identify

import "math"

// BuildPatternEvidence combines pattern signals and volume/context features.
// BuildPatternEvidence 将形态信号与量价/上下文特征融合为证据输出。
func BuildPatternEvidence(
	patterns []PatternSignal,
	candles []CandlestickWrapper,
	cfg EvidenceConfig,
) []PatternEvidence {
	if len(patterns) == 0 || len(candles) == 0 {
		return nil
	}

	ind := ComputeVolumeIndicators(candles, cfg.MFIPeriod, cfg.CMFPeriod)
	out := make([]PatternEvidence, 0, len(patterns))

	for _, p := range patterns {
		if p.Position < 0 || p.Position >= len(candles) {
			continue
		}
		ev := PatternEvidence{
			PatternType:  p.Type,
			Direction:    p.Direction,
			Position:     p.Position,
			Time:         p.Time,
			Price:        p.Price,
			BaseStrength: clamp01(p.Strength),
		}

		contextScore, ctxFactors := scoreContext(candles, p, cfg)
		volumeScore, volFactors, contradictions := scoreVolume(candles, ind, p, cfg)

		finalScore := ev.BaseStrength*cfg.BaseWeight +
			contextScore*cfg.ContextWeight +
			volumeScore*cfg.VolumeWeight
		finalScore -= contradictionPenalty(contradictions)

		ev.ContextScore = clamp01(contextScore)
		ev.VolumeScore = clamp01(volumeScore)
		ev.FinalScore = clamp01(finalScore)
		ev.ConfidenceLevel = toConfidence(ev.FinalScore)
		ev.ContextFactors = ctxFactors
		ev.VolumeFactors = volFactors
		ev.ContradictionFactors = contradictions
		out = append(out, ev)
	}

	return out
}

func scoreContext(cs []CandlestickWrapper, p PatternSignal, cfg EvidenceConfig) (float64, []FactorHit) {
	factors := make([]FactorHit, 0, 2)
	if p.Position < 3 {
		factors = append(factors, FactorHit{
			Name:      "trend_window",
			Value:     float64(p.Position),
			Threshold: 3,
			Passed:    false,
			Reason:    "insufficient candles for trend context",
		})
		return 0.4, factors
	}

	window := cfg.ContextWindow
	if window < 3 {
		window = 9
	}
	start := p.Position - window
	if start < 0 {
		start = 0
	}
	trend, metrics := AnalyzeLongTermTrend(cs[start:p.Position+1], p.Position-start+1, cfg.ContextTrendThreshold)
	score := 0.5

	if metrics != nil {
		factors = append(factors, FactorHit{
			Name:      "trend_strength",
			Value:     metrics["trendStrength"],
			Threshold: 20,
			Passed:    metrics["trendStrength"] >= 20,
			Reason:    "trend strength from AnalyzeLongTermTrend",
		})
	}

	align := false
	switch p.Direction {
	case "bullish":
		align = trend == TrendYin || trend == TrendMiddle
	case "bearish":
		align = trend == TrendYang || trend == TrendMiddle
	default:
		align = true
	}
	factors = append(factors, FactorHit{
		Name:      "trend_alignment",
		Value:     trendToValue(trend),
		Threshold: 0,
		Passed:    align,
		Reason:    "reversal patterns prefer opposite or neutral prior trend",
	})

	if align {
		score += 0.3
	} else {
		score -= 0.2
	}
	return clamp01(score), factors
}

func scoreVolume(
	cs []CandlestickWrapper,
	ind VolumeIndicatorSeries,
	p PatternSignal,
	cfg EvidenceConfig,
) (float64, []FactorHit, []string) {
	i := p.Position
	c := cs[i]
	factors := make([]FactorHit, 0, 5)
	contradictions := make([]string, 0, 2)
	score := 0.5

	avgVol := averageVolumeBefore(cs, i, cfg.VolumeLookback)
	volRatio := 0.0
	if avgVol > 0 {
		volRatio = c.Volume / avgVol
	}
	beiliangThreshold := cfg.BeiliangThreshold
	if beiliangThreshold <= 0 {
		beiliangThreshold = cfg.VolumeBoostThreshold
	}
	volBoost := volRatio >= beiliangThreshold
	factors = append(factors, FactorHit{
		Name:      "beiliang_confirm",
		Value:     volRatio,
		Threshold: beiliangThreshold,
		Passed:    volBoost,
		Reason:    "volume / ma_volume_n reaches beiliang threshold",
	})
	if volBoost {
		score += 0.25
	}

	factors = append(factors, FactorHit{
		Name:      "volume_ratio",
		Value:     volRatio,
		Threshold: 1,
		Passed:    volRatio >= 1,
		Reason:    "raw volume ratio for audit and tuning",
	})

	lowVol := avgVol > 0 && volRatio <= cfg.VolumeShrinkThreshold
	factors = append(factors, FactorHit{
		Name:      "volume_shrink",
		Value:     volRatio,
		Threshold: cfg.VolumeShrinkThreshold,
		Passed:    lowVol,
		Reason:    "low-volume warning for weak confirmation",
	})
	if lowVol {
		score -= 0.2
	}

	mfi := getSeriesValue(ind.MFI, i)
	mfiPass := (p.Direction == "bullish" && mfi < 40) ||
		(p.Direction == "bearish" && mfi > 60) ||
		p.Direction == "neutral"
	factors = append(factors, FactorHit{
		Name:      "mfi_regime",
		Value:     mfi,
		Threshold: 50,
		Passed:    mfiPass,
		Reason:    "MFI supports overbought/oversold regime checks",
	})
	if mfiPass {
		score += 0.15
	}

	cmf := getSeriesValue(ind.CMF, i)
	cmfPass := (p.Direction == "bullish" && cmf >= 0) ||
		(p.Direction == "bearish" && cmf <= 0) ||
		p.Direction == "neutral"
	factors = append(factors, FactorHit{
		Name:      "cmf_direction",
		Value:     cmf,
		Threshold: 0,
		Passed:    cmfPass,
		Reason:    "CMF confirms net inflow/outflow direction",
	})
	if cmfPass {
		score += 0.10
	}

	obvContradiction := detectOBVDivergence(cs, ind.OBV, i, p.Direction, cfg.OBVDivergenceLookback)
	if obvContradiction != "" {
		contradictions = append(contradictions, obvContradiction)
	}

	return clamp01(score), factors, contradictions
}

func detectOBVDivergence(cs []CandlestickWrapper, obv []float64, i int, direction string, lookback int) string {
	if lookback < 1 {
		lookback = 5
	}
	start := i - lookback
	if start < 0 || len(obv) <= i {
		return ""
	}

	priceDelta := cs[i].Close - cs[start].Close
	obvDelta := obv[i] - obv[start]

	if direction == "bullish" && priceDelta > 0 && obvDelta < 0 {
		return "price rises while OBV falls (bearish divergence)"
	}
	if direction == "bearish" && priceDelta < 0 && obvDelta > 0 {
		return "price falls while OBV rises (bullish divergence)"
	}
	return ""
}

func contradictionPenalty(contradictions []string) float64 {
	if len(contradictions) == 0 {
		return 0
	}
	return math.Min(0.2, float64(len(contradictions))*0.1)
}

func trendToValue(t Trend) float64 {
	switch t {
	case TrendYang:
		return 1
	case TrendYin:
		return -1
	case TrendMiddle:
		return 0
	default:
		return 0
	}
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func toConfidence(score float64) string {
	switch {
	case score >= 0.8:
		return "strong"
	case score >= 0.6:
		return "medium"
	default:
		return "weak"
	}
}

func getSeriesValue(series []float64, i int) float64 {
	if i < 0 || i >= len(series) {
		return 0
	}
	v := series[i]
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return v
}
