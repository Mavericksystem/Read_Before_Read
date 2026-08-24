package handler

import (
	"time"
)

const requesDeadline = 30 * time.Second

const maxRespnseBytes = 5 * 1024 * 1024

type analyzeRequest struct {
	URL      string `json:"url"`
	Question string `json:"question,omitempty"`
}

type successResponse struct {
	Status string `json:"status"`
	Result struct {
		Title     string `json:"title"`
		NimAnswer string `json:"nim_answer"`
	} `json:"result"`
	Meta meta `json:"meta"`
}

type meta struct {
	RequestID       string `json:"request_id"`
	TotalDurationMs int64  `json:"total_duration_ms"`
	FetchDurationMs int64  `json:"fetch_duration_ms"`
	NimDurationMs   int64  `json:"nim_duration_ms"`
}

type errorResponse struct {
	Status string  `json:"status"`
	Error  errBody `json:"error"`
}

type errBody struct {
	NimClient *nim.Client
}
