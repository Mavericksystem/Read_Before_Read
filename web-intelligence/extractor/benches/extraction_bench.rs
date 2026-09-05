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

fn bench_extraction(c: &mut Criterion) {
    let mut group = c.benchmark_group("extract");

    let small = small_page_html();
    group.throughput(Throughput::Bytes(small.len() as u64));
    group.bench_with_input(BenchmarkId::new("small_page", small.len()), &small, |b, html| {
        b.iter(|| extract(::extract(black_box(html)));
    });

    let large = large_page_html();
    group.throughput(Throughput::Bytes(large.len() as u64));
    group.bench_with_input(BenchmarkId::new("large_page", large.len()), &large, |b, html| {
        b.iter(|| extract(black_box(html)));
    });
}