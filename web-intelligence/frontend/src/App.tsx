import { useState } from "react";

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
        <div style={{ fontFamily: "monospace", padding: "2rem", maxWidth: 700 }}>
            <h1>Web Intelligence — Phase 1 vertical slice</h1>
            <form onSubmit={handleSubmit}>
                <input
                    type="url"
                    required
                    placeholder="https://example.com/article"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    style={{ width: "100%", padding: "0.5rem" }}
                />
                <button type="submit" disabled={loading} style={{ marginTop: "0.5rem" }}>
                    {loading ? "Analyzing..." : "Analyze"}
                </button>
            </form>