package karmada

import (
	"context"
	"fmt"
	ns "github.com/karmada-io/dashboard/pkg/resource/namespace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func CreateNamespace(getKubernetesClient GetKubernetesClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"create_namespace",
			mcp.WithDescription("Create a namespace resources in the Karmada control-plane"),
			mcp.WithString("name", mcp.Required(), mcp.Description("name for the namespace")),
			mcp.WithBoolean("skipAutoPropagation", mcp.Required(), mcp.Description("whether propagation the namespace automatically")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			karmadaClient, err := getKubernetesClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get Karmada client: %w", err)
			}

			paramName, ok := request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter name not found")
			}

			paramSkipAutoPropagation, ok := request.Params.Arguments["skipAutoPropagation"].(bool)
			if !ok {
				return nil, fmt.Errorf("parameter name not found")
			}
			spec := &ns.NamespaceSpec{
				Name:                paramName,
				SkipAutoPropagation: paramSkipAutoPropagation,
			}
			if err = ns.CreateNamespace(spec, karmadaClient); err != nil {
				return nil, fmt.Errorf("failed to create namespace: %w", err)
			}

			return mcp.NewToolResultText("create namespace success"), nil
		}
}
