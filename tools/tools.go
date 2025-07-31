// Package tools provides MCP tools for interacting with Elasticsearch.
// It implements various Elasticsearch operations as MCP tools that can be called
// through the Model Context Protocol interface.
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AeaZer/mcp-elasticsearch/elasticsearch"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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

// GetTools returns all available Elasticsearch tools with their schemas
func (et *ElasticsearchTools) GetTools() []mcp.Tool {
	return []mcp.Tool{
		{
			Name:        "es_cluster_info",
			Description: "Get cluster information and version details",
			InputSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{},
			},
		},
		{
			Name:        "es_cluster_health",
			Description: "Get cluster health status and metrics",
			InputSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{},
			},
		},
		{
			Name:        "es_index_create",
			Description: "Create a new index with optional settings and mappings",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"index": {
						Type:        "string",
						Description: "Index name to create",
					},
					"settings": {
						Type:        "object",
						Description: "Index settings (optional)",
					},
					"mappings": {
						Type:        "object",
						Description: "Index mappings (optional)",
					},
				},
				Required: []string{"index"},
			},
		},
		{
			Name:        "es_index_delete",
			Description: "Delete an existing index",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"index": {
						Type:        "string",
						Description: "Index name to delete",
					},
				},
				Required: []string{"index"},
			},
		},
		{
			Name:        "es_index_exists",
			Description: "Check if an index exists",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"index": {
						Type:        "string",
						Description: "Index name to check",
					},
				},
				Required: []string{"index"},
			},
		},
		{
			Name:        "es_index_list",
			Description: "List all indices with metadata",
			InputSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{},
			},
		},
		{
			Name:        "es_document_index",
			Description: "Index a document with optional ID",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"index": {
						Type:        "string",
						Description: "Index name",
					},
					"id": {
						Type:        "string",
						Description: "Document ID (optional)",
					},
					"body": {
						Type:        "object",
						Description: "Document body",
					},
				},
				Required: []string{"index", "body"},
			},
		},
		{
			Name:        "es_document_get",
			Description: "Retrieve a document by ID",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"index": {
						Type:        "string",
						Description: "Index name",
					},
					"id": {
						Type:        "string",
						Description: "Document ID",
					},
				},
				Required: []string{"index", "id"},
			},
		},
		{
			Name:        "es_document_update",
			Description: "Update an existing document",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"index": {
						Type:        "string",
						Description: "Index name",
					},
					"id": {
						Type:        "string",
						Description: "Document ID",
					},
					"body": {
						Type:        "object",
						Description: "Update body with doc or script",
					},
				},
				Required: []string{"index", "id", "body"},
			},
		},
		{
			Name:        "es_document_delete",
			Description: "Delete a document by ID",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"index": {
						Type:        "string",
						Description: "Index name",
					},
					"id": {
						Type:        "string",
						Description: "Document ID",
					},
				},
				Required: []string{"index", "id"},
			},
		},
		{
			Name:        "es_search",
			Description: "Execute search queries with filters and sorting",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"index": {
						Type:        "string",
						Description: "Index name (optional, searches all if not provided)",
					},
					"query": {
						Type:        "object",
						Description: "Search query",
					},
					"size": {
						Type:        "integer",
						Description: "Number of results to return (default: 10)",
					},
					"from": {
						Type:        "integer",
						Description: "Offset for pagination (default: 0)",
					},
					"sort": {
						Type:        "array",
						Description: "Sort specification (array of sort objects)",
						Items: &jsonschema.Schema{
							Type: "object",
						},
					},
					"_source": {
						Description: "Source filtering: boolean (true/false), array of field names, or object with includes/excludes",
						OneOf: []*jsonschema.Schema{
							{Type: "boolean"},
							{Type: "array", Items: &jsonschema.Schema{Type: "string"}},
							{Type: "object"},
						},
					},
				},
			},
		},
		{
			Name:        "es_bulk",
			Description: "Execute multiple operations in a single request",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"operations": {
						Type:        "array",
						Description: "Array of bulk operations",
						Items: &jsonschema.Schema{
							Type: "object",
						},
					},
				},
				Required: []string{"operations"},
			},
		},
	}
}

// Helper functions to create standardized MCP results

// createErrorResult creates a standardized error result
func createErrorResult(message string) mcp.CallToolResult {
	return mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Error: %s", message),
			},
		},
		IsError: true,
	}
}

// createSuccessResult creates a standardized success result with structured data
func createSuccessResult(text string, data interface{}) mcp.CallToolResult {
	content := []mcp.Content{
		&mcp.TextContent{
			Text: text,
		},
	}

	result := mcp.CallToolResult{
		Content: content,
		IsError: false,
	}

	if data != nil {
		switch v := data.(type) {
		case []interface{}:
			result.StructuredContent = map[string]interface{}{
				"data":  v,
				"count": len(v),
			}
		case []string:
			result.StructuredContent = map[string]interface{}{
				"data":  v,
				"count": len(v),
			}
		default:
			result.StructuredContent = data
		}
	}

	return result
}

// createSimpleSuccessResult creates a simple success result with only text
func createSimpleSuccessResult(text string) mcp.CallToolResult {
	return mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: text,
			},
		},
		IsError: false,
	}
}

// HandleTool handles MCP tool calls and routes them to appropriate handlers
func (et *ElasticsearchTools) HandleTool(ctx context.Context, toolName string, arguments map[string]interface{}) mcp.CallToolResult {
	switch toolName {
	case "es_cluster_info":
		return et.handleClusterInfo(ctx)
	case "es_cluster_health":
		return et.handleClusterHealth(ctx)
	case "es_index_create":
		return et.handleIndexCreate(ctx, arguments)
	case "es_index_delete":
		return et.handleIndexDelete(ctx, arguments)
	case "es_index_exists":
		return et.handleIndexExists(ctx, arguments)
	case "es_index_list":
		return et.handleIndexList(ctx)
	case "es_document_index":
		return et.handleDocumentIndex(ctx, arguments)
	case "es_document_get":
		return et.handleDocumentGet(ctx, arguments)
	case "es_document_update":
		return et.handleDocumentUpdate(ctx, arguments)
	case "es_document_delete":
		return et.handleDocumentDelete(ctx, arguments)
	case "es_search":
		return et.handleSearch(ctx, arguments)
	case "es_bulk":
		return et.handleBulk(ctx, arguments)
	default:
		return createErrorResult(fmt.Sprintf("Unknown tool: %s", toolName))
	}
}

// Individual tool handlers
func (et *ElasticsearchTools) handleClusterInfo(ctx context.Context) mcp.CallToolResult {
	info, err := et.client.Info(ctx)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to get cluster info: %v", err))
	}
	return createSuccessResult("Cluster info retrieved successfully", info)
}

func (et *ElasticsearchTools) handleClusterHealth(ctx context.Context) mcp.CallToolResult {
	health, err := et.client.Health(ctx)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to get cluster health: %v", err))
	}

	return createSuccessResult("Cluster health retrieved successfully", health)
}

func (et *ElasticsearchTools) handleIndexCreate(ctx context.Context, args map[string]interface{}) mcp.CallToolResult {
	index, ok := args["index"].(string)
	if !ok {
		return createErrorResult("Missing or invalid 'index' parameter")
	}

	// Build the request body
	body := make(map[string]interface{})
	if settings, exists := args["settings"]; exists {
		if settingsMap, ok := settings.(map[string]interface{}); ok {
			body["settings"] = settingsMap
		}
	}
	if mappings, exists := args["mappings"]; exists {
		if mappingsMap, ok := mappings.(map[string]interface{}); ok {
			body["mappings"] = mappingsMap
		}
	}

	err := et.client.CreateIndex(ctx, index, body)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to create index: %v", err))
	}

	return createSimpleSuccessResult(fmt.Sprintf("Index '%s' created successfully", index))
}

func (et *ElasticsearchTools) handleIndexDelete(ctx context.Context, args map[string]interface{}) mcp.CallToolResult {
	index, ok := args["index"].(string)
	if !ok {
		return createErrorResult("Missing or invalid 'index' parameter")
	}

	err := et.client.DeleteIndex(ctx, index)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to delete index: %v", err))
	}

	return createSimpleSuccessResult(fmt.Sprintf("Index '%s' deleted successfully", index))
}

func (et *ElasticsearchTools) handleIndexExists(ctx context.Context, args map[string]interface{}) mcp.CallToolResult {
	index, ok := args["index"].(string)
	if !ok {
		return createErrorResult("Missing or invalid 'index' parameter")
	}

	exists, err := et.client.IndexExists(ctx, index)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to check index existence: %v", err))
	}

	result := map[string]interface{}{
		"index":  index,
		"exists": exists,
	}

	return createSuccessResult(fmt.Sprintf("Index '%s' exists: %t", index, exists), result)
}

func (et *ElasticsearchTools) handleIndexList(ctx context.Context) mcp.CallToolResult {
	indices, err := et.client.ListIndices(ctx)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to list indices: %v", err))
	}

	result := map[string]interface{}{
		"indices": indices,
		"count":   len(indices),
	}

	return createSuccessResult(fmt.Sprintf("Found %d indices", len(indices)), result)
}

func (et *ElasticsearchTools) handleDocumentIndex(ctx context.Context, args map[string]interface{}) mcp.CallToolResult {
	index, ok := args["index"].(string)
	if !ok {
		return createErrorResult("Missing or invalid 'index' parameter")
	}

	body, ok := args["body"].(map[string]interface{})
	if !ok {
		return createErrorResult("Missing or invalid 'body' parameter")
	}

	id, _ := args["id"].(string) // Optional parameter

	result, err := et.client.Index(ctx, index, id, body)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to index document: %v", err))
	}

	return createSuccessResult("Document indexed successfully", result)
}

func (et *ElasticsearchTools) handleDocumentGet(ctx context.Context, args map[string]interface{}) mcp.CallToolResult {
	index, ok := args["index"].(string)
	if !ok {
		return createErrorResult("Missing or invalid 'index' parameter")
	}

	id, ok := args["id"].(string)
	if !ok {
		return createErrorResult("Missing or invalid 'id' parameter")
	}

	result, err := et.client.Get(ctx, index, id)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to get document: %v", err))
	}

	return createSuccessResult("Document retrieved successfully", result)
}

func (et *ElasticsearchTools) handleDocumentUpdate(ctx context.Context, args map[string]interface{}) mcp.CallToolResult {
	index, ok := args["index"].(string)
	if !ok {
		return createErrorResult("Missing or invalid 'index' parameter")
	}

	id, ok := args["id"].(string)
	if !ok {
		return createErrorResult("Missing or invalid 'id' parameter")
	}

	body, ok := args["body"].(map[string]interface{})
	if !ok {
		return createErrorResult("Missing or invalid 'body' parameter")
	}

	err := et.client.Update(ctx, index, id, body)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to update document: %v", err))
	}

	return createSimpleSuccessResult(fmt.Sprintf("Document '%s' updated successfully in index '%s'", id, index))
}

func (et *ElasticsearchTools) handleDocumentDelete(ctx context.Context, args map[string]interface{}) mcp.CallToolResult {
	index, ok := args["index"].(string)
	if !ok {
		return createErrorResult("Missing or invalid 'index' parameter")
	}

	id, ok := args["id"].(string)
	if !ok {
		return createErrorResult("Missing or invalid 'id' parameter")
	}

	err := et.client.Delete(ctx, index, id)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to delete document: %v", err))
	}

	return createSimpleSuccessResult(fmt.Sprintf("Document '%s' deleted successfully from index '%s'", id, index))
}

func (et *ElasticsearchTools) handleSearch(ctx context.Context, args map[string]interface{}) mcp.CallToolResult {
	// Index is optional for search
	index, _ := args["index"].(string)

	// Default query if none provided (this should be the query content, not wrapped in "query")
	query := map[string]interface{}{
		"match_all": map[string]interface{}{},
	}
	if q, exists := args["query"]; exists {
		if queryMap, ok := q.(map[string]interface{}); ok {
			query = queryMap
		}
	}

	// Default size
	size := 10
	if s, exists := args["size"]; exists {
		if sizeInt, ok := s.(float64); ok {
			size = int(sizeInt)
		}
	}

	// Default from
	from := 0
	if f, exists := args["from"]; exists {
		if fromInt, ok := f.(int64); ok {
			from = int(fromInt)
		}
	}

	// Parse sort parameter
	var sort []interface{}
	if s, exists := args["sort"]; exists {
		if sortArray, ok := s.([]interface{}); ok {
			sort = sortArray
		}
	}

	// Parse _source parameter
	var source interface{}
	if src, exists := args["_source"]; exists {
		source = src
	}

	searchRequest := &elasticsearch.SearchRequest{
		Index:  index,
		Query:  query,
		Size:   size,
		From:   from,
		Sort:   sort,
		Source: source,
	}

	result, err := et.client.Search(ctx, searchRequest)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to execute search: %v", err))
	}

	return createSuccessResult("Search executed successfully", result)
}

func (et *ElasticsearchTools) handleBulk(ctx context.Context, args map[string]interface{}) mcp.CallToolResult {
	operations, ok := args["operations"].([]interface{})
	if !ok {
		return createErrorResult("Missing or invalid 'operations' parameter")
	}

	// Convert to BulkOperation slice
	bulkOps := make([]elasticsearch.BulkOperation, len(operations))
	for i, op := range operations {
		if opMap, ok := op.(map[string]interface{}); ok {
			// Convert the operation to the expected format
			if opBytes, err := json.Marshal(opMap); err == nil {
				var bulkOp elasticsearch.BulkOperation
				if err := json.Unmarshal(opBytes, &bulkOp); err == nil {
					bulkOps[i] = bulkOp
				}
			}
		}
	}

	result, err := et.client.Bulk(ctx, bulkOps)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Failed to execute bulk operations: %v", err))
	}

	return createSuccessResult("Bulk operations executed successfully", result)
}
