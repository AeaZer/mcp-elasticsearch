# Elasticsearch MCP Server

*Read this in other languages: [English](README.md), [中文](README_zh.md)*

## Overview

An Elasticsearch MCP (Model Context Protocol) server built on [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk), providing seamless integration with Elasticsearch 7, 8, and 9 versions.

## Features

- 🔗 **Multi-Protocol Support**: Supports stdio, Streamable HTTP, and SSE protocols (SSE deprecated)
- 📊 **Multi-Version Compatibility**: Compatible with Elasticsearch 7, 8, and 9
- ⚙️ **Environment Configuration**: Configure via environment variables
- 🔧 **Rich Toolset**: Complete set of Elasticsearch operation tools
- 🌐 **Production Ready**: Docker support with optimized builds
- 🐳 **Container Ready**: Pre-built Docker images available

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
- `es_search`: Execute search queries with filters, sorting, and field selection
  - Supports: `index`, `query`, `size`, `from`, `sort`, `_source`
  - Full Elasticsearch Query DSL support

### Bulk Operations
- `es_bulk`: Execute multiple operations in a single request

## Quick Start

Choose one of the following methods to run the Elasticsearch MCP server:

### Method 1: Use Pre-built Docker Image (Recommended)

```bash
# Basic usage with local Elasticsearch
docker run --rm \
  -e ES_ADDRESSES=http://localhost:9200 \
  ghcr.io/aeazer/mcp-elasticsearch:latest

# HTTP mode for remote access
docker run -d -p 8080:8080 \
  -e MCP_PROTOCOL=http \
  -e ES_ADDRESSES=http://your-elasticsearch:9200 \
  ghcr.io/aeazer/mcp-elasticsearch:latest

# With authentication
docker run -d -p 8080:8080 \
  -e MCP_PROTOCOL=http \
  -e ES_ADDRESSES=https://your-elasticsearch:9200 \
  -e ES_USERNAME=elastic \
  -e ES_PASSWORD=your-password \
  -e ES_SSL=true \
  ghcr.io/aeazer/mcp-elasticsearch:latest
```

### Method 2: Build Docker Image

```bash
# Clone the repository
git clone https://github.com/AeaZer/mcp-elasticsearch.git
cd mcp-elasticsearch

# Build Docker image
docker build -t mcp-elasticsearch .

# Run the container
docker run -e ES_ADDRESSES=http://localhost:9200 -e ES_VERSION=8 mcp-elasticsearch
```

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

### Method 4: Desktop Application Integration (Cursor, Claude Desktop, etc.)

For integration with desktop applications that support MCP, you can configure the server through a configuration file:

#### Step 1: Build or Install the Executable

First, ensure you have the `mcp-elasticsearch` executable available:

```bash
# Option A: Install directly from GitHub (Recommended)
go install github.com/AeaZer/mcp-elasticsearch@latest

# Option B: Build from source
git clone https://github.com/AeaZer/mcp-elasticsearch.git
cd mcp-elasticsearch
go mod download
go build -o mcp-elasticsearch main.go

# Option C: Download pre-built binary (Windows only)
# Download mcp-elasticsearch.exe from GitHub Releases
# macOS and Linux users should use Option A or B
```

> **Note**: Docker is not suitable for desktop application integration as it doesn't support stdio mode properly in this context. Use native executables instead.

#### Step 2: Make Executable Available

If you used `go install` (Option A), the executable is automatically placed in `$GOPATH/bin` or `$GOBIN`, which is typically already in your PATH environment variable.

If you built from source (Option B) or downloaded the pre-built binary (Option C), ensure the `mcp-elasticsearch` executable is in your system PATH environment variable, or use the full path in your configuration.

#### Step 3: Create MCP Configuration

Create or modify the MCP configuration file for your desktop application:

**Cursor:** Create or edit `~/.cursor/mcp.json` (Linux/Mac) or `%APPDATA%\Cursor\User\mcp.json` (Windows)

**Claude Desktop:** Create or edit configuration in the appropriate location for your platform.

#### Configuration Examples

**Basic Configuration:**
```json
{
  "mcpServers": {
    "elasticsearch": {
      "command": "mcp-elasticsearch",
      "env": {
        "ES_ADDRESSES": "http://localhost:9200",
        "ES_VERSION": "8"
      }
    }
  }
}
```

**With Authentication:**
```json
{
  "mcpServers": {
    "elasticsearch": {
      "command": "mcp-elasticsearch",
      "env": {
        "ES_ADDRESSES": "https://your-elasticsearch.com:9200",
        "ES_USERNAME": "elastic",
        "ES_PASSWORD": "your-password",
        "ES_VERSION": "8",
        "ES_SSL": "true"
      }
    }
  }
}
```

**With API Key:**
```json
{
  "mcpServers": {
    "elasticsearch": {
      "command": "mcp-elasticsearch",
      "env": {
        "ES_ADDRESSES": "https://your-elasticsearch.com:9200",
        "ES_API_KEY": "your-api-key",
        "ES_VERSION": "8"
      }
    }
  }
}
```

**Elastic Cloud Configuration:**
```json
{
  "mcpServers": {
    "elasticsearch": {
      "command": "mcp-elasticsearch",
      "env": {
        "ES_CLOUD_ID": "your-cloud-id",
        "ES_USERNAME": "elastic",
        "ES_PASSWORD": "your-password",
        "ES_VERSION": "8"
      }
    }
  }
}
```

**Using Full Path (if not in PATH):**
```json
{
  "mcpServers": {
    "elasticsearch": {
      "command": "/full/path/to/mcp-elasticsearch",
      "env": {
        "ES_ADDRESSES": "http://localhost:9200",
        "ES_VERSION": "8"
      }
    }
  }
}
```

#### Step 4: Restart Desktop Application

After creating or modifying the configuration file, restart your desktop application to load the new MCP server configuration.

#### Verification

Once configured, you should see Elasticsearch tools and resources available in your desktop application. The available tools include cluster operations, index management, document operations, search capabilities, and bulk operations as listed in the "Supported Tools" section above.

## Docker Usage Examples

### Basic Stdio Mode (for LLM tool integration)
```bash
docker run -it --rm \
  -e ES_ADDRESSES=http://host.docker.internal:9200 \
  ghcr.io/aeazer/mcp-elasticsearch:latest
```

### HTTP Server Mode (for n8n, API access)
```bash
docker run -d -p 8080:8080 \
  --name mcp-elasticsearch \
  -e MCP_PROTOCOL=http \
  -e ES_ADDRESSES=http://host.docker.internal:9200 \
  ghcr.io/aeazer/mcp-elasticsearch:latest

# Test the server endpoints
curl http://localhost:8080/health    # Health check
curl http://localhost:8080/mcp       # MCP endpoint (requires proper MCP client)
```

### With Elastic Cloud
```bash
docker run -d -p 8080:8080 \
  -e MCP_PROTOCOL=http \
  -e ES_CLOUD_ID="your-cloud-id" \
  -e ES_USERNAME=elastic \
  -e ES_PASSWORD="your-password" \
  -e ES_VERSION=8 \
  ghcr.io/aeazer/mcp-elasticsearch:latest
```

### Using Docker Compose
Create a `docker-compose.yml` file:

```yaml
version: '3.8'
services:
  mcp-elasticsearch:
    image: ghcr.io/aeazer/mcp-elasticsearch:latest
    ports:
      - "8080:8080"
    environment:
      - MCP_PROTOCOL=http
      - ES_ADDRESSES=http://elasticsearch:9200
      - ES_VERSION=8
    depends_on:
      - elasticsearch
    
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports:
      - "9200:9200"
```

Run with: `docker-compose up -d`

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
| `MCP_PROTOCOL` | Protocol to use (`stdio`, `http`, or `sse` - deprecated) | `http` (in Docker), `stdio` (native) |
| `MCP_ADDRESS` | Streamable HTTP server address (HTTP mode only) | `0.0.0.0` (in Docker), `localhost` (native) |
| `MCP_PORT` | Streamable HTTP server port (HTTP mode only) | `8080` |

### Protocol Endpoints

Different protocols use different access methods:

#### Stdio Protocol
- **Access method**: Direct stdin/stdout communication
- **Use case**: LLM tool integration (Claude Desktop, etc.)
- **Endpoint**: N/A (direct process communication)

#### Streamable HTTP Protocol (Recommended)
- **MCP endpoint**: `http://host:port/mcp`
- **Health check**: `http://host:port/health`
- **Use case**: Remote access, n8n integration, API usage
- **Example**: `http://localhost:8080/mcp`

#### SSE Protocol (Deprecated)
- **MCP endpoint**: `http://host:port/sse`  
- **Use case**: Legacy SSE clients (not recommended)
- **Example**: `http://localhost:8080/sse`
- ⚠️ **Warning**: Deprecated, use HTTP protocol instead

## Usage Examples

### Stdio Mode (Default for native builds)
```bash
export ES_ADDRESSES=http://localhost:9200
export MCP_PROTOCOL=stdio
./mcp-elasticsearch
```

### Streamable HTTP Mode (Default for Docker)
```bash
export ES_ADDRESSES=http://localhost:9200
export MCP_PROTOCOL=http
export MCP_PORT=8080
./mcp-elasticsearch
```

### SSE Mode (Deprecated - Not Recommended)
⚠️ **WARNING**: SSE protocol is deprecated and not recommended for production use. Use Streamable HTTP instead.

```bash
export ES_ADDRESSES=http://localhost:9200
export MCP_PROTOCOL=sse
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
- Go 1.23 or higher
- Access to an Elasticsearch cluster
- Docker (optional, for containerized development)

### Building
```bash
go mod download
go build -o mcp-elasticsearch main.go
```

### Testing
```bash
go test ./...
```

### Building Docker Image
```bash
docker build -t mcp-elasticsearch .
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
    "body": {
      "title": "Hello World",
      "content": "This is a test document",
      "timestamp": "2024-01-01T00:00:00Z"
    }
  }
}
```

### Advanced Search with Sorting and Field Selection
```json
{
  "tool": "es_search",
  "arguments": {
    "index": "my-index",
    "query": {
      "bool": {
        "must": [
          {"match": {"title": "Hello"}}
        ],
        "filter": [
          {"range": {"timestamp": {"gte": "2024-01-01"}}}
        ]
      }
    },
    "sort": [
      {"timestamp": {"order": "desc"}},
      {"_score": {"order": "desc"}}
    ],
    "_source": ["title", "content", "timestamp"],
    "size": 20,
    "from": 0
  }
}
```

## Health Monitoring

When running in HTTP mode, the server provides multiple endpoints:

### Health Check Endpoint
```bash
# Check server health (publicly accessible)
curl http://localhost:8080/health

# Response
{"status":"healthy","server":"elasticsearch-mcp"}
```

### MCP Protocol Endpoint
```bash
# MCP communication endpoint (requires MCP client)
# URL: http://localhost:8080/mcp
# This endpoint handles MCP protocol messages and tool calls
# Not directly accessible via simple HTTP GET requests
```

### Important Notes
- **Health endpoint** (`/health`): Simple HTTP GET for monitoring
- **MCP endpoint** (`/mcp`): For MCP protocol communication only
- **SSE endpoint** (`/sse`): Deprecated, avoid using

## Error Handling

All errors are reported within the MCP tool results with `isError: true`, allowing LLMs to see and handle errors appropriately. Protocol-level errors are reserved for exceptional conditions like missing tools or server failures.

## Troubleshooting

### Container Issues
- **Container exits immediately**: Ensure you're using HTTP protocol for Docker containers
- **Cannot connect to Elasticsearch**: Use `host.docker.internal:9200` instead of `localhost:9200` in Docker
- **Permission denied**: Check Docker daemon permissions and image access

### Network Issues
- **Connection refused**: Verify Elasticsearch is running and accessible
- **SSL errors**: Set `ES_INSECURE_SKIP_VERIFY=true` for testing with self-signed certificates

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Official MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Official MCP implementation for Go
- [Elastic](https://github.com/elastic/go-elasticsearch) - Official Elasticsearch Go client
- [Model Context Protocol](https://modelcontextprotocol.io/) - Protocol specification

<div align="center">
  <sub>Built with ❤️ for the Go community</sub>
</div>