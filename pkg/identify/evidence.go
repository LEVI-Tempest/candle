package identify

// PatternSignal is the normalized pattern input used by the evidence engine.
// PatternSignal 是证据引擎使用的标准化形态输入。
type PatternSignal struct {
	Type      string  // Pattern type (形态类型)
	Direction string  // bullish/bearish/neutral
	Position  int     // Index in candlestick series (K线序号)
	Strength  float64 // Base pattern score from recognizer (形态基础分)
	Risk      float64 // Risk score from recognizer (风险分)
	Price     float64 // Trigger price (触发价)
	Time      string  // Trigger time string (触发时间)
}

// FactorHit records one explainable factor and whether it is satisfied.
// FactorHit 记录一个可解释因子及是否命中。
type FactorHit struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Threshold float64 `json:"threshold"`
	Passed    bool    `json:"passed"`
	Reason    string  `json:"reason"`
}

// PatternEvidence is the structured output for decision/audit layers.
// PatternEvidence 是用于决策与审计层的结构化输出。
type PatternEvidence struct {
	PatternType          string      `json:"pattern_type"`
	Direction            string      `json:"direction"`
	Position             int         `json:"position"`
	Time                 string      `json:"time"`
	Price                float64     `json:"price"`
	BaseStrength         float64     `json:"base_strength"`
	ContextScore         float64     `json:"context_score"`
	VolumeScore          float64     `json:"volume_score"`
	FinalScore           float64     `json:"final_score"`
	ConfidenceLevel      string      `json:"confidence_level"`
	ContextFactors       []FactorHit `json:"context_factors"`
	VolumeFactors        []FactorHit `json:"volume_factors"`
	ContradictionFactors []string    `json:"contradiction_factors"`
}

// EvidenceConfig controls scoring weights and thresholds.
// EvidenceConfig 控制评分权重和阈值。
type EvidenceConfig struct {
	VolumeLookback        int     `json:"volume_lookback"`
	MFIPeriod             int     `json:"mfi_period"`
	CMFPeriod             int     `json:"cmf_period"`
	BeiliangThreshold     float64 `json:"beiliang_threshold"`
	VolumeBoostThreshold  float64 `json:"volume_boost_threshold"`
	VolumeShrinkThreshold float64 `json:"volume_shrink_threshold"`
	BaseWeight            float64 `json:"base_weight"`
	ContextWeight         float64 `json:"context_weight"`
	VolumeWeight          float64 `json:"volume_weight"`
	// ContextWindow is the candlestick lookback for trend context scoring (default 9).
	// ContextWindow 是趋势上下文评分的回望窗口（默认 9）。
	ContextWindow int `json:"context_window"`
	// ContextTrendThreshold is the minimum angle/slope for a valid trend (default 3.0).
	// ContextTrendThreshold 是判断有效趋势的最小斜率阈值（默认 3.0）。
	ContextTrendThreshold float64 `json:"context_trend_threshold"`
	// OBVDivergenceLookback is the candlestick lookback for OBV divergence check (default 5).
	// OBVDivergenceLookback 是 OBV 背离检测的回望窗口（默认 5）。
	OBVDivergenceLookback int `json:"obv_divergence_lookback"`
}

// DefaultEvidenceConfig returns a conservative default config.
// DefaultEvidenceConfig 返回较保守的默认配置。
func DefaultEvidenceConfig() EvidenceConfig {
	return EvidenceConfig{
		VolumeLookback:        10,
		MFIPeriod:             14,
		CMFPeriod:             20,
		BeiliangThreshold:     1.5,
		VolumeBoostThreshold:  1.5,
		VolumeShrinkThreshold: 0.8,
		BaseWeight:            0.55,
		ContextWeight:         0.20,
		VolumeWeight:          0.25,
		ContextWindow:         9,
		ContextTrendThreshold: 3.0,
		OBVDivergenceLookback: 5,
	}
}
