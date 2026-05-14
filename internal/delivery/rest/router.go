package rest

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *RestHandler) *chi.Mux {
	r := chi.NewRouter()

	setupMiddleware(r)

	r.Get("/", handler.mainHandler)

	r.Route("/tasks", func(r chi.Router) {
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
