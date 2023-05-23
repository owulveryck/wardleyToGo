package main

import (
	"bytes"
	"encoding/json"
	"image"
	"net/http"

	"github.com/google/uuid"
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
)

type EvolutionInput struct {
	Component string `json:"component"`
	Evolution int    `json:"evolution"`
}

func (a *apiHandler) mapHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var input EvolutionInput
	err := dec.Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create the map
	m := wardleyToGo.NewMap(0)
	c := wardley.NewComponent(1)
	c.Label = input.Component
	c.Placement = image.Point{
		X: input.Evolution,
		Y: 50,
	}
	m.AddComponent(c)

	var buf bytes.Buffer
	// Encode the map
	e, err := svgmap.NewEncoder(&buf, image.Rect(0, 0, 1100, 200), image.Rect(30, 50, 1070, 150))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	style := svgmap.NewOctoStyle(svgmap.DefaultEvolution)
	e.Init(style)
	err = e.Encode(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	e.Close()

	// Save the map
	id := uuid.NewString() + ".svg"
	a.svgHandler.maps[id] = buf.Bytes()

	// Return the map
	result := struct {
		ImageURL string
	}{
		ImageURL: "http://localhost:3333/api/svg/" + id,
	}
	enc := json.NewEncoder(w)
	enc.Encode(result)

}

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	response := map[string]string{"message": "Welcome!"}
	json.NewEncoder(w).Encode(response)
}
