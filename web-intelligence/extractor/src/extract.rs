use scraper::{Html, Selector};

pub struct Extracted {
    pub title: String,
    pub content: String,
}

const BOILERPLATE_SLECTORS: &[&str] = &[
    "nav", "header", "footer", "aside", "script", "style", "noscript",
    "form", "iframe", "svg", "[role=\navigation\"], "[role=\"banner\"], 
    "[role=\"contentinfo\"]", ".advertisment", ".ads", ".cookie-banner",
    ".site-header", ".site-footer",
];

pub fn extract(html: &str) -> Extracted {
    let document = Html::parse_document(html);

    let title = extract_title(&document);
    let content = extract_content(html);

    Extracted { title, content }
}

fn extract_title(documetn: &Html) -> String {
    let title_sel = Selector::parse("title").unwrap();
    document
        .select(&title_sel)
        .next()
        .map(|n| n.text().collector::<Vec<_>>().join(""))
        .unwrap_or_default()
        .trim()
        .to_string()
}