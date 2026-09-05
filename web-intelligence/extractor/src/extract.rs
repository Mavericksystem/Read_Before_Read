use scraper::{Html, Selector};

pub struct Extracted {
    pub title: String,
    pub content: String,
}

const BOILERPLATE_SELECTORS: &[&str] = &[
    "nav",
    "header",
    "footer",
    "aside",
    "script",
    "style",
    "noscript",
    "form",
    "iframe",
    "svg",
    "[role=\"navigation\"]",
    "[role=\"banner\"]",
    "[role=\"contentinfo\"]",
    ".advertisement",
    ".ads",
    ".cookie-banner",
    ".site-header",
    ".site-footer",
];

pub fn extract(html: &str) -> Extracted {
    extract_content(html)
}

fn extract_title(document: &Html) -> String {
    let title_sel = Selector::parse("title").unwrap();
    document
        .select(&title_sel)
        .next()
        .map(|n| n.text().collect::<Vec<_>>().join(""))
        .unwrap_or_default()
        .trim()
        .to_string()
}

pub fn extract_content(html: &str) -> Extracted {
    let document = Html::parse_document(html);
    let title = extract_title(&document);
    let content = extract_body_content(&document);

    Extracted { title, content }
}

fn extract_body_content(document: &Html) -> String {
    for selector_str in &["article", "main", "[role=\"main\"]"] {
        if let Ok(sel) = Selector::parse(selector_str) {
            if let Some(node) = document.select(&sel).next() {
                let text = clean_text(&node.text().collect::<Vec<_>>().join(" "));
                if !text.is_empty() {
                    return text;
                }
            }
        }
    }

    let body_sel = Selector::parse("body").unwrap();
    let Some(body) = document.select(&body_sel).next() else {
        return String::new();
    };

    let mut exclude = std::collections::HashSet::new();
    for selector_str in BOILERPLATE_SELECTORS {
        if let Ok(sel) = Selector::parse(selector_str) {
            for node in document.select(&sel) {
                exclude.insert(node.text().collect::<Vec<_>>().join(" "));
            }
        }
    }

    let full_text = body.text().collect::<Vec<_>>().join(" ");
    let filtered = if exclude.is_empty() {
        full_text
    } else {
        full_text.split_whitespace().collect::<Vec<_>>().join(" ")
    };

    clean_text(&filtered)
}

fn clean_text(s: &str) -> String {
    s.split_whitespace()
        .collect::<Vec<_>>()
        .join(" ")
        .trim()
        .to_string()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn prefers_article_tag_over_body() {
        let html = r#"
            <html><body>
                <nav>Home About COntact</nav>
                <article><p>The real content goes here.</p></article>
                <footer>Copyright 2026</footer>
            </body></html>"
        "#;
        let result = extract_content(html);
        assert!(result.content.contains("real content"));
    }

    #[test]
    fn falls_back_to_body_when_no_article_tag() {
        let html = r#"<html><body><p> Just a plain page.</p></body></html>"#;
        let result = extract_content(html);
        assert!(result.content.contains("Just a plain page."));
    }
}
