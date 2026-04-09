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

  // === SECTION 2: Embedded wasm_exec.js (Go WASM runtime) ===
  if (typeof globalThis.Go === "undefined") {
    // BEGIN wasm_exec.js — Copyright 2018 The Go Authors. BSD license.
    (function () {
      var enosys = function () {
        var err = new Error("not implemented");
        err.code = "ENOSYS";
        return err;
      };

      if (!globalThis.fs) {
        var outputBuf = "";
        globalThis.fs = {
          constants: {
            O_WRONLY: -1,
            O_RDWR: -1,
            O_CREAT: -1,
            O_TRUNC: -1,
            O_APPEND: -1,
            O_EXCL: -1,
          },
          writeSync: function (fd, buf) {
            outputBuf += decoder.decode(buf);
            var nl = outputBuf.lastIndexOf("\n");
            if (nl != -1) {
              console.log(outputBuf.substr(0, nl));
              outputBuf = outputBuf.substr(nl + 1);
            }
            return buf.length;
          },
          write: function (fd, buf, offset, length, position, callback) {
            if (offset !== 0 || length !== buf.length || position !== null) {
              callback(enosys());
              return;
            }
            var n = this.writeSync(fd, buf);
            callback(null, n);
          },
          chmod: function (path, mode, callback) {
            callback(enosys());
          },
          chown: function (path, uid, gid, callback) {
            callback(enosys());
          },
          close: function (fd, callback) {
            callback(enosys());
          },
          fchmod: function (fd, mode, callback) {
            callback(enosys());
          },
          fchown: function (fd, uid, gid, callback) {
            callback(enosys());
          },
          fstat: function (fd, callback) {
            callback(enosys());
          },
          fsync: function (fd, callback) {
            callback(null);
          },
          ftruncate: function (fd, length, callback) {
            callback(enosys());
          },
          lchown: function (path, uid, gid, callback) {
            callback(enosys());
          },
          link: function (path, link, callback) {
            callback(enosys());
          },
          lstat: function (path, callback) {
            callback(enosys());
          },
          mkdir: function (path, perm, callback) {
            callback(enosys());
          },
          open: function (path, flags, mode, callback) {
            callback(enosys());
          },
          read: function (fd, buffer, offset, length, position, callback) {
            callback(enosys());
          },
          readdir: function (path, callback) {
            callback(enosys());
          },
          readlink: function (path, callback) {
            callback(enosys());
          },
          rename: function (from, to, callback) {
            callback(enosys());
          },
          rmdir: function (path, callback) {
            callback(enosys());
          },
          stat: function (path, callback) {
            callback(enosys());
          },
          symlink: function (path, link, callback) {
            callback(enosys());
          },
          truncate: function (path, length, callback) {
            callback(enosys());
          },
          unlink: function (path, callback) {
            callback(enosys());
          },
          utimes: function (path, atime, mtime, callback) {
            callback(enosys());
          },
        };
      }

      if (!globalThis.process) {
        globalThis.process = {
          getuid: function () {
            return -1;
          },
          getgid: function () {
            return -1;
          },
          geteuid: function () {
            return -1;
          },
          getegid: function () {
            return -1;
          },
          getgroups: function () {
            throw enosys();
          },
          pid: -1,
          ppid: -1,
          umask: function () {
            throw enosys();
          },
          cwd: function () {
            throw enosys();
          },
          chdir: function () {
            throw enosys();
          },
        };
      }

      if (!globalThis.crypto) {
        throw new Error(
          "globalThis.crypto is not available, polyfill required (crypto.getRandomValues only)"
        );
      }

      if (!globalThis.performance) {
        throw new Error(
          "globalThis.performance is not available, polyfill required (performance.now only)"
        );
      }

      if (!globalThis.TextEncoder) {
        throw new Error(
          "globalThis.TextEncoder is not available, polyfill required"
        );
      }

      if (!globalThis.TextDecoder) {
        throw new Error(
          "globalThis.TextDecoder is not available, polyfill required"
        );
      }

      var encoder = new TextEncoder("utf-8");
      var decoder = new TextDecoder("utf-8");

      globalThis.Go = class {
        constructor() {
          this.argv = ["js"];
          this.env = {};
          this.exit = (code) => {
            if (code !== 0) {
              console.warn("exit code:", code);
            }
          };
          this._exitPromise = new Promise((resolve) => {
            this._resolveExitPromise = resolve;
          });
          this._pendingEvent = null;
          this._scheduledTimeouts = new Map();
          this._nextCallbackTimeoutID = 1;

          const setInt64 = (addr, v) => {
            this.mem.setUint32(addr + 0, v, true);
            this.mem.setUint32(addr + 4, Math.floor(v / 4294967296), true);
          };

          const getInt64 = (addr) => {
            const low = this.mem.getUint32(addr + 0, true);
            const high = this.mem.getInt32(addr + 4, true);
            return low + high * 4294967296;
          };

          const loadValue = (addr) => {
            const f = this.mem.getFloat64(addr, true);
            if (f === 0) {
              return undefined;
            }
            if (!isNaN(f)) {
              return f;
            }

            const id = this.mem.getUint32(addr, true);
            return this._values[id];
          };

          const storeValue = (addr, v) => {
            const nanHead = 0x7ff80000;

            if (typeof v === "number" && v !== 0) {
              if (isNaN(v)) {
                this.mem.setUint32(addr + 4, nanHead, true);
                this.mem.setUint32(addr, 0, true);
                return;
              }
              this.mem.setFloat64(addr, v, true);
              return;
            }

            if (v === undefined) {
              this.mem.setFloat64(addr, 0, true);
              return;
            }

            let id = this._ids.get(v);
            if (id === undefined) {
              id = this._idPool.pop();
              if (id === undefined) {
                id = this._values.length;
              }
              this._values[id] = v;
              this._goRefCounts[id] = 0;
              this._ids.set(v, id);
            }
            this._goRefCounts[id]++;
            let typeFlag = 0;
            switch (typeof v) {
              case "object":
                if (v !== null) {
                  typeFlag = 1;
                }
                break;
              case "string":
                typeFlag = 2;
                break;
              case "symbol":
                typeFlag = 3;
                break;
              case "function":
                typeFlag = 4;
                break;
            }
            this.mem.setUint32(addr + 4, nanHead | typeFlag, true);
            this.mem.setUint32(addr, id, true);
          };

          const loadSlice = (addr) => {
            const array = getInt64(addr + 0);
            const len = getInt64(addr + 8);
            return new Uint8Array(this._inst.exports.mem.buffer, array, len);
          };

          const loadSliceOfValues = (addr) => {
            const array = getInt64(addr + 0);
            const len = getInt64(addr + 8);
            const a = new Array(len);
            for (let i = 0; i < len; i++) {
              a[i] = loadValue(array + i * 8);
            }
            return a;
          };

          const loadString = (addr) => {
            const saddr = getInt64(addr + 0);
            const len = getInt64(addr + 8);
            return decoder.decode(
              new DataView(this._inst.exports.mem.buffer, saddr, len)
            );
          };

          const timeOrigin = Date.now() - performance.now();
          this.importObject = {
            go: {
              "runtime.wasmExit": (sp) => {
                sp >>>= 0;
                const code = this.mem.getInt32(sp + 8, true);
                this.exited = true;
                delete this._inst;
                delete this._values;
                delete this._goRefCounts;
                delete this._ids;
                delete this._idPool;
                this.exit(code);
              },

              "runtime.wasmWrite": (sp) => {
                sp >>>= 0;
                const fd = getInt64(sp + 8);
                const p = getInt64(sp + 16);
                const n = this.mem.getInt32(sp + 24, true);
                fs.writeSync(
                  fd,
                  new Uint8Array(this._inst.exports.mem.buffer, p, n)
                );
              },

              "runtime.resetMemoryDataView": (sp) => {
                sp >>>= 0;
                this.mem = new DataView(this._inst.exports.mem.buffer);
              },

              "runtime.nanotime1": (sp) => {
                sp >>>= 0;
                setInt64(
                  sp + 8,
                  (timeOrigin + performance.now()) * 1000000
                );
              },

              "runtime.walltime": (sp) => {
                sp >>>= 0;
                const msec = new Date().getTime();
                setInt64(sp + 8, msec / 1000);
                this.mem.setInt32(sp + 16, (msec % 1000) * 1000000, true);
              },

              "runtime.scheduleTimeoutEvent": (sp) => {
                sp >>>= 0;
                const id = this._nextCallbackTimeoutID;
                this._nextCallbackTimeoutID++;
                this._scheduledTimeouts.set(
                  id,
                  setTimeout(() => {
                    this._resume();
                    while (this._scheduledTimeouts.has(id)) {
                      console.warn(
                        "scheduleTimeoutEvent: missed timeout event"
                      );
                      this._resume();
                    }
                  }, getInt64(sp + 8) + 1)
                );
                this.mem.setInt32(sp + 16, id, true);
              },

              "runtime.clearTimeoutEvent": (sp) => {
                sp >>>= 0;
                const id = this.mem.getInt32(sp + 8, true);
                clearTimeout(this._scheduledTimeouts.get(id));
                this._scheduledTimeouts.delete(id);
              },

              "runtime.getRandomData": (sp) => {
                sp >>>= 0;
                crypto.getRandomValues(loadSlice(sp + 8));
              },

              "syscall/js.finalizeRef": (sp) => {
                sp >>>= 0;
                const id = this.mem.getUint32(sp + 8, true);
                this._goRefCounts[id]--;
                if (this._goRefCounts[id] === 0) {
                  const v = this._values[id];
                  this._values[id] = null;
                  this._ids.delete(v);
                  this._idPool.push(id);
                }
              },

              "syscall/js.stringVal": (sp) => {
                sp >>>= 0;
                storeValue(sp + 24, loadString(sp + 8));
              },

              "syscall/js.valueGet": (sp) => {
                sp >>>= 0;
                const result = Reflect.get(
                  loadValue(sp + 8),
                  loadString(sp + 16)
                );
                sp = this._inst.exports.getsp() >>> 0;
                storeValue(sp + 32, result);
              },

              "syscall/js.valueSet": (sp) => {
                sp >>>= 0;
                Reflect.set(
                  loadValue(sp + 8),
                  loadString(sp + 16),
                  loadValue(sp + 32)
                );
              },

              "syscall/js.valueDelete": (sp) => {
                sp >>>= 0;
                Reflect.deleteProperty(
                  loadValue(sp + 8),
                  loadString(sp + 16)
                );
              },

              "syscall/js.valueIndex": (sp) => {
                sp >>>= 0;
                storeValue(
                  sp + 24,
                  Reflect.get(loadValue(sp + 8), getInt64(sp + 16))
                );
              },

              "syscall/js.valueSetIndex": (sp) => {
                sp >>>= 0;
                Reflect.set(
                  loadValue(sp + 8),
                  getInt64(sp + 16),
                  loadValue(sp + 24)
                );
              },

              "syscall/js.valueCall": (sp) => {
                sp >>>= 0;
                try {
                  const v = loadValue(sp + 8);
                  const m = Reflect.get(v, loadString(sp + 16));
                  const args = loadSliceOfValues(sp + 32);
                  const result = Reflect.apply(m, v, args);
                  sp = this._inst.exports.getsp() >>> 0;
                  storeValue(sp + 56, result);
                  this.mem.setUint8(sp + 64, 1);
                } catch (err) {
                  sp = this._inst.exports.getsp() >>> 0;
                  storeValue(sp + 56, err);
                  this.mem.setUint8(sp + 64, 0);
                }
              },

              "syscall/js.valueInvoke": (sp) => {
                sp >>>= 0;
                try {
                  const v = loadValue(sp + 8);
                  const args = loadSliceOfValues(sp + 16);
                  const result = Reflect.apply(v, undefined, args);
                  sp = this._inst.exports.getsp() >>> 0;
                  storeValue(sp + 40, result);
                  this.mem.setUint8(sp + 48, 1);
                } catch (err) {
                  sp = this._inst.exports.getsp() >>> 0;
                  storeValue(sp + 40, err);
                  this.mem.setUint8(sp + 48, 0);
                }
              },

              "syscall/js.valueNew": (sp) => {
                sp >>>= 0;
                try {
                  const v = loadValue(sp + 8);
                  const args = loadSliceOfValues(sp + 16);
                  const result = Reflect.construct(v, args);
                  sp = this._inst.exports.getsp() >>> 0;
                  storeValue(sp + 40, result);
                  this.mem.setUint8(sp + 48, 1);
                } catch (err) {
                  sp = this._inst.exports.getsp() >>> 0;
                  storeValue(sp + 40, err);
                  this.mem.setUint8(sp + 48, 0);
                }
              },

              "syscall/js.valueLength": (sp) => {
                sp >>>= 0;
                setInt64(sp + 16, parseInt(loadValue(sp + 8).length));
              },

              "syscall/js.valuePrepareString": (sp) => {
                sp >>>= 0;
                const str = encoder.encode(String(loadValue(sp + 8)));
                storeValue(sp + 16, str);
                setInt64(sp + 24, str.length);
              },

              "syscall/js.valueLoadString": (sp) => {
                sp >>>= 0;
                const str = loadValue(sp + 8);
                loadSlice(sp + 16).set(str);
              },

              "syscall/js.valueInstanceOf": (sp) => {
                sp >>>= 0;
                this.mem.setUint8(
                  sp + 24,
                  loadValue(sp + 8) instanceof loadValue(sp + 16) ? 1 : 0
                );
              },

              "syscall/js.copyBytesToGo": (sp) => {
                sp >>>= 0;
                const dst = loadSlice(sp + 8);
                const src = loadValue(sp + 32);
                if (
                  !(
                    src instanceof Uint8Array ||
                    src instanceof Uint8ClampedArray
                  )
                ) {
                  this.mem.setUint8(sp + 48, 0);
                  return;
                }
                const toCopy = src.subarray(0, dst.length);
                dst.set(toCopy);
                setInt64(sp + 40, toCopy.length);
                this.mem.setUint8(sp + 48, 1);
              },

              "syscall/js.copyBytesToJS": (sp) => {
                sp >>>= 0;
                const dst = loadValue(sp + 8);
                const src = loadSlice(sp + 16);
                if (
                  !(
                    dst instanceof Uint8Array ||
                    dst instanceof Uint8ClampedArray
                  )
                ) {
                  this.mem.setUint8(sp + 48, 0);
                  return;
                }
                const toCopy = src.subarray(0, dst.length);
                dst.set(toCopy);
                setInt64(sp + 40, toCopy.length);
                this.mem.setUint8(sp + 48, 1);
              },

              debug: (value) => {
                console.log(value);
              },
            },
          };
        }

        async run(instance) {
          if (!(instance instanceof WebAssembly.Instance)) {
            throw new Error("Go.run: WebAssembly.Instance expected");
          }
          this._inst = instance;
          this.mem = new DataView(this._inst.exports.mem.buffer);
          this._values = [NaN, 0, null, true, false, globalThis, this];
          this._goRefCounts = new Array(this._values.length).fill(Infinity);
          this._ids = new Map([
            [0, 1],
            [null, 2],
            [true, 3],
            [false, 4],
            [globalThis, 5],
            [this, 6],
          ]);
          this._idPool = [];
          this.exited = false;

          let offset = 4096;

          const strPtr = (str) => {
            const ptr = offset;
            const bytes = encoder.encode(str + "\0");
            new Uint8Array(this.mem.buffer, offset, bytes.length).set(bytes);
            offset += bytes.length;
            if (offset % 8 !== 0) {
              offset += 8 - (offset % 8);
            }
            return ptr;
          };

          const argc = this.argv.length;

          const argvPtrs = [];
          this.argv.forEach((arg) => {
            argvPtrs.push(strPtr(arg));
          });
          argvPtrs.push(0);

          const keys = Object.keys(this.env).sort();
          keys.forEach((key) => {
            argvPtrs.push(strPtr(key + "=" + this.env[key]));
          });
          argvPtrs.push(0);

          const argv = offset;
          argvPtrs.forEach((ptr) => {
            this.mem.setUint32(offset, ptr, true);
            this.mem.setUint32(offset + 4, 0, true);
            offset += 8;
          });

          const wasmMinDataAddr = 4096 + 8192;
          if (offset >= wasmMinDataAddr) {
            throw new Error(
              "total length of command line and environment variables exceeds limit"
            );
          }

          this._inst.exports.run(argc, argv);
          if (this.exited) {
            this._resolveExitPromise();
          }
          await this._exitPromise;
        }

        _resume() {
          if (this.exited) {
            throw new Error("Go program has already exited");
          }
          this._inst.exports.resume();
          if (this.exited) {
            this._resolveExitPromise();
          }
        }

        _makeFuncWrapper(id) {
          const go = this;
          return function () {
            const event = { id: id, this: this, args: arguments };
            go._pendingEvent = event;
            go._resume();
            return event.result;
          };
        }
      };
    })();
    // END wasm_exec.js
  }

  // === SECTION 3: WASM Loader ===
  var wasmPromise = null;

  function resolveWasmURL() {
    if (config.baseURL) {
      var base = config.baseURL;
      if (base.charAt(base.length - 1) !== "/") {
        base += "/";
      }
      return base + config.wasmFile;
    }
    // Derive from script src attribute
    if (currentScript && currentScript.src) {
      var parts = currentScript.src.split("/");
      parts.pop(); // remove filename
      return parts.join("/") + "/" + config.wasmFile;
    }
    // Fallback: relative to current page
    return config.wasmFile;
  }

  function loadWasm() {
    if (wasmPromise) {
      return wasmPromise;
    }

    // If generateSVG is already available (e.g., playground loaded it), skip
    if (typeof globalThis.generateSVG === "function") {
      wasmPromise = Promise.resolve();
      return wasmPromise;
    }

    var wasmURL = resolveWasmURL();
    var go = new Go();

    wasmPromise = WebAssembly.instantiateStreaming(
      fetch(wasmURL),
      go.importObject
    )
      .then(function (result) {
        go.run(result.instance);
      })
      .catch(function (err) {
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
