package identify

import (
	"fmt"

	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gonum.org/v1/gonum/stat"
)

/***
 * @author Tempest
 * @description Long-term trend identification (识别长期趋势)
 * @date 2025/05/20
 */

// Trend represents the direction of a price movement.
// 表示价格走势的方向。
type Trend string

const (
	TrendYang    Trend = "Yang"    // Upward trend (上升趋势)
	TrendYin     Trend = "Yin"     // Downward trend (下降趋势)
	TrendMiddle  Trend = "Middle"  // Sideways/oscillating trend (震荡持平)
	TrendUnknown Trend = "Unknown" // Unknown trend (未知趋势)
)

// DetermineLongTermTrend determines the long-term trend within a given time period using linear regression.
// Parameters:
//   - candles: Historical candlestick data, arranged from oldest to newest. (历史蜡烛图数据，从旧到新排列)
//   - days: Number of days used to determine the trend. (用于判断趋势的天数)
//
// Returns:
//   - Trend: The direction of the trend (Yang, Yin, or Middle). (趋势方向)
//   - error: An error if the candlestick data is insufficient or the number of days is invalid. (如果数据不足或天数无效时返回错误)
func DetermineLongTermTrend(candles []*v1.Candlestick, days int) (Trend, error) {
	if len(candles) < days {
		return TrendUnknown, fmt.Errorf("Need at least %d days of data to determine trend (需要至少 %d 天的数据)", days, days)
	}
	if days <= 0 {
		return TrendUnknown, fmt.Errorf("Number of days must be greater than 0 (天数必须大于 0)")
	}

	// Get closing prices and relative time index for the specified number of days.
	// 获取指定天数内的收盘价和相对时间索引。
	prices := make([]float64, 0, days)
	timestamps := make([]float64, 0, days)
	for i := len(candles) - days; i < len(candles); i++ {
		prices = append(prices, candles[i].Close)
		// Use relative index to avoid numerical instability from large Unix timestamps.
		// 使用相对索引，避免 Unix 时间戳过大导致的线性回归数值不稳定。
		timestamps = append(timestamps, float64(i-(len(candles)-days)))
	}

	// Use Gota and Gonum for linear regression
	// 使用 Gota 和 Gonum 进行线性回归
	priceSeries := series.New(prices, series.Float, "Price")
	timeSeries := series.New(timestamps, series.Float, "Time")
	df := dataframe.New(priceSeries, timeSeries)

	// Calculate linear regression
	// 计算线性回归
	slope, _, rSquared, err := LinearRegression(df, "Time", "Price")
	if err != nil {
		return TrendUnknown, fmt.Errorf("Linear regression failed: %v (线性回归失败)", err)
	}

	// Determine the trend based on slope and R-squared values
	// 根据斜率和R平方值确定趋势

	// Slope threshold: determines if the rate of price change is significant enough to be considered an upward or downward trend
	// 斜率阈值：决定价格变化的速率是否足够显著来判断为上升或下降趋势
	const thresholdSlope = 1e-5 // Adjustable slope threshold

	// R-squared threshold: determines the goodness of fit of the linear model, higher values indicate more pronounced trends
	// R平方阈值：决定线性模型的拟合程度，值越高表示趋势越明显
	const thresholdRSquared = 0.5 // Adjustable R-squared threshold

	// If R-squared is high enough, the linear relationship is significant
	// 如果R平方值足够高，表示线性关系显著
	if rSquared > thresholdRSquared {
		// Determine trend direction based on slope
		// 根据斜率判断趋势方向
		if slope > thresholdSlope {
			return TrendYang, nil // Upward trend (上升趋势)
		} else if slope < -thresholdSlope {
			return TrendYin, nil // Downward trend (下降趋势)
		} else {
			return TrendMiddle, nil // Sideways trend (震荡持平)
		}
	} else {
		// If R-squared is too low, the linear relationship is not significant, considered a sideways market
		// 如果R平方值太低，表示线性关系不显著，认为是震荡市场
		return TrendMiddle, nil
	}
}

// LinearRegression performs linear regression on the given dataframe and calculates R-squared.
// Parameters:
//   - df: The input dataframe containing the data for regression. (输入的数据帧)
//   - xCol: The name of the column to be used as the independent variable. (自变量列名)
//   - yCol: The name of the column to be used as the dependent variable. (因变量列名)
//
// Returns:
//   - slope: The slope of the regression line. (回归线的斜率)
//   - intercept: The y-intercept of the regression line. (回归线的截距)
//   - rSquared: The coefficient of determination (R-squared). (R平方值，决定系数)
//   - error: An error if the regression fails. (回归失败时的错误)
func LinearRegression(df dataframe.DataFrame, xCol, yCol string) (float64, float64, float64, error) {
	xSeries := df.Col(xCol)
	ySeries := df.Col(yCol)
	if xSeries.Len() != ySeries.Len() {
		return 0, 0, 0, fmt.Errorf("Input sequence lengths do not match (输入序列长度不匹配)")
	}

	x := xSeries.Float()
	y := ySeries.Float()

	if len(x) == 0 {
		return 0, 0, 0, fmt.Errorf("Input sequence is empty (输入序列为空)")
	}

	// Calculate weights (using nil means all points have equal weight)
	// 计算权重（这里使用nil表示所有点权重相等）
	var weights []float64

	// Gonum returns intercept first, then slope.
	// Gonum 返回值顺序是截距在前、斜率在后。
	intercept, slope := stat.LinearRegression(x, y, weights, false)

	// Manually calculate R-squared value
	// R-squared = 1 - (residual sum of squares / total sum of squares)
	// 手动计算R平方值
	yMean := stat.Mean(y, weights)
	var totalSumSquares, residualSumSquares float64

	for i := range x {
		// Predicted value
		// 预测值
		predicted := slope*x[i] + intercept
		// Residual squared
		// 残差平方
		residualSumSquares += (y[i] - predicted) * (y[i] - predicted)
		// Total deviation squared
		// 总偏差平方
		totalSumSquares += (y[i] - yMean) * (y[i] - yMean)
	}

	// Calculate final R-squared value
	// 计算最终R平方值
	rSquared := 1.0
	if totalSumSquares > 0 {
		rSquared = 1.0 - (residualSumSquares / totalSumSquares)
	}

	return slope, intercept, rSquared, nil
}
