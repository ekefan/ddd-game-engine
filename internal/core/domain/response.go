package domain

type Response struct {
	Round        int    `json:"round"`
	Winner       string `json:"winner"`
	RoundOutcome string`json:"round_outcome"`
	SessionEnded bool   `json:"session_ended"`
}
