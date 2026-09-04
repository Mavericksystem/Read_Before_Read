use extractor_lib::{extract, url_validate};

#[test]
fn url_validate_rejects_private_ip_ranges() {
    let cases = [
        "http://127.0.0.1",
        "http://10.0.0.5",
        "http://192.168.1.1",
        "http://169.254.169.254",
    ];
}
