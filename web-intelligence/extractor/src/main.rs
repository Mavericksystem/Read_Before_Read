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

#[derive(Serialize)]
struct Metadata {
    content_type: String,
    content_length_bytes: u64,
    fetch_duration_ms: u128,
}

#[derive(Serialize)]
struct ErrorBody {
    category: &'static str,
    message: String,
}

fn main() {
    let mut input = String::new();
    if let Err(e) = io::stdin().read_to_string(&mut input) {
        emit_error("internal", &format!("failed to read stdin: {e}"));
        std::process::exit(1);
    }

    let req: Request = match serde_json::from_str(&input) {
        Ok(r) => r,
        Err(e) => {
            emit_error("invalid_url", &format!("malformed request JSON: {e}"));
            std::process::exit(1);
        }
    };

    if !req.url.starts_with("http://") && !req.url.starts_with("https://") {
        emit_error("invalid_url", "url must start with http:// or https://");
        std::process::exit(1);
    }

    match fetch_and_extract(&req) {
        Ok(doc) => {
            let resp = Response::Ok { document: doc };
            print_json(&resp);
        }
        Err((category, message)) => {
            emit_error(category, &message);
            std::process::exit(1);
        }
    }
}

fn fetch_and_extract(req: &Request) -> Result<Document, (&'static str, String)> {
    let start = std::time::Instant::now();

    let client = reqwest::blocking::Client::builder()
        .timeout(Duration::from_millis(req.timeout_ms))
        .build()
        .map_err(|e| ("internal", e.to_string()))?;

    let resp = client
        .get(&req.url)
        .send()
        .map_err(|e|