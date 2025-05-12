package app

import (
	"context"
	"fmt"
	"github.com/karmada-io/dashboard/pkg/client"
	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/warjiang/karmada-mcp-server/pkg/karmada"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
	"syscall"
)

func NewSseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sse",
		Short: "Start sse server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		Run: func(_ *cobra.Command, _ []string) {
			if err := runSseServer(); err != nil {

			}
		},
	}

	return cmd
}

func runSseServer() error {
	klog.Infof("Starting sse server")
	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx = ctx
	hooks := &server.Hooks{}

	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		fmt.Printf("beforeAny: %s, %v, %v\n", method, id, message)
	})
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		fmt.Printf("onSuccess: %s, %v, %v, %v\n", method, id, message, result)
	})
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		fmt.Printf("onError: %s, %v, %v, %v\n", method, id, message, err)
	})
	hooks.AddBeforeInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest) {
		fmt.Printf("beforeInitialize: %v, %v\n", id, message)
	})
	hooks.AddOnRequestInitialization(func(ctx context.Context, id any, message any) error {
		fmt.Printf("AddOnRequestInitialization: %v, %v\n", id, message)
		// authorization verification and other preprocessing tasks are performed.
		return nil
	})
	hooks.AddAfterInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest, result *mcp.InitializeResult) {
		fmt.Printf("afterInitialize: %v, %v, %v\n", id, message, result)
	})
	hooks.AddAfterCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		fmt.Printf("afterCallTool: %v, %v, %v\n", id, message, result)
	})
	hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		fmt.Printf("beforeCallTool: %v, %v\n", id, message)
	})

	// Create a new MCP server
	karmadaServer := karmada.NewServer("version",
		server.WithHooks(hooks),
		server.WithLogging(),
	)

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

	sseServer := server.NewSSEServer(karmadaServer,
		server.WithBaseURL("http://localhost:5173"),
		server.WithStaticBasePath("/mcp"),
	)
	sseServer.Start("localhost:1234")
	//
	/*
		// Start listening for messages
		errC := make(chan error, 1)
		go func() {
			//in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)
			// errC <- sseServer.Listen(ctx, in, out)
			sseServer.Start("localhost:7890")
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
		}*/

	return nil
}
