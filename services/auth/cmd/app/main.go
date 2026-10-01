package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ritchieridanko/apotekly/services/auth/configs"
	"github.com/ritchieridanko/apotekly/services/auth/internal/di"
	"github.com/ritchieridanko/apotekly/services/auth/internal/infra"
	"github.com/ritchieridanko/apotekly/services/auth/internal/transport/rpc/server"
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
	diCtx, diCancel := context.WithTimeout(context.Background(), cfg.OAuth.InitTimeout)
	defer diCancel()

	ctr, err := di.Init(diCtx, cfg, inf)
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	// Server Start
	srv := ctr.Server()
	go func(srv *server.Server) {
		if err := srv.Start(); err != nil {
			log.Fatalln("[FATAL]:", err)
		}
	}(srv)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Printf("[%s] is shutting down...", strings.ToUpper(cfg.App.Name))

	sdCtx, sdCancel := context.WithTimeout(context.Background(), cfg.Server.Timeout.Shutdown)
	defer sdCancel()

	if err := srv.Shutdown(sdCtx); err != nil {
		log.Println("[WARN]:", err)
	}
}
