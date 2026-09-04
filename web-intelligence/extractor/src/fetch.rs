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
