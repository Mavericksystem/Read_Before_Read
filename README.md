# Current Status

🚧 **Under active development**

The project is currently being built from the ground up, starting with the Go and Rust application entry point and gradually introducing each architectural layer.


---

## Web Intelligence Platform 

> **Ask a webpage a question before you spend time reading it.**

An AI-powered web intelligence platform that takes a URL, extracts its meaningful content, and uses **NVIDIA Nemotron** to determine whether the page contains the information you're looking for.

The project is being built from the ground up using **Go, Rust, and TypeScript**, with a focus on production-oriented backend engineering.

## 🎯 What I'm Building

Paste a URL and optionally ask:

> "Does this page contain information about the company's revenue?"

The system returns:

- **Yes / No / Unclear**
- **Relevance score**
- Explanation based on the page
- Key points
- Important facts
- Entities and statistics
- Why the information matters

### Architecture

```text
React + TypeScript
        │
        ▼
     Go API
        │
        ▼
   Rust Extractor
        │
        ▼
  Clean Web Content
        │
        ▼
 NVIDIA Nemotron
        │
        ▼
Question + Relevance + Explanation
```

## 🛠️ Technology Stack

**Frontend**
- TypeScript
- React
- Framer Motion

**Backend**
- Go
- REST API
- JSON
- Concurrency

**Web Processing**
- Rust
- HTTP fetching
- HTML parsing
- Content extraction

**AI**
- NVIDIA NIM
- Nemotron 3 Ultra

**Deployment**
- Render

## 🔐 Privacy by Design

The application is intentionally **stateless**.

It does not store:

- URLs
- Webpage content
- User questions
- AI summaries
- User accounts
- History

Web content is processed temporarily and returned to the user.

## 🧠 Why Go + Rust?

This project is also an exploration of using two systems languages for different responsibilities:

**Go** → API, orchestration, concurrency, external services

**Rust** → high-performance webpage fetching and content extraction

The goal is to understand how these responsibilities can be separated in a real production-oriented system.

---

## 📚 Engineering Focus

This project is being used to explore:

- Backend architecture
- Go concurrency
- Rust systems programming
- API design
- Web content extraction
- AI inference integration
- Security and SSRF protection
- Testing and QA
- Performance engineering
- Observability
- Production deployment

> **Build → understand → measure → improve.**

## Current Folder Structure

```text
rust_pro/
├── .env
├── .gitignore
├── README.md
├── fake.go
├── target/                         # Workspace build output
└── web-intelligence/
      ├── backend/
      │   ├── bin/                    # Compiled Rust extractor binary
      │   ├── cmd/
      │   │   └── server/
      │   │       └── main.go         # Go HTTP server entry point
      │   ├── internal/
      │   │   ├── extractor/
      │   │   │   └── extractor.go    # Rust process orchestration
      │   │   ├── handler/
      │   │   │   └── analyze.go      # Analyze API handler
      │   │   └── nim/
      │   │       └── nim.go          # NVIDIA NIM client
      │   ├── go.mod
      │   └── package-lock.json
      ├── docs/
      │   └── .gitkeep
      ├── extractor/
      │   ├── benches/
      │   │   └── extraction_bench.rs
      │   ├── src/
      │   │   ├── encoding.rs
      │   │   ├── extract.rs
      │   │   ├── fetch.rs
      │   │   ├── lib.rs
      │   │   ├── main.rs
      │   │   └── url_validate.rs
      │   ├── tests/
      │   │   └── extract_test.rs
      │   ├── Cargo.toml
      │   └── Cargo.lock
      ├── frontend/
      │   ├── src/
      │   │   ├── App.tsx
      │   │   └── main.tsx
      │   ├── index.html
      │   ├── package.json
      │   ├── package-lock.json
      │   ├── tsconfig.json
      │   ├── tsconfig.app.json
      │   ├── tsconfig.node.json
      │   └── vite.config.ts
      ├── infrastructure/
      │   └── .gitkeep
      └── tests/
            └── .gitkeep
```

Generated directories such as `frontend/node_modules`, `frontend/dist`, and
`extractor/target` are omitted from the detailed source tree.
