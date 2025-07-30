# Elasticsearch MCP Server

*Read this in other languages: [English](README.md), [中文](README_zh.md)*

## Overview

An Elasticsearch MCP (Model Context Protocol) server built on [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go), providing seamless integration with Elasticsearch 7, 8, and 9 versions.

## Features

- 🔗 **Multi-Protocol Support**: Supports both stdio and Streamable HTTP protocols
- 📊 **Multi-Version Compatibility**: Compatible with Elasticsearch 7, 8, and 9
- ⚙️ **Environment Configuration**: Configure via environment variables
- 🔧 **Rich Toolset**: Complete set of Elasticsearch operation tools
- 🌐 **Production Ready**: Docker support with optimized builds

## Supported Tools

### Cluster Operations
- `es_cluster_info`: Get cluster information and version details
- `es_cluster_health`: Get cluster health status and metrics

### Index Management
- `es_index_create`: Create new indices with settings and mappings
- `es_index_delete`: Delete existing indices
- `es_index_exists`: Check if an index exists
- `es_index_list`: List all indices with metadata

### Document Operations
- `es_document_index`: Index documents with optional ID
- `es_document_get`: Retrieve documents by ID
- `es_document_update`: Update existing documents
- `es_document_delete`: Delete documents by ID

### Search Operations
- `es_search`: Execute search queries with filters and sorting

### Bulk Operations
- `es_bulk`: Execute multiple operations in a single request

## Quick Start

Choose one of the following methods to run the Elasticsearch MCP server:

### Method 1: Build Docker Image (Recommended)

```bash
# Clone the repository
git clone https://github.com/AeaZer/mcp-elasticsearch.git
cd mcp-elasticsearch

# Build Docker image
docker build -t mcp-elasticsearch .

# Run the container
docker run -e ES_ADDRESSES=http://localhost:9200 -e ES_VERSION=8 mcp-elasticsearch
```

### Method 2: Use Pre-built Image (Coming Soon)

```bash
# This will be available when the image is published to a registry
# docker run -e ES_ADDRESSES=http://localhost:9200 ghcr.io/aeazer/mcp-elasticsearch:latest
```

*Note: Pre-built images are not yet available. Please use Method 1 or Method 3.*

### Method 3: Compile from Source

```bash
# Clone the repository
git clone https://github.com/AeaZer/mcp-elasticsearch.git
cd mcp-elasticsearch

# Download dependencies and build
go mod download
go build -o mcp-elasticsearch main.go

# Run with environment variables
export ES_ADDRESSES=http://localhost:9200
export ES_VERSION=8
export MCP_PROTOCOL=stdio
./mcp-elasticsearch
```

## Configuration

All configuration is done via environment variables:

### Elasticsearch Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `ES_ADDRESSES` | Elasticsearch cluster addresses (comma-separated) | `http://localhost:9200` |
| `ES_USERNAME` | Username for basic authentication | - |
| `ES_PASSWORD` | Password for basic authentication | - |
| `ES_API_KEY` | API Key for authentication | - |
| `ES_CLOUD_ID` | Elastic Cloud ID | - |
| `ES_SSL` | Enable SSL/TLS | `false` |
| `ES_INSECURE_SKIP_VERIFY` | Skip SSL certificate verification | `false` |
| `ES_TIMEOUT` | Connection timeout | `30s` |
| `ES_MAX_RETRIES` | Maximum retry attempts | `3` |
| `ES_VERSION` | Target Elasticsearch version (7, 8, or 9) | `8` |

### MCP Server Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `MCP_SERVER_NAME` | Server name for MCP | `Elasticsearch MCP Server` |
| `MCP_SERVER_VERSION` | Server version | `1.0.0` |
| `MCP_PROTOCOL` | Protocol to use (`stdio` or `http`) | `stdio` |
| `MCP_ADDRESS` | Streamable HTTP server address (HTTP mode only) | `localhost` |
| `MCP_PORT` | Streamable HTTP server port (HTTP mode only) | `8080` |

## Usage Examples

### Stdio Mode (Default)
```bash
export ES_ADDRESSES=http://localhost:9200
export MCP_PROTOCOL=stdio
./mcp-elasticsearch
```

### Streamable HTTP Mode
```bash
export ES_ADDRESSES=http://localhost:9200
export MCP_PROTOCOL=http
export MCP_PORT=8080
./mcp-elasticsearch
```

### Using with Elastic Cloud
```bash
export ES_CLOUD_ID=your_cloud_id
export ES_USERNAME=elastic
export ES_PASSWORD=your_password
export ES_VERSION=8
./mcp-elasticsearch
```


## Development

### Prerequisites
- Go 1.21 or higher
- Access to an Elasticsearch cluster

### Building
```bash
go mod download
go build -o mcp-elasticsearch main.go
```

### Testing
```bash
go test ./...
```

## Tool Usage Examples

### Get Cluster Information
```json
{
  "tool": "es_cluster_info",
  "arguments": {}
}
```

### Create an Index
```json
{
  "tool": "es_index_create",
  "arguments": {
    "index": "my-index",
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 0
    },
    "mappings": {
      "properties": {
        "title": {"type": "text"},
        "timestamp": {"type": "date"}
      }
    }
  }
}
```

### Index a Document
```json
{
  "tool": "es_document_index",
  "arguments": {
    "index": "my-index",
    "id": "doc1",
    "document": {
      "title": "Hello World",
      "content": "This is a test document",
      "timestamp": "2024-01-01T00:00:00Z"
    }
  }
}
```

### Search Documents
```json
{
  "tool": "es_search",
  "arguments": {
    "index": "my-index",
    "query": {
      "match": {
        "title": "Hello"
      }
    },
    "size": 10
  }
}
```

## Error Handling

All errors are reported within the MCP tool results with `isError: true`, allowing LLMs to see and handle errors appropriately. Protocol-level errors are reserved for exceptional conditions like missing tools or server failures.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


## Acknowledgments

- [Mark3Labs MCP-Go](https://github.com/mark3labs/mcp-go) - MCP implementation for Go
- [Elastic](https://github.com/elastic/go-elasticsearch) - Official Elasticsearch Go client
- [Model Context Protocol](https://modelcontextprotocol.io/) - Protocol specification 