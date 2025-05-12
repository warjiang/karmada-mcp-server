package main

import (
	"context"
	"fmt"
	"github.com/karmada-io/dashboard/pkg/client"
	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/warjiang/karmada-mcp-server/pkg/karmada"
	"io"
	"k8s.io/client-go/kubernetes"
	"os"
	"os/signal"
	"syscall"
)

var version = "version"
var commit = "commit"
var date = "date"

var (
	rootCmd = &cobra.Command{
		Use:     "server",
		Short:   "Karmada MCP Server",
		Long:    `A Karmada MCP server that handles various tools and resources.`,
		Version: fmt.Sprintf("Version: %s\nCommit: %s\nBuild Date: %s", version, commit, date),
	}
	stdioCmd = &cobra.Command{
		Use:   "stdio",
		Short: "Start stdio server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		Run: func(_ *cobra.Command, _ []string) {
			if err := runStdioServer(); err != nil {

			}
		},
	}
	sseCmd = &cobra.Command{
		Use:   "sse",
		Short: "Start sse server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		Run: func(_ *cobra.Command, _ []string) {
			if err := runSseServer(); err != nil {

			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SetVersionTemplate("{{.Short}}\n{{.Version}}\n")

	// Add global flags that will be shared by all commands
	rootCmd.PersistentFlags().String("karmada-kubeconfig", "", "Path to the karmada control plane kubeconfig file.")
	rootCmd.PersistentFlags().String("karmada-context", "", "The name of the karmada-kubeconfig context to use.")
	rootCmd.PersistentFlags().Bool("skip-karmada-apiserver-tls-verify", false, "enable if connection with remote Karmada API server should skip TLS verify")

	// Bind flag to viper
	_ = viper.BindPFlag("karmada-kubeconfig", rootCmd.PersistentFlags().Lookup("karmada-kubeconfig"))
	_ = viper.BindPFlag("karmada-context", rootCmd.PersistentFlags().Lookup("karmada-context"))
	_ = viper.BindPFlag("skip-karmada-apiserver-tls-verify", rootCmd.PersistentFlags().Lookup("skip-karmada-apiserver-tls-verify"))

	// Add subcommands
	rootCmd.AddCommand(stdioCmd)
	rootCmd.AddCommand(sseCmd)
}

func initConfig() {
	// Initialize Viper configuration
	viper.SetEnvPrefix("karmada")
	viper.AutomaticEnv()
}

func runStdioServer() error {
	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	hooks := &server.Hooks{}

	// Create a new MCP server
	karmadaServer := karmada.NewServer(version, server.WithHooks(hooks))

	// init client for karmada apiserver
	client.InitKarmadaConfig(
		client.WithKubeconfig(viper.GetString("karmada-kubeconfig")),
		client.WithKubeContext(viper.GetString("karmada-context")),
		client.WithInsecureTLSSkipVerify(viper.GetBool("skip-karmada-apiserver-tls-verify")),
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

func runSseServer() error {
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
	karmadaServer := karmada.NewServer(version,
		server.WithHooks(hooks),
		server.WithLogging(),
	)

	// init client for karmada apiserver
	client.InitKarmadaConfig(
		client.WithKubeconfig(viper.GetString("karmada-kubeconfig")),
		client.WithKubeContext(viper.GetString("karmada-context")),
		client.WithInsecureTLSSkipVerify(viper.GetBool("skip-karmada-apiserver-tls-verify")),
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
