package main

import (
	"log"
	"time"

	"cloudru/internal/balancer"
	"cloudru/internal/configs"
	"cloudru/internal/proxy"
	"cloudru/internal/server"
)

func main() {
	cfg, err := configs.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	lb := balancer.New(cfg.Backends)
	go lb.HealthCheck(10 * time.Second)

	px := proxy.New(lb)
	srv := server.New(px, cfg.Port)

	if err := srv.Start(); err != nil {
		log.Fatalf("Server starting failed: %v", err)
	}
	log.Printf("Starting server")
}
