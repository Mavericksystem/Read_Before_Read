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