package karmada

import (
	"context"
	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	toolsets "github.com/warjiang/karmada-mcp-server/pkg/toolset"
	"k8s.io/client-go/kubernetes"
)

type GetClientFn func(context.Context) (karmadaclientset.Interface, error)

type GetKubernetesClientFn func(context.Context) (kubernetes.Interface, error)

var DefaultTools = []string{"all"}

func InitToolsetGroup(passedToolsets []string, readOnly bool, getClient GetClientFn) (*toolsets.ToolsetGroup, error) {
	// Create a new toolset group
	tsg := toolsets.NewToolsetGroup(readOnly)

	// Define all available features with their default state (disabled)
	// Create toolsets
	clusters := toolsets.NewToolset("cluster", "Karmada cluster related tools").
		AddReadTools().
		AddWriteTools()
	policies := toolsets.NewToolset("policy", "Karmada policy related tools").
		AddReadTools().
		AddWriteTools()
	// Add toolsets to the group
	tsg.AddToolset(clusters)
	tsg.AddToolset(policies)

	// Enable the requested features
	if err := tsg.EnableToolsets(passedToolsets); err != nil {
		return nil, err
	}

	return tsg, nil
}
