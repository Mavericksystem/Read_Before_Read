use crate::url_validate::{self, ValidationError};
use std::io::Read;
use std::time::Duration;

const MAX_REDIRECTS: u8 = 5;

#[derive(Debug)]
pub enum FetchError {
    Validation(ValidationError),
    TooManyRedirects,
    Timeout,
    Network(String),
    TooLarge(u64),
    BadStatus(u16),
}

impl std::fmt::Display for FetchError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            FetchError::Validation(e) => write!(f, "{e}"),
            FetchError::TooManyRedirects => write!(f, "exceeded {MAX_REDIRECTS} redirects"),
            FetchError::Timeout => write!(f, "request timed out"),
            FetchError::Network(m) => write!(f, "network error: {m}"),
            FetchError::TooLarge(max) => write!(f, "response exceeded {max} bytes"),
            FetchError::BadStatus(s) => write!(f, "upstream returned status {s}"),
        }
    }
}

pub struct FetchResult {
    pub bytes: Vec<u8>,
    pub content_type: String,
    pub final_url: String,
}

pub fn fetch(url: &str, max_bytes: u64, timeout: Duration) -> Result<FetchResult, FetchError> {
    let client = reqwest::blocking::Client::builder()
        .timeout(timeout)
        .redirect(reqwest::redirect::Policy::none()) // we handle redirects ourselves
        .build()
        .map_err(|e| FetchError::Network(e.to_string()))?;

    let mut current_url = url.to_string();

    for _ in 0..=MAX_REDIRECTS {
        url_validate::validate(&current_url).map_err(FetchError::Validation)?;

        let resp = client.get(&current_url).send().map_err(|e| {
            if e.is_timeout() {
                FetchError::Timeout
            } else {
                FetchError::Network(e.to_string())
            }
        })?;

        let content_type = resp
            .headers()
            .get("content-type")
            .and_then(|v| v.to_str().ok())
            .unwrap_or("")
            .to_string();

        let final_url = current_url.clone();
        let bytes = read_capped(resp, max_bytes)?;

        return Ok(FetchResult {
            bytes,
            content_type,
            final_url,
        });
    }

    Err(FetchError::TooManyRedirects)
}

fn read_capped(resp: reqwest::blocking::Response, max_bytes: u64) -> Result<Vec<u8>, FetchError> {
    let mut reader = resp.take(max_bytes + 1);
    let mut buf = Vec::new();
    reader
        .read_to_end(&mut buf)
        .map_err(|e| FetchError::Network(e.to_string()))?;

    if buf.len() as u64 > max_bytes {
        return Err(FetchError::TooLarge(max_bytes));
    }

    Ok(buf)
}

fn resolve_redirect(base: &str, location: &str) -> String {
    match url::Url::parse(base).and_then(|b| b.join(location)) {
        Ok(joined) => joined.to_string(),
        Err(_) => location.to_string(), // best-effort fallback; will fail validate() next loop if malformed
    }
}
