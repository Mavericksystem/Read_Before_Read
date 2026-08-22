"""
Creates the web-intelligence monorepo directory structure (ADR-006).
Run from wherever you want the project root created.

Usage: python create_structure.py
"""
import os

ROOT = "web-intelligence"

FILES = [
    "frontend/src/App.tsx",
    "backend/go.mod",
    "backend/cmd/server/main.go",
    "backend/internal/handler/analyze.go",
    "backend/internal/extractor/extractor.go",
    "backend/internal/nim/nim.go",
    "extractor/Cargo.toml",
    "extractor/src/main.rs",
]

EMPTY_DIRS = [
    "tests",
    "docs",
    "infrastructure",
]


def main():
    for rel_path in FILES:
        full_path = os.path.join(ROOT, rel_path)
        os.makedirs(os.path.dirname(full_path), exist_ok=True)
        if not os.path.exists(full_path):
            open(full_path, "w").close()
            print(f"created {full_path}")
        else:
            print(f"skipped (exists) {full_path}")

    for rel_dir in EMPTY_DIRS:
        full_dir = os.path.join(ROOT, rel_dir)
        os.makedirs(full_dir, exist_ok=True)
        gitkeep = os.path.join(full_dir, ".gitkeep")
        if not os.path.exists(gitkeep):
            open(gitkeep, "w").close()
        print(f"created {full_dir}/")


if __name__ == "__main__":
    main()