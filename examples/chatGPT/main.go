package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

const (
	EnvConfigPrefix = "WTG_CHATGPT"
)

func main() {
	help := flag.Bool("h", false, "help")
	flag.Parse()
	var spec configuration
	err := envconfig.Process(EnvConfigPrefix, &spec)
	if *help || err != nil {
		envconfig.Usage(EnvConfigPrefix, &spec)
		return

	}
	l, err := setupListener(context.Background(), &spec)
	// Close the listener when the application closes.
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	aiPlugin, err := NewAIPlugin(fmt.Sprintf("%v://%v", spec.scheme, spec.ListenAddr))
	if err != nil {
		log.Fatal(err)
	}
	plumbing, err := NewChatGPTPlumbing(aiPlugin)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", &apiHandler{}))
	mux.Handle("/", plumbing)

	fmt.Printf("ingress url: %v://%v\n", spec.scheme, spec.ListenAddr)
	// Configure CORS middleware with allowed origins
	http.Serve(l, corsHandler(mux))
}
