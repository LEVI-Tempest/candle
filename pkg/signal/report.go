package signal

import "github.com/LEVI-Tempest/Candle/pkg/identify"

// PatternReport is a user-facing pattern item in signal output.
// PatternReport 是信号输出里的用户可读形态条目。
type PatternReport struct {
	Type          string   `json:"type"`
	Direction     string   `json:"direction"`
	Position      int      `json:"position"`
	Strength      float64  `json:"strength"`
	Risk          float64  `json:"risk"`
	Price         float64  `json:"price"`
	Time          string   `json:"time"`
	VolumeState   string   `json:"volume_state"`
	DecisionScore float64  `json:"decision_score"`
	DecisionLevel string   `json:"decision_level"`
	Reason        []string `json:"reason"`
	ForwardRet3   *float64 `json:"forward_ret_3,omitempty"`
	ForwardRet5   *float64 `json:"forward_ret_5,omitempty"`
	ForwardRet10  *float64 `json:"forward_ret_10,omitempty"`
}

// Report is the structured signal payload for upper-layer agents.
// Report 是给上层 Agent 使用的结构化信号载荷。
type Report struct {
	Symbol          string                     `json:"symbol"`
	AsOf            string                     `json:"as_of"`
	Source          string                     `json:"source"`
	Trend           string                     `json:"trend"`
	Score           float64                    `json:"score"` // normalized 0-1
	DecisionScore   float64                    `json:"decision_score"`
	DecisionLevel   string                     `json:"decision_level"`
	Patterns        []PatternReport            `json:"patterns"`
	Evidence        []identify.PatternEvidence `json:"evidence"`
	CounterEvidence []string                   `json:"counter_evidence"`
	InvalidIf       []string                   `json:"invalid_if"`
}
