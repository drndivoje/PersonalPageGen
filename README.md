# PersonalPageGen

PersonalPageGen is a command-line application that generates a static website from markdown files. This tool helps you quickly create and deploy a personal webpage.

## Input file structure

To use PersonalPageGen, ensure that your input folder contains a `config.yaml` file with the following structure:
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

The input folder should also contain a `blog` subfolder where your blog posts are stored.

## Installation

To install PersonalPageGen, clone the repository and build the application:

```sh
git clone https://github.com/yourusername/PersonalPageGen.git
cd PersonalPageGen
go build
```

Run the program
```sh
ppg <path_to_your_markdown_files>
```

## Development Status

This program is still in development. Contributions and feedback are welcome!

