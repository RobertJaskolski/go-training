package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

// TYPES -------------------------------

type Storage interface {
	CreateTask(*Task) error
	GetTasks() error
	GetTaskByID(int) error
	UpdateTask(*Task) error
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
		    created_at timestamp,
		    modified_at timestamp
		)`

	_, err := s.db.Exec(query)
	return err
}

// HANDLERS -------------------------------

func (*PostgresStore) CreateTask(*Task) error {
	return nil
}

func (*PostgresStore) GetTasks() error {
	return nil
}

func (*PostgresStore) GetTaskByID(int) error {
	return nil
}

func (*PostgresStore) UpdateTask(*Task) error {
	return nil
}

func (*PostgresStore) DeleteTask(int) error {
	return nil
}
