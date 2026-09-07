import { useState } from "react";
import "./App.css";

interface AnalyzeSuccess {
    status: "success";
    result: { title: string; nim_answer: string };
    meta: { request_id: string; total_duration_ms: number };
}

interface AnalyzeError {
    status: "error";
    error: { category: string; message: string; request_id: string };
}

type AnalyzeResponse = AnalyzeSuccess | AnalyzeError;

const API_BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8080";

export default function App() {
    const [url, setUrl] = useState("");
    const [loading, setLoading] = useState(false);
    const [response, setResponse] = useState<AnalyzeResponse | null>(null);

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault();
        setLoading(true);
        setResponse(null);
        try {
            const res = await fetch(`${API_BASE}/api/v1/analyze`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ url }),
            });
            const data: AnalyzeResponse = await res.json();
            setResponse(data);
        } catch (err) {
            setResponse({
                status: "error",
                error: {
                    category: "internal",
                    message: err instanceof Error ? err.message : "unknown network error",
                    request_id: "",
                },
            });
        } finally {
            setLoading(false);
        }
    }

    return (
        <main className="app-shell">
            <nav className="topbar">
                <a className="brand" href="/" aria-label="Web Intelligence home">
                    <span className="brand-mark">WI</span>
                    <span>Web Intelligence</span>
                </a>
                <span className="status-pill"><span className="status-dot" /> AI research assistant</span>
            </nav>

            <section className="hero">
                <p className="eyebrow">Read less. Understand more.</p>
                <h1>Ask any webpage<br /><em>a better question.</em></h1>
                <p className="hero-copy">Paste an article URL and get a clear, useful answer without digging through every paragraph.</p>
            </section>

            <section className="workspace" aria-label="Webpage analyzer">
                <form className="analyzer-form" onSubmit={handleSubmit}>
                    <label htmlFor="url">Page to investigate</label>
                    <div className="input-row">
                        <div className="url-input-wrap">
                            <span className="url-icon" aria-hidden="true">↗</span>
                            <input
                                id="url"
                                type="url"
                                required
                                placeholder="https://example.com/article"
                                value={url}
                                onChange={(e) => setUrl(e.target.value)}
                            />
                        </div>
                        <button type="submit" disabled={loading}>
                            {loading ? <><span className="spinner" /> Reading</> : <>Analyze <span aria-hidden="true">→</span></>}
                        </button>
                    </div>
                    <p className="form-hint">We extract the page content temporarily and do not store your URL.</p>
                </form>

                {loading && <div className="loading-card"><span className="loader-line" /><strong>Reading the page</strong><span>Extracting the useful bits and preparing your answer...</span></div>}

                {response?.status === "success" && (
                    <article className="result-card">
                        <div className="result-header"><span className="result-label"><span className="check">✓</span> Analysis complete</span><span>{response.meta.total_duration_ms} ms</span></div>
                        <h2>{response.result.title || "Untitled page"}</h2>
                        <p className="answer">{response.result.nim_answer}</p>
                        <div className="result-footer"><span>Powered by Nemotron</span><span>Request {response.meta.request_id.slice(0, 8)}</span></div>
                    </article>
                )}

                {response?.status === "error" && (
                    <div className="error-card" role="alert"><strong>We couldn't analyze that page.</strong><span>{response.error.message}</span></div>
                )}
            </section>

            <footer><span>WEB INTELLIGENCE / 2026</span><span>Turn pages into perspective.</span></footer>
        </main>
    );
}