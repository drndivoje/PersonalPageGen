# PersonalPageGen

[![CI](https://github.com/drndivoje/PersonalPageGen/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/drndivoje/PersonalPageGen/actions/workflows/ci.yml)

A command-line tool that generates a static website from Markdown files, designed for personal blogs and websites.

**Live example:** [drnd.rocks](https://drnd.rocks/)

## Requirements

- Go 1.18+
- Docker (optional, for local preview)

## Installation

```sh
git clone https://github.com/drndivoje/PersonalPageGen.git
cd PersonalPageGen
make build
```

## Usage

```sh
./ppg <input-folder>
```

Generated files are written to `output/`.

## Input folder structure

```
input/
├── config.yml
├── index.md
├── about.md
└── blog/
    ├── first-post.md
    └── second-post.md
```

### config.yml

```yaml
domain: example.com
author: Your Name
footer: '© 2025 Your Name'
menu:
  - title: Blog
    path: blog
  - title: About
    path: about
```

| Field    | Description                                                    |
|----------|----------------------------------------------------------------|
| `domain` | Your domain name                                               |
| `author` | Your name                                                      |
| `footer` | Footer text (HTML allowed)                                     |
| `menu`   | Navigation items; `path` must match a subfolder or page name  |

### Page format

Every Markdown file begins with a `+++` header block:

```
+++
title = Page Title
date = 2025-02-24
tags = [go, programming]
+++

Your content goes here.
```

| Field   | Description                              |
|---------|------------------------------------------|
| `title` | Page title (required)                    |
| `date`  | Publish date (`YYYY-MM-DD`)              |
| `tags`  | Comma-separated list in square brackets  |

## Make targets

| Target          | Description                                      |
|-----------------|--------------------------------------------------|
| `make build`    | Compile the binary to `./ppg`                    |
| `make test`     | Run all tests                                    |
| `make coverage` | Run tests with coverage and open an HTML report  |
| `make clean`    | Remove the binary, coverage file, and output     |
| `make run`      | Build, generate the site, and start Docker       |

## Local preview

`make run` generates the site and serves it via an Nginx container with a self-signed certificate:

```sh
make run
```

Open `https://localhost` in your browser and accept the self-signed certificate.

## Status

Work in progress. Custom theming is not yet supported — to change the appearance, edit `resource/main.css` directly.
