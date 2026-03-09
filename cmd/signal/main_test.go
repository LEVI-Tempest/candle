package main

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/LEVI-Tempest/Candle/pkg/identify"
	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

func TestCalculateOverallScore(t *testing.T) {
	ev := []identify.PatternEvidence{
		{FinalScore: 0.9},
		{FinalScore: 0.8},
		{FinalScore: 0.7},
		{FinalScore: 0.1},
	}
	got := calculateOverallScore(ev)
	want := 0.8
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("expected score %.2f, got %.2f", want, got)
	}
}

func TestCollectCounterEvidence(t *testing.T) {
	ev := []identify.PatternEvidence{
		{ContradictionFactors: []string{"a", "b"}},
		{ContradictionFactors: []string{"b", "c"}},
	}
	got := collectCounterEvidence(ev)
	if len(got) != 3 {
		t.Fatalf("expected 3 items, got %d: %v", len(got), got)
	}
	if got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("unexpected order/content: %v", got)
	}
}

func TestSignalReportJSONContract(t *testing.T) {
	baseTime := time.Now().AddDate(0, 0, -12)
	candles := []*v1.Candlestick{
		{Timestamp: baseTime.Unix(), Open: 100, High: 102, Low: 99, Close: 101, Volume: 1000},
		{Timestamp: baseTime.AddDate(0, 0, 1).Unix(), Open: 101, High: 103, Low: 100, Close: 102, Volume: 1100},
		{Timestamp: baseTime.AddDate(0, 0, 2).Unix(), Open: 102, High: 103, Low: 98, Close: 99, Volume: 1300},
		{Timestamp: baseTime.AddDate(0, 0, 3).Unix(), Open: 99, High: 100, Low: 95, Close: 96, Volume: 1500},
		{Timestamp: baseTime.AddDate(0, 0, 4).Unix(), Open: 96, High: 105, Low: 95, Close: 104, Volume: 2500},
		{Timestamp: baseTime.AddDate(0, 0, 5).Unix(), Open: 104, High: 108, Low: 103, Close: 107, Volume: 2300},
		{Timestamp: baseTime.AddDate(0, 0, 6).Unix(), Open: 107, High: 108, Low: 102, Close: 103, Volume: 1700},
		{Timestamp: baseTime.AddDate(0, 0, 7).Unix(), Open: 103, High: 104, Low: 99, Close: 100, Volume: 1600},
	}

	report := buildSignalReport("XSHE:300059", "2026-03-09T09:30:00Z", "test", candles)
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal report failed: %v", err)
	}

	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		t.Fatalf("unmarshal report failed: %v", err)
	}

	requiredKeys := []string{
		"symbol", "as_of", "source", "trend", "score",
		"patterns", "evidence", "counter_evidence", "invalid_if",
	}
	for _, k := range requiredKeys {
		if _, ok := obj[k]; !ok {
			t.Fatalf("missing required key: %s", k)
		}
	}

	score, ok := obj["score"].(float64)
	if !ok {
		t.Fatalf("score should be number, got %T", obj["score"])
	}
	if score < 0 || score > 1 {
		t.Fatalf("score out of range: %f", score)
	}
}
