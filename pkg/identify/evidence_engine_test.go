package identify

import (
	"testing"

	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

func TestBuildPatternEvidence(t *testing.T) {
	candles := []CandlestickWrapper{
		NewCandlestickWrapper(&v1.Candlestick{Open: 10, High: 11, Low: 9.8, Close: 10.6, Volume: 100}),
		NewCandlestickWrapper(&v1.Candlestick{Open: 10.6, High: 10.8, Low: 10.2, Close: 10.3, Volume: 110}),
		NewCandlestickWrapper(&v1.Candlestick{Open: 10.3, High: 10.5, Low: 9.9, Close: 10.0, Volume: 120}),
		NewCandlestickWrapper(&v1.Candlestick{Open: 10.0, High: 10.2, Low: 9.4, Close: 9.7, Volume: 180}),
		NewCandlestickWrapper(&v1.Candlestick{Open: 9.8, High: 10.8, Low: 9.7, Close: 10.7, Volume: 260}),
		NewCandlestickWrapper(&v1.Candlestick{Open: 10.7, High: 11.1, Low: 10.6, Close: 11.0, Volume: 240}),
	}

	signals := []PatternSignal{
		{
			Type:      "Bullish Engulfing",
			Direction: "bullish",
			Position:  4,
			Strength:  0.8,
			Risk:      0.3,
			Price:     candles[4].Close,
			Time:      "2026-03-09 00:00:00",
		},
	}

	evidences := BuildPatternEvidence(signals, candles, DefaultEvidenceConfig())
	if len(evidences) != 1 {
		t.Fatalf("expected 1 evidence, got %d", len(evidences))
	}

	ev := evidences[0]
	if ev.PatternType != "Bullish Engulfing" {
		t.Fatalf("unexpected pattern type: %s", ev.PatternType)
	}
	if ev.FinalScore < 0 || ev.FinalScore > 1 {
		t.Fatalf("final score out of range: %f", ev.FinalScore)
	}
	if ev.ConfidenceLevel == "" {
		t.Fatal("expected non-empty confidence level")
	}
	if len(ev.VolumeFactors) == 0 {
		t.Fatal("expected volume factors to be populated")
	}
}
