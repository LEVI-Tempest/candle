package signal

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/LEVI-Tempest/Candle/pkg/identify"
)

// TrendConfig controls trend filter behavior.
// TrendConfig 控制趋势过滤行为。
type TrendConfig struct {
	Period int `json:"period"`
}

// ScoreConfig controls weighted decision scoring.
// ScoreConfig 控制加权决策评分。
type ScoreConfig struct {
	PatternWeight   float64 `json:"pattern_weight"`
	TrendWeight     float64 `json:"trend_weight"`
	VolumeWeight    float64 `json:"volume_weight"`
	StrongThreshold float64 `json:"strong_threshold"`
	MediumThreshold float64 `json:"medium_threshold"`
}

// Config is the top-level configuration for signal generation.
// Config 是信号生成的顶层配置。
type Config struct {
	Trend      TrendConfig             `json:"trend"`
	Score      ScoreConfig             `json:"score"`
	Evidence   identify.EvidenceConfig `json:"evidence"`
	LogCSVPath string                  `json:"log_csv_path"`
}

// DefaultConfig returns default values for local research workflow.
// DefaultConfig 返回适合个人研究工作流的默认值。
func DefaultConfig() Config {
	return Config{
		Trend: TrendConfig{
			Period: 20,
		},
		Score: ScoreConfig{
			PatternWeight:   60,
			TrendWeight:     20,
			VolumeWeight:    20,
			StrongThreshold: 80,
			MediumThreshold: 60,
		},
		Evidence:   identify.DefaultEvidenceConfig(),
		LogCSVPath: filepath.Join("data", "signal_log.csv"),
	}
}

// LoadConfig loads config from JSON file and merges with defaults.
// LoadConfig 从 JSON 文件加载配置并与默认值合并。
func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var userCfg Config
	if err := json.Unmarshal(raw, &userCfg); err != nil {
		return Config{}, err
	}
	mergeConfig(&cfg, &userCfg)
	return cfg, nil
}

func mergeConfig(dst, src *Config) {
	if src.Trend.Period > 0 {
		dst.Trend.Period = src.Trend.Period
	}

	if src.Score.PatternWeight > 0 {
		dst.Score.PatternWeight = src.Score.PatternWeight
	}
	if src.Score.TrendWeight > 0 {
		dst.Score.TrendWeight = src.Score.TrendWeight
	}
	if src.Score.VolumeWeight > 0 {
		dst.Score.VolumeWeight = src.Score.VolumeWeight
	}
	if src.Score.StrongThreshold > 0 {
		dst.Score.StrongThreshold = src.Score.StrongThreshold
	}
	if src.Score.MediumThreshold > 0 {
		dst.Score.MediumThreshold = src.Score.MediumThreshold
	}

	if src.Evidence.VolumeLookback > 0 {
		dst.Evidence.VolumeLookback = src.Evidence.VolumeLookback
	}
	if src.Evidence.MFIPeriod > 0 {
		dst.Evidence.MFIPeriod = src.Evidence.MFIPeriod
	}
	if src.Evidence.CMFPeriod > 0 {
		dst.Evidence.CMFPeriod = src.Evidence.CMFPeriod
	}
	if src.Evidence.VolumeBoostThreshold > 0 {
		dst.Evidence.VolumeBoostThreshold = src.Evidence.VolumeBoostThreshold
	}
	if src.Evidence.VolumeShrinkThreshold > 0 {
		dst.Evidence.VolumeShrinkThreshold = src.Evidence.VolumeShrinkThreshold
	}
	if src.Evidence.BaseWeight > 0 {
		dst.Evidence.BaseWeight = src.Evidence.BaseWeight
	}
	if src.Evidence.ContextWeight > 0 {
		dst.Evidence.ContextWeight = src.Evidence.ContextWeight
	}
	if src.Evidence.VolumeWeight > 0 {
		dst.Evidence.VolumeWeight = src.Evidence.VolumeWeight
	}

	if src.LogCSVPath != "" {
		dst.LogCSVPath = src.LogCSVPath
	}
}
