package karmada

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/karmada-io/dashboard/pkg/dataselect"
	"github.com/karmada-io/dashboard/pkg/resource/common"
	"github.com/karmada-io/dashboard/pkg/resource/deployment"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

func CreateDeployment(getKubernetesClient GetKubernetesClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"create_deployment",
			mcp.WithDescription("Create a deployment resources in the Karmada control-plane"),
			mcp.WithString("name", mcp.Required(), mcp.Description("name for deployment")),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("namespace for deployment")),
			mcp.WithString("content", mcp.Required(), mcp.Description("deployment content which in form of yaml")),
			//mcp.WithBoolean("dry-run", mcp.DefaultBool(true), mcp.Description("dry-run options, if true the resource will not be applied to apiserver directly")),
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

			paramNamespace, ok := request.Params.Arguments["namespace"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter namespace not found")
			}

			paramContent, ok := request.Params.Arguments["content"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter content not found")
			}
			deployment := appsv1.Deployment{}
			if err = yaml.Unmarshal([]byte(paramContent), &deployment); err != nil {
				klog.Errorf("unmarshal deployment error: %v", err)
				return nil, err
			}
			deployment.Name = paramName

			createResp, err := karmadaClient.AppsV1().Deployments(paramNamespace).Create(ctx, &deployment, metav1.CreateOptions{})
			if err != nil {
				klog.Errorf("create deployment error: %v", err)
				return nil, err
			}

			respBuff, err := json.Marshal(createResp)
			if err != nil {
				klog.Errorf("marshal created deployment error: %v", err)
				return nil, err
			}
			return mcp.NewToolResultText(string(respBuff)), nil
		}
}

func ListDeployment(getKubernetesClient GetKubernetesClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"list_deployment",
			mcp.WithDescription("List deployments under the specific namespace in the Karmada control-plane"),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("name of namespace")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			karmadaClient, err := getKubernetesClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get Karmada client: %w", err)
			}

			paramNamespace, ok := request.Params.Arguments["namespace"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter namespace not found")
			}
			namespace := common.NewNamespaceQuery([]string{paramNamespace})
			resp, err := deployment.GetDeploymentList(karmadaClient, namespace, dataselect.NoDataSelect)
			if err != nil {
				klog.Errorf("failed to list deployments, err: %v", err)
				return nil, err
			}
			deployList := make([]string, 0)
			for _, deployment := range resp.Deployments {
				deployList = append(deployList, deployment.ObjectMeta.Name)
			}
			r, err := json.Marshal(map[string]interface{}{
				"deployments": deployList,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to marshal deployments: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
