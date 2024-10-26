package domain

type Response struct {
	Round        int    `json:"round"`
	Winner       string `json:"winner"`
	SessionEnded bool   `json:"session_ended"`
}
