package extractor

type Request struct {
	URL              string `json:"url"`
	MaxResponseBytes int64  `json:"max_response_bytes"`
	TimeoutMs        int64  `json:"timeout_ms"`
}

type Document struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	ContentType        string `json:"content_type"`
	ContentLengthBytes int64  `json:"content_length_bytes"`
	FetchDurationMs    int64  `json:"fetch_duration_ms"`
}
