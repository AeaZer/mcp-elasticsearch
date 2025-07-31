# Elasticsearch MCP 服务器

*其他语言版本: [English](README.md), [中文](README_zh.md)*

## 概述

基于 [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) 构建的 Elasticsearch MCP (Model Context Protocol) 服务器，无缝集成 Elasticsearch 7、8、9 版本。

## 功能特性

- 🔗 **多协议支持**: 支持 stdio、Streamable HTTP 和 SSE 协议（SSE 已弃用）
- 📊 **多版本兼容**: 兼容 Elasticsearch 7、8、9 版本
- ⚙️ **环境变量配置**: 通过环境变量进行配置
- 🔧 **丰富工具集**: 完整的 Elasticsearch 操作工具
- 🌐 **生产就绪**: 支持 Docker，优化构建
- 🐳 **容器就绪**: 提供预构建的 Docker 镜像

## 支持的工具

### 集群操作
- `es_cluster_info`: 获取集群信息和版本详情
- `es_cluster_health`: 获取集群健康状态和指标

### 索引管理
- `es_index_create`: 创建新索引，支持设置和映射
- `es_index_delete`: 删除现有索引
- `es_index_exists`: 检查索引是否存在
- `es_index_list`: 列出所有索引及其元数据

### 文档操作
- `es_document_index`: 索引文档，支持可选 ID
- `es_document_get`: 通过 ID 检索文档
- `es_document_update`: 更新现有文档
- `es_document_delete`: 通过 ID 删除文档

### 搜索操作
- `es_search`: 执行搜索查询，支持过滤、排序和字段选择
  - 支持参数：`index`、`query`、`size`、`from`、`sort`、`_source`
  - 完整的 Elasticsearch Query DSL 支持

### 批量操作
- `es_bulk`: 在单个请求中执行多个操作

## 快速开始

选择以下任一方式运行 Elasticsearch MCP 服务器：

### 方式一：使用预构建 Docker 镜像（推荐）

```bash
# 基本用法，连接本地 Elasticsearch
docker run --rm \
  -e ES_ADDRESSES=http://localhost:9200 \
  ghcr.io/aeazer/mcp-elasticsearch:latest

# HTTP 模式用于远程访问
docker run -d -p 8080:8080 \
  -e MCP_PROTOCOL=http \
  -e ES_ADDRESSES=http://your-elasticsearch:9200 \
  ghcr.io/aeazer/mcp-elasticsearch:latest

# 带认证的用法
docker run -d -p 8080:8080 \
  -e MCP_PROTOCOL=http \
  -e ES_ADDRESSES=https://your-elasticsearch:9200 \
  -e ES_USERNAME=elastic \
  -e ES_PASSWORD=your-password \
  -e ES_SSL=true \
  ghcr.io/aeazer/mcp-elasticsearch:latest
```

### 方式二：构建 Docker 镜像

```bash
# 克隆仓库
git clone https://github.com/AeaZer/mcp-elasticsearch.git
cd mcp-elasticsearch

# 构建 Docker 镜像
docker build -t mcp-elasticsearch .

# 运行容器
docker run -e ES_ADDRESSES=http://localhost:9200 -e ES_VERSION=8 mcp-elasticsearch
```

### 方式三：从源码编译

```bash
# 克隆仓库
git clone https://github.com/AeaZer/mcp-elasticsearch.git
cd mcp-elasticsearch

# 下载依赖并构建
go mod download
go build -o mcp-elasticsearch main.go

# 使用环境变量运行
export ES_ADDRESSES=http://localhost:9200
export ES_VERSION=8
export MCP_PROTOCOL=stdio
./mcp-elasticsearch
```

## Docker 使用示例

### 基本 Stdio 模式（用于 LLM 工具集成）
```bash
docker run -it --rm \
  -e ES_ADDRESSES=http://host.docker.internal:9200 \
  ghcr.io/aeazer/mcp-elasticsearch:latest
```

### HTTP 服务器模式（用于 n8n、API 访问）
```bash
docker run -d -p 8080:8080 \
  --name mcp-elasticsearch \
  -e MCP_PROTOCOL=http \
  -e ES_ADDRESSES=http://host.docker.internal:9200 \
  ghcr.io/aeazer/mcp-elasticsearch:latest

# 测试服务器端点
curl http://localhost:8080/health    # 健康检查
curl http://localhost:8080/mcp       # MCP 端点（需要合适的 MCP 客户端）
```

### 使用 Elastic Cloud
```bash
docker run -d -p 8080:8080 \
  -e MCP_PROTOCOL=http \
  -e ES_CLOUD_ID="your-cloud-id" \
  -e ES_USERNAME=elastic \
  -e ES_PASSWORD="your-password" \
  -e ES_VERSION=8 \
  ghcr.io/aeazer/mcp-elasticsearch:latest
```

### 使用 Docker Compose
创建 `docker-compose.yml` 文件：

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

运行命令：`docker-compose up -d`

## 配置说明

所有配置通过环境变量完成：

### Elasticsearch 配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `ES_ADDRESSES` | Elasticsearch 集群地址（逗号分隔） | `http://localhost:9200` |
| `ES_USERNAME` | 基本认证用户名 | - |
| `ES_PASSWORD` | 基本认证密码 | - |
| `ES_API_KEY` | API Key 认证 | - |
| `ES_CLOUD_ID` | Elastic Cloud ID | - |
| `ES_SSL` | 启用 SSL/TLS | `false` |
| `ES_INSECURE_SKIP_VERIFY` | 跳过 SSL 证书验证 | `false` |
| `ES_TIMEOUT` | 连接超时时间 | `30s` |
| `ES_MAX_RETRIES` | 最大重试次数 | `3` |
| `ES_VERSION` | 目标 Elasticsearch 版本（7、8 或 9） | `8` |

### MCP 服务器配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `MCP_SERVER_NAME` | MCP 服务器名称 | `Elasticsearch MCP Server` |
| `MCP_SERVER_VERSION` | 服务器版本 | `1.0.0` |
| `MCP_PROTOCOL` | 使用的协议（`stdio`、`http` 或 `sse` - 已弃用） | `http`（Docker 中），`stdio`（本地） |
| `MCP_ADDRESS` | Streamable HTTP 服务器地址（仅 HTTP 模式） | `0.0.0.0`（Docker 中），`localhost`（本地） |
| `MCP_PORT` | Streamable HTTP 服务器端口（仅 HTTP 模式） | `8080` |

### 协议端点

不同协议使用不同的访问方式：

#### Stdio 协议
- **访问方式**: 直接的 stdin/stdout 通信
- **使用场景**: LLM 工具集成（Claude Desktop 等）
- **端点**: 无（直接进程通信）

#### Streamable HTTP 协议（推荐）
- **MCP 端点**: `http://host:port/mcp`
- **健康检查**: `http://host:port/health`
- **使用场景**: 远程访问、n8n 集成、API 使用
- **示例**: `http://localhost:8080/mcp`

#### SSE 协议（已弃用）
- **MCP 端点**: `http://host:port/sse`  
- **使用场景**: 传统 SSE 客户端（不推荐）
- **示例**: `http://localhost:8080/sse`
- ⚠️ **警告**: 已弃用，请使用 HTTP 协议

## 使用示例

### Stdio 模式（本地构建默认）
```bash
export ES_ADDRESSES=http://localhost:9200
export MCP_PROTOCOL=stdio
./mcp-elasticsearch
```

### Streamable HTTP 模式（Docker 默认）
```bash
export ES_ADDRESSES=http://localhost:9200
export MCP_PROTOCOL=http
export MCP_PORT=8080
./mcp-elasticsearch
```

### SSE 模式（已弃用 - 不建议使用）
⚠️ **警告**：SSE 协议已弃用，不建议在生产环境中使用。请使用 Streamable HTTP 协议。

```bash
export ES_ADDRESSES=http://localhost:9200
export MCP_PROTOCOL=sse
export MCP_PORT=8080
./mcp-elasticsearch
```

### 使用 Elastic Cloud
```bash
export ES_CLOUD_ID=your_cloud_id
export ES_USERNAME=elastic
export ES_PASSWORD=your_password
export ES_VERSION=8
./mcp-elasticsearch
```

## 开发

### 先决条件
- Go 1.23 或更高版本
- 可访问的 Elasticsearch 集群
- Docker（可选，用于容器化开发）

### 构建
```bash
go mod download
go build -o mcp-elasticsearch main.go
```

### 测试
```bash
go test ./...
```

### 构建 Docker 镜像
```bash
docker build -t mcp-elasticsearch .
```

## 工具使用示例

### 获取集群信息
```json
{
  "tool": "es_cluster_info",
  "arguments": {}
}
```

### 创建索引
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

### 索引文档
```json
{
  "tool": "es_document_index",
  "arguments": {
    "index": "my-index",
    "id": "doc1",
    "body": {
      "title": "你好世界",
      "content": "这是一个测试文档",
      "timestamp": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 高级搜索（带排序和字段选择）
```json
{
  "tool": "es_search",
  "arguments": {
    "index": "my-index",
    "query": {
      "bool": {
        "must": [
          {"match": {"title": "你好"}}
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

## 健康监控

在 HTTP 模式下运行时，服务器提供多个端点：

### 健康检查端点
```bash
# 检查服务器健康状态（公开访问）
curl http://localhost:8080/health

# 响应
{"status":"healthy","server":"elasticsearch-mcp"}
```

### MCP 协议端点
```bash
# MCP 通信端点（需要 MCP 客户端）
# URL: http://localhost:8080/mcp
# 此端点处理 MCP 协议消息和工具调用
# 无法通过简单的 HTTP GET 请求直接访问
```

### 重要说明
- **健康端点** (`/health`)：用于监控的简单 HTTP GET 请求
- **MCP 端点** (`/mcp`)：仅用于 MCP 协议通信
- **SSE 端点** (`/sse`)：已弃用，避免使用

## 错误处理

所有错误都在 MCP 工具结果中报告，设置 `isError: true`，允许 LLM 看到并适当处理错误。协议级别的错误仅用于异常情况，如缺少工具或服务器故障。

## 故障排除

### 容器问题
- **容器立即退出**：确保在 Docker 容器中使用 HTTP 协议
- **无法连接到 Elasticsearch**：在 Docker 中使用 `host.docker.internal:9200` 而不是 `localhost:9200`
- **权限被拒绝**：检查 Docker 守护进程权限和镜像访问权限

### 网络问题
- **连接被拒绝**：验证 Elasticsearch 是否正在运行且可访问
- **SSL 错误**：对于自签名证书，设置 `ES_INSECURE_SKIP_VERIFY=true` 进行测试

## 贡献

1. Fork 此仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 许可证

此项目使用 MIT 许可证 - 有关详细信息，请参阅 [LICENSE](LICENSE) 文件。

## 致谢

- [官方 MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - 官方 Go MCP 实现
- [Elastic](https://github.com/elastic/go-elasticsearch) - 官方 Elasticsearch Go 客户端
- [Model Context Protocol](https://modelcontextprotocol.io/) - 协议规范

<div align="center">
  <sub>用 ❤️ 为 Go 社区构建</sub>
</div>