package main

import "time"

type APIError struct {
	Error string `json:"error"`
}

type APIESuccess struct {
	Success bool `json:"success"`
}

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
}

type CreateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
