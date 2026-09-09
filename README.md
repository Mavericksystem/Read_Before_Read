# URL Extraction & Summarization Engine

**Fetch any webpage, extract the real content, and get an AI summary — powered by a Go + Rust concurrency pipeline.**

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go)
![Rust](https://img.shields.io/badge/Rust-stable-DEA584?logo=rust)
![React](https://img.shields.io/badge/React-Vite-61DAFB?logo=react)
![Status](https://img.shields.io/badge/status-local--only%20demo-yellow)
![License](https://img.shields.io/badge/license-MIT-blue)

> 📸 *Screenshot/GIF of the UI goes here — paste one URL, click Analyze, get a summary.*

---

## What is this?

A backend systems-engineering demo: paste a URL, and it fetches the page, extracts the real content (stripping ads/nav/boilerplate), and uses an LLM (NVIDIA NIM) to summarize it or answer a question about it — all running locally, built to showcase concurrent process architecture, not to be a hosted product.

**Key features:**
- 🦀 **Rust extraction engine** — fast, safe HTML parsing with SSRF-hardened fetching
- 🐹 **Go orchestration layer** — a pooled worker architecture instead of naive spawn-per-request
- 🌐 **JS-page support** — falls back to a pooled headless Chrome instance for JS-rendered pages
- 🤖 **NVIDIA NIM integration** — summarize a page, or ask it a direct question
- ⚛️ **React + Vite frontend** — minimal UI to drive the whole pipeline

---

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Architecture](#architecture)
- [Known Limitations](#known-limitations)
- [Contributing](#contributing)
- [License & Credits](#license--credits)

---

## Getting Started

**Requirements:**
- Go 1.22+
- Rust (stable, via `rustup`)
- Node.js 18+
- An NVIDIA NIM API key ([build.nvidia.com](https://build.nvidia.com))

**Install & build:**
```bash
git clone https://github.com/yourname/url-extraction-engine.git
cd url-extraction-engine/web-intelligence

# Build the Rust extractor
cd extractor && cargo build --release
cp target/release/extractor ../backend/bin/extractor   # Windows: extractor.exe

# Set up the backend
cd ../backend
export NVIDIA_NIM_API_KEY=your_key_here
go get github.com/chromedp/chromedp
go run ./cmd/server

# Set up the frontend (new terminal)
cd ../frontend
npm install
npm run dev
```

Open `http://localhost:5173`.

---

## Usage

**Via the UI:** paste a URL, hit Analyze.

**Via the API directly:**
```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/article"}'
```

With a question:
```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/article", "question": "What is the main argument?"}'
```

**Response:**
```json
{
  "status": "success",
  "result": { "title": "...", "nim_answer": "..." },
  "meta": { "request_id": "...", "total_duration_ms": 1234 }
}
```

---

## Architecture

```mermaid
flowchart TD
    A[React + Vite Frontend] -->|POST /api/v1/analyze| B[Go API Server]
    B --> C{URL validation<br/>SSRF checks<br/>Admission gate}
    C -->|rejected| Z1[400 / 503 error]
    C -->|accepted| D[Rust Worker Pool]
    D -->|content found| F[Extracted Document]
    D -->|empty content| E[Headless Chrome Pool]
    E --> D2[Re-extract via Rust]
    D2 --> F
    F --> G[NVIDIA NIM Client<br/>throttled]
    G --> H[Summary / Answer]
    H --> A

    style A fill:#61DAFB,color:#000
    style D fill:#DEA584,color:#000
    style D2 fill:#DEA584,color:#000
    style E fill:#4285F4,color:#fff
    style G fill:#76B900,color:#fff
    style B fill:#00ADD8,color:#000
```

Go and Rust communicate over a hand-rolled line-delimited JSON protocol across OS pipes, using a pool of long-lived worker processes rather than spawning one per request. Full design rationale is in the project's architecture notes (available on request).

---

## Known Limitations

- Sites with aggressive bot protection (Cloudflare-style challenges) may return `403` and fail — advanced anti-bot bypass is intentionally out of scope.
- Very slow pages time out after 20s (by design).
- Pool sizes are currently reasonable defaults, not yet load-tested.
- No persistence, no auth, no hosting — this runs locally only.

---

## Contributing

This is a personal learning project and not currently accepting external contributions, but issues/suggestions are welcome.

---

## License & Credits

MIT License. Built with [reqwest](https://github.com/seanmonstar/reqwest), [scraper](https://github.com/causal-agent/scraper), [chromedp](https://github.com/chromedp/chromedp), and the [NVIDIA NIM](https://build.nvidia.com) API.
