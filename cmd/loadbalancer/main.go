package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloudru/internal/balancer"
	"cloudru/internal/proxy"
	"cloudru/internal/server"

	"cloudru/internal/configs"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := configs.Load(configPath)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	lb := balancer.New(cfg.Backends)
	go lb.HealthCheck(10 * time.Second)

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

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error on server shutdown: %s", err)
	}
}
