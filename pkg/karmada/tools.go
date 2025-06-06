package karmada

import (
	"context"
	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	toolsets "github.com/warjiang/karmada-mcp-server/pkg/toolset"
	"k8s.io/client-go/kubernetes"
)

type GetKarmadaClientFn func(context.Context) (karmadaclientset.Interface, error)

type GetKubernetesClientFn func(context.Context) (kubernetes.Interface, error)

var DefaultTools = []string{"all"}

func InitToolsetGroup(passedToolsets []string, readOnly bool, getKarmadaClient GetKarmadaClientFn, getKubernetesClient GetKubernetesClientFn) (*toolsets.ToolsetGroup, error) {
	// Create a new toolset group
	tsg := toolsets.NewToolsetGroup(readOnly)

	// Define all available features with their default state (disabled)
	// Create toolsets
	clusters := toolsets.NewToolset("cluster", "Karmada cluster related tools").
		AddReadTools(
			toolsets.NewServerTool(ListClusters(getKarmadaClient)),
		).
		AddWriteTools()
	policies := toolsets.NewToolset("policy", "Karmada policy related tools").
		AddReadTools(
			toolsets.NewServerTool(ListPropagationPolicy(getKarmadaClient, getKubernetesClient)),
			toolsets.NewServerTool(GetPropagationPolicy(getKarmadaClient)),
		).
		AddWriteTools(
			toolsets.NewServerTool(CreatePropagationPolicy(getKarmadaClient)),
			toolsets.NewServerTool(DeletePropagationPolicy(getKarmadaClient)),
		)
	resources := toolsets.NewToolset("resource", "Karmada resource related tools").
		AddReadTools(
			toolsets.NewServerTool(ListNamespace(getKubernetesClient)),
			toolsets.NewServerTool(ListDeployment(getKubernetesClient)),
		).
		AddWriteTools(
			toolsets.NewServerTool(CreateNamespace(getKubernetesClient)),
			toolsets.NewServerTool(CreateDeployment(getKubernetesClient)),
			toolsets.NewServerTool(DeleteUnstructuredResource()),
		)
	// Add toolsets to the group
	tsg.AddToolset(clusters)
	tsg.AddToolset(policies)
	tsg.AddToolset(resources)

	// Enable the requested features
	if err := tsg.EnableToolsets(passedToolsets); err != nil {
		return nil, err
	}

	return tsg, nil
}
