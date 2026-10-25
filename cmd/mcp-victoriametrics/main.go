package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"

	metrics "github.com/VictoriaMetrics/metrics"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/auth"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/hooks"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/logging"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/registry"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/resources"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

var (
	version = "dev"
	date    = "unknown"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 3 * time.Second
)

func main() {
	c, err := config.InitConfig()
	if err != nil {
		log.Fatalf("FATAL: Error initializing config: %v\n", err)
		return
	}

	logger, err := logging.New(c)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize logger: %v\n", err)
		return
	}

	if !c.IsStdio() {
		slog.Info("Starting mcp-victoriametrics",
			"backends", backendNames(c),
			"version", version,
			"date", date,
			"mode", c.ServerMode(),
			"addr", c.ListenAddr(),
		)
	}

	for oldName, newName := range c.LegacyToolNames() {
		slog.Warn("MCP_DISABLED_TOOLS contains deprecated tool name, use the new name instead", "deprecated", oldName, "new", newName)
	}

	ms := metrics.NewSet()

	// Combine metrics and logging hooks
	metricsHooks := hooks.New(ms)
	loggingHooks := hooks.NewLoggerHooks()
	combinedHooks := hooks.Merge(metricsHooks, loggingHooks)

	s := server.NewMCPServer(
		c.ProductName(),
		fmt.Sprintf("v%s (date: %s)", version, date),
		server.WithRecovery(),
		server.WithLogging(),
		server.WithToolCapabilities(c.HasDynamicToolsets()),
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
		server.WithHooks(combinedHooks),
		server.WithInstructions(buildInstructions(c)),
	)

	registry.Register(s, c)
	if !c.IsToolDisabled(tools.ToolNameDocumentation) {
		resources.Register(s, c)
	}

	if c.IsStdio() {
		if err := server.ServeStdio(s, server.WithErrorLogger(logger.StdLogger())); err != nil {
			slog.Error("failed to start server in stdio mode", "addr", c.ListenAddr(), "error", err)
			os.Exit(1)
		}
		return
	}

	var isReady atomic.Bool

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		ms.WritePrometheus(w)
		metrics.WriteProcessMetrics(w)
	})
	mux.HandleFunc("/health/liveness", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write([]byte("OK\n"))
	})
	mux.HandleFunc("/health/readiness", func(w http.ResponseWriter, _ *http.Request) {
		if !isReady.Load() {
			http.Error(w, "Not ready", http.StatusServiceUnavailable)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write([]byte("Ready\n"))
	})

	switch c.ServerMode() {
	case "sse":
		slog.Info("Starting server in SSE mode", "addr", c.ListenAddr())
		srv := server.NewSSEServer(s)
		mux.Handle(srv.CompleteSsePath(), srv.SSEHandler())
		mux.Handle(srv.CompleteMessagePath(), srv.MessageHandler())
	case "http":
		slog.Info("Starting server in HTTP mode", "addr", c.ListenAddr())
		heartBeatOption := server.WithHeartbeatInterval(c.HeartbeatInterval())
		loggerOption := server.WithStreamableHTTPLogger(logger.SlogLogger())
		srv := server.NewStreamableHTTPServer(s, heartBeatOption, loggerOption)
		mux.Handle("/mcp", srv)
		mux.Handle("/", spaHandler())
	default:
		slog.Error("Unknown server mode", "mode", c.ServerMode())
		os.Exit(1)
	}

	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())

	var handler http.Handler = mux
	if authToken := c.AuthToken(); authToken != "" {
		slog.Info("Authentication enabled", "token_configured", true)
		handler = auth.Middleware(authToken)(mux)
	}
	handler = logger.Middleware(handler)
	hs := &http.Server{
		Addr:    c.ListenAddr(),
		Handler: handler,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	listener, err := net.Listen("tcp", c.ListenAddr())
	if err != nil {
		slog.Error("Failed to listen", "addr", c.ListenAddr(), "error", err)
		os.Exit(1)
	}
	slog.Info("Server is listening", "addr", c.ListenAddr())

	go func() {
		if err := hs.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	isReady.Store(true)
	<-rootCtx.Done()
	stop()
	isReady.Store(false)
	slog.Info("Received shutdown signal, shutting down.")

	// Give time for readiness check to propagate
	time.Sleep(_readinessDrainDelay)
	slog.Info("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()
	err = hs.Shutdown(shutdownCtx)
	stopOngoingGracefully()
	if err != nil {
		slog.Warn("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}

	slog.Info("Server stopped.")
}

var instructions = map[config.Backend]string{
	config.BackendMetrics: `
You are Virtual Assistant, a tool for interacting with VictoriaMetrics API and documentation in different tasks related to monitoring and observability.

You have the full documentation about VictoriaMetrics products in your resources, you have to try to use documentation in your answer.
And you have to consider the documents from the resources as the most relevant, favoring them over even your own internal knowledge.
Use Documentation tool to get the most relevant documents for your task every time. Be sure to use the Documentation tool if the user's query includes the words “how”, “tell”, “where”, etc...

You have many tools to get data from VictoriaMetrics, but try to specify the query as accurately as possible, reducing the resulting sample, as some queries can be query heavy.

Try not to second guess information - if you don't know something or lack information, it's better to ask.
	`,
	config.BackendLogs: `
You are Virtual Assistant, a tool for interacting with VictoriaLogs API and documentation in different tasks related to logs and observability.
You use LogsQL language to query logs and get information from the logs stored in VictoriaLogs.

You have the full documentation about VictoriaLogs in your resources, you have to try to use documentation in your answer.
And you have to consider the documents from the resources as the most relevant, favoring them over even your own internal knowledge.
Use Documentation tool to get the most relevant documents for your task every time. Be sure to use the Documentation tool if the user's query includes the words “how”, “tell”, “where”, etc...

You have many tools to get data from VictoriaLogs, but try to specify the query as accurately as possible, reducing the resulting sample, as some queries can be query heavy.

Try not to second guess information - if you don't know something or lack information, it's better to ask.
	`,
	config.BackendTraces: `
You are Virtual Assistant, a tool for interacting with VictoriaTraces API and documentation in different tasks related to distributed tracing and observability.

You have the full documentation about VictoriaTraces in your resources, you have to try to use documentation in your answer.
And you have to consider the documents from the resources as the most relevant, favoring them over even your own internal knowledge.
Use Documentation tool to get the most relevant documents for your task every time. Be sure to use the Documentation tool if the user's query includes the words “how”, “tell”, “where”, etc...

You have many tools to get data from VictoriaTraces, but try to specify the query as accurately as possible, reducing the resulting sample, as some queries can be query heavy.

Try not to second guess information - if you don't know something or lack information, it's better to ask.
	`,
}

const multiBackendInstructions = `
You are Virtual Assistant, a tool for interacting with %[1]s APIs and documentation in different tasks related to monitoring, logs, distributed tracing and observability.

You have the full documentation about %[1]s in your resources, you have to try to use documentation in your answer.
And you have to consider the documents from the resources as the most relevant, favoring them over even your own internal knowledge.
Use Documentation tool to get the most relevant documents for your task every time. Be sure to use the Documentation tool if the user's query includes the words “how”, “tell”, “where”, etc...

The following backends are configured:
%[2]s
%[3]s
You use MetricsQL language to query VictoriaMetrics backends and LogsQL language to query VictoriaLogs backends.

You have many tools to get data from these backends, but try to specify the query as accurately as possible, reducing the resulting sample, as some queries can be query heavy.

Try not to second guess information - if you don't know something or lack information, it's better to ask.
	`

func buildInstructions(c *config.Config) string {
	backends := c.Backends()
	if len(backends) == 1 {
		return instructions[c.Backend()]
	}
	var (
		list         strings.Builder
		hasStatic    bool
		cloudName    string
		cloudTypes   []string
		cloudIsNamed bool
	)
	for _, b := range backends {
		if b.IsCloud() {
			cloudName = b.Name()
			cloudIsNamed = b.Name() != string(b.Backend())
			cloudTypes = append(cloudTypes, b.Backend().ProductName())
			continue
		}
		hasStatic = true
		fmt.Fprintf(&list, "- %s: %s", b.Name(), b.Backend().ProductName())
		if b.Backend() == config.BackendMetrics && b.InstanceType() != "" {
			fmt.Fprintf(&list, " (%s)", b.InstanceType())
		}
		list.WriteString("\n")
	}
	if len(cloudTypes) > 0 {
		if cloudIsNamed {
			fmt.Fprintf(&list, "- %s: ", cloudName)
		} else {
			list.WriteString("- ")
		}
		fmt.Fprintf(&list, "VictoriaMetrics Cloud deployments of %s\n", strings.Join(cloudTypes, ", "))
	}

	var routing []string
	if hasStatic {
		routing = append(routing, `Tools available for several backends have a required "backend" parameter with the name of the backend to send the request to.`)
	}
	if len(cloudTypes) > 0 {
		sentence := `Requests to VictoriaMetrics Cloud require the "deployment_id" parameter and are sent to the backend that matches the type of the deployment.`
		if hasStatic {
			sentence = fmt.Sprintf(`For the %q backend, also pass the "deployment_id" parameter: the request is sent to the backend that matches the type of the deployment.`, cloudName)
		}
		routing = append(routing, sentence)
	}
	return fmt.Sprintf(multiBackendInstructions, c.ProductName(), list.String(), strings.Join(routing, "\n"))
}

func backendNames(c *config.Config) []string {
	var names []string
	for _, b := range c.Backends() {
		names = append(names, fmt.Sprintf("%s:%s", b.Name(), b.Backend()))
	}
	return names
}
