package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestAllTemplates(t *testing.T) {
	input := []byte(`
# 目录

- [第一节 概述](#第一节-概述)
- [第二节 快速入门](#第二节-快速入门)

## 第一节 概述

测试内容 1

## 第二节 快速入门

测试内容 2
`)

	names, err := listBuiltinTemplates()
	if err != nil {
		t.Fatal(err)
	}

	if len(names) == 0 {
		t.Fatal("no built-in templates found")
	}

	for _, tmplName := range names {
		t.Run(tmplName, func(t *testing.T) {
			output, err := parseContent(input, tmplName)
			if err != nil {
				t.Fatalf("failed to parse template %s: %v", tmplName, err)
			}
			html := string(output)
			if !strings.Contains(html, "第一节-概述") {
				t.Errorf("template %s missing expected ID", tmplName)
			}
		})
	}
}

func TestListBuiltinTemplates(t *testing.T) {
	names, err := listBuiltinTemplates()
	if err != nil {
		t.Fatalf("listBuiltinTemplates error: %v", err)
	}
	expected := []string{
		"academic",
		"dark_modern",
		"dracula",
		"github",
		"newsprint",
		"nord",
		"notion",
		"solarized_dark",
		"vitepress",
	}
	for _, exp := range expected {
		found := false
		for _, name := range names {
			if name == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected template %q in list, got %v", exp, names)
		}
	}
}

func TestPrintTemplates(t *testing.T) {
	var buf bytes.Buffer
	if err := printTemplates(&buf); err != nil {
		t.Fatalf("printTemplates failed: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Available built-in templates:") {
		t.Errorf("expected header in output, got %q", output)
	}
	if !strings.Contains(output, "github") {
		t.Errorf("expected github template in output, got %q", output)
	}
}

func TestRandomTemplate(t *testing.T) {
	input := []byte("# Hello Random\n\nContent")
	output, err := parseContent(input, "random")
	if err != nil {
		t.Fatalf("parseContent with random failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected non-empty output with random template")
	}
}

func TestRun(t *testing.T) {
	var out bytes.Buffer
	err := run("README.md", "", &out)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	expectedName := "mdp_README.html"
	if !strings.Contains(out.String(), expectedName) {
		t.Errorf("expected output to contain %q, got %q", expectedName, out.String())
	}
}
