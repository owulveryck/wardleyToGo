package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// RENDER HANDLERS - Tools for rendering maps and converting between formats

func decodeUriHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	uri, err := request.RequireString("uri")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get uri parameter", err), nil
	}

	// Extract base64 data from URI
	encodedData, err := extractBase64FromURI(uri)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to extract base64 data from URI", err), nil
	}

	// Decode the map data
	m, err := decodeMapFromGzippedBase64(encodedData)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to decode map data from URI", err), nil
	}

	// Return JSON representation
	jsonData, err := MarshalMap(m)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal decoded map to JSON", err), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
