package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"web-intelligence/backend/internal/extractor"
	"web-intelligence/backend/internal/nim"
)

const requestDeadline = 30 * time.Second

const maxResponseBytes = 5 * 1024 * 1024

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
	Category  string `json:"category"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type Analyzer struct {
	NimClient *nim.Client
}

func (a *Analyzer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := newRequestID()

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "validation", "only POST is supported", requestID)
		return
	}

	var req analyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation", "malformed JSON body", requestID)
		return
	}
	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "validation", "url is required", requestID)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestDeadline)
	defer cancel()

	fetchStart := time.Now()
	doc, err := extractor.Run(ctx, extractor.Request{
		URL:              req.URL,
		MaxResponseBytes: maxResponseBytes,
		TimeoutMs:        extractor.DefaultTimeout.Milliseconds(),
	})
	fetchDuration := time.Since(fetchStart)

	if err != nil {
		if extErr, ok := err.(*extractor.Error); ok {
			writeError(w, statusFor(extErr.Category), extErr.Category, extErr.Message, requestID)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal", "extraction failed", requestID)
		return
	}

	nimStart := time.Now()
	answer, err := a.NimClient.Analyze(ctx, doc.Title, doc.Content, req.Question)
	nimDuration := time.Since(nimStart)

	if err != nil {
		log.Printf("nim error: %v", err)
		writeError(w, http.StatusBadGateway, "ai", "AI analysis failed", requestID)
		return
	}

	resp := successResponse{Status: "success"}
	resp.Result.Title = doc.Title
	resp.Result.NimAnswer = answer
	resp.Meta = meta{
		RequestID:       requestID,
		TotalDurationMs: time.Since(start).Milliseconds(),
		FetchDurationMs: fetchDuration.Milliseconds(),
		NimDurationMs:   nimDuration.Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func statusFor(category string) int {
	switch category {
	case "invalid_url":
		return http.StatusBadRequest
	case "timeout":
		return http.StatusGatewayTimeout
	case "too_large", "unsupported_content_type", "no_content_extracted":
		return http.StatusUnprocessableEntity
	case "fetch_failure":
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

func writeError(w http.ResponseWriter, status int, category, message, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Status: "error",
		Error: errBody{
			Category:  category,
			Message:   message,
			RequestID: requestID,
		},
	})
}

func newRequestID() string {
	return time.Now().Format("20060102150405.000000")
}
