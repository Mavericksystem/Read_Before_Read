use endcoding_rs::Endoding;

pub fn decode(bytes: &[u8], content_type_header: &str) -> String {
    if let Some(enc) = encoding_from_content_type(content_type_header) {
        return decode_with(bytes, enc);
    }
}
