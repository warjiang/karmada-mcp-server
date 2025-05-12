package app

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/warjiang/karmada-mcp-server/pkg/environment"
	"github.com/warjiang/karmada-mcp-server/pkg/karmada"
	"io"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
	"syscall"
)

type StdioServerConfig struct {
	// Version of the server
	Version string

	// EnabledToolsets is a list of toolsets to enable
	// See: https://github.com/github/github-mcp-server?tab=readme-ov-file#tool-configuration
	EnabledToolsets []string
	// ReadOnly indicates if we should only register read-only tools
	ReadOnly bool
}

func NewStdioCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stdio",
		Short: "Start stdio server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			// If you're wondering why we're not using viper.GetStringSlice("toolsets"),
			// it's because viper doesn't handle comma-separated values correctly for env
			// vars when using GetStringSlice.
			// https://github.com/spf13/viper/issues/380
			var enabledToolsets []string
			if err := viper.UnmarshalKey("toolsets", &enabledToolsets); err != nil {
				return fmt.Errorf("failed to unmarshal toolsets: %w", err)
			}
			stdioServerConfig := StdioServerConfig{
				Version:         environment.Version(),
				EnabledToolsets: enabledToolsets,
				ReadOnly:        viper.GetBool("read-only"),
			}
			return runStdioServer(stdioServerConfig)
		},
	}
	return cmd
}

func runStdioServer(cfg StdioServerConfig) error {
	klog.Info("Starting mcp server in stdio mode")

	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	karmadaServer, err := karmada.NewMCPServer(karmada.MCPServerConfig{
		Version:         cfg.Version,
		EnabledToolsets: cfg.EnabledToolsets,
		ReadOnly:        cfg.ReadOnly,
	})
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	stdioServer := server.NewStdioServer(karmadaServer)

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)

		errC <- stdioServer.Listen(ctx, in, out)
	}()

	// Output karmada-mcp-server string
	klog.Info("Karmada MCP Server running on stdio\n")

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		klog.Info("shutting down server...")
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}
