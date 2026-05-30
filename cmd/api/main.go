package main

import (
	"todo/internal/config"
	"todo/internal/delivery/rest"
	"todo/internal/repository/sqlite"
	"todo/internal/server"
	"todo/internal/service"
	"todo/pkg/hasher"
	"todo/pkg/logger"
)

func main() {
	/* Инициализация конфига */
	cfg := config.Init()

	/* Инициализация логгера*/
	logger.Init(cfg.AppMode == "dev")

	/* Репозиторий */
	db, err := sqlite.New(cfg.DBPath)
	if err != nil {
		logger.Error.Fatalf("ошибка запуска репозитория: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error.Printf("ошибка закрытия БД: %v", err)
		}
	}()

	taskRepo := sqlite.NewTask(db)
	userRepo := sqlite.NewUser(db)

	/* Hasher */
	hasher := hasher.New()

	/* Сервис */
	taskService := service.NewTask(taskRepo, userRepo)
	userService := service.NewUser(userRepo, hasher)
	jwtService := service.NewJWT(
		service.JWTConfig{
			JWTSecret:     cfg.JWTSecret,
			JWTExpiration: cfg.JWTExpiration,
		})

	/* REST API */
	restHandler := rest.NewHandler(taskService, userService, jwtService)
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
