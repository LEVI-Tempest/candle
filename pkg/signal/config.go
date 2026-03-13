package signal

import (
	"encoding/json"
	"fmt"
	"math"
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
		if err := validateConfig(cfg); err != nil {
			return Config{}, err
		}
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
	if err := validateConfig(cfg); err != nil {
		return Config{}, err
	}
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
	if src.Evidence.BeiliangThreshold > 0 {
		dst.Evidence.BeiliangThreshold = src.Evidence.BeiliangThreshold
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
	if src.Evidence.ContextWindow > 0 {
		dst.Evidence.ContextWindow = src.Evidence.ContextWindow
	}
	if src.Evidence.ContextTrendThreshold > 0 {
		dst.Evidence.ContextTrendThreshold = src.Evidence.ContextTrendThreshold
	}
	if src.Evidence.OBVDivergenceLookback > 0 {
		dst.Evidence.OBVDivergenceLookback = src.Evidence.OBVDivergenceLookback
	}

	if src.LogCSVPath != "" {
		dst.LogCSVPath = src.LogCSVPath
	}
}

func validateConfig(cfg Config) error {
	if cfg.Trend.Period < 2 {
		return fmt.Errorf("trend.period must be >= 2")
	}
	totalWeight := cfg.Score.PatternWeight + cfg.Score.TrendWeight + cfg.Score.VolumeWeight
	if math.Abs(totalWeight-100) > 1e-9 {
		return fmt.Errorf("score weights must sum to 100, got %.2f", totalWeight)
	}
	if cfg.Score.StrongThreshold < cfg.Score.MediumThreshold {
		return fmt.Errorf("strong_threshold must be >= medium_threshold")
	}
	if cfg.Score.StrongThreshold > 100 || cfg.Score.StrongThreshold < 0 {
		return fmt.Errorf("strong_threshold must be within [0,100]")
	}
	if cfg.Score.MediumThreshold > 100 || cfg.Score.MediumThreshold < 0 {
		return fmt.Errorf("medium_threshold must be within [0,100]")
	}
	if cfg.Evidence.BaseWeight+cfg.Evidence.ContextWeight+cfg.Evidence.VolumeWeight <= 0 {
		return fmt.Errorf("evidence weights sum must be > 0")
	}
	if cfg.Evidence.BeiliangThreshold <= 0 {
		return fmt.Errorf("evidence.beiliang_threshold must be > 0")
	}
	return nil
}
