package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// TYPES -------------------------------

type Storage interface {
	CreateTask(*CreateTaskDTO) (*Task, error)
	GetTasks() ([]*Task, error)
	GetTaskByID(int) (*Task, error)
	UpdateTask(int, *UpdateTaskDTO) (*Task, error)
	DeleteTask(int) error
}

type PostgresStore struct {
	db *sql.DB
}

// METHODS -------------------------------

func NewPostgresStore() (*PostgresStore, error) {
	connectStr := "user=postgres dbname=todo password=root sslmode=disable"
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateTasksTable()
}

// HANDLERS -------------------------------

func (s *PostgresStore) CreateTasksTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
    		id serial primary key,
    		title varChar(50),
		    description varchar(356),
		    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		    modified_at timestamp DEFAULT CURRENT_TIMESTAMP
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateTask(dto *CreateTaskDTO) (*Task, error) {
	query := `
		INSERT INTO tasks 
		    (title, description) 
		VALUES 
		    ($1, $2)
		RETURNING id, title, description, created_at, modified_at
	`
	rows, err := s.db.Query(query, dto.Title, dto.Description)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		task := new(Task)
		if err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.ModifiedAt); err != nil {
			return nil, err
		}

		fmt.Println(task)

		return task, nil
	}

	return nil, nil
}

func (s *PostgresStore) GetTasks() ([]*Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks ORDER BY tasks.id")
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	for rows.Next() {
		task := new(Task)
		if err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.ModifiedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *PostgresStore) GetTaskByID(id int) (*Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE tasks.id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		task := new(Task)
		if err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.ModifiedAt); err != nil {
			return nil, err
		}

		return task, nil
	}

	return nil, fmt.Errorf("task %d not found", id)
}

func (s *PostgresStore) UpdateTask(id int, dto *UpdateTaskDTO) (*Task, error) {
	dtoValues := make(map[string]string)
	if dto.Title != "" {
		dtoValues["title"] = dto.Title
	}

	if dto.Description != "" {
		dtoValues["description"] = dto.Description
	}

	query := "UPDATE tasks SET modified_at = CURRENT_TIMESTAMP, "
	i := 0
	for key, value := range dtoValues {
		query += fmt.Sprintf("%v = '%v'", key, value)
		if i != len(dtoValues)-1 {
			query += ","
		}
		query += " "
		i++
	}
	query += "WHERE tasks.id = $1 RETURNING id, title, description, created_at, modified_at"
	rows, err := s.db.Query(query, id)
	for rows.Next() {
		task := new(Task)
		if err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.ModifiedAt); err != nil {
			return nil, err
		}

		return task, nil
	}

	return nil, nil
}

func (s *PostgresStore) DeleteTask(id int) error {
	rows, err := s.db.Query("DELETE FROM tasks WHERE tasks.id = $1 RETURNING id", id)
	for rows.Next() {
		task := new(Task)
		if err = rows.Scan(&task.ID); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("task %d not found", id)
}
