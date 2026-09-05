mod encoding;
mod extract;
mod fetch;
mod url_validate;

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
    final_url: String,
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

    match run(&req) {
        Ok(doc) => print_json(&Response::Ok { document: doc }),
        Err((category, message)) => {
            emit_error(category, &message);
            std::process::exit(1);
        }
    }
}

fn run(req: &Request) -> Result<Document, (&'static str, String)> {
    let start = std::time::Instant::now();

    url_validate::validate(&req.url).map_err(|e| ("invalid_url", e.to_string()))?;

    let fetch_result = fetch::fetch(
        &req.url,
        req.max_response_bytes,
        Duration::from_millis(req.timeout_ms),
    )
    .map_err(map_fetch_error)?;

    if !fetch_result.content_type.contains("text/html") && !fetch_result.content_type.is_empty() {
        return Err((
            "unsupported_content_type",
            format!(
                "got {}, only text/html supported",
                fetch_result.content_type
            ),
        ));
    }

    let html = encoding::decode(&fetch_result.bytes, &fetch_result.content_type);

    let extracted = extract::extract(&html);

    if extracted.content.trim().is_empty() {
        return Err((
            "no_content_extracted",
            "no text content found — page may require JavaScript (out of scope, ADR-009)".into(),
        ));
    }

    Ok(Document {
        title: extracted.title,
        content: extracted.content,
        metadata: Metadata {
            content_type: if fetch_result.content_type.is_empty() {
                "text/html".to_string()
            } else {
                fetch_result.content_type
            },
            content_length_bytes: fetch_result.bytes.len() as u64,
            fetch_duration_ms: start.elapsed().as_millis(),
            final_url: fetch_result.final_url,
        },
    })
}

fn map_fetch_error(e: fetch::FetchError) -> (&'static str, String) {
    use fetch::FetchError;
    let msg = e.to_string();
    let category = match e {
        FetchError::Validation(_) => "invalid_url",
        FetchError::TooManyRedirects => "fetch_failed",
        FetchError::Timeout => "timeout",
        FetchError::Network(_) => "fetch_failed",
        FetchError::TooLarge(_) => "too_large",
        FetchError::BadStatus(_) => "fetch_failed",
    };
    (category, msg)
}

fn emit_error(category: &'static str, message: &str) {
    print_json(&Response::Error {
        error: ErrorBody {
            category,
            message: message.to_string(),
        },
    });
}

fn print_json<T: Serialize>(v: &T) {
    let out = serde_json::to_string(v).unwrap_or_else(|_| {
        r#"{"status":"error","error":{"category":"internal","message":"failed to serialize response"}}"#
            .to_string()
    });
    let mut stdout = io::stdout();
    let _ = stdout.write_all(out.as_bytes());
    let _ = stdout.flush();
}
