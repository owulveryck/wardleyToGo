# wardley.js

A drop-in JavaScript library that automatically renders Wardley Maps from fenced code blocks in HTML pages. Works like [Mermaid.js](https://mermaid.js.org/) but for Wardley Maps using the WTG2 syntax.

## Quick Start

### Using GitHub Pages (recommended)

The WASM binary and library are built and published to GitHub Pages by CI on every push to `main`. You can use them directly:

```html
<script src="https://owulveryck.github.io/wardleyToGo/wardley.js"></script>
```

The library automatically loads `wtg.wasm` from the same directory as the script, so no additional configuration is needed.

Then write your maps in standard fenced code blocks. When rendered as HTML (by any Markdown processor — Hugo, Jekyll, GitHub Pages, etc.), these become `<pre><code class="language-wtg2">` elements that `wardley.js` auto-discovers and renders:

```html
<pre><code class="language-wtg2">
title: Tea Shop

anchor Customer

Cup of Tea : III.5
Tea : II.8
Hot Water : III.7
Cup : IV.5
Water : IV.8
Kettle : III.5
Power : IV.9

Customer -> Cup of Tea -> Tea
Cup of Tea -> Hot Water -> Water
Hot Water -> Kettle -> Power
Cup of Tea -> Cup
</code></pre>
```

### Self-hosting

Download both files and serve them from the same directory:

- `wardley.js` — the library (includes the Go WASM runtime)
- `wtg.wasm` — the compiled rendering engine

```html
<script src="/path/to/wardley.js"></script>
```

## How It Works

1. You include `wardley.js` via a `<script>` tag
2. On page load, it scans the DOM for `<pre><code class="language-wtg2">` elements
3. It loads the `wtg.wasm` binary (Go compiled to WebAssembly)
4. Each code block is replaced with a rendered SVG Wardley Map
5. Errors are shown inline with the original source code

## Configuration

Set `window.WardleyMap` **before** loading the script to customise defaults:

```html
<script>
window.WardleyMap = {
  baseURL: "/assets/wardley/",  // directory containing wtg.wasm
  defaultWidth: 1300,           // default SVG width (px)
  defaultHeight: 900,           // default SVG height (px)
  withAnnotations: false,       // show evolution stage annotations
  autoRender: true,             // render on page load
  selector: "pre code.language-wtg2",  // CSS selector for code blocks
  wasmFile: "wtg.wasm"          // WASM binary filename
};
</script>
<script src="https://owulveryck.github.io/wardleyToGo/wardley.js"></script>
```

### Per-block options

Override width, height, or annotations on individual blocks using `data-` attributes on the `<pre>` element:

```html
<pre data-width="800" data-height="500" data-annotations="true">
  <code class="language-wtg2">
    title: Small map
    anchor User
    Need : III.5
    User -> Need
  </code>
</pre>
```

## JavaScript API

For dynamic content (e.g., single-page apps), use the programmatic API:

```javascript
// Re-scan and render all wtg2 code blocks
WardleyMap.renderAll();

// Render a specific element
WardleyMap.render(document.getElementById("my-map"));

// Wait for WASM to be loaded
WardleyMap.ready.then(function() {
  console.log("WASM is ready");
});
```

## Events

Each rendered block dispatches a `wardley:rendered` event:

```javascript
document.addEventListener("wardley:rendered", function(e) {
  console.log(e.detail.success);  // true or false
  console.log(e.detail.source);   // original WTG2 source
  if (!e.detail.success) {
    console.log(e.detail.error);  // error message
  }
});
```

## Use in Markdown

Most Markdown processors convert fenced code blocks into `<pre><code class="language-...">`. This means you can write Wardley Maps in your Markdown files:

````markdown
# My Architecture Document

Here is the current state of our platform:

```wtg2
title: Platform Map

anchor User

Frontend : III.5
API Gateway : III.8
Auth Service : IV.2
Database : IV.7

User -> Frontend -> API Gateway -> Auth Service
API Gateway -> Database
```
````

Include `wardley.js` in your site template, and the maps render automatically.

### Hugo

Add to your base template (`layouts/_default/baseof.html`):

```html
<script src="https://owulveryck.github.io/wardleyToGo/wardley.js"></script>
```

### Jekyll

Add to `_includes/head.html` or your layout:

```html
<script src="https://owulveryck.github.io/wardleyToGo/wardley.js"></script>
```

### Plain HTML

```html
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body>
  <pre><code class="language-wtg2">
title: Hello Wardley
anchor User
Service : III.5
User -> Service
  </code></pre>
  <script src="https://owulveryck.github.io/wardleyToGo/wardley.js"></script>
</body>
</html>
```

## WTG2 Syntax Reference

See the [grammar file](../../wtg2/grammar.bnf) and [test examples](../../parser/wtg2/testdata/) for the full WTG2 syntax. Key elements:

```
title: Map Title                    # metadata
anchor Customer                     # user/anchor node (top of value chain)
Component Name : III.5              # component at evolution position
Component : III.5 >> IV.2           # component with evolution movement
Component : III.5 (buy)             # component with type (build/buy/outsource)

Customer -> Component               # dependency edge
A -> B -> C                         # chained dependencies
A -[label]-> B                      # labelled edge
A <-> B                             # bidirectional edge

pipeline Component {                # pipeline grouping
  SubA : III.2
  SubB : IV.1
}

group "Label" {                     # visual group
  Component1
  Component2
}
```

Evolution positions use Roman numerals I-IV with a decimal (0-9):
- `I.0` to `I.9` — Genesis
- `II.0` to `II.9` — Custom-Built
- `III.0` to `III.9` — Product
- `IV.0` to `IV.9` — Commodity

## Building from Source

```bash
cd examples/webassembly

# Build everything (requires Go 1.25+)
make all

# Serve locally
cd assets && python3 -m http.server 8080

# Open http://localhost:8080/wardley-example.html
```

## CI/CD

The [pages workflow](../../.github/workflows/pages.yml) runs on every push to `main`:

1. Builds `wtg.wasm` from `examples/webassembly/wasm/`
2. Copies `wardley.js` to the deployment directory
3. Deploys to GitHub Pages at `https://owulveryck.github.io/wardleyToGo/`

The WASM binary is ~9 MB uncompressed but compresses well with gzip (~2-3 MB). Most web servers and CDNs serve it compressed automatically.
