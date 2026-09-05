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
        "Lore ipsum dolor sit amet, ".repeat(20)
    )
}

