package stdio

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/warjiang/karmada-mcp-server/pkg/environment"
	"github.com/warjiang/karmada-mcp-server/pkg/karmada"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
	"syscall"
)

func NewStdioCommand() *cobra.Command {
	opts := newStdioServerOptions()
	cmd := &cobra.Command{
		Use:   "stdio",
		Short: "Start stdio server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			stdioServerConfig := StdioServerOptions{
				Version:         environment.Version(),
				EnabledToolsets: opts.EnabledToolsets,
				ReadOnly:        opts.ReadOnly,
			}
			return runStdioServer(stdioServerConfig)
		},
	}
	opts.AddFlags(cmd.Flags())
	return cmd
}

func runStdioServer(opts StdioServerOptions) error {
	klog.Info("Starting mcp server in stdio mode")

	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	karmadaServer, err := karmada.NewMCPServer(karmada.MCPServerConfig{
		Version:         opts.Version,
		EnabledToolsets: opts.EnabledToolsets,
		ReadOnly:        opts.ReadOnly,
	})
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	stdioServer := server.NewStdioServer(karmadaServer)

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		klog.Info("mcp server in stdio mode started")
		errC <- stdioServer.Listen(ctx, os.Stdin, os.Stdout)
	}()

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		klog.Info("shutting down stdio server...")
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}
