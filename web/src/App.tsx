import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { TooltipProvider } from "@/components/ui/tooltip"
import { ClientConfigs } from "@/components/tabs/client-configs"
import { MCPInspector } from "@/components/tabs/mcp-inspector"
import {Settings, Wrench, MessageCircle, Github, ExternalLink, SquareCode} from "lucide-react"
import { Button } from "@/components/ui/button"

export default function App() {
  return (
    <TooltipProvider>
      <div className="min-h-screen bg-background">
        {/* Header */}
        <header className="border-b bg-background sticky top-0 z-50">
          <div className="container mx-auto px-4 h-14 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <img
                src="/logo.png"
                alt="VictoriaMetrics"
                className="h-8 w-8"
              />
              <div className="flex flex-col">
                <span className="font-semibold text-lg leading-tight">VictoriaMetrics MCP Server</span>
                <span className="text-xs text-muted-foreground leading-tight">Model Context Protocol Server for VictoriaMetrics API</span>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Button variant="ghost" size="sm" asChild>
                <a
                  href="https://github.com/VictoriaMetrics/mcp-victoriametrics/blob/main/README.md"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-1.5"
                >
                  Docs and source code
                  <ExternalLink className="h-3 w-3" />
                </a>
              </Button>
            </div>
          </div>
        </header>

        {/* Main Content */}
        <main className="container mx-auto px-4 py-6">
          <Tabs defaultValue="setup" className="space-y-6">
            <TabsList className="grid w-full max-w-md grid-cols-2">
              <TabsTrigger value="setup" className="flex items-center gap-2">
                <Settings className="h-4 w-4" />
                Setup
              </TabsTrigger>
              <TabsTrigger value="inspector" className="flex items-center gap-2">
                <Wrench className="h-4 w-4" />
                Inspector
              </TabsTrigger>
            </TabsList>

            <TabsContent value="setup">
              <ClientConfigs />
            </TabsContent>

            <TabsContent value="inspector">
              <MCPInspector />
            </TabsContent>
          </Tabs>
        </main>
      </div>
    </TooltipProvider>
  )
}
