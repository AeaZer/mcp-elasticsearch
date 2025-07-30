// Package elasticsearch provides Elasticsearch client functionality with support for multiple versions.
// It implements a comprehensive interface for interacting with Elasticsearch clusters
// through the official Go client library.
package elasticsearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AeaZer/mcp-elasticsearch/config"
	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client defines the interface for Elasticsearch operations.
// It abstracts the underlying Elasticsearch client to provide a consistent API
// for all supported Elasticsearch versions (7, 8, 9).
type Client interface {
	Info(ctx context.Context) (*InfoResponse, error)
	Health(ctx context.Context) (*HealthResponse, error)

	CreateIndex(ctx context.Context, index string, body map[string]interface{}) error
	DeleteIndex(ctx context.Context, index string) error
	IndexExists(ctx context.Context, index string) (bool, error)
	ListIndices(ctx context.Context) ([]IndexInfo, error)

	Index(ctx context.Context, index, docID string, body map[string]interface{}) (*IndexResponse, error)
	Get(ctx context.Context, index, docID string) (*GetResponse, error)
	Delete(ctx context.Context, index, docID string) error
	Update(ctx context.Context, index, docID string, body map[string]interface{}) error

	Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)

	Bulk(ctx context.Context, operations []BulkOperation) (*BulkResponse, error)

	Close() error
}

// ESClient implements the Client interface using the official Elasticsearch Go client.
// It provides a concrete implementation that can work with Elasticsearch 7, 8, and 9.
type ESClient struct {
	client  *elasticsearch8.Client      // The underlying Elasticsearch client
	config  *config.ElasticsearchConfig // Configuration for the client
	version string                      // Elasticsearch version string
}

// NewClient creates a new Elasticsearch client with the provided configuration.
// It configures authentication, SSL settings, and retry policies based on the config.
//
// Parameters:
//   - cfg: Elasticsearch configuration containing connection details
//   - version: Target Elasticsearch version string
//
// Returns:
//   - Client: Configured Elasticsearch client interface
//   - error: Any error that occurred during client creation
func NewClient(cfg *config.ElasticsearchConfig, version string) (Client, error) {
	// Configure the Elasticsearch client with connection settings
	esConfig := elasticsearch8.Config{
		Addresses: cfg.Addresses,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.InsecureSkipVerify,
			},
		},
		MaxRetries:    cfg.MaxRetries,
		RetryOnStatus: []int{502, 503, 504, 429},
		Logger:        &esLogger{},
	}

	// Configure authentication based on provided credentials
	if cfg.Username != "" && cfg.Password != "" {
		esConfig.Username = cfg.Username
		esConfig.Password = cfg.Password
	}

	if cfg.APIKey != "" {
		esConfig.APIKey = cfg.APIKey
	}

	if cfg.CloudID != "" {
		esConfig.CloudID = cfg.CloudID
	}

	// Create the Elasticsearch client
	client, err := elasticsearch8.NewClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	return &ESClient{
		client:  client,
		config:  cfg,
		version: version,
	}, nil
}

// Info retrieves cluster information from Elasticsearch.
// This includes cluster name, version, and other basic information.
func (c *ESClient) Info(ctx context.Context) (*InfoResponse, error) {
	res, err := c.client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %v", res)
	}

	var info InfoResponse
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &info, nil
}

// Health retrieves the cluster health status from Elasticsearch.
// This provides information about cluster status, number of nodes, etc.
func (c *ESClient) Health(ctx context.Context) (*HealthResponse, error) {
	res, err := c.client.Cluster.Health()
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster health: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %v", res)
	}

	var health HealthResponse
	if err := json.NewDecoder(res.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &health, nil
}

// CreateIndex creates a new index in Elasticsearch with the specified configuration.
//
// Parameters:
//   - ctx: Context for request cancellation
//   - index: Name of the index to create
//   - body: Index configuration (mappings, settings, etc.)
func (c *ESClient) CreateIndex(ctx context.Context, index string, body map[string]interface{}) error {
	var req esapi.IndicesCreateRequest
	req.Index = index

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to serialize request body: %w", err)
		}
		req.Body = &bodyReader{data: bodyBytes}
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

// DeleteIndex removes an index from Elasticsearch.
//
// Parameters:
//   - ctx: Context for request cancellation
//   - index: Name of the index to delete
func (c *ESClient) DeleteIndex(ctx context.Context, index string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{index},
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to delete index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

// IndexExists checks whether an index exists in Elasticsearch.
//
// Parameters:
//   - ctx: Context for request cancellation
//   - index: Name of the index to check
//
// Returns:
//   - bool: true if the index exists, false otherwise
//   - error: Any error that occurred during the check
func (c *ESClient) IndexExists(ctx context.Context, index string) (bool, error) {
	req := esapi.IndicesExistsRequest{
		Index: []string{index},
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return false, fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

// ListIndices retrieves a list of all indices in the Elasticsearch cluster.
//
// Returns:
//   - []IndexInfo: List of index information
//   - error: Any error that occurred during the operation
func (c *ESClient) ListIndices(ctx context.Context) ([]IndexInfo, error) {
	req := esapi.CatIndicesRequest{
		Format: "json",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to list indices: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var indices []IndexInfo
	if err := json.NewDecoder(res.Body).Decode(&indices); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return indices, nil
}

// Index adds or updates a document in Elasticsearch.
//
// Parameters:
//   - ctx: Context for request cancellation
//   - index: Name of the target index
//   - docID: Document ID (empty string for auto-generation)
//   - body: Document content as a map
//
// Returns:
//   - *IndexResponse: Response containing operation details
//   - error: Any error that occurred during indexing
func (c *ESClient) Index(ctx context.Context, index, docID string, body map[string]interface{}) (*IndexResponse, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: docID,
		Body:       &bodyReader{data: bodyBytes},
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var indexResp IndexResponse
	if err := json.NewDecoder(res.Body).Decode(&indexResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &indexResp, nil
}

// Get retrieves a document from Elasticsearch by its ID.
//
// Parameters:
//   - ctx: Context for request cancellation
//   - index: Name of the source index
//   - docID: Document ID to retrieve
//
// Returns:
//   - *GetResponse: Response containing the document
//   - error: Any error that occurred during retrieval
func (c *ESClient) Get(ctx context.Context, index, docID string) (*GetResponse, error) {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: docID,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var getResp GetResponse
	if err := json.NewDecoder(res.Body).Decode(&getResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &getResp, nil
}

// Delete removes a document from Elasticsearch.
//
// Parameters:
//   - ctx: Context for request cancellation
//   - index: Name of the source index
//   - docID: Document ID to delete
func (c *ESClient) Delete(ctx context.Context, index, docID string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: docID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

// Update partially updates a document in Elasticsearch.
//
// Parameters:
//   - ctx: Context for request cancellation
//   - index: Name of the target index
//   - docID: Document ID to update
//   - body: Partial document content for update
func (c *ESClient) Update(ctx context.Context, index, docID string, body map[string]interface{}) error {
	updateBody := map[string]interface{}{
		"doc": body,
	}

	bodyBytes, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("failed to serialize update body: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: docID,
		Body:       &bodyReader{data: bodyBytes},
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

// Search executes a search query against Elasticsearch.
//
// Parameters:
//   - ctx: Context for request cancellation
//   - req: Search request containing query, index, pagination, etc.
//
// Returns:
//   - *SearchResponse: Search results with hits and metadata
//   - error: Any error that occurred during search
func (c *ESClient) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	bodyBytes, err := json.Marshal(req.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize search query: %w", err)
	}

	esReq := esapi.SearchRequest{
		Index: []string{req.Index},
		Body:  &bodyReader{data: bodyBytes},
		Size:  &req.Size,
		From:  &req.From,
	}

	res, err := esReq.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	return &searchResp, nil
}

// Bulk performs multiple operations in a single request.
// This is more efficient than individual operations for large datasets.
//
// Parameters:
//   - ctx: Context for request cancellation
//   - operations: List of bulk operations to perform
//
// Returns:
//   - *BulkResponse: Results of all bulk operations
//   - error: Any error that occurred during bulk operation
func (c *ESClient) Bulk(ctx context.Context, operations []BulkOperation) (*BulkResponse, error) {
	// Build the bulk request body in NDJSON format
	body := ""
	for _, op := range operations {
		// Create the action line based on operation type
		action := map[string]interface{}{
			op.Operation: map[string]interface{}{
				"_index": op.Index,
			},
		}
		if op.Type != "" {
			action[op.Operation].(map[string]interface{})["_type"] = op.Type
		}
		if op.ID != "" {
			action[op.Operation].(map[string]interface{})["_id"] = op.ID
		}

		actionBytes, err := json.Marshal(action)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize bulk operation: %w", err)
		}
		body += string(actionBytes) + "\n"

		// Add the source document if it exists (not needed for delete operations)
		if op.Body != nil && op.Operation != "delete" {
			sourceBytes, err := json.Marshal(op.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to serialize bulk operation source: %w", err)
			}
			body += string(sourceBytes) + "\n"
		}
	}

	req := esapi.BulkRequest{
		Body:    &bodyReader{data: []byte(body)},
		Refresh: "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("bulk operation failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var bulkResp BulkResponse
	if err := json.NewDecoder(res.Body).Decode(&bulkResp); err != nil {
		return nil, fmt.Errorf("failed to parse bulk response: %w", err)
	}

	return &bulkResp, nil
}

// Close gracefully closes the Elasticsearch client connection.
// Note: The official Elasticsearch Go client doesn't require explicit closing.
func (c *ESClient) Close() error {
	// The Elasticsearch Go client doesn't require explicit closing
	return nil
}
