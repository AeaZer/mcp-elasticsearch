# Elasticsearch MCP 服务器

*其他语言版本: [English](README.md), [中文](README_zh.md)*

## 概述

基于 [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) 构建的 Elasticsearch MCP (Model Context Protocol) 服务器，无缝集成 Elasticsearch 7、8、9 版本。

## 功能特性

- 🔗 **多协议支持**: 支持 stdio 和 Streamable HTTP 协议
- 📊 **多版本兼容**: 兼容 Elasticsearch 7、8、9 版本
- ⚙️ **环境变量配置**: 通过环境变量进行配置
- 🔧 **丰富工具集**: 完整的 Elasticsearch 操作工具
- 🌐 **生产就绪**: 支持 Docker，优化构建

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
- `es_search`: 执行搜索查询，支持过滤和排序

### 批量操作
- `es_bulk`: 在单个请求中执行多个操作

## 快速开始

选择以下任一方式运行 Elasticsearch MCP 服务器：

### 方式一：构建 Docker 镜像（推荐）

```bash
# 克隆仓库
git clone https://github.com/AeaZer/mcp-elasticsearch.git
cd mcp-elasticsearch

# 构建 Docker 镜像
docker build -t mcp-elasticsearch .

# 运行容器
docker run -e ES_ADDRESSES=http://localhost:9200 -e ES_VERSION=8 mcp-elasticsearch
```

### 方式二：使用预构建镜像（即将推出）

```bash
# 当镜像发布到仓库后将可用
# docker run -e ES_ADDRESSES=http://localhost:9200 ghcr.io/aeazer/mcp-elasticsearch:latest
```

*注意：预构建镜像尚未发布。请使用方式一或方式三。*

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
| `MCP_PROTOCOL` | 使用的协议（`stdio` 或 `http`） | `stdio` |
| `MCP_ADDRESS` | Streamable HTTP 服务器地址（仅 HTTP 模式） | `localhost` |
| `MCP_PORT` | Streamable HTTP 服务器端口（仅 HTTP 模式） | `8080` |

## 使用示例

### Stdio 模式（默认）
```bash
export ES_ADDRESSES=http://localhost:9200
export MCP_PROTOCOL=stdio
./mcp-elasticsearch
```

### Streamable HTTP 模式
```bash
export ES_ADDRESSES=http://localhost:9200
export MCP_PROTOCOL=http
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
- Go 1.21 或更高版本
- 可访问的 Elasticsearch 集群

### 构建
```bash
go mod download
go build -o mcp-elasticsearch main.go
```

### 测试
```bash
go test ./...
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
    "document": {
      "title": "你好世界",
      "content": "这是一个测试文档",
      "timestamp": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 搜索文档
```json
{
  "tool": "es_search",
  "arguments": {
    "index": "my-index",
    "query": {
      "match": {
        "title": "你好"
      }
    },
    "size": 10
  }
}
```

## 错误处理

所有错误都在 MCP 工具结果中报告，设置 `isError: true`，允许 LLM 看到并适当处理错误。协议级别的错误仅用于异常情况，如缺少工具或服务器故障。

## 贡献

1. Fork 此仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 许可证

此项目使用 MIT 许可证 - 有关详细信息，请参阅 [LICENSE](LICENSE) 文件。

## 致谢

- [Mark3Labs MCP-Go](https://github.com/mark3labs/mcp-go) - Go 的 MCP 实现
- [Elastic](https://github.com/elastic/go-elasticsearch) - 官方 Elasticsearch Go 客户端
- [Model Context Protocol](https://modelcontextprotocol.io/) - 协议规范 