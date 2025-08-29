package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// handleMapRequest handles HTTP requests for the /map endpoint
func handleMapRequest(w http.ResponseWriter, r *http.Request) {
	// Get the base64 encoded map data from query parameter
	encodedData := r.URL.Query().Get("wardley_map_json_base64")
	if encodedData == "" {
		http.Error(w, "Missing wardley_map_json_base64 query parameter", http.StatusBadRequest)
		return
	}

	// Get the output format (default to SVG)
	outputFormat := r.URL.Query().Get("output")
	if outputFormat == "" {
		outputFormat = "svg"
	}

	// Validate output format
	if outputFormat != "svg" && outputFormat != "json" {
		http.Error(w, "Invalid output format. Must be 'svg' or 'json'", http.StatusBadRequest)
		return
	}

	// Decode the map data
	m, err := decodeMapFromGzippedBase64(encodedData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode map data: %v", err), http.StatusBadRequest)
		return
	}

	// Generate output in requested format
	content, err := generateOutput(m, outputFormat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate output: %v", err), http.StatusInternalServerError)
		return
	}

	// Set appropriate content type
	switch outputFormat {
	case "json":
		w.Header().Set("Content-Type", "application/json")
	case "svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	}

	// Allow CORS for web applications
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle OPTIONS request for CORS preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Write the content
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}

// startWebServer starts the HTTP server for serving Wardley maps
func startWebServer() {
	port, err := strconv.Atoi(uriPort)
	if err != nil {
		log.Printf("Invalid port number for web server: %s, web server disabled", uriPort)
		return
	}

	addr := fmt.Sprintf("%s:%d", uriHost, port)

	// Set up routes
	http.HandleFunc("/map", handleMapRequest)

	// Add a simple health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"wardley-map-server"}`))
	})

	// Add a root endpoint that shows usage
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		usage := `<!DOCTYPE html>
<html>
<head>
    <title>Wardley Map Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        code { background-color: #f4f4f4; padding: 2px 4px; border-radius: 3px; }
        pre { background-color: #f4f4f4; padding: 10px; border-radius: 5px; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>Wardley Map Server</h1>
    <p>This server renders Wardley maps from base64-encoded, gzipped JSON data.</p>
    
    <h2>Usage</h2>
    <p>Send a GET request to:</p>
    <pre><code>/map?wardley_map_json_base64=&lt;encoded_data&gt;&amp;output=&lt;format&gt;</code></pre>
    
    <h3>Parameters:</h3>
    <ul>
        <li><strong>wardley_map_json_base64</strong> (required): Base64-encoded, gzipped JSON representation of the Wardley map</li>
        <li><strong>output</strong> (optional): Output format - "svg" (default) or "json"</li>
    </ul>
    
    <h3>Examples:</h3>
    <ul>
        <li><code>/map?wardley_map_json_base64=&lt;data&gt;</code> - Returns SVG</li>
        <li><code>/map?wardley_map_json_base64=&lt;data&gt;&amp;output=svg</code> - Returns SVG</li>
        <li><code>/map?wardley_map_json_base64=&lt;data&gt;&amp;output=json</code> - Returns JSON</li>
    </ul>
    
    <h3>Health Check:</h3>
    <p><a href="/health">/health</a> - Returns server status</p>
</body>
</html>`
		w.Write([]byte(usage))
	})

	log.Printf("Starting Wardley Map web server on %s", addr)
	log.Printf("Server endpoints:")
	log.Printf("  GET /map?wardley_map_json_base64=<data>&output=<format>")
	log.Printf("  GET /health")
	log.Printf("  GET / (usage information)")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("Web server failed to start: %v", err)
	}
}
