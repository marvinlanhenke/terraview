# Terraview Documentation

This directory contains a Hugo documentation site for Terraview.

## Prerequisites

- Hugo Extended `>= 0.110.0`
- Node.js 18+
- Go 1.25+

## Local development

```bash
cd docs/hugo-pages
npm install
hugo mod get
hugo server --disableFastRender
```

Visit <http://localhost:1313/> after the server starts.

## Build

```bash
cd docs/hugo-pages
npm install
hugo mod get
hugo
```

The generated site is written to `docs/hugo-pages/public/`.

## Structure

```text
docs/hugo-pages/
|-- hugo.toml
|-- go.mod
|-- package.json
`-- content/en/
    |-- _index.md
    `-- docs/
        |-- _index.md
        |-- getting-started/
        |-- usage/
        `-- reference/
```

## Editing content

- Add new pages under `content/en/docs/`
- Give each page frontmatter with `title`, `linkTitle`, `weight`, and `description`
- Use relative links between pages to keep navigation simple
