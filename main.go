// Package main implements the Elasticsearch MCP (Model Context Protocol) server
// that provides an interface for interacting with Elasticsearch clusters
// through both stdio and HTTP protocols.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AeaZer/mcp-elasticsearch/config"
	"github.com/AeaZer/mcp-elasticsearch/server"
)

func main() {
	// Configure logging format with timestamp and file location
	log.SetFlags(log.Ldate | log.Ltime | log.LstdFlags)

	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Log startup information
	log.Printf("Starting Elasticsearch MCP Server...")
	log.Printf("Protocol: %s", cfg.Server.Protocol)
	log.Printf("Elasticsearch addresses: %v", cfg.Elasticsearch.Addresses)
	log.Printf("Elasticsearch version: %s", cfg.GetElasticsearchVersion())

	// Create the Elasticsearch MCP server instance
	mcpServer, err := server.NewElasticsearchMCPServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		defer close(serverErrChan)
		if err := mcpServer.Start(); err != nil {
			serverErrChan <- err
		}
	}()

	// Wait for either a signal or server error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v, shutting down server...", sig)
	case err := <-serverErrChan:
		if err != nil {
			log.Printf("Server error: %v", err)
		}
	}

	// Gracefully shutdown the server
	if err := mcpServer.Stop(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Printf("Server stopped successfully")
}
