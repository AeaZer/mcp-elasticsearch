FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o mcp-elasticsearch main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /app/mcp-elasticsearch .

# Use HTTP protocol for container environments
ENV MCP_PROTOCOL=http
ENV MCP_ADDRESS=0.0.0.0
ENV MCP_PORT=8080
ENV ES_ADDRESSES=http://localhost:9200
ENV ES_VERSION=8

EXPOSE 8080

# Health check for HTTP mode
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD if [ "$MCP_PROTOCOL" = "http" ]; then curl -f http://localhost:${MCP_PORT}/health || exit 1; else exit 0; fi

# Add labels
LABEL org.opencontainers.image.title="Elasticsearch MCP Server"
LABEL org.opencontainers.image.description="Model Context Protocol server for Elasticsearch"
LABEL org.opencontainers.image.url="https://github.com/AeaZer/mcp-elasticsearch"

CMD ["./mcp-elasticsearch"] 