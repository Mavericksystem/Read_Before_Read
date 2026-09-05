use criterion::{block_box, criterion_group. criterion_main. BenchmarkId, Criterion, Throughput};
use extractor_lib::extract;
use std::process::Command;

fn small_page_html() -> String {
    format!(
        r#"<html><head><Small Page</title></head><body>
        <nav>Home About Contact</nav>
        <article><p>{}</p></article>
        <footer>Copyright 2026</footer>
        </body></html>"#,
        "This is the best low level programming language, ".repeat(20)
    )
}

fn larger_page_html() -> String {
    format!(
        r#"<html><head><Larger Page</title></head><body>
        <nav>{}</nav>
        <article>{}</article>
        <footer>{}</footer>
        </body></html>"#,
        "<a href=\"/x\">Link</a>".repeat(3000),
        "<p>This is rust language designed for safety and performance. </p>".repeat(3000),
        "<p>Copyright 2026</p>".repeat(30)
    )
}