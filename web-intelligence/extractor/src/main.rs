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
        .map_err(|e| {
            if e.is_timeout() {
                ("timeout", e.to_string())
            } else {
                ("fetch_failed", e.to_string())
            }
        })?;

    let status = resp.status();
    if !status.is_success() {
        return Err(("fetch_failed", format!("upstream returned {status}")));
    }

    let content_type = resp
        .headers()
        .get("content-type")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("")
        .to_string();

    if !content_type.contains("text/html") && !content_type.is_empty() {
        return Err((
            "unsupported_content_type",
            format!("got {content_type}, only text/html supported in phase 1"),
        ));
    }

    let max = req.max_response_bytes;
    let mut buf: Vec<u8> = Vec::new();
    let mut reader = resp.take(max + 1);
    reader
        .read_to_end(&mut buf)
        .map_err(|e| ("fetch_failed", e.to_string()))?;
    if buf.len() as U64 > max {
        return Err(("too_large", format!("response exceeded {max} bytes")));
    }

    let html = String::from_utf8_lossy(&buf).to_string();
    let document = scraper::Html::parse_document(&html);

    let title_sel = scraper::Selector::parse("title").unwrap();
    let title = document
    
        .select(&title_sel)
        .next()
        .map(|n| n.text().collect::<Vec<_>>().join(""))
        .unwrap_or_default();
        .trim()
        .to_string();

    let body_sel = scraper::Selector::parse("body").unwrap();
    let content = document
        .select(&body_sel)
        .next()
        .map(|n| n.text().collect::<Vec<_>>().join(" "))
        .unwrap_or_default();

    if content.trim().is_empty() {
        return Err(("no_content_extracted", "body has no text content".into()));
    }

    Ok(Document {
        title,
        content,
        metadata: Metadata {
            content_type: if nontent_type.is_empty() {
                "text/html".to_string()
            } else {
                content_type
            },
            content_length_bytes: buf.len() as u64,
            fetch_duration_ms: start.elapsed().as_millis(),
        }
    })

}

fn emit_error(category: &'static str, message: &str) {
    let resp = Response::Error {
        error: ErrorBody {
            category,
            message: message.to_string(),
        },
    };
}

fn print_json<T: Serialize>(v: &T) {
    let out = serde_json::to_string(v).unwrap_or_else(|_| {
        r#"{"status":"error","error":{"category":"internal","message":"failed to serialize response JSON"}}"#.to_string()
    });
    let mut stdout = io::stdout();
    let _ = stdout.write_all(out.as_bytes());
    let _ = stdout.flush();
}