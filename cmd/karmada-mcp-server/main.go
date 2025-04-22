package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
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
}

func initConfig() {
	// Initialize Viper configuration
	viper.SetEnvPrefix("karmada")
	viper.AutomaticEnv()
}

func runStdioServer() error {
	// Create a new MCP server
	server.NewMCPServer(
		"karmada-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
