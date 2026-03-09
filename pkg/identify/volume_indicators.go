package identify

import (
	"math"

	talib "github.com/markcheno/go-talib"
)

// VolumeIndicatorSeries stores volume-related indicators aligned by candle index.
// VolumeIndicatorSeries 保存与K线索引对齐的量价指标序列。
type VolumeIndicatorSeries struct {
	OBV []float64 // On-Balance Volume
	AD  []float64 // Accumulation/Distribution
	MFI []float64 // Money Flow Index
	CMF []float64 // Chaikin Money Flow
	VPT []float64 // Volume Price Trend
}

// ComputeVolumeIndicators computes the core volume indicators for downstream scoring.
// ComputeVolumeIndicators 为后续评分计算核心量价指标。
func ComputeVolumeIndicators(cs []CandlestickWrapper, mfiPeriod, cmfPeriod int) VolumeIndicatorSeries {
	n := len(cs)
	empty := VolumeIndicatorSeries{
		OBV: make([]float64, n),
		AD:  make([]float64, n),
		MFI: make([]float64, n),
		CMF: make([]float64, n),
		VPT: make([]float64, n),
	}
	if n == 0 {
		return empty
	}
	if mfiPeriod < 2 {
		mfiPeriod = 14
	}
	if cmfPeriod < 2 {
		cmfPeriod = 20
	}
	effectiveMFIPeriod := mfiPeriod
	if n <= effectiveMFIPeriod {
		effectiveMFIPeriod = n - 1
	}
	if effectiveMFIPeriod < 2 {
		effectiveMFIPeriod = 2
	}

	high := make([]float64, n)
	low := make([]float64, n)
	closep := make([]float64, n)
	volume := make([]float64, n)
	for i, c := range cs {
		high[i] = c.High
		low[i] = c.Low
		closep[i] = c.Close
		volume[i] = c.Volume
	}

	out := VolumeIndicatorSeries{
		OBV: talib.Obv(closep, volume),
		AD:  talib.Ad(high, low, closep, volume),
		MFI: talib.Mfi(high, low, closep, volume, effectiveMFIPeriod),
		CMF: computeCMF(high, low, closep, volume, cmfPeriod),
		VPT: computeVPT(closep, volume),
	}

	normalizeIndicatorNaN(out.MFI)
	normalizeIndicatorNaN(out.CMF)
	normalizeIndicatorNaN(out.VPT)
	return out
}

func computeCMF(high, low, closep, volume []float64, period int) []float64 {
	n := len(closep)
	out := make([]float64, n)
	for i := range out {
		out[i] = math.NaN()
	}
	if n == 0 {
		return out
	}

	for i := 0; i < n; i++ {
		start := i - period + 1
		if start < 0 {
			start = 0
		}
		sumMFV := 0.0
		sumVol := 0.0
		for j := start; j <= i; j++ {
			hlRange := high[j] - low[j]
			if hlRange == 0 {
				continue
			}
			mfm := ((closep[j] - low[j]) - (high[j] - closep[j])) / hlRange
			mfv := mfm * volume[j]
			sumMFV += mfv
			sumVol += volume[j]
		}
		if sumVol == 0 {
			out[i] = 0
			continue
		}
		out[i] = sumMFV / sumVol
	}
	return out
}

func computeVPT(closep, volume []float64) []float64 {
	n := len(closep)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	out[0] = 0
	for i := 1; i < n; i++ {
		prevClose := closep[i-1]
		if prevClose == 0 {
			out[i] = out[i-1]
			continue
		}
		out[i] = out[i-1] + volume[i]*((closep[i]-prevClose)/prevClose)
	}
	return out
}

func normalizeIndicatorNaN(series []float64) {
	for i, v := range series {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			series[i] = 0
		}
	}
}
