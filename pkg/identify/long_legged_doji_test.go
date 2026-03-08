package identify

import (
	"testing"

	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

func TestLongLeggedDoji(t *testing.T) {
	// 长腿十字星样本：实体很小，双影线都很长 | Long-legged doji sample: tiny body with long upper/lower shadows
	c := &v1.Candlestick{
		Open:  100.0,
		Close: 100.3,
		High:  112.0,
		Low:   88.0,
	}
	w := []CandlestickWrapper{NewCandlestickWrapper(c)}
	if !LongLeggedDoji(w) {
		t.Fatalf("expected LongLeggedDoji to be detected, but got false")
	}

	// 反例：实体过大，不应识别 | Negative sample: body too large, should not be detected
	c2 := &v1.Candlestick{
		Open:  100.0,
		Close: 104.0,
		High:  112.0,
		Low:   88.0,
	}
	w2 := []CandlestickWrapper{NewCandlestickWrapper(c2)}
	if LongLeggedDoji(w2) {
		t.Fatalf("expected LongLeggedDoji to be false for large body candle")
	}
}
