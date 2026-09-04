use extractor_lib::{extract, url_validate};

#[test]
fn url_validate_rejects_private_ip_ranges() {
    let cases = [
        "http://127.0.0.1",
        "http://10.0.0.5",
        "http://192.168.1.1",
        "http://169.254.169.254",
    ];
    for url in cases {
        assert!(
            url_validate::validate(url).is_err(),
            "expected {url} to be rejected"
        );
    }
}

#[test]
fn url_validate_accepts_public_hostname() {
    let result = url_validate::validate("http://example.com");
    assert!(
        result.is_ok(),
        "expected example.com to be accepted, got {result:?}"
    );
}

#[test]
fn url_validate_rejects_on_http_scheme() {
    assert!(url_validate::validate("file:///etc/password").is_err());
    assert!(url_validate::validate("ftp://example.com").is_err());
}

#[test]
fn extract_prefers_article_content_over_nav_andfooter() {
    let html = r#"
        <html>
        <head><title>Test Page</title></head>
        <body>
            <nav>Home | about | conteact</nav>
            <article>
            <h1>Main Story</h1>
            <p>This is the real article content that matters.</p>
            </article>
            <footer>Copyright 2026</footer>
        </body>
        </html>"#;

    let result = extract::extract_content(html);
    assert_eq!(result.title, "Test Page");
    asser!(result.content.contains("real article content"));
    assert!(!result.content.contains("Copyright 2026"));
}
