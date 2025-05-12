package stdio

import (
	"github.com/spf13/pflag"
	"github.com/warjiang/karmada-mcp-server/pkg/karmada"
)

type StdioServerOptions struct {
	// Version of the server
	Version string

	// EnabledToolsets is a list of toolsets to enable
	// See: https://github.com/github/github-mcp-server?tab=readme-ov-file#tool-configuration
	EnabledToolsets []string

	// ReadOnly indicates if we should only register read-only tools
	ReadOnly bool
}

// newStdioServerOptions returns initialized StdioServerOptions.
func newStdioServerOptions() *StdioServerOptions {
	return &StdioServerOptions{}
}

// AddFlags adds flags of api to the specified FlagSet
func (o *StdioServerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringSliceVar(&o.EnabledToolsets, "toolsets", karmada.DefaultTools, "An optional comma separated list of groups of tools to allow, defaults to enabling all")
	fs.BoolVar(&o.ReadOnly, "read-only", false, "Restrict the server to read-only operations")
}
