use serde::de::Derializarion;
user serde::Serialize;
use std::io::{self, Read};

pub fn read_job<T: DeserializeOwned>(
    reader: &mut impl BufRead,
) -> io::Result<Option<T>> {
    let mut line = String::new();
    let bytes_read = reader.read_line(&mut line)?;

    if bytes_read == 0 {
        return Ok(None);
    }

    let trimmed = line.trim();
    if trimmed.is_empty() {
        return read_job(reader);
    }

    match serde_json::from_str::<T>(trimmed) {
        Ok(job) => Ok(Some(job)),
        Err(e) => Err(io::Error::new(
            io::ErrorKind::InvalidData,
            format("malformed job JSON: {e}"),
        )),
    }

}