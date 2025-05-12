package karmada

import "github.com/mark3labs/mcp-go/server"

// NewServer creates a new GitHub MCP server with the specified GH client and logger.
func NewServer(version string, opts ...server.ServerOption) *server.MCPServer {
	// Add default options
	defaultOpts := []server.ServerOption{
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	}
	opts = append(defaultOpts, opts...)

	// Create a new MCP server
	s := server.NewMCPServer(
		"karmada-mcp-server",
		version,
		opts...,
	)
	return s
}
