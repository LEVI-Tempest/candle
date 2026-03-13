package charting

import (
	"strings"
	"testing"
	"time"

	"github.com/LEVI-Tempest/Candle/pkg/identify"
	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

func Test_Example(t *testing.T) {
	if err := (KlineExamples{}).Examples(); err != nil {
		t.Fatalf("examples render failed: %v", err)
	}
}

func TestEnhancedKlineWithPatternRecognition(t *testing.T) {
	// Create test candlestick data with various patterns
	// 创建包含各种形态的测试蜡烛图数据
	testData := createTestCandlestickData()

	// Create enhanced kline chart
	// 创建增强K线图
	enhancedKline := NewEnhancedKline()

	// Load test data
	// 加载测试数据
	enhancedKline.LoadData(testData)
	t.Logf("Loaded %d candlesticks", len(enhancedKline.Data))

	// Auto-detect patterns
	// 自动检测形态
	enhancedKline.AutoDetectPatterns()
	t.Logf("Detected %d patterns", len(enhancedKline.Patterns))

	// Print detected patterns
	// 打印检测到的形态
	for _, pattern := range enhancedKline.Patterns {
		t.Logf("Pattern: %s at position %d, price %.2f, strength %.2f, risk %.2f",
			pattern.Type, pattern.Position, pattern.Price, pattern.Strength, pattern.Risk)
	}

	// Get pattern summary
	// 获取形态摘要
	summary := enhancedKline.GetPatternSummary()
	t.Log("Pattern Summary:")
	for patternType, count := range summary {
		t.Logf("  %s: %d occurrences", patternType, count)
	}

	// Create chart with patterns
	// 创建带有形态的图表
	enhancedKline.CreateChart("Enhanced Candlestick Chart with Pattern Recognition")

	// Render to file
	// 渲染到文件
	err := enhancedKline.RenderToFile("enhanced_kline_with_patterns.html")
	if err != nil {
		t.Errorf("Failed to render chart to file: %v", err)
	} else {
		t.Log("Chart rendered to enhanced_kline_with_patterns.html")
	}

	// Verify that patterns were detected
	// 验证形态被检测到
	if len(enhancedKline.Patterns) == 0 {
		t.Error("Expected to detect some patterns, but none were found")
	}
	if len(enhancedKline.Evidences) == 0 {
		t.Error("Expected structured evidences, but none were found")
	}
}

func TestPatternColorAndSymbol(t *testing.T) {
	// Test pattern color assignment
	// 测试形态颜色分配
	testCases := []struct {
		patternType    string
		expectedColor  string
		expectedSymbol string
	}{
		{"Hammer", "#00da3c", "triangle"},
		{"Hanging Man", "#ec0000", "triangleDown"},
		{"Doji", "#ffaa00", "diamond"},
		{"Marubozu", "#9900cc", "star"},
		{"Rising Window", "#0066cc", "arrow"},
		{"Unknown Pattern", "#666666", "circle"},
	}

	for _, tc := range testCases {
		color := getPatternColor(tc.patternType)
		symbol := getPatternSymbol(tc.patternType)

		if color != tc.expectedColor {
			t.Errorf("Pattern %s: expected color %s, got %s", tc.patternType, tc.expectedColor, color)
		}

		if symbol != tc.expectedSymbol {
			t.Errorf("Pattern %s: expected symbol %s, got %s", tc.patternType, tc.expectedSymbol, symbol)
		}
	}
}

func TestFormatPatternLabelIncludesScoreAndReasons(t *testing.T) {
	ev := identify.PatternEvidence{
		PatternType: "Hammer",
		FinalScore:  0.84,
		VolumeFactors: []identify.FactorHit{
			{Name: "beiliang_confirm", Passed: true},
			{Name: "mfi_regime", Passed: true},
		},
	}

	label := formatPatternLabel("Hammer", ev)
	if !strings.Contains(label, "锤头") {
		t.Fatalf("label missing pattern name: %s", label)
	}
	if !strings.Contains(label, "84") {
		t.Fatalf("label missing score: %s", label)
	}
	if !strings.Contains(label, "beiliang_confirm") || !strings.Contains(label, "mfi_regime") {
		t.Fatalf("label missing reasons: %s", label)
	}
}

// createTestCandlestickData creates test data with various candlestick patterns
// 创建包含各种蜡烛图形态的测试数据
func createTestCandlestickData() []*v1.Candlestick {
	baseTime := time.Now().AddDate(0, 0, -30) // 30 days ago

	return []*v1.Candlestick{
		// Normal candles
		{Timestamp: baseTime.Unix(), Open: 100, High: 105, Low: 98, Close: 103, Volume: 1000},
		{Timestamp: baseTime.AddDate(0, 0, 1).Unix(), Open: 103, High: 108, Low: 101, Close: 106, Volume: 1200},

		// Doji pattern (十字星)
		{Timestamp: baseTime.AddDate(0, 0, 2).Unix(), Open: 106, High: 110, Low: 102, Close: 106.1, Volume: 800},

		// Hammer pattern (锤头线)
		{Timestamp: baseTime.AddDate(0, 0, 3).Unix(), Open: 105, High: 106, Low: 95, Close: 104, Volume: 1500},

		// Marubozu pattern (光头光脚)
		{Timestamp: baseTime.AddDate(0, 0, 4).Unix(), Open: 104, High: 115, Low: 104, Close: 115, Volume: 2000},

		// Shooting Star pattern (流星线)
		{Timestamp: baseTime.AddDate(0, 0, 5).Unix(), Open: 115, High: 125, Low: 114, Close: 116, Volume: 1800},

		// Bearish Engulfing setup - first candle (bullish)
		{Timestamp: baseTime.AddDate(0, 0, 6).Unix(), Open: 116, High: 120, Low: 115, Close: 119, Volume: 1300},
		// Bearish Engulfing - second candle (bearish, engulfs previous)
		{Timestamp: baseTime.AddDate(0, 0, 7).Unix(), Open: 121, High: 122, Low: 113, Close: 114, Volume: 2200},

		// Bullish Engulfing setup - first candle (bearish)
		{Timestamp: baseTime.AddDate(0, 0, 8).Unix(), Open: 114, High: 116, Low: 110, Close: 111, Volume: 1600},
		// Bullish Engulfing - second candle (bullish, engulfs previous)
		{Timestamp: baseTime.AddDate(0, 0, 9).Unix(), Open: 109, High: 118, Low: 108, Close: 117, Volume: 2100},

		// Morning Star setup - first candle (bearish)
		{Timestamp: baseTime.AddDate(0, 0, 10).Unix(), Open: 117, High: 118, Low: 110, Close: 111, Volume: 1700},
		// Morning Star - second candle (small body, gap down)
		{Timestamp: baseTime.AddDate(0, 0, 11).Unix(), Open: 108, High: 109, Low: 106, Close: 107, Volume: 900},
		// Morning Star - third candle (bullish, closes well into first)
		{Timestamp: baseTime.AddDate(0, 0, 12).Unix(), Open: 109, High: 120, Low: 108, Close: 118, Volume: 2300},

		// Three White Soldiers - first candle
		{Timestamp: baseTime.AddDate(0, 0, 13).Unix(), Open: 118, High: 125, Low: 117, Close: 124, Volume: 1800},
		// Three White Soldiers - second candle
		{Timestamp: baseTime.AddDate(0, 0, 14).Unix(), Open: 122, High: 130, Low: 121, Close: 129, Volume: 1900},
		// Three White Soldiers - third candle
		{Timestamp: baseTime.AddDate(0, 0, 15).Unix(), Open: 127, High: 135, Low: 126, Close: 134, Volume: 2000},

		// Spinning Top (陀螺线)
		{Timestamp: baseTime.AddDate(0, 0, 16).Unix(), Open: 134, High: 140, Low: 128, Close: 135, Volume: 1400},

		// More normal candles to complete the dataset
		{Timestamp: baseTime.AddDate(0, 0, 17).Unix(), Open: 135, High: 138, Low: 132, Close: 136, Volume: 1100},
		{Timestamp: baseTime.AddDate(0, 0, 18).Unix(), Open: 136, High: 140, Low: 134, Close: 139, Volume: 1300},
		{Timestamp: baseTime.AddDate(0, 0, 19).Unix(), Open: 139, High: 142, Low: 137, Close: 141, Volume: 1200},
		{Timestamp: baseTime.AddDate(0, 0, 20).Unix(), Open: 141, High: 145, Low: 139, Close: 143, Volume: 1500},
	}
}
