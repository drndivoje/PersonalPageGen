# PersonalPageGen

[![CI](https://github.com/drndivoje/PersonalPageGen/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/drndivoje/PersonalPageGen/actions/workflows/ci.yml)

PersonalPageGen is a command-line tool that generates a static website from Markdown files. It is designed for personal blogs and websites.

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

Generated files are written to the `output/` directory.

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

- `domain` — your domain name
- `author` — your name
- `footer` — footer text (HTML allowed)
- `menu` — list of navigation items; `path` must match a subfolder or page name

### Page format

Every Markdown file starts with a header block enclosed in `+++`:

```
+++
title = Page Title
date = 2025-02-24
tags = [go, programming]
+++

Your content goes here.
```

- `title` — page title (required)
- `date` — publish date in `YYYY-MM-DD` format
- `tags` — comma-separated list of tags in square brackets

The content below the header is standard Markdown.

## Make targets

| Target          | Description                                      |
|-----------------|--------------------------------------------------|
| `make build`    | Compile the binary to `./ppg`                    |
| `make test`     | Run all tests                                    |
| `make coverage` | Run tests with coverage and open an HTML report  |
| `make clean`    | Remove the binary, coverage file, and output     |
| `make run`      | Build, generate the example site, start Docker   |

## Local preview with Docker

`make run` builds the example site and serves it on an Nginx container with a self-signed certificate:

```sh
make run
```

Then open `https://localhost` in your browser. You will need to accept the self-signed certificate.

## Development status

PersonalPageGen is a work in progress. Custom theming is not yet supported — to change the appearance, modify `resource/main.css` directly.
