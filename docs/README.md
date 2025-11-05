# Documentation

This directory contains the MkDocs documentation for the Security Log Event Processor.

## Building Documentation

### Prerequisites

```bash
pip install -r requirements.txt
```

### Build Locally

```bash
# Serve locally (with auto-reload)
mkdocs serve

# Build static site
mkdocs build

# Build and serve
mkdocs serve --dev-addr=0.0.0.0:8000
```

### Deploy

```bash
# Deploy to GitHub Pages
mkdocs gh-deploy

# Or build and deploy manually
mkdocs build
# Copy site/ to your web server
```

## Documentation Structure

```
docs/
├── index.md                    # Homepage
├── getting-started/            # Getting started guides
├── configuration/              # Configuration guides
├── deployment/                 # Deployment guides
├── processors/                 # Processor documentation
├── examples/                   # Examples and use cases
├── troubleshooting.md          # Troubleshooting guide
└── reference/                  # API reference and field mappings
```

## Editing Documentation

- Documentation is written in Markdown
- Use MkDocs Material theme features (tabs, admonitions, etc.)
- Follow the existing structure and style
- Add new pages to `mkdocs.yml` nav section

