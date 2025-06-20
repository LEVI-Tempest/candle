package main

import (
	"fmt"
	"log"
	"time"

	"github.com/LEVI-Tempest/Candle/pkg/charting"
	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

func main() {
	fmt.Println("🎯 Precise Pattern Position Markers Demo")
	fmt.Println("=========================================")
	fmt.Println("This demo creates an HTML chart with patterns marked at their EXACT positions!")
	fmt.Println()

	// Create sample data with clear, identifiable patterns at specific positions
	// 创建在特定位置包含清晰可识别形态的示例数据
	candleData := createPrecisePatternData()

	// Create enhanced kline chart
	// 创建增强K线图
	enhancedKline := charting.NewEnhancedKline()

	// Load the data
	// 加载数据
	enhancedKline.LoadData(candleData)
	fmt.Printf("📊 Loaded %d candlesticks\n", len(enhancedKline.Data))

	// Auto-detect patterns
	// 自动检测形态
	enhancedKline.AutoDetectPatterns()
	fmt.Printf("🔍 Detected %d patterns\n\n", len(enhancedKline.Patterns))

	// Show detailed pattern information with exact positions
	// 显示包含确切位置的详细形态信息
	fmt.Println("📍 Detected Patterns with Exact Positions:")
	fmt.Println("-------------------------------------------")
	for i, pattern := range enhancedKline.Patterns {
		emoji := getPatternEmoji(pattern.Type)
		fmt.Printf("%d. %s %s\n", i+1, emoji, pattern.Type)
		fmt.Printf("   📍 Position: Day %d (X-axis position)\n", pattern.Position)
		fmt.Printf("   💰 Price: %.2f (Y-axis position)\n", pattern.Price)
		fmt.Printf("   💪 Strength: %.1f/1.0\n", pattern.Strength)
		fmt.Printf("   ⚠️  Risk: %.1f/1.0\n", pattern.Risk)
		fmt.Printf("   🕐 Time: %s\n", pattern.Time)
		fmt.Println()
	}

	// Create chart with precise pattern position markers
	// 创建带有精确形态位置标记的图表
	enhancedKline.CreateChart("🎯 Candlestick Chart with Precise Pattern Position Markers")

	// Render to HTML file with enhanced JavaScript markers
	// 渲染到HTML文件并添加增强的JavaScript标记
	filename := "precise_pattern_markers.html"
	err := enhancedKline.RenderToFile(filename)
	if err != nil {
		log.Fatalf("❌ Failed to render chart: %v", err)
	}

	fmt.Printf("✅ Enhanced chart successfully created: %s\n", filename)
	fmt.Println()
	fmt.Println("🌐 What you'll see in the HTML chart:")
	fmt.Println("   • Interactive candlestick chart with zoom/pan")
	fmt.Println("   • Pattern markers at EXACT positions on the chart")
	fmt.Println("   • Pattern names and strength displayed directly on chart")
	fmt.Println("   • Color-coded pattern types for easy identification")
	fmt.Println("   • Pattern legend in the top-right corner")
	fmt.Println("   • Support/resistance lines for high-strength patterns")
	fmt.Println()
	fmt.Println("🎨 Pattern Markers:")
	fmt.Println("   • Each pattern is marked at its exact candlestick position")
	fmt.Println("   • Markers show pattern name and strength rating")
	fmt.Println("   • Colors indicate pattern type (bullish/bearish/neutral)")
	fmt.Println("   • JavaScript enhancement adds precise positioning")
	fmt.Println()
	fmt.Println("📖 How to use:")
	fmt.Println("   1. Open the HTML file in your web browser")
	fmt.Println("   2. Look for colored text markers above specific candlesticks")
	fmt.Println("   3. Each marker shows exactly where a pattern was detected")
	fmt.Println("   4. Use the legend to understand pattern types and counts")
	fmt.Println("   5. Zoom in to see pattern details more clearly")
}

// getPatternEmoji returns an emoji for each pattern type
// 为每种形态类型返回表情符号
func getPatternEmoji(patternType string) string {
	switch patternType {
	// Bullish patterns
	case "Hammer", "Inverted Hammer", "Bullish Engulfing", "Piercing Line", "Morning Star", "Three White Soldiers":
		return "🟢"
	// Bearish patterns
	case "Hanging Man", "Shooting Star", "Bearish Engulfing", "Dark Cloud Cover", "Evening Star", "Three Black Crows":
		return "🔴"
	// Neutral/Reversal patterns
	case "Doji", "Spinning Top", "Tweezer Tops", "Tweezer Bottoms":
		return "🟠"
	// Gap patterns
	case "Rising Window", "Falling Window":
		return "🔵"
	// Strong patterns
	case "Marubozu":
		return "🟣"
	default:
		return "⚪"
	}
}

// createPrecisePatternData creates sample data with patterns at known positions
// 创建在已知位置包含形态的示例数据
func createPrecisePatternData() []*v1.Candlestick {
	baseTime := time.Now().AddDate(0, 0, -25) // 25 days ago
	
	return []*v1.Candlestick{
		// Position 0-2: Normal uptrend
		{Timestamp: baseTime.Unix(), Open: 100, High: 105, Low: 98, Close: 103, Volume: 1000},
		{Timestamp: baseTime.AddDate(0, 0, 1).Unix(), Open: 103, High: 108, Low: 101, Close: 106, Volume: 1200},
		{Timestamp: baseTime.AddDate(0, 0, 2).Unix(), Open: 106, High: 110, Low: 104, Close: 109, Volume: 1100},
		
		// Position 3: 🟠 Doji - Market indecision (十字星 - 市场犹豫)
		{Timestamp: baseTime.AddDate(0, 0, 3).Unix(), Open: 109, High: 114, Low: 104, Close: 109.3, Volume: 800},
		
		// Position 4: Pullback
		{Timestamp: baseTime.AddDate(0, 0, 4).Unix(), Open: 109, High: 111, Low: 102, Close: 104, Volume: 1300},
		
		// Position 5: 🟢 Hammer - Bullish reversal signal (锤头线 - 看涨反转信号)
		{Timestamp: baseTime.AddDate(0, 0, 5).Unix(), Open: 104, High: 106, Low: 92, Close: 105, Volume: 1800},
		
		// Position 6: 🟣 Bullish Marubozu - Strong buying (看涨光头光脚 - 强烈买入)
		{Timestamp: baseTime.AddDate(0, 0, 6).Unix(), Open: 105, High: 122, Low: 105, Close: 122, Volume: 2500},
		
		// Position 7: Continuation
		{Timestamp: baseTime.AddDate(0, 0, 7).Unix(), Open: 122, High: 128, Low: 120, Close: 126, Volume: 1600},
		
		// Position 8: 🔴 Shooting Star - Bearish reversal warning (流星线 - 看跌反转警告)
		{Timestamp: baseTime.AddDate(0, 0, 8).Unix(), Open: 126, High: 142, Low: 125, Close: 128, Volume: 2000},
		
		// Position 9-10: 🔴 Bearish Engulfing Pattern (看跌吞噬形态)
		{Timestamp: baseTime.AddDate(0, 0, 9).Unix(), Open: 128, High: 132, Low: 127, Close: 131, Volume: 1400},  // First candle (bullish)
		{Timestamp: baseTime.AddDate(0, 0, 10).Unix(), Open: 133, High: 134, Low: 120, Close: 122, Volume: 2800}, // Second candle (bearish engulfing)
		
		// Position 11-12: Downtrend
		{Timestamp: baseTime.AddDate(0, 0, 11).Unix(), Open: 122, High: 124, Low: 115, Close: 117, Volume: 1900},
		{Timestamp: baseTime.AddDate(0, 0, 12).Unix(), Open: 117, High: 119, Low: 110, Close: 112, Volume: 1700},
		
		// Position 13-15: 🟢 Morning Star Pattern - Bullish reversal (启明星形态 - 看涨反转)
		{Timestamp: baseTime.AddDate(0, 0, 13).Unix(), Open: 112, High: 113, Low: 105, Close: 107, Volume: 1800}, // First candle (bearish)
		{Timestamp: baseTime.AddDate(0, 0, 14).Unix(), Open: 104, High: 105, Low: 102, Close: 103, Volume: 900},  // Second candle (small, gap down)
		{Timestamp: baseTime.AddDate(0, 0, 15).Unix(), Open: 105, High: 118, Low: 104, Close: 116, Volume: 2400}, // Third candle (bullish)
		
		// Position 16-18: 🟢 Three White Soldiers - Strong bullish continuation (红三兵 - 强烈看涨延续)
		{Timestamp: baseTime.AddDate(0, 0, 16).Unix(), Open: 116, High: 124, Low: 115, Close: 123, Volume: 2000}, // First soldier
		{Timestamp: baseTime.AddDate(0, 0, 17).Unix(), Open: 121, High: 130, Low: 120, Close: 129, Volume: 2100}, // Second soldier
		{Timestamp: baseTime.AddDate(0, 0, 18).Unix(), Open: 127, High: 136, Low: 126, Close: 135, Volume: 2200}, // Third soldier
		
		// Position 19: Peak
		{Timestamp: baseTime.AddDate(0, 0, 19).Unix(), Open: 135, High: 142, Low: 133, Close: 140, Volume: 1800},
		
		// Position 20: 🟠 Spinning Top - Indecision at top (陀螺线 - 顶部犹豫)
		{Timestamp: baseTime.AddDate(0, 0, 20).Unix(), Open: 140, High: 147, Low: 133, Close: 141, Volume: 1500},
		
		// Position 21-23: 🔴 Evening Star Pattern - Bearish reversal (黄昏之星 - 看跌反转)
		{Timestamp: baseTime.AddDate(0, 0, 21).Unix(), Open: 141, High: 148, Low: 140, Close: 147, Volume: 1700}, // First candle (bullish)
		{Timestamp: baseTime.AddDate(0, 0, 22).Unix(), Open: 149, High: 150, Low: 148, Close: 149.2, Volume: 800}, // Second candle (small, gap up)
		{Timestamp: baseTime.AddDate(0, 0, 23).Unix(), Open: 148, High: 149, Low: 138, Close: 141, Volume: 2600}, // Third candle (bearish)
		
		// Position 24: Final candle
		{Timestamp: baseTime.AddDate(0, 0, 24).Unix(), Open: 141, High: 145, Low: 138, Close: 142, Volume: 1400},
	}
}
