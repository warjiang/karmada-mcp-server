package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/warjiang/karmada-mcp-server/cmd/karmada-mcp-server/app"
	"github.com/warjiang/karmada-mcp-server/pkg/environment"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:     "server",
		Short:   "Karmada MCP Server",
		Long:    `A Karmada MCP server that handles various tools and resources.`,
		Version: fmt.Sprintf("Version: %s\nCommit: %s\nBuild Date: %s", environment.Version(), environment.GitCommit(), environment.BuildDate()),
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
	rootCmd.AddCommand(app.StdioCmd)
	rootCmd.AddCommand(app.SseCmd)
}

func initConfig() {
	// Initialize Viper configuration
	viper.SetEnvPrefix("karmada")
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
