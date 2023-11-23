package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// TYPES -------------------------------

type Storage interface {
	CreateTask(*CreateTaskDTO) error
	GetTasks() ([]*Task, error)
	GetTaskByID(int) (*Task, error)
	UpdateTask(*UpdateTaskDTO) error
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

// HANDLERS -------------------------------

func (s *PostgresStore) CreateTask(dto *CreateTaskDTO) error {
	query := `
		INSERT INTO tasks 
		    (title, description) 
		VALUES 
		    ($1, $2)
	`
	_, err := s.db.Query(query, dto.Title, dto.Description)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetTasks() ([]*Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	for rows.Next() {
		task := new(Task)
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.ModifiedAt); err != nil {
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
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.ModifiedAt); err != nil {
			return nil, err
		}

		return task, nil
	}

	return nil, fmt.Errorf("task %d not found", id)
}

func (*PostgresStore) UpdateTask(dto *UpdateTaskDTO) error {
	return nil
}

func (*PostgresStore) DeleteTask(int) error {
	return nil
}
