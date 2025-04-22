package karmada

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/karmada-io/dashboard/pkg/dataselect"
	"github.com/karmada-io/dashboard/pkg/resource/cluster"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ListClusters(getClient GetClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"list_clusters",
			mcp.WithDescription("List all clusters in the Karmada control plane."),
		),

		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			karmadaClient, err := getClient(ctx)
			ds := dataselect.NoDataSelect
			if err != nil {
				return nil, fmt.Errorf("failed to get Karmada client: %w", err)
			}

			result, err := cluster.GetClusterList(karmadaClient, ds)
			if err != nil {
				return nil, fmt.Errorf("failed to get cluster list: %w", err)
			}

			clusters := make([]string, 0)
			for _, c := range result.Clusters {
				clusters = append(clusters, c.ObjectMeta.Name)
			}

			r, err := json.Marshal(map[string]interface{}{
				"clusters": clusters,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to marshal user: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
