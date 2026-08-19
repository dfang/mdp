<h1 align="center">mdp (Markdown Preview)</h1>

*<p align="center">A Markdown Preview CLI Tool Built in Go</p>*


## About

Use the `mdp` CLI tool to preview a Markdown file in the browser.


### Features:

- Cross platform:  Linux / Macos / Windows.

- Convert a Markdown file into an HTML file.

- Support custom templates.

- `mdp` uses [bluemonday](https://github.com/microcosm-cc/bluemonday) to sanitize the Markdown content creating a safe HTML file.

## Installation

### Requirements:

- [Go](https://go.dev/)

### How to install:

- Run: 

  ```
  $ go install github.com/dfang/mdp@latest
  ```

## Usage

### Options:

```
$ mdp
Usage: mdp [options] <markdown_file>

Options:
  -t string
        Alternate template name
```

### Examples:

#### Preview a Markdown file in the browser:

```
$ mdp MyFile.md
```

#### Use a custom template file to generate the HTML:

A custom template file may be used to provide custom styles and layouts. 

`mdp` accepts a Go [html/template](https://pkg.go.dev/html/template) file provided using the `-t` flag.

```
$ mdp -t template.tmpl MyFile.md
```
The content of the markdown file must be defined in the template file inside de `<body>` tags as `{{ .Body }}`.

  Ex.: `custom.tmpl`

  ```
  ...
    </head>
    <body>
      {{ .Body }}
    </body>
    <footer>
      <p>© 2022</p>
    </footer>
  ```
