import { useState, useEffect, useCallback, useRef } from 'react'
import { Client } from '@modelcontextprotocol/sdk/client/index.js'
import { StreamableHTTPClientTransport } from '@modelcontextprotocol/sdk/client/streamableHttp.js'

export interface MCPTool {
  name: string
  description?: string
  inputSchema?: Record<string, unknown>
}

export interface MCPPrompt {
  name: string
  description?: string
  arguments?: Array<{
    name: string
    description?: string
    required?: boolean
  }>
}

export interface MCPServerInfo {
  name: string
  version: string
  protocolVersion?: string
  capabilities?: Record<string, unknown>
  instructions?: string
}

export interface UseMCPClientResult {
  client: Client | null
  serverInfo: MCPServerInfo | null
  tools: MCPTool[]
  prompts: MCPPrompt[]
  isConnected: boolean
  isConnecting: boolean
  error: Error | null
  connect: () => Promise<void>
  disconnect: () => Promise<void>
  callTool: (name: string, args: Record<string, unknown>) => Promise<void>
  getPrompt: (name: string, args: Record<string, string>) => Promise<void>
}

export function useMCPClient(): UseMCPClientResult {
  const [client, setClient] = useState<Client | null>(null)
  const [serverInfo, setServerInfo] = useState<MCPServerInfo | null>(null)
  const [tools, setTools] = useState<MCPTool[]>([])
  const [prompts, setPrompts] = useState<MCPPrompt[]>([])
  const [isConnected, setIsConnected] = useState(false)
  const [isConnecting, setIsConnecting] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  const clientRef = useRef<Client | null>(null)

  const getMCPEndpoint = useCallback(() => {
    const baseUrl = window.location.origin
    return `${baseUrl}/mcp`
  }, [])

  const connect = useCallback(async () => {
    if (isConnecting || isConnected) return

    setIsConnecting(true)
    setError(null)

    try {
      const endpoint = getMCPEndpoint()

      // Create client
      const mcpClient = new Client({
        name: 'mcp-victoriametrics-web',
        version: '1.0.0',
      }, {
        capabilities: {}
      })

      const transport = new StreamableHTTPClientTransport(new URL(endpoint))
      await mcpClient.connect(transport)

      clientRef.current = mcpClient
      setClient(mcpClient)

      // Get server info
      const info = mcpClient.getServerVersion()
      setServerInfo({
        name: info?.name || 'Unknown',
        version: info?.version || 'Unknown',
        protocolVersion: info?.title,
        capabilities: mcpClient.getServerCapabilities() as Record<string, unknown>,
      })

      // Fetch tools
      try {
        const toolsResult = await mcpClient.listTools()
        setTools(toolsResult.tools.map(t => ({
          name: t.name,
          description: t.description,
          inputSchema: t.inputSchema as Record<string, unknown>,
        })))
      } catch {
        console.warn('Failed to fetch tools')
      }

      // Fetch prompts
      try {
        const promptsResult = await mcpClient.listPrompts()
        setPrompts(promptsResult.prompts.map(p => ({
          name: p.name,
          description: p.description,
          arguments: p.arguments,
        })))
      } catch {
        console.warn('Failed to fetch prompts')
      }

      setIsConnected(true)
    } catch (err) {
      setError(err instanceof Error ? err : new Error(String(err)))
      setIsConnected(false)
    } finally {
      setIsConnecting(false)
    }
  }, [isConnecting, isConnected, getMCPEndpoint])

  const disconnect = useCallback(async () => {
    if (clientRef.current) {
      try {
        await clientRef.current.close()
      } catch {
        // Ignore close errors
      }
      clientRef.current = null
    }
    setClient(null)
    setServerInfo(null)
    setTools([])
    setPrompts([])
    setIsConnected(false)
    setError(null)
  }, [])

  const callTool = useCallback(async (name: string, args: Record<string, unknown>) => {
    if (!clientRef.current) {
      throw new Error('Not connected to MCP server')
    }
    const result = await clientRef.current.callTool({ name, arguments: args })
    return result
  }, [])

  const getPrompt = useCallback(async (name: string, args: Record<string, string>) => {
    if (!clientRef.current) {
      throw new Error('Not connected to MCP server')
    }
    const result = await clientRef.current.getPrompt({ name, arguments: args })
    return result
  }, [])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      disconnect()
    }
  }, [disconnect])

  return {
    client,
    serverInfo,
    tools,
    prompts,
    isConnected,
    isConnecting,
    error,
    connect,
    disconnect,
    callTool,
    getPrompt,
  }
}
