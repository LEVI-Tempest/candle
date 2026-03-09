package signal

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/LEVI-Tempest/Candle/pkg/charting"
	"github.com/LEVI-Tempest/Candle/pkg/identify"
	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

const (
	VolumeConfirm    = "confirm"
	VolumeNeutral    = "neutral"
	VolumeContradict = "contradict"
)

// BuildReport generates structured signal output with trend filter and decision score.
// BuildReport 生成包含趋势过滤和决策分的结构化信号输出。
func BuildReport(symbol, asOf, source string, candles []*v1.Candlestick, cfg Config) Report {
	ek := charting.NewEnhancedKline()
	ek.LoadData(candles)
	ek.AutoDetectPatterns()

	signals := toPatternSignals(ek.Patterns)
	evidence := identify.BuildPatternEvidence(signals, ek.Data, cfg.Evidence)
	sort.Slice(evidence, func(i, j int) bool {
		if evidence[i].FinalScore == evidence[j].FinalScore {
			return evidence[i].Position < evidence[j].Position
		}
		return evidence[i].FinalScore > evidence[j].FinalScore
	})

	trend := determineTrendByMA20(ek.Data, cfg.Trend.Period)
	patternReports := make([]PatternReport, 0, len(ek.Patterns))
	evidenceByKey := make(map[string]identify.PatternEvidence, len(evidence))
	for _, ev := range evidence {
		evidenceByKey[evidenceKey(ev.PatternType, ev.Position)] = ev
	}

	for _, p := range ek.Patterns {
		ev, ok := evidenceByKey[evidenceKey(p.Type, p.Position)]
		if !ok {
			continue
		}
		volumeState, reason := volumeStateAndReason(ev)
		score := decisionScore(cfg.Score, ev.BaseStrength, trendMatchScore(p.Type, trend), volumeStateScore(volumeState))
		level := decisionLevel(score, cfg.Score.StrongThreshold, cfg.Score.MediumThreshold)
		r3, r5, r10 := forwardReturns(candles, p.Position)

		patternReports = append(patternReports, PatternReport{
			Type:          p.Type,
			Direction:     patternDirection(p.Type),
			Position:      p.Position,
			Strength:      p.Strength,
			Risk:          p.Risk,
			Price:         p.Price,
			Time:          p.Time,
			VolumeState:   volumeState,
			DecisionScore: score,
			DecisionLevel: level,
			Reason:        reason,
			ForwardRet3:   r3,
			ForwardRet5:   r5,
			ForwardRet10:  r10,
		})
	}

	sort.Slice(patternReports, func(i, j int) bool {
		if patternReports[i].DecisionScore == patternReports[j].DecisionScore {
			return patternReports[i].Position < patternReports[j].Position
		}
		return patternReports[i].DecisionScore > patternReports[j].DecisionScore
	})

	return Report{
		Symbol:          symbol,
		AsOf:            asOf,
		Source:          source,
		Trend:           trend,
		Score:           normalizedScore(patternReports),
		DecisionScore:   topDecisionScore(patternReports),
		DecisionLevel:   decisionLevel(topDecisionScore(patternReports), cfg.Score.StrongThreshold, cfg.Score.MediumThreshold),
		Patterns:        patternReports,
		Evidence:        evidence,
		CounterEvidence: collectCounterEvidence(evidence),
		InvalidIf: []string{
			"data source has missing/incorrect OHLCV records",
			"next trading sessions show no volume confirmation",
			"price breaks pattern invalidation level with high volatility",
		},
	}
}

// AppendSignalLogCSV appends signal rows to a CSV file for personal replay.
// AppendSignalLogCSV 将信号记录追加到 CSV 文件，便于个人复盘。
func AppendSignalLogCSV(path string, report Report) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	needHeader := false
	if _, err := os.Stat(path); err != nil {
		needHeader = true
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if needHeader {
		if err := w.Write([]string{
			"time", "symbol", "pattern", "trend", "volume_state", "decision_score", "decision_level",
			"forward_ret_3", "forward_ret_5", "forward_ret_10", "reason",
		}); err != nil {
			return err
		}
	}

	for _, p := range report.Patterns {
		row := []string{
			report.AsOf,
			report.Symbol,
			p.Type,
			report.Trend,
			p.VolumeState,
			fmt.Sprintf("%.2f", p.DecisionScore),
			p.DecisionLevel,
			formatFloatPtr(p.ForwardRet3),
			formatFloatPtr(p.ForwardRet5),
			formatFloatPtr(p.ForwardRet10),
			stringsJoin(p.Reason, " | "),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return w.Error()
}

func toPatternSignals(patterns []charting.Pattern) []identify.PatternSignal {
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

func determineTrendByMA20(cs []identify.CandlestickWrapper, period int) string {
	if period < 2 {
		period = 20
	}
	if len(cs) < period+1 {
		return "unknown"
	}
	lastMA := movingAverageClose(cs, len(cs)-period, len(cs))
	prevMA := movingAverageClose(cs, len(cs)-period-1, len(cs)-1)
	lastClose := cs[len(cs)-1].Close

	switch {
	case lastClose > lastMA && lastMA > prevMA:
		return "up"
	case lastClose < lastMA && lastMA < prevMA:
		return "down"
	default:
		return "sideway"
	}
}

func movingAverageClose(cs []identify.CandlestickWrapper, start, end int) float64 {
	if start < 0 {
		start = 0
	}
	if end > len(cs) {
		end = len(cs)
	}
	if start >= end {
		return 0
	}
	sum := 0.0
	for i := start; i < end; i++ {
		sum += cs[i].Close
	}
	return sum / float64(end-start)
}

func trendMatchScore(patternType, trend string) float64 {
	dir := patternDirection(patternType)
	switch dir {
	case "neutral":
		return 0.5
	case "bullish":
		if trend == "down" {
			return 1.0
		}
		if trend == "sideway" {
			return 0.5
		}
		return 0.25
	case "bearish":
		if trend == "up" {
			return 1.0
		}
		if trend == "sideway" {
			return 0.5
		}
		return 0.25
	default:
		return 0.0
	}
}

func volumeStateAndReason(ev identify.PatternEvidence) (string, []string) {
	if len(ev.ContradictionFactors) > 0 {
		return VolumeContradict, []string{ev.ContradictionFactors[0]}
	}
	passed := 0
	reasons := make([]string, 0, 2)
	for _, f := range ev.VolumeFactors {
		if f.Passed {
			passed++
			if len(reasons) < 2 {
				reasons = append(reasons, f.Name)
			}
		}
	}
	if len(reasons) == 0 {
		for _, f := range ev.ContextFactors {
			if f.Passed && len(reasons) < 2 {
				reasons = append(reasons, f.Name)
			}
		}
	}
	switch {
	case passed >= 2:
		return VolumeConfirm, reasons
	case passed == 1:
		return VolumeNeutral, reasons
	default:
		return VolumeContradict, reasons
	}
}

func volumeStateScore(state string) float64 {
	switch state {
	case VolumeConfirm:
		return 1.0
	case VolumeNeutral:
		return 0.5
	default:
		return 0.0
	}
}

func decisionScore(cfg ScoreConfig, baseStrength, trendScore, volumeScore float64) float64 {
	score := baseStrength*cfg.PatternWeight + trendScore*cfg.TrendWeight + volumeScore*cfg.VolumeWeight
	if math.IsNaN(score) || math.IsInf(score, 0) {
		return 0
	}
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func decisionLevel(score, strong, medium float64) string {
	if score >= strong {
		return "strong"
	}
	if score >= medium {
		return "medium"
	}
	return "weak"
}

func normalizedScore(patterns []PatternReport) float64 {
	if len(patterns) == 0 {
		return 0
	}
	return topDecisionScore(patterns) / 100.0
}

func topDecisionScore(patterns []PatternReport) float64 {
	if len(patterns) == 0 {
		return 0
	}
	top := 3
	if len(patterns) < top {
		top = len(patterns)
	}
	sum := 0.0
	for i := 0; i < top; i++ {
		sum += patterns[i].DecisionScore
	}
	return sum / float64(top)
}

func collectCounterEvidence(evidence []identify.PatternEvidence) []string {
	out := make([]string, 0, 8)
	seen := make(map[string]struct{})
	for _, ev := range evidence {
		for _, c := range ev.ContradictionFactors {
			if _, ok := seen[c]; ok {
				continue
			}
			seen[c] = struct{}{}
			out = append(out, c)
		}
	}
	sort.Strings(out)
	return out
}

func forwardReturns(candles []*v1.Candlestick, pos int) (*float64, *float64, *float64) {
	return forwardReturn(candles, pos, 3), forwardReturn(candles, pos, 5), forwardReturn(candles, pos, 10)
}

func forwardReturn(candles []*v1.Candlestick, pos, horizon int) *float64 {
	if pos < 0 || pos >= len(candles) {
		return nil
	}
	target := pos + horizon
	if target >= len(candles) {
		return nil
	}
	entry := candles[pos].Close
	if entry == 0 {
		return nil
	}
	ret := (candles[target].Close - entry) / entry * 100
	return &ret
}

func evidenceKey(patternType string, pos int) string {
	return patternType + "#" + strconv.Itoa(pos)
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

func formatFloatPtr(v *float64) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%.4f", *v)
}

func stringsJoin(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += sep + parts[i]
	}
	return out
}
