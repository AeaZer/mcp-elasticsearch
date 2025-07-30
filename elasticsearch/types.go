// Package elasticsearch defines types and structures used for Elasticsearch operations.
// It contains response types, request types, and utility types for the MCP server.
package elasticsearch

import (
	"io"
	"log"
	"net/http"
	"time"
)

// InfoResponse represents the response from Elasticsearch cluster info API
type InfoResponse struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	ClusterUUID string `json:"cluster_uuid"`
	Version     struct {
		Number                           string `json:"number"`
		BuildFlavor                      string `json:"build_flavor"`
		BuildType                        string `json:"build_type"`
		BuildHash                        string `json:"build_hash"`
		BuildDate                        string `json:"build_date"`
		BuildSnapshot                    bool   `json:"build_snapshot"`
		LuceneVersion                    string `json:"lucene_version"`
		MinimumWireCompatibilityVersion  string `json:"minimum_wire_compatibility_version"`
		MinimumIndexCompatibilityVersion string `json:"minimum_index_compatibility_version"`
	} `json:"version"`
	TagLine string `json:"tagline"`
}

// HealthResponse represents the response from Elasticsearch cluster health API
type HealthResponse struct {
	ClusterName                 string  `json:"cluster_name"`
	Status                      string  `json:"status"`
	TimedOut                    bool    `json:"timed_out"`
	NumberOfNodes               int     `json:"number_of_nodes"`
	NumberOfDataNodes           int     `json:"number_of_data_nodes"`
	ActivePrimaryShards         int     `json:"active_primary_shards"`
	ActiveShards                int     `json:"active_shards"`
	RelocatingShards            int     `json:"relocating_shards"`
	InitializingShards          int     `json:"initializing_shards"`
	UnassignedShards            int     `json:"unassigned_shards"`
	DelayedUnassignedShards     int     `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        int     `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       int     `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis int     `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float64 `json:"active_shards_percent_as_number"`
}

// IndexInfo contains information about an Elasticsearch index
type IndexInfo struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}

// IndexResponse represents the response from document indexing operations
type IndexResponse struct {
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	ID      string `json:"_id"`
	Version int    `json:"_version"`
	Result  string `json:"result"`
	Shards  struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	SeqNo       int `json:"_seq_no"`
	PrimaryTerm int `json:"_primary_term"`
}

// GetResponse represents the response from document retrieval operations
type GetResponse struct {
	Index   string                 `json:"_index"`
	Type    string                 `json:"_type"`
	ID      string                 `json:"_id"`
	Version int                    `json:"_version"`
	SeqNo   int                    `json:"_seq_no"`
	Found   bool                   `json:"found"`
	Source  map[string]interface{} `json:"_source"`
}

// SearchRequest represents a search request to Elasticsearch
type SearchRequest struct {
	Index  string                 `json:"index,omitempty"`
	Query  map[string]interface{} `json:"query,omitempty"`
	Size   int                    `json:"size,omitempty"`
	From   int                    `json:"from,omitempty"`
	Sort   []interface{}          `json:"sort,omitempty"`
	Source interface{}            `json:"_source,omitempty"`
}

// SearchResponse represents the response from search operations
type SearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64     `json:"max_score"`
		Hits     []SearchHit `json:"hits"`
	} `json:"hits"`
}

// SearchHit represents a single search result
type SearchHit struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	ID     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

// BulkOperation represents a single operation in a bulk request
type BulkOperation struct {
	Operation string                 `json:"operation"`
	Index     string                 `json:"index"`
	Type      string                 `json:"type,omitempty"`
	ID        string                 `json:"id,omitempty"`
	Body      map[string]interface{} `json:"body,omitempty"`
}

// BulkResponse represents the response from bulk operations
type BulkResponse struct {
	Took   int                           `json:"took"`
	Errors bool                          `json:"errors"`
	Items  []map[string]BulkItemResponse `json:"items"`
}

// BulkItemResponse represents the response for a single item in a bulk operation
type BulkItemResponse struct {
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	ID      string `json:"_id"`
	Version int    `json:"_version"`
	Result  string `json:"result"`
	Status  int    `json:"status"`
	Error   *struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error,omitempty"`
}

// bodyReader implements io.Reader interface for request bodies
type bodyReader struct {
	data []byte
	pos  int
}

func (br *bodyReader) Read(p []byte) (int, error) {
	if br.pos >= len(br.data) {
		return 0, io.EOF
	}

	n := copy(p, br.data[br.pos:])
	br.pos += n
	return n, nil
}

func (br *bodyReader) Close() error {
	return nil
}

// esLogger implements a simple logger for the Elasticsearch client
type esLogger struct{}

// LogRoundTrip logs HTTP request/response information for debugging
func (l *esLogger) LogRoundTrip(req *http.Request, res *http.Response, err error, start time.Time, dur time.Duration) error {
	if err != nil {
		log.Printf("Elasticsearch request failed: %v", err)
	} else {
		log.Printf("Elasticsearch %s %s %d %s", req.Method, req.URL.Path, res.StatusCode, dur)
	}
	return nil
}

// RequestBodyEnabled returns true to enable request body logging
func (l *esLogger) RequestBodyEnabled() bool {
	return false
}

// ResponseBodyEnabled returns true to enable response body logging
func (l *esLogger) ResponseBodyEnabled() bool {
	return false
}

// closeBody is a helper function to safely close response bodies
func closeBody(body io.Closer) {
	if err := body.Close(); err != nil {
		log.Printf("Failed to close response body: %v", err)
	}
}
