package main

import (
	"math"
	"testing"

	"github.com/LEVI-Tempest/Candle/pkg/identify"
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
