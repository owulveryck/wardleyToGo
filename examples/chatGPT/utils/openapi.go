package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/openapi"
)

func generateOpenAPI(defFile string, config *load.Config) ([]byte, error) {
	buildInstances := load.Instances([]string{defFile}, config)
	insts := cue.Build(buildInstances)
	b, err := openapi.Gen(insts[0], nil)
	if err != nil {
		return nil, err
	}
	var s spec
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}

	s.OpenAPIVersion = "3.0.1"
	s.Info = info{
		Title:       "Wardley To Go",
		Description: "This plugin draw wardley maps",
		Version:     "v1",
	}
	s.Servers = []server{{
		URL: "http://localhost:3333",
	}}
	s.Paths = paths{
		"/api/map": path{
			Post: post{
				Summary:     "create a complete map with a set of component linked on a value chain and each component placed on its evolution phase",
				OperationID: "map",
				RequestBody: requestBody{
					Required: true,
					Content: map[string]contentBody{
						"application/json": {
							Schema: schema{
								Ref: "#/components/schemas/mapRequest",
							},
						},
					},
				},
				Responses: map[int]response{
					200: {
						Description: "Map URL",
						Content: map[string]responseContent{
							"application/json": {
								Schema: schema{
									Ref: "#/components/schemas/mapURL",
								},
							},
						},
					},
				},
			},
		},
		"/api/mapEvolution": path{
			Post: post{
				Summary:     "create a map with a single component on a specific evolution phase",
				OperationID: "evolutionMap",
				RequestBody: requestBody{
					Required: true,
					Content: map[string]contentBody{
						"application/json": {
							Schema: schema{
								Ref: "#/components/schemas/evolutionRequest",
							},
						},
					},
				},
				Responses: map[int]response{
					200: {
						Description: "Map URL",
						Content: map[string]responseContent{
							"application/json": {
								Schema: schema{
									Ref: "#/components/schemas/mapURL",
								},
							},
						},
					},
				},
			},
		},
		"/api/mapValueChain": path{
			Post: post{
				Summary:     "create a simple map with a value chain for a set of components",
				OperationID: "valueChainMap",
				RequestBody: requestBody{
					Required: true,
					Content: map[string]contentBody{
						"application/json": {
							Schema: schema{
								Ref: "#/components/schemas/valueChainRequest",
							},
						},
					},
				},
				Responses: map[int]response{
					200: {
						Description: "Map URL",
						Content: map[string]responseContent{
							"application/json": {
								Schema: schema{
									Ref: "#/components/schemas/mapURL",
								},
							},
						},
					},
				},
			},
		},
	}
	return json.MarshalIndent(s, "", "    ")
}

func main() {
	src := bytes.NewBufferString(`
	#mapURL: {
		imageURL: string
	}
	#evolutionRequest: {
		// The component to add to the map
		component: string
		// stages 4 entries to describe each stage of eveolution depending on the type of the component.
		// If the component is an asset or an activity, the array can be: [genesis,custom,product,commodity]
		// If the component is a practice, the array can be: [novel,emerging,good,best]
		// If the component is a general knowledge, the array can be: [concent,hypothesis,theory,accepted]
		// If the component is some data, the array can be: [unmodeled,divergent,convergent,modeled]
		stages: [...string]
		// evolution is the position on the evolution axis between 0 and 100. 
		// From 0 to 17 the compoenent is in stage 1 (genesis for an asset or a an activity, novel for a practice, concept for some general knowledge)
		// From 18 to 40 the component is in stage 2 (custom for an asset or an activity, emerging for a practice, hypothesis for some general knowledge) 
		// From 40 to 70 the component is in stage 3 (product for and asset or an activity, good for a practice, or theory for some general knowledge)
		// From 70 to 99 the component is in stage 4 (commodity for an asset of an activity, best for a practice, accepted for some general knwoledge)
		evolution: number & < 100
	}
	// the value chain in a set of couples
	// a couple is a relation between a component and one of its dependency
	// a couple may contain information about the visibility of the dependency from the point of view of the component
	#valueChainRequest: {
		couple: [...#link]
	}
	#link: {
		// the component that has a requirement
		component: string
		// the dependency allowing the source component to fulfill its need
		dependency: string
		// visibility is how important the dependency is for the component. 
		// the greater the less important the dependency is
		visibility: number & < 5 & > 1 | *1
	}
	// #mapRequest is the representation of a map with the value chain and the evolution for each component
	#mapRequest: {
		// the valueChain linking the components together
		valueChain: #valueChainRequest
		// for each component: its evolution phase.
		components: [...#evolutionRequest]

	}
	`)
	b, err := generateOpenAPI("-", &load.Config{
		Stdin: src,
	})
	if err != nil {
		log.Fatal(err)
	}
	// This contains the OpenAPI.json definition
	fmt.Println(string(b))
}

type info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type server struct {
	URL string `json:"url"`
}

type spec struct {
	OpenAPIVersion string      `json:"openapi"`
	Info           info        `json:"info"`
	Servers        []server    `json:"servers"`
	Paths          paths       `json:"paths"`
	Components     interface{} `json:"components"`
}

type requestBody struct {
	Required bool                   `json:"required"`
	Content  map[string]contentBody `json:"content"`
}

type contentBody struct {
	Schema schema `json:"schema"`
}

type schema struct {
	Ref string `json:"$ref"`
}

type response struct {
	Description string                     `json:"description"`
	Content     map[string]responseContent `json:"content"`
}

type responseContent struct {
	Schema schema `json:"schema"`
}

type schemaResponse struct {
	Type   string `json:"type"`
	Format string `json:"format"`
}

type post struct {
	Summary     string           `json:"summary"`
	OperationID string           `json:"operationId"`
	RequestBody requestBody      `json:"requestBody"`
	Responses   map[int]response `json:"responses"`
}

type path struct {
	Post post `json:"post"`
}

type paths map[string]path
