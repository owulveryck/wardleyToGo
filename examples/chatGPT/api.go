package main

import "net/http"

type apiHandler struct {
	address    string
	basePath   string
	svgHandler *SVGHandler
}

// func (a apiHandler) GetMux() *http.ServeMux {
func (a apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("/mapEvolution", a.mapHandler)
	mux.Handle("/svg/", a.svgHandler)
	mux.HandleFunc("/", rootHandler)
	mux.ServeHTTP(w, r)
	// return mux
}
