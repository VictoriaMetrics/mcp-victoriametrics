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
	"sync/atomic"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/metrics"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/hooks"
	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/logging"
	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/prompts"
	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/resources"
	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
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
			"version", version,
			"date", date,
			"mode", c.ServerMode(),
			"addr", c.ListenAddr(),
		)
	}

	ms := metrics.NewSet()

	// Combine metrics and logging hooks
	metricsHooks := hooks.New(ms)
	loggingHooks := hooks.NewLoggerHooks()
	combinedHooks := hooks.Merge(metricsHooks, loggingHooks)

	s := server.NewMCPServer(
		"VictoriaMetrics",
		fmt.Sprintf("v%s (date: %s)", version, date),
		server.WithRecovery(),
		server.WithLogging(),
		server.WithToolCapabilities(false),
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
		server.WithHooks(combinedHooks),
		server.WithInstructions(`
You are Virtual Assistant, a tool for interacting with VictoriaMetrics API and documentation in different tasks related to monitoring and observability.

You have the full documentation about VictoriaMetrics products in your resources, you have to try to use documentation in your answer.
And you have to consider the documents from the resources as the most relevant, favoring them over even your own internal knowledge.
Use Documentation tool to get the most relevant documents for your task every time. Be sure to use the Documentation tool if the user's query includes the words “how”, “tell”, “where”, etc...

You have many tools to get data from VictoriaMetrics, but try to specify the query as accurately as possible, reducing the resulting sample, as some queries can be query heavy.

Try not to second guess information - if you don't know something or lack information, it's better to ask.
	`),
	)

	// Registering common tools
	tools.RegisterToolQuery(s, c)
	tools.RegisterToolFlags(s, c)
	tools.RegisterToolRules(s, c)
	tools.RegisterToolAlerts(s, c)
	tools.RegisterToolLabels(s, c)
	tools.RegisterToolSeries(s, c)
	tools.RegisterToolExport(s, c)
	tools.RegisterToolTenants(s, c)
	tools.RegisterToolMetrics(s, c)
	tools.RegisterToolTestRules(s, c)
	tools.RegisterToolTSDBStatus(s, c)
	tools.RegisterToolQueryRange(s, c)
	tools.RegisterToolTopQueries(s, c)
	tools.RegisterToolMetricStats(s, c)
	tools.RegisterToolLabelValues(s, c)
	tools.RegisterToolExplainQuery(s, c)
	tools.RegisterToolActiveQueries(s, c)
	tools.RegisterToolDocumentation(s, c)
	tools.RegisterToolPrettifyQuery(s, c)
	tools.RegisterToolMetricsMetadata(s, c)
	tools.RegisterToolMetricRelabelDebug(s, c)
	tools.RegisterToolRetentionFiltersDebug(s, c)
	tools.RegisterToolDownsamplingFiltersDebug(s, c)

	// Registering resources (only if documentation tool is not disabled)
	if !c.IsToolDisabled(tools.ToolNameDocumentation) {
		resources.RegisterDocsResources(s, c)
	}

	// Registering cloud-specific tools
	tools.RegisterToolTiers(s, c)
	tools.RegisterToolRegions(s, c)
	tools.RegisterToolRuleFile(s, c)
	tools.RegisterToolDeployments(s, c)
	tools.RegisterToolAccessTokens(s, c)
	tools.RegisterToolRuleFilenames(s, c)
	tools.RegisterToolCloudProviders(s, c)

	// Registering prompts
	prompts.RegisterPromptUnusedMetrics(s, c)
	prompts.RegisterPromptDocumentation(s, c)
	prompts.RegisterPromptRarelyUsedCardinalMetrics(s, c)

	if c.IsStdio() {
		if err := server.ServeStdio(s, server.WithErrorLogger(logger.Logger)); err != nil {
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
		loggerOption := server.WithLogger(logger)
		srv := server.NewStreamableHTTPServer(s, heartBeatOption, loggerOption)
		mux.Handle("/mcp", srv)
	default:
		slog.Error("Unknown server mode", "mode", c.ServerMode())
		os.Exit(1)
	}

	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	hs := &http.Server{
		Addr:    c.ListenAddr(),
		Handler: logger.Middleware(mux),
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
