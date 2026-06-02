// Command worker starts a Conductor task runner that registers each
// container/Kubernetes agent tool (podman, kind, k3s) as a worker, then polls
// the Conductor server for work until interrupted.
//
// Configuration is read from the environment:
//
//	CONDUCTOR_SERVER_URL   Conductor API base URL (default http://localhost:8080/api)
//	CONDUCTOR_AUTH_KEY     Orkes application key id    (optional; omit for unauthenticated OSS)
//	CONDUCTOR_AUTH_SECRET  Orkes application key secret (optional)
//	WORKER_BATCH_SIZE      tasks polled per cycle, per worker (default 1)
//	WORKER_POLL_INTERVAL   poll interval, e.g. 100ms, 1s (default 1s)
package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/settings"
	"github.com/conductor-sdk/conductor-go/sdk/worker"

	"github.com/AGenNext/Agent-Tools/internal/tools"
)

func main() {
	serverURL := getenv("CONDUCTOR_SERVER_URL", "http://localhost:8080/api")
	authKey := os.Getenv("CONDUCTOR_AUTH_KEY")
	authSecret := os.Getenv("CONDUCTOR_AUTH_SECRET")

	apiClient := client.NewAPIClient(
		settings.NewAuthenticationSettings(authKey, authSecret),
		settings.NewHttpSettings(serverURL),
	)

	batchSize := getenvInt("WORKER_BATCH_SIZE", 1)
	pollInterval := getenvDuration("WORKER_POLL_INTERVAL", time.Second)

	taskRunner := worker.NewTaskRunnerWithApiClient(apiClient)
	for taskName, tool := range tools.Registry {
		taskRunner.StartWorker(taskName, tool.Worker, batchSize, pollInterval)
		log.Printf("registered agent-tool worker %q v%s (capabilities=%v, batch=%d, poll=%s)",
			taskName, tool.Version, tool.Capabilities, batchSize, pollInterval)
	}

	log.Printf("polling %s for tasks; press Ctrl+C to stop", serverURL)

	// Block until a termination signal arrives, then stop polling cleanly.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down workers")
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getenvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
