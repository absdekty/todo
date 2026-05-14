package rest

import (
	"time"
)

type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"desc"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskCreate struct {
	Title       string `json:"title"`
	Description string `json:"desc"`
}

type TaskUpdateFull struct {
	Title       string `json:"title"`
	Description string `json:"desc"`
	Completed   bool   `json:"completed"`
}
