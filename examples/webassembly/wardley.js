// wardley.js — Auto-render library for Wardley Maps
// Finds <pre><code class="language-wtg2"> blocks and renders them as SVG.
// Usage: <script src="wardley.js"></script>
//
// Configuration (set before loading this script):
//   window.WardleyMap = { baseURL: "/assets/", defaultWidth: 1300, ... }

(function () {
  "use strict";

  // Capture script element before any async work (may be null with async/defer)
  var currentScript = document.currentScript;

  // === SECTION 1: Configuration ===
  var userConfig = window.WardleyMap || {};
  var config = {
    baseURL: userConfig.baseURL || "",
    selector: userConfig.selector || "pre code.language-wtg2",
    defaultWidth: userConfig.defaultWidth || 1300,
    defaultHeight: userConfig.defaultHeight || 900,
    withAnnotations: userConfig.withAnnotations || false,
    autoRender: userConfig.autoRender !== false,
    wasmFile: userConfig.wasmFile || "wtg.wasm",
  };

  // === SECTION 2: Dynamic wasm_exec.js loader ===
  // Load wasm_exec.js from the same directory as this script (or baseURL).
  // This keeps it in sync with the Go version used to build the WASM binary.
  function resolveBaseURL() {
    if (config.baseURL) {
      var base = config.baseURL;
      if (base.charAt(base.length - 1) !== "/") {
        base += "/";
      }
      return base;
    }
    if (currentScript && currentScript.src) {
      var parts = currentScript.src.split("/");
      parts.pop();
      return parts.join("/") + "/";
    }
    return "";
  }

  var wasmExecPromise = null;

  function loadWasmExec() {
    if (typeof globalThis.Go !== "undefined") {
      return Promise.resolve();
    }
    if (wasmExecPromise) {
      return wasmExecPromise;
    }
    wasmExecPromise = new Promise(function (resolve, reject) {
      var script = document.createElement("script");
      script.src = resolveBaseURL() + "wasm_exec.js";
      script.onload = resolve;
      script.onerror = function () {
        reject(new Error("Failed to load wasm_exec.js from " + script.src));
      };
      document.head.appendChild(script);
    });
    return wasmExecPromise;
  }

  // === SECTION 3: WASM Loader ===
  var wasmPromise = null;

  function loadWasm() {
    if (wasmPromise) {
      return wasmPromise;
    }

    // If generateSVG is already available (e.g., playground loaded it), skip
    if (typeof globalThis.generateSVG === "function") {
      wasmPromise = Promise.resolve();
      return wasmPromise;
    }

    wasmPromise = loadWasmExec().then(function () {
      var wasmURL = resolveBaseURL() + config.wasmFile;
      var go = new Go();

      return WebAssembly.instantiateStreaming(
        fetch(wasmURL),
        go.importObject
      )
        .then(function (result) {
          go.run(result.instance);
        })
        .catch(function () {
          // Fallback for servers that don't set application/wasm MIME type
          return fetch(wasmURL)
            .then(function (resp) {
              return resp.arrayBuffer();
            })
            .then(function (bytes) {
              return WebAssembly.instantiate(bytes, go.importObject);
            })
            .then(function (result) {
              go.run(result.instance);
            });
        });
    });

    return wasmPromise;
  }

  // === SECTION 4: CSS Injection ===
  function injectStyles() {
    var style = document.createElement("style");
    style.textContent = [
      ".wardley-map { display: block; max-width: 100%; }",
      ".wardley-map svg { max-width: 100%; height: auto; }",
      ".wardley-map-loading { padding: 20px; color: #666; font-family: sans-serif; font-size: 14px; }",
      ".wardley-map-error { padding: 20px; color: #c00; background: #fff0f0; border: 1px solid #fcc; border-radius: 4px; font-family: sans-serif; font-size: 14px; }",
      ".wardley-map-error details { margin-top: 10px; }",
      ".wardley-map-error pre { background: #f5f5f5; padding: 10px; overflow-x: auto; font-size: 12px; }",
    ].join("\n");
    document.head.appendChild(style);
  }

  // === SECTION 5: DOM Scanner ===
  function findCodeBlocks() {
    return Array.prototype.slice.call(
      document.querySelectorAll(config.selector)
    );
  }

  // === SECTION 6: Renderer ===
  function getElementConfig(codeEl) {
    var preEl = codeEl.parentElement;
    // Check data attributes on both <code> and <pre>
    function attr(name) {
      return (
        codeEl.getAttribute("data-" + name) ||
        (preEl && preEl.getAttribute("data-" + name)) ||
        null
      );
    }
    return {
      width: parseInt(attr("width"), 10) || config.defaultWidth,
      height: parseInt(attr("height"), 10) || config.defaultHeight,
      annotations:
        attr("annotations") === "true" ? true : config.withAnnotations,
    };
  }

  function renderBlock(codeEl) {
    var source = codeEl.textContent.trim();
    if (!source) return;

    var preEl = codeEl.parentElement;
    var target = preEl && preEl.tagName === "PRE" ? preEl : codeEl;

    var container = document.createElement("div");
    container.className = "wardley-map wardley-map-loading";
    container.textContent = "Rendering Wardley Map...";

    target.parentNode.replaceChild(container, target);

    var elConfig = getElementConfig(codeEl);

    loadWasm()
      .then(function () {
        var result = generateSVG(
          source,
          elConfig.width,
          elConfig.height,
          elConfig.annotations
        );

        container.className = "wardley-map";

        if (typeof result === "string" && result.indexOf("<svg") !== -1) {
          container.innerHTML = result;
          container.classList.add("wardley-map-rendered");
          container.dispatchEvent(
            new CustomEvent("wardley:rendered", {
              bubbles: true,
              detail: { source: source, success: true },
            })
          );
        } else {
          showError(container, source, String(result));
        }
      })
      .catch(function (err) {
        container.className = "wardley-map";
        showError(container, source, "Failed to load WASM: " + err.message);
      });
  }

  function showError(container, source, message) {
    container.classList.add("wardley-map-error");
    container.innerHTML =
      "<strong>Wardley Map Error</strong><br>" +
      escapeHTML(message) +
      "<details><summary>Source</summary><pre>" +
      escapeHTML(source) +
      "</pre></details>";
    container.dispatchEvent(
      new CustomEvent("wardley:rendered", {
        bubbles: true,
        detail: { source: source, success: false, error: message },
      })
    );
  }

  function escapeHTML(str) {
    var div = document.createElement("div");
    div.appendChild(document.createTextNode(str));
    return div.innerHTML;
  }

  // === SECTION 7: Initialization & Public API ===
  function renderAll() {
    var blocks = findCodeBlocks();
    if (blocks.length === 0) return;

    // Render blocks one at a time using setTimeout to yield to the browser
    var i = 0;
    function next() {
      if (i < blocks.length) {
        renderBlock(blocks[i]);
        i++;
        setTimeout(next, 0);
      }
    }
    next();
  }

  function renderElement(el) {
    var codeEl = el;
    if (el.tagName === "PRE") {
      codeEl = el.querySelector("code") || el;
    }
    renderBlock(codeEl);
  }

  function init() {
    injectStyles();
    if (config.autoRender) {
      renderAll();
    }
  }

  // Expose public API
  window.WardleyMap = window.WardleyMap || {};
  window.WardleyMap.render = renderElement;
  window.WardleyMap.renderAll = renderAll;
  Object.defineProperty(window.WardleyMap, "ready", {
    get: function () {
      return loadWasm();
    },
  });
  // Preserve user config
  window.WardleyMap.config = config;

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
