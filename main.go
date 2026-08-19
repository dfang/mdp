package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed templates/*
var templateFS embed.FS

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if version == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				if commit == "none" {
					commit = setting.Value
					if len(commit) > 7 {
						commit = commit[:7]
					}
				}
			case "vcs.time":
				if date == "unknown" {
					date = setting.Value
				}
			}
		}
	}
}

var idRegex = regexp.MustCompile(`^[\p{L}\p{N}\p{M}\p{Pd}\p{Pc}_.:-]+$`)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>{{ .Title }}</title>
    <style>
      html {
        scroll-behavior: smooth;
      }
      body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
        line-height: 1.6;
        color: #24292e;
        max-width: 860px;
        margin: 0 auto;
        padding: 32px 20px;
      }
      ul, ol {
        padding-left: 2em;
        margin: 0.8em 0;
      }
      li {
        margin: 0.3em 0;
      }
      li > ul, li > ol {
        margin: 0.2em 0;
      }
      input[type="checkbox"] {
        margin-right: 0.35em;
        vertical-align: middle;
      }
    </style>
  </head>
  <body>
  {{ .Body }}
  <script>
    document.addEventListener("DOMContentLoaded", function() {
      function normalize(str) {
        return (str || "").toLowerCase().replace(/[\s\-_:：、，,.。（）()\[\]【】/\\#]/g, "");
      }

      function resolveAnchor(targetId) {
        if (!targetId) return;
        var decoded = decodeURIComponent(targetId.replace(/^#/, ""));
        var el = document.getElementById(decoded);
        if (!el) {
          try {
            el = document.querySelector('[name="' + CSS.escape(decoded) + '"]');
          } catch(e) {}
        }
        if (!el) {
          var targetNorm = normalize(decoded);
          var headings = document.querySelectorAll("h1, h2, h3, h4, h5, h6, a[name], [id]");
          for (var i = 0; i < headings.length; i++) {
            var h = headings[i];
            var hNorm = normalize(h.id || h.getAttribute("name") || h.textContent);
            if (hNorm === targetNorm) {
              el = h;
              break;
            }
          }
        }
        if (el) {
          el.scrollIntoView({ behavior: "smooth" });
        }
      }

      if (window.location.hash) {
        setTimeout(function() { resolveAnchor(window.location.hash); }, 100);
      }

      document.addEventListener("click", function(e) {
        var a = e.target.closest("a");
        if (a && a.getAttribute("href") && a.getAttribute("href").indexOf("#") === 0) {
          var href = a.getAttribute("href");
          e.preventDefault();
          history.pushState(null, null, href);
          resolveAnchor(href);
        }
      });

      // Mermaid diagrams
      var mermaidBlocks = [];
      document.querySelectorAll("pre code.language-mermaid, pre.language-mermaid code").forEach(function(codeEl, i) {
        var pre = codeEl.closest("pre");
        var container = document.createElement("div");
        container.className = "mermaid";
        container.setAttribute("data-mermaid-src", codeEl.textContent);
        container.style.textAlign = "center";
        container.style.margin = "24px 0";
        container.style.overflowX = "auto";
        pre.parentNode.replaceChild(container, pre);
        mermaidBlocks.push(container);
      });

      if (typeof mermaid !== "undefined" && mermaidBlocks.length > 0) {
        mermaid.initialize({
          startOnLoad: false,
          theme: "default",
          securityLevel: "loose"
        });
        mermaidBlocks.forEach(function(el, index) {
          var code = el.getAttribute("data-mermaid-src");
          if (!code) return;
          var id = "mermaid-svg-" + index + "-" + Math.floor(Math.random() * 10000);
          mermaid.render(id, code).then(function(result) {
            el.innerHTML = result.svg;
          }).catch(function(err) {
            console.warn("Mermaid render error:", err);
          });
        });
      }
    });
  </script>
  <script src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"></script>
  </body>
</html>
`
)

type content struct {
	Title string
	Body  template.HTML
}

func listBuiltinTemplates() ([]string, error) {
	entries, err := templateFS.ReadDir("templates")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".html.tmpl") {
			name := strings.TrimSuffix(entry.Name(), ".html.tmpl")
			names = append(names, name)
		}
	}
	return names, nil
}

func getRandomTemplate() (string, error) {
	names, err := listBuiltinTemplates()
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("no built-in templates found")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return names[r.Intn(len(names))], nil
}

var templateDescriptions = map[string]string{
	"academic":       "Academic paper style with serif typography and centered title",
	"dark_modern":    "Futuristic dark glassmorphic UI with neon gradient glows",
	"dracula":        "Iconic Dracula dark theme with purple, pink, and cyan accents",
	"github":         "GitHub Markdown style with auto/manual light & dark mode",
	"newsprint":      "Warm ivory newsprint / editorial paper reading style",
	"nord":           "Arctic bluish-gray Nord theme with light & dark mode",
	"notion":         "Notion minimalistic notebook style with light & dark mode",
	"solarized_dark": "Classic Solarized palette with light & dark mode",
	"vitepress":      "VitePress / Vue modern tech docs style with light & dark mode",
}

func printTemplates(out io.Writer) error {
	names, err := listBuiltinTemplates()
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Available built-in templates:")
	for _, name := range names {
		desc := templateDescriptions[name]
		if desc != "" {
			fmt.Fprintf(out, "  - %-15s : %s\n", name, desc)
		} else {
			fmt.Fprintf(out, "  - %s\n", name)
		}
	}
	return nil
}

func printVersion(out io.Writer) {
	fmt.Fprintf(out, "mdp %s (commit: %s, date: %s)\n", version, commit, date)
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "-v", "--version":
			printVersion(os.Stdout)
			return
		case "templates", "list", "list-templates":
			if err := printTemplates(os.Stdout); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
	}

	tFname := flag.String("t", "", "Alternate template name or custom template file")
	randomTmpl := flag.Bool("r", false, "Use a random built-in template")
	flag.BoolVar(randomTmpl, "random", false, "Use a random built-in template")
	listTmpl := flag.Bool("l", false, "List available built-in templates")
	flag.BoolVar(listTmpl, "list", false, "List available built-in templates")
	showVersion := flag.Bool("v", false, "Show version information")
	flag.BoolVar(showVersion, "version", false, "Show version information")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <markdown_file>\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "       %s templates\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "       %s version\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		printVersion(os.Stdout)
		return
	}

	if *listTmpl {
		if err := printTemplates(os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	templateName := *tFname
	if *randomTmpl {
		templateName = "random"
	}

	if err := run(filename, templateName, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename, tFname string, out io.Writer) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData, err := parseContent(input, tFname)
	if err != nil {
		return err
	}

	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	nameOnly := strings.TrimSuffix(base, ext)

	tmpDir := "/tmp"
	if runtime.GOOS == "windows" {
		tmpDir = os.TempDir()
	}

	outName := filepath.Join(tmpDir, fmt.Sprintf("mdp_%s.html", nameOnly))
	fmt.Fprintln(out, outName)

	if err := os.WriteFile(outName, htmlData, 0644); err != nil {
		return err
	}

	return preview(outName)
}

var markdownConverter = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
		parser.WithAttribute(),
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

type unicodeIDs struct {
	values map[string]bool
}

func newUnicodeIDs() parser.IDs {
	return &unicodeIDs{values: make(map[string]bool)}
}

func sanitizeAnchorName(text string) string {
	var b strings.Builder
	lastDash := true
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(unicode.ToLower(r))
			lastDash = false
		} else if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func (s *unicodeIDs) Generate(value []byte, kind ast.NodeKind) []byte {
	id := sanitizeAnchorName(string(value))
	if id == "" {
		id = "section"
	}
	base := id
	count := 1
	for s.values[id] {
		id = fmt.Sprintf("%s-%d", base, count)
		count++
	}
	s.values[id] = true
	return []byte(id)
}

func (s *unicodeIDs) Put(value []byte) {
	s.values[string(value)] = true
}

func parseContent(input []byte, tFname string) ([]byte, error) {
	pc := parser.NewContext(parser.WithIDs(newUnicodeIDs()))
	var outBuf bytes.Buffer
	if err := markdownConverter.Convert(input, &outBuf, parser.WithContext(pc)); err != nil {
		return nil, err
	}

	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("id").Matching(idRegex).OnElements("h1", "h2", "h3", "h4", "h5", "h6", "a", "section", "div", "span")
	policy.AllowElements("input")
	policy.AllowAttrs("type").Matching(regexp.MustCompile(`^checkbox$`)).OnElements("input")
	policy.AllowAttrs("checked", "disabled").OnElements("input")
	policy.AllowAttrs("class").OnElements("li", "ul", "ol", "div", "span", "table", "pre", "code")
	body := policy.SanitizeBytes(outBuf.Bytes())

	var t *template.Template
	var err error

	if tFname == "random" {
		selected, err := getRandomTemplate()
		if err != nil {
			return nil, err
		}
		tFname = selected
	}

	if tFname == "" {
		t, err = template.New("mdp").Parse(defaultTemplate)
	} else if _, statErr := os.Stat(tFname); statErr == nil {
		t, err = template.ParseFiles(tFname)
	} else {
		tmplName := tFname
		if !strings.HasSuffix(tmplName, ".html.tmpl") {
			tmplName += ".html.tmpl"
		}
		embedPath := filepath.Join("templates", filepath.Base(tmplName))
		if _, readErr := templateFS.ReadFile(embedPath); readErr == nil {
			t, err = template.ParseFS(templateFS, embedPath)
		} else {
			return nil, fmt.Errorf("template %q not found (neither local file nor built-in template)", tFname)
		}
	}
	if err != nil {
		return nil, err
	}

	c := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(body),
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func preview(fname string) error {
	var cName string
	var cParams []string

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	cParams = append(cParams, fname)

	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	return exec.Command(cPath, cParams...).Run()
}
