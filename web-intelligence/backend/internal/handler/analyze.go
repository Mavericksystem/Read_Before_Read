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
