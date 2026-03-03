package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HACK3R911/go-loadbalancer/internal/balancer"
	"github.com/HACK3R911/go-loadbalancer/internal/configs"
	"github.com/HACK3R911/go-loadbalancer/internal/proxy"
	"github.com/HACK3R911/go-loadbalancer/internal/server"
)

func main() {
	cfg, err := configs.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	lb := balancer.New(cfg.Backends)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go lb.HealthCheck(ctx, 10*time.Second)

	px := proxy.New(lb)
	srv := new(server.Server)

	//Старт сервера
	go func() {
		if err := srv.Run(cfg.Port, px); err != nil {
			log.Fatalf("error running server: %s", err.Error())
		}
	}()

	log.Println("Loadbalancer is started ")

	//gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Loadbalancer is stopped")

	cancel()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error on server shutdown: %s", err)
	}
}
