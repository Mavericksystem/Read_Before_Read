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
