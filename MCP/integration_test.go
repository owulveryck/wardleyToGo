package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func TestWebServerIntegration(t *testing.T) {
	// Start the web server in a goroutine
	go func() {
		startWebServer()
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Create a test map
	m := wardleyToGo.NewMap(1)
	m.Title = "Integration Test Map"

	comp := wardley.NewComponent(1)
	comp.Label = "Integration Component"
	comp.Placement.X = 40
	comp.Placement.Y = 60
	m.AddComponent(comp)

	// Generate the encoded data
	encodedData, err := encodeMapToGzippedBase64(m)
	if err != nil {
		t.Fatalf("Failed to encode map data: %v", err)
	}

	// Test health endpoint
	t.Run("Health Check", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8585/health")
		if err != nil {
			t.Skipf("Web server not available: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var health map[string]string
		if err := json.Unmarshal(body, &health); err != nil {
			t.Fatalf("Failed to parse health response: %v", err)
		}

		if health["status"] != "ok" {
			t.Errorf("Expected status 'ok', got %s", health["status"])
		}
	})

	// Test root endpoint
	t.Run("Root Endpoint", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8585/")
		if err != nil {
			t.Skipf("Web server not available: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "text/html" {
			t.Errorf("Expected Content-Type text/html, got %s", contentType)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		bodyStr := string(body)
		if !strings.Contains(bodyStr, "Wardley Map Server") {
			t.Error("Root page doesn't contain expected title")
		}
	})

	// Test map endpoint with SVG output
	t.Run("Map SVG Output", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:8585/map?wardley_map_json_base64=%s&output=svg", encodedData)
		resp, err := http.Get(url)
		if err != nil {
			t.Skipf("Web server not available: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "image/svg+xml" {
			t.Errorf("Expected Content-Type image/svg+xml, got %s", contentType)
		}

		// Check CORS headers
		corsOrigin := resp.Header.Get("Access-Control-Allow-Origin")
		if corsOrigin != "*" {
			t.Errorf("Expected CORS origin *, got %s", corsOrigin)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		bodyStr := string(body)
		if !strings.Contains(bodyStr, "<svg") {
			t.Error("SVG response doesn't contain <svg tag")
		}

		if !strings.Contains(bodyStr, "Integration Component") {
			t.Error("SVG doesn't contain the component name")
		}
	})

	// Test map endpoint with JSON output
	t.Run("Map JSON Output", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:8585/map?wardley_map_json_base64=%s&output=json", encodedData)
		resp, err := http.Get(url)
		if err != nil {
			t.Skipf("Web server not available: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var jsonMap JSONMap
		if err := json.Unmarshal(body, &jsonMap); err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}

		if jsonMap.Title != "Integration Test Map" {
			t.Errorf("Expected title 'Integration Test Map', got %s", jsonMap.Title)
		}

		if len(jsonMap.Components) != 1 {
			t.Errorf("Expected 1 component, got %d", len(jsonMap.Components))
		}

		if jsonMap.Components[0].Name != "Integration Component" {
			t.Errorf("Expected component name 'Integration Component', got %s", jsonMap.Components[0].Name)
		}
	})

	// Test map endpoint with default output (should be SVG)
	t.Run("Map Default Output", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:8585/map?wardley_map_json_base64=%s", encodedData)
		resp, err := http.Get(url)
		if err != nil {
			t.Skipf("Web server not available: %v", err)
		}
		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		if contentType != "image/svg+xml" {
			t.Errorf("Expected default Content-Type image/svg+xml, got %s", contentType)
		}
	})

	// Test error cases
	t.Run("Missing Parameter", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8585/map")
		if err != nil {
			t.Skipf("Web server not available: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Output Format", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:8585/map?wardley_map_json_base64=%s&output=invalid", encodedData)
		resp, err := http.Get(url)
		if err != nil {
			t.Skipf("Web server not available: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Base64 Data", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8585/map?wardley_map_json_base64=invalid_base64_data")
		if err != nil {
			t.Skipf("Web server not available: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}
