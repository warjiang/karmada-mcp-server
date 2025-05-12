package app

import (
	"context"
	"fmt"
	"github.com/karmada-io/dashboard/pkg/client"
	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/warjiang/karmada-mcp-server/pkg/karmada"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
	"syscall"
)

func NewStdioCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stdio",
		Short: "Start stdio server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		Run: func(_ *cobra.Command, _ []string) {
			if err := runStdioServer(); err != nil {

			}
		},
	}
	return cmd
}

func runStdioServer() error {
	klog.Info("Starting mcp server in stdio mode")
	// 4 3 2 1
	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	hooks := &server.Hooks{}

	// Create a new MCP server
	karmadaServer := karmada.NewServer("version", server.WithHooks(hooks))

	karmadaClient := client.InClusterKarmadaClient()
	getClient := func(_ context.Context) (karmadaclientset.Interface, error) {
		return karmadaClient, nil // closing over client
	}

	k8sClient := client.InClusterClientForKarmadaAPIServer()
	getKubernetesClient := func(_ context.Context) (kubernetes.Interface, error) {
		return k8sClient, nil // closing over client
	}

	{
		tool, handler := karmada.ListClusters(getClient)
		karmadaServer.AddTool(tool, handler)
	}

	{
		tool, handler := karmada.CreateNamespace(getKubernetesClient)
		karmadaServer.AddTool(tool, handler)
	}

	stdioServer := server.NewStdioServer(karmadaServer)

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)

		errC <- stdioServer.Listen(ctx, in, out)
	}()

	// Output karmada-mcp-server string
	_, _ = fmt.Fprintf(os.Stderr, "Karmada MCP Server running on stdio\n")

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		fmt.Errorf("shutting down server...")
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}
