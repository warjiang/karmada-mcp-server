# karmada-mcp-server
MCP Server for karmada



```yaml
# for stdio mode
{
  "mcpServers": {
    "karmada-mcp-server": {
      "name": "karmada-mcp-server",
      "type": "stdio",
      "command": "/path/to/karmada-mcp-server",
      "args": [
        "stdio",
        "--karmada-kubeconfig=/Users/warjiang/.kube/karmada.config",
        "--karmada-context=karmada-apiserver",
        "--skip-karmada-apiserver-tls-verify"
      ]
    }
  }
}


# for sse mode
{
  "mcpServers": {
    "karmada-mcp-server": {
      "name": "karmada-mcp-server",
      "type": "sse",
      "baseUrl": "http://localhost:1234/mcp/sse"
    }
  }
}

```