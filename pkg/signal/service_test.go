package signal

import (
	"encoding/json"
	"testing"
	"time"

	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

func TestBuildReportContainsContractFields(t *testing.T) {
	base := time.Now().AddDate(0, 0, -15)
	candles := []*v1.Candlestick{
		{Timestamp: base.Unix(), Open: 100, High: 102, Low: 99, Close: 101, Volume: 1000},
		{Timestamp: base.AddDate(0, 0, 1).Unix(), Open: 101, High: 103, Low: 100, Close: 102, Volume: 1100},
		{Timestamp: base.AddDate(0, 0, 2).Unix(), Open: 102, High: 103, Low: 98, Close: 99, Volume: 1300},
		{Timestamp: base.AddDate(0, 0, 3).Unix(), Open: 99, High: 100, Low: 95, Close: 96, Volume: 1500},
		{Timestamp: base.AddDate(0, 0, 4).Unix(), Open: 96, High: 106, Low: 95, Close: 105, Volume: 2600},
		{Timestamp: base.AddDate(0, 0, 5).Unix(), Open: 105, High: 108, Low: 103, Close: 107, Volume: 2300},
		{Timestamp: base.AddDate(0, 0, 6).Unix(), Open: 107, High: 109, Low: 106, Close: 108, Volume: 2200},
		{Timestamp: base.AddDate(0, 0, 7).Unix(), Open: 108, High: 109, Low: 103, Close: 104, Volume: 1800},
		{Timestamp: base.AddDate(0, 0, 8).Unix(), Open: 104, High: 105, Low: 100, Close: 101, Volume: 1700},
		{Timestamp: base.AddDate(0, 0, 9).Unix(), Open: 101, High: 102, Low: 98, Close: 99, Volume: 1600},
	}

	cfg := DefaultConfig()
	report := BuildReport("XSHE:300059", "2026-03-09T09:30:00Z", "test", candles, cfg)

	if report.Symbol == "" || report.AsOf == "" || report.Source == "" {
		t.Fatalf("missing base fields: %+v", report)
	}
	if report.Score < 0 || report.Score > 1 {
		t.Fatalf("normalized score out of range: %f", report.Score)
	}
	if report.DecisionScore < 0 || report.DecisionScore > 100 {
		t.Fatalf("decision score out of range: %f", report.DecisionScore)
	}
	if report.DecisionLevel == "" {
		t.Fatal("decision level should not be empty")
	}

	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal report failed: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal report failed: %v", err)
	}
	required := []string{
		"symbol", "as_of", "source", "trend", "score",
		"decision_score", "decision_level", "patterns",
		"evidence", "counter_evidence", "invalid_if",
	}
	for _, k := range required {
		if _, ok := payload[k]; !ok {
			t.Fatalf("missing required key: %s", k)
		}
	}
}

func TestDetermineTrendByMA20(t *testing.T) {
	cfg := DefaultConfig()
	now := time.Now().AddDate(0, 0, -40)
	candles := make([]*v1.Candlestick, 0, 30)
	for i := 0; i < 30; i++ {
		price := 100.0 + float64(i)
		candles = append(candles, &v1.Candlestick{
			Timestamp: now.AddDate(0, 0, i).Unix(),
			Open:      price - 0.5,
			High:      price + 1,
			Low:       price - 1,
			Close:     price,
			Volume:    1000 + float64(i*10),
		})
	}

	report := BuildReport("XSHE:000001", "2026-03-09T09:30:00Z", "test", candles, cfg)
	if report.Trend != "up" {
		t.Fatalf("expected trend up, got %s", report.Trend)
	}
}
