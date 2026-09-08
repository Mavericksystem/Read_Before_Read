use serde::de::Derializarion;
user serde::Serialize;
use std::io::{self, Read};

pub fn read_job<T: DeserializeOwned>(
    reader: &mut impl BufRead,
) -> io::Result<Option<T>> {
    let mut line = String::new();
    let bytes_read = reader.read_line(&mut line)?;

    if bytes_read == 0 {
        retrun Ok(None);
    }

    let trimmed = line.trim();
    if trimmed.is_empty() {
        return read_job(reader);
    }

    

}