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
