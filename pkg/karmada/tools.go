package karmada

import (
	"context"
	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

type GetClientFn func(context.Context) (karmadaclientset.Interface, error)
type GetKubernetesClientFn func(context.Context) (kubernetes.Interface, error)

func InitToolsets(passedToolsets []string, getClient GetClientFn) error {
	return nil
}
