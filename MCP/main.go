package main

import (
	"flag"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Parse command line flags
	disableWeb := flag.Bool("no-web", false, "Disable web server (MCP server only)")
	enablePrompts := flag.Bool("prompt", false, "Enable prompt capabilities in MCP server")
	flag.Parse()

	// Start the web server in a goroutine unless disabled
	if !*disableWeb {
		go startWebServer()
	}

	// Create a new MCP server
	/*
		🤖 AI WORKFLOW GUIDE 🤖

		This MCP server provides tools for creating and editing Wardley Maps with a complete workflow designed for AI agents.
		The server also runs a web server (localhost:8585) that serves maps from shareable URIs.

		📋 CORE WORKFLOW PRINCIPLES:
		1. ALWAYS use 'json' output when you plan to make more changes (intermediate operations)
		2. Use 'uri' output to create shareable links that display SVG maps on the web server
		3. Use 'svg' output only for final display when no more changes are needed
		4. Each step should return JSON that is understandable and can be passed to the next tool

		🔄 COMPLETE WORKFLOWS:

		A) NEW MAP CREATION:
		   1. create_map(output="json") → Get starting JSON
		   2. add_elements(map_json=..., output="json") → Add components and anchors
		   3. add_links(map_json=..., output="json") → Add dependencies
		   4. auto_layout(map_json=..., output="json") → Position elements (optional)
		   5. Final step: Use output="uri" to create shareable link

		B) EDITING EXISTING MAP:
		   1. decode_uri(uri="...") → Extract JSON from shareable URI
		   2. Modify using any tools with output="json"
		   3. Final step: Use output="uri" to create new shareable link

		C) QUICK DISPLAY:
		   - Any tool with output="svg" for immediate visualization
		   - Visit the URI in a browser to see the interactive map

		🎯 KEY TOOLS BY PURPOSE:
		- 🚀 CREATE: create_map
		- 🔧 ADD: add_elements, add_links
		- 📐 POSITION: move_elements, auto_layout
		- 🎨 CONFIGURE: configure_evolution
		- 🔄 CONVERT: decode_uri
		- 🗑️ REMOVE: remove_elements

		⚠️ IMPORTANT: The JSON format is the universal interchange format. All tools understand it.
		   URIs contain compressed JSON and can be decoded back to JSON for further editing.
	*/
	s := server.NewMCPServer(
		"Wardley Map Generator & Web Server 🗺️",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// 📐 POSITION: move_elements tool
	moveElementsTool := mcp.NewTool("move_elements",
		mcp.WithDescription("📐 POSITION: Reposition specific elements by name. Handles single or multiple moves. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("moves",
			mcp.Required(),
			mcp.Description("JSON array of move operations. Each move must have: name (string, element name), x (0-100, evolution: 0=genesis/left, 100=commodity/right), y (0-100, visibility: 0=visible/bottom, 100=invisible/top). Example: [{'name':'Customer','x':20,'y':15}, {'name':'Service','x':60,'y':50}]"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 🚀 CREATE: create_map tool
	createMapTool := mcp.NewTool("create_map",
		mcp.WithDescription("🚀 CREATE: Create a new empty Wardley map. Starting point for all workflows. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("title",
			mcp.Description("Title for the map (default: 'New Wardley Map')"),
		),
		mcp.WithNumber("map_id",
			mcp.Description("Unique identifier for the map (default: 1). Use different IDs to create multiple independent maps."),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 🔧 ADD: add_elements tool (unified components and anchors)
	addElementsTool := mcp.NewTool("add_elements",
		mcp.WithDescription("🔧 ADD: Add or update components and anchors. Handles single or multiple elements. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("elements",
			mcp.Required(),
			mcp.Description("JSON array of elements. Each element must have: name (string), x (0-100, evolution: 0=genesis/left, 100=commodity/right), y (0-100, visibility: 0=visible/bottom, 100=invisible/top), element_type ('component' or 'anchor', default: 'component'), type (optional for components: 'regular', 'build', 'buy', 'outsource', 'dataproduct'). Example: [{'name':'Customer','x':15,'y':15,'element_type':'anchor'}, {'name':'Service','x':50,'y':50,'element_type':'component','type':'regular'}]"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 🔧 ADD: add_links tool
	addLinksTool := mcp.NewTool("add_links",
		mcp.WithDescription("🔧 ADD: Add dependency relationships between elements. Handles single or multiple links. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("links",
			mcp.Required(),
			mcp.Description("JSON array of dependency links. Each link must have: from (string, source element name), to (string, target element name), type (optional: 'regular', 'evolved_component', 'evolved'). The 'from' element depends on the 'to' element. Example: [{'from':'Customer','to':'Service','type':'regular'}]"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 🎨 CONFIGURE: configure_evolution tool
	configureEvolutionTool := mcp.NewTool("configure_evolution",
		mcp.WithDescription("🎨 CONFIGURE: Customize the four evolution stage labels displayed on the X-axis. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("stage1",
			mcp.Required(),
			mcp.Description("Label for stage 1 (Genesis/Concept): Novel, experimental, uncertain"),
		),
		mcp.WithString("stage2",
			mcp.Required(),
			mcp.Description("Label for stage 2 (Custom/Emerging): Bespoke, tailored solutions"),
		),
		mcp.WithString("stage3",
			mcp.Required(),
			mcp.Description("Label for stage 3 (Product/Converging): Packaged, feature-complete products"),
		),
		mcp.WithString("stage4",
			mcp.Required(),
			mcp.Description("Label for stage 4 (Commodity/Accepted): Most evolved, standardized, utility-like"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 📐 POSITION: auto_layout tool
	autoLayoutTool := mcp.NewTool("auto_layout",
		mcp.WithDescription("📐 POSITION: Automatically arrange all elements based on value chain depth analysis. Calculates dependency paths and positions elements in horizontal layers. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 🔄 CONVERT: decode_uri tool
	decodeUriTool := mcp.NewTool("decode_uri",
		mcp.WithDescription("🔄 CONVERT: Extract map JSON from a shareable URI for editing. Essential for modifying existing maps. Returns JSON that can be passed to other tools."),
		mcp.WithString("uri",
			mcp.Required(),
			mcp.Description("The complete URI containing the base64-encoded map data (e.g., 'http://localhost:8585/map?wardley_map_json_base64=...')"),
		),
	)

	// 🗑️ REMOVE: remove_elements tool
	removeElementsTool := mcp.NewTool("remove_elements",
		mcp.WithDescription("🗑️ REMOVE: Remove components, anchors, or links from the map. Handles single or multiple removals. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("remove_type",
			mcp.Required(),
			mcp.Description("Type of removal: 'elements' to remove components/anchors, 'links' to remove connections"),
		),
		mcp.WithString("items",
			mcp.Required(),
			mcp.Description("JSON array of items to remove. For elements: [{'name':'ElementName'}]. For links: [{'from':'SourceName','to':'TargetName'}]. Example: [{'name':'Customer'}] or [{'from':'Customer','to':'Service'}]"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// Add tool handlers
	s.AddTool(createMapTool, createMapHandler)
	s.AddTool(addElementsTool, addElementsHandler)
	s.AddTool(addLinksTool, addLinksHandler)
	s.AddTool(moveElementsTool, moveElementsHandler)
	s.AddTool(autoLayoutTool, autoLayoutHandler)
	s.AddTool(configureEvolutionTool, configureEvolutionHandler)
	s.AddTool(decodeUriTool, decodeUriHandler)
	s.AddTool(removeElementsTool, removeElementsHandler)

	// Add workflow prompts only if -prompt flag is set
	if *enablePrompts {
		addWorkflowPrompts(s)
	}

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
