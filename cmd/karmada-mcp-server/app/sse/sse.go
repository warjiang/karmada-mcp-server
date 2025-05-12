package sse

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

func NewSseCommand() *cobra.Command {
	opts := newSseServerOptions()
	cmd := &cobra.Command{
		Use:   "sse",
		Short: "Start sse server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			sseServerConfig := SseServerOptions{
				Version:         environment.Version(),
				EnabledToolsets: opts.EnabledToolsets,
				ReadOnly:        opts.ReadOnly,
			}
			return runSseServer(sseServerConfig)
		},
	}

	return cmd
}

func runSseServer(opts SseServerOptions) error {
	klog.Info("Starting mcp server in sse mode")

	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx = ctx

	karmadaServer, err := karmada.NewMCPServer(karmada.MCPServerConfig{
		Version:         opts.Version,
		EnabledToolsets: opts.EnabledToolsets,
		ReadOnly:        opts.ReadOnly,
	})
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	sseServer := server.NewSSEServer(karmadaServer,
		server.WithBaseURL("http://localhost:5173"),
		server.WithStaticBasePath("/mcp"),
	)
	sseServer.Start("localhost:1234")

	return nil
}
