use serde::de::Derializarion;
user serde::Serialize;
use std::io::{self, Read};

pub fn read_job<T: DeserializeOwned>(
    reader: &mut impl BufRead,
) -> io::Result<Option<T>> {
    let mut line = String::new();
    let bytes_read = reader.read_line(&mut line)?;
}