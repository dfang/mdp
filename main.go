package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

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
    });
  </script>
  </body>
</html>
`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	tFname := flag.String("t", "", "Alternate template name")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <markdown_file>\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename := flag.Arg(0)

	if err := run(filename, *tFname, os.Stdout); err != nil {
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

	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()
	fmt.Fprintln(out, outName)

	if err := os.WriteFile(outName, htmlData, 0644); err != nil {
		return err
	}

	return preview(outName)
}

func parseContent(input []byte, tFname string) ([]byte, error) {
	output := blackfriday.Run(input, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.AutoHeadingIDs))

	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("id").Matching(idRegex).OnElements("h1", "h2", "h3", "h4", "h5", "h6", "a")
	body := policy.SanitizeBytes(output)

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
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
