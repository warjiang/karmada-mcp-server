package main

import (
	"context"
	"fmt"
	"github.com/karmada-io/dashboard/pkg/client"
	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
