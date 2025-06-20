package main

import (
	"fmt"
	"log"
	"time"

	"github.com/LEVI-Tempest/Candle/pkg/charting"
	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

func main() {
	fmt.Println("🧪 Testing Pattern Markers")
	fmt.Println("==========================")

	// Create simple test data with obvious patterns
	// 创建包含明显形态的简单测试数据
	testData := createSimpleTestData()

	// Create enhanced kline chart
	// 创建增强K线图
	chart := charting.NewEnhancedKline()

	// Load data
	// 加载数据
	chart.LoadData(testData)
	fmt.Printf("📊 Loaded %d candlesticks\n", len(chart.Data))

	// Detect patterns
	// 检测形态
	chart.AutoDetectPatterns()
	fmt.Printf("🔍 Detected %d patterns\n", len(chart.Patterns))

	// Show detected patterns
	// 显示检测到的形态
	fmt.Println("\n📋 Detected Patterns:")
	for i, pattern := range chart.Patterns {
		fmt.Printf("%d. %s at position %d, price %.2f, strength %.1f\n", 
			i+1, pattern.Type, pattern.Position, pattern.Price, pattern.Strength)
	}

	// Create chart
	// 创建图表
	chart.CreateChart("🧪 Test Chart with Pattern Markers")

	// Render to file
	// 渲染到文件
	filename := "test_markers.html"
	err := chart.RenderToFile(filename)
	if err != nil {
		log.Fatalf("❌ Failed to create chart: %v", err)
	}

	fmt.Printf("\n✅ Test chart created: %s\n", filename)
	fmt.Println("\n📖 Instructions:")
	fmt.Println("1. Open test_markers.html in your web browser")
	fmt.Println("2. Open browser developer tools (F12)")
	fmt.Println("3. Check the Console tab for JavaScript messages")
	fmt.Println("4. Look for pattern markers on the chart")
	fmt.Println("5. Check if the legend appears in the top-right corner")
}

// createSimpleTestData creates very simple test data with clear patterns
// 创建包含清晰形态的非常简单的测试数据
func createSimpleTestData() []*v1.Candlestick {
	baseTime := time.Now().AddDate(0, 0, -10) // 10 days ago
	
	return []*v1.Candlestick{
		// Normal candles
		{Timestamp: baseTime.Unix(), Open: 100, High: 105, Low: 98, Close: 103, Volume: 1000},
		{Timestamp: baseTime.AddDate(0, 0, 1).Unix(), Open: 103, High: 108, Low: 101, Close: 106, Volume: 1200},
		
		// Clear Doji pattern
		{Timestamp: baseTime.AddDate(0, 0, 2).Unix(), Open: 106, High: 110, Low: 102, Close: 106.1, Volume: 800},
		
		// Clear Hammer pattern
		{Timestamp: baseTime.AddDate(0, 0, 3).Unix(), Open: 105, High: 107, Low: 90, Close: 106, Volume: 1500},
		
		// Clear Marubozu pattern
		{Timestamp: baseTime.AddDate(0, 0, 4).Unix(), Open: 106, High: 120, Low: 106, Close: 120, Volume: 2000},
		
		// More normal candles
		{Timestamp: baseTime.AddDate(0, 0, 5).Unix(), Open: 120, High: 125, Low: 118, Close: 123, Volume: 1300},
		{Timestamp: baseTime.AddDate(0, 0, 6).Unix(), Open: 123, High: 128, Low: 121, Close: 126, Volume: 1400},
		{Timestamp: baseTime.AddDate(0, 0, 7).Unix(), Open: 126, High: 130, Low: 124, Close: 128, Volume: 1200},
		{Timestamp: baseTime.AddDate(0, 0, 8).Unix(), Open: 128, High: 132, Low: 126, Close: 130, Volume: 1100},
		{Timestamp: baseTime.AddDate(0, 0, 9).Unix(), Open: 130, High: 135, Low: 128, Close: 133, Volume: 1000},
	}
}
