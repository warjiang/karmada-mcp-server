package karmada

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/karmada-io/dashboard/pkg/dataselect"
	"github.com/karmada-io/dashboard/pkg/resource/common"
	"github.com/karmada-io/dashboard/pkg/resource/propagationpolicy"
	"github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

func CreatePropagationPolicy(getKarmadaClient GetKarmadaClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"create_propagationpolicy",
			mcp.WithDescription("Create a propagationpolicy resources in the Karmada control-plane"),
			mcp.WithString("name", mcp.Required(), mcp.Description("name for propagationpolicy")),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("namespace for propagationpolicy")),
			mcp.WithString("content", mcp.Required(), mcp.Description(`propagationpolicy content which in form of yaml, one propagationpolicy yaml file likes:
apiVersion: policy.karmada.io/v1alpha1
kind: PropagationPolicy
metadata:
  name: nginx-propagation
spec:
  resourceSelectors:
    - apiVersion: apps/v1
      kind: Deployment
      name: nginx
  placement:
    clusterAffinity:
      clusterNames:
        - member1
        - member2
    replicaScheduling:
      replicaDivisionPreference: Weighted
      replicaSchedulingType: Divided
      weightPreference:
        staticWeightList:
          - targetCluster:
              clusterNames:
                - member1
            weight: 1
          - targetCluster:
              clusterNames:
                - member2
            weight: 1
`)),
			//mcp.WithBoolean("dry-run", mcp.DefaultBool(true), mcp.Description("dry-run options, if true the resource will not be applied to apiserver directly")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			karmadaClient, err := getKarmadaClient(ctx)
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
			propagationPolicy := v1alpha1.PropagationPolicy{}
			if err = yaml.Unmarshal([]byte(paramContent), &propagationPolicy); err != nil {
				klog.Errorf("unmarshal propagationpolicy error: %v", err)
				return nil, err
			}
			propagationPolicy.Name = paramName

			createResp, err := karmadaClient.PolicyV1alpha1().PropagationPolicies(paramNamespace).Create(ctx, &propagationPolicy, metav1.CreateOptions{})
			if err != nil {
				klog.Errorf("create propagationpolicy error: %v", err)
				return nil, err
			}

			respBuff, err := json.Marshal(createResp)
			if err != nil {
				klog.Errorf("marshal created propagationpolicy error: %v", err)
				return nil, err
			}
			return mcp.NewToolResultText(string(respBuff)), nil
		}
}

func ListPropagationPolicy(getKarmadaClient GetKarmadaClientFn, getKubernetesClient GetKubernetesClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"list_propagationpolicy",
			mcp.WithDescription("List propagationpolicies under the specific namespace in the Karmada control-plane"),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("name of namespace")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			karmadaClient, err := getKarmadaClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get Karmada client: %v", err)
			}

			kubernetesClient, err := getKubernetesClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get Kubernetes client: %v", err)
			}

			paramNamespace, ok := request.Params.Arguments["namespace"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter namespace not found")
			}

			namespace := common.NewNamespaceQuery([]string{paramNamespace})
			dataSelect := dataselect.NoDataSelect
			resp, err := propagationpolicy.GetPropagationPolicyList(karmadaClient, kubernetesClient, namespace, dataSelect)

			if err != nil {
				klog.Errorf("failed to list propagationpolicies, err: %v", err)
				return nil, err
			}
			propagationPolicyList := make([]string, 0)
			for _, propagationPolicy := range resp.PropagationPolicys {
				propagationPolicyList = append(propagationPolicyList, propagationPolicy.ObjectMeta.Name)
			}
			r, err := json.Marshal(map[string]interface{}{
				"propagationPolicies": propagationPolicyList,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to marshal propagationpolicies: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func GetPropagationPolicy(getKarmadaClient GetKarmadaClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_propagationpolicy",
			mcp.WithDescription("Get propagationpolicy detailed yaml under the specific namespace in the Karmada control-plane"),
			mcp.WithString("name", mcp.Required(), mcp.Description("name for propagationpolicy")),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("name of namespace")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			karmadaClient, err := getKarmadaClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get Karmada client: %v", err)
			}

			paramName, ok := request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter name not found")
			}

			paramNamespace, ok := request.Params.Arguments["namespace"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter namespace not found")
			}

			resp, err := propagationpolicy.GetPropagationPolicyDetail(karmadaClient, paramNamespace, paramName)
			if err != nil {
				klog.Errorf("failed to get propagationpolicy, err: %v", err)
				return nil, err
			}
			r, err := json.Marshal(resp)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal propagationpolicy: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func DeletePropagationPolicy(getKarmadaClient GetKarmadaClientFn) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"delete_propagationpolicy",
			mcp.WithDescription("Delete propagationpolicy under the specific namespace in the Karmada control-plane"),
			mcp.WithString("name", mcp.Required(), mcp.Description("name for propagationpolicy")),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("name of namespace")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			karmadaClient, err := getKarmadaClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get Karmada client: %v", err)
			}

			paramName, ok := request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter name not found")
			}

			paramNamespace, ok := request.Params.Arguments["namespace"].(string)
			if !ok {
				return nil, fmt.Errorf("parameter namespace not found")
			}

			err = karmadaClient.PolicyV1alpha1().PropagationPolicies(paramNamespace).Delete(ctx, paramName, metav1.DeleteOptions{})
			if err != nil {
				klog.ErrorS(err, "Failed to delete propagationpolicy")
				return nil, err
			}

			return mcp.NewToolResultText("delete propagationpolicy success"), nil
		}
}
