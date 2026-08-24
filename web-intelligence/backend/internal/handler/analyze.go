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
