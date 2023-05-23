package main

import (
	"net/http"
	"path/filepath"
)

type SVGHandler struct {
	maps map[string][]byte
}

func (s *SVGHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	svgFile := filepath.Base(r.URL.Path)
	if content, ok := s.maps[svgFile]; ok {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(content)
		return
	}
	http.Error(w, "image not found", http.StatusNotFound)
}
