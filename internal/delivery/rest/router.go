package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo/pkg/logger"
)

func NewRouter(handler *RestHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 10))

	r.Get("/", handler.mainHandler)

	r.Get("/tasks", handler.getAllTasks)        // Получить Конкретную задачу
	r.Get("/tasks/{id}", handler.getTaskByID)   // Список всех задач
	r.Post("/tasks", handler.createTask)        // Создать задачу
	r.Put("/tasks/{id}", handler.putTask)       // Полностью обновить задачу
	r.Patch("/tasks/{id}", handler.patchTask)   // Частично обновить задачу
	r.Delete("/tasks/{id}", handler.deleteTask) // Удалить задачу

	return r
}

func StartServer(r *chi.Mux, addr string) {
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		IdleTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
	}

	go func() {
		logger.Info.Println("сервер слушает :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error.Fatalf("ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info.Println("начало graceful shutdown..")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error.Printf("ошибка graceful shutdown: %v", err)
	}

	logger.Info.Println("успешное завершение программы - graceful shutdown!")
}
