package main

import "net/http"

type apiHandler struct{}

// func (a apiHandler) GetMux() *http.ServeMux {
func (a apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("/mapEvolution", mapHandler)
	mux.HandleFunc("/map.svg", svgHandler)
	mux.HandleFunc("/", RootHandler)
	mux.ServeHTTP(w, r)
	// return mux
}
