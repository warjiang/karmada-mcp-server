package main

import (
	"fmt"
	"github.com/karmada-io/dashboard/pkg/client"
	"github.com/karmada-io/karmada/pkg/sharedcli/klogflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/warjiang/karmada-mcp-server/cmd/karmada-mcp-server/app"
	"github.com/warjiang/karmada-mcp-server/pkg/environment"
	cliflag "k8s.io/component-base/cli/flag"
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
	cobra.OnInitialize(initConfig, initKarmada)
	rootCmd.SetVersionTemplate("{{.Short}}\n{{.Version}}\n")

	// Add global flags that will be shared by all commands
	rootCmd.PersistentFlags().String("karmada-kubeconfig", "", "Path to the karmada control plane kubeconfig file.")
	rootCmd.PersistentFlags().String("karmada-context", "", "The name of the karmada-kubeconfig context to use.")
	rootCmd.PersistentFlags().Bool("skip-karmada-apiserver-tls-verify", false, "enable if connection with remote Karmada API server should skip TLS verify")

	fss := cliflag.NamedFlagSets{}
	logsFlagSet := fss.FlagSet("logs")
	klogflag.Add(logsFlagSet)
	rootCmd.PersistentFlags().AddFlagSet(logsFlagSet)

	// Bind flag to viper
	_ = viper.BindPFlag("karmada-kubeconfig", rootCmd.PersistentFlags().Lookup("karmada-kubeconfig"))
	_ = viper.BindPFlag("karmada-context", rootCmd.PersistentFlags().Lookup("karmada-context"))
	_ = viper.BindPFlag("skip-karmada-apiserver-tls-verify", rootCmd.PersistentFlags().Lookup("skip-karmada-apiserver-tls-verify"))

	// Add subcommands
	rootCmd.AddCommand(app.NewStdioCommand())
	rootCmd.AddCommand(app.NewSseCommand())

}

func initConfig() {
	// Initialize Viper configuration
	viper.SetEnvPrefix("karmada")
	viper.AutomaticEnv()
}

func initKarmada() {
	// Initialize client for karmada apiserver
	client.InitKarmadaConfig(
		client.WithKubeconfig(viper.GetString("karmada-kubeconfig")),
		client.WithKubeContext(viper.GetString("karmada-context")),
		client.WithInsecureTLSSkipVerify(viper.GetBool("skip-karmada-apiserver-tls-verify")),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
