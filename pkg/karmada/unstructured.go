package karmada

import (
	"context"
	"fmt"
	"github.com/karmada-io/dashboard/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
)

func DeleteUnstructuredResource() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"delete_unstructured_resource",
			mcp.WithDescription("Delete unstructured resources in the Karmada control-plane"),
			mcp.WithString("namespace",
				mcp.Description("namespace for scoped resources, only required for namespace-scoped resources"),
			),
			mcp.WithString("kind",
				mcp.Required(),
				mcp.Description("resource kind in lower case"),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("resources name"),
			),
			mcp.WithBoolean("deleteNow",
				mcp.DefaultBool(true),
				mcp.Description("whether waiting for resources be deleted successfully")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// todo update by karmada-dashboard
			verber, err := client.VerberClient(nil)

			paramNamespace, _ := request.Params.Arguments["namespace"].(string)
			paramKind, ok := request.Params.Arguments["kind"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter kind not found")
			}
			paramName, ok := request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter name not found")
			}
			paramDeleteNow, ok := request.Params.Arguments["deleteNow"].(bool)
			if !ok {
				return nil, fmt.Errorf("parameter deleteNow not found")
			}

			if err = verber.Delete(paramKind, paramNamespace, paramName, paramDeleteNow); err != nil {
				klog.ErrorS(err, "Failed to delete resource")
				errMsg := ""
				if paramNamespace != "" {
					errMsg = fmt.Sprintf("Karmada: failed to delete %s/%s %s-resource", paramNamespace, paramName, paramKind)
				} else {
					errMsg = fmt.Sprintf("Karmada: failed to delete %s %s-resource", paramName, paramKind)
				}
				return mcp.NewToolResultText(errMsg), err
			}

			err = retry.OnError(
				retry.DefaultRetry,
				func(err error) bool {
					return errors.IsNotFound(err)
				},
				func() error {
					_, getErr := verber.Get(paramKind, paramNamespace, paramName)
					return getErr
				})
			if !errors.IsNotFound(err) {
				klog.ErrorS(err, "Wait for verber delete resource failed")
				return mcp.NewToolResultText("Wait for verber delete resource failed"), err
			}

			return mcp.NewToolResultText("delete resource success"), nil
		}
}
