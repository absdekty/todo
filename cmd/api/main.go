package main

import (
	"todo/internal/config"
	"todo/internal/delivery/rest"
	"todo/internal/repository/sqlite"
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

	/* REST */
	restHandler := rest.NewHandler(serv)
	restRouter := rest.NewRouter(restHandler)
	rest.StartServer(restRouter, ":"+cfg.RestPort)
}
