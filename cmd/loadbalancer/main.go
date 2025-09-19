package main

import (
	"cloudru/internal/configs"
	"cloudru/internal/logger"
	"cloudru/internal/server"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO viper config
	// TODO /healthcheck
	// TODO prometh метрики

	cfg, err := configs.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	//"health_check %b \n", cfg.HealthCheck.Enable,
	//"rate_limit %b", cfg.RateLimit.Enable,
	log.Infof(`Конфигурация загружена
		port %s,
		algorithm %s
		backends_count %d`,
		cfg.Port, cfg.Algorithm, len(cfg.Backends),
	)

	//lb := balancer.New(cfg.Backends)
	//go lb.HealthCheck(10 * time.Second)
	//
	//px := proxy.New(lb)
	srv := server.New(log)

	//Старт сервера
	go func() {
		if err := srv.Run(cfg.Port); err != nil {
			log.Fatalf("Ошибка запуска сервера: %s", err.Error())
		}
	}()

	//gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Ошибка при завершении работы сервера: %s", err)
	}

	log.Info("Сервер успешно остановлен")
}
