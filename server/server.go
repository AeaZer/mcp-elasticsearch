// Package server implements the MCP server functionality for Elasticsearch integration.
// It provides both stdio and Streamable HTTP protocol support for the Model Context Protocol.
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AeaZer/mcp-elasticsearch/config"
	"github.com/AeaZer/mcp-elasticsearch/elasticsearch"
	"github.com/AeaZer/mcp-elasticsearch/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ElasticsearchMCPServer represents the main MCP server for Elasticsearch operations
type ElasticsearchMCPServer struct {
	config    *config.Config            // Server configuration
	esClient  elasticsearch.Client      // Elasticsearch client instance
	esTools   *tools.ElasticsearchTools // Elasticsearch tools collection
	mcpServer *mcp.Server               // Underlying MCP server
}

// NewElasticsearchMCPServer creates a new instance of the Elasticsearch MCP server
// with the provided configuration. It initializes the Elasticsearch client,
// creates the tools collection, and sets up the MCP server with all tools registered.
func NewElasticsearchMCPServer(cfg *config.Config) (*ElasticsearchMCPServer, error) {
	// Create Elasticsearch client with the provided configuration
	esClient, err := elasticsearch.NewClient(&cfg.Elasticsearch, cfg.GetElasticsearchVersion())
	if err != nil {
		log.Printf("ERROR: Failed to connect to Elasticsearch: %v", err)
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}
	log.Printf("Connected to Elasticsearch")

	// Create the tools collection with the Elasticsearch client
	esTools := tools.NewElasticsearchTools(esClient)

	// Create the MCP server
	impl := &mcp.Implementation{
		Name:    cfg.Server.Name,
		Version: cfg.Server.Version,
	}
	mcpServer := mcp.NewServer(impl, &mcp.ServerOptions{
		Instructions: "Elasticsearch MCP Server - provides tools for interacting with Elasticsearch clusters",
	})

	// Register all Elasticsearch tools with the MCP server
	if err = registerTools(mcpServer, esTools); err != nil {
		log.Printf("ERROR: Failed to register tools: %v", err)
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}
	log.Printf("Registered %d tools", len(esTools.GetTools()))

	return &ElasticsearchMCPServer{
		config:    cfg,
		esClient:  esClient,
		esTools:   esTools,
		mcpServer: mcpServer,
	}, nil
}

// registerTools registers all Elasticsearch tools with the MCP server.
// Each tool is registered with its corresponding handler function.
func registerTools(mcpServer *mcp.Server, esTools *tools.ElasticsearchTools) error {
	toolsList := esTools.GetTools()

	// Create tool handlers for each tool
	for _, tool := range toolsList {
		// Create a closure to capture the tool name for the handler
		toolName := tool.Name

		handler := func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResult, error) {
			log.Printf("Tool call: %s", toolName)

			result := esTools.HandleTool(ctx, toolName, params.Arguments)

			if result.IsError {
				log.Printf("Tool %s failed", toolName)
			}

			return &result, nil
		}

		// Add tool to the server
		mcpServer.AddTool(&tool, handler)
	}

	return nil
}

// Start launches the MCP server using the configured protocol (stdio, http, or sse).
// The method blocks until the server is stopped or encounters an error.
func (s *ElasticsearchMCPServer) Start() error {
	switch s.config.Server.Protocol {
	case "stdio":
		return s.startStdioServer()
	case "http":
		return s.startStreamableHTTP()
	case "sse":
		log.Printf("WARNING: SSE protocol is deprecated")
		return s.startSSEServer()
	default:
		err := fmt.Errorf("unsupported protocol: %s", s.config.Server.Protocol)
		log.Printf("ERROR: %v", err)
		return err
	}
}

// startStdioServer starts the MCP server using the stdio protocol.
// This mode is typically used for direct integration with LLM tools.
func (s *ElasticsearchMCPServer) startStdioServer() error {
	log.Printf("Starting MCP server on stdio")

	// Create stdio transport
	transport := mcp.NewStdioTransport()

	// Run the server
	err := s.mcpServer.Run(context.Background(), transport)
	if err != nil {
		log.Printf("ERROR: Stdio server failed: %v", err)
	}
	return err
}

// startStreamableHTTP starts the MCP server using the Streamable HTTP protocol.
// This mode allows the server to be accessed over HTTP for remote connections.
func (s *ElasticsearchMCPServer) startStreamableHTTP() error {
	address := fmt.Sprintf("%s:%d", s.config.Server.Address, s.config.Server.Port)

	// Create HTTP handler
	mcpHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		log.Printf("MCP request from %s", r.RemoteAddr)
		return s.mcpServer
	}, nil)

	// Wrap handler with error logging
	loggingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the actual MCP handler
		mcpHandler.ServeHTTP(w, r)
	})

	// Start HTTP server
	http.Handle("/mcp", loggingHandler)

	// Add health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","server":"elasticsearch-mcp"}`))
	})

	log.Printf("HTTP server listening on %s", address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Printf("ERROR: HTTP server failed: %v", err)
	}
	return err
}

// startSSEServer starts the MCP server using the SSE (Server-Sent Events) protocol.
// WARNING: This protocol is deprecated and not recommended for production use.
func (s *ElasticsearchMCPServer) startSSEServer() error {
	address := fmt.Sprintf("%s:%d", s.config.Server.Address, s.config.Server.Port)

	// Create SSE handler
	sseHandler := mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
		log.Printf("SSE request from %s", r.RemoteAddr)
		return s.mcpServer
	})

	// Start HTTP server with SSE endpoints
	http.Handle("/sse", sseHandler)

	log.Printf("SSE server listening on %s", address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Printf("ERROR: SSE server failed: %v", err)
	}
	return err
}

// Stop gracefully shuts down the MCP server and cleans up resources.
func (s *ElasticsearchMCPServer) Stop() error {
	// Close Elasticsearch client connection
	if s.esClient != nil {
		err := s.esClient.Close()
		if err != nil {
			log.Printf("ERROR: Failed to close Elasticsearch client: %v", err)
			return err
		}
	}

	return nil
}

// GetInfo returns information about the server configuration and status.
// This can be useful for debugging and monitoring purposes.
func (s *ElasticsearchMCPServer) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":     s.config.Server.Name,
		"version":  s.config.Server.Version,
		"protocol": s.config.Server.Protocol,
		"elasticsearch": map[string]interface{}{
			"addresses": s.config.Elasticsearch.Addresses,
			"version":   s.config.GetElasticsearchVersion(),
		},
	}
}
