// Package config provides configuration management for the Elasticsearch MCP server.
// It supports loading configuration from environment variables with sensible defaults.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains the complete configuration for the Elasticsearch MCP server
type Config struct {
	// Elasticsearch connection and client configuration
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`

	// MCP server configuration
	Server ServerConfig `mapstructure:"server"`
}

// ElasticsearchConfig contains all Elasticsearch connection settings
type ElasticsearchConfig struct {
	// Addresses is a list of Elasticsearch cluster addresses
	Addresses []string `mapstructure:"addresses"`

	// Username for basic authentication
	Username string `mapstructure:"username"`

	// Password for basic authentication
	Password string `mapstructure:"password"`

	// APIKey for API key authentication (alternative to username/password)
	APIKey string `mapstructure:"api_key"`

	// CloudID for Elastic Cloud connections
	CloudID string `mapstructure:"cloud_id"`

	// SSL enables SSL/TLS connections
	SSL bool `mapstructure:"ssl"`

	// InsecureSkipVerify bypasses certificate verification (for development only)
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`

	// Timeout for HTTP requests to Elasticsearch
	Timeout time.Duration `mapstructure:"timeout"`

	// MaxRetries specifies the maximum number of retry attempts for failed requests
	MaxRetries int `mapstructure:"max_retries"`
}

// ServerConfig contains MCP server settings
type ServerConfig struct {
	// Name of the MCP server
	Name string `mapstructure:"name"`

	// Version of the MCP server
	Version string `mapstructure:"version"`

	// Protocol specifies the communication protocol (stdio, http, or sse)
	// Note: SSE protocol is deprecated and not recommended for production use
	Protocol string `mapstructure:"protocol"`

	// Address for HTTP server (only used when protocol is http)
	Address string `mapstructure:"address"`

	// Port for HTTP server (only used when protocol is http)
	Port int `mapstructure:"port"`
}

// LoadConfig loads configuration from environment variables with default values
func LoadConfig() (*Config, error) {
	config := &Config{
		Elasticsearch: ElasticsearchConfig{
			Addresses:          getEnvStringSlice("ES_ADDRESSES", []string{"http://127.0.0.1:9200"}),
			Username:           getEnvString("ES_USERNAME", ""),
			Password:           getEnvString("ES_PASSWORD", ""),
			APIKey:             getEnvString("ES_API_KEY", ""),
			CloudID:            getEnvString("ES_CLOUD_ID", ""),
			SSL:                getEnvBool("ES_SSL", false),
			InsecureSkipVerify: getEnvBool("ES_INSECURE_SKIP_VERIFY", false),
			Timeout:            getEnvDuration("ES_TIMEOUT", 30*time.Second),
			MaxRetries:         getEnvInt("ES_MAX_RETRIES", 3),
		},
		Server: ServerConfig{
			Name:     getEnvString("MCP_SERVER_NAME", "Elasticsearch MCP Server"),
			Version:  getEnvString("MCP_SERVER_VERSION", "1.0.0"),
			Protocol: getEnvString("MCP_PROTOCOL", "stdio"),
			Address:  getEnvString("MCP_ADDRESS", "localhost"),
			Port:     getEnvInt("MCP_PORT", 8080),
		},
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate checks if the configuration is valid and returns an error if not
func (c *Config) Validate() error {
	if len(c.Elasticsearch.Addresses) == 0 {
		return fmt.Errorf("at least one Elasticsearch address must be specified")
	}

	if c.Server.Protocol != "stdio" && c.Server.Protocol != "http" && c.Server.Protocol != "sse" {
		return fmt.Errorf("unsupported protocol: %s, supported protocols: stdio, http, sse (deprecated)", c.Server.Protocol)
	}

	if (c.Server.Protocol == "http" || c.Server.Protocol == "sse") && c.Server.Port <= 0 {
		return fmt.Errorf("valid port number is required for HTTP/SSE protocol")
	}

	return nil
}

// GetElasticsearchVersion returns the Elasticsearch version from environment or default
func (c *Config) GetElasticsearchVersion() string {
	// Version can be specified via environment variable, defaults to v8
	return getEnvString("ES_VERSION", "8")
}

// Environment variable helper functions

// getEnvString returns the environment variable value or the default value if not set
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvStringSlice returns a comma-separated environment variable as a slice
func getEnvStringSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")
}

// getEnvInt returns the environment variable value as an integer or the default value
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// getEnvBool returns the environment variable value as a boolean or the default value
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

// getEnvDuration returns the environment variable value as a duration or the default value
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}
