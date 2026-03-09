package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/LEVI-Tempest/Candle/pkg/charting"
	"github.com/LEVI-Tempest/Candle/pkg/datasource"
	"github.com/LEVI-Tempest/Candle/pkg/identify"
	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

type signalPattern struct {
	Type      string  `json:"type"`
	Direction string  `json:"direction"`
	Position  int     `json:"position"`
	Strength  float64 `json:"strength"`
	Risk      float64 `json:"risk"`
	Price     float64 `json:"price"`
	Time      string  `json:"time"`
}

type signalReport struct {
	Symbol          string                     `json:"symbol"`
	AsOf            string                     `json:"as_of"`
	Source          string                     `json:"source"`
	Trend           string                     `json:"trend"`
	Score           float64                    `json:"score"`
	Patterns        []signalPattern            `json:"patterns"`
	Evidence        []identify.PatternEvidence `json:"evidence"`
	CounterEvidence []string                   `json:"counter_evidence"`
	InvalidIf       []string                   `json:"invalid_if"`
}

type candleInputEnvelope struct {
	Symbol string            `json:"symbol"`
	Source string            `json:"source"`
	Data   []*v1.Candlestick `json:"data"`
}

func main() {
	inputPath := flag.String("input", "", "Input JSON file path. Supports []Candlestick or {symbol,source,data}.")
	outputPath := flag.String("output", "", "Output JSON file path. Empty prints to stdout.")
	asOf := flag.String("as-of", time.Now().Format(time.RFC3339), "As-of timestamp (RFC3339).")
	symbol := flag.String("symbol", "XSHE:300059", "Symbol for reporting context.")

	fetch := flag.Bool("fetch", false, "Fetch data from Tsanghi instead of --input.")
	exchange := flag.String("exchange", datasource.TsanghiXSHE, "Exchange: XSHE | XSHG | XHKG")
	ticker := flag.String("ticker", "300059", "Ticker code")
	token := flag.String("token", "demo", "Tsanghi API token")
	limit := flag.Int("limit", 120, "Number of candles to fetch")
	flag.Parse()

	candles, source, detectedSymbol, err := loadCandles(*inputPath, *fetch, *exchange, *ticker, *token, *limit)
	if err != nil {
		exitf("load candles failed: %v", err)
	}

	if detectedSymbol != "" {
		*symbol = detectedSymbol
	}
	if len(candles) == 0 {
		exitf("no candles available")
	}

	report := buildSignalReport(*symbol, *asOf, source, candles)
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		exitf("marshal report failed: %v", err)
	}

	if *outputPath == "" {
		fmt.Println(string(data))
		return
	}
	if err := os.WriteFile(*outputPath, data, 0o644); err != nil {
		exitf("write output failed: %v", err)
	}
}

func loadCandles(
	inputPath string,
	fetch bool,
	exchange, ticker, token string,
	limit int,
) ([]*v1.Candlestick, string, string, error) {
	if fetch {
		client := datasource.NewTsanghiClient(token)
		candles, err := client.Fetch(exchange, ticker, &datasource.FetchOptions{
			Limit: limit,
			Order: 2,
		})
		if err != nil {
			return nil, "", "", err
		}
		return candles, "tsanghi", fmt.Sprintf("%s:%s", exchange, ticker), nil
	}
	if inputPath == "" {
		return nil, "", "", fmt.Errorf("either --input or --fetch is required")
	}

	raw, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, "", "", err
	}

	var envelope candleInputEnvelope
	if err := json.Unmarshal(raw, &envelope); err == nil && len(envelope.Data) > 0 {
		src := envelope.Source
		if src == "" {
			src = "file"
		}
		return envelope.Data, src, envelope.Symbol, nil
	}

	var candles []*v1.Candlestick
	if err := json.Unmarshal(raw, &candles); err != nil {
		return nil, "", "", fmt.Errorf("input must be []Candlestick or {symbol,source,data}: %w", err)
	}
	return candles, "file", "", nil
}

func buildSignalReport(symbol, asOf, source string, candles []*v1.Candlestick) signalReport {
	ek := charting.NewEnhancedKline()
	ek.LoadData(candles)
	ek.AutoDetectPatterns()

	patterns := make([]signalPattern, 0, len(ek.Patterns))
	for _, p := range ek.Patterns {
		patterns = append(patterns, signalPattern{
			Type:      p.Type,
			Direction: patternDirection(p.Type),
			Position:  p.Position,
			Strength:  p.Strength,
			Risk:      p.Risk,
			Price:     p.Price,
			Time:      p.Time,
		})
	}

	trend := "unknown"
	if len(ek.Data) >= 3 {
		period := 20
		if len(ek.Data) < period {
			period = len(ek.Data)
		}
		t, _ := identify.AnalyzeLongTermTrend(ek.Data, period, 3.0)
		trend = normalizeTrend(t)
	}

	evidence := ek.Evidences
	sort.Slice(evidence, func(i, j int) bool {
		if evidence[i].FinalScore == evidence[j].FinalScore {
			return evidence[i].Position < evidence[j].Position
		}
		return evidence[i].FinalScore > evidence[j].FinalScore
	})

	return signalReport{
		Symbol:          symbol,
		AsOf:            asOf,
		Source:          source,
		Trend:           trend,
		Score:           calculateOverallScore(evidence),
		Patterns:        patterns,
		Evidence:        evidence,
		CounterEvidence: collectCounterEvidence(evidence),
		InvalidIf: []string{
			"data source has missing/incorrect OHLCV records",
			"next trading sessions show no volume confirmation",
			"price breaks pattern invalidation level with high volatility",
		},
	}
}

func calculateOverallScore(evidence []identify.PatternEvidence) float64 {
	if len(evidence) == 0 {
		return 0
	}
	top := 3
	if len(evidence) < top {
		top = len(evidence)
	}
	sum := 0.0
	for i := 0; i < top; i++ {
		sum += evidence[i].FinalScore
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

func normalizeTrend(t identify.Trend) string {
	switch t {
	case identify.TrendYang:
		return "up"
	case identify.TrendYin:
		return "down"
	case identify.TrendMiddle:
		return "range"
	default:
		return "unknown"
	}
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

func exitf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "signal: "+format+"\n", args...)
	os.Exit(1)
}
