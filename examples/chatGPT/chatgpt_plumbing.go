package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type ChatGPTPlumbing struct {
	aiPlugin        *AIPlugin
	aiPluginPayload []byte
	openAPIFile     string
	openAPIContent  []byte
}

// func (chatgptplumbing *ChatGPTPlumbing) GetMux() *http.ServeMux {
func (chatgptplumbing *ChatGPTPlumbing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/ai-plugin.json", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "%s", chatgptplumbing.aiPluginPayload)
	})
	mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "%s", chatgptplumbing.openAPIContent)

	})
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, _ *http.Request) {
		generateLogo(w)
	})

	//return mux
	mux.ServeHTTP(w, r)
}

func NewChatGPTPlumbing(aiPlugin *AIPlugin) (*ChatGPTPlumbing, error) {
	chatGPTPlumbing := &ChatGPTPlumbing{
		aiPlugin: aiPlugin,
	}
	var err error
	chatGPTPlumbing.aiPluginPayload, err = json.MarshalIndent(aiPlugin, " ", " ")
	if err != nil {
		return nil, err
	}

	openapiTmplContent, err := ioutil.ReadFile("openapi.tmpl")
	if err != nil {
		return nil, err
	}
	openAPITmpl, err := template.New("openapi").Parse(string(openapiTmplContent))
	if err != nil {
		return nil, err
	}
	var openAPI bytes.Buffer
	err = openAPITmpl.Execute(&openAPI, aiPlugin)
	if err != nil {
		return nil, err
	}
	chatGPTPlumbing.openAPIContent = openAPI.Bytes()
	return chatGPTPlumbing, nil
}
