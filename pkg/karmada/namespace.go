package karmada

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/karmada-io/dashboard/pkg/dataselect"
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

func ListNamespace(getKubernetesClient GetKubernetesClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"list_namespace",
			mcp.WithDescription("Return all namespace resources in the Karmada control-plane"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			karmadaClient, err := getKubernetesClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get Karmada client: %w", err)
			}

			resp, err := ns.GetNamespaceList(karmadaClient, dataselect.NoDataSelect)
			if err != nil {
				return nil, fmt.Errorf("failed to list namespace: %w", err)
			}
			nsList := make([]string, 0)
			for _, namespace := range resp.Namespaces {
				nsList = append(nsList, namespace.ObjectMeta.Name)
			}

			r, err := json.Marshal(map[string]interface{}{
				"namespaces": nsList,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to marshal namespaces: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
