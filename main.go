// Candle - Japanese Candlestick Charting
// 日本蜡烛图技术 - 形态识别与可视化
//
// Ref: Japanese Candlestick Charting Techniques
// https://book.douban.com/subject/2124790/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/LEVI-Tempest/Candle/pkg/charting"
	"github.com/LEVI-Tempest/Candle/pkg/datasource"
	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

func main() {
	// CLI flags | 命令行参数
	example := flag.String("example", "chart", "Demo: chart | fetch")
	output := flag.String("output", "candle_chart.html", "Output HTML filename")
	// fetch 专用
	exchange := flag.String("exchange", "XSHE", "Exchange: XSHE(深圳) | XSHG(上海)")
	ticker := flag.String("ticker", "300059", "Stock ticker code")
	token := flag.String("token", "demo", "Tsanghi API token")
	limit := flag.Int("limit", 60, "Number of days to fetch")
	flag.Parse()

	switch *example {
	case "chart":
		runChartDemo(*output)
	case "fetch":
		runFetchDemo(*output, *exchange, *ticker, *token, *limit)
	default:
		fmt.Fprintf(os.Stderr, "Unknown example: %s. Use: chart | fetch\n", *example)
		os.Exit(1)
	}
}

// runFetchDemo fetches data from Tsanghi API and generates chart
// runFetchDemo 从 Tsanghi API 拉取数据并生成图表
func runFetchDemo(outputFile, exchange, ticker, token string, limit int) {
	fmt.Println("🕯️  Candle - Fetch & Chart Demo")
	fmt.Println("=================================")
	fmt.Printf("Fetching %s %s (%d days)...\n", exchange, ticker, limit)

	client := datasource.NewTsanghiClient(token)
	candles, err := client.Fetch(exchange, ticker, &datasource.FetchOptions{
		Limit: limit,
		Order: 2,
	})
	if err != nil {
		log.Fatalf("❌ Fetch failed: %v", err)
	}
	if len(candles) == 0 {
		log.Fatal("❌ No data returned")
	}

	fmt.Printf("📊 Fetched %d candlesticks\n", len(candles))

	ek := charting.NewEnhancedKline()
	ek.LoadData(candles)
	ek.AutoDetectPatterns()
	fmt.Printf("🔍 Detected %d patterns\n\n", len(ek.Patterns))

	ek.CreateChart(fmt.Sprintf("🕯️ %s %s - Candlestick Chart", exchange, ticker))
	if err := ek.RenderToFile(outputFile); err != nil {
		log.Fatalf("❌ Render failed: %v", err)
	}
	fmt.Printf("✅ Chart saved: %s\n", outputFile)
}

// runChartDemo creates a candlestick chart with pattern markers
// 运行图表示例：生成带形态标记的蜡烛图
func runChartDemo(outputFile string) {
	fmt.Println("🕯️  Candle - Japanese Candlestick Chart Demo")
	fmt.Println("=============================================")
	fmt.Println()

	// 1. Create sample data with identifiable patterns | 创建包含可识别形态的示例数据
	candleData := createDemoData()
	fmt.Printf("📊 Loaded %d candlesticks\n", len(candleData))

	// 2. Create and load enhanced kline chart | 创建并加载增强K线图
	ek := charting.NewEnhancedKline()
	ek.LoadData(candleData)

	// 3. Auto-detect patterns | 自动检测形态
	ek.AutoDetectPatterns()
	fmt.Printf("🔍 Detected %d patterns\n\n", len(ek.Patterns))

	// 4. Print detected patterns | 输出检测到的形态
	if len(ek.Patterns) > 0 {
		fmt.Println("📍 Detected Patterns:")
		fmt.Println("--------------------")
		for i, p := range ek.Patterns {
			fmt.Printf("%d. %s | Pos:%d | Price:%.2f | Strength:%.1f\n",
				i+1, p.Type, p.Position, p.Price, p.Strength)
		}
		fmt.Println()
	}

	// 5. Create chart and render to HTML | 创建图表并渲染到 HTML
	ek.CreateChart("🕯️ Candlestick Chart with Pattern Markers")
	if err := ek.RenderToFile(outputFile); err != nil {
		log.Fatalf("❌ Failed to render chart: %v", err)
	}

	fmt.Printf("✅ Chart saved: %s\n", outputFile)
	fmt.Println()
	fmt.Println("📖 Usage: Open the HTML file in your browser to view the interactive chart.")
}

// createDemoData returns sample candlestick data with clear patterns
// 创建包含明确形态的示例蜡烛数据
func createDemoData() []*v1.Candlestick {
	baseTime := time.Now().AddDate(0, 0, -25)

	return []*v1.Candlestick{
		{Timestamp: baseTime.Unix(), Open: 100, High: 105, Low: 98, Close: 103, Volume: 1000},
		{Timestamp: baseTime.AddDate(0, 0, 1).Unix(), Open: 103, High: 108, Low: 101, Close: 106, Volume: 1200},
		{Timestamp: baseTime.AddDate(0, 0, 2).Unix(), Open: 106, High: 110, Low: 104, Close: 109, Volume: 1100},
		// Doji
		{Timestamp: baseTime.AddDate(0, 0, 3).Unix(), Open: 109, High: 114, Low: 104, Close: 109.3, Volume: 800},
		{Timestamp: baseTime.AddDate(0, 0, 4).Unix(), Open: 109, High: 111, Low: 102, Close: 104, Volume: 1300},
		// Hammer
		{Timestamp: baseTime.AddDate(0, 0, 5).Unix(), Open: 104, High: 106, Low: 92, Close: 105, Volume: 1800},
		// Marubozu
		{Timestamp: baseTime.AddDate(0, 0, 6).Unix(), Open: 105, High: 122, Low: 105, Close: 122, Volume: 2500},
		{Timestamp: baseTime.AddDate(0, 0, 7).Unix(), Open: 122, High: 128, Low: 120, Close: 126, Volume: 1600},
		// Shooting Star
		{Timestamp: baseTime.AddDate(0, 0, 8).Unix(), Open: 126, High: 142, Low: 125, Close: 128, Volume: 2000},
		// Bearish Engulfing
		{Timestamp: baseTime.AddDate(0, 0, 9).Unix(), Open: 128, High: 132, Low: 127, Close: 131, Volume: 1400},
		{Timestamp: baseTime.AddDate(0, 0, 10).Unix(), Open: 133, High: 134, Low: 120, Close: 122, Volume: 2800},
		{Timestamp: baseTime.AddDate(0, 0, 11).Unix(), Open: 122, High: 124, Low: 115, Close: 117, Volume: 1900},
		{Timestamp: baseTime.AddDate(0, 0, 12).Unix(), Open: 117, High: 119, Low: 110, Close: 112, Volume: 1700},
		// Morning Star
		{Timestamp: baseTime.AddDate(0, 0, 13).Unix(), Open: 112, High: 113, Low: 105, Close: 107, Volume: 1800},
		{Timestamp: baseTime.AddDate(0, 0, 14).Unix(), Open: 104, High: 105, Low: 102, Close: 103, Volume: 900},
		{Timestamp: baseTime.AddDate(0, 0, 15).Unix(), Open: 105, High: 118, Low: 104, Close: 116, Volume: 2400},
		// Three White Soldiers
		{Timestamp: baseTime.AddDate(0, 0, 16).Unix(), Open: 116, High: 124, Low: 115, Close: 123, Volume: 2000},
		{Timestamp: baseTime.AddDate(0, 0, 17).Unix(), Open: 121, High: 130, Low: 120, Close: 129, Volume: 2100},
		{Timestamp: baseTime.AddDate(0, 0, 18).Unix(), Open: 127, High: 136, Low: 126, Close: 135, Volume: 2200},
		{Timestamp: baseTime.AddDate(0, 0, 19).Unix(), Open: 135, High: 142, Low: 133, Close: 140, Volume: 1800},
		// Spinning Top
		{Timestamp: baseTime.AddDate(0, 0, 20).Unix(), Open: 140, High: 147, Low: 133, Close: 141, Volume: 1500},
		// Evening Star
		{Timestamp: baseTime.AddDate(0, 0, 21).Unix(), Open: 141, High: 148, Low: 140, Close: 147, Volume: 1700},
		{Timestamp: baseTime.AddDate(0, 0, 22).Unix(), Open: 149, High: 150, Low: 148, Close: 149.2, Volume: 800},
		{Timestamp: baseTime.AddDate(0, 0, 23).Unix(), Open: 148, High: 149, Low: 138, Close: 141, Volume: 2600},
		{Timestamp: baseTime.AddDate(0, 0, 24).Unix(), Open: 141, High: 145, Low: 138, Close: 142, Volume: 1400},
	}
}
