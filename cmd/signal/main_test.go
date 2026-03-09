package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCandlesFromArrayInput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "candles.json")
	content := `[
		{"timestamp":1700000000,"open":10,"high":11,"low":9,"close":10.5,"volume":100},
		{"timestamp":1700086400,"open":10.5,"high":11.2,"low":10.1,"close":11,"volume":120}
	]`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	candles, source, symbol, err := loadCandles(path, false, "", "", "", 0)
	if err != nil {
		t.Fatalf("load candles failed: %v", err)
	}
	if len(candles) != 2 {
		t.Fatalf("expected 2 candles, got %d", len(candles))
	}
	if source != "file" {
		t.Fatalf("expected source=file, got %s", source)
	}
	if symbol != "" {
		t.Fatalf("expected empty symbol, got %s", symbol)
	}
}
