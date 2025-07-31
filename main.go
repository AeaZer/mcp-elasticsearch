// Package main implements the Elasticsearch MCP (Model Context Protocol) server
// that provides an interface for interacting with Elasticsearch clusters
// through both stdio and Streamable HTTP protocols.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AeaZer/mcp-elasticsearch/config"
	"github.com/AeaZer/mcp-elasticsearch/server"
)

func main() {
	// Configure logging format with timestamp and file location
	log.SetFlags(log.Ldate | log.Ltime | log.LstdFlags)

	log.Printf("Starting Elasticsearch MCP Server...")

	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v", err)
	}

	// Log essential startup information
	log.Printf("Protocol: %s", cfg.Server.Protocol)
	if cfg.Server.Protocol == "http" || cfg.Server.Protocol == "sse" {
		log.Printf("Listen Address: %s:%d", cfg.Server.Address, cfg.Server.Port)
	}
	log.Printf("Elasticsearch: %v", cfg.Elasticsearch.Addresses)

	// Create the Elasticsearch MCP server instance
	mcpServer, err := server.NewElasticsearchMCPServer(cfg)
	if err != nil {
		log.Fatalf("FATAL: Failed to create MCP server: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		defer close(serverErrChan)
		if err := mcpServer.Start(); err != nil {
			log.Printf("ERROR: Server failed: %v", err)
			serverErrChan <- err
		}
	}()

	// Give the server a moment to start
	time.Sleep(1 * time.Second)

	if cfg.Server.Protocol == "http" || cfg.Server.Protocol == "sse" {
		endpoint := "mcp"
		if cfg.Server.Protocol == "sse" {
			endpoint = "sse"
		}
		log.Printf("MCP Server running at: http://%s:%d/%s", cfg.Server.Address, cfg.Server.Port, endpoint)
	}

	// Wait for either a signal or server error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v, shutting down...", sig)
	case err := <-serverErrChan:
		if err != nil {
			log.Printf("Server error: %v", err)
		}
	}

	// Gracefully shutdown the server
	if err := mcpServer.Stop(); err != nil {
		log.Printf("ERROR during shutdown: %v", err)
	}

	log.Printf("Server stopped")
}
