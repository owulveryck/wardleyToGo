package main

import (
	"net/http"
	"path/filepath"
)

type SVGHandler struct {
	maps storage
}

func (s *SVGHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	svgFile := filepath.Base(r.URL.Path)
	content, err := s.maps.get(svgFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if content == nil {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(content)
}

type storage interface {
	save(key string, value []byte) error
	get(key string) ([]byte, error)
}

type memoryStorage map[string][]byte

func (m memoryStorage) save(key string, value []byte) error {
	m[key] = value
	return nil
}

func (m memoryStorage) get(key string) ([]byte, error) {
	if val, ok := m[key]; ok {
		return val, nil
	}
	return nil, nil
}
