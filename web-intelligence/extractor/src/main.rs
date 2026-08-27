use serde::{Deserialize, Serialize};
use std::io::{self, Read, Write};
use std::time::Duration;

#[derive(Deserialize)]
struct Request {
    url: String,
    max_response_bytes: u64,
    timeout_ms: u64,
}

#[derive(Serialize)]
#[serde(tag = "status")]
enum Response {
    #[serde(rename = "ok")]
    Ok { document: Document },
    #[serde(rename = "error")]
    Error { error: ErrorBody },
}

#[derive(Serialize)]
struct Document {
    title: String,
    content: String,
    metadata: Metadata,
}
