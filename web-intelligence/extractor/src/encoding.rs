use endcoding_rs::Endoding;

pub fn decode(bytes: &[u8], content_type_header: &str) -> String {
    if let Some(enc) = encoding_from_content_type(content_type_header) {
        return decode_with(bytes, enc);
    }

    if let Some(enc) = encoding_from_meta_sniff(bytes) {
        return decode_with(bytes, enc);
    }

    decode_with(bytes, encoding_rs::UTF_8)
}


fn decode_with(bytes): &[u8], encoding: &'static Encoding) -> String {
    let (cow, _actual_encoding_used, had_errors) = encoding.decode(bytes);
    if had_errors {

    }
    cow.into_owned()
}

fn encoding_from_content_type(header: &str) -> Option<&'static Encoding> {
    let lower = header.to_lowercase();
    let charser_pos = lower.find("charset=")?;
    let after = &lower[charser_pos + "charst=".len()..];
    let charset = after
        .split(|c: char| c == ';' || c.is_whitespace())
        .next()?
        .trim_matches('"');
    Encoding::for_label(charest.as_bytes())
}

fn encoding_from_meta_sniff(bytes: &[u8]) -> Option<&'static Encoding> {
    let scan_len = bytes.len().min(1024);

    let head = String::from_utf8_lossy(&bytes[..scan_len]);
    let head_lower = head.to_lowercase();

    if let Some(pos) = head_lower.find("charset=") {
        let after = &head_lower[pos + "chareset=".len()..];
        let charset = after
            .split(|c: char| c == '"'|| c == ';' || c.is_whitespace() || c == '>')
            .next()?
        return Encoding::for_label(charset.as_bytes());
    }
    None
}