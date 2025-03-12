# PersonalPageGen

PersonalPageGen is a command-line application that generates a static website from markdown files. This tool helps you quickly create and deploy a personal blog as a static website.

## Input file structure


The input folder is where all your content for the static website in markdown format is stored. To use PersonalPageGen, ensure that your input folder contains a `config.yaml` file with the following structure:
```yaml
domain: <your domain>
author: <your name>
menu:
  - title: About
    path: about
  - title: Blog
    path: blog
```
- `domain`: Your domain, for example, `example.com`.
- `author`: Your name.
- `menu`: A list of menu items that will be rendered as a horizontal menu bar. `title` is the label for each menu item, and `path` is the path to the corresponding page.

The sample page in markdown format look like:
```
+++
title=Page title
date=2025-02-24
+++

Your content goes here
The header of the page is enclosed between `+++` characters. The header contains the following properties:
- **title**: The title of the page.
- **date**: The date when the page is published.


Since PersonalPageGen is designed for personal websites, which typically include blog pages, the input folder should also contain a `blog` subfolder where your blog posts are stored.

## Installation

To install PersonalPageGen, clone the repository and build the application:

```sh
git clone https://github.com/yourusername/PersonalPageGen.git
cd PersonalPageGen
go build
```

Generate the static web site from example folder
```sh
ppg example
```

or run

```sh
go run . example
```
the site will be inside the output folder.

## Development Status

PersonalPageGen is still a work in progress. Theming is not yet supported, and custom CSS must be manually adapted. If you want to apply your own CSS styles, you need to modify `resource/main.css`.

