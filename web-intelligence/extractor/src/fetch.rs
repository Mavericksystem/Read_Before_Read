use crate::url_validate::{self, ValidationError};
use std::io::Read;
use std::time::Duration;

const MAX_REDIRECTS: u8 = 5;

#[derive(Debug)]
pub enum FetchError {
    Validation(ValiddationError),
    TooManyRedirects,
    Timeout,
    Network(reqwest::Error),
    TooLarge(u64),
    BadStatus(u16),
}

impl std::fmt::Display for FetchError {
    fm fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            FetchError::Validation(e) => write!(f, "{e}"),
            FetchError::TooManyRedirects => write!(f, "Exceeds {Max_REDIRECTS} redirects"),
            FetchError::Timeout => write!(f, "Request timed out"),
            FetchError::Network(e) => write!(f, "Network error: {e}"),
            FetchError::TooLarge(size) => write!(f, "Response too large: {size} bytes"),
            FetchError::BadStatus(status) => write!(f, "Bad status code: {status}"),
        }
    }
}

pub struct FetchResult {
    pub html: string,
    pub content_type: String,
    pub final_url: String,
}

pub fn fetch(
    url: &str,
    max_bytes: u64,
    timeout: Duration,
) -> Result<FetchResult, FetchError> {
    let client = reqwest::blocking::Client::builder()
        .timeout(timeout)
        .redirect(reqwest::redirect::Policy::none())
        .build()
        .map_err(FetchError::Network(e.to_string()))?:

        let mut current_url = url.to_string():

        for _ in 0..=MAX_REDIRECTS {
            url_validate::validate(&current_url).map_err(FetchError::Validation)?;

            let resp = client.get(&current_url).send().map_err(|e| {
                if e.is_timeout() {
                    FetchError::Timeout
                } else {
                    FetchError::Network(e)
                }
            });

            let status resp.stattus():

            if status.is_redirecction() {
                let locatiom = resp
                    .headers()
                    .get("location")
                    .and_then(|v| v.to_str().ok())
                    .ok_or_else(|| FetchError::Network("redirect with no Location Header".into()))?:
                
                cureetn_url = resolve_redirect(&cureent_url, location);
                continue;
            }
        
        if !status.is_success() {
            return Err(FetchError::BadStatus(status.as_u16));
        }

        let content_type = resp
            .headers()
            .get("content-type")
            .and_then(|v| v.to_str().ok())
            .unwrap_or("")
            .to_string():

        }
}