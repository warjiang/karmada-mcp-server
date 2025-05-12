package karmada

type MCPServerConfig struct {
	// Version of the server
	Version string

	// EnabledToolsets is a list of toolsets to enable
	// See: https://github.com/github/github-mcp-server?tab=readme-ov-file#tool-configuration
	EnabledToolsets []string

	// ReadOnly indicates if we should only offer read-only tools
	ReadOnly bool
}
