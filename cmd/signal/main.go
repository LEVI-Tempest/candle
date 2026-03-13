package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/LEVI-Tempest/Candle/pkg/datasource"
	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
	"github.com/LEVI-Tempest/Candle/pkg/signal"
)

type candleInputEnvelope struct {
	Symbol string            `json:"symbol"`
	Source string            `json:"source"`
	Data   []*v1.Candlestick `json:"data"`
}

func main() {
	inputPath := flag.String("input", "", "Input JSON file path. Supports []Candlestick or {symbol,source,data}.")
	outputPath := flag.String("output", "", "Output JSON file path. Empty prints to stdout.")
	configPath := flag.String("config", "", "Signal config JSON path.")
	logCSVPath := flag.String("log-csv", "", "Override signal log CSV path.")
	validateSchema := flag.Bool("validate-schema", false, "Validate output against JSON schema before printing.")
	schemaPath := flag.String("schema", "docs/signal.schema.json", "Path to JSON schema used with --validate-schema.")

	asOf := flag.String("as-of", time.Now().Format(time.RFC3339), "As-of timestamp (RFC3339).")
	symbol := flag.String("symbol", "XSHE:300059", "Symbol for reporting context.")

	fetch := flag.Bool("fetch", false, "Fetch data from Tsanghi instead of --input.")
	exchange := flag.String("exchange", datasource.TsanghiXSHE, "Exchange: XSHE | XSHG | XHKG")
	ticker := flag.String("ticker", "300059", "Ticker code")
	token := flag.String("token", "demo", "Tsanghi API token")
	limit := flag.Int("limit", 120, "Number of candles to fetch")
	flag.Parse()

	cfg, err := signal.LoadConfig(*configPath)
	if err != nil {
		exitf("load config failed: %v", err)
	}
	if *logCSVPath != "" {
		cfg.LogCSVPath = *logCSVPath
	}

	candles, source, detectedSymbol, err := loadCandles(*inputPath, *fetch, *exchange, *ticker, *token, *limit)
	if err != nil {
		exitf("load candles failed: %v", err)
	}
	if detectedSymbol != "" {
		*symbol = detectedSymbol
	}
	if len(candles) == 0 {
		exitf("no candles available")
	}

	report := signal.BuildReport(*symbol, *asOf, source, candles, cfg)
	if *validateSchema {
		if err := signal.ValidateReportSchema(report, *schemaPath); err != nil {
			exitf("schema validation failed: %v", err)
		}
	}
	if err := signal.AppendSignalLogCSV(cfg.LogCSVPath, report); err != nil {
		exitf("append signal log failed: %v", err)
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		exitf("marshal report failed: %v", err)
	}
	if *outputPath == "" {
		fmt.Println(string(data))
		return
	}
	if err := os.WriteFile(*outputPath, data, 0o644); err != nil {
		exitf("write output failed: %v", err)
	}
}

func loadCandles(
	inputPath string,
	fetch bool,
	exchange, ticker, token string,
	limit int,
) ([]*v1.Candlestick, string, string, error) {
	if fetch {
		client := datasource.NewTsanghiClient(token)
		candles, err := client.Fetch(exchange, ticker, &datasource.FetchOptions{
			Limit: limit,
			Order: 2,
		})
		if err != nil {
			return nil, "", "", err
		}
		return candles, "tsanghi", fmt.Sprintf("%s:%s", exchange, ticker), nil
	}
	if inputPath == "" {
		return nil, "", "", fmt.Errorf("either --input or --fetch is required")
	}

	raw, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, "", "", err
	}

	var envelope candleInputEnvelope
	if err := json.Unmarshal(raw, &envelope); err == nil && len(envelope.Data) > 0 {
		src := envelope.Source
		if src == "" {
			src = "file"
		}
		return envelope.Data, src, envelope.Symbol, nil
	}

	var candles []*v1.Candlestick
	if err := json.Unmarshal(raw, &candles); err != nil {
		return nil, "", "", fmt.Errorf("input must be []Candlestick or {symbol,source,data}: %w", err)
	}
	return candles, "file", "", nil
}

func exitf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "signal: "+format+"\n", args...)
	os.Exit(1)
}
