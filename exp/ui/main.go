//go:build !js

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	fs := http.FileServer(http.Dir("."))
	handler := cacheHeaders(gzipHandler(fs))
	fmt.Println("Serving on http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("Server error:", err)
	}
}

var gzipWriterPool = sync.Pool{
	New: func() any { w, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed); return w },
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	sniffDone bool
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.sniffDone {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", http.DetectContentType(b))
		}
		w.Header().Del("Content-Length")
		w.sniffDone = true
	}
	return w.Writer.Write(b)
}

func gzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz := gzipWriterPool.Get().(*gzip.Writer)
		gz.Reset(w)
		defer func() {
			_ = gz.Close()
			gzipWriterPool.Put(gz)
		}()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		next.ServeHTTP(&gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}

func cacheHeaders(next http.Handler) http.Handler {
	longCache := map[string]bool{".wasm": true, ".js": true, ".css": true, ".png": true, ".ico": true}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(r.URL.Path)
		if longCache[ext] {
			w.Header().Set("Cache-Control", "public, max-age=86400")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		next.ServeHTTP(w, r)
	})
}
