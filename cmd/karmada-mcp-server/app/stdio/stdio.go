package stdio

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/warjiang/karmada-mcp-server/pkg/environment"
	"github.com/warjiang/karmada-mcp-server/pkg/karmada"
	"io"
	"k8s.io/klog/v2"
	"os"
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

	karmadaServer, err := karmada.NewMCPServer(karmada.MCPServerConfig{
		Version:         opts.Version,
		EnabledToolsets: opts.EnabledToolsets,
		ReadOnly:        opts.ReadOnly,
	})
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	stdioServer := server.NewStdioServer(karmadaServer)
	in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)
	ctx := context.TODO()
	stdioServer.Listen(ctx, in, out)
	return nil
}
