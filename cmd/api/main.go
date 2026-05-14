package main

import (
	"todo/internal/config"
	"todo/internal/delivery/rest"
	"todo/internal/repository/sqlite"
	"todo/internal/server"
	"todo/internal/service"
	"todo/pkg/logger"
)

func main() {
	/* Инициализация конфига */
	cfg := config.Init()

	/* Инициализация логгера*/
	logger.Init(cfg.AppMode == "dev")

	/* Репозиторий */
	repo, err := repository.New(cfg.DBPath)
	if err != nil {
		logger.Error.Fatalf("ошибка запуска репозитория: %v", err)
	}

	/* Сервис */
	serv := service.New(repo)

	/* REST API */
	restHandler := rest.NewHandler(serv)
	restRouter := rest.NewRouter(restHandler)

	restServer := server.NewRESTServer(restRouter,
		server.RESTServerConfig{
			Addr:         ":" + cfg.RestPort,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
			GSTime:       cfg.Shutdown})

	if err := restServer.Run(); err != nil {
		logger.Error.Printf("ошибка остановки сервера")
	}
}
