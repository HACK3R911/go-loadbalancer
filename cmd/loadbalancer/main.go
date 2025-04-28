package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloudru/internal/balancer"
	"cloudru/internal/configs"
	"cloudru/internal/proxy"
	"cloudru/internal/ratelimit"
	"cloudru/internal/server"
)

func main() {
	cfg, err := configs.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация rate limiter
	rl := ratelimit.New(ratelimit.Config{
		Capacity:   cfg.RateLimit.Capacity,
		RefillRate: cfg.RateLimit.RefillRate,
	})

	lb := balancer.New(cfg.Backends)
	go lb.HealthCheck(10 * time.Second)

	px := proxy.New(lb)
	//
	handler := rl.Middleware(px)

	srv := new(server.Server)

	//gracefull shutdown
	go func() {
		if err := srv.Run(cfg.Port, handler); err != nil {
			log.Fatalf("error running server: %s", err.Error())
		}
	}()

	log.Println("Loadbalancer is started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Loadbalancer is stopped")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error on server shutdown: %s", err)
	}
}
