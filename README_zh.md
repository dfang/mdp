<h1 align="center">mdp (Markdown Preview)</h1>

*<p align="center">基于 Go 开发的高性能、优雅的 Markdown 实时预览命令行工具</p>*

<p align="center">
  <a href="README.md">English</a> •
  <a href="README_zh.md">简体中文</a>
</p>

## 简介

`mdp` 是一个轻量级、跨平台的命令行工具，可将 Markdown 文件快速转换为美观独立的 HTML 页面，并自动在你的默认浏览器中打开预览。

### 特性

- 🚀 **零配置 & 单二进制文件**：所有模板通过 Go `embed` 内置打包，无需额外模板文件或运行时依赖。
- 🎨 **9 款精美内置主题**：包含 GitHub、Notion、Dracula、Nord、VitePress、学术风、复古报纸风等。
- 🌓 **明亮 / 暗色模式**：支持暗色模式的主题配有悬浮 ☀️/🌙 切换按钮，支持跟随系统偏好与本地状态持久化存储。
- 📋 **代码一键复制**：代码块自带复制按钮与即时反馈。
- ⚓ **全能锚点跳转**：支持平滑滚动，完美兼容中文、英文及复合字符标题 ID 解析与跳转。
- 🔝 **回到顶部按钮**：悬浮快速置顶按钮，长文档阅读体验极佳。
- 🖨️ **打印与 PDF 导出优化**：内置 `@media print` 样式，在浏览器中按 `Cmd + P` / `Ctrl + P` 即可一键导出干净整洁的 PDF。
- 🔒 **HTML 安全过滤**：集成 [bluemonday](https://github.com/microcosm-cc/bluemonday) 防止 XSS 注入。
- 💻 **全平台支持**：支持 Linux、macOS 与 Windows。

---

## 安装

### 环境要求
- [Go](https://go.dev/) (1.18 或更高版本)

### 通过 `go install` 安装：
```bash
go install github.com/dfang/mdp@latest
```

### 直接下载二进制包：
前往 [GitHub Releases](https://github.com/dfang/mdp/releases) 页面下载适合你操作系统的预编译包。

---

## 命令行参数

```text
Usage: mdp [options] <markdown_file>
       mdp templates
       mdp version

Options:
  -l, -list
        列出所有内置模板及其样式说明
  -r, -random
        随机使用一个内置模板
  -t string
        指定内置模板名称或自定义模板文件路径
  -v, -version
        显示版本信息
```

---

## 内置模板一览

运行 `mdp templates` 或 `mdp -l` 查看所有可用模板：

| 模板名称 | 明暗模式 | 样式特点 |
| :--- | :---: | :--- |
| `academic` | 浅色 | 经典学术论文风格，衬线字体排版，居中标题 |
| `dark_modern` | 深色 | 未来感暗色磨砂玻璃（Glassmorphism）UI，霓虹渐变光晕 |
| `dracula` | 深色 | 经典 Dracula 暗黑主题，紫色、粉色、青色点缀 |
| `github` | ☀️/🌙 双模式 | 官方 GitHub Markdown 样式，支持自动/手动明暗切换 |
| `newsprint` | 浅色 | 暖象牙色复古报刊阅读质感 |
| `nord` | ☀️/🌙 双模式 | 极地冷灰蓝 Nord 极简北欧主题，支持明暗切换 |
| `notion` | ☀️/🌙 双模式 | Notion 极简笔记风格，优雅清晰，支持明暗切换 |
| `solarized_dark`| ☀️/🌙 双模式 | 经典 Solarized 配色调色板，支持明暗切换 |
| `vitepress` | ☀️/🌙 双模式 | VitePress / Vue 现代技术文档风格，支持明暗切换 |

---

## 使用示例

### 1. 默认预览 Markdown 文件：
```bash
mdp README.md
```

### 2. 查看所有内置模板：
```bash
mdp templates
# 或者
mdp -l
```

### 3. 指定内置模板预览：
```bash
mdp -t github README.md
mdp -t vitepress README.md
mdp -t nord README.md
mdp -t notion README.md
```

### 4. 随机选择模板预览：
```bash
mdp -r README.md
```

### 5. 使用自定义 Go 模板文件：
你可以传入自定义的 Go `html/template` 文件进行个性化定制：

```bash
mdp -t custom.tmpl README.md
```

自定义模板中只需在 `<body>` 内渲染 `{{ .Body }}` 即可：

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
