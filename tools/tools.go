// Package tools provides MCP tools for interacting with Elasticsearch.
// It implements various Elasticsearch operations as MCP tools that can be called
// through the Model Context Protocol interface.
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AeaZer/mcp-elasticsearch/elasticsearch"
	"github.com/mark3labs/mcp-go/mcp"
)

// ElasticsearchTools represents a collection of Elasticsearch-related MCP tools
type ElasticsearchTools struct {
	client elasticsearch.Client // Elasticsearch client for performing operations
}

// NewElasticsearchTools creates a new instance of ElasticsearchTools with the provided client
func NewElasticsearchTools(client elasticsearch.Client) *ElasticsearchTools {
	return &ElasticsearchTools{
		client: client,
	}
}

// GetTools returns all available Elasticsearch tools for the MCP server
func (t *ElasticsearchTools) GetTools() []mcp.Tool {
	return []mcp.Tool{
		// Cluster management tools
		{
			Name:        "es_cluster_info",
			Description: "Get Elasticsearch cluster information including version and basic configuration",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
			},
		},
		{
			Name:        "es_cluster_health",
			Description: "Get Elasticsearch cluster health status and statistics",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
			},
		},

		// Index management tools
		{
			Name:        "es_index_create",
			Description: "Create a new Elasticsearch index with optional settings and mappings",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"index": map[string]interface{}{
						"type":        "string",
						"description": "Name of the index to create",
					},
					"settings": map[string]interface{}{
						"type":        "object",
						"description": "Index settings (optional)",
					},
					"mappings": map[string]interface{}{
						"type":        "object",
						"description": "Index field mappings (optional)",
					},
				},
				Required: []string{"index"},
			},
		},
		{
			Name:        "es_index_delete",
			Description: "Delete an existing Elasticsearch index",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"index": map[string]interface{}{
						"type":        "string",
						"description": "Name of the index to delete",
					},
				},
				Required: []string{"index"},
			},
		},
		{
			Name:        "es_index_exists",
			Description: "Check if an Elasticsearch index exists",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"index": map[string]interface{}{
						"type":        "string",
						"description": "Name of the index to check",
					},
				},
				Required: []string{"index"},
			},
		},
		{
			Name:        "es_index_list",
			Description: "List all Elasticsearch indices with their metadata",
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
			},
		},

		// Document management tools
		{
			Name:        "es_document_index",
			Description: "Index a document in Elasticsearch with optional document ID",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"index": map[string]interface{}{
						"type":        "string",
						"description": "Name of the index to store the document",
					},
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Document ID (optional, will be auto-generated if not provided)",
					},
					"document": map[string]interface{}{
						"type":        "object",
						"description": "Document content to index",
					},
				},
				Required: []string{"index", "document"},
			},
		},
		{
			Name:        "es_document_get",
			Description: "Retrieve a document from Elasticsearch by its ID",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"index": map[string]interface{}{
						"type":        "string",
						"description": "Name of the index containing the document",
					},
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Document ID to retrieve",
					},
				},
				Required: []string{"index", "id"},
			},
		},
		{
			Name:        "es_document_update",
			Description: "Update an existing document in Elasticsearch",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"index": map[string]interface{}{
						"type":        "string",
						"description": "Name of the index containing the document",
					},
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Document ID to update",
					},
					"document": map[string]interface{}{
						"type":        "object",
						"description": "Updated document content or partial update",
					},
				},
				Required: []string{"index", "id", "document"},
			},
		},
		{
			Name:        "es_document_delete",
			Description: "Delete a document from Elasticsearch by its ID",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"index": map[string]interface{}{
						"type":        "string",
						"description": "Name of the index containing the document",
					},
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Document ID to delete",
					},
				},
				Required: []string{"index", "id"},
			},
		},

		// Search tools
		{
			Name:        "es_search",
			Description: "Execute a search query against Elasticsearch with optional filters and sorting",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"index": map[string]interface{}{
						"type":        "string",
						"description": "Name of the index to search (optional, searches all indices if not specified)",
					},
					"query": map[string]interface{}{
						"type":        "object",
						"description": "Elasticsearch query DSL",
					},
					"size": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results to return (default: 10)",
					},
					"from": map[string]interface{}{
						"type":        "integer",
						"description": "Number of results to skip for pagination (default: 0)",
					},
				},
				Required: []string{"query"},
			},
		},

		// Bulk operations
		{
			Name:        "es_bulk",
			Description: "Execute bulk operations (index, update, delete) in a single request for better performance",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"operations": map[string]interface{}{
						"type":        "array",
						"description": "Array of bulk operations to execute",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"operation": map[string]interface{}{
									"type":        "string",
									"description": "Type of operation: index, update, delete",
								},
								"index": map[string]interface{}{
									"type":        "string",
									"description": "Target index name",
								},
								"id": map[string]interface{}{
									"type":        "string",
									"description": "Document ID (optional for index operation)",
								},
								"body": map[string]interface{}{
									"type":        "object",
									"description": "Document body for index/update operations",
								},
							},
							"required": []string{"operation", "index"},
						},
					},
				},
				Required: []string{"operations"},
			},
		},
	}
}

// HandleTool processes MCP tool calls and routes them to the appropriate Elasticsearch operations
func (t *ElasticsearchTools) HandleTool(ctx context.Context, toolName string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	switch toolName {
	case "es_cluster_info":
		return t.handleClusterInfo(ctx)
	case "es_cluster_health":
		return t.handleClusterHealth(ctx)
	case "es_index_create":
		return t.handleIndexCreate(ctx, arguments)
	case "es_index_delete":
		return t.handleIndexDelete(ctx, arguments)
	case "es_index_exists":
		return t.handleIndexExists(ctx, arguments)
	case "es_index_list":
		return t.handleIndexList(ctx)
	case "es_document_index":
		return t.handleDocumentIndex(ctx, arguments)
	case "es_document_get":
		return t.handleDocumentGet(ctx, arguments)
	case "es_document_update":
		return t.handleDocumentUpdate(ctx, arguments)
	case "es_document_delete":
		return t.handleDocumentDelete(ctx, arguments)
	case "es_search":
		return t.handleSearch(ctx, arguments)
	case "es_bulk":
		return t.handleBulk(ctx, arguments)
	default:
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Unknown tool: %s", toolName),
				},
			},
			IsError: true,
		}, nil
	}
}

// Helper function to create error result
func createErrorResult(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: message,
			},
		},
		IsError: true,
	}
}

// Helper function to create success result with structured content
func createSuccessResult(text string, structuredContent interface{}) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: text,
			},
		},
		StructuredContent: structuredContent,
	}
}

// Helper function to create simple success result
func createSimpleSuccessResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}
}

// Cluster information handlers

func (t *ElasticsearchTools) handleClusterInfo(ctx context.Context) (*mcp.CallToolResult, error) {
	info, err := t.client.Info(ctx)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to get cluster info: %v", err)), nil
	}

	response, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to format cluster info: %v", err)), nil
	}

	return createSuccessResult(string(response), info), nil
}

func (t *ElasticsearchTools) handleClusterHealth(ctx context.Context) (*mcp.CallToolResult, error) {
	health, err := t.client.Health(ctx)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to get cluster health: %v", err)), nil
	}

	response, err := json.MarshalIndent(health, "", "  ")
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to format cluster health: %v", err)), nil
	}

	return createSuccessResult(string(response), health), nil
}

// Index management handlers

func (t *ElasticsearchTools) handleIndexCreate(ctx context.Context, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	indexName, ok := arguments["index"].(string)
	if !ok {
		return createErrorResult("Index name is required and must be a string"), nil
	}

	body := make(map[string]interface{})
	if settings, ok := arguments["settings"].(map[string]interface{}); ok {
		body["settings"] = settings
	}
	if mappings, ok := arguments["mappings"].(map[string]interface{}); ok {
		body["mappings"] = mappings
	}

	err := t.client.CreateIndex(ctx, indexName, body)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to create index '%s': %v", indexName, err)), nil
	}

	return createSimpleSuccessResult(fmt.Sprintf("Successfully created index '%s'", indexName)), nil
}

func (t *ElasticsearchTools) handleIndexDelete(ctx context.Context, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	indexName, ok := arguments["index"].(string)
	if !ok {
		return createErrorResult("Index name is required and must be a string"), nil
	}

	err := t.client.DeleteIndex(ctx, indexName)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to delete index '%s': %v", indexName, err)), nil
	}

	return createSimpleSuccessResult(fmt.Sprintf("Successfully deleted index '%s'", indexName)), nil
}

func (t *ElasticsearchTools) handleIndexExists(ctx context.Context, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	indexName, ok := arguments["index"].(string)
	if !ok {
		return createErrorResult("Index name is required and must be a string"), nil
	}

	exists, err := t.client.IndexExists(ctx, indexName)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to check if index '%s' exists: %v", indexName, err)), nil
	}

	result := map[string]interface{}{
		"index":  indexName,
		"exists": exists,
	}

	return createSuccessResult(fmt.Sprintf("Index '%s' exists: %t", indexName, exists), result), nil
}

func (t *ElasticsearchTools) handleIndexList(ctx context.Context) (*mcp.CallToolResult, error) {
	indices, err := t.client.ListIndices(ctx)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to list indices: %v", err)), nil
	}

	response, err := json.MarshalIndent(indices, "", "  ")
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to format indices list: %v", err)), nil
	}

	return createSuccessResult(string(response), indices), nil
}

// Document management handlers

func (t *ElasticsearchTools) handleDocumentIndex(ctx context.Context, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	indexName, ok := arguments["index"].(string)
	if !ok {
		return createErrorResult("Index name is required and must be a string"), nil
	}

	document, ok := arguments["document"].(map[string]interface{})
	if !ok {
		return createErrorResult("Document is required and must be an object"), nil
	}

	docID, _ := arguments["id"].(string) // Optional

	response, err := t.client.Index(ctx, indexName, docID, document)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to index document: %v", err)), nil
	}

	responseData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to format index response: %v", err)), nil
	}

	return createSuccessResult(string(responseData), response), nil
}

func (t *ElasticsearchTools) handleDocumentGet(ctx context.Context, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	indexName, ok := arguments["index"].(string)
	if !ok {
		return createErrorResult("Index name is required and must be a string"), nil
	}

	docID, ok := arguments["id"].(string)
	if !ok {
		return createErrorResult("Document ID is required and must be a string"), nil
	}

	response, err := t.client.Get(ctx, indexName, docID)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to get document: %v", err)), nil
	}

	responseData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to format get response: %v", err)), nil
	}

	return createSuccessResult(string(responseData), response), nil
}

func (t *ElasticsearchTools) handleDocumentUpdate(ctx context.Context, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	indexName, ok := arguments["index"].(string)
	if !ok {
		return createErrorResult("Index name is required and must be a string"), nil
	}

	docID, ok := arguments["id"].(string)
	if !ok {
		return createErrorResult("Document ID is required and must be a string"), nil
	}

	document, ok := arguments["document"].(map[string]interface{})
	if !ok {
		return createErrorResult("Document is required and must be an object"), nil
	}

	err := t.client.Update(ctx, indexName, docID, document)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to update document: %v", err)), nil
	}

	return createSimpleSuccessResult(fmt.Sprintf("Successfully updated document '%s' in index '%s'", docID, indexName)), nil
}

func (t *ElasticsearchTools) handleDocumentDelete(ctx context.Context, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	indexName, ok := arguments["index"].(string)
	if !ok {
		return createErrorResult("Index name is required and must be a string"), nil
	}

	docID, ok := arguments["id"].(string)
	if !ok {
		return createErrorResult("Document ID is required and must be a string"), nil
	}

	err := t.client.Delete(ctx, indexName, docID)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to delete document: %v", err)), nil
	}

	return createSimpleSuccessResult(fmt.Sprintf("Successfully deleted document '%s' from index '%s'", docID, indexName)), nil
}

// Search handlers

func (t *ElasticsearchTools) handleSearch(ctx context.Context, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	query, ok := arguments["query"].(map[string]interface{})
	if !ok {
		return createErrorResult("Query is required and must be an object"), nil
	}

	searchReq := &elasticsearch.SearchRequest{
		Query: query,
	}

	if index, ok := arguments["index"].(string); ok {
		searchReq.Index = index
	}

	if size, ok := arguments["size"].(float64); ok {
		searchReq.Size = int(size)
	} else {
		searchReq.Size = 10 // Default size
	}

	if from, ok := arguments["from"].(float64); ok {
		searchReq.From = int(from)
	}

	response, err := t.client.Search(ctx, searchReq)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to execute search: %v", err)), nil
	}

	responseData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to format search response: %v", err)), nil
	}

	return createSuccessResult(string(responseData), response), nil
}

// Bulk operation handlers

func (t *ElasticsearchTools) handleBulk(ctx context.Context, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	operationsData, ok := arguments["operations"].([]interface{})
	if !ok {
		return createErrorResult("Operations array is required"), nil
	}

	var operations []elasticsearch.BulkOperation
	for i, opData := range operationsData {
		opMap, ok := opData.(map[string]interface{})
		if !ok {
			return createErrorResult(fmt.Sprintf("Operation %d must be an object", i)), nil
		}

		operation := elasticsearch.BulkOperation{}

		if op, ok := opMap["operation"].(string); ok {
			operation.Operation = op
		} else {
			return createErrorResult(fmt.Sprintf("Operation %d must have an 'operation' field", i)), nil
		}

		if index, ok := opMap["index"].(string); ok {
			operation.Index = index
		} else {
			return createErrorResult(fmt.Sprintf("Operation %d must have an 'index' field", i)), nil
		}

		if id, ok := opMap["id"].(string); ok {
			operation.ID = id
		}

		if body, ok := opMap["body"].(map[string]interface{}); ok {
			operation.Body = body
		}

		operations = append(operations, operation)
	}

	response, err := t.client.Bulk(ctx, operations)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to execute bulk operations: %v", err)), nil
	}

	responseData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to format bulk response: %v", err)), nil
	}

	return createSuccessResult(string(responseData), response), nil
}
