package karmada

import (
	"context"
	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
)

type GetClientFn func(context.Context) (karmadaclientset.Interface, error)

func InitToolsets(passedToolsets []string, getClient GetClientFn) error {
	return nil
}
