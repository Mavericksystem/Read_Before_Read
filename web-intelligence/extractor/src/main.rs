use serde::{Deserialize, Serialize};
use std::io::{self, Read, Write};
use std::time::Duration;

#[derive(Deserialize)]
struct Request {
    url: String,
    max_response_bytes: u64,
    timeout_ms: u64,
}
