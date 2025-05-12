package karmada

import (
	"context"
	"fmt"
	"github.com/karmada-io/dashboard/pkg/client"
	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/kubernetes"
)

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

func NewMCPServer(cfg MCPServerConfig) (*server.MCPServer, error) {
	// init karmada client
	// init kubernetes client

	hooks := &server.Hooks{}
	/*
		hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
			fmt.Printf("beforeAny: %s, %v, %v\n", method, id, message)
		})
		hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
			fmt.Printf("onSuccess: %s, %v, %v, %v\n", method, id, message, result)
		})
		hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
			fmt.Printf("onError: %s, %v, %v, %v\n", method, id, message, err)
		})
		hooks.AddBeforeInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest) {
			fmt.Printf("beforeInitialize: %v, %v\n", id, message)
		})
		hooks.AddOnRequestInitialization(func(ctx context.Context, id any, message any) error {
			fmt.Printf("AddOnRequestInitialization: %v, %v\n", id, message)
			// authorization verification and other preprocessing tasks are performed.
			return nil
		})
		hooks.AddAfterInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest, result *mcp.InitializeResult) {
			fmt.Printf("afterInitialize: %v, %v, %v\n", id, message, result)
		})
		hooks.AddAfterCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
			fmt.Printf("afterCallTool: %v, %v, %v\n", id, message, result)
		})
		hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
			fmt.Printf("beforeCallTool: %v, %v\n", id, message)
		})
	*/

	// Create karmada MCP server
	karmadaServer := NewServer(cfg.Version, server.WithHooks(hooks))

	karmadaClient := client.InClusterKarmadaClient()
	getKarmadaClient := func(_ context.Context) (karmadaclientset.Interface, error) {
		return karmadaClient, nil // closing over client
	}

	k8sClient := client.InClusterClientForKarmadaAPIServer()
	getKubernetesClient := func(_ context.Context) (kubernetes.Interface, error) {
		return k8sClient, nil // closing over client
	}

	enabledToolsets := cfg.EnabledToolsets
	// Create default toolsets
	toolsets, err := InitToolsetGroup(
		enabledToolsets,
		cfg.ReadOnly,
		getKarmadaClient, getKubernetesClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize toolsets: %w", err)
	}

	// Register the tools with the server
	toolsets.RegisterTools(karmadaServer)

	return karmadaServer, nil
}
