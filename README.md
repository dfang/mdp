<h1 align="center">mdp (Markdown Preview)</h1>

*<p align="center">A Fast, Beautiful Markdown Preview CLI Tool Built in Go</p>*

## About

`mdp` is a lightweight, cross-platform CLI tool that converts Markdown files into standalone, beautifully styled HTML pages and opens them instantly in your default browser.

### Features

- 🚀 **Zero Configuration & Single Binary**: All templates are embedded via Go `embed`, no external template files required.
- 🎨 **9 Curated Built-in Themes**: From GitHub and Notion to Dracula and Nord.
- 🌓 **Dark / Light Mode Support**: Dual-mode themes include a floating ☀️/🌙 toggle button, system preference auto-detection, and local storage state persistence.
- 📋 **One-Click Code Copy**: Built-in copy buttons with instant feedback for code blocks.
- ⚓ **Universal Anchor Navigation**: Smooth scroll and automatic header ID resolution for English, Chinese, and complex anchors.
- 🔝 **Back to Top Button**: Floating scroll-to-top button for easy reading of long documents.
- 🖨️ **Print / PDF Optimization**: Clean print stylesheets (`@media print`) for exporting crisp PDFs via `Cmd + P`.
- 🔒 **HTML Sanitization**: Powered by [bluemonday](https://github.com/microcosm-cc/bluemonday) for secure output.
- 💻 **Cross-Platform**: Linux, macOS, and Windows.

---

## Installation

### Requirements
- [Go](https://go.dev/) (1.18 or higher)

### Install via `go install`:
```bash
go install github.com/dfang/mdp@latest
```

---

## Usage

```text
Usage: mdp [options] <markdown_file>
       mdp templates

Options:
  -l, -list
        List available built-in templates with descriptions
  -r, -random
        Use a random built-in template
  -t string
        Alternate template name or custom template file
```

---

## Built-in Templates

Use `mdp templates` or `mdp -l` to list all available templates:

| Template Name | Dark/Light Mode | Style Description |
| :--- | :---: | :--- |
| `academic` | Light | Academic paper style with serif typography and centered title |
| `dark_modern` | Dark | Futuristic dark glassmorphic UI with neon gradient glows |
| `dracula` | Dark | Iconic Dracula dark theme with purple, pink, and cyan accents |
| `github` | ☀️/🌙 Dual | GitHub Markdown style with auto/manual light & dark mode |
| `newsprint` | Light | Warm ivory newsprint / editorial paper reading style |
| `nord` | ☀️/🌙 Dual | Arctic bluish-gray Nord theme with light & dark mode |
| `notion` | ☀️/🌙 Dual | Notion minimalistic notebook style with light & dark mode |
| `solarized_dark`| ☀️/🌙 Dual | Classic Solarized palette with light & dark mode |
| `vitepress` | ☀️/🌙 Dual | VitePress / Vue modern tech docs style with light & dark mode |

---

## Examples

### 1. Preview a Markdown file:
```bash
mdp README.md
```

### 2. List available templates with style descriptions:
```bash
mdp templates
# or
mdp -l
```

### 3. Preview with a specific built-in template:
```bash
mdp -t github README.md
mdp -t vitepress README.md
mdp -t nord README.md
mdp -t notion README.md
```

### 4. Preview with a random template:
```bash
mdp -r README.md
```

### 5. Use a custom Go template file:
A custom Go `html/template` file may be used to provide your own styles:

```bash
mdp -t custom.tmpl README.md
```

The template must render the content inside the `<body>` tags as `{{ .Body }}`:

```html
<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Title }}</title>
  </head>
  <body>
    <article>
      {{ .Body }}
    </article>
  </body>
</html>
```
