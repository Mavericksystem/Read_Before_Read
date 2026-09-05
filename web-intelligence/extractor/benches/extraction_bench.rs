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

    group.finish();
}

fn bench_process_strtup(c: &mut Criterion) {
    let binary = env!("CARGO_BIN_EXE_extractor");

    c.bench_function("subprocess_spawn_overhead", |b| {
        b.iter(|| {
            use std::io::Write;
            let mut child = Command::new(binary)
                .stdin(std::process::Stdio::piped())
                .stdout(std::process::Stdio::piped())
                .stderr(std::process::Stdio::null())
                .spawn()
                .expect("Failed to spawn extractor binary");

            let req = r#"{"url""not-a-calid-url", "max_response_bytes", "timeout_ms":100}"#;
            child
                .stdin
                .as_mut()
                .unwrap()
                .write_all(req.as_bytes())
                .unwrap();
            let _ = child.wait_with_output().unwrap();
        });
    });
}

critrion_group!(benches, bench_extraction, bench_process_strtup);
criterion_main!(benches);