package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// addWorkflowPrompts adds simple prompts that guide AI to use MCP tools
func addWorkflowPrompts(s *server.MCPServer) {
	// Simple map creation prompt
	s.AddPrompt(mcp.NewPrompt("create_wardley_map",
		mcp.WithPromptDescription("🗺️ Step-by-step guide for creating Wardley Maps using MCP tools"),
		mcp.WithArgument("title", mcp.ArgumentDescription("Map title")),
		mcp.WithArgument("components", mcp.ArgumentDescription("Comma-separated components")),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		title := request.Params.Arguments["title"]
		components := request.Params.Arguments["components"]

		// Always provide defaults - never return nil
		if title == "" {
			title = "New Wardley Map"
		}
		if components == "" {
			components = "User Interface, Service, Database"
		}

		guide := fmt.Sprintf(`I'll guide you through creating "%s" with components: %s

Execute these MCP tool calls in order:

**STEP 1:** get_empty_map
- title: "%s"
- output: "json"

**STEP 2:** add_anchor
- map_json: [result from step 1]
- anchor_name: "User"
- x: 15, y: 15
- output: "json"

**STEP 3:** add_components
- map_json: [result from step 2]
- components: [array from: %s]
- output: "json"

**STEP 4:** add_links
- map_json: [result from step 3]
- links: [dependency array]
- output: "json"

**STEP 5:** auto_value_chain (REQUIRED)
- map_json: [result from step 4]
- output: "json"

**STEP 6:** auto_value_chain (final URI)
- map_json: [result from step 5]
- output: "uri"

Start with step 1!`, title, components, title, components)

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Create: %s", title),
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleAssistant, mcp.NewTextContent(guide)),
			},
		), nil
	})

	// Simple map editing prompt
	s.AddPrompt(mcp.NewPrompt("edit_wardley_map",
		mcp.WithPromptDescription("✏️ Step-by-step guide for editing Wardley Maps using MCP tools"),
		mcp.WithArgument("uri", mcp.ArgumentDescription("Existing map URI")),
		mcp.WithArgument("changes", mcp.ArgumentDescription("Changes to make")),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		uri := request.Params.Arguments["uri"]
		changes := request.Params.Arguments["changes"]

		// Always provide defaults - never return nil
		if uri == "" {
			uri = "http://localhost:8585/map?wardley_map_json_base64=..."
		}
		if changes == "" {
			changes = "modify the map"
		}

		guide := fmt.Sprintf(`I'll help you edit: %s
Changes: %s

Execute these MCP tool calls in order:

**STEP 1:** decode_map_from_uri
- uri: "%s"

**STEP 2:** Apply changes using:
- add_components (new components)
- add_links (new dependencies)
- move_component (reposition)
- unlink_components (remove links)
(Use output="json" for all)

**STEP 3:** auto_value_chain (REQUIRED)
- map_json: [result from step 2]
- output: "json"

**STEP 4:** auto_value_chain (final URI)
- map_json: [result from step 3]
- output: "uri"

Start with step 1!`, uri, changes, uri)

		return mcp.NewGetPromptResult(
			"Edit Wardley Map",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleAssistant, mcp.NewTextContent(guide)),
			},
		), nil
	})
}
