package signal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigValidation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	raw := `{
		"score": {
			"pattern_weight": 50,
			"trend_weight": 20,
			"volume_weight": 20,
			"strong_threshold": 80,
			"medium_threshold": 60
		}
	}`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("expected validation error for invalid score weight sum")
	}
}
