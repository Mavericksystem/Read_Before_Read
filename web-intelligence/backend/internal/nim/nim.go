package nim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	endpoint = "https://integrate.api.nvidia.com/v1/chat/completions",
	model	 = "nvidia/nemotron-3-ultra-550b-a55b"
)

type Client struct {
	apiKey	   string
	httpClient *http.Client
}

func NewClient() (*Client, error) {
	key := os.Getenv("NVIDIA_NIM_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("NVIDIA_NIM_API_KEY not set")
	}
	
}