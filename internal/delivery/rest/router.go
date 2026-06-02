package rest

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *RestHandler) *chi.Mux {
	r := chi.NewRouter()

	setupMiddleware(r, handler.metrics)

	r.Post("/register", handler.RegisterUser) // POST /register - зарегистрировать пользователя
	r.Post("/login", handler.LoginUser)       // POST /login - аутентифицироваться, получить токен
	r.Post("/refresh", handler.RefreshTokens) // POST /refresh - аутентификация по рефреш токену
	r.Post("/logout", handler.LogoutUser)     // POST /refresh - аутентификация по рефреш токену

	r.Get("/", handler.mainHandler)

	r.Get("/metrics", handler.GetMetrics) // GET /metrics - получить актуальные метрики

	r.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware)

		r.Get("/", handler.getAllTasks) // GET /tasks - список всех задач
		r.Post("/", handler.createTask) // POST /tasks - создать задачу

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.getTaskByID)   // GET /tasks/{id} - получить задачу
			r.Put("/", handler.putTask)       // PUT /tasks/{id} - полностью обновить
			r.Patch("/", handler.patchTask)   // PATCH /tasks/{id} - частично обновить
			r.Delete("/", handler.deleteTask) // DELETE /tasks/{id} - удалить задачу
		})
	})

	return r
}
