// Package server implements the MCP server functionality for Elasticsearch integration.
// It provides both stdio and Streamable HTTP protocol support for the Model Context Protocol.
package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/AeaZer/mcp-elasticsearch/config"
	"github.com/AeaZer/mcp-elasticsearch/elasticsearch"
	"github.com/AeaZer/mcp-elasticsearch/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ElasticsearchMCPServer represents the main MCP server for Elasticsearch operations
type ElasticsearchMCPServer struct {
	config    *config.Config            // Server configuration
	esClient  elasticsearch.Client      // Elasticsearch client instance
	esTools   *tools.ElasticsearchTools // Elasticsearch tools collection
	mcpServer *server.MCPServer         // Underlying MCP server
}

// NewElasticsearchMCPServer creates a new instance of the Elasticsearch MCP server
// with the provided configuration. It initializes the Elasticsearch client,
// creates the tools collection, and sets up the MCP server with all tools registered.
func NewElasticsearchMCPServer(cfg *config.Config) (*ElasticsearchMCPServer, error) {
	// Create Elasticsearch client with the provided configuration
	esClient, err := elasticsearch.NewClient(&cfg.Elasticsearch, cfg.GetElasticsearchVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// Create the tools collection with the Elasticsearch client
	esTools := tools.NewElasticsearchTools(esClient)

	// Create the MCP server with tool capabilities enabled
	mcpServer := server.NewMCPServer(
		cfg.Server.Name,
		cfg.Server.Version,
		server.WithToolCapabilities(true),
	)

	// Register all Elasticsearch tools with the MCP server
	if err := registerTools(mcpServer, esTools); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	return &ElasticsearchMCPServer{
		config:    cfg,
		esClient:  esClient,
		esTools:   esTools,
		mcpServer: mcpServer,
	}, nil
}

// registerTools registers all Elasticsearch tools with the MCP server.
// Each tool is registered with its corresponding handler function.
func registerTools(mcpServer *server.MCPServer, esTools *tools.ElasticsearchTools) error {
	toolsList := esTools.GetTools()

	for _, tool := range toolsList {
		// Create a closure to capture the tool name for the handler
		toolName := tool.Name
		mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return esTools.HandleTool(ctx, toolName, request.GetArguments())
		})
	}

	return nil
}

// Start launches the MCP server using the configured protocol (stdio or http).
// The method blocks until the server is stopped or encounters an error.
func (s *ElasticsearchMCPServer) Start() error {
	switch s.config.Server.Protocol {
	case "stdio":
		return s.startStdioServer()
	case "http":
		return s.startStreamableHTTP()
	default:
		return fmt.Errorf("unsupported protocol: %s", s.config.Server.Protocol)
	}
}

// startStdioServer starts the MCP server using the stdio protocol.
// This method blocks until the server is stopped.
func (s *ElasticsearchMCPServer) startStdioServer() error {
	log.Printf("Starting Elasticsearch MCP Server (stdio protocol)")

	// Start the server using stdio protocol
	if err := server.ServeStdio(s.mcpServer); err != nil {
		return fmt.Errorf("stdio server failed to start: %w", err)
	}

	return nil
}

// startStreamableHTTP starts the MCP server using the HTTP protocol.
// This method blocks until the server is stopped.
func (s *ElasticsearchMCPServer) startStreamableHTTP() error {
	address := fmt.Sprintf("%s:%d", s.config.Server.Address, s.config.Server.Port)
	log.Printf("Starting Elasticsearch MCP Server (HTTP protocol) listening on: %s", address)

	// Create HTTP handler for the MCP server
	handler := server.NewStreamableHTTPServer(s.mcpServer)

	// Create and configure the HTTP server
	httpServer := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	// Start the HTTP server
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server failed to start: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the MCP server and closes the Elasticsearch client.
// It ensures all resources are properly cleaned up.
func (s *ElasticsearchMCPServer) Stop() error {
	log.Printf("Stopping Elasticsearch MCP Server")

	// Close the Elasticsearch client connection
	if err := s.esClient.Close(); err != nil {
		log.Printf("Failed to close Elasticsearch client: %v", err)
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
