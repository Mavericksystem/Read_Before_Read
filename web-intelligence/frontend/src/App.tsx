import { useState } from "react";

interface AnalyzeSuccess {
    status: "success";
    result: { title: string; nim_answer: string };
    meta: { request_id: string; total_duration_ms: number };
}