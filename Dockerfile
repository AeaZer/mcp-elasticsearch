FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o mcp-elasticsearch main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/mcp-elasticsearch .

ENV MCP_PROTOCOL=stdio
ENV ES_ADDRESSES=http://localhost:9200
ENV ES_VERSION=8

EXPOSE 8080

CMD ["./mcp-elasticsearch"] 