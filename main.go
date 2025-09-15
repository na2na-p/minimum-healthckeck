package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

func logJSON(level, message string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Fields:    fields,
	}
	data, _ := json.Marshal(entry)
	fmt.Fprintln(os.Stderr, string(data))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "OK")
	logJSON("info", "Served health check request", map[string]interface{}{
		"path":   r.URL.Path,
		"method": r.Method,
		"remote": r.RemoteAddr,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logJSON("info", "Starting server", map[string]interface{}{
			"port": port,
			"url":  fmt.Sprintf("http://localhost:%s/healthz", port),
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logJSON("error", "Could not start server", map[string]interface{}{
				"error": err.Error(),
			})
			os.Exit(1)
		}
	}()

	<-stop

	logJSON("info", "Shutting down server", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logJSON("error", "Server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	logJSON("info", "Server exiting gracefully", nil)
}
