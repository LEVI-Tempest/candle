package identify

import (
	"testing"

	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

func TestAnalyzeVolumePriceSignals_HighVolumeLongLeggedDoji(t *testing.T) {
	candles := []CandlestickWrapper{
		NewCandlestickWrapper(&v1.Candlestick{Open: 100, Close: 101, High: 102, Low: 99, Volume: 1000}),
		NewCandlestickWrapper(&v1.Candlestick{Open: 101, Close: 100.5, High: 102, Low: 100, Volume: 1100}),
		NewCandlestickWrapper(&v1.Candlestick{Open: 100.5, Close: 100.4, High: 112, Low: 88, Volume: 2200}),
	}

	sigs := AnalyzeVolumePriceSignals(candles, 2)
	found := false
	for _, s := range sigs {
		if s.Type == "High-Volume Long-Legged Doji" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected High-Volume Long-Legged Doji signal, got %+v", sigs)
	}
}
