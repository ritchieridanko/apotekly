package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ritchieridanko/apotekly/services/notification/configs"
	"github.com/ritchieridanko/apotekly/services/notification/internal/di"
	"github.com/ritchieridanko/apotekly/services/notification/internal/infra"
)

const (
	maxRetries      int           = 10
	maxRetryBackoff time.Duration = 30 * time.Second
)

func main() {
	// Config Initialization
	cfg, err := configs.Init("./configs")
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	// Infra Initialization
	inf, err := infra.Init(cfg)
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}
	defer func(inf *infra.Infra) {
		if err := inf.Close(); err != nil {
			log.Println("[WARN]:", err)
		}
	}(inf)

	// DI Initialization
	ctr, err := di.Init(cfg, inf)
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}
	defer func(ctr *di.Container) {
		if err := ctr.Close(); err != nil {
			log.Println("[WARN]:", err)
		}
	}(ctr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Event Subscribers & Processors Run
	workers := []func(context.Context) error{
		ctr.RunSubscriberAC,
		ctr.RunEventProcessor1,
		ctr.RunEventProcessor2,
	}
	for _, worker := range workers {
		runWorker(&wg, worker, ctx)
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Printf("[%s] is shutting down...", strings.ToUpper(cfg.App.Name))

	cancel()
	wg.Wait()
}

func runWorker(wg *sync.WaitGroup, fn func(context.Context) error, ctx context.Context) {
	wg.Add(1)

	go func(ctx context.Context, fn func(context.Context) error) {
		defer wg.Done()

		for attempt := 0; attempt < maxRetries; attempt++ {
			if err := fn(ctx); err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					log.Println("[INFO]: Worker stopped gracefully")
					return
				}

				// Exponential Backoff
				backoff := min(
					time.Second*time.Duration(1<<attempt),
					maxRetryBackoff,
				)
				log.Printf("[WARN]: Worker failed (attempt %d/%d): %v. Retrying in %v", attempt+1, maxRetries, err, backoff)

				select {
				case <-time.After(backoff):
					continue
				case <-ctx.Done():
					return
				}
			}
			return
		}

		log.Printf("[ERROR]: Worker exceeded max retries (%d). Giving up.", maxRetries)
	}(ctx, fn)
}
