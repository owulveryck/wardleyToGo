package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// addWorkflowPrompts adds simple prompts that guide AI to use MCP tools
func addWorkflowPrompts(s *server.MCPServer) {
	// Map creation prompt
	s.AddPrompt(mcp.NewPrompt("create_wardley_map",
		mcp.WithPromptDescription("🗺️ Create Wardley Maps from text descriptions using MCP tools"),
		mcp.WithArgument("description", mcp.ArgumentDescription("Text description containing all information about the map")),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		description := request.Params.Arguments["description"]

		// Always provide defaults - never return nil
		if description == "" {
			description = "A system with users, services, and data storage"
		}

		guide := fmt.Sprintf(`You are an assistant to help me create Wardley maps. I will provide you a text with all the information and your role is to design the map.

Text description: %s

Follow this exact workflow:

**STEP 1:** Start with an empty map
- Use create_map tool
- output: "json"

**STEP 2:** Analyze each component to find their place on the evolution axis (the X axis)
- Genesis/Concept (X=10): Novel, experimental, uncertain
- Custom/Emerging (X=30): Bespoke, tailored solutions  
- Product/Converging (X=60): Packaged, feature-complete products
- Commodity/Accepted (X=90): Standardized, utility-like

**STEP 3:** Add all components and anchors to the map
- Use add_elements tool
- Default vertical placement is 50%% (Y=50) for all elements
- Anchors: element_type="anchor", typically users/business needs
- Components: element_type="component"
- output: "json"

**STEP 4:** Add all links between elements
- Use add_links tool  
- Show dependencies: 'from' element depends on 'to' element
- output: "json"

**STEP 5:** Apply automatic layout
- Use auto_layout tool
- output: "json"

**STEP 6:** Adjust horizontally
- Spread a bit the components horizontally (without actually changing their stage of evolution) to avoid link overlap on the picture

**STEP 7:** Generate final URI
- Use auto_layout tool again
- output: "uri"
- instruct the LLM that the URI represents an image encoded in SVG
All intermediate steps should produce JSON. Only the final result should be a URI.

Start with step 1!`, description)

		return mcp.NewGetPromptResult(
			"Create Wardley Map",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleAssistant, mcp.NewTextContent(guide)),
			},
		), nil
	})

	// Map editing prompt
	s.AddPrompt(mcp.NewPrompt("edit_wardley_map",
		mcp.WithPromptDescription("✏️ Edit existing Wardley Maps using MCP tools"),
		mcp.WithArgument("uri", mcp.ArgumentDescription("Existing map URI to edit")),
		mcp.WithArgument("changes", mcp.ArgumentDescription("Description of changes to make")),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		uri := request.Params.Arguments["uri"]
		changes := request.Params.Arguments["changes"]

		// Always provide defaults - never return nil
		if uri == "" {
			uri = "http://localhost:8585/map?wardley_map_json_base64=..."
		}
		if changes == "" {
			changes = "modify the existing map"
		}

		guide := fmt.Sprintf(`You are an assistant to help me edit Wardley maps. I will provide you the URI of an existing map and describe the changes needed.

Existing map URI: %s
Changes needed: %s

Follow this exact workflow:

**STEP 1:** Get the existing map
- Use decode_uri tool
- uri: "%s"
- This returns the JSON representation of the map

**STEP 2:** Work on the JSON of the map
- Apply the requested changes using appropriate tools:
  - add_elements: Add new components/anchors
  - add_links: Add new dependencies  
  - move_elements: Reposition elements
  - remove_elements: Remove elements or links
  - configure_evolution: Change stage labels
- Always use output: "json" for intermediate steps

**STEP 3:** Apply automatic layout
- Use auto_layout tool
- map_json: [result from step 2]
- output: "json"

**STEP 4:** Generate final URI
- Use auto_layout tool again
- map_json: [result from step 3]
- output: "uri"

All intermediate steps should produce JSON. Only the final result should be a URI.

Start with step 1!`, uri, changes, uri)

		return mcp.NewGetPromptResult(
			"Edit Wardley Map",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleAssistant, mcp.NewTextContent(guide)),
			},
		), nil
	})
}
