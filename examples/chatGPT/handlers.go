package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func svgHandler(w http.ResponseWriter, _ *http.Request) {
	content, err := ioutil.ReadFile("sample.svg")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(content)

}

func mapHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("handling content")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(b))
	image := struct {
		ImageURL string
	}{
		ImageURL: "http://localhost:3333/map.svg",
	}
	enc := json.NewEncoder(w)
	enc.Encode(image)

}

func RootHandler(w http.ResponseWriter, _ *http.Request) {
	response := map[string]string{"message": "Welcome!"}
	json.NewEncoder(w).Encode(response)
}
