package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/metrics"

	"github.com/VictoriaMetrics-Community/mcp-victoriatraces/cmd/mcp-victoriatraces/config"
	"github.com/VictoriaMetrics-Community/mcp-victoriatraces/cmd/mcp-victoriatraces/hooks"
	"github.com/VictoriaMetrics-Community/mcp-victoriatraces/cmd/mcp-victoriatraces/prompts"
	"github.com/VictoriaMetrics-Community/mcp-victoriatraces/cmd/mcp-victoriatraces/resources"
	"github.com/VictoriaMetrics-Community/mcp-victoriatraces/cmd/mcp-victoriatraces/tools"
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
		fmt.Printf("Error initializing config: %v\n", err)
		return
	}

	if !c.IsStdio() {
		log.Printf("Starting mcp-victoriatraces version %s (date: %s)", version, date)
	}

	ms := metrics.NewSet()
	s := server.NewMCPServer(
		"VictoriaTraces",
		fmt.Sprintf("v%s (date: %s)", version, date),
		server.WithRecovery(),
		server.WithLogging(),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithHooks(hooks.New(ms)),
		server.WithInstructions(`
You are Virtual Assistant, a tool for interacting with VictoriaTraces API and documentation in different tasks related to distributed tracing and observability.

You have the full documentation about VictoriaTraces in your resources, you have to try to use documentation in your answer.
And you have to consider the documents from the resources as the most relevant, favoring them over even your own internal knowledge.
Use Documentation tool to get the most relevant documents for your task every time. Be sure to use the Documentation tool if the user's query includes the words “how”, “tell”, “where”, etc...

You have many tools to get data from VictoriaTraces, but try to specify the query as accurately as possible, reducing the resulting sample, as some queries can be query heavy.

Try not to second guess information - if you don't know something or lack information, it's better to ask.
	`),
	)

	resources.RegisterDocsResources(s, c)

	tools.RegisterToolServices(s, c)
	tools.RegisterToolDocumentation(s, c)
	tools.RegisterToolTrace(s, c)

	prompts.RegisterPromptDocumentation(s, c)

	if c.IsStdio() {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("failed to start server in stdio mode on %s: %v", c.ListenAddr(), err)
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
		log.Printf("Starting server in SSE mode on %s", c.ListenAddr())
		srv := server.NewSSEServer(s)
		mux.Handle(srv.CompleteSsePath(), srv.SSEHandler())
		mux.Handle(srv.CompleteMessagePath(), srv.MessageHandler())
	case "http":
		log.Printf("Starting server in HTTP mode on %s", c.ListenAddr())
		heartBeatOption := server.WithHeartbeatInterval(c.HeartbeatInterval())
		srv := server.NewStreamableHTTPServer(s, heartBeatOption)
		mux.Handle("/mcp", srv)
	default:
		log.Fatalf("Unknown server mode: %s", c.ServerMode())
	}

	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	hs := &http.Server{
		Addr:    c.ListenAddr(),
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	listener, err := net.Listen("tcp", c.ListenAddr())
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", c.ListenAddr(), err)
	}
	log.Printf("Server is listening on %s", c.ListenAddr())

	go func() {
		if err := hs.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	isReady.Store(true)
	<-rootCtx.Done()
	stop()
	isReady.Store(false)
	log.Println("Received shutdown signal, shutting down.")

	// Give time for readiness check to propagate
	time.Sleep(_readinessDrainDelay)
	log.Println("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()
	err = hs.Shutdown(shutdownCtx)
	stopOngoingGracefully()
	if err != nil {
		log.Println("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}

	log.Println("Server stopped.")
}
